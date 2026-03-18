package http

import (
	"app/internal/module/methodology/application"
	"app/internal/module/methodology/infrastructure"

	"gorm.io/gorm"
)

type MethodologyDDD struct {
	Repo       application.MethodologyRepo
	AppService *application.MethodologyAppService
}

// InitMethodologyWire initializes methodology DDD components
func InitMethodologyWire(db *gorm.DB) *MethodologyDDD {
	repo := infrastructure.NewMethodologyRepository(db)
	appService := application.NewMethodologyAppService(repo)
	return &MethodologyDDD{
		Repo:       repo,
		AppService: appService,
	}
}
