package project

import (
	"app/internal/shared/entity"
	"context"
)

// ProjectRepo 项目仓储接口
type ProjectRepo interface {
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

type ProjectMembersRepo interface {
	Create(ctx context.Context, projectMembers *ProjectMembers) error
	Update(ctx context.Context, projectMembers *ProjectMembers) error
	Delete(ctx context.Context, id int64, userID int64) error
	FindByID(ctx context.Context, id int64) (*ProjectMembers, error)
	FindByProjectID(ctx context.Context, projectID int64) ([]*ProjectMembers, error)
	FindByUserID(ctx context.Context, userID int64) ([]*ProjectMembers, error)
	FindByProjectAndUser(ctx context.Context, projectID, userID int64) (*ProjectMembers, error)
	FindPage(ctx context.Context, query *ProjectMembersPageQuery) (*entity.PaginationResult[*ProjectMembers], error)
}
