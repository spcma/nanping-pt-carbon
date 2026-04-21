# Scheduler 定时任务模块 - HTTP API 文档

## 概述

定时任务（Scheduler）模块用于管理和执行定时任务，支持 Cron 表达式自定义执行频率，提供任务的增删改查、启用/禁用等功能。

**主要特性**：
- 支持自定义参数传入（使用 map 结构）
- 任务函数通过 key 获取参数值
- 参数持久化到数据库
- 支持动态创建和更新任务

---

## 基础信息

- **基础路径**: `/api/scheduler`
- **认证方式**: Token 认证（所有接口均需登录）
- **数据格式**: JSON
- **Cron 格式**: 支持秒级精度（6位：秒 分 时 日 月 周）

---

## Cron 表达式说明

### 格式
```
秒 分 时 日 月 周
```

### 示例

| Cron 表达式 | 说明 |
|------------|------|
| `0 * * * * *` | 每分钟执行 |
| `0 */5 * * * *` | 每5分钟执行 |
| `0 0 * * * *` | 每小时执行 |
| `0 0 2 * * *` | 每天凌晨2点执行 |
| `0 30 1 * * *` | 每天凌晨1:30执行 |
| `0 0 0 * * 1` | 每周一凌晨0点执行 |
| `0 0 0 1 * *` | 每月1号凌晨0点执行 |
| `0 0 0 1 1 *` | 每年1月1号凌晨0点执行 |

### 特殊字符

- `*`：任意值
- `,`：枚举值（如 `1,3,5`）
- `-`：范围值（如 `1-5`）
- `/`：步长（如 `*/5` 表示每5个单位）

---

## API 接口

### 1. 创建定时任务（带自定义参数）

**请求**：
```http
POST /api/scheduler/tasks
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "my_custom_task",
  "cron_spec": "0 */5 * * * *",
  "description": "每5分钟执行一次的自定义任务",
  "enabled": true,
  "task_type": "daily_log_output",
  "params": {
    "message": "自定义消息",
    "interval": 300,
    "retry_count": 3
  }
}
```

**请求参数说明**：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 任务名称，唯一标识 |
| cron_spec | string | 是 | Cron 表达式 |
| description | string | 否 | 任务描述 |
| enabled | boolean | 否 | 是否启用，默认 true |
| task_type | string | 是 | 任务类型，必须与注册的任务函数名称一致 |
| params | object | 否 | 自定义参数，map 结构，任务执行时通过 key 获取 |

**响应**：
```json
{
  "code": 200,
  "data": null,
  "message": "success"
}
```

**使用示例 - 碳月报汇总任务**：
```json
{
  "name": "carbon_monthly_agg",
  "cron_spec": "0 0 2 1 * *",
  "description": "每月1号凌晨2点执行碳月报汇总",
  "enabled": true,
  "task_type": "carbon_report_monthly_aggregation",
  "params": {
    "year": 2024,
    "month": 12
  }
}
```

---

### 2. 列出所有定时任务

**请求**：
```http
GET /api/scheduler/tasks
Authorization: Bearer <token>
```

**响应**：
```json
{
  "code": 200,
  "data": [
    {
      "name": "example_every_minute",
      "cron_spec": "0 * * * * *",
      "description": "示例任务：每分钟执行一次",
      "enabled": true,
      "last_run": "2026-04-15 10:30:00",
      "next_run": "2026-04-15 10:31:00",
      "total_runs": 150
    },
    {
      "name": "example_daily_2am",
      "cron_spec": "0 0 2 * * *",
      "description": "示例任务：每天凌晨2点执行",
      "enabled": true,
      "last_run": "2026-04-15 02:00:00",
      "next_run": "2026-04-16 02:00:00",
      "total_runs": 30
    }
  ],
  "message": "success"
}
```

---

### 3. 获取单个任务状态

**请求**：
```http
GET /api/scheduler/tasks/example_every_minute
Authorization: Bearer <token>
```

**响应**：
```json
{
  "code": 200,
  "data": {
    "name": "example_every_minute",
    "cron_spec": "0 * * * * *",
    "description": "示例任务：每分钟执行一次",
    "enabled": true,
    "last_run": "2026-04-15 10:30:00",
    "next_run": "2026-04-15 10:31:00",
    "total_runs": 150
  },
  "message": "success"
}
```

---

### 4. 移除定时任务

**请求**：
```http
DELETE /api/scheduler/tasks/example_every_minute
Authorization: Bearer <token>
```

**响应**：
```json
{
  "code": 200,
  "data": null,
  "message": "success"
}
```

**说明**：
- 任务被移除后将不再执行
- 可以通过重新添加任务来恢复

---

### 5. 启用定时任务

**请求**：
```http
PUT /api/scheduler/tasks/example_every_5min/enable
Authorization: Bearer <token>
```

**响应**：
```json
{
  "code": 200,
  "data": null,
  "message": "success"
}
```

---

### 6. 禁用定时任务

**请求**：
```http
PUT /api/scheduler/tasks/example_every_minute/disable
Authorization: Bearer <token>
```

**响应**：
```json
{
  "code": 200,
  "data": null,
  "message": "success"
}
```

---

## 错误码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（Token 无效或过期） |
| 404 | 任务不存在 |
| 500 | 服务器内部错误 |

**错误响应示例**：
```json
{
  "code": 404,
  "message": "task not_found not found",
  "data": null
}
```

---

## cURL 示例

### 列出所有任务
```bash
curl -X GET http://localhost:8080/api/scheduler/tasks \
  -H "Authorization: Bearer <your_token>"
```

### 获取任务状态
```bash
curl -X GET http://localhost:8080/api/scheduler/tasks/example_every_minute \
  -H "Authorization: Bearer <your_token>"
```

### 禁用任务
```bash
curl -X PUT http://localhost:8080/api/scheduler/tasks/example_every_minute/disable \
  -H "Authorization: Bearer <your_token>"
```

### 启用任务
```bash
curl -X PUT http://localhost:8080/api/scheduler/tasks/example_every_5min/enable \
  -H "Authorization: Bearer <your_token>"
```

### 删除任务
```bash
curl -X DELETE http://localhost:8080/api/scheduler/tasks/example_every_minute \
  -H "Authorization: Bearer <your_token>"
```

---

## 路由结构

```
/api/scheduler/
├── tasks                      # 任务管理
│   ├── GET                    # 获取所有任务
│   └── POST                   # 创建新任务（带参数）
│
└── tasks/:name                # 单个任务操作
    ├── GET                    # 获取任务状态
    ├── DELETE                 # 删除任务
    ├── PUT /enable            # 启用任务
    └── PUT /disable           # 禁用任务
```

---

## 开发指南

### 1. 注册自定义任务（支持参数）

在你的模块中注册自定义定时任务，任务函数接收 `params` 参数：

```go
import (
    "app/internal/module/scheduler"
    "context"
)

func init() {
    scheduler.RegisterTask("my_custom_task", func(ctx context.Context, params map[string]interface{}) error {
        // 从参数中获取配置
        if message, ok := params["message"]; ok {
            // 使用参数值
            logger.Info("Got message", zap.String("message", message.(string)))
        }
        
        // 实现你的业务逻辑
        return nil
    })
}
```

**参数获取示例**：

```go
// 获取字符串类型参数
if name, ok := params["name"].(string); ok {
    // 使用 name
}

// 获取数字类型参数（JSON 解析后为 float64）
if count, ok := params["count"].(float64); ok {
    intCount := int(count)
    // 使用 intCount
}

// 获取布尔类型参数
if enabled, ok := params["enabled"].(bool); ok {
    // 使用 enabled
}

// 安全获取参数，带默认值
retryCount := 3
if rc, ok := params["retry_count"].(float64); ok {
    retryCount = int(rc)
}
```

### 2. 添加任务到调度器（带参数）

```go
import (
    "app/internal/module/scheduler"
    "context"
)

func registerMyTask() {
    sched := scheduler.Default()
    
    sched.AddTask(&scheduler.TaskConfig{
        Name:     "my_task",
        CronSpec: "0 */10 * * * *", // 每10分钟
        TaskFunc: func(ctx context.Context, params map[string]interface{}) error {
            // 从参数中获取配置
            interval := 600 // 默认值
            if intervalParam, ok := params["interval"].(float64); ok {
                interval = int(intervalParam)
            }
            
            // 任务逻辑
            return nil
        },
        Description: "我的自定义任务",
        Enabled:     true,
        Params: map[string]interface{}{
            "interval": 600,
            "timeout":  30,
        },
    })
}
```

### 3. 任务执行上下文

每个任务都会接收到一个 `context.Context` 和 `params map[string]interface{}`，可以用于：
- 控制任务超时
- 传递请求级别的值
- 优雅关闭任务
- **通过 key 获取自定义参数**

---

## 注意事项

1. **Cron 表达式**：使用6位格式（包含秒），确保表达式正确
2. **任务幂等性**：定时任务应该是幂等的，避免重复执行造成问题
3. **错误处理**：任务执行失败不会影响其他任务，调度器会自动恢复
4. **并发执行**：如果上次任务未执行完，下次任务仍会按时启动
5. **日志记录**：所有任务执行都会记录日志，便于监控和调试
6. **性能考虑**：避免在任务中执行长时间运行的操作，必要时使用异步处理
7. **参数类型**：JSON 中的数字会被解析为 `float64`，使用时需要类型转换
8. **参数安全**：获取参数时建议使用类型断言检查，并提供默认值

---

## 内置示例任务

系统默认注册了3个示例任务：

1. **example_every_minute**：每分钟执行一次（已启用）
2. **example_daily_2am**：每天凌晨2点执行（已启用）
3. **example_every_5min**：每5分钟执行一次（默认禁用）

这些任务仅用于演示，可以在 `routes.go` 中修改或删除。

---

## 完整示例

### 示例1：创建带参数的任务

**请求**：
```bash
curl -X POST http://localhost:8080/api/scheduler/tasks \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "data_sync_task",
    "cron_spec": "0 */30 * * * *",
    "description": "每30分钟同步一次数据",
    "enabled": true,
    "task_type": "daily_log_output",
    "params": {
      "source": "mysql",
      "target": "postgresql",
      "batch_size": 1000,
      "timeout": 300,
      "retry_count": 3
    }
  }'
```

### 示例2：在任务函数中使用参数

```go
scheduler.RegisterTask("data_sync_task", func(ctx context.Context, params map[string]interface{}) error {
    // 获取参数，带默认值
    source := "mysql"
    if s, ok := params["source"].(string); ok {
        source = s
    }
    
    target := "postgresql"
    if t, ok := params["target"].(string); ok {
        target = t
    }
    
    batchSize := 1000
    if bs, ok := params["batch_size"].(float64); ok {
        batchSize = int(bs)
    }
    
    timeout := 300
    if to, ok := params["timeout"].(float64); ok {
        timeout = int(to)
    }
    
    retryCount := 3
    if rc, ok := params["retry_count"].(float64); ok {
        retryCount = int(rc)
    }
    
    logger.Info("Starting data sync",
        zap.String("source", source),
        zap.String("target", target),
        zap.Int("batch_size", batchSize),
        zap.Int("timeout", timeout),
        zap.Int("retry_count", retryCount),
    )
    
    // 执行数据同步逻辑
    // ...
    
    return nil
})
```

### 示例3：碳月报汇总任务（灵活配置年月）

**请求**：
```bash
curl -X POST http://localhost:8080/api/scheduler/tasks \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "carbon_agg_2024_12",
    "cron_spec": "0 0 10 15 1 *",
    "description": "2024年12月碳月报汇总",
    "enabled": true,
    "task_type": "carbon_report_monthly_aggregation",
    "params": {
      "year": 2024,
      "month": 12
    }
  }'
```

**说明**：
- 如果 `params` 中指定了 `year` 和 `month`，则使用指定的年月
- 如果未指定，则默认使用上个月的年月
