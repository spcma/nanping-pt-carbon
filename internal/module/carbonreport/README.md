# CarbonReport 模块使用说明

## DDD 架构概述

CarbonReport 模块已按照 DDD（领域驱动设计）架构重构，分为以下层次：

### 目录结构

```
internal/module/carbonreport/
├── application/          # 应用层
│   ├── dto.go           # 数据传输对象
│   └── file_app_service.go  # 应用服务（业务逻辑编排）
├── domain/              # 领域层
│   ├── entity.go        # 实体定义
│   ├── repository.go    # 仓储接口
│   ├── service.go       # 领域服务接口
│   ├── value_object.go  # 值对象
│   ├── error.go         # 领域错误
│   ├── file_service_impl.go    # 文件领域服务实现
│   └── upload_service_impl.go  # 上传领域服务实现
├── persistence/         # 基础设施层（持久化）
│   ├── file_repository.go      # 文件仓储实现
│   └── session_repository.go   # 会话仓储实现
├── transport/           # 接口层
│   └── http/
│       ├── file_handler.go     # HTTP 处理器
│       ├── routes.go           # 路由配置
│       └── module.go           # 模块入口
└── wire/              # 依赖注入
    └── wire.go        # DDD 组件初始化
```

## 使用示例

### 1. 初始化模块

```go
import (
    "app/internal/module/carbonreport/transport/http"
    "app/internal/rpc"
)

// 创建客户端和会话
client, session, err := domain.CreateFsClient()
if err != nil {
    panic(err)
}

// 初始化模块
module, err := http.NewModule(client, session)
if err != nil {
    panic(err)
}
defer module.Close()
```

### 2. 注册路由

```go
// 在主路由管理器中注册 carbonreport 模块
router.RegisterModule(module.Handlers)
```

### 3. 使用应用服务

```go
import (
    "app/internal/module/carbonreport/application"
    "context"
)

// 获取应用服务
appService := module.DDD.AppService

// 检查目录
checkResult, err := appService.CheckDir(context.Background(), application.CheckDirCommand{
    Path: "/my/path",
})

// 创建目录
_, err = appService.CreateDir(context.Background(), application.CreateDirCommand{
    Path:      "/my/new/dir",
    Recursive: true,
})

// 列出目录
listResult, err := appService.ListDir(context.Background(), application.ListDirCommand{
    Path: "/my/path",
})

// 保存文件
saveResult, err := appService.SaveFile(context.Background(), application.SaveFileCommand{
    Content:  "file content",
    Dir:      "/my/dir",
    Filename: "test.txt",
})
```

## API 端点

所有 API 都需要认证（AuthTypeRequired）：

- `POST /api/v1/dir/check` - 检查目录
- `POST /api/v1/dir/create` - 创建目录
- `GET /api/v1/dir/list` - 列出目录
- `DELETE /api/v1/dir/delete` - 删除目录
- `GET /api/v1/file/read` - 读取文件
- `POST /api/v1/file/save` - 保存文件
- `POST /api/v1/file/upload` - 上传文件
- `GET /api/v1/file/download` - 下载文件
- `DELETE /api/v1/file/delete` - 删除文件

## 架构说明

### 各层职责

1. **Domain Layer（领域层）**
   - 定义实体（Entity）、值对象（Value Object）
   - 定义仓储接口（Repository Interface）
   - 定义领域服务接口和实现
   - 包含核心业务逻辑

2. **Application Layer（应用层）**
   - 定义应用服务（App Service）
   - 编排领域对象完成业务流程
   - 定义命令和查询对象（Command/Query）
   - 定义 DTO（数据传输对象）

3. **Persistence Layer（持久化层）**
   - 实现仓储接口
   - 处理数据库/文件系统操作
   - 依赖注入时传入 Domain 层定义的接口

4. **Transport Layer（接口层）**
   - HTTP Handler 处理请求
   - 参数绑定和验证
   - 调用应用服务
   - 返回响应

5. **Wire（依赖注入）**
   - 组装各层组件
   - 管理依赖关系

### 数据流向

```
HTTP Request → Handler → AppService → Domain Service → Repository → Persistence
                ↑                                                            ↓
             Response ← DTO ← AppService ← Domain Entity ← Data
```

## 与 IAM 模块对比

CarbonReport 模块完全遵循 IAM 模块的 DDD 架构模式：

| 层级 | IAM 模块 | CarbonReport 模块 |
|------|---------|------------------|
| Application | SysRoleAppService | FileAppService |
| Domain | SysRole (实体) | File (实体) |
| Persistence | SysRoleRepository | FileRepository |
| Transport | SysRoleHandler | FileHandler |
| Wire | InitSysRoleDDD | InitFileDDD |

## 扩展指南

### 添加新的领域对象

1. 在 `domain/entity.go` 中定义实体
2. 在 `domain/repository.go` 中定义仓储接口
3. 在 `persistence/` 中实现仓储
4. 在 `domain/service.go` 中定义领域服务（如需要）
5. 在 `application/` 中定义应用服务和 DTO

### 添加新的 API 端点

1. 在 `application/file_app_service.go` 中添加应用方法
2. 在 `transport/http/file_handler.go` 中添加 Handler
3. 在 `transport/http/routes.go` 中注册路由

## 注意事项

1. **Context 传递**：所有仓储和应用服务方法都接受 `context.Context` 参数
2. **错误处理**：使用领域层定义的 errors（`domain/error.go`）
3. **事务管理**：如需事务，使用 `shared/transaction` 包
4. **日志记录**：使用 `shared/logger` 包进行日志记录
5. **依赖注入**：通过 `wire/wire.go` 统一管理依赖

## 迁移说明

旧的 `service/` 和 `infrastructure/` 目录已重命名为 `service_backup/` 和 `infrastructure_backup/`，可以参考但不应再使用。

新代码应该完全基于 DDD 架构编写。
