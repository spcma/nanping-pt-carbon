package carbonreportmonth

import (
	"app/internal/shared/db"
	shared_http "app/internal/shared/http"

	"github.com/gin-gonic/gin"
)

// carbonReportMonthRoutes CarbonReportMonth 模块路由注册器
type carbonReportMonthRoutes struct {
	carbonReportMonthHandler *CarbonReportMonthHandler
}

// NewCarbonReportMonthRoutes 创建 CarbonReportMonth 模块的路由注册器
func NewCarbonReportMonthRoutes() shared_http.RouteRegistry {
	dbInst := db.Default()

	//	初始化 carbon_report_month 模块
	carbonReportMonthRepo := NewCarbonReportMonthRepository(dbInst)
	carbonReportMonthService := NewCarbonReportMonthAppService(carbonReportMonthRepo)
	carbonReportMonthHandler := NewCarbonReportMonthHandler(carbonReportMonthService)

	return &carbonReportMonthRoutes{
		carbonReportMonthHandler: carbonReportMonthHandler,
	}
}

// RegisterRoutes 注册 CarbonReportMonth 模块的所有路由（实现 RouteRegistry 接口）
func (r *carbonReportMonthRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {

	// 统一认证中间件
	authTypeRequiredRoute := group.Group("")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		authTypeRequiredRoute.Use(mw)
	}
	{
		// 碳报告月报管理路由 - /api/carbon-report-month/*
		carbonReportMonthGroup := authTypeRequiredRoute.Group("/carbonReportMonth")
		{
			carbonReportMonthGroup.POST("", r.carbonReportMonthHandler.Create)
			carbonReportMonthGroup.PUT("", r.carbonReportMonthHandler.Update)
			carbonReportMonthGroup.DELETE("", r.carbonReportMonthHandler.Delete)
			carbonReportMonthGroup.GET("", r.carbonReportMonthHandler.GetByID) // 仅 ID 查询
		}

		// 碳报告月报列表路由 - /api/carbon-report-months/*
		carbonReportMonthsGroup := authTypeRequiredRoute.Group("/carbonReportMonths")
		{
			carbonReportMonthsGroup.GET("page", r.carbonReportMonthHandler.GetPage)
		}
	}
}
