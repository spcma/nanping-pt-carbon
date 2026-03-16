package domain

import "context"

// ProjectPageQuery 项目分页查询对象
type ProjectPageQuery struct {
	PageNum   int64  `json:"pageNum" binding:"required,min=1"`
	PageSize  int64  `json:"pageSize" binding:"required,min=1,max=100"`
	Name      string `json:"name"`   // 项目名模糊匹配
	Code      string `json:"code"`   // 项目编码精确匹配
	Status    string `json:"status"` // 状态
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder" binding:"oneof=asc desc"` // "asc" or "desc"
}

// ProjectRepository 项目仓储接口
type ProjectRepository interface {
	Create(ctx context.Context, project *Project) error
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id int64, userID int64) error // 逻辑删除
	FindByID(ctx context.Context, id int64) (*Project, error)
	FindByCode(ctx context.Context, code string) (*Project, error)
	FindPage(ctx context.Context, query ProjectPageQuery) ([]*Project, int64, error)
	FindListByStatus(ctx context.Context, status ProjectStatus) ([]*Project, error)
}
