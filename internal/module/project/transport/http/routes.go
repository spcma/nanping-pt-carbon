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
	projectWire := InitProjectWire(r.db)
	projectMembersWire := InitProjectMembersWire(r.db)

	// 创建 handlers
	handlers := &Handlers{
		ProjectHandler:        NewProjectHandler(projectWire.AppService),
		ProjectMembersHandler: NewProjectMembersHandler(projectMembersWire.AppService),
	}

	// 项目管理路由 - /api/project/*
	projectGroup := group.Group("/project")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		projectGroup.Use(mw)
	}
	projectGroup.POST("", handlers.ProjectHandler.Create)
	projectGroup.PUT("", handlers.ProjectHandler.Update)
	projectGroup.DELETE("", handlers.ProjectHandler.Delete)
	projectGroup.GET("", handlers.ProjectHandler.GetByCond)
	projectGroup.PUT("status", handlers.ProjectHandler.ChangeStatus)

	// 项目列表路由 - /api/projects/*
	projectsGroup := group.Group("/projects")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		projectsGroup.Use(mw)
	}
	projectsGroup.GET("page", handlers.ProjectHandler.GetPage)

	// 项目成员管理路由 - /api/projectMember/*
	projectMemberGroup := group.Group("/projectMember")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		projectMemberGroup.Use(mw)
	}
	projectMemberGroup.POST("", handlers.ProjectMembersHandler.Create)
	projectMemberGroup.PUT("", handlers.ProjectMembersHandler.Update)
	projectMemberGroup.DELETE("", handlers.ProjectMembersHandler.Delete)

	// 项目成员列表路由 - /api/projectMembers/*
	projectMembersGroup := group.Group("/projectMembers")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		projectMembersGroup.Use(mw)
	}
	projectMembersGroup.GET("page", handlers.ProjectMembersHandler.GetPage)
	projectMembersGroup.GET("list", handlers.ProjectMembersHandler.GetList)
}

// Handlers 包含所有 HTTP 处理器
type Handlers struct {
	ProjectHandler        *ProjectHandler
	ProjectMembersHandler *ProjectMembersHandler
}
