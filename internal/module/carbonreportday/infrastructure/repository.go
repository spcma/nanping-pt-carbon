package infrastructure

import (
	"app/internal/module/carbonreportday/domain"
	"app/internal/shared/db"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CarbonReportDayRepository 碳日报仓储实现
type CarbonReportDayRepository struct{}

// NewCarbonReportDayRepository 创建碳日报仓储实例
func NewCarbonReportDayRepository() *CarbonReportDayRepository {
	return &CarbonReportDayRepository{}
}

func (r *CarbonReportDayRepository) Create(ctx context.Context, report *domain.CarbonReportDay) error {
	return db.GetDB(ctx).WithContext(ctx).Create(report).Error
}

func (r *CarbonReportDayRepository) Update(ctx context.Context, report *domain.CarbonReportDay) error {
	return db.GetDB(ctx).WithContext(ctx).Updates(report).Error
}

func (r *CarbonReportDayRepository) Delete(ctx context.Context, id, uid int64) error {
	updates := map[string]interface{}{
		"deleteBy":   uid,
		"deleteTime": timeutil.Now(),
	}
	return db.GetDB(ctx).WithContext(ctx).
		Model(&domain.CarbonReportDay{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *CarbonReportDayRepository) FindByID(ctx context.Context, id int64) (*domain.CarbonReportDay, error) {
	var report domain.CarbonReportDay
	err := db.GetDB(ctx).WithContext(ctx).
		Table("carbon_report_day").
		Where("id = ?", id).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&report).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &report, nil
}

func (r *CarbonReportDayRepository) FindList(ctx context.Context) ([]*domain.CarbonReportDay, error) {
	var reports []*domain.CarbonReportDay
	err := db.GetDB(ctx).WithContext(ctx).Find(&reports).Error
	return reports, err
}

func (r *CarbonReportDayRepository) FindPage(ctx context.Context, query *domain.CarbonReportDayPageQuery) (*entity.PaginationResult[*domain.CarbonReportDay], error) {
	helper := db.NewPaginationHelper[*domain.CarbonReportDay](db.GetDB(ctx))

	result, err := helper.PageQuery(query.PageNum, query.PageSize, func(dq *gorm.DB) *gorm.DB {
		dq = dq.WithContext(ctx).
			Table("carbon_report_day").
			Where(entity.FieldDeleteBy + " = 0")

		// 日期范围过滤
		if query.StartDate != "" {
			dq = dq.Where("collection_date >= ?", query.StartDate)
		}
		if query.EndDate != "" {
			dq = dq.Where("collection_date <= ?", query.EndDate)
		}

		// 排序
		if query.SortBy != "" {
			order := query.SortBy + " " + query.SortOrder
			dq = dq.Order(order)
		} else {
			dq = dq.Order("id DESC")
		}

		return dq
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindByMonth 根据年月查询该月的所有日报数据
func (r *CarbonReportDayRepository) FindByMonth(ctx context.Context, year int, month int) ([]*domain.CarbonReportDay, error) {
	var reports []*domain.CarbonReportDay

	// 计算该月的开始和结束日期
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.Local)

	err := db.GetDB(ctx).WithContext(ctx).
		Table("carbon_report_day").
		Where(entity.FieldDeleteBy+" = 0").
		Where("collection_date >= ? AND collection_date < ?", startDate, endDate).
		Order("collection_date ASC").
		Find(&reports).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return reports, nil
}
