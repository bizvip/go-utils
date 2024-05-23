/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 * 使用repo作为包名，可以避免与其他repository模式的代码冲突包名字（虽然冲突可以使用别名） *
 ******************************************************************************/

package repo

import (
	"errors"
	"math"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/bizvip/go-utils/db/mysql"
	"github.com/bizvip/go-utils/logs"
)

type BaseRepo[T any] struct{ Orm *gorm.DB }

type IBaseRepo[T any] interface {
	Exec(sql string, values ...interface{}) (*gorm.DB, error)
	CountAll() (int64, error)
	Insert(model *T) error
	UpdateById(id uint64, updateValues *T) error
	DeleteById(id uint64) error
	DeleteBy(condition map[string]interface{}, hardDelete bool) error
	FindByID(id uint64) (*T, error)
	FindBy(condition map[string]interface{}) (*T, error)
	SelectBy(condition map[string]interface{}, results *[]*T, opts ...SelOpt) error
	InsertOrIgnore(model *T, condition map[string]interface{}) (int64, error)
	InsertOrUpdate(insertItem *T, condition map[string]interface{}, updateValues *T) error
	UpsertByID(model *T, updateFields []string) error
	Upsert(model *T, condition map[string]interface{}) error
	GetByPage(page int, pageSize int) ([]*T, error)
	UpdateBy(condition map[string]interface{}, updateValues *T) (int64, error)
}

// NewBaseRepo 暂未使用接口返回
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
func (r *BaseRepo[T]) UpdateById(id uint64, updateValues *T) error {
	result := r.Orm.Model(updateValues).Where("id = ?", id).Updates(updateValues)
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

// FindByID 按照ID读取一条记录 本repo的find均代表找1条记录 select均代表多条记录
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

// Pagination 定义
type Pagination[T any] struct {
	TotalRecords int64 `json:"totalRecords"`
	TotalPages   int   `json:"totalPages"`
	CurrentPage  uint  `json:"currentPage"`
	PageSize     uint  `json:"pageSize"`
	Records      []*T  `json:"records"`
}

// GetByPage 根据分页获取记录
func (r *BaseRepo[T]) GetByPage(pageNum uint, pageSize uint, opts ...SelOpt) (*Pagination[T], error) {
	if pageNum < 1 || pageSize < 1 {
		return nil, errors.New("pageNum 和 pageSize 必须大于 0")
	}

	var results []*T
	var totalRecords int64

	// 创建基础查询
	query := r.Orm.Model(new(T))

	// 应用Where条件
	for _, opt := range opts {
		query = opt(query)
	}

	// 计算总记录数
	err := query.Count(&totalRecords).Error
	if err != nil {
		return nil, err
	}

	// 计算分页偏移量
	offset := (pageNum - 1) * pageSize

	// 应用分页和排序
	query = query.Offset(int(offset)).Limit(int(pageSize))
	for _, opt := range opts {
		query = opt(query)
	}

	// 执行查询
	err = query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationResult := &Pagination[T]{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  pageNum,
		PageSize:     pageSize,
		Records:      results,
	}

	return paginationResult, nil
}

// InsertOrIgnore 无显式事务 先查找，存在则忽略，否则插入 (并发性也可以由数据库相同的unique key来保证)
func (r *BaseRepo[T]) InsertOrIgnore(model *T, condition map[string]interface{}) (int64, error) {
	var existsT T
	err := r.Orm.Where(condition).First(&existsT).Error
	if err == nil {
		return 0, nil
	}
	if !errors.Is(gorm.ErrRecordNotFound, err) {
		return 0, err
	}
	// 使用 ON CONFLICT DO NOTHING 确保并发安全
	result := r.Orm.Clauses(clause.OnConflict{DoNothing: true}).Create(model)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// InsertOrUpdate 事务版本 先查找，不存在则插入，存在则更新
func (r *BaseRepo[T]) InsertOrUpdate(insertItem *T, condition map[string]interface{}, updateValues *T) error {
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
			if err = tx.Create(insertItem).Error; err != nil {
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

// UpsertByID 非显式事务(onConflict和clauses)，固定根据id查找记录，如果存在则更新，如果不存在则创建
func (r *BaseRepo[T]) UpsertByID(model *T, updateFields []string) error {
	result := r.Orm.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(updateFields)}).Create(model)
	return result.Error
}

// Upsert 非显式事务(onConflict和clauses)，根据condition查找记录，如果存在则更新，如果不存在则创建
func (r *BaseRepo[T]) Upsert(model *T, condition map[string]interface{}) error {
	// 内联函数获取结构体字段名
	getStructFields := func(v interface{}) []string {
		val := reflect.ValueOf(v).Elem()
		typ := val.Type()
		fields := make([]string, val.NumField())

		for i := 0; i < val.NumField(); i++ {
			fields[i] = typ.Field(i).Name
		}
		return fields
	}

	// 内联函数获取map的键名
	getMapKeys := func(m map[string]interface{}) []string {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		return keys
	}

	// 内联函数将字段名转换为clause.Column类型
	getColumnClauses := func(fields []string) []clause.Column {
		columns := make([]clause.Column, len(fields))
		for i, field := range fields {
			columns[i] = clause.Column{Name: field}
		}
		return columns
	}

	// 获取条件字段名和更新字段名
	conditionFields := getMapKeys(condition)
	updateFields := getStructFields(model)

	// 创建或更新记录
	tx := r.Orm.Clauses(clause.OnConflict{Columns: getColumnClauses(conditionFields),
		DoUpdates: clause.AssignmentColumns(updateFields)}).Create(model)

	return tx.Error
}

// UpdateBy 根据条件更新记录
func (r *BaseRepo[T]) UpdateBy(condition map[string]interface{}, updateValues *T) (int64, error) {
	// 尝试更新记录
	result := r.Orm.Model(updateValues).Where(condition).Updates(updateValues)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
