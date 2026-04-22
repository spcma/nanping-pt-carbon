package carbonreportday

import (
	shared_http "app/internal/shared/http"

	"github.com/gin-gonic/gin"
)

// 全局服务实例，供定时任务使用
var defaultService *CarbonReportDayService

// DefaultService 获取默认的碳报告月报服务实例
func DefaultService() *CarbonReportDayService {
	return defaultService
}

func setDefaultService(service *CarbonReportDayService) {
	defaultService = service
}

// carbonReportDayRoutes CarbonReportDay 模块路由注册器
type carbonReportDayRoutes struct {
	carbonReportDayHandler *CarbonReportDayHandler
}

// NewCarbonReportDayRoutes 创建 CarbonReportDay 模块的路由注册器
func NewCarbonReportDayRoutes() shared_http.RouteRegistry {
	repo := NewCarbonReportDayRepository()
	service := NewCarbonReportDayService(repo)
	handler := NewCarbonReportDayHandler(service)

	setDefaultService(service)

	return &carbonReportDayRoutes{
		carbonReportDayHandler: handler,
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
