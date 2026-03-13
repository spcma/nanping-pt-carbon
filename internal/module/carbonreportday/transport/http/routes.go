package http

import (
	"app/internal/module/carbonreportday/wire"
	shared_http "app/internal/shared/http"
	"gorm.io/gorm"
)

// RegisterRoutes 注册 CarbonReportDay 模块的所有路由
func RegisterRoutes(db *gorm.DB) []shared_http.RouteGroupConfig {
	// 初始化 DDD 组件
	carbonReportDayDDD := wire.InitCarbonReportDayDDD(db)

	// 创建 handlers
	handlers := &Handlers{
		CarbonReportDayHandler: NewCarbonReportDayHandler(carbonReportDayDDD.AppService),
	}

	return handlers.registerRoutes()
}

// Handlers 包含所有 HTTP 处理器
type Handlers struct {
	CarbonReportDayHandler *CarbonReportDayHandler
}

// registerRoutes 注册 CarbonReportDay 模块的所有路由（内部方法）
// registerRoutes 注册 CarbonReportDay 模块的所有路由（内部方法）
func (h *Handlers) registerRoutes() []shared_http.RouteGroupConfig {
	return []shared_http.RouteGroupConfig{
		// 需要认证的碳报告日报管理路由
		{
			Prefix: "/carbon-report-day",
			Handlers: []shared_http.RouteHandler{
				{Method: "POST", Path: "", Handler: h.CarbonReportDayHandler.Create},
				{Method: "PUT", Path: "/:id", Handler: h.CarbonReportDayHandler.Update},
				{Method: "DELETE", Path: "/:id", Handler: h.CarbonReportDayHandler.Delete},
				{Method: "GET", Path: "/:id", Handler: h.CarbonReportDayHandler.GetByID},
			},
			AuthType: shared_http.AuthTypeRequired,
		},
		{
			Prefix: "/carbon-report-days",
			Handlers: []shared_http.RouteHandler{
				{Method: "GET", Path: "", Handler: h.CarbonReportDayHandler.GetPage},
			},
			AuthType: shared_http.AuthTypeRequired,
		},
	}
}
