package infrastructure

import (
	"app/internal/module/iam/domain"
	db "app/internal/shared/db"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"

	"gorm.io/gorm"
)

type RoleRepository struct{}

func NewRoleRepository(_db *gorm.DB) *RoleRepository {
	return &RoleRepository{}
}

func (r *RoleRepository) Create(ctx context.Context, role *domain.Role) error {
	return db.GetDB(ctx).WithContext(ctx).Create(role).Error
}

func (r *RoleRepository) Update(ctx context.Context, role *domain.Role) error {
	return db.GetDB(ctx).WithContext(ctx).Updates(role).Error
}

func (r *RoleRepository) UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error {
	return db.GetDB(ctx).WithContext(ctx).Model(&domain.Role{}).Where("id = ?", id).Updates(updates).Error
}

func (r *RoleRepository) Delete(ctx context.Context, id, uid int64) error {
	// 实现逻辑删除，更新状态而不是物理删除
	updates := map[string]interface{}{
		entity.FieldDeleteBy:   0,
		entity.FieldDeleteTime: timeutil.Now(),
	}
	return db.GetDB(ctx).WithContext(ctx).Model(&domain.Role{}).Where("id = ?", id).Updates(updates).Error
}

func (r *RoleRepository) FindByID(ctx context.Context, id int64) (*domain.Role, error) {
	var role domain.Role
	err := db.GetDB(ctx).WithContext(ctx).
		Table("sys_role").
		Where("id = ?", id).
		Where(entity.FieldDeleteBy + " = 0").First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindList(ctx context.Context) ([]*domain.Role, error) {
	var roles []*domain.Role
	err := db.GetDB(ctx).WithContext(ctx).
		Table("sys_role").
		Where(entity.FieldDeleteBy + " = 0").
		Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) FindPage(ctx context.Context, query *domain.SysRolePageQuery) (*entity.PaginationResult[*domain.Role], error) {

	helper := db.NewPaginationHelper[*domain.Role](db.GetDB(ctx))
	result, err := helper.PageQuery(query.PageNum, query.PageSize, func(dq *gorm.DB) *gorm.DB {
		dq = dq.WithContext(ctx).
			Table("sys_role").
			Where(entity.FieldDeleteBy + " = 0")

		return dq
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *RoleRepository) FindByCode(ctx context.Context, code string) (*domain.Role, error) {
	var role domain.Role
	err := db.GetDB(ctx).WithContext(ctx).
		Table("sys_role").
		Where("code = ?", code).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindListByCodes(ctx context.Context, codes []string) ([]*domain.Role, error) {
	var roles []*domain.Role
	err := db.GetDB(ctx).WithContext(ctx).
		Table("sys_role").
		Where("code IN ?", codes).
		Where(entity.FieldDeleteBy + " = 0").
		Find(&roles).Error
	return roles, err
}
