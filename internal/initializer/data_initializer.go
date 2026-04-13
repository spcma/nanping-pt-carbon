package initializer

import (
	users_domain "app/internal/module/iam/domain"
	users_infrastructure "app/internal/module/iam/infrastructure"
	methodology_domain "app/internal/module/methodology"
	project_domain "app/internal/module/project"
	"app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/entity"
	"app/internal/shared/logger"
	"app/internal/shared/timeutil"
	"context"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
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

func (i *DataInitializer) Add_Methodology_20260323(c *gin.Context) {
	ctx := http.Ctx(c)
	projectRepo := methodology_domain.NewMethodologyRepository()

	var datas []methodology_domain.Methodology
	datas = append(datas,
		methodology_domain.Methodology{
			Name:        "项目1",
			Code:        "P1",
			Description: "项目1描述",
			Icon:        "",
			BaseEntity: entity.BaseEntity{
				CreateBy: 1,
			},
			StartDate: timeutil.Now(),
			EndDate:   timeutil.Now(),
		},
		methodology_domain.Methodology{
			Name:        "项目2",
			Code:        "P2",
			Description: "项目2描述",
			Icon:        "",
			BaseEntity: entity.BaseEntity{
				CreateBy: 1,
			},
			StartDate: timeutil.Now(),
			EndDate:   timeutil.Now(),
		},
		methodology_domain.Methodology{
			Name:        "项目3",
			Code:        "P3",
			Description: "项目3描述",
			Icon:        "",
			BaseEntity: entity.BaseEntity{
				CreateBy: 1,
			},
			StartDate: timeutil.Now(),
			EndDate:   timeutil.Now(),
		},
	)

	for _, data := range datas {
		p1, err := methodology_domain.NewMethodology(data.Name, data.Code, data.Icon, data.Description, 1, data.StartDate, data.EndDate)
		if err != nil {
			logger.RuntimeL.Error("new methodology error", zap.Any("data", data), zap.Error(err))
			response.BadRequest(c, fmt.Errorf("new methodology error: %w", err).Error())
			return
		}
		err = projectRepo.Create(ctx, p1)
		if err != nil {
			logger.RuntimeL.Error("add methodology error", zap.Any("data", data), zap.Error(err))
			response.BadRequest(c, fmt.Errorf("add methodology error: %w", err).Error())
			return
		}
	}

	response.Success(c, "success")
}

func (i *DataInitializer) Add_Project_20260323(c *gin.Context) {
	ctx := http.Ctx(c)
	projectRepo := project_domain.NewProjectRepository(i.db)

	var datas []project_domain.Project
	datas = append(datas,
		project_domain.Project{
			Name:        "项目1",
			Code:        "P1",
			Description: "项目1描述",
			Icon:        "",
			BaseEntity: entity.BaseEntity{
				CreateBy: 1,
			},
			StartDate: timeutil.Now(),
			EndDate:   timeutil.Now(),
		},
		project_domain.Project{
			Name:        "项目2",
			Code:        "P2",
			Description: "项目2描述",
			Icon:        "",
			BaseEntity: entity.BaseEntity{
				CreateBy: 1,
			},
			StartDate: timeutil.Now(),
			EndDate:   timeutil.Now(),
		},
		project_domain.Project{
			Name:        "项目3",
			Code:        "P3",
			Description: "项目3描述",
			Icon:        "",
			BaseEntity: entity.BaseEntity{
				CreateBy: 1,
			},
			StartDate: timeutil.Now(),
			EndDate:   timeutil.Now(),
		},
	)

	for _, data := range datas {
		p1, err := project_domain.NewProject(data.Name, data.Code, data.Icon, data.Description, 1, data.StartDate, data.EndDate)
		if err != nil {
			logger.RuntimeL.Error("add project error", zap.Any("data", data), zap.Error(err))
			response.BadRequest(c, fmt.Errorf("add project error: %w", err).Error())
			return
		}
		err = projectRepo.Create(ctx, p1)
		if err != nil {
			logger.RuntimeL.Error("add project error", zap.Any("data", data), zap.Error(err))
			response.BadRequest(c, fmt.Errorf("add project error: %w", err).Error())
			return
		}
	}

	response.Success(c, "success")
}
