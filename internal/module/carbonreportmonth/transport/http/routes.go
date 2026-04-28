package transport

import (
	carbonreportday_app "app/internal/module/carbonreportday/application"
	"app/internal/module/carbonreportmonth/application"
	"app/internal/module/carbonreportmonth/infrastructure"
	shared_http "app/internal/shared/http"

	"github.com/gin-gonic/gin"
)

// CarbonReportMonthRoutes 碳月报路由注册器
type CarbonReportMonthRoutes struct {
	handler *CarbonReportMonthHandler
}

// NewCarbonReportMonthRoutes 创建碳月报路由注册器
func NewCarbonReportMonthRoutes() *CarbonReportMonthRoutes {
	// 依赖注入
	repo := infrastructure.NewCarbonReportMonthRepository()

	// 获取碳日报服务（通过默认实例）
	// 创建适配器, 调用碳日报服务
	dayService := carbonreportday_app.Service()
	adapter := infrastructure.NewAdapter(dayService)

	appService := application.NewCarbonReportMonthAppService(repo, adapter)
	handler := NewCarbonReportMonthHandler(appService)

	return &CarbonReportMonthRoutes{
		handler: handler,
	}
}

// RegisterRoutes 注册路由
func (r *CarbonReportMonthRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {
	// 使用认证中间件
	authMiddleware := middlewares[shared_http.AuthTypeRequired]

	// 碳月报管理路由 - /api/carbonReportMonth/*
	carbonReportMonthGroup := group.Group("/carbonReportMonth")
	if authMiddleware != nil {
		carbonReportMonthGroup.Use(authMiddleware)
	}
	{
		carbonReportMonthGroup.POST("", r.handler.Create)       // 创建碳月报
		carbonReportMonthGroup.PUT("/:id", r.handler.Update)    // 更新碳月报
		carbonReportMonthGroup.DELETE("/:id", r.handler.Delete) // 删除碳月报
		carbonReportMonthGroup.GET("/:id", r.handler.GetByID)   // 根据ID查询
	}

	// 碳月报列表路由 - /api/carbonReportMonths/*
	carbonReportMonthsGroup := group.Group("/carbonReportMonths")
	if authMiddleware != nil {
		carbonReportMonthsGroup.Use(authMiddleware)
	}
	{
		carbonReportMonthsGroup.GET("/page", r.handler.GetPage) // 分页查询
	}
}
