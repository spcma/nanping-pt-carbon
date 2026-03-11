package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

// ==================== WebSocket 命令处理器 ====================

// CommandHandler WebSocket 命令处理器
type CommandHandler struct {
	mu sync.Mutex
}

// NewCommandHandler 创建命令处理器
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{}
}

// HandleCommand 处理单个命令
func (h *CommandHandler) HandleCommand(conn *websocket.Conn, cmd string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	fmt.Println("33", cmd)

	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return nil
	}

	fmt.Println("parts", parts)

	command := parts[0]
	args := parts[1:]

	var result interface{}
	var err error

	switch command {
	case "help", "h":
		result = h.showHelp()
	case "init":
		err = h.cmdInit()
		result = "NPFS 服务已初始化"
	case "close":
		h.cmdClose()
		result = "NPFS 服务已关闭"
	case "checkdir":
		if len(args) < 1 {
			return fmt.Errorf("用法：checkdir <路径>")
		}
		exists, err2 := CheckDirExists(args[0])
		result = map[string]interface{}{"exists": exists}
		err = err2
	case "createdir":
		if len(args) < 1 {
			return fmt.Errorf("用法：createdir <路径> [recursive=true]")
		}
		recursive := true
		if len(args) > 1 {
			recursive = args[1] == "true"
		}
		created, err2 := CreateDirIfNotExists(args[0], recursive)
		result = map[string]interface{}{"created": created}
		err = err2
	case "listdir", "ls":
		fmt.Println("正在列出目录...")
		if len(args) < 1 {
			return fmt.Errorf("用法：listdir <路径>")
		}
		links, err2 := ListDirectory(args[0])
		result = links
		err = err2
	case "delete", "rm":
		if len(args) < 1 {
			return fmt.Errorf("用法：delete <路径> [recursive=true] [force=true]")
		}
		recursive := true
		force := true
		if len(args) > 1 {
			recursive = args[1] == "true"
		}
		if len(args) > 2 {
			force = args[2] == "true"
		}
		err = DeleteFile(args[0], recursive, force)
		result = "删除成功"
	case "readfile":
		if len(args) < 1 {
			return fmt.Errorf("用法：readfile <文件路径>")
		}
		data, size, err2 := ReadFileFromNpfs(args[0])
		result = map[string]interface{}{
			"size": size,
			"data": string(data),
		}
		err = err2
	case "savefile":
		if len(args) < 2 {
			return fmt.Errorf("用法：savefile <NPFS 路径> <本地路径>")
		}
		err = SaveFileToLocal(args[0], args[1])
		result = "保存成功"
	case "upload":
		if len(args) < 3 {
			return fmt.Errorf("用法：upload <本地路径> <NPFS 目录> <文件名>")
		}
		ipfsid, err2 := SaveLocalFileToNpfs(args[0], args[1], args[2])
		result = map[string]interface{}{"ipfsid": ipfsid}
		err = err2
	case "savecontent":
		if len(args) < 3 {
			return fmt.Errorf("用法：savecontent <内容> <NPFS 目录> <文件名>")
		}
		ipfsid, err2 := SaveContentToNpfs(args[0], args[1], args[2])
		result = map[string]interface{}{"ipfsid": ipfsid}
		err = err2
	case "batch":
		count := 10
		if len(args) > 0 {
			fmt.Sscanf(args[0], "%d", &count)
		}
		go BatchCreateFilesExample(count)
		result = fmt.Sprintf("正在批量创建 %d 个文件...", count)
	case "test":
		go TestNpFsOperations()
		result = "正在运行测试..."
	default:
		return fmt.Errorf("未知命令：%s", command)
	}

	// 发送响应
	response := map[string]interface{}{
		"command": command,
		"result":  result,
		"success": err == nil,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	return conn.WriteJSON(response)
}

// showHelp 显示帮助信息
func (h *CommandHandler) showHelp() map[string]interface{} {
	return map[string]interface{}{
		"commands": []map[string]string{
			{"cmd": "help, h", "desc": "显示帮助信息"},
			{"cmd": "init", "desc": "初始化 NPFS 服务"},
			{"cmd": "close", "desc": "关闭 NPFS 服务"},
			{"cmd": "checkdir <路径>", "desc": "检查目录是否存在"},
			{"cmd": "createdir <路径> [recursive]", "desc": "创建目录"},
			{"cmd": "listdir, ls <路径>", "desc": "列出目录内容"},
			{"cmd": "delete, rm <路径> [recursive] [force]", "desc": "删除文件/目录"},
			{"cmd": "readfile <NPFS 路径>", "desc": "读取 NPFS 文件"},
			{"cmd": "savefile <NPFS 路径> <本地路径>", "desc": "保存 NPFS 文件到本地"},
			{"cmd": "upload <本地路径> <NPFS 目录> <文件名>", "desc": "上传本地文件到 NPFS"},
			{"cmd": "savecontent <内容> <NPFS 目录> <文件名>", "desc": "保存文本内容到 NPFS"},
			{"cmd": "batch [数量]", "desc": "批量创建测试文件"},
			{"cmd": "test", "desc": "运行完整测试"},
		},
	}
}

// cmdInit 初始化命令
func (h *CommandHandler) cmdInit() error {
	InitNpFsService()
	return nil
}

// cmdClose 关闭命令
func (h *CommandHandler) cmdClose() {
	CloseNpFsService()
}

// ==================== WebSocket 服务器 ====================

// WSServer WebSocket 服务器
type WSServer struct {
	addr    string
	handler *CommandHandler
}

// NewWSServer 创建 WebSocket 服务器
func NewWSServer(addr string) *WSServer {
	return &WSServer{
		addr:    addr,
		handler: NewCommandHandler(),
	}
}

var upgrader = websocket.FastHTTPUpgrader{
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true // 允许所有来源
	},
}

// ServeWebSocket 处理 WebSocket 连接
func (ws *WSServer) ServeWebSocket(ctx *fasthttp.RequestCtx) {
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		defer conn.Close()

		fmt.Printf("客户端已连接：%s\n", conn.RemoteAddr())

		//	初始化NPFS服务
		InitNpFsService()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("读取消息错误：%v\n", err)
				break
			}

			cmd := strings.TrimSpace(string(message))
			if cmd == "" {
				continue
			}

			fmt.Printf("收到命令：%s\n", cmd, "new", cmd)

			// 处理命令
			err = ws.handler.HandleCommand(conn, cmd)
			if err != nil {
				sendError(conn, err.Error())
			}
		}

		fmt.Printf("客户端断开连接：%s\n", conn.RemoteAddr())
	})

	if err != nil {
		fmt.Printf("WebSocket 升级错误：%v\n", err)
	}
}

// sendError 发送错误响应
func sendError(conn *websocket.Conn, errMsg string) {
	response := map[string]interface{}{
		"success": false,
		"error":   errMsg,
	}
	conn.WriteJSON(response)
}

// Start 启动 WebSocket 服务器
func (ws *WSServer) Start() error {
	fmt.Println("=====================================")
	fmt.Println("   NPFS WebSocket 服务器")
	fmt.Println("=====================================")
	fmt.Printf("监听地址：%s\n", ws.addr)
	fmt.Println("使用 WebSocket 客户端连接发送命令")
	fmt.Println("按 Ctrl+C 停止服务器")
	fmt.Println("=====================================")

	// 同时提供 HTTP 接口
	httpServer := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			switch string(ctx.Path()) {
			case "/ws":
				ws.ServeWebSocket(ctx)
			case "/health":
				ctx.SetStatusCode(fasthttp.StatusOK)
				ctx.SetBodyString("OK")
			default:
				ctx.SetStatusCode(fasthttp.StatusNotFound)
				ctx.SetBodyString("404 Not Found")
			}
		},
	}

	return httpServer.ListenAndServe(ws.addr)
}

// RunInteractive 运行交互式命令行（通过标准输入）
func (h *CommandHandler) RunInteractive() {
	fmt.Println("=====================================")
	fmt.Println("   NPFS 交互式命令行")
	fmt.Println("=====================================")
	fmt.Println("输入 'help' 查看可用命令")
	fmt.Println("输入 'exit' 或 'quit' 退出程序")
	fmt.Println("=====================================")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("npfs> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// 检查退出命令
		if line == "exit" || line == "quit" {
			fmt.Println("再见！")
			CloseNpFsService()
			break
		}

		// 处理命令
		parts := strings.Fields(line)
		command := parts[0]
		args := parts[1:]

		var result interface{}
		var err error

		switch command {
		case "help", "h":
			help := h.showHelp()
			printHelp(help)
			continue
		case "init":
			err = h.cmdInit()
			result = "✓ NPFS 服务已初始化"
		case "close":
			h.cmdClose()
			result = "✓ NPFS 服务已关闭"
		case "checkdir":
			if len(args) < 1 {
				fmt.Println("✗ 错误：用法：checkdir <路径>")
				continue
			}
			exists, err2 := CheckDirExists(args[0])
			if err2 != nil {
				fmt.Printf("✗ 错误：%v\n", err2)
				continue
			}
			if exists {
				fmt.Printf("✓ 目录 '%s' 存在\n", args[0])
			} else {
				fmt.Printf("✗ 目录 '%s' 不存在\n", args[0])
			}
			continue
		case "createdir":
			if len(args) < 1 {
				fmt.Println("✗ 错误：用法：createdir <路径> [recursive=true]")
				continue
			}
			recursive := true
			if len(args) > 1 {
				recursive = args[1] == "true"
			}
			created, err2 := CreateDirIfNotExists(args[0], recursive)
			if err2 != nil {
				fmt.Printf("✗ 错误：%v\n", err2)
				continue
			}
			if created {
				fmt.Printf("✓ 目录 '%s' 已创建\n", args[0])
			} else {
				fmt.Printf("✓ 目录 '%s' 已存在\n", args[0])
			}
			continue
		case "listdir", "ls":
			if len(args) < 1 {
				fmt.Println("✗ 错误：用法：listdir <路径>")
				continue
			}
			links, err2 := ListDirectory(args[0])
			if err2 != nil {
				fmt.Printf("✗ 错误：%v\n", err2)
				continue
			}
			fmt.Printf("目录 '%s' 的内容:\n", args[0])
			for _, link := range links {
				typeStr := "文件"
				if link.IsDir() {
					typeStr = "目录"
				}
				fmt.Printf("  [%s] %s (大小：%d bytes)\n", typeStr, link.Name, link.Size)
			}
			continue
		case "delete", "rm":
			if len(args) < 1 {
				fmt.Println("✗ 错误：用法：delete <路径> [recursive=true] [force=true]")
				continue
			}
			recursive := true
			force := true
			if len(args) > 1 {
				recursive = args[1] == "true"
			}
			if len(args) > 2 {
				force = args[2] == "true"
			}
			err = DeleteFile(args[0], recursive, force)
			if err != nil {
				fmt.Printf("✗ 错误：%v\n", err)
				continue
			}
			fmt.Printf("✓ '%s' 已删除\n", args[0])
			continue
		case "readfile":
			if len(args) < 1 {
				fmt.Println("✗ 错误：用法：readfile <文件路径>")
				continue
			}
			data, size, err2 := ReadFileFromNpfs(args[0])
			if err2 != nil {
				fmt.Printf("✗ 错误：%v\n", err2)
				continue
			}
			fmt.Printf("✓ 文件大小：%d bytes\n", size)
			if size < 1000 {
				fmt.Printf("内容:\n%s\n", string(data))
			} else {
				fmt.Println("(文件太大，仅显示前 1000 字符)")
				fmt.Printf("%.1000s\n", string(data))
			}
			continue
		case "savefile":
			if len(args) < 2 {
				fmt.Println("✗ 错误：用法：savefile <NPFS 路径> <本地路径>")
				continue
			}
			err = SaveFileToLocal(args[0], args[1])
			if err != nil {
				fmt.Printf("✗ 错误：%v\n", err)
				continue
			}
			fmt.Printf("✓ 文件已保存到 '%s'\n", args[1])
			continue
		case "upload":
			if len(args) < 3 {
				fmt.Println("✗ 错误：用法：upload <本地路径> <NPFS 目录> <文件名>")
				continue
			}
			ipfsid, err2 := SaveLocalFileToNpfs(args[0], args[1], args[2])
			if err2 != nil {
				fmt.Printf("✗ 错误：%v\n", err2)
				continue
			}
			fmt.Printf("✓ 文件已上传，IPFS ID: %s\n", ipfsid)
			continue
		case "savecontent":
			if len(args) < 3 {
				fmt.Println("✗ 错误：用法：savecontent <内容> <NPFS 目录> <文件名>")
				continue
			}
			ipfsid, err2 := SaveContentToNpfs(args[0], args[1], args[2])
			if err2 != nil {
				fmt.Printf("✗ 错误：%v\n", err2)
				continue
			}
			fmt.Printf("✓ 内容已保存，IPFS ID: %s\n", ipfsid)
			continue
		case "batch":
			count := 10
			if len(args) > 0 {
				fmt.Sscanf(args[0], "%d", &count)
			}
			fmt.Printf("正在批量创建 %d 个文件...\n", count)
			BatchCreateFilesExample(count)
			fmt.Println("✓ 批量创建完成")
			continue
		case "test":
			fmt.Println("正在运行完整测试...")
			TestNpFsOperations()
			fmt.Println("✓ 测试完成")
			continue
		default:
			fmt.Printf("✗ 未知命令：%s\n", command)
			fmt.Println("输入 'help' 查看可用命令")
			continue
		}

		if err != nil {
			fmt.Printf("✗ 错误：%v\n", err)
		} else {
			fmt.Println(result)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "读取输入错误：%v\n", err)
	}
}

// printHelp 打印帮助信息
func printHelp(help map[string]interface{}) {
	fmt.Println("\n可用命令:")
	fmt.Println("─────────────────────────────────────────")
	commands := help["commands"].([]map[string]string)
	for _, cmd := range commands {
		fmt.Printf("  %-35s %s\n", cmd["cmd"], cmd["desc"])
	}
	fmt.Println("─────────────────────────────────────────\n")
}

// ==================== JSON 命令解析 ====================

// CommandRequest 命令请求结构
type CommandRequest struct {
	Command string      `json:"command"`
	Args    []string    `json:"args,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// CommandResponse 命令响应结构
type CommandResponse struct {
	Command string      `json:"command"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
	Success bool        `json:"success"`
}

// ParseCommand 解析 JSON 命令
func ParseCommand(jsonStr string) (*CommandRequest, error) {
	var req CommandRequest
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

// FormatResponse 格式化响应
func FormatResponse(req *CommandRequest, result interface{}, err error) string {
	resp := CommandResponse{
		Command: req.Command,
		Result:  result,
		Success: err == nil,
	}

	if err != nil {
		resp.Error = err.Error()
	}

	data, _ := json.MarshalIndent(resp, "", "  ")
	return string(data)
}

// ==================== 启动模式 ====================

// RunMode 运行模式
type RunMode int

const (
	// ModeInteractive 交互式命令行模式
	ModeInteractive RunMode = iota
	// ModeWebSocket WebSocket 服务器模式
	ModeWebSocket
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Mode RunMode
	Addr string // WebSocket 监听地址，例如 ":8080"
}

// StartServer 启动 NPFS 服务
func StartServer(config ServerConfig) error {
	handler := NewCommandHandler()

	switch config.Mode {
	case ModeInteractive:
		fmt.Println("提示：可以在命令行中输入 'init' 初始化 NPFS 服务")
		fmt.Println()
		handler.RunInteractive()
		return nil

	case ModeWebSocket:
		if config.Addr == "" {
			config.Addr = ":19870" // 默认地址
		}
		wsServer := NewWSServer(config.Addr)
		return wsServer.Start()

	default:
		return fmt.Errorf("未知的运行模式：%d", config.Mode)
	}
}
