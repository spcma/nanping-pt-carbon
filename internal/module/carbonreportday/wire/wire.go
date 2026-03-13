package wire

import (
	"app/internal/module/carbonreportday/application"
	"app/internal/module/carbonreportday/domain"
	"app/internal/module/carbonreportday/persistence"

	"gorm.io/gorm"
)

type CarbonReportDayDDD struct {
	Repo       domain.CarbonReportDayRepository
	AppService *application.CarbonReportDayAppService
}

// InitCarbonReportDayDDD initializes carbon report day DDD components
func InitCarbonReportDayDDD(db *gorm.DB) *CarbonReportDayDDD {
	repo := persistence.NewCarbonReportDayRepository(db)
	appService := application.NewCarbonReportDayAppService(repo)
	return &CarbonReportDayDDD{
		Repo:       repo,
		AppService: appService,
	}
}
