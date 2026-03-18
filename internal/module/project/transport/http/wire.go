package http

import (
	"app/internal/module/project/application"
	"app/internal/module/project/infrastructure"

	"gorm.io/gorm"
)

type ProjectWire struct {
	Repo    application.ProjectRepo
	Service *application.ProjectAppService
}

// InitProjectWire initializes project DDD components
func InitProjectWire(db *gorm.DB) *ProjectWire {
	projectRepo := infrastructure.NewProjectRepository(db)
	appService := application.NewProjectService(projectRepo)
	return &ProjectWire{
		Repo:    projectRepo,
		Service: appService,
	}
}

type ProjectMembersWire struct {
	Repo    application.ProjectMembersRepo
	Service *application.ProjectMembersService
}

// InitProjectMembersWire initializes project members DDD components
func InitProjectMembersWire(db *gorm.DB) *ProjectMembersWire {
	projectMembersRepo := infrastructure.NewProjectMembersRepository(db)
	appService := application.NewProjectMembersService(projectMembersRepo)
	return &ProjectMembersWire{
		Repo:    projectMembersRepo,
		Service: appService,
	}
}
