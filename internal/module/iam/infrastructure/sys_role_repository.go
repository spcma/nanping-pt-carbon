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

type RoleRepository struct {
	*db.BaseRepository[domain.SysRole]
}

func NewRoleRepository(_db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		BaseRepository: db.NewBaseRepository[domain.SysRole](_db),
	}
}

func (r *RoleRepository) Create(ctx context.Context, role *domain.SysRole) error {
	return r.GetDB(ctx).WithContext(ctx).Create(role).Error
}

func (r *RoleRepository) Update(ctx context.Context, role *domain.SysRole) error {
	return r.GetDB(ctx).WithContext(ctx).Updates(role).Error
}

func (r *RoleRepository) UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.GetDB(ctx).WithContext(ctx).Model(&domain.SysRole{}).Where("id = ?", id).Updates(updates).Error
}

func (r *RoleRepository) Delete(ctx context.Context, id, uid int64) error {
	// 实现逻辑删除，更新状态而不是物理删除
	updates := map[string]interface{}{
		entity.FieldDeleteBy:   0,
		entity.FieldDeleteTime: timeutil.Now(),
	}
	return r.GetDB(ctx).WithContext(ctx).Model(&domain.SysRole{}).Where("id = ?", id).Updates(updates).Error
}

func (r *RoleRepository) FindByID(ctx context.Context, id int64) (*domain.SysRole, error) {
	var role domain.SysRole
	err := r.GetDB(ctx).WithContext(ctx).
		Table("SysRole").
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

func (r *RoleRepository) FindList(ctx context.Context) ([]*domain.SysRole, error) {
	var SysRole []*domain.SysRole
	err := r.GetDB(ctx).WithContext(ctx).
		Table("SysRole").
		Where(entity.FieldDeleteBy + " = 0").
		Find(&SysRole).Error
	return SysRole, err
}

func (r *RoleRepository) FindPage(ctx context.Context, query *domain.SysRolePageQuery) (*entity.PaginationResult[domain.SysRole], error) {
	txDB := r.GetDB(ctx)
	// 使用通用分页助手
	helper := db.NewPaginationHelper[domain.SysRole](txDB)
	result, err := helper.PageQuery(int(query.PageNum), int(query.PageSize), func(dq *gorm.DB) *gorm.DB {
		// 构建基础查询
		dq = txDB.WithContext(ctx).
			Table("SysRole").
			Where(entity.FieldDeleteBy + " = 0")

		return dq
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *RoleRepository) FindByCode(ctx context.Context, code string) (*domain.SysRole, error) {
	var role domain.SysRole
	err := r.GetDB(ctx).WithContext(ctx).
		Table("SysRole").
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

func (r *RoleRepository) FindListByCodes(ctx context.Context, codes []string) ([]*domain.SysRole, error) {
	var roles []*domain.SysRole
	err := r.GetDB(ctx).WithContext(ctx).
		Table("SysRole").
		Where("code IN ?", codes).
		Where(entity.FieldDeleteBy + " = 0").
		Find(&roles).Error
	return roles, err
}
