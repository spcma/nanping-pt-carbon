package server

import (
	"app/internal/bootstrap"
	"app/internal/config"
	"app/internal/infrastructure/db"
	"app/internal/shared/cache"
	idgen "app/internal/shared/idgen"
	"app/internal/shared/logger"
	transport_http "app/internal/transport/http"
	"fmt"
	"github.com/spf13/cast"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server 服务器实例
type Server struct {
	config *config.Config
	db     *gorm.DB
	router *gin.Engine
}

// Initialize 初始化服务器所有组件
func Initialize() (*Server, error) {
	// 加载配置文件
	if err := config.Init("./config/config.yaml"); err != nil {
		return nil, fmt.Errorf("failed to initialize config: %v", err)
	}

	if err := logger.Initialize(&config.GlobalConfig.Logger); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	factory, err := idgen.NewIdgenGenerateFactory(1)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize idgen: %v", err)
	}
	idgen.SetDefault(factory)

	// 初始化数据库
	dbInstance, err := initDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	// 初始化基础数据（超级管理员等）
	if err := initData(dbInstance); err != nil {
		return nil, fmt.Errorf("failed to initialize data: %v", err)
	}

	// 初始化Redis
	redisClient, err := initRedis()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %v", err)
	}

	// 初始化 HTTP 路由
	router := transport_http.InitRouter(dbInstance, redisClient)

	return &Server{
		config: config.GlobalConfig,
		db:     dbInstance,
		router: router,
	}, nil
}

// initData 初始化基础数据
func initData(dbInstance *gorm.DB) error {
	initializer := bootstrap.NewDataInitializer(dbInstance)
	return initializer.Initialize()
}

// initDatabase 初始化数据库连接
func initDatabase() (*gorm.DB, error) {
	dbConfig := db.Config{
		Driver:     "postgres",
		Host:       config.GlobalConfig.Database.Host,
		Port:       config.GlobalConfig.Database.Port,
		User:       config.GlobalConfig.Database.User,
		Password:   config.GlobalConfig.Database.Password,
		DbName:     config.GlobalConfig.Database.DBName,
		SearchPath: config.GlobalConfig.Database.SearchPath,
	}

	dbInstance := db.NewGormDB(dbConfig)
	if dbInstance == nil {
		return nil, fmt.Errorf("failed to create database instance")
	}

	return dbInstance, nil
}

func initRedis() (*cache.RedisClient, error) {
	redisClient, err := cache.NewRedisClient(&cache.RedisConnConfig{
		Host:     config.GlobalConfig.Redis.Host + ":" + cast.ToString(config.GlobalConfig.Redis.Port),
		Password: config.GlobalConfig.Redis.Password,
		DB:       config.GlobalConfig.Redis.DB,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %v", err)
	}

	return redisClient, nil
}

// Run 启动服务器并阻塞直到收到退出信号
func (s *Server) Run() error {
	port := s.config.Server.Port

	if err := s.router.Run(":" + port); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}

// WaitForShutdown 等待关闭信号并优雅关闭
func (s *Server) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
}

// Close 关闭服务器资源
func (s *Server) Close() error {
	// 在这里添加资源清理逻辑
	// 例如：关闭数据库连接、清理缓存等
	if s.db != nil {
		// dbInstance.Close() // 如果 GormDB 有 Close 方法
	}

	fmt.Println("Server resources closed")
	return nil
}
