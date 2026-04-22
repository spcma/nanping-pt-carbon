package http

import (
	"app/internal/module/iam/application"
	"app/internal/module/iam/infrastructure"
	"app/internal/shared/db"
	shared_http "app/internal/shared/http"
	"app/internal/shared/token"

	"github.com/gin-gonic/gin"
)

// iamRoutes IAM 模块路由注册器
type iamRoutes struct {
	userHandler *UserHandler
	roleHandler *RoleHandler
	authHandler *AuthHandler
}

// NewIAMRoutes 创建 IAM 模块的路由注册器
func NewIAMRoutes() shared_http.RouteRegistry {
	dbInst := db.Default()
	tokenManager := token.Default()

	//	初始化 users 模块
	usersRepo := infrastructure.NewUserRepository(dbInst)
	usersService := application.NewUserService(usersRepo)
	userHandler := NewUserHandler(usersService)

	//	初始化 roles 模块
	roleRepo := infrastructure.NewRoleRepository(dbInst)
	roleService := application.NewRoleAppService(roleRepo)
	roleHandler := NewRoleHandler(roleService)

	//	初始化 auth 模块
	authHandler := NewAuthHandler(usersService, tokenManager)

	return &iamRoutes{
		userHandler: userHandler,
		roleHandler: roleHandler,
		authHandler: authHandler,
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
		optionalAuthRoute.GET("/pub/page", r.userHandler.GetPublicPage)
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
			sysUserGroup.POST("", r.userHandler.Create)
			sysUserGroup.PUT("", r.userHandler.Update)
			sysUserGroup.DELETE("", r.userHandler.Delete)
			sysUserGroup.GET("", r.userHandler.GetById)         // 仅 ID 查询
			sysUserGroup.GET("query", r.userHandler.GetByQuery) // 综合查询
			sysUserGroup.PUT("/password", r.userHandler.ChangePassword)
			sysUserGroup.PUT("/password/reset", r.userHandler.ResetPassword)
			sysUserGroup.PUT("/status", r.userHandler.ChangeStatus)
		}

		// 系统用户列表 - /api/users/*
		sysUsersGroup := authTypeRequiredRoute.Group("/users")
		{
			sysUsersGroup.GET("list", r.userHandler.GetList)
			sysUsersGroup.GET("page", r.userHandler.GetPage)
		}

		// 系统角色管理 - /api/roles/*
		sysRoleGroup := authTypeRequiredRoute.Group("/roles")
		{
			sysRoleGroup.POST("", r.roleHandler.Create)
			sysRoleGroup.PUT("", r.roleHandler.Update)
			sysRoleGroup.DELETE("", r.roleHandler.Delete)
			sysRoleGroup.GET("", r.roleHandler.GetByID) // 仅 ID 查询
			//sysRoleGroup.GET("query", r.rolesHandler.GetByQuery) // 综合查询
			sysRoleGroup.PUT("/status", r.roleHandler.ChangeStatus)
		}

		// 系统角色列表 - /api/roles/*
		sysRolesGroup := authTypeRequiredRoute.Group("/roles")
		{
			sysRolesGroup.GET("page", r.roleHandler.GetPage)
		}
	}
}

// Handlers 包含所有 HTTP 处理器
type Handlers struct {
	UserHandler *UserHandler
	RoleHandler *RoleHandler
	AuthHandler *AuthHandler // 认证处理器
}
