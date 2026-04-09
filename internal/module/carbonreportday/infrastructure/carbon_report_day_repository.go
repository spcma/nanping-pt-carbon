package infrastructure

import (
	"app/internal/module/carbonreportday/domain"
	"app/internal/shared/db"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"

	"gorm.io/gorm"
)

type CarbonReportDayRepository struct {
	*db.BaseRepository[domain.CarbonReportDay]
}

func NewCarbonReportDayRepository(_db *gorm.DB) *CarbonReportDayRepository {
	return &CarbonReportDayRepository{
		BaseRepository: db.NewBaseRepository[domain.CarbonReportDay](_db),
	}
}

func (u *CarbonReportDayRepository) Create(ctx context.Context, CarbonReportDay *domain.CarbonReportDay) error {
	return u.GetDB(ctx).WithContext(ctx).Create(CarbonReportDay).Error
}

func (u *CarbonReportDayRepository) Update(ctx context.Context, CarbonReportDay *domain.CarbonReportDay) error {
	return u.GetDB(ctx).WithContext(ctx).Updates(CarbonReportDay).Error
}

func (u *CarbonReportDayRepository) UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error {
	return u.GetDB(ctx).WithContext(ctx).Model(&domain.CarbonReportDay{}).Where("id = ? AND "+entity.FieldDeleteBy+" = 0", id).Updates(updates).Error
}

func (u *CarbonReportDayRepository) Delete(ctx context.Context, id, uid int64) error {
	// 实现逻辑删除，更新状态而不是物理删除
	updates := map[string]interface{}{
		"deleteBy":   uid,
		"deleteTime": timeutil.Now(),
	}
	return u.GetDB(ctx).WithContext(ctx).Model(&domain.CarbonReportDay{}).Where("id = ?", id).Updates(updates).Error
}

func (u *CarbonReportDayRepository) FindByID(ctx context.Context, id int64) (*domain.CarbonReportDay, error) {
	var CarbonReportDay domain.CarbonReportDay
	err := u.GetDB(ctx).WithContext(ctx).
		Table("CarbonReportDays").
		Where("id = ?", id).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&CarbonReportDay).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &CarbonReportDay, nil
}

func (u *CarbonReportDayRepository) FindList(ctx context.Context) ([]*domain.CarbonReportDay, error) {
	var CarbonReportDays []*domain.CarbonReportDay
	err := u.GetDB(ctx).WithContext(ctx).Find(&CarbonReportDays).Error
	return CarbonReportDays, err
}

func (u *CarbonReportDayRepository) FindPage(ctx context.Context, query *domain.CarbonReportDayPageQuery) (*entity.PaginationResult[domain.CarbonReportDay], error) {
	// 使用通用分页助手
	helper := db.NewPaginationHelper[domain.CarbonReportDay](u.GetDB(ctx))
	result, err := helper.PageQuery(int(query.PageNum), int(query.PageSize), func(dq *gorm.DB) *gorm.DB {
		// 构建基础查询 - 使用 delete_by 条件
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
