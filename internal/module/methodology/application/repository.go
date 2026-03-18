package application

import (
	"app/internal/module/methodology/domain"
	"app/internal/shared/entity"
	"context"
)

// ===== Repository Ports（仓储端口） =====

type MethodologyRepo interface {
	Create(ctx context.Context, methodology *domain.Methodology) error
	Update(ctx context.Context, methodology *domain.Methodology) error
	Delete(ctx context.Context, id, uid int64) error
	FindByID(ctx context.Context, id int64) (*domain.Methodology, error)
	FindByQuery(ctx context.Context, query *domain.MethodologyQuery) (*domain.Methodology, error)
	FindList(ctx context.Context) ([]*domain.Methodology, error)
	FindPage(ctx context.Context, query *domain.MethodologyPageQuery) (*entity.PaginationResult[*domain.Methodology], error)
	FindListByStatus(ctx context.Context, status domain.MethodologyStatus) ([]*domain.Methodology, error)
}
