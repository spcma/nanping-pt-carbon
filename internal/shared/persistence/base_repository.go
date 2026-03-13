package persistence

import (
	"context"

	"gorm.io/gorm"
)

// BaseRepository 基础仓储实现
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// DB 获取数据库连接
func (r *BaseRepository) DB() *gorm.DB {
	return r.db
}

// WithContext 获取带上下的数据库连接
func (r *BaseRepository) WithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Transaction 执行事务
func (r *BaseRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}
