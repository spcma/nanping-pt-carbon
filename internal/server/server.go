package server

import (
	"app/internal/config"
	"app/internal/initializer"
	ipfs_application "app/internal/module/ipfs/application"
	"app/internal/module/scheduler"
	"app/internal/shared/cache"
	"app/internal/shared/db"
	idgen "app/internal/shared/idgen"
	"app/internal/shared/logger"
	"app/internal/shared/token"
	transport_http "app/internal/transport/http"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server 服务器实例
type Server struct {
	config    *config.Config
	db        *gorm.DB
	remoteDb  *gorm.DB
	router    *gin.Engine
	scheduler *scheduler.Scheduler
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

	dbConfig := db.Config{
		Driver:     "postgres",
		Host:       config.GlobalConfig.Database.Host,
		Port:       config.GlobalConfig.Database.Port,
		User:       config.GlobalConfig.Database.User,
		Password:   config.GlobalConfig.Database.Password,
		DbName:     config.GlobalConfig.Database.DBName,
		SearchPath: config.GlobalConfig.Database.SearchPath,
		Name:       config.GlobalConfig.Database.Name,
	}

	// 初始化数据源1
	dbInstance, err := initDatabase(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	// 初始化基础数据（超级管理员等）
	if err := initData(dbInstance); err != nil {
		return nil, fmt.Errorf("failed to initialize data: %v", err)
	}

	dbConfig2 := db.Config{
		Driver:     "postgres",
		Host:       config.GlobalConfig.RemoteDatabase.Host,
		Port:       config.GlobalConfig.RemoteDatabase.Port,
		User:       config.GlobalConfig.RemoteDatabase.User,
		Password:   config.GlobalConfig.RemoteDatabase.Password,
		DbName:     config.GlobalConfig.RemoteDatabase.DBName,
		SearchPath: config.GlobalConfig.RemoteDatabase.SearchPath,
		Name:       config.GlobalConfig.RemoteDatabase.Name,
	}

	// 初始化数据源2
	dbInstance2, err := initDatabase(dbConfig2)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize remote_database: %v", err)
	}

	db.RegisterDB("remote", dbInstance2)

	// 初始化Redis
	redisClient, err := initRedis()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %v", err)
	}
	redisClient.SetDefault()

	// 初始化 Token 管理器并设置为默认实例
	if config.GlobalConfig.Token.Expire <= 0 {
		config.GlobalConfig.Token.Expire = 7 * 24 * 60 * 60 // 默认7天
	}
	tokenManager, err := token.NewManager(token.ConfigEx{
		Type:        token.TokenType_Snowflake,
		RedisClient: redisClient.GetClient(),
		ExpireTime:  time.Duration(config.GlobalConfig.Token.Expire) * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create token manager: %v", err)
	}
	token.SetDefault(tokenManager)

	// 初始化定时任务调度器
	// 注意: Start() 必须在所有模块初始化完成后调用,避免任务执行时依赖的 Service 未就绪
	sched := scheduler.Default()

	ipfs_application.RegisterTask()

	// 初始化 HTTP 路由
	// 注意: 此步骤会注册所有模块的路由,并初始化各模块的 Service
	// 必须在调度器启动之前完成,确保调度任务可以安全访问其他模块的 Service
	router := transport_http.InitRouter()

	sched.Start()

	return &Server{
		config:    config.GlobalConfig,
		db:        dbInstance,
		remoteDb:  dbInstance2,
		router:    router,
		scheduler: sched,
	}, nil
}

// initData 初始化基础数据
func initData(dbInstance *gorm.DB) error {
	return initializer.NewDataInitializer(dbInstance).Initialize()
}

// initDatabase 初始化数据库连接
func initDatabase(config db.Config) (*gorm.DB, error) {

	dbInstance := db.NewGormDB(config)
	if dbInstance == nil {
		return nil, fmt.Errorf("failed to create database instance")
	}

	return dbInstance, nil
}

func initRedis() (*cache.RedisClient, error) {
	redisClient, err := cache.NewRedisClient(&cache.RedisConnConfig{
		Addr:     config.GlobalConfig.Redis.Addr,
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
	// 关闭定时任务调度器
	if s.scheduler != nil {
		s.scheduler.Stop()
	}

	// 在这里添加资源清理逻辑
	// 例如：关闭数据库连接、清理缓存等
	if s.db != nil {
		// dbInstance.Close() // 如果 GormDB 有 Close 方法
	}

	fmt.Println("Server resources closed")
	return nil
}
