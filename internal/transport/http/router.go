package http

import (
	initializer_http "app/internal/initializer"
	carbonreportday_http "app/internal/module/carbonreportday"
	carbonreportmonth_http "app/internal/module/carbonreportmonth"
	iam_http "app/internal/module/iam/transport/http"
	ipfs_http "app/internal/module/ipfs/transport/http"
	methodology_http "app/internal/module/methodology"
	project_http "app/internal/module/project"
	scheduler_http "app/internal/module/scheduler"
	shared_http "app/internal/shared/http"
	"app/internal/shared/logger"
	"app/internal/shared/token"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由器（主入口）
func InitRouter() *gin.Engine {
	// 1. 创建基础路由引擎, 并配置中间件
	router := gin.Default()

	// 2. 使用全局中间件
	router.Use(CORSMiddleware(), LogMiddleware(), Recovery())

	// 3. 注册所有路由
	registerAllRoutes(router)

	logger.HTTPL.Info("Router initialized")

	return router
}

// registerAllRoutes 注册所有路由
func registerAllRoutes(router *gin.Engine) {
	// API 根路由组
	apiGroup := router.Group("/api")

	// 创建统一的认证中间件映射（使用默认 Token 管理器）
	tm := token.Default()
	middlewares := map[shared_http.AuthType]gin.HandlerFunc{
		shared_http.AuthTypeNone:     nil,
		shared_http.AuthTypeOptional: OptionalAuthMiddleware(tm),
		shared_http.AuthTypeRequired: AuthMiddleware(tm),
	}

	// 获取所有模块的路由注册器并注册
	registries := getAllRouteRegistries()
	for _, registry := range registries {
		registry.RegisterRoutes(apiGroup, middlewares)
	}

	logger.HTTPL.Info("Router initialized with module-based registration")
}

// getAllRouteRegistries 获取所有模块的路由注册器
//
//	新的模块添加时，请添加到该列表中完成注册
func getAllRouteRegistries() []shared_http.RouteRegistry {
	var registries []shared_http.RouteRegistry

	// 收集各模块的路由注册器（使用默认实例）
	registries = append(registries,
		iam_http.NewIAMRoutes(),
		project_http.NewProjectRoutes(),
		methodology_http.NewMethodologyRoutes(),
		carbonreportday_http.NewCarbonReportDayRoutes(),
		carbonreportmonth_http.NewCarbonReportMonthRoutes(),
		scheduler_http.NewSchedulerRoutes(),
		ipfs_http.NewIpfsRoutes(),
	)

	//	初始化数据路由
	registries = append(registries, initializer_http.NewInitializerRoutes())

	return registries
}
