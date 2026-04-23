package initializer

import (
	users_domain "app/internal/module/iam/domain"
	users_infrastructure "app/internal/module/iam/infrastructure"
	methodology_domain "app/internal/module/methodology"
	project_domain "app/internal/module/project"
	"app/internal/module/scheduler"
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

	// 初始化默认调度任务配置
	if err := i.initDefaultScheduledTasks(); err != nil {
		return fmt.Errorf("failed to init default scheduled tasks: %w", err)
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

// initDefaultScheduledTasks 初始化默认调度任务配置
func (i *DataInitializer) initDefaultScheduledTasks() error {
	ctx := context.Background()
	repo := scheduler.NewScheduledTaskRepository()

	// 定义默认任务配置
	defaultTasks := []*scheduler.ScheduledTask{
		{
			Name:        "carbon_report_monthly_aggregation",
			CronSpec:    "0 0 1 3 * *", // 每月3号凌晨1点
			Description: "每月3号汇总上月碳日报数据生成月报",
			Enabled:     true,
			TaskType:    "carbon_report_monthly_aggregation",
			BaseEntity: entity.BaseEntity{
				CreateBy: 1,
			},
		},
		{
			Name:        "daily_log_output",
			CronSpec:    "0 0 0 * * *", // 每天凌晨0点
			Description: "每天输出调度任务运行日志",
			Enabled:     true,
			TaskType:    "daily_log_output",
			BaseEntity: entity.BaseEntity{
				CreateBy: 1,
			},
		},
	}

	for _, task := range defaultTasks {
		// 检查任务是否已存在
		existing, err := repo.FindByName(ctx, task.Name)
		if err != nil {
			logger.InitL.Warn("Failed to check existing task",
				zap.String("task_name", task.Name),
				zap.Error(err),
			)
			continue
		}

		if existing != nil {
			logger.InitL.Info("Task already exists, skipping",
				zap.String("task_name", task.Name),
			)
			continue
		}

		// 创建任务配置
		if err := repo.Create(ctx, task); err != nil {
			logger.InitL.Warn("Failed to create default task",
				zap.String("task_name", task.Name),
				zap.Error(err),
			)
			continue
		}

		logger.InitL.Info("Default task created",
			zap.String("task_name", task.Name),
			zap.String("cron_spec", task.CronSpec),
		)
	}

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
