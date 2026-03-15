# DDD 架构快速参考

## 📁 目录结构速查

```
internal/
├── domain/              ← 核心业务逻辑（实体、接口定义）
├── application/         ← 应用服务（业务流程编排）
├── infrastructure/      ← 技术实现（仓储实现、NPFS 客户端）
└── interfaces/http/     ← HTTP 接口（Gin 路由、Handler）
```

## 🎯 各层职责

| 层级 | 包名 | 职责 | 示例 |
|------|------|------|------|
| **Domain** | `domain` | 核心业务逻辑 | Entity, Repository 接口 |
| **Application** | `application` | 流程编排 | FileApplicationService |
| **Infrastructure** | `infrastructure` | 技术实现 | NpfsFileRepository |
| **Interfaces** | `http` | 对外接口 | FileHandler |

## 🔄 调用流程

```
HTTP Request 
  ↓
Interface (handler.go)
  ↓
Application (file_service.go)
  ↓
Domain (repository interface)
  ↓
Infrastructure (repository impl)
  ↓
NPFS API
```

## 📡 API 端点

### 健康检查
- `GET /health`

### 目录操作
- `POST /api/v1/dir/check` - 检查目录
- `POST /api/v1/dir/create` - 创建目录
- `GET /api/v1/dir/list` - 列出目录
- `DELETE /api/v1/dir/delete` - 删除目录

### 文件操作
- `GET /api/v1/file/read` - 读取文件
- `POST /api/v1/file/save` - 保存文件
- `POST /api/v1/file/upload` - 上传文件
- `GET /api/v1/file/download` - 下载文件
- `DELETE /api/v1/file/delete` - 删除文件

## 💻 常用命令

### 构建
```bash
go build -o nanping-pt-carbon.exe
```

### 运行
```bash
go run main.go -mode http -addr :19870
```

### 测试
```bash
curl http://localhost:19870/health
```

## 🔑 关键代码位置

| 功能 | 文件路径 |
|------|----------|
| HTTP 路由配置 | `internal/interfaces/http/router.go` |
| HTTP Handler | `internal/interfaces/http/handler.go` |
| 应用服务 | `internal/application/file_service.go` |
| DTO 定义 | `internal/application/dto.go` |
| 仓储实现 | `internal/infrastructure/file_repository.go` |
| 领域实体 | `internal/domain/entity.go` |
| 主入口 | `main.go` |

## 📝 请求示例

### 检查目录
```bash
curl -X POST http://localhost:19870/api/v1/dir/check \
  -H "Content-Type: application/json" \
  -d '{"path":"/test"}'
```

### 创建目录
```bash
curl -X POST http://localhost:19870/api/v1/dir/create \
  -H "Content-Type: application/json" \
  -d '{"path":"/test","recursive":true}'
```

### 保存文件
```bash
curl -X POST http://localhost:19870/api/v1/file/save \
  -H "Content-Type: application/json" \
  -d '{"content":"Hello World","dir":"/test","filename":"hello.txt"}'
```

### 上传文件
```bash
curl -X POST http://localhost:19870/api/v1/file/upload \
  -F "file=@localfile.txt" \
  -F "dir=/upload" \
  -F "filename=test.txt"
```

## 🏗️ DDD 核心概念

### Entity（实体）
- 有唯一标识的对象
- 例：`File`, `Directory`, `NpfsSession`

### Value Object（值对象）
- 无唯一标识，通过属性值判断相等
- 例：`FilePath`, `FileName`, `IPFSID`

### Repository（仓储）
- 封装数据访问逻辑
- 接口在 Domain 层定义
- 实现在 Infrastructure 层

### Service（服务）
- Application Service: 协调领域对象
- Domain Service: 纯业务逻辑
- Interface Service: 对外接口适配

## ✅ 改造完成清单

- [x] Domain 层：实体、值对象、仓储接口、领域服务
- [x] Application 层：应用服务、DTO
- [x] Infrastructure 层：仓储实现、NPFS 客户端
- [x] Interfaces 层：HTTP Handler、Gin 路由、中间件
- [x] main.go：依赖注入、启动 HTTP 服务器
- [x] Gin 框架集成
- [x] CORS 中间件
- [x] Logger 中间件
- [x] 统一响应格式
- [x] 错误处理
- [x] 启动脚本
- [x] 测试脚本
- [x] 架构文档
