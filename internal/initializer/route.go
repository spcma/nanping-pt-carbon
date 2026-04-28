package initializer

import (
	"app/internal/shared/db"
	shared_http "app/internal/shared/http"

	"github.com/gin-gonic/gin"
)

var dataInitializer *DataInitializer

type Routes struct {
	handler *DataInitializer
}

func NewInitializerRoutes() shared_http.RouteRegistry {
	dbInst := db.Default()
	handler := NewDataInitializer(dbInst)
	return &Routes{
		handler: handler,
	}
}

func (i *Routes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {
	//	统一认证中间件
	authTypeRequiredRoute := group.Group("")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		authTypeRequiredRoute.Use(mw)
	}
	{
		// 项目管理路由 - /api/initializer/*
		projectGroup := authTypeRequiredRoute.Group("/initializer")
		_ = projectGroup
	}
}
