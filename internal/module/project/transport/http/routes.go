package transport

import (
	"app/internal/module/project/application"
	"app/internal/module/project/infrastructure"
	shared_http "app/internal/shared/http"

	"github.com/gin-gonic/gin"
)

// ProjectRoutes 项目路由注册器
type ProjectRoutes struct {
	handler       *ProjectHandler
	memberHandler *ProjectMembersHandler
}

// NewProjectRoutes 创建项目路由注册器
func NewProjectRoutes() *ProjectRoutes {
	// 依赖注入：组装各层
	repo := infrastructure.NewProjectRepository()
	appService := application.NewProjectAppService(repo)
	handler := NewProjectHandler(appService)

	// 项目成员依赖注入
	memberRepo := infrastructure.NewProjectMembersRepository()
	memberAppService := application.NewProjectMembersAppService(memberRepo)
	memberHandler := NewProjectMembersHandler(memberAppService)

	return &ProjectRoutes{
		handler:       handler,
		memberHandler: memberHandler,
	}
}

// RegisterRoutes 注册路由
func (r *ProjectRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {
	// 项目路由组
	projectGroup := group.Group("/projects")

	// 使用认证中间件
	authMiddleware := middlewares[shared_http.AuthTypeRequired]
	if authMiddleware != nil {
		projectGroup.Use(authMiddleware)
	}

	// 注册项目路由
	projectGroup.POST("", r.handler.Create)             // 创建项目
	projectGroup.PUT("", r.handler.Update)              // 更新项目
	projectGroup.DELETE("", r.handler.Delete)           // 删除项目
	projectGroup.GET("/id", r.handler.GetByID)          // 根据ID查询
	projectGroup.GET("/query", r.handler.GetByQuery)    // 条件查询
	projectGroup.GET("/list", r.handler.GetList)        // 列表
	projectGroup.GET("/page", r.handler.GetPage)        // 分页
	projectGroup.PUT("/status", r.handler.ChangeStatus) // 变更状态
	projectGroup.PUT("/activate", r.handler.Activate)   // 激活项目
	projectGroup.PUT("/complete", r.handler.Complete)   // 完成项目
	projectGroup.PUT("/cancel", r.handler.Cancel)       // 取消项目

	// 项目成员路由组
	memberGroup := group.Group("/project-members")
	if authMiddleware != nil {
		memberGroup.Use(authMiddleware)
	}

	// 注册项目成员路由
	memberGroup.POST("", r.memberHandler.Create)          // 添加成员
	memberGroup.PUT("", r.memberHandler.Update)           // 更新成员
	memberGroup.DELETE("", r.memberHandler.Delete)        // 删除成员
	memberGroup.GET("/list", r.memberHandler.GetList)     // 成员列表
	memberGroup.GET("/user", r.memberHandler.GetByUserID) // 用户参与的项目
	memberGroup.GET("/page", r.memberHandler.GetPage)     // 分页查询
}
