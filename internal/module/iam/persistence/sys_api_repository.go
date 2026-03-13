package persistence

import (
	"app/internal/module/iam/domain"
	"app/internal/shared/db"
	"context"
	"gorm.io/gorm"
)

type sysApiRepository struct {
	db *gorm.DB
}

// NewSysApiRepository creates API repository
func NewSysApiRepository(db *gorm.DB) domain.SysApiRepository {
	return &sysApiRepository{db: db}
}

func (r *sysApiRepository) Create(ctx context.Context, api *domain.SysApi) error {
	return db.WithContext(ctx, r.db).Create(api).Error
}

func (r *sysApiRepository) Update(ctx context.Context, api *domain.SysApi) error {
	return db.WithContext(ctx, r.db).Save(api).Error
}

func (r *sysApiRepository) Delete(ctx context.Context, id int64) error {
	return db.WithContext(ctx, r.db).Delete(&domain.SysApi{}, id).Error
}

func (r *sysApiRepository) FindByID(ctx context.Context, id int64) (*domain.SysApi, error) {
	var api domain.SysApi
	err := db.WithContext(ctx, r.db).First(&api, id).Error
	return &api, err
}

func (r *sysApiRepository) FindByCode(ctx context.Context, code string) (*domain.SysApi, error) {
	var api domain.SysApi
	err := db.WithContext(ctx, r.db).Where("code = ?", code).First(&api).Error
	return &api, err
}

func (r *sysApiRepository) FindPage(ctx context.Context, pageNum, pageSize int64, name string) ([]*domain.SysApi, int64, error) {
	var apis []*domain.SysApi
	var total int64

	tx := db.WithContext(ctx, r.db).Model(&domain.SysApi{})
	if name != "" {
		tx = tx.Where("name LIKE ?", "%"+name+"%")
	}

	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pageNum - 1) * pageSize
	err := tx.Offset(int(offset)).Limit(int(pageSize)).Find(&apis).Error
	return apis, total, err
}

func (r *sysApiRepository) FindAll(ctx context.Context) ([]*domain.SysApi, error) {
	var apis []*domain.SysApi
	err := db.WithContext(ctx, r.db).Find(&apis).Error
	return apis, err
}
