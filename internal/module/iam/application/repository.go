package application

import (
	"app/internal/module/iam/domain"
	"app/internal/shared/entity"
	"context"
)

// ===== Repository Ports =====

type UserRepo interface {
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error // 添加部分更新方法
	Delete(ctx context.Context, id, uid int64) error
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error) // 添加按用户名查找
	FindByQuery(ctx context.Context, query *domain.UserQuery) (*domain.User, error)
	FindList(ctx context.Context) ([]*domain.User, error)
	FindPage(ctx context.Context, query *domain.UsersPageQuery) (*entity.PaginationResult[*domain.User], error)
}

type RoleRepo interface {
	Create(ctx context.Context, role *domain.Role) error
	Update(ctx context.Context, role *domain.Role) error
	UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error // 添加部分更新方法
	Delete(ctx context.Context, id, uid int64) error
	FindByID(ctx context.Context, id int64) (*domain.Role, error)
	FindList(ctx context.Context) ([]*domain.Role, error)
	FindPage(ctx context.Context, query *domain.SysRolePageQuery) (*entity.PaginationResult[*domain.Role], error)
	FindByCode(ctx context.Context, code string) (*domain.Role, error)
	FindListByCodes(ctx context.Context, codes []string) ([]*domain.Role, error)
}

type UserRoleRepo interface {
	Create(ctx context.Context, userRole *domain.UserRole) error
	Delete(ctx context.Context, userId int64) error
	FindByUserID(ctx context.Context, userID int64) ([]*domain.UserRole, error)
	FindByRoleID(ctx context.Context, roleID int64) ([]*domain.UserRole, error)
	DeleteByUser(ctx context.Context, userId int64) error
	DeleteByUserID(ctx context.Context, userID int64) error
	DeleteByUserRoleID(ctx context.Context, userID, roleID int64) error
	DeleteByRole(ctx context.Context, roleId int64) error
	FindByUserAndRole(ctx context.Context, userId, roleId int64) (*domain.UserRole, error)
	FindRoleCodesByUserID(ctx context.Context, userID int64) ([]string, error)
}

// ApiRepository API repository interface
type ApiRepository interface {
	Create(ctx context.Context, api *domain.Api) error
	Update(ctx context.Context, api *domain.Api) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*domain.Api, error)
	FindByCode(ctx context.Context, code string) (*domain.Api, error)
	FindPage(ctx context.Context, pageNum, pageSize int64, name string) ([]*domain.Api, int64, error)
	FindAll(ctx context.Context) ([]*domain.Api, error)
}

// UserRoleRepository user role repository interface
type UserRoleRepository interface {
	Create(ctx context.Context, userRole *domain.UserRole) error
	Delete(ctx context.Context, userId int64) error
	DeleteByUser(ctx context.Context, userId int64) error
	DeleteByRole(ctx context.Context, roleId int64) error
	DeleteByUserRoleID(ctx context.Context, userId, roleId int64) error
	FindByUserID(ctx context.Context, userId int64) ([]*domain.UserRole, error)
	FindByRoleID(ctx context.Context, roleId int64) ([]*domain.UserRole, error)
	FindByUserAndRole(ctx context.Context, userId, roleId int64) (*domain.UserRole, error)
}

type RolePermissionRepo interface {
	HasPermission(ctx context.Context, roleCode, permissionCode string) (bool, error)
}
