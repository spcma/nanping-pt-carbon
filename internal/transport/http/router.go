package http

import (
	"app/internal/config"
	carbonreportday_http "app/internal/module/carbonreportday/transport/http"
	iam_http "app/internal/module/iam/transport/http"
	ipfs_http "app/internal/module/ipfs"
	methodology_http "app/internal/module/methodology/transport/http"
	project_http "app/internal/module/project/transport/http"
	shared_http "app/internal/shared/http"
	"app/internal/shared/logger"
	"app/internal/shared/token"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由器（主入口）
func InitRouter(db *gorm.DB) *gin.Engine {
	// 1. 创建基础路由引擎, 并配置中间件
	router := gin.Default()
	// 使用全局中间件
	router.Use(CORSMiddleware(), LogMiddleware(), Recovery())

	// 2. 初始化 JWT Manager
	jwtManager := initJWTManager()

	// 3. 注册所有路由
	registerAllRoutes(router, db, jwtManager)

	logger.Info("http", "Router initialized")

	return router
}

// initJWTManager 初始化 JWT 管理器
func initJWTManager() *token.JWTManager {
	jwtConfig := token.Config{
		SecretKey:     config.GlobalConfig.JWT.Secret,
		ExpireTime:    0, // 使用默认值
		RefreshTime:   0, // 使用默认值
		Issuer:        "",
		BlacklistTime: 0, // 使用默认值
	}
	return token.NewJWTManager(jwtConfig)
}

// registerAllRoutes 注册所有路由
func registerAllRoutes(router *gin.Engine, db *gorm.DB, jwtManager *token.JWTManager) {
	// API 根路由组
	apiGroup := router.Group("/api")

	// 获取所有模块的路由配置
	allConfigs := getAllRouteConfigs(db, jwtManager)

	// 按认证类型分组并应用中间件
	registerRoutesWithAuth(apiGroup, allConfigs, jwtManager)
}

// getAllRouteConfigs 获取所有模块的路由配置
func getAllRouteConfigs(db *gorm.DB, jwtManager *token.JWTManager) []shared_http.RouteGroupConfig {
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
func registerRoutesWithAuth(apiGroup *gin.RouterGroup, routeConfigs []shared_http.RouteGroupConfig, jwtManager *token.JWTManager) {
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
