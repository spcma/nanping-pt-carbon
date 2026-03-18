package application

import (
	"app/internal/module/project/domain"
	"context"
)

// ===== Service Ports（服务端口 - 给外部模块用） =====

type ProjectService interface {
	GetProject(ctx context.Context, id int64) (*domain.Project, error)
	GetProjectByCode(ctx context.Context, code string) (*domain.Project, error)
}
