package domain

import (
	"context"
)

// SysApiRepository API repository interface
type SysApiRepository interface {
	Create(ctx context.Context, api *SysApi) error
	Update(ctx context.Context, api *SysApi) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*SysApi, error)
	FindByCode(ctx context.Context, code string) (*SysApi, error)
	FindPage(ctx context.Context, pageNum, pageSize int64, name string) ([]*SysApi, int64, error)
	FindAll(ctx context.Context) ([]*SysApi, error)
}
