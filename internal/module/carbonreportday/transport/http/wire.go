package http

import (
	"app/internal/module/carbonreportday/application"
	"app/internal/module/carbonreportday/infrastructure"

	"gorm.io/gorm"
)

type CarbonReportDayWire struct {
	Repo    application.CarbonReportDayRepo
	Service *application.CarbonReportDayAppService
}

// InitCarbonReportDayWire initializes carbon report day DDD components
func InitCarbonReportDayWire(db *gorm.DB) *CarbonReportDayWire {
	repo := infrastructure.NewCarbonReportDayRepository(db)
	service := application.NewCarbonReportDayAppService(repo)

	return &CarbonReportDayWire{
		Repo:    repo,
		Service: service,
	}
}
