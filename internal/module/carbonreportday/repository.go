package carbonreportday

import (
	"app/internal/shared/entity"
	"context"
)

// ===== Repository Ports =====

type CarbonReportDayRepo interface {
	Create(ctx context.Context, user *CarbonReportDay) error
	Update(ctx context.Context, user *CarbonReportDay) error
	UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error // 添加部分更新方法
	Delete(ctx context.Context, id, uid int64) error
	FindByID(ctx context.Context, id int64) (*CarbonReportDay, error)
	FindList(ctx context.Context) ([]*CarbonReportDay, error)
	FindPage(ctx context.Context, query *CarbonReportDayPageQuery) (*entity.PaginationResult[CarbonReportDay], error)
}
