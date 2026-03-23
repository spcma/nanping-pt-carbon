package initializer

import (
	shared_http "app/internal/shared/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var dataInitializer *DataInitializer

type Routes struct {
	handler *DataInitializer
}

func NewInitializerRoutes(db *gorm.DB) shared_http.RouteRegistry {
	handler := NewDataInitializer(db)
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
		{
			projectGroup.GET("/project/20260323", i.handler.Add_Project_20260323)
			projectGroup.GET("/methodology/20260323", i.handler.Add_Methodology_20260323)
		}
	}
}
