package carbonreportday

import (
	"app/internal/shared/db"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"

	"gorm.io/gorm"
)

type CarbonReportDayRepository struct {
}

func NewCarbonReportDayRepository(_db *gorm.DB) *CarbonReportDayRepository {
	return &CarbonReportDayRepository{}
}

func (u *CarbonReportDayRepository) Create(ctx context.Context, CarbonReportDay *CarbonReportDay) error {
	return db.GetDB(ctx).WithContext(ctx).Create(CarbonReportDay).Error
}

func (u *CarbonReportDayRepository) Update(ctx context.Context, CarbonReportDay *CarbonReportDay) error {
	return db.GetDB(ctx).WithContext(ctx).Updates(CarbonReportDay).Error
}

func (u *CarbonReportDayRepository) UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error {
	return db.GetDB(ctx).WithContext(ctx).
		Model(&CarbonReportDay{}).
		Where("id = ? AND "+entity.FieldDeleteBy+" = 0", id).
		Updates(updates).Error
}

func (u *CarbonReportDayRepository) Delete(ctx context.Context, id, uid int64) error {
	// 实现逻辑删除，更新状态而不是物理删除
	updates := map[string]interface{}{
		"deleteBy":   uid,
		"deleteTime": timeutil.Now(),
	}
	return db.GetDB(ctx).WithContext(ctx).
		Model(&CarbonReportDay{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (u *CarbonReportDayRepository) FindByID(ctx context.Context, id int64) (*CarbonReportDay, error) {
	var carbonReportDay CarbonReportDay
	err := db.GetDB(ctx).WithContext(ctx).
		Table("CarbonReportDays").
		Where("id = ?", id).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&carbonReportDay).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &carbonReportDay, nil
}

func (u *CarbonReportDayRepository) FindList(ctx context.Context) ([]*CarbonReportDay, error) {
	var CarbonReportDays []*CarbonReportDay
	err := db.GetDB(ctx).WithContext(ctx).Find(&CarbonReportDays).Error
	return CarbonReportDays, err
}

func (u *CarbonReportDayRepository) FindPage(ctx context.Context, query *CarbonReportDayPageQuery) (*entity.PaginationResult[CarbonReportDay], error) {

	pageHelper := db.NewPaginationHelper[CarbonReportDay](db.GetDB(ctx))

	result, err := pageHelper.PageQuery(query.PageNum, query.PageSize, func(dq *gorm.DB) *gorm.DB {
		dq = dq.WithContext(ctx).
			Table("carbon_report_day").
			Where(entity.FieldDeleteBy + " = 0")

		return dq
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
