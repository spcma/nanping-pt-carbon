# Scheduler 定时任务模块 - HTTP API 文档

## 概述

定时任务（Scheduler）模块用于管理和执行定时任务，支持 Cron 表达式自定义执行频率，提供任务的增删改查、启用/禁用等功能。

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

### 1. 列出所有定时任务

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

### 2. 获取单个任务状态

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

### 3. 移除定时任务

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

### 4. 启用定时任务

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

### 5. 禁用定时任务

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
├── tasks                      # 任务列表
│   └── GET                    # 获取所有任务
│
└── tasks/:name                # 单个任务操作
    ├── GET                    # 获取任务状态
    ├── DELETE                 # 删除任务
    ├── PUT /enable            # 启用任务
    └── PUT /disable           # 禁用任务
```

---

## 开发指南

### 1. 注册自定义任务

在你的模块中注册自定义定时任务：

```go
import (
    "app/internal/module/scheduler"
    "context"
)

func init() {
    scheduler.RegisterTask("my_custom_task", func(ctx context.Context) error {
        // 实现你的业务逻辑
        return nil
    })
}
```

### 2. 添加任务到调度器

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
        TaskFunc: func(ctx context.Context) error {
            // 任务逻辑
            return nil
        },
        Description: "我的自定义任务",
        Enabled:     true,
    })
}
```

### 3. 任务执行上下文

每个任务都会接收到一个 `context.Context`，可以用于：
- 控制任务超时
- 传递请求级别的值
- 优雅关闭任务

---

## 注意事项

1. **Cron 表达式**：使用6位格式（包含秒），确保表达式正确
2. **任务幂等性**：定时任务应该是幂等的，避免重复执行造成问题
3. **错误处理**：任务执行失败不会影响其他任务，调度器会自动恢复
4. **并发执行**：如果上次任务未执行完，下次任务仍会按时启动
5. **日志记录**：所有任务执行都会记录日志，便于监控和调试
6. **性能考虑**：避免在任务中执行长时间运行的操作，必要时使用异步处理

---

## 内置示例任务

系统默认注册了3个示例任务：

1. **example_every_minute**：每分钟执行一次（已启用）
2. **example_daily_2am**：每天凌晨2点执行（已启用）
3. **example_every_5min**：每5分钟执行一次（默认禁用）

这些任务仅用于演示，可以在 `routes.go` 中修改或删除。
