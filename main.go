package main

import (
	"app/internal/service"
	"flag"
	"fmt"
	"os"
)

func main() {
	// 定义命令行参数
	mode := flag.String("mode", "http", "运行模式：interactive(交互式) 或 websocket 或 http")
	addr := flag.String("addr", ":19870", "WebSocket/HTTP 服务器监听地址")
	httpAddr := flag.String("http-addr", ":19870", "HTTP 服务器监听地址 (仅在 http 模式下使用)")
	flag.Parse()

	// 检查是否启动 HTTP 服务器
	for _, arg := range os.Args {
		if arg == "-http-server" {
			err := service.StartHTTPServer(*httpAddr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "错误：%v\n", err)
				os.Exit(1)
			}
			return
		}
	}

	// 解析运行模式
	var runMode service.RunMode
	switch *mode {
	case "interactive", "cli":
		runMode = service.ModeInteractive
	case "websocket", "ws":
		runMode = service.ModeWebSocket
	case "http":
		err := service.StartHTTPServer(*httpAddr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误：%v\n", err)
			os.Exit(1)
		}
		return
	default:
		fmt.Fprintf(os.Stderr, "错误：未知的运行模式 '%s'\n", *mode)
		fmt.Println("可用的模式:")
		fmt.Println("  interactive - 交互式命令行模式 (默认)")
		fmt.Println("  websocket   - WebSocket 服务器模式")
		fmt.Println("  http        - HTTP REST API 模式")
		os.Exit(1)
	}

	// 启动服务
	config := service.ServerConfig{
		Mode: runMode,
		Addr: *addr,
	}

	err := service.StartServer(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误：%v\n", err)
		os.Exit(1)
	}
}
