package http

import (
	"app/internal/module/project/wire"
	shared_http "app/internal/shared/http"

	"gorm.io/gorm"
)

// RegisterRoutes 注册 Project 模块的所有路由
func RegisterRoutes(db *gorm.DB) []shared_http.RouteGroupConfig {
	// 初始化 DDD 组件
	projectDDD := wire.InitProjectDDD(db)

	// 创建 handlers
	handlers := &Handlers{
		ProjectHandler: NewProjectHandler(projectDDD.AppService),
	}

	return handlers.registerRoutes()
}

// Handlers 包含所有 HTTP 处理器
type Handlers struct {
	ProjectHandler *ProjectHandler
}

// registerRoutes 注册 Project 模块的所有路由（内部方法）
// registerRoutes 注册 Project 模块的所有路由（内部方法）
func (h *Handlers) registerRoutes() []shared_http.RouteGroupConfig {
	return []shared_http.RouteGroupConfig{
		// 需要认证的项目管理路由
		{
			Prefix: "/project",
			Handlers: []shared_http.RouteHandler{
				{Method: "POST", Path: "", Handler: h.ProjectHandler.Create},
				{Method: "PUT", Path: "/:id", Handler: h.ProjectHandler.Update},
				{Method: "DELETE", Path: "/:id", Handler: h.ProjectHandler.Delete},
				{Method: "GET", Path: "/:id", Handler: h.ProjectHandler.GetByID},
				{Method: "PUT", Path: "/:id/status", Handler: h.ProjectHandler.ChangeStatus},
			},
			AuthType: shared_http.AuthTypeRequired,
		},
		{
			Prefix: "/projects",
			Handlers: []shared_http.RouteHandler{
				{Method: "GET", Path: "", Handler: h.ProjectHandler.GetPage},
			},
			AuthType: shared_http.AuthTypeRequired,
		},
		{
			Prefix: "/project/code",
			Handlers: []shared_http.RouteHandler{
				{Method: "GET", Path: "/:code", Handler: h.ProjectHandler.GetByCode},
			},
			AuthType: shared_http.AuthTypeRequired,
		},
	}
}
