package transport

import (
	"app/internal/module/methodology/application"
	"app/internal/module/methodology/infrastructure"
	shared_http "app/internal/shared/http"

	"github.com/gin-gonic/gin"
)

// MethodologyRoutes 方法学路由注册器
type MethodologyRoutes struct {
	handler *MethodologyHandler
}

// NewMethodologyRoutes 创建方法学路由注册器
func NewMethodologyRoutes() *MethodologyRoutes {
	// 依赖注入：组装各层
	repo := infrastructure.NewMethodologyRepository()
	appService := application.NewMethodologyAppService(repo)
	handler := NewMethodologyHandler(appService)

	return &MethodologyRoutes{
		handler: handler,
	}
}

// RegisterRoutes 注册路由
func (r *MethodologyRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {
	// 使用认证中间件
	authMiddleware := middlewares[shared_http.AuthTypeRequired]

	// 方法学管理路由 - /api/methodology/*
	methodologyGroup := group.Group("/methodology")
	if authMiddleware != nil {
		methodologyGroup.Use(authMiddleware)
	}
	{
		methodologyGroup.POST("", r.handler.Create)               // 创建方法学
		methodologyGroup.PUT("", r.handler.Update)                // 更新方法学
		methodologyGroup.DELETE("", r.handler.Delete)             // 删除方法学
		methodologyGroup.GET("", r.handler.GetById)               // 根据ID查询
		methodologyGroup.GET("/query", r.handler.GetByQuery)      // 综合查询
		methodologyGroup.PUT("/status", r.handler.ChangeStatus)   // 变更状态
		methodologyGroup.PUT("/activate", r.handler.Activate)     // 启用
		methodologyGroup.PUT("/deactivate", r.handler.Deactivate) // 禁用
	}

	// 方法学列表路由 - /api/methodologies/*
	methodologiesGroup := group.Group("/methodologys")
	if authMiddleware != nil {
		methodologiesGroup.Use(authMiddleware)
	}
	{
		methodologiesGroup.GET("/list", r.handler.GetList) // 列表
		methodologiesGroup.GET("/page", r.handler.GetPage) // 分页
	}
}
