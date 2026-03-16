package main

import (
	"app/internal/server"
	"fmt"
	"log"
)

func main() {
	// 初始化服务器
	srv, err := server.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// 启动服务器（阻塞）
	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	// 等待关闭信号
	srv.WaitForShutdown()

	// 关闭资源
	if err := srv.Close(); err != nil {
		log.Printf("Error closing server resources: %v", err)
	}

	fmt.Println("Server shutdown complete")
}
