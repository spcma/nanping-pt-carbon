package http

import (
	"app/internal/module/methodology/wire"
	shared_http "app/internal/shared/http"
	"gorm.io/gorm"
)

// RegisterRoutes 注册 Methodology 模块的所有路由
func RegisterRoutes(db *gorm.DB) []shared_http.RouteGroupConfig {
	// 初始化 DDD 组件
	methodologyDDD := wire.InitMethodologyDDD(db)

	// 创建 handlers
	handlers := &Handlers{
		MethodologyHandler: NewMethodologyHandler(methodologyDDD.AppService),
	}

	return handlers.registerRoutes()
}

// Handlers 包含所有 HTTP 处理器
type Handlers struct {
	MethodologyHandler *MethodologyHandler
}

// registerRoutes 注册 Methodology 模块的所有路由（内部方法）
// registerRoutes 注册 Methodology 模块的所有路由（内部方法）
func (h *Handlers) registerRoutes() []shared_http.RouteGroupConfig {
	return []shared_http.RouteGroupConfig{
		// 需要认证的方法学管理路由
		{
			Prefix: "/methodology",
			Handlers: []shared_http.RouteHandler{
				{Method: "POST", Path: "", Handler: h.MethodologyHandler.Create},
				{Method: "PUT", Path: "/:id", Handler: h.MethodologyHandler.Update},
				{Method: "DELETE", Path: "/:id", Handler: h.MethodologyHandler.Delete},
				{Method: "GET", Path: "/:id", Handler: h.MethodologyHandler.GetByID},
				{Method: "PUT", Path: "/:id/status", Handler: h.MethodologyHandler.ChangeStatus},
			},
			AuthType: shared_http.AuthTypeRequired,
		},
		{
			Prefix: "/methodologies",
			Handlers: []shared_http.RouteHandler{
				{Method: "GET", Path: "", Handler: h.MethodologyHandler.GetPage},
			},
			AuthType: shared_http.AuthTypeRequired,
		},
		{
			Prefix: "/methodology/code",
			Handlers: []shared_http.RouteHandler{
				{Method: "GET", Path: "/:code", Handler: h.MethodologyHandler.GetByCode},
			},
			AuthType: shared_http.AuthTypeRequired,
		},
	}
}
