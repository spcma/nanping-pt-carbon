package http

import (
	shared_http "app/internal/shared/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// projectRoutes Project 模块路由注册器
type projectRoutes struct {
	db *gorm.DB
}

// NewProjectRoutes 创建 Project 模块的路由注册器
func NewProjectRoutes(db *gorm.DB) shared_http.RouteRegistry {
	return &projectRoutes{
		db: db,
	}
}

// RegisterRoutes 注册 Project 模块的所有路由（实现 RouteRegistry 接口）
func (r *projectRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {
	// 初始化 DDD 组件
	projectDDD := InitProjectWire(r.db)

	// 创建 handlers
	handlers := &Handlers{
		ProjectHandler: NewProjectHandler(projectDDD.AppService),
	}

	// 项目管理路由 - /api/project/*
	projectGroup := group.Group("/project")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		projectGroup.Use(mw)
	}
	projectGroup.POST("", handlers.ProjectHandler.Create)
	projectGroup.PUT("/:id", handlers.ProjectHandler.Update)
	projectGroup.DELETE("/:id", handlers.ProjectHandler.Delete)
	projectGroup.GET("/:id", handlers.ProjectHandler.GetByID)
	projectGroup.PUT("/:id/status", handlers.ProjectHandler.ChangeStatus)

	// 项目列表路由 - /api/projects/*
	projectsGroup := group.Group("/projects")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		projectsGroup.Use(mw)
	}
	projectsGroup.GET("", handlers.ProjectHandler.GetPage)

	// 项目代码查询路由 - /api/project/code/*
	projectCodeGroup := group.Group("/project/code")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		projectCodeGroup.Use(mw)
	}
	projectCodeGroup.GET("/:code", handlers.ProjectHandler.GetByCode)
}

// Handlers 包含所有 HTTP 处理器
type Handlers struct {
	ProjectHandler *ProjectHandler
}
