package transport

import (
	"app/internal/module/carbonreportday/application"
	"app/internal/module/carbonreportday/infrastructure"
	shared_http "app/internal/shared/http"

	"github.com/gin-gonic/gin"
)

// CarbonReportDayRoutes 碳日报路由注册器
type CarbonReportDayRoutes struct {
	handler *CarbonReportDayHandler
}

// NewCarbonReportDayRoutes 创建碳日报路由注册器
func NewCarbonReportDayRoutes() *CarbonReportDayRoutes {
	// 依赖注入：组装各层
	repo := infrastructure.NewCarbonReportDayRepository()
	appService := application.NewCarbonReportDayService(repo)
	handler := NewCarbonReportDayHandler(appService)

	return &CarbonReportDayRoutes{
		handler: handler,
	}
}

// RegisterRoutes 注册路由
func (r *CarbonReportDayRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {
	// 使用认证中间件
	authMiddleware := middlewares[shared_http.AuthTypeRequired]

	// 碳日报管理路由 - /api/carbonReportDay/*
	carbonReportDayGroup := group.Group("/carbonReportDay")
	if authMiddleware != nil {
		carbonReportDayGroup.Use(authMiddleware)
	}
	{
		carbonReportDayGroup.POST("", r.handler.Create)       // 创建碳日报
		carbonReportDayGroup.PUT("/:id", r.handler.Update)    // 更新碳日报
		carbonReportDayGroup.DELETE("/:id", r.handler.Delete) // 删除碳日报
		carbonReportDayGroup.GET("/:id", r.handler.GetByID)   // 根据ID查询
	}

	// 碳日报列表路由 - /api/carbonReportDays/*
	carbonReportDaysGroup := group.Group("/carbonReportDays")
	if authMiddleware != nil {
		carbonReportDaysGroup.Use(authMiddleware)
	}
	{
		carbonReportDaysGroup.GET("/page", r.handler.GetPage) // 分页查询
	}
}
