// internal/infrastructure/db/transaction/unit_of_work.go
package db

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// 工作单元接口
type UnitOfWork interface {
	Execute(ctx context.Context, fn func(ctx context.Context) error) error
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	GetDB(ctx context.Context) *gorm.DB
}

// 统一工作单元实现
type GormUnitOfWork struct {
	db     *gorm.DB // 原始数据库连接（用于开启事务）
	txDB   *gorm.DB // 当前事务 DB（Begin 后设置）
	dbName string   // 数据源名称
	hasTx  bool     // 是否有活动事务
}

func NewGormUnitOfWork(db *gorm.DB) *GormUnitOfWork {
	return &GormUnitOfWork{
		db:     db,
		dbName: "default",
	}
}

// NewGormUnitOfWorkWithName 创建指定数据源的工作单元
func NewGormUnitOfWorkWithName(db *gorm.DB, dbName string) *GormUnitOfWork {
	return &GormUnitOfWork{
		db:     db,
		dbName: dbName,
	}
}

// 核心事务执行方法
func (uow *GormUnitOfWork) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	// 开始事务
	txCtx, err := uow.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// 确保回滚
	defer func() {
		if r := recover(); r != nil {
			_ = uow.Rollback(txCtx)
			panic(r)
		}
	}()

	// 执行业务逻辑
	if err := fn(txCtx); err != nil {
		_ = uow.Rollback(txCtx)
		return err
	}

	// 提交事务
	return uow.Commit(txCtx)
}

func (uow *GormUnitOfWork) Begin(ctx context.Context) (context.Context, error) {
	if uow.hasTx {
		return ctx, fmt.Errorf("transaction already started")
	}

	tx := uow.db.Begin()
	if tx.Error != nil {
		return ctx, tx.Error
	}

	// 设置事务状态
	uow.txDB = tx
	uow.hasTx = true

	// 将事务对象存入context，供 Repository 层使用
	txCtx := context.WithValue(ctx, "tx_db", tx)
	return txCtx, nil
}

func (uow *GormUnitOfWork) Commit(ctx context.Context) error {
	if !uow.hasTx || uow.txDB == nil {
		return fmt.Errorf("no active transaction")
	}

	err := uow.txDB.Commit().Error
	// 清理事务状态
	uow.txDB = nil
	uow.hasTx = false
	return err
}

func (uow *GormUnitOfWork) Rollback(ctx context.Context) error {
	if !uow.hasTx || uow.txDB == nil {
		return fmt.Errorf("no active transaction")
	}

	err := uow.txDB.Rollback().Error
	// 清理事务状态
	uow.txDB = nil
	uow.hasTx = false
	return err
}

// GetDB 获取当前数据库实例
// 如果有活动事务，返回事务 DB
// 否则返回原始 DB
func (uow *GormUnitOfWork) GetDB(ctx context.Context) *gorm.DB {
	if uow.hasTx && uow.txDB != nil {
		return uow.txDB
	}
	return uow.db
}

// GetDBName 获取数据源名称
func (uow *GormUnitOfWork) GetDBName() string {
	return uow.dbName
}

// HasTransaction 检查是否有活动事务
func (uow *GormUnitOfWork) HasTransaction() bool {
	return uow.hasTx
}
