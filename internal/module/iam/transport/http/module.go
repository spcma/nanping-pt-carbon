package http

import (
	"app/internal/module/iam/wire"
	shared_http "app/internal/shared/http"
	"app/internal/shared/token"
	"gorm.io/gorm"
)

// IAMModule IAM 模块实现
type IAMModule struct {
	SysUserDDD     *wire.SysUserDDD
	SysRoleDDD     *wire.SysRoleDDD
	SysApiDDD      *wire.SysApiDDD
	SysUserRoleDDD *wire.SysUserRoleDDD
}

// Name 返回模块名称
func (m *IAMModule) Name() string {
	return "IAM"
}

// InitWithDeps 初始化模块（支持外部依赖注入）
func (m *IAMModule) InitWithDeps(db *gorm.DB, deps ...any) any {
	// 从依赖列表中提取 jwtManager
	var jwtManager *token.JWTManager
	for _, dep := range deps {
		if jm, ok := dep.(*token.JWTManager); ok {
			jwtManager = jm
			break
		}
	}

	m.SysUserDDD = wire.InitSysUserDDD(db)
	m.SysRoleDDD = wire.InitSysRoleDDD(db)
	m.SysApiDDD = wire.InitSysApiDDD(db)
	m.SysUserRoleDDD = wire.InitSysUserRoleDDD(db)

	// 创建并返回 Handlers 实例，注入 jwtManager
	return &Handlers{
		SysUserHandler: NewSysUserHandler(m.SysUserDDD.AppService),
		SysRoleHandler: NewSysRoleHandler(m.SysRoleDDD.AppService),
		AuthHandler:    NewAuthHandler(m.SysUserDDD.AppService, jwtManager),
	}
}

// Init 初始化模块（默认实现，用于兼容性）
func (m *IAMModule) Init(db *gorm.DB) any {
	return m.InitWithDeps(db)
}

// RegisterRoutes 注册路由（由 ModuleManager 调用）
func (m *IAMModule) RegisterRoutes(handlers any) []shared_http.RouteGroupConfig {
	h, ok := handlers.(*Handlers)
	if !ok {
		return nil
	}
	return h.RegisterRoutes()
}
