package transaction

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// TransactionManager 事务管理器
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager 创建事务管理器
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &TransactionManager{db: db}
}

// ExecuteInTransaction 在事务中执行操作
func (tm *TransactionManager) ExecuteInTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			if rbErr := tx.Rollback().Error; rbErr != nil {
				fmt.Printf("failed to rollback transaction after panic: %v\n", rbErr)
			}
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("transaction execution failed: %v, rollback error: %w", err, rbErr)
		}
		return fmt.Errorf("transaction execution failed: %w", err)
	}

	if commitErr := tx.Commit().Error; commitErr != nil {
		return fmt.Errorf("failed to commit transaction: %w", commitErr)
	}

	return nil
}

// RunInTransaction 简化版事务执行函数（包级别函数）
func RunInTransaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	if db == nil {
		return errors.New("database connection cannot be nil")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
