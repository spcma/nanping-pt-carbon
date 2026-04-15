# 定时任务模块使用说明

## 快速开始

### 1. 模块已自动启动

定时任务模块已在服务器启动时自动初始化和启动，无需额外配置。

### 2. 查看现有任务

启动服务器后，访问以下API查看所有任务：

```bash
curl -X GET http://localhost:8080/api/scheduler/tasks \
  -H "Authorization: Bearer <your_token>"
```

### 3. 内置示例任务

系统默认包含3个示例任务：
- `example_every_minute` - 每分钟执行（已启用）
- `example_daily_2am` - 每天凌晨2点执行（已启用）
- `example_every_5min` - 每5分钟执行（默认禁用）

---

## 自定义任务开发

### 方式一：使用任务注册表（推荐）

在你的模块中注册任务函数：

```go
package yourmodule

import (
    "app/internal/module/scheduler"
    "app/internal/shared/logger"
    "context"
    
    "go.uber.org/zap"
)

func init() {
    // 注册任务函数
    scheduler.RegisterTask("my_task", func(ctx context.Context) error {
        logger.SchedulerL.Info("My task is running")
        
        // 实现你的业务逻辑
        // 例如：
        // - 查询数据库
        // - 生成报告
        // - 发送通知
        // - 数据同步
        // - 清理过期数据
        
        return nil
    })
}
```

### 方式二：直接添加到调度器

```go
package yourmodule

import (
    "app/internal/module/scheduler"
    "app/internal/shared/logger"
    "context"
    
    "go.uber.org/zap"
)

func registerMyTask() {
    sched := scheduler.Default()
    
    err := sched.AddTask(&scheduler.TaskConfig{
        Name:     "my_custom_task",
        CronSpec: "0 */10 * * * *", // 每10分钟执行
        TaskFunc: func(ctx context.Context) error {
            logger.SchedulerL.Info("Custom task executing",
                zap.String("info", "task details"),
            )
            
            // 业务逻辑
            return nil
        },
        Description: "我的自定义任务",
        Enabled:     true,
    })
    
    if err != nil {
        logger.SchedulerL.Error("Failed to add task",
            zap.Error(err),
        )
    }
}
```

---

## Cron 表达式指南

### 格式（6位）
```
秒 分 时 日 月 周
```

### 常用示例

| 表达式 | 说明 |
|--------|------|
| `0 * * * * *` | 每分钟 |
| `0 */5 * * * *` | 每5分钟 |
| `0 0 * * * *` | 每小时 |
| `0 30 2 * * *` | 每天凌晨2:30 |
| `0 0 0 * * 1` | 每周一0点 |
| `0 0 0 1 * *` | 每月1号0点 |
| `0 0 0 1 1 *` | 每年1月1日0点 |
| `0 0 9-18 * * *` | 每天9点到18点的整点 |
| `0 */30 9-18 * * *` | 工作日9-18点，每30分钟 |

### 特殊字符

- `*` - 任意值
- `,` - 枚举（如 `1,3,5`）
- `-` - 范围（如 `1-5`）
- `/` - 步长（如 `*/5`）

---

## 任务管理 API

### 列出所有任务
```http
GET /api/scheduler/tasks
```

### 查看任务状态
```http
GET /api/scheduler/tasks/{task_name}
```

### 启用任务
```http
PUT /api/scheduler/tasks/{task_name}/enable
```

### 禁用任务
```http
PUT /api/scheduler/tasks/{task_name}/disable
```

### 删除任务
```http
DELETE /api/scheduler/tasks/{task_name}
```

---

## 最佳实践

### 1. 任务幂等性
确保任务可以安全地重复执行：

```go
scheduler.RegisterTask("safe_task", func(ctx context.Context) error {
    // 使用事务确保数据一致性
    // 检查任务是否已执行
    // 避免重复处理
    return nil
})
```

### 2. 错误处理
任务失败不会影响其他任务，但应记录错误：

```go
scheduler.RegisterTask("robust_task", func(ctx context.Context) error {
    if err := doSomething(); err != nil {
        logger.SchedulerL.Error("Task failed",
            zap.Error(err),
        )
        return err // 返回错误以便记录
    }
    return nil
})
```

### 3. 超时控制
使用 context 控制任务执行时间：

```go
scheduler.RegisterTask("timeout_task", func(ctx context.Context) error {
    // 创建带超时的 context
    taskCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
    defer cancel()
    
    // 使用 taskCtx 执行操作
    return doLongTask(taskCtx)
})
```

### 4. 数据库访问
在任务中安全地访问数据库：

```go
import "app/internal/shared/db"

scheduler.RegisterTask("db_task", func(ctx context.Context) error {
    dbInst := db.Default()
    
    // 使用 context 执行查询
    var results []YourModel
    err := dbInst.WithContext(ctx).
        Where("status = ?", "pending").
        Find(&results).Error
    
    return err
})
```

### 5. 日志记录
详细记录任务执行情况：

```go
scheduler.RegisterTask("logged_task", func(ctx context.Context) error {
    logger.SchedulerL.Info("Task started")
    
    if err := execute(); err != nil {
        logger.SchedulerL.Error("Task failed",
            zap.Error(err),
        )
        return err
    }
    
    logger.SchedulerL.Info("Task completed successfully")
    return nil
})
```

---

## 监控和调试

### 1. 查看日志

定时任务日志记录在 `logs/other/scheduler/` 目录下。

### 2. 查看任务状态

通过 API 查看任务的执行统计：
- 最后执行时间
- 下次执行时间
- 总执行次数

### 3. 手动测试

在开发阶段，可以先禁用任务，然后通过 API 手动触发测试。

---

## 常见问题

### Q: 如何临时禁用任务？
A: 使用 API 禁用任务，无需删除：
```bash
curl -X PUT http://localhost:8080/api/scheduler/tasks/my_task/disable \
  -H "Authorization: Bearer <token>"
```

### Q: 任务执行失败会影响其他任务吗？
A: 不会。每个任务独立执行，失败不会影响其他任务。

### Q: 如何确保任务不重复执行？
A: 实现幂等性逻辑，或在执行前检查任务状态。

### Q: 任务可以传递参数吗？
A: 任务函数通过闭包捕获所需变量，或通过全局配置读取。

### Q: 如何调试任务？
A: 
1. 查看详细日志
2. 降低执行频率
3. 在开发环境测试
4. 使用 context 传递调试信息

---

## 示例：碳报告自动生成

```go
package scheduler

import (
    "app/internal/shared/db"
    "app/internal/shared/logger"
    "context"
    "time"
    
    "go.uber.org/zap"
)

func init() {
    RegisterTask("carbon_report_daily", func(ctx context.Context) error {
        logger.SchedulerL.Info("Generating daily carbon report")
        
        // 获取昨天的日期
        yesterday := time.Now().AddDate(0, 0, -1)
        
        // 查询昨天的数据
        dbInst := db.Default()
        // ... 查询和处理逻辑
        
        logger.SchedulerL.Info("Daily carbon report generated",
            zap.Time("report_date", yesterday),
        )
        
        return nil
    })
}
```

---

## 注意事项

1. ⚠️ **避免长时间阻塞**：任务应尽快完成，必要时使用异步处理
2. ⚠️ **注意并发**：如果上次任务未完成，下次仍会启动
3. ⚠️ **资源清理**：确保任务完成后清理资源（关闭连接等）
4. ⚠️ **时间敏感**：注意时区和夏令时问题
5. ⚠️ **依赖检查**：确保任务所需的依赖（数据库、外部服务）可用

---

## 更多信息

- API 文档：[api.md](./api.md)
- HTTP 测试：[scheduler.http](./scheduler.http)
