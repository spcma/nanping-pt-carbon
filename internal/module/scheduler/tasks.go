package scheduler

import (
	"app/internal/shared/logger"
	"context"

	"go.uber.org/zap"
)

// ExampleCustomTask 自定义任务示例
// 你可以在任何地方调用此函数来注册自定义任务
func ExampleCustomTask() {
	RegisterTask("my_custom_task", func(ctx context.Context, params map[string]interface{}) error {
		logger.SchedulerL.Info("Executing my custom task",
			zap.String("info", "This is a custom task example"),
		)

		// 从参数中获取配置
		if paramValue, ok := params["example_param"]; ok {
			logger.SchedulerL.Info("Got parameter from config",
				zap.Any("example_param", paramValue),
			)
		}

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
