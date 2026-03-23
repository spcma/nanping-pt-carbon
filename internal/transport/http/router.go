package http

import (
	"app/internal/config"
	initializer_http "app/internal/initializer"
	carbonreportday_http "app/internal/module/carbonreportday/transport/http"
	iam_http "app/internal/module/iam/transport/http"
	ipfs_http "app/internal/module/ipfs/transport/http"
	methodology_http "app/internal/module/methodology/transport/http"
	project_http "app/internal/module/project/transport/http"
	"app/internal/shared/cache"
	shared_http "app/internal/shared/http"
	"app/internal/shared/logger"
	"app/internal/shared/token"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由器（主入口）
func InitRouter(db *gorm.DB, remoteDB *gorm.DB, redisClient *cache.RedisClient) *gin.Engine {
	// 1. 创建基础路由引擎, 并配置中间件
	router := gin.Default()

	// 使用全局中间件
	router.Use(CORSMiddleware(), LogMiddleware(), Recovery())

	var tokenManager token.Manager

	// 2. 使用工厂方法创建 Token 管理器（支持 JWT 和雪花 ID 两种模式）
	if config.GlobalConfig.Token.Expire <= 0 {
		config.GlobalConfig.Token.Expire = 7 * 24 * 60 * 60 // 默认7天
	}

	tokenManager, err := token.NewManager(token.ConfigEx{
		Type:        token.TokenTypeSnowflake,
		RedisClient: redisClient.GetClient(),
		ExpireTime:  time.Duration(config.GlobalConfig.Token.Expire) * time.Second,
	})
	if err != nil {
		logger.HttpLogger.Error("Failed to create token manager: " + err.Error())
		panic(err)
	}

	// 3. 注册所有路由
	registerAllRoutes(router, db, remoteDB, tokenManager)

	logger.HttpLogger.Info("Router initialized")

	return router
}

// registerAllRoutes 注册所有路由
func registerAllRoutes(router *gin.Engine, db *gorm.DB, remoteDB *gorm.DB, tokenManager token.Manager) {
	// API 根路由组
	apiGroup := router.Group("/api")

	// 创建统一的认证中间件映射
	middlewares := map[shared_http.AuthType]gin.HandlerFunc{
		shared_http.AuthTypeNone:     nil,
		shared_http.AuthTypeOptional: OptionalAuthMiddleware(tokenManager),
		shared_http.AuthTypeRequired: AuthMiddleware(tokenManager),
	}

	// 获取所有模块的路由注册器并注册
	registries := getAllRouteRegistries(db, remoteDB, tokenManager)
	for _, registry := range registries {
		registry.RegisterRoutes(apiGroup, middlewares)
	}

	logger.HttpLogger.Info("Router initialized with module-based registration")
}

// getAllRouteRegistries 获取所有模块的路由注册器
//
//	新的模块添加时，请添加到该列表中完成注册
func getAllRouteRegistries(db *gorm.DB, remoteDB *gorm.DB, tokenManager token.Manager) []shared_http.RouteRegistry {
	var registries []shared_http.RouteRegistry

	// 收集各模块的路由注册器
	registries = append(registries, iam_http.NewIAMRoutes(db, tokenManager))
	registries = append(registries, project_http.NewProjectRoutes(db))
	registries = append(registries, methodology_http.NewMethodologyRoutes(db))
	registries = append(registries, carbonreportday_http.NewCarbonReportDayRoutes(db))
	registries = append(registries, ipfs_http.NewIpfsRoutes(db, remoteDB))

	//	初始化数据路由
	registries = append(registries, initializer_http.NewInitializerRoutes(db))

	return registries
}
