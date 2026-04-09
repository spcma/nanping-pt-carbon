package application

import (
	"app/internal/module/iam/domain"
	"app/internal/shared/entity"
	"context"
)

// ===== Repository Ports =====

type UserRepo interface {
	Create(ctx context.Context, user *domain.Users) error
	Update(ctx context.Context, user *domain.Users) error
	UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error // 添加部分更新方法
	Delete(ctx context.Context, id, uid int64) error
	FindByID(ctx context.Context, id int64) (*domain.Users, error)
	FindByUsername(ctx context.Context, username string) (*domain.Users, error) // 添加按用户名查找
	FindByQuery(ctx context.Context, query *domain.UserQuery) (*domain.Users, error)
	FindList(ctx context.Context) ([]*domain.Users, error)
	FindPage(ctx context.Context, query *domain.UsersPageQuery) (*entity.PaginationResult[domain.Users], error)
}

type RoleRepo interface {
	Create(ctx context.Context, role *domain.SysRole) error
	Update(ctx context.Context, role *domain.SysRole) error
	UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error // 添加部分更新方法
	Delete(ctx context.Context, id, uid int64) error
	FindByID(ctx context.Context, id int64) (*domain.SysRole, error)
	FindList(ctx context.Context) ([]*domain.SysRole, error)
	FindPage(ctx context.Context, query *domain.SysRolePageQuery) (*entity.PaginationResult[domain.SysRole], error)
	FindByCode(ctx context.Context, code string) (*domain.SysRole, error)
	FindListByCodes(ctx context.Context, codes []string) ([]*domain.SysRole, error)
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

// SysApiRepository API repository interface
type SysApiRepository interface {
	Create(ctx context.Context, api *domain.SysApi) error
	Update(ctx context.Context, api *domain.SysApi) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*domain.SysApi, error)
	FindByCode(ctx context.Context, code string) (*domain.SysApi, error)
	FindPage(ctx context.Context, pageNum, pageSize int64, name string) ([]*domain.SysApi, int64, error)
	FindAll(ctx context.Context) ([]*domain.SysApi, error)
}

// SysUserRoleRepository user role repository interface
type SysUserRoleRepository interface {
	Create(ctx context.Context, userRole *domain.SysUserRole) error
	Delete(ctx context.Context, userId int64) error
	DeleteByUser(ctx context.Context, userId int64) error
	DeleteByRole(ctx context.Context, roleId int64) error
	DeleteByUserRoleID(ctx context.Context, userId, roleId int64) error
	FindByUserID(ctx context.Context, userId int64) ([]*domain.SysUserRole, error)
	FindByRoleID(ctx context.Context, roleId int64) ([]*domain.SysUserRole, error)
	FindByUserAndRole(ctx context.Context, userId, roleId int64) (*domain.SysUserRole, error)
}

type RolePermissionRepo interface {
	HasPermission(ctx context.Context, roleCode, permissionCode string) (bool, error)
}
