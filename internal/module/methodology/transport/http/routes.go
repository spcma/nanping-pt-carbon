package http

import (
	shared_http "app/internal/shared/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// methodologyRoutes Methodology 模块路由注册器
type methodologyRoutes struct {
	db *gorm.DB
}

// NewMethodologyRoutes 创建 Methodology 模块的路由注册器
func NewMethodologyRoutes(db *gorm.DB) shared_http.RouteRegistry {
	return &methodologyRoutes{
		db: db,
	}
}

// RegisterRoutes 注册 Methodology 模块的所有路由（实现 RouteRegistry 接口）
func (r *methodologyRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {
	// 初始化 wire 组件
	methodologyWire := InitMethodologyWire(r.db)

	// 创建 handlers
	methodologyHandler := NewMethodologyHandler(methodologyWire.AppService)

	// 统一认证中间件
	authTypeRequiredRoute := group.Group("")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		authTypeRequiredRoute.Use(mw)
	}
	{
		// 方法学管理路由 - /api/methodology/*
		methodologyGroup := authTypeRequiredRoute.Group("/methodology")
		{
			methodologyGroup.POST("", methodologyHandler.Create)
			methodologyGroup.PUT("", methodologyHandler.Update)
			methodologyGroup.DELETE("", methodologyHandler.Delete)
			methodologyGroup.GET("", methodologyHandler.GetById)         // 仅 ID 查询
			methodologyGroup.GET("query", methodologyHandler.GetByQuery) // 综合查询
			methodologyGroup.PUT("status", methodologyHandler.ChangeStatus)
		}

		// 方法学列表路由 - /api/methodologies/*
		methodologiesGroup := authTypeRequiredRoute.Group("/methodologies")
		{
			methodologiesGroup.GET("list", methodologyHandler.GetList)
			methodologiesGroup.GET("page", methodologyHandler.GetPage)
		}
	}
}
