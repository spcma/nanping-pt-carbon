package application

import (
	"app/internal/module/carbonreportday/domain"
	"app/internal/shared/entity"
	"context"
)

// ===== Repository Ports =====

type CarbonReportDayRepo interface {
	Create(ctx context.Context, user *domain.CarbonReportDay) error
	Update(ctx context.Context, user *domain.CarbonReportDay) error
	UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error // 添加部分更新方法
	Delete(ctx context.Context, id, uid int64) error
	FindByID(ctx context.Context, id int64) (*domain.CarbonReportDay, error)
	FindList(ctx context.Context) ([]*domain.CarbonReportDay, error)
	FindPage(ctx context.Context, pageNum, pageSize int) (*entity.PaginationResult[*domain.CarbonReportDay], error)
}
