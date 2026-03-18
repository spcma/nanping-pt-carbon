package http

import (
	"app/internal/module/iam/application"
	"app/internal/module/iam/infrastructure"

	"gorm.io/gorm"
)

type SysUserWire struct {
	Repo    application.UserRepo
	Service *application.UsersService
}

// InitSysUserWire initializes system user DDD components
func InitSysUserWire(db *gorm.DB) *SysUserWire {
	repo := infrastructure.NewUserRepository(db)
	appService := application.NewUsersService(repo)
	return &SysUserWire{
		Repo:    repo,
		Service: appService,
	}
}

type SysRoleWire struct {
	Repo    application.RoleRepo
	Service *application.SysRoleAppService
}

// InitSysRoleWire initializes system role DDD components
func InitSysRoleWire(db *gorm.DB) *SysRoleWire {
	repo := infrastructure.NewRoleRepository(db)
	appService := application.NewSysRoleAppService(repo)
	return &SysRoleWire{
		Repo:    repo,
		Service: appService,
	}
}

type SysApiWire struct {
	Repo       application.SysApiRepository
	AppService *application.SysApiAppService
}

// InitSysApiWire initializes system API DDD components
func InitSysApiWire(db *gorm.DB) *SysApiWire {
	repo := infrastructure.NewApiRepository(db)
	appService := application.NewSysApiAppService(repo)
	return &SysApiWire{
		Repo:       repo,
		AppService: appService,
	}
}

type SysUserRoleWire struct {
	Repo       application.SysUserRoleRepository
	AppService *application.SysUserRoleAppService
}

// InitSysUserRoleWire initializes system user role DDD components
func InitSysUserRoleWire(db *gorm.DB) *SysUserRoleWire {
	repo := infrastructure.NewUserRoleRepository(db)
	appService := application.NewSysUserRoleAppService(repo)
	return &SysUserRoleWire{
		Repo:       repo,
		AppService: appService,
	}
}
