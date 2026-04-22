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

func (r *UserRoleRepository) Create(ctx context.Context, userRole *domain.UserRole) error {
	return r.db.WithContext(ctx).Create(userRole).Error
}

func (r *UserRoleRepository) FindByUserID(ctx context.Context, userID int64) ([]*domain.UserRole, error) {
	var userRoles []*domain.UserRole
	err := r.db.WithContext(ctx).Where("sys_user_id = ?", userID).Find(&userRoles).Error
	return userRoles, err
}

func (r *UserRoleRepository) FindByRoleID(ctx context.Context, roleID int64) ([]*domain.UserRole, error) {
	var userRoles []*domain.UserRole
	err := r.db.WithContext(ctx).Where("sys_role_id = ?", roleID).Find(&userRoles).Error
	return userRoles, err
}

func (r *UserRoleRepository) DeleteByUser(ctx context.Context, userId int64) error {
	return r.db.WithContext(ctx).Where("sys_user_id = ?", userId).Delete(&domain.UserRole{}).Error
}

func (r *UserRoleRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Where("sys_user_id = ?", userID).Delete(&domain.UserRole{}).Error
}

func (r *UserRoleRepository) Delete(ctx context.Context, userId int64) error {
	return r.db.WithContext(ctx).Where("sys_user_id = ?", userId).Delete(&domain.UserRole{}).Error
}

func (r *UserRoleRepository) DeleteByUserRoleID(ctx context.Context, userID, roleID int64) error {
	return r.db.WithContext(ctx).Where("sys_user_id = ? AND sys_role_id = ?", userID, roleID).Delete(&domain.UserRole{}).Error
}

func (r *UserRoleRepository) FindRoleCodesByUserID(ctx context.Context, userID int64) ([]string, error) {
	var roleCodes []string
	err := r.db.WithContext(ctx).
		Table("sys_user_role ur").
		Joins("JOIN sys_role r ON ur.sys_role_id = r.id").
		Where("ur.sys_user_id = ? AND r.status = ?", userID, domain.RoleStatusNormal).
		Pluck("r.code", &roleCodes).Error
	return roleCodes, err
}

func (r *UserRoleRepository) DeleteByRole(ctx context.Context, roleId int64) error {
	return r.db.WithContext(ctx).Where("sys_role_id = ?", roleId).Delete(&domain.UserRole{}).Error
}

func (r *UserRoleRepository) FindByUserAndRole(ctx context.Context, userId, roleId int64) (*domain.UserRole, error) {
	var userRole domain.UserRole
	err := r.db.WithContext(ctx).
		Where("sys_user_id = ? AND sys_role_id = ?", userId, roleId).
		First(&userRole).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &userRole, nil
}
