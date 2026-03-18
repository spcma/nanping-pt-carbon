package application

import (
	"app/internal/module/ipfs/domain"
	"app/internal/shared/entity"
	"context"
)

// ===== Repository Ports =====

type IpfsDetailRepo interface {
	Create(ctx context.Context, ipfsDetail *domain.IpfsDetail) error
	Update(ctx context.Context, ipfsDetail *domain.IpfsDetail) error
	UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error
	Delete(ctx context.Context, id, uid int64) error
	FindByID(ctx context.Context, id int64) (*domain.IpfsDetail, error)
	FindByDeviceCode(ctx context.Context, deviceCode string) ([]*domain.IpfsDetail, error)
	FindByFilename(ctx context.Context, filename string) (*domain.IpfsDetail, error)
	FindList(ctx context.Context) ([]*domain.IpfsDetail, error)
	FindPage(ctx context.Context, query *domain.IpfsDetailPageQuery) (*entity.PaginationResult[*domain.IpfsDetail], error)
}
