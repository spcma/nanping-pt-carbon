package db

import (
	"app/internal/shared/entity"
	"fmt"

	"gorm.io/gorm"
)

// PaginationHelper 分页查询助手
type PaginationHelper[T any] struct {
	db *gorm.DB
}

// NewPaginationHelper 创建分页助手实例
func NewPaginationHelper[T any](db *gorm.DB) *PaginationHelper[T] {
	return &PaginationHelper[T]{db: db}
}

// PageQuery 执行分页查询
// T: 实体类型
// queryBuilder: 查询构建函数，用于添加where条件等
func (h *PaginationHelper[T]) PageQuery(pageNum, pageSize int, queryBuilder func(*gorm.DB) *gorm.DB) (*entity.PaginationResult[T], error) {

	// 参数校验
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 构建基础查询
	baseQuery := h.db.Model(new(T))
	if queryBuilder != nil {
		baseQuery = queryBuilder(baseQuery)
	}

	// 获取总数
	var total int64
	countQuery := baseQuery.Session(&gorm.Session{})
	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, fmt.Errorf("count query failed: %w", err)
	}

	// 计算偏移量并执行分页查询
	offset := (pageNum - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	var results []T
	err = baseQuery.Offset(offset).Limit(pageSize).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("page query failed: %w", err)
	}

	return entity.NewPaginationResult(results, total, pageNum, pageSize), nil
}

// PageQueryWithTransaction 在事务中执行分页查询
func (h *PaginationHelper[T]) PageQueryWithTransaction(tx *gorm.DB, pageNum, pageSize int, queryBuilder func(*gorm.DB) *gorm.DB) (*entity.PaginationResult[T], error) {

	helper := &PaginationHelper[T]{db: tx}
	return helper.PageQuery(pageNum, pageSize, queryBuilder)
}
