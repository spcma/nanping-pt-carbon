// internal/infrastructure/db/transaction/unit_of_work.go
package transaction

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
	db *gorm.DB
}

func NewGormUnitOfWork(db *gorm.DB) *GormUnitOfWork {
	return &GormUnitOfWork{db: db}
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
			uow.Rollback(txCtx)
			panic(r)
		}
	}()

	// 执行业务逻辑
	if err := fn(txCtx); err != nil {
		uow.Rollback(txCtx)
		return err
	}

	// 提交事务
	return uow.Commit(txCtx)
}

func (uow *GormUnitOfWork) Begin(ctx context.Context) (context.Context, error) {
	tx := uow.db.Begin()
	if tx.Error != nil {
		return ctx, tx.Error
	}

	// 将事务对象存入context
	txCtx := context.WithValue(ctx, "tx_db", tx)
	return txCtx, nil
}

func (uow *GormUnitOfWork) Commit(ctx context.Context) error {
	if tx := uow.GetDB(ctx); tx != nil {
		return tx.Commit().Error
	}
	return fmt.Errorf("no transaction found")
}

func (uow *GormUnitOfWork) Rollback(ctx context.Context) error {
	if tx := uow.GetDB(ctx); tx != nil {
		return tx.Rollback().Error
	}
	return fmt.Errorf("no transaction found")
}

func (uow *GormUnitOfWork) GetDB(ctx context.Context) *gorm.DB {
	// 优先使用事务中的DB
	if tx, ok := ctx.Value("tx_db").(*gorm.DB); ok {
		return tx
	}
	// 如果没有事务，返回普通DB
	return uow.db
}
