package initializer

import (
	users_domain "app/internal/module/iam/domain"
	users_infrastructure "app/internal/module/iam/infrastructure"
	methodology_domain "app/internal/module/methodology/domain"
	methodology_infrastructure "app/internal/module/methodology/infrastructure"
	project_domain "app/internal/module/project/domain"
	project_infrastructure "app/internal/module/project/infrastructure"
	"app/internal/shared/logger"
	"app/internal/shared/timeutil"
	"context"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DataInitializer 数据初始化器
type DataInitializer struct {
	db *gorm.DB
}

// NewDataInitializer 创建数据初始化器
func NewDataInitializer(db *gorm.DB) *DataInitializer {
	return &DataInitializer{
		db: db,
	}
}

// Initialize 初始化所有基础数据
func (i *DataInitializer) Initialize() error {
	logger.Info("initialize", "starting data initialization...")

	// 初始化超级管理员用户
	if err := i.initSuperAdminUser(); err != nil {
		return fmt.Errorf("failed to init super admin user: %w", err)
	}

	//i.initProject_20260318()

	i.initMethodology_20260318()

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
		logger.InitLogger.Info("admin user already exists",
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

	logger.InitLogger.Info("super admin user created successfully",
		zap.Int64("user_id", user.Id),
		zap.String("username", user.Username),
	)
	return nil
}

func (i *DataInitializer) initProject_20260318() error {
	ctx := context.Background()
	projectRepo := project_infrastructure.NewProjectRepository(i.db)

	p1, err := project_domain.NewProject("项目1", "P1", "", "项目1描述", 1, timeutil.Now(), timeutil.Now())
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	projectRepo.Create(ctx, p1)

	p2, err := project_domain.NewProject("项目2", "P2", "", "项目2描述", 1, timeutil.Now(), timeutil.Now())
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	projectRepo.Create(ctx, p2)

	p3, err := project_domain.NewProject("项目3", "P3", "", "项目3描述", 1, timeutil.Now(), timeutil.Now())
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	projectRepo.Create(ctx, p3)

	return nil
}

func (i *DataInitializer) initMethodology_20260318() error {
	ctx := context.Background()
	projectRepo := methodology_infrastructure.NewMethodologyRepository(i.db)

	p1, err := methodology_domain.NewMethodology("项目1", "P1", "", "项目1描述", 1, timeutil.Now(), timeutil.Now())
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	projectRepo.Create(ctx, p1)

	p2, err := methodology_domain.NewMethodology("项目2", "P2", "", "项目2描述", 1, timeutil.Now(), timeutil.Now())
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	projectRepo.Create(ctx, p2)

	p3, err := methodology_domain.NewMethodology("项目3", "P3", "", "项目3描述", 1, timeutil.Now(), timeutil.Now())
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	projectRepo.Create(ctx, p3)

	return nil
}
