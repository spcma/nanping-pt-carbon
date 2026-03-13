# NPFS DDD 架构说明

## 项目结构

本项目已改造为标准的 DDD (领域驱动设计) 四层架构：

```
internal/
├── domain/                 # 领域层
│   ├── entity.go          # 领域实体（File, Directory, NpfsSession）
│   ├── value_object.go    # 值对象（FilePath, FileName 等）
│   ├── repository.go      # 仓储接口定义
│   ├── service.go         # 领域服务接口
│   └── error.go           # 领域错误定义
│
├── application/           # 应用层
│   ├── dto.go            # 数据传输对象（DTO）
│   └── file_service.go   # 应用服务（协调领域对象完成业务逻辑）
│
├── infrastructure/        # 基础设施层
│   ├── file_repository.go     # 文件仓储实现（NPFS 具体实现）
│   ├── session_repository.go  # 会话仓储实现
│   └── fs_client.go           # NPFS 客户端工具
│
└── interfaces/            # 接口层
    └── http/
        ├── handler.go         # HTTP 处理器（处理 HTTP 请求）
        ├── router.go          # Gin 路由配置
        └── middleware/
            └── cors.go        # 中间件（CORS, Logger）
```

## 各层职责

### 1. Domain Layer（领域层）
- **位置**: `internal/domain/`
- **职责**: 包含核心业务逻辑和领域模型
- **特点**: 
  - 不依赖任何其他层
  - 定义实体、值对象、仓储接口、领域服务
  - 最纯粹的业务逻辑

**核心实体**:
- `File`: 文件实体
- `Directory`: 目录实体
- `NpfsSession`: NPFS 会话实体

**仓储接口**:
- `FileRepository`: 文件操作接口
- `SessionRepository`: 会话管理接口

### 2. Application Layer（应用层）
- **位置**: `internal/application/`
- **职责**: 协调领域对象完成具体的业务任务
- **特点**:
  - 依赖 Domain 层
  - 定义 DTO（数据传输对象）
  - 不包含核心业务逻辑，只负责流程编排

**应用服务**:
- `FileApplicationService`: 文件应用服务
  - `CheckDirExists()`: 检查目录是否存在
  - `CreateDirectory()`: 创建目录
  - `ListDirectory()`: 列出目录内容
  - `DeleteFile()`: 删除文件
  - `ReadFile()`: 读取文件
  - `SaveContent()`: 保存内容到文件
  - `SaveLocalFile()`: 保存本地文件到 NPFS

### 3. Infrastructure Layer（基础设施层）
- **位置**: `internal/infrastructure/`
- **职责**: 提供技术实现细节
- **特点**:
  - 实现 Domain 层定义的接口
  - 与外部系统交互（NPFS、数据库等）
  - 提供通用工具函数

**仓储实现**:
- `NpfsFileRepository`: 基于 NPFS 的文件仓储
- `NpfsSessionRepository`: NPFS 会话仓储

**工具函数**:
- `CreateFsClient()`: 创建 NPFS 客户端
- `ReadLocalFile()`: 读取本地文件

### 4. Interfaces Layer（接口层）
- **位置**: `internal/interfaces/`
- **职责**: 对外暴露的接口适配器
- **特点**:
  - 依赖 Application 层
  - 将外部请求转换为应用层能理解的形式
  - HTTP API、WebSocket 等

**HTTP Handler**:
- `FileHandler`: 文件操作 HTTP 处理器
  - `CheckDirHandler()`: 检查目录
  - `CreateDirHandler()`: 创建目录
  - `ListDirHandler()`: 列出目录
  - `DeleteDirHandler()`: 删除目录
  - `ReadFileHandler()`: 读取文件
  - `SaveFileHandler()`: 保存文件
  - `UploadFileHandler()`: 上传文件
  - `DownloadFileHandler()`: 下载文件
  - `DeleteFileHandler()`: 删除文件

**路由配置**:
- `/health`: 健康检查
- `/api/v1/dir/*`: 目录操作
- `/api/v1/file/*`: 文件操作

## 依赖关系

```
Interfaces → Application → Domain ← Infrastructure
                        ↑              ↓
                        └──────────────┘
```

- **箭头方向**: 表示依赖方向
- **核心原则**: 依赖倒置，Domain 层不依赖任何其他层
- **Infrastructure**: 同时依赖 Domain 并实现 Domain 的接口

## 数据流示例

以"保存文件"为例：

1. **HTTP 请求** → `interfaces/http/handler.go:SaveFileHandler()`
2. **解析请求** → 转换为 `application.SaveFileRequest`
3. **调用应用服务** → `application.FileApplicationService.SaveContent()`
4. **调用仓储接口** → `domain.FileRepository.SaveContent()`
5. **基础设施实现** → `infrastructure.NpfsFileRepository.SaveContent()`
6. **NPFS 操作** → 实际保存到 NPFS 文件系统
7. **返回结果** → 逐层返回到 HTTP 响应

## API 接口

### 健康检查
```bash
GET /health
```

### 目录操作
```bash
POST /api/v1/dir/check          # 检查目录
POST /api/v1/dir/create         # 创建目录
GET  /api/v1/dir/list?path=xxx  # 列出目录
DELETE /api/v1/dir/delete       # 删除目录
```

### 文件操作
```bash
GET    /api/v1/file/read?path=xxx      # 读取文件
POST   /api/v1/file/save               # 保存文件（从内容）
POST   /api/v1/file/upload             # 上传文件（multipart/form-data）
GET    /api/v1/file/download?path=xxx  # 下载文件
DELETE /api/v1/file/delete?path=xxx    # 删除文件
```

## 启动服务器

### 方式 1: 直接运行
```bash
go run main.go -mode http -addr :19870
```

### 方式 2: 构建后运行
```bash
go build -o nanping-pt-carbon.exe
.\nanping-pt-carbon.exe -mode http -addr :19870
```

### 方式 3: 使用脚本
```bash
.\scripts\start_ddd_http_server.bat
```

## 测试 API

使用提供的测试脚本：
```bash
.\scripts\test_ddd_http_api.bat
```

或使用 Postman/curl 手动测试：
```bash
# 健康检查
curl http://localhost:19870/health

# 检查目录
curl -X POST http://localhost:19870/api/v1/dir/check \
     -H "Content-Type: application/json" \
     -d '{"path":"/test"}'

# 创建目录
curl -X POST http://localhost:19870/api/v1/dir/create \
     -H "Content-Type: application/json" \
     -d '{"path":"/test","recursive":true}'

# 保存文件
curl -X POST http://localhost:19870/api/v1/file/save \
     -H "Content-Type: application/json" \
     -d '{"content":"Hello World","dir":"/test","filename":"hello.txt"}'
```

## DDD 优势

1. **清晰的职责分离**: 每层都有明确的职责
2. **可测试性**: 各层可以独立测试
3. **可维护性**: 代码组织清晰，易于理解和修改
4. **可扩展性**: 容易添加新功能而不影响现有代码
5. **业务聚焦**: Domain 层专注于核心业务逻辑
6. **技术无关**: Domain 层不依赖具体技术实现

## 后续改进建议

1. 添加命令查询责任分离（CQRS）模式
2. 引入领域事件（Domain Events）
3. 添加单元测试和集成测试
4. 实现仓储的缓存层
5. 添加认证授权中间件
6. 实现请求限流和防抖
