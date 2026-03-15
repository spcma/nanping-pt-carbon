package infrastructure

import (
	"app/internal/module/iam/application"
	"app/internal/module/iam/domain"
	"context"
	"errors"
	"gorm.io/gorm"
)

type UserRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) application.UserRoleRepo {
	return &UserRoleRepository{db: db}
}

func (r *UserRoleRepository) Create(ctx context.Context, userRole *domain.SysUserRole) error {
	return r.db.WithContext(ctx).Create(userRole).Error
}

func (r *UserRoleRepository) FindByUserID(ctx context.Context, userID int64) ([]*domain.SysUserRole, error) {
	var userRoles []*domain.SysUserRole
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&userRoles).Error
	return userRoles, err
}

func (r *UserRoleRepository) FindByRoleID(ctx context.Context, roleID int64) ([]*domain.SysUserRole, error) {
	var userRoles []*domain.SysUserRole
	err := r.db.WithContext(ctx).Where("role_id = ?", roleID).Find(&userRoles).Error
	return userRoles, err
}

func (r *UserRoleRepository) DeleteByUser(ctx context.Context, userId int64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userId).Delete(&domain.SysUserRole{}).Error
}

func (r *UserRoleRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.SysUserRole{}).Error
}

func (r *UserRoleRepository) Delete(ctx context.Context, userId int64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userId).Delete(&domain.SysUserRole{}).Error
}

func (r *UserRoleRepository) DeleteByUserRoleID(ctx context.Context, userID, roleID int64) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&domain.SysUserRole{}).Error
}

func (r *UserRoleRepository) FindRoleCodesByUserID(ctx context.Context, userID int64) ([]string, error) {
	var roleCodes []string
	err := r.db.WithContext(ctx).
		Table("user_roles ur").
		Joins("JOIN roles r ON ur.role_id = r.id").
		Where("ur.user_id = ? AND r.status = ?", userID, domain.RoleStatusNormal).
		Pluck("r.code", &roleCodes).Error
	return roleCodes, err
}

func (r *UserRoleRepository) DeleteByRole(ctx context.Context, roleId int64) error {
	return r.db.WithContext(ctx).Where("role_id = ?", roleId).Delete(&domain.SysUserRole{}).Error
}

func (r *UserRoleRepository) FindByUserAndRole(ctx context.Context, userId, roleId int64) (*domain.SysUserRole, error) {
	var userRole domain.SysUserRole
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userId, roleId).
		First(&userRole).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &userRole, nil
}
