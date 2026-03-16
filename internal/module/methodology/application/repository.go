package application

import (
	"app/internal/module/methodology/domain"
	"context"
)

// ===== Repository Ports（仓储端口） =====

type MethodologyRepo interface {
	Create(ctx context.Context, methodology *domain.Methodology) error
	Update(ctx context.Context, methodology *domain.Methodology) error
	Delete(ctx context.Context, id, uid int64) error
	FindByID(ctx context.Context, id int64) (*domain.Methodology, error)
	FindByCode(ctx context.Context, code string) (*domain.Methodology, error)
	FindPage(ctx context.Context, query domain.MethodologyPageQuery) ([]*domain.Methodology, int64, error)
	FindListByStatus(ctx context.Context, status domain.MethodologyStatus) ([]*domain.Methodology, error)
}
