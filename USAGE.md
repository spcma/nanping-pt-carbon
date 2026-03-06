# NPFS WebSocket 使用指南

## 功能概述

本项目实现了通过 WebSocket 调用 NPFS（Nanping 文件系统）的方法，提供两种运行模式：
1. **交互式命令行模式** - 本地终端交互
2. **WebSocket 服务器模式** - 远程调用

## 可用命令

### 基础命令
- `help` / `h` - 显示帮助信息
- `init` - 初始化 NPFS 服务
- `close` - 关闭 NPFS 服务
- `exit` / `quit` - 退出程序

### 目录操作
- `checkdir <路径>` - 检查目录是否存在
- `createdir <路径> [recursive]` - 创建目录
  - `recursive`: true/false (默认 true)
- `listdir` / `ls <路径>` - 列出目录内容
- `delete` / `rm <路径> [recursive] [force]` - 删除文件/目录
  - `recursive`: true/false (默认 true)
  - `force`: true/false (默认 true)

### 文件操作
- `readfile <NPFS 路径>` - 读取 NPFS 文件内容
- `savefile <NPFS 路径> <本地路径>` - 保存 NPFS 文件到本地
- `upload <本地路径> <NPFS 目录> <文件名>` - 上传本地文件到 NPFS
- `savecontent <内容> <NPFS 目录> <文件名>` - 保存文本内容到 NPFS

### 测试命令
- `batch [数量]` - 批量创建测试文件
- `test` - 运行完整测试

## 使用方法

### 方式一：交互式命令行模式（默认）

```bash
# 直接运行
go run .

# 或使用批处理脚本
interactive.bat
```

然后在提示符后输入命令，例如：
```
npfs> help
npfs> init
npfs> createdir /my_files
npfs> upload C:\test.txt /my_files test.txt
npfs> ls /my_files
npfs> exit
```

### 方式二：WebSocket 服务器模式

#### 1. 启动服务器

```bash
# 使用默认端口 (:8080)
go run . -mode websocket

# 指定端口
go run . -mode websocket -addr :9000

# 或使用批处理脚本
start_server.bat
```

#### 2. 连接客户端

在另一个终端窗口运行：

```bash
# 使用默认地址
go run client_example.go ws://127.0.0.1:8080/ws

# 指定地址
go run client_example.go ws://127.0.0.1:9000/ws

# 或使用批处理脚本
client.bat
```

#### 3. 发送命令示例

连接成功后，可以发送 JSON 格式的命令：

```
help
init
createdir /test
upload ./hello.txt /test hello.txt
readfile /test/hello.txt
```

## 命令示例

### 1. 初始化和创建目录
```bash
init
createdir /storage
createdir /storage/docs true
```

### 2. 上传文件
```bash
upload C:\Users\test\document.pdf /storage/docs document.pdf
```

### 3. 保存文本内容
```bash
savecontent "Hello World" /storage/test hello.txt
```

### 4. 查看目录
```bash
ls /storage
ls /storage/docs
```

### 5. 读取文件
```bash
readfile /storage/test/hello.txt
```

### 6. 下载文件到本地
```bash
savefile /storage/test/hello.txt C:\temp\hello.txt
```

### 7. 删除文件
```bash
delete /storage/test/hello.txt
rm /storage/docs true false
```

### 8. 批量测试
```bash
batch 50
test
```

## WebSocket 协议

### 请求格式
发送纯文本命令或 JSON 格式：

```json
{
  "command": "upload",
  "args": ["./test.txt", "/files", "test.txt"]
}
```

### 响应格式
```json
{
  "command": "upload",
  "result": {
    "ipfsid": "QmXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  },
  "success": true
}
```

错误时：
```json
{
  "command": "upload",
  "error": "文件不存在",
  "success": false
}
```

## API 端点

- **WebSocket**: `ws://<host>:<port>/ws`
- **健康检查**: `http://<host>:<port>/health`

## 注意事项

1. 使用前必须先执行 `init` 初始化 NPFS 服务
2. 路径格式支持相对路径和绝对路径
3. 大文件读取时会限制显示前 1000 字符
4. WebSocket 服务器支持多个客户端同时连接
5. 所有命令都是线程安全的

## 故障排除

### 问题：无法连接 WebSocket
- 检查服务器是否已启动
- 确认端口没有被占用
- 检查防火墙设置

### 问题：命令执行失败
- 确认已执行 `init` 初始化
- 检查路径是否正确
- 查看错误信息详情

### 问题：文件上传失败
- 确认本地文件路径正确
- 检查是否有足够权限
- 确认 NPFS 服务正常运行
