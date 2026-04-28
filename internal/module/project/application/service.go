package application

import (
	"app/internal/module/project/domain"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"

	"github.com/dromara/carbon/v2"
)

// ProjectAppService 项目应用服务
type ProjectAppService struct {
	repo domain.ProjectRepository
}

// NewProjectAppService 创建项目应用服务
func NewProjectAppService(repo domain.ProjectRepository) *ProjectAppService {
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
	// 检查项目编码是否已存在
	existing, err := s.repo.FindByCode(ctx, cmd.Code)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, errors.New("项目编码已存在")
	}

	// 解析时间字符串为 timeutil.Time
	startDate, err := parseTime(cmd.StartDate)
	if err != nil {
		return 0, errors.New("开始日期格式错误")
	}
	endDate, err := parseTime(cmd.EndDate)
	if err != nil {
		return 0, errors.New("结束日期格式错误")
	}

	// 调用领域层创建聚合根
	project, err := domain.NewProject(
		cmd.Name,
		cmd.Code,
		cmd.Icon,
		cmd.Description,
		cmd.UserID,
		startDate.(timeutil.Time),
		endDate.(timeutil.Time),
	)
	if err != nil {
		return 0, err
	}

	// 持久化
	err = s.repo.Create(ctx, project)
	if err != nil {
		return 0, err
	}

	return project.Id, nil
}

// UpdateProjectCommand 更新项目命令
type UpdateProjectCommand struct {
	ID          int64   `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	UserID      int64   `json:"userId"`
}

// Update 更新项目
func (s *ProjectAppService) Update(ctx context.Context, cmd UpdateProjectCommand) error {
	// 获取聚合根
	project, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("项目不存在")
	}

	// 调用领域行为
	err = project.UpdateInfo(cmd.Name, cmd.Description, cmd.UserID)
	if err != nil {
		return err
	}

	// 持久化
	return s.repo.Update(ctx, project)
}

// Delete 删除项目
func (s *ProjectAppService) Delete(ctx context.Context, id int64, userID int64) error {
	project, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("项目不存在")
	}

	// 调用领域行为
	err = project.Delete(userID)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, project)
}

// GetByID 根据ID获取项目
func (s *ProjectAppService) GetByID(ctx context.Context, id int64) (*domain.Project, error) {
	return s.repo.FindByID(ctx, id)
}

// GetByQuery 根据查询条件获取项目
func (s *ProjectAppService) GetByQuery(ctx context.Context, query *domain.ProjectQuery) (*domain.Project, error) {
	if query == nil {
		return nil, errors.New("query is nil")
	}

	if query.ID > 0 {
		return s.repo.FindByID(ctx, query.ID)
	}
	if query.Code != "" {
		return s.repo.FindByCode(ctx, query.Code)
	}

	return s.repo.FindByQuery(ctx, query)
}

// GetList 获取项目列表
func (s *ProjectAppService) GetList(ctx context.Context) ([]*domain.Project, error) {
	return s.repo.FindList(ctx)
}

// GetPage 分页查询项目
func (s *ProjectAppService) GetPage(ctx context.Context, query *domain.ProjectPageQuery) (*entity.PaginationResult[*domain.Project], error) {
	return s.repo.FindPage(ctx, query)
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

	// 调用领域行为
	err = project.ChangeStatus(cmd.Status, cmd.UserID)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, project)
}

// ActivateProject 激活项目
func (s *ProjectAppService) ActivateProject(ctx context.Context, id int64, userID int64) error {
	project, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("项目不存在")
	}

	err = project.Activate(userID)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, project)
}

// CompleteProject 完成项目
func (s *ProjectAppService) CompleteProject(ctx context.Context, id int64, userID int64) error {
	project, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("项目不存在")
	}

	err = project.Complete(userID)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, project)
}

// CancelProject 取消项目
func (s *ProjectAppService) CancelProject(ctx context.Context, id int64, userID int64) error {
	project, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("项目不存在")
	}

	err = project.Cancel(userID)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, project)
}

// parseTime 辅助函数：解析时间字符串
func parseTime(timeStr string) (interface{}, error) {
	parse := carbon.Parse(timeStr, carbon.Shanghai)

	if parse.IsValid() {
		return timeutil.Now(parse.StdTime()), nil
	}

	// TODO: 实现时间解析逻辑
	// 这里需要根据实际的 timeutil.Time 类型来实现
	return nil, errors.New("not implemented")
}
