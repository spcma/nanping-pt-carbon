// WebSocket 客户端示例 - 用于测试 NPFS WebSocket 服务器
// 注意：此客户端使用标准库的 websocket 包
// 运行前需要安装：go get golang.org/x/net/websocket
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/websocket"
)

func main_1() {
	if len(os.Args) < 2 {
		fmt.Println("用法：go run client_example.go [ws://地址]")
		fmt.Println("示例：go run client_example.go ws://127.0.0.1:8080/ws")
		os.Exit(1)
	}

	serverAddr := os.Args[1]
	if !strings.HasPrefix(serverAddr, "ws") {
		serverAddr = "ws://" + serverAddr
	}
	if !strings.HasSuffix(serverAddr, "/ws") {
		serverAddr = strings.TrimSuffix(serverAddr, "/") + "/ws"
	}

	fmt.Printf("连接到 NPFS WebSocket 服务器：%s\n", serverAddr)

	// 创建 WebSocket 连接
	conn, err := websocket.Dial(serverAddr, "", "http://localhost")
	if err != nil {
		fmt.Printf("连接失败：%v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("✓ 连接成功！")
	fmt.Println("输入命令发送，输入 'exit' 或 'quit' 退出")
	fmt.Println("─────────────────────────────────────")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
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
			break
		}

		// 发送命令
		err = websocket.Message.Send(conn, line)
		if err != nil {
			fmt.Printf("发送失败：%v\n", err)
			break
		}

		// 接收响应
		var response string
		err = websocket.Message.Receive(conn, &response)
		if err != nil {
			fmt.Printf("接收失败：%v\n", err)
			break
		}

		// 格式化并打印响应
		printResponse(response)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "读取输入错误：%v\n", err)
	}
}

// printResponse 格式化打印 JSON 响应
func printResponse(jsonStr string) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		fmt.Println(jsonStr)
		return
	}

	// 美化输出
	prettyData, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(prettyData))
	fmt.Println()
}
