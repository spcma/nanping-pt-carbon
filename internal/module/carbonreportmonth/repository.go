package carbonreportmonth

import (
	"app/internal/shared/db"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// ===== Repository Ports =====

type CarbonReportMonthRepo interface {
	Create(ctx context.Context, carbonReportMonth *CarbonReportMonth) error
	Update(ctx context.Context, carbonReportMonth *CarbonReportMonth) error
	UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error // 添加部分更新方法
	Delete(ctx context.Context, id, uid int64) error
	FindByID(ctx context.Context, id int64) (*CarbonReportMonth, error)
	FindList(ctx context.Context) ([]*CarbonReportMonth, error)
	FindPage(ctx context.Context, query *CarbonReportMonthPageQuery) (*entity.PaginationResult[*CarbonReportMonth], error)
	// FindByMonth 根据年月查询月报
	FindByMonth(ctx context.Context, year int, month int) (*CarbonReportMonth, error)
}

type CarbonReportMonthRepository struct {
}

func NewCarbonReportMonthRepository(_db *gorm.DB) *CarbonReportMonthRepository {
	return &CarbonReportMonthRepository{}
}

func (r *CarbonReportMonthRepository) Create(ctx context.Context, carbonReportMonth *CarbonReportMonth) error {
	return db.GetDB(ctx).WithContext(ctx).Create(carbonReportMonth).Error
}

func (r *CarbonReportMonthRepository) Update(ctx context.Context, carbonReportMonth *CarbonReportMonth) error {
	return db.GetDB(ctx).WithContext(ctx).Updates(carbonReportMonth).Error
}

func (r *CarbonReportMonthRepository) UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error {
	return db.GetDB(ctx).WithContext(ctx).
		Model(&CarbonReportMonth{}).
		Where("id = ? AND "+entity.FieldDeleteBy+" = 0", id).
		Updates(updates).Error
}

func (r *CarbonReportMonthRepository) Delete(ctx context.Context, id, uid int64) error {
	// 实现逻辑删除，更新状态而不是物理删除
	updates := map[string]interface{}{
		"deleteBy":   uid,
		"deleteTime": timeutil.Now(),
	}
	return db.GetDB(ctx).WithContext(ctx).
		Model(&CarbonReportMonth{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *CarbonReportMonthRepository) FindByID(ctx context.Context, id int64) (*CarbonReportMonth, error) {
	var carbonReportMonth CarbonReportMonth
	err := db.GetDB(ctx).WithContext(ctx).
		Table("carbon_report_month").
		Where("id = ?", id).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&carbonReportMonth).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &carbonReportMonth, nil
}

func (r *CarbonReportMonthRepository) FindList(ctx context.Context) ([]*CarbonReportMonth, error) {
	var carbonReportMonths []*CarbonReportMonth
	err := db.GetDB(ctx).WithContext(ctx).Find(&carbonReportMonths).Error
	return carbonReportMonths, err
}

func (r *CarbonReportMonthRepository) FindPage(ctx context.Context, query *CarbonReportMonthPageQuery) (*entity.PaginationResult[*CarbonReportMonth], error) {

	pageHelper := db.NewPaginationHelper[*CarbonReportMonth](db.GetDB(ctx))

	result, err := pageHelper.PageQuery(query.PageNum, query.PageSize, func(dq *gorm.DB) *gorm.DB {
		dq = dq.WithContext(ctx).
			Table("carbon_report_month").
			Where(entity.FieldDeleteBy + " = 0")

		return dq
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindByMonth 根据年月查询月报
func (r *CarbonReportMonthRepository) FindByMonth(ctx context.Context, year int, month int) (*CarbonReportMonth, error) {
	var carbonReportMonth CarbonReportMonth

	date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	startDate := date
	endDate := date.AddDate(0, 1, 0)

	err := db.GetDB(ctx).WithContext(ctx).
		Table("carbon_report_month").
		Where(entity.FieldDeleteBy+" = 0").
		Where("collection_date >= ? AND collection_date < ?", startDate, endDate).
		Take(&carbonReportMonth).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &carbonReportMonth, nil
}
