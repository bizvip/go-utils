package repository

import (
	"errors"
	"reflect"

	"gorm.io/gorm"

	"github.com/bizvip/go-utils/db/mysql"
)

type BaseRepo[T any] struct{ Orm *gorm.DB }
type IBaseRepo interface {
	Insert(model interface{}) error
	InsertOrUpdate(model interface{}, condition map[string]interface{}, forUpdateValues interface{}) error
	UpdateById(model interface{}, id uint64) error
	DeleteById(model interface{}, id uint64) error
	DeleteBy(condition map[string]interface{}, model interface{}, hardDelete bool) error
	SelectById(model interface{}, id uint64) error
	SelectBy(condition map[string]interface{}) ([]interface{}, error)
	SelectOne(condition map[string]interface{}, model interface{}) error
}

func NewBaseRepo[T any]() *BaseRepo[T] {
	orm := mysql.GetOrmInstance()
	if orm == nil {
		panic("base repo error : mysql orm instance is nil")
	}
	return &BaseRepo[T]{orm}
}
func (r *BaseRepo[T]) Exec(sql string, values ...interface{}) (*gorm.DB, error) {
	var result *gorm.DB
	if len(values) == 0 {
		result = r.Orm.Exec(sql)
	} else {
		result = r.Orm.Exec(sql, values)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return result, nil
}
func (r *BaseRepo[T]) CountAll(model interface{}) (int64, error) {
	var count int64
	result := r.Orm.Model(model).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
func (r *BaseRepo[T]) Insert(model interface{}) error {
	result := r.Orm.Create(model)
	return result.Error
}
func (r *BaseRepo[T]) UpdateById(model interface{}, id uint64) error {
	result := r.Orm.Model(model).Where("id = ?", id).Updates(model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no record found to update")
	}
	return nil
}
func (r *BaseRepo[T]) DeleteById(model interface{}, id uint64) error {
	result := r.Orm.Where("id = ?", id).Delete(model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no record found to delete")
	}
	return nil
}
func (r *BaseRepo[T]) DeleteBy(condition map[string]interface{}, model interface{}, hardDelete bool) error {
	var result *gorm.DB
	if hardDelete {
		result = r.Orm.Unscoped().Where(condition).Delete(model)
	} else {
		result = r.Orm.Where(condition).Delete(model)
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *BaseRepo[T]) SelectById(model interface{}, id uint64) error {
	result := r.Orm.First(model, "id = ?", id)
	return result.Error
}

func (r *BaseRepo[T]) InsertOrIgnore(model interface{}, condition map[string]interface{}) (int64, error) {
	tx := r.Orm.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	// 尝试查找符合条件的记录
	var count int64
	err := tx.Model(model).Where(condition).Count(&count).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if count > 0 {
		// 如果记录已存在，则提交事务并返回 0，无错误
		if err = tx.Commit().Error; err != nil {
			return 0, err
		}
		return 0, nil
	}

	// 如果未找到记录，则创建新记录
	result := tx.Create(model)
	if result.Error != nil {
		tx.Rollback()
		return 0, result.Error
	}

	if err = tx.Commit().Error; err != nil {
		return 0, err
	}

	return result.RowsAffected, nil
}

func (r *BaseRepo[T]) InsertOrUpdate(model interface{}, condition map[string]interface{}, updateValues interface{}) error {
	tx := r.Orm.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 使用反射来获取model的类型
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// 创建一个新的model类型的实例
	resultPtr := reflect.New(modelType).Interface()

	// 先尝试查找
	err := tx.Where(condition).First(resultPtr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没找到，创建新记录
			if err = tx.Create(model).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return err
		}
	} else {
		// 找到了，执行更新
		if err = tx.Model(resultPtr).Where(condition).Updates(updateValues).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// FindBy 查找多条
func (r *BaseRepo[T]) FindBy(condition map[string]interface{}, results *[]*T, limit int) error {
	orderBy, hasOrderBy := condition["orderBy"].(string)
	if hasOrderBy {
		delete(condition, "orderBy")
	}
	query := r.Orm.Where(condition)
	if hasOrderBy {
		query = query.Order(orderBy)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(results).Error
	if err != nil {
		return err
	}
	return nil
}

// FindOne 查找1条
func (r *BaseRepo[T]) FindOne(condition map[string]interface{}, model interface{}) error {
	err := r.Orm.Where(condition).First(model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("no record found")
		}
		return err
	}
	return nil
}
func (r *BaseRepo[T]) GetByPage(model interface{}, page int, pageSize int) (interface{}, error) {
	if page < 1 || pageSize < 1 {
		return nil, errors.New("page 和 pageSize 必须大于 0")
	}

	// 获取 model 的实际类型，处理指针类型的情况
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem() // 获取指针指向的元素类型
	}

	// 创建元素类型为 modelType 的切片的指针
	valuesPtr := reflect.New(reflect.SliceOf(modelType))
	values := valuesPtr.Elem() // 获取切片的实际值

	// 计算跳过的记录数
	offset := (page - 1) * pageSize
	result := r.Orm.Offset(offset).Limit(pageSize).Find(values.Addr().Interface())
	if result.Error != nil {
		return nil, result.Error
	}

	return values.Interface(), nil
}
