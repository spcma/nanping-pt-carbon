package http

import (
	shared_http "app/internal/shared/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// projectRoutes Project 模块路由注册器
type projectRoutes struct {
	projectHandler        *ProjectHandler
	projectMembersHandler *ProjectMembersHandler
}

// NewProjectRoutes 创建 Project 模块的路由注册器
func NewProjectRoutes(db *gorm.DB) shared_http.RouteRegistry {
	//	初始化application
	projectWire := InitProjectWire(db)
	projectMembersWire := InitProjectMembersWire(db)

	//	初始化handler
	projectHandler := NewProjectHandler(projectWire.Service)
	projectMembersHandler := NewProjectMembersHandler(projectMembersWire.Service)

	return &projectRoutes{
		projectHandler:        projectHandler,
		projectMembersHandler: projectMembersHandler,
	}
}

// RegisterRoutes 注册 Project 模块的所有路由（实现 RouteRegistry 接口）
func (r *projectRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {

	//	统一认证中间件
	authTypeRequiredRoute := group.Group("")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		authTypeRequiredRoute.Use(mw)
	}
	{
		// 项目管理路由 - /api/project/*
		projectGroup := authTypeRequiredRoute.Group("/project")
		{
			projectGroup.POST("", r.projectHandler.Create)
			projectGroup.PUT("", r.projectHandler.Update)
			projectGroup.DELETE("", r.projectHandler.Delete)
			projectGroup.GET("", r.projectHandler.GetById)         // 仅ID查询
			projectGroup.GET("query", r.projectHandler.GetByQuery) // 综合查询
			projectGroup.PUT("status", r.projectHandler.ChangeStatus)
		}
		// 项目列表路由 - /api/projects/*
		projectsGroup := authTypeRequiredRoute.Group("/projects")
		{
			projectsGroup.GET("list", r.projectHandler.GetList)
			projectsGroup.GET("page", r.projectHandler.GetPage)
		}
		// 项目成员管理路由 - /api/projectMember/*
		projectMemberGroup := authTypeRequiredRoute.Group("/projectMember")
		{
			projectMemberGroup.POST("", r.projectMembersHandler.Create)
			projectMemberGroup.PUT("", r.projectMembersHandler.Update)
			projectMemberGroup.DELETE("", r.projectMembersHandler.Delete)
		}
		// 项目成员列表路由 - /api/projectMembers/*
		projectMembersGroup := authTypeRequiredRoute.Group("/projectMembers")
		{
			projectMembersGroup.GET("page", r.projectMembersHandler.GetPage)
			projectMembersGroup.GET("list", r.projectMembersHandler.GetList)
		}
	}
}
