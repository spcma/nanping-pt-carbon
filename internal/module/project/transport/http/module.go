package http

import (
	"app/internal/module/project/wire"
	shared_http "app/internal/shared/http"

	"gorm.io/gorm"
)

// ProjectModule 项目模块实现
type ProjectModule struct {
	ProjectDDD *wire.ProjectDDD
}

// Name 返回模块名称
func (m *ProjectModule) Name() string {
	return "Project"
}

// InitWithDeps 初始化模块（支持外部依赖注入）
func (m *ProjectModule) InitWithDeps(db *gorm.DB, deps ...any) any {
	m.ProjectDDD = wire.InitProjectDDD(db)

	// 创建并返回 Handlers 实例
	return &Handlers{
		ProjectHandler: NewProjectHandler(m.ProjectDDD.AppService),
	}
}

// Init 初始化模块（默认实现，用于兼容性）
func (m *ProjectModule) Init(db *gorm.DB) any {
	return m.InitWithDeps(db)
}

// RegisterRoutes 注册路由（由 ModuleManager 调用）
func (m *ProjectModule) RegisterRoutes(handlers any) []shared_http.RouteGroupConfig {
	h, ok := handlers.(*Handlers)
	if !ok {
		return nil
	}
	return h.RegisterRoutes()
}
