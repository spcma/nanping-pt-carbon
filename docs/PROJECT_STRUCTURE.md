# 项目架构说明（保留 Shared 包）

## 📁 完整的项目目录结构

```
internal/
├── domain/                 # 领域层 - 核心业务逻辑
│   ├── entity.go          # 领域实体（使用 shared/entity）
│   ├── value_object.go    # 值对象（复用 shared/valueobject）
│   ├── repository.go      # 仓储接口定义
│   ├── service.go         # 领域服务接口
│   └── error.go           # 领域错误（使用 shared/error）
│
├── application/           # 应用层 - 业务流程编排
│   ├── dto.go            # 数据传输对象
│   ├── file_service.go   # 文件应用服务
│   └── service_registry.go # 服务注册表
│
├── infrastructure/        # 基础设施层 - 技术实现
│   ├── file_repository.go     # NPFS 文件仓储实现
│   ├── session_repository.go  # 会话仓储实现
│   └── fs_client.go           # NPFS 客户端工具
│
├── interfaces/            # 接口层 - 对外接口适配器
│   └── http/
│       ├── handler.go         # HTTP 处理器
│       ├── router.go          # Gin 路由配置
│       └── middleware/
│           └── cors.go        # 中间件
│
├── shared/                # 公共功能包 - 跨层复用的通用组件 ⭐
│   ├── init.go            # 全局初始化（GlobalDB, GlobalEventBus）
│   │
│   ├── entity/            # 通用实体基类
│   │   ├── base_entity.go     # 基础实体（审计字段）
│   │   └── pagination.go      # 分页功能
│   │
│   ├── valueobject/       # 通用值对象
│   │   ├── coordinate.go      # 坐标值对象
│   │   ├── license.go         # 许可证值对象
│   │   └── time_range.go      # 时间范围值对象
│   │
│   ├── error/             # 错误处理
│   │   └── error.go         # 统一错误结构
│   │
│   ├── db/                # 数据库访问（GORM）
│   │   ├── config.go
│   │   └── gorm.go
│   │
│   ├── cache/             # 缓存（Redis）
│   │   ├── cache.go
│   │   └── redis.go
│   │
│   ├── event/             # 事件定义
│   │   └── ...
│   │
│   ├── eventbus/          # 事件总线
│   │   └── ...
│   │
│   ├── http/              # HTTP 工具
│   │   └── route_registry.go  # 路由注册配置
│   │
│   ├── logger/            # 日志系统
│   │   ├── logger.go
│   │   ├── middleware.go
│   │   └── ...
│   │
│   ├── crypto/            # 加密工具
│   ├── idgen/             # ID 生成器
│   ├── token/             # Token 管理
│   ├── validator/         # 验证器
│   ├── transaction/       # 事务管理
│   ├── persistence/       # 持久化工具
│   └── timeutil/          # 时间工具
│
├── rpc/                   # RPC 客户端（NPFS 专用）
│   ├── backend.go
│   ├── lapi.go
│   ├── mdb.go
│   ├── rpc.go
│   ├── stub.go
│   └── vars.go
│
├── paltform/              # 平台相关代码
├── service/               # 旧版服务代码（可逐步迁移）
└── transport/             # 传输层
```

## 🎯 Shared 包的职责

### Shared 包是**跨层复用的公共基础设施**，提供：

1. **通用抽象** (`entity`, `valueobject`, `error`)
   - 所有领域实体都继承自 `shared/entity.BaseEntity`
   - 使用 `shared/error.Error` 统一定义错误
   - 复用 `shared/valueobject` 中的值对象

2. **技术实现** (`db`, `cache`, `logger`)
   - GORM 数据库连接和配置
   - Redis 缓存客户端
   - 统一的日志系统

3. **横切关注点** (`event`, `eventbus`, `transaction`)
   - 事件驱动架构支持
   - 事务管理
   - 中间件支持

4. **工具函数** (`crypto`, `idgen`, `token`, `validator`, `timeutil`)
   - 加密解密
   - ID 生成
   - Token 验证
   - 数据验证
   - 时间处理

## 🔄 Shared 包与 DDD 各层的关系

```
┌─────────────────────────────────────────┐
│         Interfaces Layer                │
│  (interfaces/http)                      │
│  - Handler                              │
│  - Router (使用 shared/http)            │
│  - Middleware (使用 shared/logger)      │
└────────────┬────────────────────────────┘
             │ 调用
             ↓
┌─────────────────────────────────────────┐
│         Application Layer               │
│  (application)                          │
│  - Application Services                 │
│  - DTOs                                 │
│  - 使用 shared/eventbus 发布事件        │
└────────────┬────────────────────────────┘
             │ 调用
             ↓
┌─────────────────────────────────────────┐
│         Domain Layer                    │
│  (domain)                               │
│  - Entities (继承 shared/entity)        │
│  - Value Objects (复用 shared/vo)       │
│  - Repository Interfaces                │
│  - Domain Errors (使用 shared/error)    │
└────────────┬────────────────────────────┘
             │ 依赖倒置
             ↓
┌─────────────────────────────────────────┐
│      Infrastructure Layer               │
│  (infrastructure)                       │
│  - Repository Implementations           │
│  - 使用 shared/db 进行数据持久化        │
│  - 使用 shared/cache 进行缓存           │
│  - 使用 shared/rpc 调用外部服务         │
└─────────────────────────────────────────┘
             ↑
             │ 使用
             │
┌─────────────────────────────────────────┐
│         Shared Package                  │
│  (internal/shared)                      │
│  - entity/, valueobject/, error/        │
│  - db/, cache/, logger/                 │
│  - event/, eventbus/                    │
│  - crypto/, idgen/, token/, ...         │
└─────────────────────────────────────────┘
```

## 📦 Shared 包的使用规范

### ✅ 正确使用方式

1. **Domain 层** 可以使用：
   - `shared/entity` - 继承 BaseEntity
   - `shared/valueobject` - 复用值对象
   - `shared/error` - 定义领域错误

2. **Application 层** 可以使用：
   - `shared/eventbus` - 发布领域事件
   - `shared/logger` - 记录日志
   - `shared/validator` - 验证数据

3. **Infrastructure 层** 可以使用：
   - `shared/db` - 数据库操作
   - `shared/cache` - 缓存操作
   - `shared/logger` - 记录日志
   - `shared/transaction` - 事务管理

4. **Interfaces 层** 可以使用：
   - `shared/http` - 路由注册
   - `shared/logger` - 日志中间件
   - `shared/token` - 认证中间件

### ❌ 错误使用方式

1. **Domain 层不应该依赖**：
   - `shared/db` - 会破坏领域层的纯粹性
   - `shared/cache` - 技术细节不应侵入领域层
   - `shared/logger` - 领域层应该专注于业务逻辑

2. **Shared 包不应该依赖**：
   - Domain/Application/Infrastructure/Interfaces 层
   - Shared 包应该是最低层的通用组件

## 🔧 关键集成点

### 1. 实体继承 shared/entity.BaseEntity

```go
// domain/entity.go
type File struct {
    entity.BaseEntity  // 继承审计字段
    
    Path     string `json:"path"`
    Name     string `json:"name"`
    Size     int64  `json:"size"`
    Hash     string `json:"hash"`
    IsDir    bool   `json:"is_dir"`
}
```

### 2. 错误使用 shared/error.Error

```go
// domain/error.go
var ErrFileNotFound = sharederror.NewError(
    "FILE_NOT_FOUND", 
    "文件不存在",
)
```

### 3. 应用服务使用 shared/eventbus

```go
// application/file_service.go
func (s *FileApplicationService) SaveContent(...) {
    // ... 业务逻辑
    
    // 发布事件
    s.eventBus.Publish("file.saved", FileSavedEvent{...})
}
```

### 4. 基础设施使用 shared/db

```go
// infrastructure/file_repository.go
type NpfsFileRepository struct {
    rpcClient *rpc.LApiStub
    db        *gorm.DB  // 可选的数据库支持
}

func (r *NpfsFileRepository) SaveToFileDB(file *domain.File) error {
    return r.db.Create(file).Error
}
```

## 🚀 启动流程

```go
// main.go
func main() {
    // 1. 初始化 shared 包
    db, err := shared.Init(persistence.Config{...})
    
    // 2. 初始化基础设施
    client, session := infrastructure.CreateFsClient()
    fileRepo := infrastructure.NewNpfsFileRepository(client, session, db)
    
    // 3. 初始化应用层
    fileAppService := application.NewFileApplicationService(fileRepo)
    
    // 4. 设置接口层
    router := http.SetupRouter(fileAppService)
    
    // 5. 启动服务
    router.Run(":19870")
}
```

## 📊 代码统计

| 层级 | 包路径 | 主要功能 | 依赖 Shared 模块 |
|------|--------|----------|------------------|
| **Domain** | `domain` | 核心业务逻辑 | entity, valueobject, error |
| **Application** | `application` | 流程编排 | eventbus, logger, validator |
| **Infrastructure** | `infrastructure` | 技术实现 | db, cache, logger |
| **Interfaces** | `interfaces/http` | HTTP 接口 | http, logger, token |
| **Shared** | `shared` | 公共组件 | - |

## 🎯 架构优势

1. **代码复用**: Shared 包提供跨层复用的组件，避免重复造轮子
2. **统一标准**: 统一的实体基类、错误处理、日志格式
3. **易于扩展**: 新增功能时可以直接使用 Shared 包的成熟组件
4. **清晰分层**: DDD 分层明确，Shared 包作为底层支撑
5. **技术隔离**: 技术细节封装在 Shared 和 Infrastructure 层

## ⚠️ 注意事项

1. **避免循环依赖**: Shared 包不能依赖任何上层包
2. **保持领域纯粹性**: Domain 层只使用 Shared 的基础抽象，不使用技术实现
3. **适度使用**: 不要过度设计，简单场景可以直接使用 Shared 包
4. **文档化**: Shared 包应该有完善的文档和示例

## 📚 下一步优化建议

1. [ ] 将 `service/` 目录的旧代码逐步迁移到新的 DDD 结构
2. [ ] 为 Shared 包添加完整的单元测试
3. [ ] 使用 `shared/http.RouteRegistry` 重构路由注册
4. [ ] 集成 `shared/token` 实现 JWT 认证
5. [ ] 使用 `shared/validator` 统一参数验证
6. [ ] 使用 `shared/logger` 替换所有 fmt.Println
