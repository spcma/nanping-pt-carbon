package domain

import "context"

// MethodologyPageQuery 方法学分页查询对象
type MethodologyPageQuery struct {
	PageNum   int64  `json:"pageNum" binding:"required,min=1"`
	PageSize  int64  `json:"pageSize" binding:"required,min=1,max=100"`
	Name      string `json:"name"`   // 方法学名模糊匹配
	Code      string `json:"code"`   // 方法学编码精确匹配
	Status    string `json:"status"` // 状态
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder" binding:"oneof=asc desc"` // "asc" or "desc"
}

// MethodologyRepository 方法学仓储接口
type MethodologyRepository interface {
	Create(ctx context.Context, methodology *Methodology) error
	Update(ctx context.Context, methodology *Methodology) error
	Delete(ctx context.Context, id int64, userID int64) error
	FindByID(ctx context.Context, id int64) (*Methodology, error)
	FindByCode(ctx context.Context, code string) (*Methodology, error)
	FindPage(ctx context.Context, query MethodologyPageQuery) ([]*Methodology, int64, error)
	FindListByStatus(ctx context.Context, status MethodologyStatus) ([]*Methodology, error)
}
