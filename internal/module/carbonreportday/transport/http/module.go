package http

import (
	"app/internal/module/carbonreportday/wire"
	shared_http "app/internal/shared/http"

	"gorm.io/gorm"
)

// CarbonReportDayModule 碳报告日报模块实现
type CarbonReportDayModule struct {
	CarbonReportDayDDD *wire.CarbonReportDayDDD
}

// Name 返回模块名称
func (m *CarbonReportDayModule) Name() string {
	return "CarbonReportDay"
}

// InitWithDeps 初始化模块（支持外部依赖注入）
func (m *CarbonReportDayModule) InitWithDeps(db *gorm.DB, deps ...any) any {
	m.CarbonReportDayDDD = wire.InitCarbonReportDayDDD(db)

	// 创建并返回 Handlers 实例
	return &Handlers{
		CarbonReportDayHandler: NewCarbonReportDayHandler(m.CarbonReportDayDDD.AppService),
	}
}

// Init 初始化模块（默认实现，用于兼容性）
func (m *CarbonReportDayModule) Init(db *gorm.DB) any {
	return m.InitWithDeps(db)
}

// RegisterRoutes 注册路由（由 ModuleManager 调用）
func (m *CarbonReportDayModule) RegisterRoutes(handlers any) []shared_http.RouteGroupConfig {
	h, ok := handlers.(*Handlers)
	if !ok {
		return nil
	}
	return h.RegisterRoutes()
}
