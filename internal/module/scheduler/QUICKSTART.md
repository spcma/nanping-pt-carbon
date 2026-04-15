# 定时任务模块 - 快速开始

## ✅ 模块已完成

定时任务模块已经完整实现并集成到项目中，支持自定义 Cron 表达式设置执行频率。

---

## 📁 文件结构

```
internal/module/scheduler/
├── scheduler.go      # 核心调度器实现
├── handler.go        # HTTP 处理器
├── routes.go         # 路由注册和示例任务
├── tasks.go          # 业务任务注册示例
├── api.md            # API 文档
├── scheduler.http    # HTTP 测试文件
└── README.md         # 详细使用说明
```

---

## 🚀 快速使用

### 1. 启动服务器

定时任务模块会随服务器自动启动：

```bash
go run cmd/main.go
```

### 2. 查看任务列表

```bash
curl -X GET http://localhost:8080/api/scheduler/tasks \
  -H "Authorization: Bearer <your_token>"
```

### 3. 内置示例任务

系统默认包含3个示例任务：
- ✅ `example_every_minute` - 每分钟执行
- ✅ `example_daily_2am` - 每天凌晨2点执行
- ⭕ `example_every_5min` - 每5分钟执行（默认禁用）

---

## 🎯 自定义任务

### 方法一：使用任务注册表（推荐）

```go
import (
    "app/internal/module/scheduler"
    "context"
)

func init() {
    scheduler.RegisterTask("my_task", func(ctx context.Context) error {
        // 你的业务逻辑
        return nil
    })
}
```

### 方法二：直接添加到调度器

```go
import (
    "app/internal/module/scheduler"
    "context"
)

func registerTask() {
    sched := scheduler.Default()
    
    sched.AddTask(&scheduler.TaskConfig{
        Name:     "my_custom_task",
        CronSpec: "0 */10 * * * *", // 每10分钟
        TaskFunc: func(ctx context.Context) error {
            // 你的业务逻辑
            return nil
        },
        Description: "我的自定义任务",
        Enabled:     true,
    })
}
```

---

## ⏰ Cron 表达式

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

---

## 📡 管理 API

| 操作 | 方法 | 路径 |
|------|------|------|
| 列出所有任务 | GET | `/api/scheduler/tasks` |
| 查看任务状态 | GET | `/api/scheduler/tasks/{name}` |
| 启用任务 | PUT | `/api/scheduler/tasks/{name}/enable` |
| 禁用任务 | PUT | `/api/scheduler/tasks/{name}/disable` |
| 删除任务 | DELETE | `/api/scheduler/tasks/{name}` |

---

## 🔧 已完成的功能

- ✅ 核心调度器（基于 robfig/cron/v3）
- ✅ 支持秒级精度的 Cron 表达式
- ✅ 任务动态添加/删除/启用/禁用
- ✅ 任务执行统计（执行次数、最后/下次执行时间）
- ✅ HTTP API 管理接口
- ✅ 任务注册表机制
- ✅ 自动错误恢复（panic 保护）
- ✅ 详细的日志记录
- ✅ 优雅启动和关闭
- ✅ 内置示例任务
- ✅ 完整的 API 文档
- ✅ HTTP 测试文件
- ✅ 集成到服务器生命周期

---

## 📝 下一步

1. **运行项目**：启动服务器查看示例任务执行
2. **查看日志**：检查 `logs/other/scheduler/` 目录
3. **测试 API**：使用 `scheduler.http` 文件测试
4. **开发任务**：参考 `tasks.go` 实现你的业务任务
5. **阅读文档**：查看 [README.md](./README.md) 了解详细用法

---

## 💡 提示

- 所有任务执行都有详细的日志记录
- 任务失败不会影响其他任务
- 可以通过 API 动态管理任务，无需重启服务器
- 使用任务注册表可以方便地在各模块中注册任务
- 示例任务可以在 `routes.go` 中修改或删除

---

## 📚 更多文档

- [详细使用说明](./README.md)
- [API 文档](./api.md)
- [HTTP 测试](./scheduler.http)
