# Project Module - 项目管理模块

## 概述

Project 模块是基於 DDD（领域驱动设计）架构的项目管理功能模块，提供完整的项目 CRUD 操作和状态管理。

## 目录结构

```
internal/module/project/
├── application/              # 应用层
│   └── project_app_service.go    # 应用服务
├── domain/                   # 领域层
│   ├── project.go                # 项目聚合根
│   └── repository.go             # 仓储接口
├── persistence/              # 基础设施层
│   └── project_repository.go     # 仓储实现
├── transport/http/           # 接口层
│   ├── module.go                 # 模块定义
│   ├── project_handler.go        # HTTP 处理器
│   └── routes.go                 # 路由配置
├── wire/                     # 依赖注入
│   └── wire.go                   # Wire 配置
└── service_backup/           # 测试和示例
    └── api.http                  # API 测试文件
```

## 功能特性

### 1. 项目实体

- **项目名称** (name): 项目的显示名称
- **项目编码** (code): 项目的唯一标识符
- **项目状态** (status): 
  - `0` - 待启动 (Pending)
  - `1` - 进行中 (Active)
  - `2` - 已完成 (Completed)
  - `3` - 已取消 (Cancelled)
- **项目描述** (description): 项目的详细描述

### 2. API 端点

所有 API 都需要 JWT Token 认证。

#### 创建项目
```http
POST /api/project
Content-Type: application/json
Authorization: Bearer {{token}}

{
  "name": "项目名称",
  "code": "PROJECT_CODE",
  "description": "项目描述"
}
```

#### 获取项目列表
```http
GET /api/projects?pageNum=1&pageSize=10
Authorization: Bearer {{token}}
```

查询参数：
- `pageNum`: 页码（默认 1）
- `pageSize`: 每页数量（默认 10，最大 100）
- `name`: 按项目名称模糊查询
- `code`: 按项目编码精确查询
- `status`: 按状态筛选
- `sortBy`: 排序字段
- `sortOrder`: 排序方式（asc/desc）

#### 根据 ID 获取项目
```http
GET /api/project/:id
Authorization: Bearer {{token}}
```

#### 根据编码获取项目
```http
GET /api/project/code/:code
Authorization: Bearer {{token}}
```

#### 更新项目
```http
PUT /api/project/:id
Content-Type: application/json
Authorization: Bearer {{token}}

{
  "name": "更新后的名称",
  "description": "更新后的描述"
}
```

#### 变更项目状态
```http
PUT /api/project/:id/status
Content-Type: application/json
Authorization: Bearer {{token}}

{
  "status": "2"
}
```

#### 删除项目
```http
DELETE /api/project/:id
Authorization: Bearer {{token}}
```

## 使用示例

### 1. 在代码中使用

```go
import (
    "app/internal/module/project/application"
    "app/internal/module/project/domain"
)

// 创建项目
cmd := application.CreateProjectCommand{
    Name:        "新项目",
    Code:        "NEW001",
    Description: "这是一个新项目",
    UserID:      currentUserId,
}

projectId, err := projectAppService.CreateProject(ctx, cmd)

// 获取项目
project, err := projectAppService.GetProjectByID(ctx, projectId)

// 更新项目
updateCmd := application.UpdateProjectCommand{
    ID:          projectId,
    Name:        "更新后的名称",
    Description: "更新后的描述",
    UserID:      currentUserId,
}
err = projectAppService.UpdateProject(ctx, updateCmd)

// 变更状态
statusCmd := application.ChangeProjectStatusCommand{
    ID:     projectId,
    Status: domain.ProjectStatusCompleted,
    UserID: currentUserId,
}
err = projectAppService.ChangeProjectStatus(ctx, statusCmd)

// 删除项目
err = projectAppService.DeleteProject(ctx, projectId, currentUserId)
```

### 2. 使用测试脚本

运行测试脚本前，请确保：
1. 服务器已启动
2. 已获取有效的 JWT Token

```bash
# 编辑 test_project_api.bat，设置正确的 TOKEN
# 然后运行
scripts\test_project_api.bat
```

### 3. 使用 HTTP 测试文件

使用 VS Code 的 REST Client 插件或其他 HTTP 客户端工具打开：
```
internal/module/project/service_backup/api.http
```

## 数据库表结构

```sql
CREATE TABLE project (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    code        VARCHAR(100) NOT NULL UNIQUE,
    status      VARCHAR(10) NOT NULL DEFAULT '1',
    description TEXT,
    create_by   BIGINT,
    create_time TIMESTAMP NOT NULL,
    update_by   BIGINT,
    update_time TIMESTAMP,
    delete_by   BIGINT,
    delete_time TIMESTAMP
);

-- 索引
CREATE INDEX idx_project_code ON project(code);
CREATE INDEX idx_project_status ON project(status);
CREATE INDEX idx_project_name ON project(name);
```

## 扩展开发

### 添加新的业务逻辑

1. 在 `domain/project.go` 中添加领域方法
2. 在 `application/project_app_service.go` 中添加应用服务方法
3. 在 `transport/http/project_handler.go` 中添加 HTTP 处理器
4. 在 `transport/http/routes.go` 中注册新路由

### 自定义仓储方法

1. 在 `domain/repository.go` 中扩展接口
2. 在 `persistence/project_repository.go` 中实现方法

## 注意事项

1. **逻辑删除**: 项目采用逻辑删除，删除时只更新 `delete_by` 和 `delete_time` 字段
2. **状态管理**: 项目状态变更是通过领域方法实现的，确保业务规则一致性
3. **审计字段**: 自动记录创建人、创建时间、更新人、更新时间
4. **分页查询**: 默认返回脱敏数据，可根据需要扩展公开字段

## 故障排除

### 常见问题

1. **Token 无效**: 确保使用有效的 JWT Token
2. **项目编码重复**: 项目编码必须唯一
3. **状态值无效**: 只能使用预定义的状态值（0/1/2/3）

### 日志查看

查看应用日志以获取更多错误信息：
```bash
# 日志文件位置取决于 config.yaml 配置
tail -f logs/app.log
```

## 参考文档

- [DDD 架构说明](../../../docs/DDD_ARCHITECTURE.md)
- [HTTP API 文档](../../../docs/NPFS_HTTP_API.md)
- [IAM 模块示例](../iam/)
