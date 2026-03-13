package http

import (
	"app/internal/module/methodology/wire"
	shared_http "app/internal/shared/http"

	"gorm.io/gorm"
)

// MethodologyModule 方法学模块实现
type MethodologyModule struct {
	MethodologyDDD *wire.MethodologyDDD
}

// Name 返回模块名称
func (m *MethodologyModule) Name() string {
	return "Methodology"
}

// InitWithDeps 初始化模块（支持外部依赖注入）
func (m *MethodologyModule) InitWithDeps(db *gorm.DB, deps ...any) any {
	m.MethodologyDDD = wire.InitMethodologyDDD(db)

	// 创建并返回 Handlers 实例
	return &Handlers{
		MethodologyHandler: NewMethodologyHandler(m.MethodologyDDD.AppService),
	}
}

// Init 初始化模块（默认实现，用于兼容性）
func (m *MethodologyModule) Init(db *gorm.DB) any {
	return m.InitWithDeps(db)
}

// RegisterRoutes 注册路由（由 ModuleManager 调用）
func (m *MethodologyModule) RegisterRoutes(handlers any) []shared_http.RouteGroupConfig {
	h, ok := handlers.(*Handlers)
	if !ok {
		return nil
	}
	return h.RegisterRoutes()
}
