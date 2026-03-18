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
	authGroup := group.Group("")
	if mw := middlewares[shared_http.AuthTypeNone]; mw != nil {
		authGroup.Use(mw)
	}
	{
		authGroup.POST("/register", handlers.AuthHandler.Register)
		authGroup.POST("/login", handlers.AuthHandler.Login)
	}

	// 2. 公开查询路由（支持可选认证） - /api/users/*
	usersGroup := group.Group("/users")
	if mw := middlewares[shared_http.AuthTypeOptional]; mw != nil {
		usersGroup.Use(mw)
	}
	{
		usersGroup.GET("/pub/page", handlers.SysUserHandler.GetPublicPage)
	}

	// 3. 需要认证的系统用户管理路由 - /api/user/*
	sysUserGroup := group.Group("/user")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		sysUserGroup.Use(mw)
	}
	{
		sysUserGroup.POST("", handlers.SysUserHandler.Create)
		sysUserGroup.PUT("", handlers.SysUserHandler.Update)
		sysUserGroup.DELETE("", handlers.SysUserHandler.Delete)
		sysUserGroup.GET("", handlers.SysUserHandler.GetByCond)
		sysUserGroup.PUT("/password", handlers.SysUserHandler.ChangePassword)
		sysUserGroup.PUT("/status", handlers.SysUserHandler.ChangeStatus)
	}

	sysUsersGroup := group.Group("/users")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		sysUserGroup.Use(mw)
	}
	{
		sysUsersGroup.GET("list", handlers.SysUserHandler.GetList)
		sysUsersGroup.GET("page", handlers.SysUserHandler.GetPage)
	}

	// 5. 需要认证的系统角色管理路由 - /api/roles/*
	sysRoleGroup := group.Group("/roles")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		sysRoleGroup.Use(mw)
	}
	{
		sysRoleGroup.POST("", handlers.SysRoleHandler.Create)
		sysRoleGroup.PUT("/:id", handlers.SysRoleHandler.Update)
		sysRoleGroup.DELETE("/:id", handlers.SysRoleHandler.Delete)
		sysRoleGroup.GET("/:id", handlers.SysRoleHandler.GetByID)
		sysRoleGroup.PUT("/:id/status", handlers.SysRoleHandler.ChangeStatus)
	}

	// 6. 系统角色列表路由 - /api/roles/*
	sysRolesGroup := group.Group("/roles")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		sysRolesGroup.Use(mw)
	}
	sysRolesGroup.GET("page", handlers.SysRoleHandler.GetPage)
}

// Handlers 包含所有 HTTP 处理器
type Handlers struct {
	SysUserHandler *SysUserHandler
	SysRoleHandler *SysRoleHandler
	AuthHandler    *AuthHandler // 认证处理器
}
