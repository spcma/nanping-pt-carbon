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

// CreateProjectCommand 创建项目命令
type CreateProjectCommand struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	UserID      int64  `json:"userId"`
}

// UpdateProjectCommand 更新项目命令
type UpdateProjectCommand struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      int64  `json:"userId"`
}

// ChangeProjectStatusCommand 变更项目状态命令
type ChangeProjectStatusCommand struct {
	ID     int64                `json:"id"`
	Status domain.ProjectStatus `json:"status"`
	UserID int64                `json:"userId"`
}

// ProjectAppService 项目应用服务
type ProjectAppService struct {
	repo ProjectRepo
}

// NewProjectAppService 创建项目应用服务
func NewProjectAppService(repo ProjectRepo) *ProjectAppService {
	return &ProjectAppService{repo: repo}
}

// CreateProject 创建项目
func (s *ProjectAppService) CreateProject(ctx context.Context, cmd CreateProjectCommand) (int64, error) {
	project, err := domain.NewProject(cmd.Name, cmd.Code, cmd.Description, cmd.UserID)
	if err != nil {
		return 0, err
	}
	err = s.repo.Create(ctx, project)
	if err != nil {
		return 0, err
	}
	return project.Id, nil
}

// UpdateProject 更新项目
func (s *ProjectAppService) UpdateProject(ctx context.Context, cmd UpdateProjectCommand) error {
	project, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	return project.UpdateInfo(cmd.Name, cmd.Description, cmd.UserID)
}

// DeleteProject 删除项目
func (s *ProjectAppService) DeleteProject(ctx context.Context, id int64, userID int64) error {
	project, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 执行领域方法
	if err := project.Delete(userID); err != nil {
		return err
	}

	// 持久化
	return s.repo.Update(ctx, project)
}

// GetProjectByID 根据 ID 获取项目
func (s *ProjectAppService) GetProjectByID(ctx context.Context, id int64) (*domain.Project, error) {
	return s.repo.FindByID(ctx, id)
}

// GetProjectByCode 根据编码获取项目
func (s *ProjectAppService) GetProjectByCode(ctx context.Context, code string) (*domain.Project, error) {
	return s.repo.FindByCode(ctx, code)
}

// GetProjectPage 分页查询项目
func (s *ProjectAppService) GetProjectPage(ctx context.Context, query *domain.ProjectPageQuery) ([]*domain.Project, int64, error) {
	return s.repo.FindPage(ctx, query)
}

// ChangeProjectStatus 变更项目状态
func (s *ProjectAppService) ChangeProjectStatus(ctx context.Context, cmd ChangeProjectStatusCommand) error {
	project, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	return project.ChangeStatus(cmd.Status, cmd.UserID)
}

// ===== 实现 Service Ports =====

func (s *ProjectAppService) GetProject(ctx context.Context, id int64) (*domain.Project, error) {
	return s.GetProjectByID(ctx, id)
}

func (s *ProjectAppService) GetProjectByCodeService(ctx context.Context, code string) (*domain.Project, error) {
	return s.GetProjectByCode(ctx, code)
}
