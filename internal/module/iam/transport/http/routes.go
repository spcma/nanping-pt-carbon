package http

import (
	shared_http "app/internal/shared/http"
	"app/internal/shared/token"

	"gorm.io/gorm"
)

// RegisterRoutes 注册 IAM 模块的所有路由
func RegisterRoutes(db *gorm.DB, jwtManager token.Manager) []shared_http.RouteGroupConfig {
	sysUserWire := InitSysUserWire(db)
	sysRoleWire := InitSysRoleWire(db)

	handlers := &Handlers{
		SysUserHandler: NewSysUserHandler(sysUserWire.Service),
		SysRoleHandler: NewSysRoleHandler(sysRoleWire.Service),
		AuthHandler:    NewAuthHandler(sysUserWire.Service, jwtManager),
	}

	return handlers.RegisterRoutes()
}

// Handlers 包含所有 HTTP 处理器
type Handlers struct {
	SysUserHandler *SysUserHandler
	SysRoleHandler *SysRoleHandler
	AuthHandler    *AuthHandler // 认证处理器
}

// RegisterRoutes 注册 IAM 模块的所有路由（内部方法）
func (h *Handlers) RegisterRoutes() []shared_http.RouteGroupConfig {
	return []shared_http.RouteGroupConfig{
		// 1. 公开路由（不需要认证）
		{
			Prefix: "/auth",
			Handlers: []shared_http.RouteHandler{
				{Method: "POST", Path: "/register", Handler: h.AuthHandler.Register},
				{Method: "POST", Path: "/login", Handler: h.AuthHandler.Login},
			},
			AuthType: shared_http.AuthTypeNone,
		},
		// 2. 公开查询路由（支持可选认证）
		{
			Prefix: "/users",
			Handlers: []shared_http.RouteHandler{
				{Method: "GET", Path: "", Handler: h.SysUserHandler.GetPublicPage},
			},
			AuthType: shared_http.AuthTypeOptional,
		},
		// 3. 需要认证的系统用户管理路由
		{
			Prefix: "/sys/user",
			Handlers: []shared_http.RouteHandler{
				{Method: "POST", Path: "", Handler: h.SysUserHandler.Create},
				{Method: "PUT", Path: "/:id", Handler: h.SysUserHandler.Update},
				{Method: "DELETE", Path: "/:id", Handler: h.SysUserHandler.Delete},
				{Method: "GET", Path: "/:id", Handler: h.SysUserHandler.GetByID},
				{Method: "PUT", Path: "/:id/password", Handler: h.SysUserHandler.ChangePassword},
				{Method: "PUT", Path: "/:id/status", Handler: h.SysUserHandler.ChangeStatus},
			},
			AuthType: shared_http.AuthTypeRequired,
		},
		{
			Prefix: "/sys/users",
			Handlers: []shared_http.RouteHandler{
				{Method: "GET", Path: "", Handler: h.SysUserHandler.GetPage},
			},
			AuthType: shared_http.AuthTypeRequired,
		},
		// 4. 需要认证的系统角色管理路由
		{
			Prefix: "/sys/role",
			Handlers: []shared_http.RouteHandler{
				{Method: "POST", Path: "", Handler: h.SysRoleHandler.Create},
				{Method: "PUT", Path: "/:id", Handler: h.SysRoleHandler.Update},
				{Method: "DELETE", Path: "/:id", Handler: h.SysRoleHandler.Delete},
				{Method: "GET", Path: "/:id", Handler: h.SysRoleHandler.GetByID},
				{Method: "PUT", Path: "/:id/status", Handler: h.SysRoleHandler.ChangeStatus},
			},
			AuthType: shared_http.AuthTypeRequired,
		},
		{
			Prefix: "/sys/roles",
			Handlers: []shared_http.RouteHandler{
				{Method: "GET", Path: "", Handler: h.SysRoleHandler.GetPage},
			},
			AuthType: shared_http.AuthTypeRequired,
		},
	}
}
