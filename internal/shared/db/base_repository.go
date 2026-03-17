package db

import (
	"app/internal/shared/entity"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// BaseRepository 通用Repository基类
// 支持事务感知和通用CRUD操作
type BaseRepository[T any] struct {
	db *gorm.DB
}

// TransactionalRepository 支持事务的仓储接口
type TransactionalRepository interface {
	GetDB(ctx context.Context) *gorm.DB
}

// NewBaseRepository 创建基础Repository
func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

// GetDB 获取正确的数据库连接（事务或普通）
// 实现 TransactionalRepository 接口
func (r *BaseRepository[T]) GetDB(ctx context.Context) *gorm.DB {
	// 检查context中是否有事务DB
	if txDB, ok := ctx.Value("tx_db").(*gorm.DB); ok {
		return txDB
	}
	// 否则使用普通DB
	return r.db
}

const (
	//	默认分页参数
	//	默认页码: 1
	//	默认每页数量: 10
	DefaultPageNum  = 1
	DefaultPageSize = 10
)

// PageQuery 通用分页查询方法
func (r *BaseRepository[T]) PageQuery(ctx context.Context, pageNum, pageSize int, buildQuery func(*gorm.DB) *gorm.DB) (*entity.PaginationResult[T], error) {
	// 使用事务感知的DB连接
	txDB := r.GetDB(ctx)

	// 参数校验
	if pageNum <= 0 {
		pageNum = DefaultPageNum
	}
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}

	// 构建查询
	query := txDB.WithContext(ctx).Model(new(T))
	if buildQuery != nil {
		query = buildQuery(query)
	}

	// 获取总数
	var total int64
	countQuery := query.Session(&gorm.Session{})
	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, fmt.Errorf("count query failed: %w", err)
	}

	// 执行分页查询
	offset := (pageNum - 1) * pageSize
	var results []T
	err = query.Offset(offset).Limit(pageSize).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("page query failed: %w", err)
	}

	return entity.NewPaginationResult(results, total, pageNum, pageSize), nil
}

// WithSoftDelete 软删除过滤条件
func (r *BaseRepository[T]) WithSoftDelete() func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where("deleted_at IS NULL")
	}
}

// WithOrderBy 排序条件
func WithOrderBy(field string, ascending bool) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		order := field
		if !ascending {
			order += " DESC"
		}
		return query.Order(order)
	}
}

// WithSoftDelete 添加软删除过滤条件
func WithSoftDelete(query *gorm.DB) *gorm.DB {
	return query.Where(entity.FieldDeleteBy + " = 0")
}

// CombineConditions 组合多个查询条件
func CombineConditions(conditions ...func(*gorm.DB) *gorm.DB) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		for _, condition := range conditions {
			if condition != nil {
				query = condition(query)
			}
		}
		return query
	}
}
