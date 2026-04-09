package http

import (
	"app/internal/module/project/application"
	"app/internal/module/project/infrastructure"
	"app/internal/shared/db"
	shared_http "app/internal/shared/http"

	"github.com/gin-gonic/gin"
)

// projectRoutes Project 模块路由注册器
type projectRoutes struct {
	projectHandler        *ProjectHandler
	projectMembersHandler *ProjectMembersHandler
}

// NewProjectRoutes 创建 Project 模块的路由注册器
func NewProjectRoutes() shared_http.RouteRegistry {
	dbInst := db.Default()

	//	初始化 project 模块
	projectRepo := infrastructure.NewProjectRepository(dbInst)
	appService := application.NewProjectService(projectRepo)
	projectHandler := NewProjectHandler(appService)

	//	初始化 projectMembers 模块
	projectMembersRepo := infrastructure.NewProjectMembersRepository(dbInst)
	projectMembersAppService := application.NewProjectMembersService(projectMembersRepo)
	projectMembersHandler := NewProjectMembersHandler(projectMembersAppService)

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
