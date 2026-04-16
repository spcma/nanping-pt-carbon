# 定时任务模块使用说明

## 核心特性

✅ **自动持久化**: 任务配置自动保存到数据库，应用重启后自动加载
✅ **动态管理**: 通过 API 动态添加、修改、启用/禁用任务
✅ **任务注册表**: 支持注册任务函数，与任务配置解耦
✅ **完整监控**: 查看任务执行状态、最后执行时间、下次执行时间等

---

## 快速开始

### 1. 模块已自动启动

定时任务模块已在服务器启动时自动初始化和启动，无需额外配置。

**首次启动流程：**
1. 系统自动创建 `sys_scheduled_task` 表
2. 初始化默认任务配置到数据库
3. 注册所有任务函数到注册表
4. 从数据库加载已启用的任务并启动调度

**重启流程：**
1. 注册所有任务函数到注册表
2. 从数据库加载已启用的任务配置
3. 根据配置自动重新注册并启动任务

### 2. 查看现有任务

启动服务器后，访问以下API查看所有任务：

```bash
curl -X GET http://localhost:8080/api/scheduler/tasks \
  -H "Authorization: Bearer <your_token>"
```

### 3. 默认任务配置

系统首次启动时会自动创建以下默认任务：
- `carbon_report_monthly_aggregation` - 每月3号凌晨1点执行，汇总上月碳日报数据生成月报（已启用）
- `daily_log_output` - 每天凌晨0点执行，输出调度任务运行日志（已启用）

**注意**: 示例任务不再自动创建，需要时可手动添加。

---

## 自定义任务开发

### 步骤一：注册任务函数

在你的模块中注册任务函数到全局注册表：

```go
package yourmodule

import (
    "app/internal/module/scheduler"
    "app/internal/shared/logger"
    "context"
    
    "go.uber.org/zap"
)

func init() {
    // 注册任务函数（不直接添加到调度器）
    scheduler.RegisterTask("my_task_type", func(ctx context.Context) error {
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

**重要**: `RegisterTask` 只是注册任务函数，不会立即启动任务。任务的执行时间、是否启用等配置存储在数据库中。

### 步骤二：添加任务配置

通过 API 或直接在数据库中添加任务配置：

#### 方式 A: 通过 API 添加（推荐）

```bash
curl -X POST http://localhost:8080/api/scheduler/tasks \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my_custom_task",
    "cron_spec": "0 */10 * * * *",
    "description": "每10分钟执行的自定义任务",
    "enabled": true,
    "task_type": "my_task_type"
  }'
```

系统会自动：
1. 将任务配置保存到数据库
2. 从注册表中查找对应的任务函数
3. 注册并启动任务

#### 方式 B: 直接在数据库中添加

```sql
INSERT INTO sys_scheduled_task (name, cron_spec, description, enabled, task_type, create_by, create_time, update_time, delete_by, delete_time)
VALUES ('my_custom_task', '0 */10 * * * *', '每10分钟执行的自定义任务', true, 'my_task_type', 1, NOW(), NOW(), 0, '1970-01-01 00:00:00+08');
```

然后重启应用，系统会自动加载并启动该任务。

### 完整示例：碳报告每日生成

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

### Q: 应用重启后任务会自动加载吗？
A: **是的！** 这是本模块的核心特性：
- 所有任务配置都保存在数据库的 `sys_scheduled_task` 表中
- 应用启动时，系统会：
  1. 注册所有任务函数到注册表
  2. 从数据库读取所有启用的任务配置
  3. 根据配置自动重新注册并启动任务
- 你无需任何额外操作，任务会自动恢复运行

### Q: 如何修改任务的执行时间？
A: 有两种方式：
1. **通过 API**（推荐）：调用更新接口修改 cron 表达式
2. **直接修改数据库**：更新 `sys_scheduled_task` 表中的 `cron_spec` 字段，然后重启应用

### Q: 如何永久删除一个任务？
A: 
1. 通过 API 删除任务（会逻辑删除数据库记录）
2. 或直接删除数据库中的记录
3. 重启应用后该任务将不再加载

---

## 示例：碳报告自动生成

### 1. 注册任务函数

```go
package scheduler

import (
    "app/internal/module/carbonreportday"
    "app/internal/shared/logger"
    "context"
    "time"
    
    "go.uber.org/zap"
)

func init() {
    // 注册任务函数
    RegisterTask("carbon_report_daily", func(ctx context.Context) error {
        logger.SchedulerL.Info("Generating daily carbon report")
        
        // 获取昨天的日期
        yesterday := time.Now().AddDate(0, 0, -1)
        year := yesterday.Year()
        month := int(yesterday.Month())
        day := yesterday.Day()
        
        // 调用业务服务生成报告
        service := carbonreportday.DefaultService()
        if service == nil {
            logger.SchedulerL.Error("Carbon report service not initialized")
            return nil
        }
        
        err := service.GenerateDailyReport(ctx, year, month, day)
        if err != nil {
            logger.SchedulerL.Error("Failed to generate daily report",
                zap.Int("year", year),
                zap.Int("month", month),
                zap.Int("day", day),
                zap.Error(err),
            )
            return err
        }
        
        logger.SchedulerL.Info("Daily carbon report generated successfully",
            zap.Int("year", year),
            zap.Int("month", month),
            zap.Int("day", day),
        )
        
        return nil
    })
}
```

### 2. 添加任务配置

```bash
curl -X POST http://localhost:8080/api/scheduler/tasks \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "carbon_report_daily",
    "cron_spec": "0 0 1 * * *",
    "description": "每天凌晨1点生成昨日碳报告",
    "enabled": true,
    "task_type": "carbon_report_daily"
  }'
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

---

## 数据库表结构

### sys_scheduled_task 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键ID |
| name | VARCHAR | 任务名称（唯一） |
| cron_spec | VARCHAR | Cron表达式 |
| description | VARCHAR | 任务描述 |
| enabled | BOOLEAN | 是否启用 |
| task_type | VARCHAR | 任务类型（对应注册表中的函数名） |
| create_by | BIGINT | 创建人ID |
| update_by | BIGINT | 更新人ID |
| delete_by | BIGINT | 删除人ID（0表示未删除） |
| create_time | TIMESTAMP | 创建时间 |
| update_time | TIMESTAMP | 更新时间 |
| delete_time | TIMESTAMP | 删除时间 |

**注意**: `task_type` 必须与通过 `RegisterTask` 注册的任务函数名称一致。
