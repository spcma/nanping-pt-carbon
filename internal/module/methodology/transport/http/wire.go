package http

import (
	"app/internal/module/methodology/application"
	"app/internal/module/methodology/domain"
	"app/internal/module/methodology/infrastructure"

	"gorm.io/gorm"
)

type MethodologyDDD struct {
	Repo       domain.MethodologyRepository
	AppService *application.MethodologyAppService
}

// InitMethodologyDDD initializes methodology DDD components
func InitMethodologyDDD(db *gorm.DB) *MethodologyDDD {
	repo := infrastructure.NewMethodologyRepository(db)
	appService := application.NewMethodologyAppService(repo)
	return &MethodologyDDD{
		Repo:       repo,
		AppService: appService,
	}
}
