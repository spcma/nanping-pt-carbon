# NPFS HTTP API 使用文档

## 概述

基于 Gin 框架实现的 NPFS 文件系统 HTTP RESTful API。

## 启动服务

### 方式一：使用启动脚本
```bash
start_http_server.bat
```

### 方式二：命令行启动
```bash
# 使用默认端口 :8080
go run . -http-server

# 指定端口
go run . -http-server -http-addr :9000

# 或使用 mode 参数
go run . -mode http -http-addr :8080
```

## API 端点

### 健康检查
```
GET /health
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "ok"
  }
}
```

---

### 目录操作

#### 1. 检查目录是否存在
```
POST /api/v1/dir/check
Content-Type: application/x-www-form-urlencoded

参数:
- path: 目录路径
```

**请求示例：**
```bash
curl -X POST http://localhost:8080/api/v1/dir/check \
  -d "path=/test"
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "path": "/test",
    "exists": true
  }
}
```

---

#### 2. 创建目录
```
POST /api/v1/dir/create
Content-Type: application/json

参数:
- path: 目录路径 (必填)
- recursive: 是否递归创建父目录 (可选，默认 false)
```

**请求示例：**
```bash
curl -X POST http://localhost:8080/api/v1/dir/create \
  -H "Content-Type: application/json" \
  -d '{"path":"/my/files","recursive":true}'
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "path": "/my/files",
    "created": true
  }
}
```

---

#### 3. 列出目录内容
```
GET /api/v1/dir/list?path={目录路径}
```

**请求示例：**
```bash
curl http://localhost:8080/api/v1/dir/list?path=/
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "path": "/",
    "files": [
      {
        "name": "test",
        "type": "directory",
        "size": 0
      },
      {
        "name": "hello.txt",
        "type": "file",
        "size": 12
      }
    ]
  }
}
```

---

#### 4. 删除目录
```
DELETE /api/v1/dir/delete?path={目录路径}&recursive={true|false}&force={true|false}
```

**请求示例：**
```bash
curl -X DELETE "http://localhost:8080/api/v1/dir/delete?path=/test&recursive=true&force=true"
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "path": "/test",
    "deleted": true
  }
}
```

---

### 文件操作

#### 1. 读取文件
```
GET /api/v1/file/read?path={文件路径}
```

**请求示例：**
```bash
curl http://localhost:8080/api/v1/file/read?path=/my_files/hello.txt
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "path": "/my_files/hello.txt",
    "size": 12,
    "content": "Hello World!"
  }
}
```

---

#### 2. 保存文件（从文本内容）
```
POST /api/v1/file/save
Content-Type: application/json

参数:
- content: 文件内容 (必填)
- dir: 目标目录 (必填)
- filename: 文件名 (必填)
```

**请求示例：**
```bash
curl -X POST http://localhost:8080/api/v1/file/save \
  -H "Content-Type: application/json" \
  -d '{"content":"Hello World!","dir":"/my_files","filename":"hello.txt"}'
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "ipfsid": "QmWfVY9y3xjsixTgbd9AorQxH7VtMpzfx2HaWtsoUYecaX",
    "path": "/my_files/hello.txt"
  }
}
```

---

#### 3. 上传文件（multipart/form-data）
```
POST /api/v1/file/upload
Content-Type: multipart/form-data

参数:
- file: 文件对象 (必填)
- dir: 目标目录 (必填)
- filename: 文件名 (可选，默认使用原文件名)
```

**请求示例：**
```bash
curl -X POST http://localhost:8080/api/v1/file/upload \
  -F "file=@./local.txt" \
  -F "dir=/uploads" \
  -F "filename=upload.txt"
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "ipfsid": "QmWfVY9y3xjsixTgbd9AorQxH7VtMpzfx2HaWtsoUYecaX",
    "path": "/uploads/upload.txt"
  }
}
```

---

#### 4. 下载文件
```
GET /api/v1/file/download?path={文件路径}&filename={下载文件名}
```

**请求示例：**
```bash
# 下载文件
curl -O http://localhost:8080/api/v1/file/download?path=/my_files/hello.txt

# 指定下载文件名
curl -L -o downloaded.txt "http://localhost:8080/api/v1/file/download?path=/my_files/hello.txt&filename=test.txt"
```

**响应：**
- Content-Type: application/octet-stream
- Content-Disposition: attachment; filename*=UTF-8''{filename}
- Content-Length: {文件大小}

---

#### 5. 删除文件
```
DELETE /api/v1/file/delete?path={文件路径}&force={true|false}
```

**请求示例：**
```bash
curl -X DELETE "http://localhost:8080/api/v1/file/delete?path=/my_files/hello.txt&force=true"
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "path": "/my_files/hello.txt",
    "deleted": true
  }
}
```

---

## 错误响应

所有 API 错误都使用统一的响应格式：

```json
{
  "code": {错误代码},
  "message": "{错误描述}",
  "data": null
}
```

**常见错误代码：**
- `0`: 成功
- `400`: 请求参数错误
- `500`: 服务器内部错误

---

## 使用示例

### Go 客户端示例

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func main() {
    // 1. 创建目录
    body := `{"path":"/test","recursive":true}`
    resp, _ := http.Post(
        "http://localhost:8080/api/v1/dir/create",
        "application/json",
        bytes.NewReader([]byte(body))
    )
    
    var result Response
    json.NewDecoder(resp.Body).Decode(&result)
    fmt.Printf("创建目录：%+v\n", result)
    
    // 2. 保存文件
    body = `{"content":"Hello","dir":"/test","filename":"hello.txt"}`
    resp, _ = http.Post(
        "http://localhost:8080/api/v1/file/save",
        "application/json",
        bytes.NewReader([]byte(body))
    )
    json.NewDecoder(resp.Body).Decode(&result)
    fmt.Printf("保存文件：%+v\n", result)
    
    // 3. 列出目录
    resp, _ = http.Get("http://localhost:8080/api/v1/dir/list?path=/test")
    json.NewDecoder(resp.Body).Decode(&result)
    fmt.Printf("目录列表：%+v\n", result)
}
```

### JavaScript 客户端示例

```javascript
const API_BASE = 'http://localhost:8080/api/v1';

// 创建目录
async function createDir(path, recursive = true) {
    const response = await fetch(`${API_BASE}/dir/create`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ path, recursive })
    });
    return await response.json();
}

// 保存文件
async function saveFile(content, dir, filename) {
    const response = await fetch(`${API_BASE}/file/save`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content, dir, filename })
    });
    return await response.json();
}

// 列出目录
async function listDir(path) {
    const response = await fetch(`${API_BASE}/dir/list?path=${path}`);
    return await response.json();
}

// 使用示例
(async () => {
    await createDir('/my_files');
    await saveFile('Hello World!', '/my_files', 'hello.txt');
    const files = await listDir('/my_files');
    console.log(files);
})();
```

---

## Web UI

打开浏览器访问 `np_fs_web.html` 文件即可使用图形化界面管理文件。

功能包括：
- ✅ 目录检查、创建、列表
- ✅ 文件保存、读取
- ✅ 文件上传、下载
- ✅ 实时结果显示

---

## 注意事项

1. **服务依赖**: 确保 NPFS 后端服务已在 `127.0.0.1:4800` 运行
2. **会话管理**: HTTP 服务启动时会自动创建会话，关闭时自动清理
3. **路径格式**: NPFS 路径以 `/` 开头，例如 `/my_files/test.txt`
4. **文件大小**: 大文件读取时会自动限制显示前 10000 字节
5. **并发安全**: 支持并发请求，每个请求共享同一个 NPFS 会话

---

## 技术栈

- **框架**: Gin v1.12.0
- **语言**: Go 1.25.0
- **底层 RPC**: HPROSE over WebSocket
- **文件系统**: NPFS (Nanping File System)
