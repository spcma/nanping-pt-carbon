package persistence

import (
	"app/internal/module/iam/domain"
	"app/internal/shared/db"
	"context"
	"gorm.io/gorm"
)

type sysUserRoleRepository struct {
	db *gorm.DB
}

// NewSysUserRoleRepository creates user role repository
func NewSysUserRoleRepository(db *gorm.DB) domain.SysUserRoleRepository {
	return &sysUserRoleRepository{db: db}
}

func (r *sysUserRoleRepository) Create(ctx context.Context, userRole *domain.SysUserRole) error {
	return db.WithContext(ctx, r.db).Create(userRole).Error
}

func (r *sysUserRoleRepository) Delete(ctx context.Context, userId, roleId int64) error {
	return db.WithContext(ctx, r.db).Where("sys_user_id = ? AND sys_role_id = ?", userId, roleId).Delete(&domain.SysUserRole{}).Error
}

func (r *sysUserRoleRepository) DeleteByUser(ctx context.Context, userId int64) error {
	return db.WithContext(ctx, r.db).Where("sys_user_id = ?", userId).Delete(&domain.SysUserRole{}).Error
}

func (r *sysUserRoleRepository) DeleteByRole(ctx context.Context, roleId int64) error {
	return db.WithContext(ctx, r.db).Where("sys_role_id = ?", roleId).Delete(&domain.SysUserRole{}).Error
}

func (r *sysUserRoleRepository) FindByUserID(ctx context.Context, userId int64) ([]*domain.SysUserRole, error) {
	var userRoles []*domain.SysUserRole
	err := db.WithContext(ctx, r.db).Where("sys_user_id = ?", userId).Find(&userRoles).Error
	return userRoles, err
}

func (r *sysUserRoleRepository) FindByRoleID(ctx context.Context, roleId int64) ([]*domain.SysUserRole, error) {
	var userRoles []*domain.SysUserRole
	err := db.WithContext(ctx, r.db).Where("sys_role_id = ?", roleId).Find(&userRoles).Error
	return userRoles, err
}

func (r *sysUserRoleRepository) FindByUserAndRole(ctx context.Context, userId, roleId int64) (*domain.SysUserRole, error) {
	var userRole domain.SysUserRole
	err := db.WithContext(ctx, r.db).Where("sys_user_id = ? AND sys_role_id = ?", userId, roleId).First(&userRole).Error
	return &userRole, err
}
