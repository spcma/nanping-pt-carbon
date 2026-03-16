package http

import (
	"app/internal/config"
	carbonreportday_http "app/internal/module/carbonreportday/transport/http"
	iam_http "app/internal/module/iam/transport/http"
	ipfs_http "app/internal/module/ipfs"
	methodology_http "app/internal/module/methodology/transport/http"
	project_http "app/internal/module/project/transport/http"
	"app/internal/shared/cache"
	shared_http "app/internal/shared/http"
	"app/internal/shared/logger"
	"app/internal/shared/token"
	"gorm.io/gorm"
	"time"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由器（主入口）
func InitRouter(db *gorm.DB, redisClient *cache.RedisClient) *gin.Engine {
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
		logger.Error("http", "Failed to create token manager: "+err.Error())
		panic(err)
	}

	// 3. 注册所有路由
	registerAllRoutes(router, db, tokenManager)

	logger.Info("http", "Router initialized")

	return router
}

// registerAllRoutes 注册所有路由
func registerAllRoutes(router *gin.Engine, db *gorm.DB, jwtManager token.Manager) {
	// API 根路由组
	apiGroup := router.Group("/api")

	// 获取所有模块的路由配置
	allConfigs := getAllRouteConfigs(db, jwtManager)

	// 按认证类型分组并应用中间件
	registerRoutesWithAuth(apiGroup, allConfigs, jwtManager)
}

// getAllRouteConfigs 获取所有模块的路由配置
func getAllRouteConfigs(db *gorm.DB, jwtManager token.Manager) []shared_http.RouteGroupConfig {
	var allConfigs []shared_http.RouteGroupConfig

	// 收集各模块的路由配置
	allConfigs = append(allConfigs, iam_http.RegisterRoutes(db, jwtManager)...)
	allConfigs = append(allConfigs, project_http.RegisterRoutes(db)...)
	allConfigs = append(allConfigs, methodology_http.RegisterRoutes(db)...)
	allConfigs = append(allConfigs, carbonreportday_http.RegisterRoutes(db)...)

	// CarbonReport 模块需要特殊的 RPC 客户端初始化
	carbonReportConfigs := registerCarbonReportRoutes()
	allConfigs = append(allConfigs, carbonReportConfigs...)

	return allConfigs
}

// registerRoutesWithAuth 根据认证类型注册路由并应用中间件
func registerRoutesWithAuth(apiGroup *gin.RouterGroup, routeConfigs []shared_http.RouteGroupConfig, jwtManager token.Manager) {
	// 先验证路由配置
	if err := shared_http.ValidateRouteConfigs(routeConfigs); err != nil {
		logger.Error("http", "Invalid route config: "+err.Error())
		panic(err)
	}

	// 按认证类型和 prefix 分组
	groupedRoutes := make(map[shared_http.AuthType]map[string][]shared_http.RouteHandler)

	for _, routeConfig := range routeConfigs {
		if groupedRoutes[routeConfig.AuthType] == nil {
			groupedRoutes[routeConfig.AuthType] = make(map[string][]shared_http.RouteHandler)
		}
		groupedRoutes[routeConfig.AuthType][routeConfig.Prefix] = append(
			groupedRoutes[routeConfig.AuthType][routeConfig.Prefix],
			routeConfig.Handlers...,
		)
	}

	// 注册不同认证类型的路由
	registerAuthTypeRoutes(apiGroup, groupedRoutes[shared_http.AuthTypeNone], nil)
	registerAuthTypeRoutes(apiGroup, groupedRoutes[shared_http.AuthTypeOptional], OptionalAuthMiddleware(jwtManager))
	registerAuthTypeRoutes(apiGroup, groupedRoutes[shared_http.AuthTypeRequired], AuthMiddleware(jwtManager))
}

// registerAuthTypeRoutes 注册指定认证类型的路由
func registerAuthTypeRoutes(apiGroup *gin.RouterGroup, routes map[string][]shared_http.RouteHandler, middleware gin.HandlerFunc) {
	for prefix, handlers := range routes {
		group := apiGroup.Group(prefix)
		if middleware != nil {
			group.Use(middleware)
		}
		for _, handler := range handlers {
			group.Handle(handler.Method, handler.Path, handler.Handler)
		}
	}
}

// registerCarbonReportRoutes 注册 IPFS 模块的路由（需要 RPC 客户端）
func registerCarbonReportRoutes() []shared_http.RouteGroupConfig {
	// 创建 RPC 客户端和 session
	client, session, err := ipfs_http.CreateFsClient()
	if err != nil {
		logger.Error("http", "Failed to create IPFS client: "+err.Error())
		return nil // 如果初始化失败，返回空路由配置
	}

	// 调用 IPFS 的 RegisterRoutes
	return ipfs_http.RegisterRoutes(client, session)
}
