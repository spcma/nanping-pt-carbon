package application

import (
	"app/internal/module/project/domain"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"
)

// ProjectAppService 项目应用服务
type ProjectAppService struct {
	repo ProjectRepo
}

// NewProjectService 创建项目应用服务
func NewProjectService(repo ProjectRepo) *ProjectAppService {
	return &ProjectAppService{repo: repo}
}

// CreateProjectCommand 创建项目命令
type CreateProjectCommand struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	UserID      int64  `json:"userId"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
}

// Create 创建项目
func (s *ProjectAppService) Create(ctx context.Context, cmd CreateProjectCommand) (int64, error) {
	findByCode, err := s.repo.FindByCode(ctx, cmd.Code)
	if err != nil {
		return 0, err
	}

	if findByCode != nil {
		return 0, errors.New("项目编码已存在")
	}

	project, err := domain.NewProject(cmd.Name, cmd.Code, "", cmd.Description, cmd.UserID, timeutil.Now(), timeutil.Now())
	if err != nil {
		return 0, err
	}

	err = s.repo.Create(ctx, project)
	if err != nil {
		return 0, err
	}

	return project.Id, nil
}

// UpdateProjectCommand 更新项目命令（使用指针字段支持部分更新）
type UpdateProjectCommand struct {
	ID          int64   `json:"id"`
	Name        *string `json:"name"`        // 指针表示可选，nil 表示不更新
	Description *string `json:"description"` // 指针表示可选
	UserID      int64   `json:"userId"`
}

// Update 更新项目（支持部分字段更新）
func (s *ProjectAppService) Update(ctx context.Context, cmd UpdateProjectCommand) error {
	project, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if project == nil {
		return errors.New("项目不存在")
	}

	// 调用领域层的 UpdateInfo，传入指针字段
	err = project.UpdateInfo(cmd.Name, cmd.Description, cmd.UserID)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, project)
}

// Delete 删除项目
func (s *ProjectAppService) Delete(ctx context.Context, id int64, userID int64) error {
	project, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := project.Delete(userID); err != nil {
		return err
	}

	return s.repo.Update(ctx, project)
}

type ProjectQuery struct {
	ID     int64                `json:"id" form:"id"`
	Code   string               `json:"code" form:"code"`
	Name   string               `json:"name" form:"name"`
	Status domain.ProjectStatus `json:"status" form:"status"`
}

// GetByQuery 根据查询条件获取项目
func (s *ProjectAppService) GetByQuery(ctx context.Context, query *ProjectQuery) (*domain.Project, error) {
	if query == nil {
		return nil, errors.New("query is nil")
	}

	if query.ID > 0 {
		return s.repo.FindByID(ctx, query.ID)
	}
	if query.Code != "" {
		return s.repo.FindByCode(ctx, query.Code)
	}

	domainQuery := domain.ProjectQuery{
		ID:     query.ID,
		Code:   query.Code,
		Name:   query.Name,
		Status: query.Status,
	}
	return s.repo.FindByQuery(ctx, domainQuery)
}

func (s *ProjectAppService) GetList(ctx context.Context) ([]*domain.Project, error) {
	list, err := s.repo.FindList(ctx)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetPage 分页查询项目
func (s *ProjectAppService) GetPage(ctx context.Context, query *domain.ProjectPageQuery) (*entity.PaginationResult[*domain.Project], error) {
	res, err := s.repo.FindPage(ctx, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ChangeProjectStatusCommand 变更项目状态命令
type ChangeProjectStatusCommand struct {
	ID     int64                `json:"id"`
	Status domain.ProjectStatus `json:"status"`
	UserID int64                `json:"userId"`
}

// ChangeStatus 变更项目状态
func (s *ProjectAppService) ChangeStatus(ctx context.Context, cmd ChangeProjectStatusCommand) error {
	project, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if project == nil {
		return errors.New("项目不存在")
	}

	err = project.ChangeStatus(cmd.Status, cmd.UserID)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, project)
}
