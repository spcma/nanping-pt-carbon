package http

import (
	"app/internal/module/carbonreportday/application"
	"app/internal/module/carbonreportday/domain"
	"app/internal/module/carbonreportday/infrastructure"

	"gorm.io/gorm"
)

type CarbonReportDayDDD struct {
	Repo       domain.CarbonReportDayRepository
	AppService *application.CarbonReportDayAppService
}

// InitCarbonReportDayDDD initializes carbon report day DDD components
func InitCarbonReportDayDDD(db *gorm.DB) *CarbonReportDayDDD {
	repo := infrastructure.NewCarbonReportDayRepository(db)
	appService := application.NewCarbonReportDayAppService(repo)
	return &CarbonReportDayDDD{
		Repo:       repo,
		AppService: appService,
	}
}
