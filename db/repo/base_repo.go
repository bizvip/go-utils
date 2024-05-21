/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 * 使用repo作为包名，可以避免与其他repository模式的代码冲突包名字（虽然冲突可以使用别名） *
 ******************************************************************************/

package repo

import (
	"errors"
	"reflect"

	"gorm.io/gorm"

	"github.com/bizvip/go-utils/db/mysql"
	"github.com/bizvip/go-utils/logs"
)

type BaseRepo[T any] struct{ Orm *gorm.DB }

// type IBaseRepo interface {
// 	Insert(model interface{}) error
// 	InsertOrUpdate(model interface{}, condition map[string]interface{}, forUpdateValues interface{}) error
// 	UpdateById(model interface{}, id uint64) error
// 	DeleteById(model interface{}, id uint64) error
// 	DeleteBy(condition map[string]interface{}, model interface{}, hardDelete bool) error
// 	SelectById(model interface{}, id uint64) error
// 	SelectBy(condition map[string]interface{}) ([]interface{}, error)
// 	SelectOne(condition map[string]interface{}, model interface{}) error
// }

func NewBaseRepo[T any]() *BaseRepo[T] {
	orm := mysql.GetOrmInstance()
	if orm == nil {
		logs.Logger().Error("base repo error : mysql orm instance is nil")
		return nil
	}
	return &BaseRepo[T]{orm}
}

// Exec 执行原生sql语句
func (r *BaseRepo[T]) Exec(sql string, values ...interface{}) (*gorm.DB, error) {
	var result *gorm.DB
	if len(values) == 0 {
		result = r.Orm.Exec(sql)
	} else {
		result = r.Orm.Exec(sql, values...)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return result, nil
}

// CountAll 指定表无条件统计全部数量
func (r *BaseRepo[T]) CountAll() (int64, error) {
	var count int64
	var model T
	result := r.Orm.Model(&model).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// Insert 插入一条记录
func (r *BaseRepo[T]) Insert(model *T) error {
	result := r.Orm.Create(model)
	return result.Error
}

// UpdateById 按照ID更新一条
func (r *BaseRepo[T]) UpdateById(id uint64, updateValues map[string]interface{}) error {
	var model T
	result := r.Orm.Model(&model).Where("id = ?", id).Updates(updateValues)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no record found to update")
	}
	return nil
}

// DeleteById 按照ID删除一条
func (r *BaseRepo[T]) DeleteById(id uint64) error {
	var model T
	result := r.Orm.Where("id = ?", id).Delete(&model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no record found to delete")
	}
	return nil
}

// DeleteBy 根据给定条件删除记录，可选是否硬删除（仅对于有软删除的表）
func (r *BaseRepo[T]) DeleteBy(condition map[string]interface{}, hardDelete bool) error {
	var model T
	var result *gorm.DB
	if hardDelete {
		result = r.Orm.Unscoped().Where(condition).Delete(&model)
	} else {
		result = r.Orm.Where(condition).Delete(&model)
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindByID 按照ID读取一条记录
func (r *BaseRepo[T]) FindByID(id uint64) (*T, error) {
	var model T
	result := r.Orm.First(&model, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &model, nil
}

// FindBy 根据条件查找一条记录
func (r *BaseRepo[T]) FindBy(condition map[string]interface{}) (*T, error) {
	var model T
	result := r.Orm.Where(condition).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return &model, nil
}

// SelOpt SelectBy 条件配置
type SelOpt func(*gorm.DB) *gorm.DB

func WithOrderBy(orderBy string) SelOpt {
	return func(q *gorm.DB) *gorm.DB { return q.Order(orderBy) }
}
func WithLimit(limit int) SelOpt {
	return func(q *gorm.DB) *gorm.DB { return q.Limit(limit) }
}

// SelectBy 按照条件查找多条 可使用链式方法添加order和limit等参数
func (r *BaseRepo[T]) SelectBy(condition map[string]interface{}, results *[]*T, opts ...SelOpt) error {
	query := r.Orm.Where(condition)
	for _, opt := range opts {
		query = opt(query)
	}
	err := query.Find(results).Error
	if err != nil {
		return err
	}
	return nil
}

// InsertOrIgnore 事务版本  先查找，存在则忽略，否则插入
func (r *BaseRepo[T]) InsertOrIgnore(model *T, condition map[string]interface{}) (int64, error) {
	tx := r.Orm.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	var count int64
	err := tx.Model(new(T)).Where(condition).Count(&count).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if count > 0 {
		if err = tx.Commit().Error; err != nil {
			return 0, err
		}
		return 0, nil
	}
	result := tx.Create(&model)
	if result.Error != nil {
		tx.Rollback()
		return 0, result.Error
	}
	if err = tx.Commit().Error; err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}

// InsertOrUpdate 事务版本 先查找，不存在则插入，插入则更新
func (r *BaseRepo[T]) InsertOrUpdate(model *T, condition map[string]interface{}, updateValues map[string]interface{}) error {
	tx := r.Orm.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var item T
	err := tx.Where(condition).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Create(&model).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return err
		}
	} else {
		if err = tx.Model(&item).Updates(updateValues).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// Upsert 同上方法，无事务，根据条件查找记录，如果存在则更新，如果不存在则创建
func (r *BaseRepo[T]) Upsert(model *T, condition map[string]interface{}) error {
	var existingRecord T
	tx := r.Orm.Where(condition).First(&existingRecord)
	if tx.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, tx.Error) {
			return r.Orm.Create(model).Error
		}
		return tx.Error
	}
	// 如果记录存在，更新记录
	return r.Orm.Model(&existingRecord).Updates(model).Error
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

// UpdateBy 根据条件更新记录
func (r *BaseRepo[T]) UpdateBy(condition map[string]interface{}, updateValues map[string]interface{}) (int64, error) {
	var model T
	result := r.Orm.Model(&model).Where(condition).Updates(updateValues)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, nil
	}
	return result.RowsAffected, nil
}
