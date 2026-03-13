package persistence

import (
	"app/internal/module/methodology/domain"
	"app/internal/shared/db"
	"app/internal/shared/timeutil"
	"context"
	"strings"

	"gorm.io/gorm"
)

type methodologyRepository struct {
	db *gorm.DB
}

// NewMethodologyRepository creates methodology repository
func NewMethodologyRepository(db *gorm.DB) domain.MethodologyRepository {
	return &methodologyRepository{db: db}
}

func (r *methodologyRepository) Create(ctx context.Context, methodology *domain.Methodology) error {
	return db.WithContext(ctx, r.db).Create(methodology).Error
}

func (r *methodologyRepository) Update(ctx context.Context, methodology *domain.Methodology) error {
	return db.WithContext(ctx, r.db).Save(methodology).Error
}

func (r *methodologyRepository) Delete(ctx context.Context, id int64, userID int64) error {
	// 逻辑删除：只更新 delete_user 和 delete_time 字段
	updates := map[string]interface{}{
		"delete_by":   userID,
		"delete_time": timeutil.New(),
	}
	return db.WithContext(ctx, r.db).Model(&domain.Methodology{}).Where("id = ?", id).Updates(updates).Error
}

func (r *methodologyRepository) FindByID(ctx context.Context, id int64) (*domain.Methodology, error) {
	var methodology domain.Methodology
	err := db.WithContext(ctx, r.db).Where("id = ? AND delete_time IS NULL", id).First(&methodology).Error
	return &methodology, err
}

func (r *methodologyRepository) FindByCode(ctx context.Context, code string) (*domain.Methodology, error) {
	var methodology domain.Methodology
	err := db.WithContext(ctx, r.db).Where("code = ? AND delete_time IS NULL", code).First(&methodology).Error
	return &methodology, err
}

func (r *methodologyRepository) FindPage(ctx context.Context, query domain.MethodologyPageQuery) ([]*domain.Methodology, int64, error) {
	var methodologies []*domain.Methodology
	var total int64

	tx := db.WithContext(ctx, r.db).Model(&domain.Methodology{}).Where("delete_time IS NULL")

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
	err := tx.Offset(int(offset)).Limit(int(query.PageSize)).Find(&methodologies).Error
	return methodologies, total, err
}

func (r *methodologyRepository) FindListByStatus(ctx context.Context, status domain.MethodologyStatus) ([]*domain.Methodology, error) {
	var methodologies []*domain.Methodology
	err := db.WithContext(ctx, r.db).Where("status = ? AND delete_time IS NULL", status).Find(&methodologies).Error
	return methodologies, err
}
