package domain

import (
	"context"
)

// SysUserRoleRepository user role repository interface
type SysUserRoleRepository interface {
	Create(ctx context.Context, userRole *SysUserRole) error
	Delete(ctx context.Context, userId int64) error
	DeleteByUser(ctx context.Context, userId int64) error
	DeleteByRole(ctx context.Context, roleId int64) error
	DeleteByUserRoleID(ctx context.Context, userId, roleId int64) error
	FindByUserID(ctx context.Context, userId int64) ([]*SysUserRole, error)
	FindByRoleID(ctx context.Context, roleId int64) ([]*SysUserRole, error)
	FindByUserAndRole(ctx context.Context, userId, roleId int64) (*SysUserRole, error)
}
