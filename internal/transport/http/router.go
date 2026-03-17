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
	"time"

	"gorm.io/gorm"

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
		logger.HttpLogger.Error("Failed to create token manager: " + err.Error())
		panic(err)
	}

	// 3. 注册所有路由
	registerAllRoutes(router, db, tokenManager)

	logger.HttpLogger.Info("Router initialized")

	return router
}

// registerAllRoutes 注册所有路由
func registerAllRoutes(router *gin.Engine, db *gorm.DB, tokenManager token.Manager) {
	// API 根路由组
	apiGroup := router.Group("/api")

	// 获取所有模块的路由配置
	allConfigs := getAllRouteConfigs(db, tokenManager)

	// 按认证类型分组并应用中间件
	registerRoutesWithAuth(apiGroup, allConfigs, tokenManager)
}

// getAllRouteConfigs 获取所有模块的路由配置
func getAllRouteConfigs(db *gorm.DB, tokenManager token.Manager) []shared_http.RouteGroupConfig {
	var allConfigs []shared_http.RouteGroupConfig

	// 收集各模块的路由配置
	allConfigs = append(allConfigs, iam_http.RegisterIAMRoutes(db, tokenManager)...)
	allConfigs = append(allConfigs, project_http.RegisterRoutes(db)...)
	allConfigs = append(allConfigs, methodology_http.RegisterRoutes(db)...)
	allConfigs = append(allConfigs, carbonreportday_http.RegisterRoutes(db)...)

	// CarbonReport 模块需要特殊的 RPC 客户端初始化
	carbonReportConfigs := registerCarbonReportRoutes()
	allConfigs = append(allConfigs, carbonReportConfigs...)

	return allConfigs
}

// registerRoutesWithAuth 根据认证类型注册路由并应用中间件
func registerRoutesWithAuth(apiGroup *gin.RouterGroup, routeConfigs []shared_http.RouteGroupConfig, tokenManager token.Manager) {
	// 先验证路由配置
	if err := shared_http.ValidateRouteConfigs(routeConfigs); err != nil {
		logger.Error("http", "Invalid route config: "+err.Error())
		panic(err)
	}

	// 按认证类型分组（不再按 prefix 二次分组）
	noneRoutes := make([]shared_http.RouteHandler, 0)
	optionalRoutes := make([]shared_http.RouteHandler, 0)
	requiredRoutes := make([]shared_http.RouteHandler, 0)

	// 遍历所有路由配置，将 prefix 拼接到 path 中，并按认证类型分组
	for _, routeConfig := range routeConfigs {
		for _, handler := range routeConfig.Handlers {
			// 拼接完整路径
			fullPath := routeConfig.Prefix + handler.Path
			fullHandler := shared_http.RouteHandler{
				Method:  handler.Method,
				Path:    fullPath,
				Handler: handler.Handler,
			}

			// 根据认证类型添加到对应组
			switch routeConfig.AuthType {
			case shared_http.AuthTypeNone:
				noneRoutes = append(noneRoutes, fullHandler)
			case shared_http.AuthTypeOptional:
				optionalRoutes = append(optionalRoutes, fullHandler)
			case shared_http.AuthTypeRequired:
				requiredRoutes = append(requiredRoutes, fullHandler)
			}
		}
	}

	// 注册不同认证类型的路由（直接平铺注册，不再按 prefix 分组）
	registerFlatRoutes(apiGroup, noneRoutes, nil)                                      // 无认证
	registerFlatRoutes(apiGroup, optionalRoutes, OptionalAuthMiddleware(tokenManager)) // 可选认证
	registerFlatRoutes(apiGroup, requiredRoutes, AuthMiddleware(tokenManager))         // 必选认证
}

// registerFlatRoutes 注册平铺的路由（不再按 prefix 分组）
func registerFlatRoutes(apiGroup *gin.RouterGroup, routes []shared_http.RouteHandler, middleware gin.HandlerFunc) {
	// 如果中间件不为空，先在组级别应用
	group := apiGroup
	if middleware != nil {
		group = apiGroup.Group("")
		group.Use(middleware)
	}

	// 直接注册所有路由
	for _, handler := range routes {
		group.Handle(handler.Method, handler.Path, handler.Handler)
	}
}

// registerCarbonReportRoutes 注册 IPFS 模块的路由（需要 RPC 客户端）
func registerCarbonReportRoutes() []shared_http.RouteGroupConfig {
	// 创建 RPC 客户端和 session
	client, session, err := ipfs_http.CreateFsClient()
	if err != nil {
		logger.Error("http", "Failed to create IPFS client: "+err.Error())
		//return nil // 如果初始化失败，返回空路由配置
	}

	// 调用 IPFS 的 RegisterIAMRoutes
	return ipfs_http.RegisterRoutes(client, session)
}
