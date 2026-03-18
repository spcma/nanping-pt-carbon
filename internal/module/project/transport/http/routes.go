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

// Handlers 包含所有 HTTP 处理器
type Handlers struct {
	ProjectHandler        *ProjectHandler
	ProjectMembersHandler *ProjectMembersHandler
}

// RegisterRoutes 注册 Project 模块的所有路由（实现 RouteRegistry 接口）
func (r *projectRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {
	//	初始化application
	projectWire := InitProjectWire(r.db)
	projectMembersWire := InitProjectMembersWire(r.db)

	//	初始化handler
	projectHandler := NewProjectHandler(projectWire.Service)
	projectMembersHandler := NewProjectMembersHandler(projectMembersWire.Service)

	//	统一认证中间件
	authTypeRequiredRoute := group.Group("")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		authTypeRequiredRoute.Use(mw)
	}
	{
		// 项目管理路由 - /api/project/*
		projectGroup := authTypeRequiredRoute.Group("/project")
		{
			projectGroup.POST("", projectHandler.Create)
			projectGroup.PUT("", projectHandler.Update)
			projectGroup.DELETE("", projectHandler.Delete)
			projectGroup.GET("", projectHandler.GetById)         // 仅ID查询
			projectGroup.GET("query", projectHandler.GetByQuery) // 综合查询
			projectGroup.PUT("status", projectHandler.ChangeStatus)
		}
		// 项目列表路由 - /api/projects/*
		projectsGroup := authTypeRequiredRoute.Group("/projects")
		{
			projectsGroup.GET("list", projectMembersHandler.GetList)
			projectsGroup.GET("page", projectMembersHandler.GetPage)
		}
		// 项目成员管理路由 - /api/projectMember/*
		projectMemberGroup := authTypeRequiredRoute.Group("/projectMember")
		{
			projectMemberGroup.POST("", projectMembersHandler.Create)
			projectMemberGroup.PUT("", projectMembersHandler.Update)
			projectMemberGroup.DELETE("", projectMembersHandler.Delete)
		}
		// 项目成员列表路由 - /api/projectMembers/*
		projectMembersGroup := authTypeRequiredRoute.Group("/projectMembers")
		{
			projectMembersGroup.GET("page", projectMembersHandler.GetPage)
			projectMembersGroup.GET("list", projectMembersHandler.GetList)
		}
	}
}
