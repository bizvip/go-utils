package repo

import (
	"errors"
	"math"
)

// Pagination 定义
type Pagination[T any] struct {
	TotalRecords int64  `json:"totalRecords"`
	TotalPages   int64  `json:"totalPages"`
	CurrentPage  uint32 `json:"currentPage"`
	PageSize     uint32 `json:"pageSize"`
	Records      []*T   `json:"records"`
}

// GetByPage 根据分页获取记录
func (r *BaseRepo[T]) GetByPage(pageNum uint32, pageSize uint32, opts ...SelOpt) (*Pagination[T], error) {
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
	totalPages := int64(math.Ceil(float64(totalRecords) / float64(pageSize)))

	paginationResult := &Pagination[T]{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  pageNum,
		PageSize:     pageSize,
		Records:      results,
	}

	return paginationResult, nil
}
