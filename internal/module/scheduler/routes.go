package scheduler

import (
	"app/internal/module/carbonreportmonth"
	"app/internal/shared/db"
	shared_http "app/internal/shared/http"
	"app/internal/shared/logger"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var defaultService *Scheduler

func setDefaultService(sc *Scheduler) {
	defaultService = sc
}

func DefaultService() *Scheduler {
	return defaultService
}

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
// 注意: 此函数会在 InitRouter() 中被调用,此时所有业务模块的 Service 已经初始化完成
// 调度任务的 Start() 会在所有路由注册完成后才调用,确保任务执行时依赖的 Service 已就绪
func NewSchedulerRoutes() shared_http.RouteRegistry {
	scheduler := Default()

	setDefaultService(scheduler)

	handler := NewSchedulerHandler(scheduler)

	// 设置仓储
	repo := NewScheduledTaskRepository(db.Default())
	scheduler.SetRepository(repo)

	// 注册业务任务函数到注册表(不直接添加到调度器)
	registerTaskFunctions()

	// 从数据库加载已启用的任务
	if err := scheduler.LoadTasksFromDatabase(defaultRegistry); err != nil {
		logger.SchedulerL.Error("Failed to load tasks from database",
			zap.Error(err),
		)
	}

	return &routes{
		handler: handler,
	}
}

// registerTaskFunctions 注册任务函数到注册表(不直接添加到调度器)
func registerTaskFunctions() {
	// 注册每月碳日报汇总任务
	RegisterTask("carbon_report_monthly_aggregation", func(ctx context.Context, params map[string]interface{}) error {
		logger.SchedulerL.Info("执行每月碳日报汇总任务")

		// 从参数中获取年份和月份，如果没有则使用上个月
		var year, month int
		if yearParam, ok := params["year"]; ok {
			if y, valid := yearParam.(float64); valid {
				year = int(y)
			}
		}
		if monthParam, ok := params["month"]; ok {
			if m, valid := monthParam.(float64); valid {
				month = int(m)
			}
		}

		// 如果参数中没有指定，则使用上个月
		if year == 0 || month == 0 {
			now := time.Now()
			lastMonth := now.AddDate(0, -1, 0)
			year = lastMonth.Year()
			month = int(lastMonth.Month())
		}

		// 调用汇总服务
		service := carbonreportmonth.DefaultService()
		if service == nil {
			logger.SchedulerL.Error("碳报告月报服务未初始化")
			return nil
		}

		err := service.AggregateMonthlyReport(ctx, year, month)
		if err != nil {
			logger.SchedulerL.Error("碳日报汇总任务执行失败",
				zap.Int("year", year),
				zap.Int("month", month),
				zap.Error(err),
			)
			return err
		}

		logger.SchedulerL.Info("碳日报汇总任务执行成功",
			zap.Int("year", year),
			zap.Int("month", month),
		)
		return nil
	})

	// 注册每天日志输出任务
	RegisterTask("daily_log_output", func(ctx context.Context, params map[string]interface{}) error {
		now := time.Now()
		logger.SchedulerL.Info("调度任务运行中",
			zap.String("current_time", now.Format("2006-01-02 15:04:05")),
			zap.String("message", "调度任务运行中"),
		)

		// 示例：从参数中获取自定义消息
		if customMsg, ok := params["message"]; ok {
			logger.SchedulerL.Info("Custom message from params",
				zap.Any("message", customMsg),
			)
		}

		return nil
	})

	logger.SchedulerL.Info("Task functions registered")
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
