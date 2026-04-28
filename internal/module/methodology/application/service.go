package application

import (
	"app/internal/module/methodology/domain"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"
)

// MethodologyAppService 方法学应用服务
type MethodologyAppService struct {
	repo domain.MethodologyRepository
}

// NewMethodologyAppService 创建方法学应用服务
func NewMethodologyAppService(repo domain.MethodologyRepository) *MethodologyAppService {
	return &MethodologyAppService{repo: repo}
}

// CreateMethodologyCommand 创建方法学命令
type CreateMethodologyCommand struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	UserID      int64  `json:"userId"`
}

// Create 创建方法学
func (s *MethodologyAppService) Create(ctx context.Context, cmd CreateMethodologyCommand) (int64, error) {
	// 调用领域层创建聚合根
	methodology, err := domain.NewMethodology(
		cmd.Name,
		cmd.Code,
		"",
		cmd.Description,
		cmd.UserID,
		timeutil.Now(),
		timeutil.Now(),
	)
	if err != nil {
		return 0, err
	}

	// 持久化
	err = s.repo.Create(ctx, methodology)
	if err != nil {
		return 0, err
	}

	return methodology.Id, nil
}

// UpdateMethodologyCommand 更新方法学命令
type UpdateMethodologyCommand struct {
	ID          int64   `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	UserID      int64   `json:"userId"`
}

// Update 更新方法学
func (s *MethodologyAppService) Update(ctx context.Context, cmd UpdateMethodologyCommand) error {
	// 获取聚合根
	methodology, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if methodology == nil {
		return errors.New("方法学不存在")
	}

	// 调用领域行为
	return methodology.UpdateInfo(cmd.Name, cmd.Description, cmd.UserID)
}

// Delete 删除方法学
func (s *MethodologyAppService) Delete(ctx context.Context, id int64, userID int64) error {
	// 获取聚合根
	methodology, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if methodology == nil {
		return errors.New("方法学不存在")
	}

	// 调用领域行为
	return methodology.Delete(userID)
}

// GetByID 根据 ID 获取方法学
func (s *MethodologyAppService) GetByID(ctx context.Context, id int64) (*domain.Methodology, error) {
	return s.repo.FindByID(ctx, id)
}

// GetByQuery 综合查询
func (s *MethodologyAppService) GetByQuery(ctx context.Context, query *domain.MethodologyQuery) (*domain.Methodology, error) {
	return s.repo.FindByQuery(ctx, query)
}

// GetList 获取方法学列表
func (s *MethodologyAppService) GetList(ctx context.Context) ([]*domain.Methodology, error) {
	return s.repo.FindList(ctx)
}

// GetPage 分页查询方法学
func (s *MethodologyAppService) GetPage(ctx context.Context, query *domain.MethodologyPageQuery) (*entity.PaginationResult[*domain.Methodology], error) {
	return s.repo.FindPage(ctx, query)
}

// ChangeMethodologyStatusCommand 变更方法学状态命令
type ChangeMethodologyStatusCommand struct {
	ID     int64                    `json:"id"`
	Status domain.MethodologyStatus `json:"status"`
	UserID int64                    `json:"userId"`
}

// ChangeStatus 变更方法学状态
func (s *MethodologyAppService) ChangeStatus(ctx context.Context, cmd ChangeMethodologyStatusCommand) error {
	// 获取聚合根
	methodology, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if methodology == nil {
		return errors.New("方法学不存在")
	}

	// 调用领域行为
	return methodology.ChangeStatus(cmd.Status, cmd.UserID)
}

// Activate 启用方法学
func (s *MethodologyAppService) Activate(ctx context.Context, id int64, userID int64) error {
	methodology, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if methodology == nil {
		return errors.New("方法学不存在")
	}

	return methodology.Activate(userID)
}

// Deactivate 禁用方法学
func (s *MethodologyAppService) Deactivate(ctx context.Context, id int64, userID int64) error {
	methodology, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if methodology == nil {
		return errors.New("方法学不存在")
	}

	return methodology.Deactivate(userID)
}
