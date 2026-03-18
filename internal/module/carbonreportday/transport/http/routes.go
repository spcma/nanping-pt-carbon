package http

import (
	shared_http "app/internal/shared/http"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// carbonReportDayRoutes CarbonReportDay 模块路由注册器
type carbonReportDayRoutes struct {
	db *gorm.DB
}

// NewCarbonReportDayRoutes 创建 CarbonReportDay 模块的路由注册器
func NewCarbonReportDayRoutes(db *gorm.DB) shared_http.RouteRegistry {
	return &carbonReportDayRoutes{
		db: db,
	}
}

// RegisterRoutes 注册 CarbonReportDay 模块的所有路由（实现 RouteRegistry 接口）
func (r *carbonReportDayRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {
	// 初始化 DDD 组件
	carbonReportDayDDD := InitCarbonReportDayWire(r.db)

	// 创建 handlers
	handlers := &Handlers{
		CarbonReportDayHandler: NewCarbonReportDayHandler(carbonReportDayDDD.Service),
	}

	// 统一认证中间件
	authTypeRequiredRoute := group.Group("")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		authTypeRequiredRoute.Use(mw)
	}
	{
		// 碳报告日报管理路由 - /api/carbon-report-day/*
		carbonReportDayGroup := authTypeRequiredRoute.Group("/carbonReportDay")
		{
			carbonReportDayGroup.POST("", handlers.CarbonReportDayHandler.Create)
			carbonReportDayGroup.PUT("", handlers.CarbonReportDayHandler.Update)
			carbonReportDayGroup.DELETE("", handlers.CarbonReportDayHandler.Delete)
			carbonReportDayGroup.GET("", handlers.CarbonReportDayHandler.GetById) // 仅 ID 查询
		}

		// 碳报告日报列表路由 - /api/carbon-report-days/*
		carbonReportDaysGroup := authTypeRequiredRoute.Group("/carbonReportDays")
		{
			carbonReportDaysGroup.GET("page", handlers.CarbonReportDayHandler.GetPage)
		}
	}
}

// Handlers 包含所有 HTTP 处理器
type Handlers struct {
	CarbonReportDayHandler *CarbonReportDayHandler
}
