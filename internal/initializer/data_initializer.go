package initializer

import (
	users_domain "app/internal/module/iam/domain"
	users_infrastructure "app/internal/module/iam/infrastructure"
	"app/internal/shared/logger"
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DataInitializer 数据初始化器
type DataInitializer struct {
	db *gorm.DB
}

// NewDataInitializer 创建数据初始化器
func NewDataInitializer(db *gorm.DB) *DataInitializer {
	once := sync.Once{}
	once.Do(func() {
		dataInitializer = &DataInitializer{
			db: db,
		}
	})

	return dataInitializer
}

// Initialize 初始化所有基础数据
func (i *DataInitializer) Initialize() error {
	logger.Info("initialize", "starting data initialization...")

	// 初始化超级管理员用户
	if err := i.initSuperAdminUser(); err != nil {
		return fmt.Errorf("failed to init super admin user: %w", err)
	}

	logger.Info("initialize", "data initialization completed successfully")
	return nil
}

// initSuperAdminUser 初始化超级管理员用户
func (i *DataInitializer) initSuperAdminUser() error {
	ctx := context.Background()
	userRepo := users_infrastructure.NewUserRepository(i.db)

	// 检查用户是否已存在
	existingUser, err := userRepo.FindByUsername(ctx, "admin")
	if err == nil && existingUser != nil {
		logger.InitL.Info("admin user already exists",
			zap.String("username", existingUser.Username),
			zap.Int64("id", existingUser.Id))
		return nil
	}

	// 创建超级管理员用户
	// 默认密码：Admin@123
	user, err := users_domain.NewUser("admin", "系统管理员", "admin@2026", "", 0)
	if err != nil {
		return fmt.Errorf("failed to create super admin user: %w", err)
	}

	if err := userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to save super admin user: %w", err)
	}

	logger.InitL.Info("super admin user created successfully",
		zap.Int64("user_id", user.Id),
		zap.String("username", user.Username),
	)
	return nil
}
