# NPFS HTTP 服务快速开始

## 🚀 快速启动

### 方式一：使用批处理脚本（推荐）
```bash
start_http_server.bat
```

### 方式二：命令行
```bash
go run . -http-server
```

服务将运行在 `http://localhost:8080`

---

## 📝 功能列表

✅ **目录操作**
- 检查目录是否存在
- 创建目录（支持递归）
- 列出目录内容
- 删除目录

✅ **文件操作**
- 读取文件内容
- 保存文本内容到文件
- 上传本地文件
- 下载文件到本地
- 删除文件

---

## 🔧 使用方式

### 1. Web 界面（最简单）

打开浏览器访问 `np_fs_web.html` 文件，提供图形化操作界面。

### 2. REST API

使用 curl、Postman 或其他 HTTP 客户端调用 API。

**示例：**
```bash
# 创建目录
curl -X POST http://localhost:8080/api/v1/dir/create \
  -H "Content-Type: application/json" \
  -d '{"path":"/test","recursive":true}'

# 保存文件
curl -X POST http://localhost:8080/api/v1/file/save \
  -H "Content-Type: application/json" \
  -d '{"content":"Hello World","dir":"/test","filename":"hello.txt"}'

# 读取文件
curl http://localhost:8080/api/v1/file/read?path=/test/hello.txt

# 列出目录
curl http://localhost:8080/api/v1/dir/list?path=/test
```

### 3. Go 客户端

参考 `np_fs_http_example.go` 文件中的完整示例代码。

---

## 🧪 测试

运行快速测试脚本：
```bash
test_http_api.bat
```

---

## 📚 文档

详细 API 文档请查看：`NPFS_HTTP_API.md`

---

## 🛠️ 技术栈

- **Web 框架**: Gin v1.12.0
- **Go 版本**: 1.25.0
- **底层通信**: HPROSE RPC
- **文件系统**: NPFS

---

## ⚠️ 注意事项

1. 确保 NPFS 后端服务已在运行（默认端口 4800）
2. HTTP 服务默认监听 8080 端口
3. 所有路径都以 `/` 开头，例如 `/my_files/test.txt`

---

## 📞 常见问题

**Q: 无法启动服务？**
A: 检查端口是否被占用，可以修改 `-http-addr` 参数

**Q: API 返回错误？**
A: 确保 NPFS 后端服务正常运行

**Q: 如何修改端口？**
A: 使用 `-http-addr :9000` 参数指定其他端口

---

## 📋 API 端点概览

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /health | 健康检查 |
| POST | /api/v1/dir/check | 检查目录 |
| POST | /api/v1/dir/create | 创建目录 |
| GET | /api/v1/dir/list | 列出目录 |
| DELETE | /api/v1/dir/delete | 删除目录 |
| GET | /api/v1/file/read | 读取文件 |
| POST | /api/v1/file/save | 保存文件 |
| POST | /api/v1/file/upload | 上传文件 |
| GET | /api/v1/file/download | 下载文件 |
| DELETE | /api/v1/file/delete | 删除文件 |

---

## 🎯 下一步

1. 查看 `NPFS_HTTP_API.md` 了解完整 API 文档
2. 使用 `np_fs_web.html` Web 界面体验功能
3. 参考 `np_fs_http_example.go` 编写自己的客户端
