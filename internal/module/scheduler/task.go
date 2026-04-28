package scheduler

import (
	"app/internal/shared/logger"
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// TaskFunc 定时任务函数类型
type TaskFunc func(ctx context.Context, params map[string]interface{}) error

// TaskConfig 定时任务配置
type TaskConfig struct {
	Name        string                 // 任务名称
	CronSpec    string                 // Cron 表达式
	TaskFunc    TaskFunc               // 任务执行函数
	Description string                 // 任务描述
	Enabled     bool                   // 是否启用
	Params      map[string]interface{} // 任务参数
}

var (
	// 全局任务注册表
	defaultRegistry *TaskRegistry
	_once           sync.Once
)

// TaskRegistry 任务注册表
type TaskRegistry struct {
	tasks map[string]TaskFunc
}

// RegistryStore 获取默认任务注册表
func RegistryStore() *TaskRegistry {
	_once.Do(func() {
		defaultRegistry = &TaskRegistry{
			tasks: make(map[string]TaskFunc),
		}
	})
	return defaultRegistry
}

// RegisterTask 便捷方法：注册任务函数
func RegisterTask(name string, fn TaskFunc) {
	defaultRegistry.Register(name, fn)
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

// registerTaskFunctions 注册任务函数到注册表(不直接添加到调度器)
func RegisterTaskFunctions() {

	// 报时
	RegisterTask("report.time", func(ctx context.Context, params map[string]interface{}) error {
		cst := time.Now()

		logger.SchedulerL.Info("调度任务运行中",
			zap.String("current_time", cst.Format("2006-01-02 15:04:05")),
			zap.String("message", "调度任务运行中"),
		)

		if customMsg, ok := params["message"]; ok {
			logger.SchedulerL.Info("Custom message from params",
				zap.Any("message", customMsg),
			)
		}

		logger.SchedulerL.Info("调度任务运行完毕",
			zap.Duration("cost", time.Since(cst)),
		)

		return nil
	})

	logger.SchedulerL.Info("Task functions registered")
}
