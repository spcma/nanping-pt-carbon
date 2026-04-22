package carbonreportmonth

import (
	"app/internal/module/carbonreportday"
	shared_http "app/internal/shared/http"

	"github.com/gin-gonic/gin"
)

// 全局服务实例，供定时任务使用
var defaultService *CarbonReportMonthService

// DefaultService 获取默认的碳报告月报服务实例
func DefaultService() *CarbonReportMonthService {
	return defaultService
}

func setDefaultService(service *CarbonReportMonthService) {
	defaultService = service
}

// carbonReportMonthRoutes CarbonReportMonth 模块路由注册器
type carbonReportMonthRoutes struct {
	carbonReportMonthHandler *CarbonReportMonthHandler
}

// NewCarbonReportMonthRoutes 创建 CarbonReportMonth 模块的路由注册器
func NewCarbonReportMonthRoutes() shared_http.RouteRegistry {

	repo := NewCarbonReportMonthRepository()
	crdService := carbonreportday.DefaultService()
	service := NewCarbonReportMonthService(repo, crdService)
	handler := NewCarbonReportMonthHandler(service)

	// 设置全局服务实例
	setDefaultService(service)

	return &carbonReportMonthRoutes{
		carbonReportMonthHandler: handler,
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
