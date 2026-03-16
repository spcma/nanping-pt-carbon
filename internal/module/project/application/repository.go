package application

import (
	"app/internal/module/project/domain"
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
	FindPage(ctx context.Context, query *domain.ProjectPageQuery) ([]*domain.Project, int64, error)
	FindListByStatus(ctx context.Context, status domain.ProjectStatus) ([]*domain.Project, error)
}
