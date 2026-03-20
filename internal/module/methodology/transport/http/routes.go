package http

import (
	"app/internal/module/methodology/application"
	"app/internal/module/methodology/infrastructure"
	shared_http "app/internal/shared/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// methodologyRoutes Methodology 模块路由注册器
type methodologyRoutes struct {
	methodologyHandler *MethodologyHandler
}

// NewMethodologyRoutes 创建 Methodology 模块的路由注册器
func NewMethodologyRoutes(db *gorm.DB) shared_http.RouteRegistry {
	repo := infrastructure.NewMethodologyRepository(db)
	appService := application.NewMethodologyAppService(repo)
	methodologyHandler := NewMethodologyHandler(appService)

	return &methodologyRoutes{
		methodologyHandler: methodologyHandler,
	}
}

// RegisterRoutes 注册 Methodology 模块的所有路由（实现 RouteRegistry 接口）
func (r *methodologyRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {

	// 统一认证中间件
	authTypeRequiredRoute := group.Group("")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		authTypeRequiredRoute.Use(mw)
	}
	{
		// 方法学管理路由 - /api/methodology/*
		methodologyGroup := authTypeRequiredRoute.Group("/methodology")
		{
			methodologyGroup.POST("", r.methodologyHandler.Create)
			methodologyGroup.PUT("", r.methodologyHandler.Update)
			methodologyGroup.DELETE("", r.methodologyHandler.Delete)
			methodologyGroup.GET("", r.methodologyHandler.GetById)         // 仅 ID 查询
			methodologyGroup.GET("query", r.methodologyHandler.GetByQuery) // 综合查询
			methodologyGroup.PUT("status", r.methodologyHandler.ChangeStatus)
		}

		// 方法学列表路由 - /api/methodologies/*
		methodologiesGroup := authTypeRequiredRoute.Group("/methodologies")
		{
			methodologiesGroup.GET("list", r.methodologyHandler.GetList)
			methodologiesGroup.GET("page", r.methodologyHandler.GetPage)
		}
	}
}
