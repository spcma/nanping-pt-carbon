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

	// 1. 公开路由（不需要认证） - /api/auth/*
	authGroup := group.Group("/auth")
	if mw := middlewares[shared_http.AuthTypeNone]; mw != nil {
		authGroup.Use(mw)
	}
	authGroup.POST("/register", handlers.AuthHandler.Register)
	authGroup.POST("/login", handlers.AuthHandler.Login)

	// 2. 公开查询路由（支持可选认证） - /api/users/*
	usersGroup := group.Group("/users")
	if mw := middlewares[shared_http.AuthTypeOptional]; mw != nil {
		usersGroup.Use(mw)
	}
	usersGroup.GET("", handlers.SysUserHandler.GetPublicPage)

	// 3. 需要认证的系统用户管理路由 - /api/sys/user/*
	sysUserGroup := group.Group("/sys/user")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		sysUserGroup.Use(mw)
	}
	sysUserGroup.POST("", handlers.SysUserHandler.Create)
	sysUserGroup.PUT("/:id", handlers.SysUserHandler.Update)
	sysUserGroup.DELETE("/:id", handlers.SysUserHandler.Delete)
	sysUserGroup.GET("/:id", handlers.SysUserHandler.GetByID)
	sysUserGroup.PUT("/:id/password", handlers.SysUserHandler.ChangePassword)
	sysUserGroup.PUT("/:id/status", handlers.SysUserHandler.ChangeStatus)

	// 4. 系统用户列表路由 - /api/sys/users/*
	sysUsersGroup := group.Group("/sys/users")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		sysUsersGroup.Use(mw)
	}
	sysUsersGroup.GET("", handlers.SysUserHandler.GetPage)

	// 5. 需要认证的系统角色管理路由 - /api/sys/role/*
	sysRoleGroup := group.Group("/sys/role")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		sysRoleGroup.Use(mw)
	}
	sysRoleGroup.POST("", handlers.SysRoleHandler.Create)
	sysRoleGroup.PUT("/:id", handlers.SysRoleHandler.Update)
	sysRoleGroup.DELETE("/:id", handlers.SysRoleHandler.Delete)
	sysRoleGroup.GET("/:id", handlers.SysRoleHandler.GetByID)
	sysRoleGroup.PUT("/:id/status", handlers.SysRoleHandler.ChangeStatus)

	// 6. 系统角色列表路由 - /api/sys/roles/*
	sysRolesGroup := group.Group("/sys/roles")
	if mw := middlewares[shared_http.AuthTypeRequired]; mw != nil {
		sysRolesGroup.Use(mw)
	}
	sysRolesGroup.GET("", handlers.SysRoleHandler.GetPage)
}

// Handlers 包含所有 HTTP 处理器
type Handlers struct {
	SysUserHandler *SysUserHandler
	SysRoleHandler *SysRoleHandler
	AuthHandler    *AuthHandler // 认证处理器
}
