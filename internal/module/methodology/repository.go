package methodology

import (
	"app/internal/shared/entity"
	"context"
)

// ===== Repository Ports（仓储端口） =====

type MethodologyRepo interface {
	Create(ctx context.Context, methodology *Methodology) error
	Update(ctx context.Context, methodology *Methodology) error
	Delete(ctx context.Context, id, uid int64) error
	FindByID(ctx context.Context, id int64) (*Methodology, error)
	FindByQuery(ctx context.Context, query *MethodologyQuery) (*Methodology, error)
	FindList(ctx context.Context) ([]*Methodology, error)
	FindPage(ctx context.Context, query *MethodologyPageQuery) (*entity.PaginationResult[Methodology], error)
	FindListByStatus(ctx context.Context, status MethodologyStatus) ([]*Methodology, error)
}
