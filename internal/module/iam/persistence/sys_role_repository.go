package persistence

import (
	"app/internal/module/iam/domain"
	"app/internal/shared/db"
	"context"
	"gorm.io/gorm"
	"strings"
)

type sysRoleRepository struct {
	db *gorm.DB
}

// NewSysRoleRepository creates role repository
func NewSysRoleRepository(db *gorm.DB) domain.SysRoleRepository {
	return &sysRoleRepository{db: db}
}

func (r *sysRoleRepository) Create(ctx context.Context, role *domain.SysRole) error {
	return db.WithContext(ctx, r.db).Create(role).Error
}

func (r *sysRoleRepository) Update(ctx context.Context, role *domain.SysRole) error {
	return db.WithContext(ctx, r.db).Save(role).Error
}

func (r *sysRoleRepository) Delete(ctx context.Context, id int64) error {
	return db.WithContext(ctx, r.db).Delete(&domain.SysRole{}, id).Error
}

func (r *sysRoleRepository) FindByID(ctx context.Context, id int64) (*domain.SysRole, error) {
	var role domain.SysRole
	err := db.WithContext(ctx, r.db).First(&role, id).Error
	return &role, err
}

func (r *sysRoleRepository) FindByCode(ctx context.Context, code string) (*domain.SysRole, error) {
	var role domain.SysRole
	err := db.WithContext(ctx, r.db).Where("code = ?", code).First(&role).Error
	return &role, err
}

func (r *sysRoleRepository) FindPage(ctx context.Context, query domain.SysRolePageQuery) ([]*domain.SysRole, int64, error) {
	var roles []*domain.SysRole
	var total int64

	tx := db.WithContext(ctx, r.db).Model(&domain.SysRole{})

	// 动态构建查询条件
	if query.Name != "" {
		tx = tx.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Code != "" {
		tx = tx.Where("code = ?", query.Code)
	}
	if query.Status != "" {
		tx = tx.Where("status = ?", query.Status)
	}

	// 计数
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	if query.SortBy != "" {
		order := query.SortBy + " " + strings.ToUpper(query.SortOrder)
		tx = tx.Order(order)
	} else {
		tx = tx.Order("id DESC") // 默认按 ID 降序
	}

	// 分页
	offset := (query.PageNum - 1) * query.PageSize
	err := tx.Offset(int(offset)).Limit(int(query.PageSize)).Find(&roles).Error
	return roles, total, err
}

func (r *sysRoleRepository) FindListByCodes(ctx context.Context, codes []string) ([]*domain.SysRole, error) {
	var roles []*domain.SysRole
	err := db.WithContext(ctx, r.db).Where("code IN ?", codes).Find(&roles).Error
	return roles, err
}
