package domain

import (
	"app/internal/shared/entity"
	"context"
)

// MethodologyRepository 方法学仓储接口（领域层定义）
type MethodologyRepository interface {
	Create(ctx context.Context, methodology *Methodology) error
	Update(ctx context.Context, methodology *Methodology) error
	Delete(ctx context.Context, id int64, userID int64) error

	FindByID(ctx context.Context, id int64) (*Methodology, error)
	FindByQuery(ctx context.Context, query *MethodologyQuery) (*Methodology, error)
	FindList(ctx context.Context) ([]*Methodology, error)
	FindPage(ctx context.Context, query *MethodologyPageQuery) (*entity.PaginationResult[*Methodology], error)
	FindListByStatus(ctx context.Context, status MethodologyStatus) ([]*Methodology, error)
}

// MethodologyQuery 方法学查询条件
type MethodologyQuery struct {
	ID     *int64             `json:"id"`
	Code   *string            `json:"code"`
	Name   *string            `json:"name"`
	Status *MethodologyStatus `json:"status"`
}

// MethodologyPageQuery 方法学分页查询对象
type MethodologyPageQuery struct {
	entity.PaginationQuery
	Name      string `json:"name"`   // 方法学名模糊匹配
	Code      string `json:"code"`   // 方法学编码精确匹配
	Status    string `json:"status"` // 状态
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"` // "asc" or "desc"
}
