package http

import (
	shared_http "app/internal/shared/http"
	"app/internal/shared/token"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// iamRoutes IAM 模块路由注册器
type iamRoutes struct {
	db           *gorm.DB
	tokenManager token.Manager
}

// NewIAMRoutes 创建 IAM 模块的路由注册器
func NewIAMRoutes(db *gorm.DB, tokenManager token.Manager) shared_http.RouteRegistry {
	return &iamRoutes{
		db:           db,
		tokenManager: tokenManager,
	}
}

// RegisterRoutes 注册 IAM 模块的所有路由（实现 RouteRegistry 接口）
func (r *iamRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {
	sysUserWire := InitSysUserWire(r.db)
	sysRoleWire := InitSysRoleWire(r.db)

	handlers := &Handlers{
		SysUserHandler: NewSysUserHandler(sysUserWire.Service),
		SysRoleHandler: NewSysRoleHandler(sysRoleWire.Service),
		AuthHandler:    NewAuthHandler(sysUserWire.Service, r.tokenManager),
	}

	// 1. 公开路由（不需要认证） - /api/*
	publicRoute := group.Group("")
	if mw := middlewares[shared_http.AuthTypeNone]; mw != nil {
		publicRoute.Use(mw)
	}
	{
		publicRoute.POST("/register", handlers.AuthHandler.Register)
		publicRoute.POST("/login", handlers.AuthHandler.Login)
	}

	// 2. 公开查询路由（支持可选认证） - /api/users/*
	optionalAuthRoute := group.Group("/users")
	if mw := middlewares[shared_http.AuthTypeOptional]; mw != nil {
		optionalAuthRoute.Use(mw)
	}
	{
		optionalAuthRoute.GET("/pub/page", handlers.SysUserHandler.GetPublicPage)
	}

	// 3. 需要认证的系统用户管理路由 - /api/user/*
	authTypeRequiredRoute := group.Group("")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		authTypeRequiredRoute.Use(mw)
	}
	{
		// 系统用户管理 - /api/user/*
		sysUserGroup := authTypeRequiredRoute.Group("/user")
		{
			sysUserGroup.POST("", handlers.SysUserHandler.Create)
			sysUserGroup.PUT("", handlers.SysUserHandler.Update)
			sysUserGroup.DELETE("", handlers.SysUserHandler.Delete)
			sysUserGroup.GET("", handlers.SysUserHandler.GetById)         // 仅 ID 查询
			sysUserGroup.GET("query", handlers.SysUserHandler.GetByQuery) // 综合查询
			sysUserGroup.PUT("/password", handlers.SysUserHandler.ChangePassword)
			sysUserGroup.PUT("/status", handlers.SysUserHandler.ChangeStatus)
		}

		// 系统用户列表 - /api/users/*
		sysUsersGroup := authTypeRequiredRoute.Group("/users")
		{
			sysUsersGroup.GET("list", handlers.SysUserHandler.GetList)
			sysUsersGroup.GET("page", handlers.SysUserHandler.GetPage)
		}

		// 系统角色管理 - /api/roles/*
		sysRoleGroup := authTypeRequiredRoute.Group("/roles")
		{
			sysRoleGroup.POST("", handlers.SysRoleHandler.Create)
			sysRoleGroup.PUT("", handlers.SysRoleHandler.Update)
			sysRoleGroup.DELETE("", handlers.SysRoleHandler.Delete)
			sysRoleGroup.GET("", handlers.SysRoleHandler.GetById)         // 仅 ID 查询
			sysRoleGroup.GET("query", handlers.SysRoleHandler.GetByQuery) // 综合查询
			sysRoleGroup.PUT("/status", handlers.SysRoleHandler.ChangeStatus)
		}

		// 系统角色列表 - /api/roles/*
		sysRolesGroup := authTypeRequiredRoute.Group("/roles")
		{
			sysRolesGroup.GET("page", handlers.SysRoleHandler.GetPage)
		}
	}
}

// Handlers 包含所有 HTTP 处理器
type Handlers struct {
	SysUserHandler *SysUserHandler
	SysRoleHandler *SysRoleHandler
	AuthHandler    *AuthHandler // 认证处理器
}
