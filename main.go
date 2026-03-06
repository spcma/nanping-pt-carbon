package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// 定义命令行参数
	mode := flag.String("mode", "websocket", "运行模式：interactive(交互式) 或 websocket")
	addr := flag.String("addr", ":19870", "WebSocket 服务器监听地址 (仅在 websocket 模式下使用)")
	flag.Parse()

	// 解析运行模式
	var runMode RunMode
	switch *mode {
	case "interactive", "cli":
		runMode = ModeInteractive
	case "websocket", "ws":
		runMode = ModeWebSocket
	default:
		fmt.Fprintf(os.Stderr, "错误：未知的运行模式 '%s'\n", *mode)
		fmt.Println("可用的模式:")
		fmt.Println("  interactive - 交互式命令行模式 (默认)")
		fmt.Println("  websocket   - WebSocket 服务器模式")
		os.Exit(1)
	}

	// 启动服务
	config := ServerConfig{
		Mode: runMode,
		Addr: *addr,
	}

	err := StartServer(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误：%v\n", err)
		os.Exit(1)
	}
}
