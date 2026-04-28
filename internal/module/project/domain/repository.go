package domain

import (
	"app/internal/shared/entity"
	"context"
)

// ProjectRepository 项目仓储接口（领域层定义）
// 注意：接口定义在领域层，实现在基础设施层
type ProjectRepository interface {
	// 基础 CRUD
	Create(ctx context.Context, project *Project) error
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id int64, userID int64) error // 逻辑删除

	// 高频单条件查询
	FindByID(ctx context.Context, id int64) (*Project, error)
	FindByCode(ctx context.Context, code string) (*Project, error)

	// 通用查询（支持多条件组合）
	FindByQuery(ctx context.Context, query *ProjectQuery) (*Project, error)

	// 列表/分页
	FindList(ctx context.Context) ([]*Project, error)
	FindPage(ctx context.Context, query *ProjectPageQuery) (*entity.PaginationResult[*Project], error)
	FindListByStatus(ctx context.Context, status ProjectStatus) ([]*Project, error)
}

// ProjectQuery 项目查询条件（用于单条记录的多条件查询）
type ProjectQuery struct {
	ID     int64         `json:"id" form:"id"`
	Code   string        `json:"code" form:"code"`
	Name   string        `json:"name" form:"name"`
	Status ProjectStatus `json:"status" form:"status"`
}

// ProjectPageQuery 项目分页查询对象
type ProjectPageQuery struct {
	entity.PaginationQuery
	ID        int64         `json:"id" form:"id"`
	Name      string        `json:"name" form:"name"`
	Code      string        `json:"code" form:"code"`
	Status    ProjectStatus `json:"status" form:"status"`
	SortBy    string        `json:"sortBy"`
	SortOrder string        `json:"sortOrder"` // "asc" or "desc"
}
