package main

import (
	"app/internal/config"
	"app/internal/shared/db"
	"app/internal/shared/logger"
	transport_http "app/internal/transport/http"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//	加载配置文件
	if err := config.Init("./config/config.yaml"); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// 初始化日志系统
	if err := logger.Initialize(&config.GlobalConfig.Logger); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	// 初始化全局 Sugar Logger
	logger.InitGlobalLoggers()
	logger.RuntimeLogger.Info("Logger initialized successfully")

	dbConfig := db.Config{
		Driver:     "postgres", // 默认postgres数据库
		Host:       config.GlobalConfig.Database.Host,
		Port:       config.GlobalConfig.Database.Port,
		User:       config.GlobalConfig.Database.User,
		Password:   config.GlobalConfig.Database.Password,
		DbName:     config.GlobalConfig.Database.DBName,
		SearchPath: config.GlobalConfig.Database.SearchPath,
	}
	dbInstance := db.NewGormDB(dbConfig)

	// 初始化 http router
	router := transport_http.InitRouter(dbInstance)
	if err := router.Run(":" + config.GlobalConfig.Server.Port); err != nil {
		logger.Error("runtime", fmt.Sprintf("Failed to start server: %v", err))
		log.Fatalf("Failed to start server: %v", err)
	}
	logger.Info("runtime", fmt.Sprintf("Starting server on port %s...", config.GlobalConfig.Server.Port))

	fmt.Println("Application started successfully")

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("runtime", "Shutting down server...")
}

func Run() {
	//// 定义命令行参数
	//mode := flag.String("mode", "http", "运行模式：http (HTTP REST API)")
	//addr := flag.String("addr", ":19870", "HTTP 服务器监听地址")
	//flag.Parse()
	//
	//// 检查是否启动 HTTP 服务器
	//if *mode == "http" {
	//	err := startHTTPServer(*addr)
	//	if err != nil {
	//		fmt.Fprintf(os.Stderr, "错误：%v\n", err)
	//		os.Exit(1)
	//	}
	//	return
	//}
	//
	//fmt.Fprintf(os.Stderr, "错误：未知的运行模式 '%s'\n", *mode)
	//fmt.Println("可用的模式:")
	//fmt.Println("  http        - HTTP REST API 模式")
	//os.Exit(1)
}

// startHTTPServer 启动 HTTP 服务器
func startHTTPServer(addr string) error {
	// 初始化基础设施层
	//client, session, err := infrastructure.CreateFsClient()
	//if err != nil {
	//	return fmt.Errorf("创建 FS 客户端失败：%v", err)
	//}
	//defer func() {
	//	if client != nil {
	//		client.Logout(session, "")
	//	}
	//}()
	//
	//// 初始化仓储层
	//fileRepo := infrastructure.NewNpfsFileRepository(client, session)
	//sessionRepo := infrastructure.NewNpfsSessionRepository(client, session)
	//
	//// 初始化应用层
	//fileAppService := application.NewFileApplicationService(fileRepo, sessionRepo)
	//
	//// 设置路由
	//router := http.SetupRouter(fileAppService)
	//
	//fmt.Println("=====================================")
	//fmt.Println("   NPFS HTTP 服务器 (DDD 架构)")
	//fmt.Println("=====================================")
	//fmt.Printf("监听地址：http://%s\n", addr)
	//fmt.Println("API 端点:")
	//fmt.Println("  - GET  /health                    健康检查")
	//fmt.Println("  - POST /api/v1/dir/check          检查目录")
	//fmt.Println("  - POST /api/v1/dir/create         创建目录")
	//fmt.Println("  - GET  /api/v1/dir/list           列出目录")
	//fmt.Println("  - DELETE /api/v1/dir/delete       删除目录")
	//fmt.Println("  - GET  /api/v1/file/read          读取文件")
	//fmt.Println("  - POST /api/v1/file/save          保存文件")
	//fmt.Println("  - POST /api/v1/file/upload        上传文件")
	//fmt.Println("  - GET  /api/v1/file/download      下载文件")
	//fmt.Println("  - DELETE /api/v1/file/delete      删除文件")
	//fmt.Println("=====================================")
	//
	//log.Fatal(router.Run(addr))
	return nil
}
