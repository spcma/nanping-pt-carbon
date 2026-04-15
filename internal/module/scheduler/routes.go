package scheduler

import (
	shared_http "app/internal/shared/http"
	"app/internal/shared/logger"
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TaskRegistry 任务注册表
type TaskRegistry struct {
	tasks map[string]TaskFunc
}

// NewTaskRegistry 创建任务注册表
func NewTaskRegistry() *TaskRegistry {
	return &TaskRegistry{
		tasks: make(map[string]TaskFunc),
	}
}

// Register 注册任务函数
func (r *TaskRegistry) Register(name string, fn TaskFunc) {
	r.tasks[name] = fn
	logger.SchedulerL.Info("Task function registered",
		zap.String("task_name", name),
	)
}

// Get 获取任务函数
func (r *TaskRegistry) Get(name string) (TaskFunc, bool) {
	fn, exists := r.tasks[name]
	return fn, exists
}

// 全局任务注册表
var defaultRegistry = NewTaskRegistry()

// GetRegistry 获取默认任务注册表
func GetRegistry() *TaskRegistry {
	return defaultRegistry
}

// RegisterTask 便捷方法：注册任务函数
func RegisterTask(name string, fn TaskFunc) {
	defaultRegistry.Register(name, fn)
}

// routes 路由注册器
type routes struct {
	handler *SchedulerHandler
}

// NewSchedulerRoutes 创建定时任务模块的路由注册器
func NewSchedulerRoutes() shared_http.RouteRegistry {
	scheduler := Default()
	handler := NewSchedulerHandler(scheduler)

	// 注册示例任务
	registerExampleTasks(scheduler)

	return &routes{
		handler: handler,
	}
}

// registerExampleTasks 注册示例任务
func registerExampleTasks(scheduler *Scheduler) {
	// 示例任务1：每分钟执行一次
	scheduler.AddTask(&TaskConfig{
		Name:     "example_every_minute",
		CronSpec: "0 * * * * *", // 每分钟
		TaskFunc: func(ctx context.Context) error {
			logger.SchedulerL.Info("Example task: every minute")
			return nil
		},
		Description: "示例任务：每分钟执行一次",
		Enabled:     true,
	})

	// 示例任务2：每天凌晨2点执行
	scheduler.AddTask(&TaskConfig{
		Name:     "example_daily_2am",
		CronSpec: "0 0 2 * * *", // 每天凌晨2点
		TaskFunc: func(ctx context.Context) error {
			logger.SchedulerL.Info("Example task: daily at 2 AM")
			return nil
		},
		Description: "示例任务：每天凌晨2点执行",
		Enabled:     true,
	})

	// 示例任务3：每5分钟执行一次
	scheduler.AddTask(&TaskConfig{
		Name:     "example_every_5min",
		CronSpec: "0 */5 * * * *", // 每5分钟
		TaskFunc: func(ctx context.Context) error {
			logger.SchedulerL.Info("Example task: every 5 minutes")
			return nil
		},
		Description: "示例任务：每5分钟执行一次",
		Enabled:     false, // 默认禁用
	})
}

// RegisterRoutes 注册路由
func (r *routes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {
	// 统一认证中间件
	authTypeRequiredRoute := group.Group("")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		authTypeRequiredRoute.Use(mw)
	}
	{
		// 定时任务管理路由 - /api/scheduler/*
		schedulerGroup := authTypeRequiredRoute.Group("/scheduler")
		{
			// 任务列表
			schedulerGroup.GET("/tasks", r.handler.ListTasks)

			// 单个任务操作
			schedulerGroup.GET("/tasks/:name", r.handler.GetTaskStatus)
			schedulerGroup.DELETE("/tasks/:name", r.handler.RemoveTask)
			schedulerGroup.PUT("/tasks/:name/enable", r.handler.EnableTask)
			schedulerGroup.PUT("/tasks/:name/disable", r.handler.DisableTask)
		}
	}
}
