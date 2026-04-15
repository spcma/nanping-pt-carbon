package scheduler

import (
	"app/internal/shared/logger"
	"context"

	"go.uber.org/zap"
)

// RegisterBusinessTasks 注册业务定时任务
// 在各个模块初始化时调用此函数注册具体的业务任务
func RegisterBusinessTasks() {
	// 示例：注册碳报告自动生成任务
	RegisterTask("carbon_report_daily_generation", func(ctx context.Context) error {
		logger.SchedulerL.Info("Executing carbon report daily generation task")
		// TODO: 实现具体的业务逻辑
		// 例如：自动生成昨天的碳报告日报
		return nil
	})

	// 示例：注册数据清理任务
	RegisterTask("data_cleanup_weekly", func(ctx context.Context) error {
		logger.SchedulerL.Info("Executing data cleanup task")
		// TODO: 实现具体的业务逻辑
		// 例如：清理过期的临时数据
		return nil
	})

	// 示例：注册数据同步任务
	RegisterTask("data_sync_hourly", func(ctx context.Context) error {
		logger.SchedulerL.Info("Executing data sync task")
		// TODO: 实现具体的业务逻辑
		// 例如：同步外部系统数据
		return nil
	})

	logger.SchedulerL.Info("Business tasks registered")
}

// ExampleCustomTask 自定义任务示例
// 你可以在任何地方调用此函数来注册自定义任务
func ExampleCustomTask() {
	RegisterTask("my_custom_task", func(ctx context.Context) error {
		logger.SchedulerL.Info("Executing my custom task",
			zap.String("info", "This is a custom task example"),
		)

		// 在这里实现你的业务逻辑
		// 例如：
		// 1. 查询数据库
		// 2. 调用外部API
		// 3. 生成报告
		// 4. 发送通知
		// 等等...

		return nil
	})
}
