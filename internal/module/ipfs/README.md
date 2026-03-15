# IPFS 模块

简化的 IPFS 文件系统模块，提供文件和目录操作功能。

## 目录结构

```
internal/module/ipfs/
├── ipfs.go      # 核心服务实现
└── routes.go    # 路由注册
```

## 功能特性

- ✅ 不使用 DDD 架构，所有代码在一个包中
- ✅ 简洁的服务层设计
- ✅ 直接的路由注册方式
- ✅ 支持 IPFS 文件系统操作

## API 接口

### 目录操作

- `POST /api/v1/dir/check` - 检查目录
- `POST /api/v1/dir/create` - 创建目录
- `GET /api/v1/dir/list` - 列出目录内容
- `DELETE /api/v1/dir/delete` - 删除目录

### 文件操作

- `GET /api/v1/file/read` - 读取文件
- `POST /api/v1/file/save` - 保存文件
- `POST /api/v1/file/upload` - 上传文件
- `GET /api/v1/file/download` - 下载文件
- `DELETE /api/v1/file/delete` - 删除文件

## 使用示例

```go
import "app/internal/module/ipfs"

// 在 router.go 中注册路由
func registerRoutes() {
    client, session, err := ipfs.CreateFsClient()
    if err != nil {
        // 错误处理
        return
    }
    
    routes := ipfs.RegisterRoutes(client, session)
    // 注册到路由器
}
```

## 待实现功能

当前实现提供了基础框架，具体的 IPFS 操作逻辑需要在以下方法中实现：

- `CheckDir` - IPFS 目录检查
- `CreateDir` - IPFS 创建目录
- `ListDir` - IPFS 列出目录
- `DeleteFile` - IPFS 删除文件
- `ReadFile` - IPFS 读取文件
- `SaveFile` - IPFS 保存文件
- `UploadFile` - IPFS 上传文件
- `DownloadFile` - IPFS 下载文件
