package wire

import (
	"app/internal/module/project/application"
	"app/internal/module/project/domain"
	"app/internal/module/project/infrastructure"

	"gorm.io/gorm"
)

type ProjectDDD struct {
	Repo       domain.ProjectRepository
	AppService *application.ProjectAppService
}

// InitProjectDDD initializes project DDD components
func InitProjectDDD(db *gorm.DB) *ProjectDDD {
	repo := infrastructure.NewProjectRepository(db)
	appService := application.NewProjectAppService(repo)
	return &ProjectDDD{
		Repo:       repo,
		AppService: appService,
	}
}
