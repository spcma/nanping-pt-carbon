package http

import (
	"app/internal/module/iam/application"
	"app/internal/module/iam/infrastructure"
	shared_http "app/internal/shared/http"
	"app/internal/shared/token"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// iamRoutes IAM 模块路由注册器
type iamRoutes struct {
	usersHandler *UsersHandler
	rolesHandler *RolesHandler
	authHandler  *AuthHandler
}

// NewIAMRoutes 创建 IAM 模块的路由注册器
func NewIAMRoutes(db *gorm.DB, tokenManager token.Manager) shared_http.RouteRegistry {
	//	初始化 users 模块
	usersRepo := infrastructure.NewUserRepository(db)
	usersService := application.NewUsersService(usersRepo)
	usersHandler := NewUsersHandler(usersService)

	//	初始化 roles 模块
	roleRepo := infrastructure.NewRoleRepository(db)
	roleService := application.NewSysRoleAppService(roleRepo)
	rolesHandler := NewRolesHandler(roleService)

	//	初始化 auth 模块
	authHandler := NewAuthHandler(usersService, tokenManager)

	return &iamRoutes{
		usersHandler: usersHandler,
		rolesHandler: rolesHandler,
		authHandler:  authHandler,
	}
}

// RegisterRoutes 注册 IAM 模块的所有路由（实现 RouteRegistry 接口）
func (r *iamRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {

	// 1. 公开路由（不需要认证） - /api/*
	publicRoute := group.Group("")
	if mw := middlewares[shared_http.AuthTypeNone]; mw != nil {
		publicRoute.Use(mw)
	}
	{
		publicRoute.POST("/register", r.authHandler.Register)
		publicRoute.POST("/login", r.authHandler.Login)
	}

	// 2. 公开查询路由（支持可选认证） - /api/users/*
	optionalAuthRoute := group.Group("/users")
	if mw := middlewares[shared_http.AuthTypeOptional]; mw != nil {
		optionalAuthRoute.Use(mw)
	}
	{
		optionalAuthRoute.GET("/pub/page", r.usersHandler.GetPublicPage)
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
			sysUserGroup.POST("", r.usersHandler.Create)
			sysUserGroup.PUT("", r.usersHandler.Update)
			sysUserGroup.DELETE("", r.usersHandler.Delete)
			sysUserGroup.GET("", r.usersHandler.GetById)         // 仅 ID 查询
			sysUserGroup.GET("query", r.usersHandler.GetByQuery) // 综合查询
			sysUserGroup.PUT("/password", r.usersHandler.ChangePassword)
			sysUserGroup.PUT("/status", r.usersHandler.ChangeStatus)
		}

		// 系统用户列表 - /api/users/*
		sysUsersGroup := authTypeRequiredRoute.Group("/users")
		{
			sysUsersGroup.GET("list", r.usersHandler.GetList)
			sysUsersGroup.GET("page", r.usersHandler.GetPage)
		}

		// 系统角色管理 - /api/roles/*
		sysRoleGroup := authTypeRequiredRoute.Group("/roles")
		{
			sysRoleGroup.POST("", r.rolesHandler.Create)
			sysRoleGroup.PUT("", r.rolesHandler.Update)
			sysRoleGroup.DELETE("", r.rolesHandler.Delete)
			sysRoleGroup.GET("", r.rolesHandler.GetByID) // 仅 ID 查询
			//sysRoleGroup.GET("query", r.rolesHandler.GetByQuery) // 综合查询
			sysRoleGroup.PUT("/status", r.rolesHandler.ChangeStatus)
		}

		// 系统角色列表 - /api/roles/*
		sysRolesGroup := authTypeRequiredRoute.Group("/roles")
		{
			sysRolesGroup.GET("page", r.rolesHandler.GetPage)
		}
	}
}

// Handlers 包含所有 HTTP 处理器
type Handlers struct {
	SysUserHandler *UsersHandler
	SysRoleHandler *RolesHandler
	AuthHandler    *AuthHandler // 认证处理器
}
