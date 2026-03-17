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
	// 初始化 DDD 组件
	methodologyDDD := InitMethodologyDDD(r.db)

	// 创建 handlers
	handlers := &Handlers{
		MethodologyHandler: NewMethodologyHandler(methodologyDDD.AppService),
	}

	// 方法学管理路由 - /api/methodology/*
	methodologyGroup := group.Group("/methodology")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		methodologyGroup.Use(mw)
	}
	methodologyGroup.POST("", handlers.MethodologyHandler.Create)
	methodologyGroup.PUT("/:id", handlers.MethodologyHandler.Update)
	methodologyGroup.DELETE("/:id", handlers.MethodologyHandler.Delete)
	methodologyGroup.GET("/:id", handlers.MethodologyHandler.GetByID)
	methodologyGroup.PUT("/:id/status", handlers.MethodologyHandler.ChangeStatus)

	// 方法学列表路由 - /api/methodologies/*
	methodologiesGroup := group.Group("/methodologies")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		methodologiesGroup.Use(mw)
	}
	methodologiesGroup.GET("", handlers.MethodologyHandler.GetPage)

	// 方法学代码查询路由 - /api/methodology/code/*
	methodologyCodeGroup := group.Group("/methodology/code")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		methodologyCodeGroup.Use(mw)
	}
	methodologyCodeGroup.GET("/:code", handlers.MethodologyHandler.GetByCode)
}

// Handlers 包含所有 HTTP 处理器
type Handlers struct {
	MethodologyHandler *MethodologyHandler
}
