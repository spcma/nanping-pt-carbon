package persistence

import (
	"app/internal/module/carbonreportday/domain"
	"app/internal/shared/db"
	"app/internal/shared/timeutil"
	"context"
	"strings"

	"gorm.io/gorm"
)

type carbonReportDayRepository struct {
	db *gorm.DB
}

// NewCarbonReportDayRepository creates carbon report day repository
func NewCarbonReportDayRepository(db *gorm.DB) domain.CarbonReportDayRepository {
	return &carbonReportDayRepository{db: db}
}

func (r *carbonReportDayRepository) Create(ctx context.Context, report *domain.CarbonReportDay) error {
	return db.WithContext(ctx, r.db).Create(report).Error
}

func (r *carbonReportDayRepository) Update(ctx context.Context, report *domain.CarbonReportDay) error {
	return db.WithContext(ctx, r.db).Save(report).Error
}

func (r *carbonReportDayRepository) Delete(ctx context.Context, id int64, userID int64) error {
	// 逻辑删除：只更新 delete_user 和 delete_time 字段
	updates := map[string]interface{}{
		"delete_by":   userID,
		"delete_time": timeutil.New(),
	}
	return db.WithContext(ctx, r.db).Model(&domain.CarbonReportDay{}).Where("id = ?", id).Updates(updates).Error
}

func (r *carbonReportDayRepository) FindByID(ctx context.Context, id int64) (*domain.CarbonReportDay, error) {
	var report domain.CarbonReportDay
	err := db.WithContext(ctx, r.db).Where("id = ? AND delete_time IS NULL", id).First(&report).Error
	return &report, err
}

func (r *carbonReportDayRepository) FindPage(ctx context.Context, query domain.CarbonReportDayPageQuery) ([]*domain.CarbonReportDay, int64, error) {
	var reports []*domain.CarbonReportDay
	var total int64

	tx := db.WithContext(ctx, r.db).Model(&domain.CarbonReportDay{}).Where("delete_time IS NULL")

	// TODO: 根据实际业务需求添加查询条件
	// if query.StartDate != "" {
	// 	tx = tx.Where("report_date >= ?", query.StartDate)
	// }
	// if query.EndDate != "" {
	// 	tx = tx.Where("report_date <= ?", query.EndDate)
	// }

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
	err := tx.Offset(int(offset)).Limit(int(query.PageSize)).Find(&reports).Error
	return reports, total, err
}
