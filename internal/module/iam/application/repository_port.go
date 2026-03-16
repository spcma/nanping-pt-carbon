package application

import (
	"app/internal/module/iam/domain"
	"app/internal/shared/entity"
	"context"
)

// ===== Repository Ports =====

type UserRepo interface {
	Create(ctx context.Context, user *domain.SysUser) error
	Update(ctx context.Context, user *domain.SysUser) error
	UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error // 添加部分更新方法
	Delete(ctx context.Context, id, uid int64) error
	FindByID(ctx context.Context, id int64) (*domain.SysUser, error)
	FindByUsername(ctx context.Context, username string) (*domain.SysUser, error) // 添加按用户名查找
	FindList(ctx context.Context) ([]*domain.SysUser, error)
	FindPage(ctx context.Context, pageNum, pageSize int) (*entity.PaginationResult[*domain.SysUser], error)
}

type RoleRepo interface {
	Create(ctx context.Context, role *domain.SysRole) error
	Update(ctx context.Context, role *domain.SysRole) error
	UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error // 添加部分更新方法
	Delete(ctx context.Context, id, uid int64) error
	FindByID(ctx context.Context, id int64) (*domain.SysRole, error)
	FindList(ctx context.Context) ([]*domain.SysRole, error)
	FindPage(ctx context.Context, pageNum, pageSize int) (*entity.PaginationResult[*domain.SysRole], error)
	FindByCode(ctx context.Context, code string) (*domain.SysRole, error)
}

type UserRoleRepo interface {
	Create(ctx context.Context, userRole *domain.SysUserRole) error
	Delete(ctx context.Context, userId int64) error
	FindByUserID(ctx context.Context, userID int64) ([]*domain.SysUserRole, error)
	FindByRoleID(ctx context.Context, roleID int64) ([]*domain.SysUserRole, error)
	DeleteByUser(ctx context.Context, userId int64) error
	DeleteByUserID(ctx context.Context, userID int64) error
	DeleteByUserRoleID(ctx context.Context, userID, roleID int64) error
	DeleteByRole(ctx context.Context, roleId int64) error
	FindByUserAndRole(ctx context.Context, userId, roleId int64) (*domain.SysUserRole, error)
	FindRoleCodesByUserID(ctx context.Context, userID int64) ([]string, error)
}

type RolePermissionRepo interface {
	HasPermission(ctx context.Context, roleCode, permissionCode string) (bool, error)
}
