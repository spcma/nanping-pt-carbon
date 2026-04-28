package infrastructure

import (
	"app/internal/module/carbonreportmonth/domain"
	"app/internal/shared/db"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CarbonReportMonthRepositoryImpl 碳月报仓储实现
type CarbonReportMonthRepositoryImpl struct{}

// NewCarbonReportMonthRepository 创建碳月报仓储实例
func NewCarbonReportMonthRepository() *CarbonReportMonthRepositoryImpl {
	return &CarbonReportMonthRepositoryImpl{}
}

func (r *CarbonReportMonthRepositoryImpl) Create(ctx context.Context, report *domain.CarbonReportMonth) error {
	return db.GetDB(ctx).WithContext(ctx).Create(report).Error
}

func (r *CarbonReportMonthRepositoryImpl) Update(ctx context.Context, report *domain.CarbonReportMonth) error {
	return db.GetDB(ctx).WithContext(ctx).Updates(report).Error
}

func (r *CarbonReportMonthRepositoryImpl) Delete(ctx context.Context, id, uid int64) error {
	updates := map[string]interface{}{
		"deleteBy":   uid,
		"deleteTime": timeutil.Now(),
	}
	return db.GetDB(ctx).WithContext(ctx).
		Model(&domain.CarbonReportMonth{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *CarbonReportMonthRepositoryImpl) FindByID(ctx context.Context, id int64) (*domain.CarbonReportMonth, error) {
	var report domain.CarbonReportMonth
	err := db.GetDB(ctx).WithContext(ctx).
		Table("carbon_report_month").
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

func (r *CarbonReportMonthRepositoryImpl) FindList(ctx context.Context) ([]*domain.CarbonReportMonth, error) {
	var reports []*domain.CarbonReportMonth
	err := db.GetDB(ctx).WithContext(ctx).Find(&reports).Error
	return reports, err
}

func (r *CarbonReportMonthRepositoryImpl) FindPage(ctx context.Context, query *domain.CarbonReportMonthPageQuery) (*entity.PaginationResult[*domain.CarbonReportMonth], error) {
	helper := db.NewPaginationHelper[*domain.CarbonReportMonth](db.GetDB(ctx))

	result, err := helper.PageQuery(query.PageNum, query.PageSize, func(dq *gorm.DB) *gorm.DB {
		dq = dq.WithContext(ctx).
			Table("carbon_report_month").
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

// FindByMonth 根据年月查询月报
func (r *CarbonReportMonthRepositoryImpl) FindByMonth(ctx context.Context, year int, month int) (*domain.CarbonReportMonth, error) {
	var report domain.CarbonReportMonth

	date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	startDate := date
	endDate := date.AddDate(0, 1, 0)

	err := db.GetDB(ctx).WithContext(ctx).
		Table("carbon_report_month").
		Where(entity.FieldDeleteBy+" = 0").
		Where("collection_date >= ? AND collection_date < ?", startDate, endDate).
		Take(&report).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &report, nil
}
