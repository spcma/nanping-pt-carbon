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
	// 1. 创建基础路由引擎
	router := createRouterEngine()

	// 2. 初始化 JWT Manager
	jwtManager := initJWTManager()

	// 3. 注册所有路由
	registerAllRoutes(router, db, jwtManager)

	logger.Info("http", "Router initialized")

	return router
}

// createRouterEngine 创建路由引擎并配置中间件
func createRouterEngine() *gin.Engine {
	router := gin.Default()

	// 使用全局中间件
	router.Use(CORSMiddleware(), LogMiddleware(), Recovery())

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

	// 1. 完全公开路由（不需要 token，无任何认证）
	pubRouterGroup := apiGroup.Group("")
	registerAuthGroup(pubRouterGroup, shared_http.AuthTypeNone, db, jwtManager)

	// 2. 支持可选认证的路由（有 token 返回增强信息，无 token 返回基础信息）
	optionalAuthGroup := apiGroup.Group("")
	optionalAuthGroup.Use(OptionalAuthMiddleware(jwtManager))
	registerAuthGroup(optionalAuthGroup, shared_http.AuthTypeOptional, db, jwtManager)

	// 3. 需要 Token 验证的路由（强制要求有效 token）
	authenticated := apiGroup.Group("")
	authenticated.Use(AuthMiddleware(jwtManager))
	registerAuthGroup(authenticated, shared_http.AuthTypeRequired, db, jwtManager)
}

// registerAuthGroup 注册指定认证类型的所有路由
// authType 参数明确指定当前路由组需要的认证类型
func registerAuthGroup(r *gin.RouterGroup, requiredAuthType shared_http.AuthType, db *gorm.DB, jwtManager *token.JWTManager) {
	// 获取所有模块的路由配置
	var allConfigs []shared_http.RouteGroupConfig

	// 直接调用各模块的 RegisterRoutes 函数
	allConfigs = append(allConfigs, iam_http.RegisterRoutes(db, jwtManager)...)
	allConfigs = append(allConfigs, project_http.RegisterRoutes(db)...)
	allConfigs = append(allConfigs, methodology_http.RegisterRoutes(db)...)
	allConfigs = append(allConfigs, carbonreportday_http.RegisterRoutes(db)...)

	// CarbonReport 模块需要特殊的 RPC 客户端初始化
	carbonReportConfigs := registerCarbonReportRoutes()
	allConfigs = append(allConfigs, carbonReportConfigs...)

	// 按 prefix 分组并注册路由
	registerRouteGroupsByAuthType(r, allConfigs, requiredAuthType)
}

// registerRouteGroupsByAuthType 根据认证类型批量注册路由组
// requiredAuthType 指定当前路由组需要的认证类型，只注册匹配的路由
func registerRouteGroupsByAuthType(r *gin.RouterGroup, routeConfigs []shared_http.RouteGroupConfig, requiredAuthType shared_http.AuthType) {
	// 先验证路由配置
	if err := shared_http.ValidateRouteConfigs(routeConfigs); err != nil {
		logger.Error("http", "Invalid route config: "+err.Error())
		panic(err)
	}

	// 按 prefix 分组合并路由
	prefixMap := make(map[string][]shared_http.RouteHandler)
	for _, config := range routeConfigs {
		// 只注册 AuthType 与当前路由组匹配的路由
		if config.AuthType == requiredAuthType {
			prefixMap[config.Prefix] = append(prefixMap[config.Prefix], config.Handlers...)
		}
	}

	// 批量注册路由
	for prefix, handlers := range prefixMap {
		group := r.Group(prefix)
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
