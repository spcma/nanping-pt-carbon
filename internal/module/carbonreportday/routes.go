package carbonreportday

import (
	shared_http "app/internal/shared/http"

	"github.com/gin-gonic/gin"
)

// carbonReportDayRoutes CarbonReportDay 模块路由注册器
type carbonReportDayRoutes struct {
	carbonReportDayHandler *CarbonReportDayHandler
}

// NewCarbonReportDayRoutes 创建 CarbonReportDay 模块的路由注册器
func NewCarbonReportDayRoutes() shared_http.RouteRegistry {
	//	初始化 carbon_report_day 模块
	carbonReportDayRepo := NewCarbonReportDayRepository()
	carbonReportDayService := NewCarbonReportDayAppService(carbonReportDayRepo)
	carbonReportDayHandler := NewCarbonReportDayHandler(carbonReportDayService)

	return &carbonReportDayRoutes{
		carbonReportDayHandler: carbonReportDayHandler,
	}
}

// RegisterRoutes 注册 CarbonReportDay 模块的所有路由（实现 RouteRegistry 接口）
func (r *carbonReportDayRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {

	// 统一认证中间件
	authTypeRequiredRoute := group.Group("")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		authTypeRequiredRoute.Use(mw)
	}
	{
		// 碳报告日报管理路由 - /api/carbon-report-day/*
		carbonReportDayGroup := authTypeRequiredRoute.Group("/carbonReportDay")
		{
			carbonReportDayGroup.POST("", r.carbonReportDayHandler.Create)
			carbonReportDayGroup.PUT("", r.carbonReportDayHandler.Update)
			carbonReportDayGroup.DELETE("", r.carbonReportDayHandler.Delete)
			carbonReportDayGroup.GET("", r.carbonReportDayHandler.GetByID) // 仅 ID 查询
		}

		// 碳报告日报列表路由 - /api/carbon-report-days/*
		carbonReportDaysGroup := authTypeRequiredRoute.Group("/carbonReportDays")
		{
			carbonReportDaysGroup.GET("page", r.carbonReportDayHandler.GetPage)
		}
	}
}
