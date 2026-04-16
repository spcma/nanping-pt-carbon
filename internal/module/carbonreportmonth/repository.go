package carbonreportmonth

import (
	"app/internal/shared/entity"
	"context"
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
