package http

import (
	"app/internal/module/project/application"
	"app/internal/module/project/infrastructure"

	"gorm.io/gorm"
)

type ProjectWire struct {
	Repo       application.ProjectRepository
	AppService *application.ProjectAppService
}

// InitProjectWire initializes project DDD components
func InitProjectWire(db *gorm.DB) *ProjectWire {
	projectRepo := infrastructure.NewProjectRepository(db)
	appService := application.NewProjectAppService(projectRepo)
	return &ProjectWire{
		Repo:       projectRepo,
		AppService: appService,
	}
}
