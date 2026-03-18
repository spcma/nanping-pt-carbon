package application

import (
	"app/internal/module/project/domain"
	"app/internal/shared/entity"
	"context"
)

// ProjectRepo 项目仓储接口
type ProjectRepo interface {
	Create(ctx context.Context, project *domain.Project) error
	Update(ctx context.Context, project *domain.Project) error
	Delete(ctx context.Context, id int64, userID int64) error // 逻辑删除
	FindByID(ctx context.Context, id int64) (*domain.Project, error)
	FindByCode(ctx context.Context, code string) (*domain.Project, error)
	FindList(ctx context.Context) ([]*domain.Project, error)
	FindPage(ctx context.Context, query *domain.ProjectPageQuery) (*entity.PaginationResult[*domain.Project], error)
	FindListByStatus(ctx context.Context, status domain.ProjectStatus) ([]*domain.Project, error)
}

type ProjectMembersRepo interface {
	Create(ctx context.Context, projectMembers *domain.ProjectMembers) error
	Update(ctx context.Context, projectMembers *domain.ProjectMembers) error
	Delete(ctx context.Context, id int64, userID int64) error
	FindByID(ctx context.Context, id int64) (*domain.ProjectMembers, error)
	FindByProjectID(ctx context.Context, projectID int64) ([]*domain.ProjectMembers, error)
	FindByUserID(ctx context.Context, userID int64) ([]*domain.ProjectMembers, error)
	FindByProjectAndUser(ctx context.Context, projectID, userID int64) (*domain.ProjectMembers, error)
	FindPage(ctx context.Context, query *domain.ProjectMembersPageQuery) (*entity.PaginationResult[*domain.ProjectMembers], error)
}
