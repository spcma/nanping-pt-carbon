package wire

import (
	"app/internal/module/iam/application"
	"app/internal/module/iam/domain"
	"app/internal/module/iam/infrastructure"

	"gorm.io/gorm"
)

type SysUserDDD struct {
	Repo       domain.SysUserRepository
	AppService *application.SysUserAppService
}

// InitSysUserDDD initializes system user DDD components
func InitSysUserDDD(db *gorm.DB) *SysUserDDD {
	repo := infrastructure.NewUserRepository(db)
	appService := application.NewSysUserAppService(repo)
	return &SysUserDDD{
		Repo:       repo,
		AppService: appService,
	}
}

type SysRoleDDD struct {
	Repo       domain.SysRoleRepository
	AppService *application.SysRoleAppService
}

// InitSysRoleDDD initializes system role DDD components
func InitSysRoleDDD(db *gorm.DB) *SysRoleDDD {
	repo := infrastructure.NewRoleRepository(db)
	appService := application.NewSysRoleAppService(repo)
	return &SysRoleDDD{
		Repo:       repo,
		AppService: appService,
	}
}

type SysApiDDD struct {
	Repo       domain.SysApiRepository
	AppService *application.SysApiAppService
}

// InitSysApiDDD initializes system API DDD components
func InitSysApiDDD(db *gorm.DB) *SysApiDDD {
	repo := infrastructure.NewApiRepository(db)
	appService := application.NewSysApiAppService(repo)
	return &SysApiDDD{
		Repo:       repo,
		AppService: appService,
	}
}

type SysUserRoleDDD struct {
	Repo       domain.SysUserRoleRepository
	AppService *application.SysUserRoleAppService
}

// InitSysUserRoleDDD initializes system user role DDD components
func InitSysUserRoleDDD(db *gorm.DB) *SysUserRoleDDD {
	repo := infrastructure.NewUserRoleRepository(db)
	appService := application.NewSysUserRoleAppService(repo)
	return &SysUserRoleDDD{
		Repo:       repo,
		AppService: appService,
	}
}
