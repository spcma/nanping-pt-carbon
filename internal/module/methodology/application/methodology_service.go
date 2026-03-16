package application

import (
	"app/internal/module/methodology/domain"
	"context"
)

// CreateMethodologyCommand 创建方法学命令
type CreateMethodologyCommand struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	UserID      int64  `json:"userId"`
}

// UpdateMethodologyCommand 更新方法学命令
type UpdateMethodologyCommand struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      int64  `json:"userId"`
}

// ChangeMethodologyStatusCommand 变更方法学状态命令
type ChangeMethodologyStatusCommand struct {
	ID     int64                    `json:"id"`
	Status domain.MethodologyStatus `json:"status"`
	UserID int64                    `json:"userId"`
}

// MethodologyAppService 方法学应用服务
type MethodologyAppService struct {
	repo MethodologyRepo
}

// NewMethodologyAppService 创建方法学应用服务
func NewMethodologyAppService(repo MethodologyRepo) *MethodologyAppService {
	return &MethodologyAppService{repo: repo}
}

// CreateMethodology 创建方法学
func (s *MethodologyAppService) CreateMethodology(ctx context.Context, cmd CreateMethodologyCommand) (int64, error) {
	methodology, err := domain.NewMethodology(cmd.Name, cmd.Code, cmd.Description, cmd.UserID)
	if err != nil {
		return 0, err
	}
	err = s.repo.Create(ctx, methodology)
	if err != nil {
		return 0, err
	}
	return methodology.Id, nil
}

// UpdateMethodology 更新方法学
func (s *MethodologyAppService) UpdateMethodology(ctx context.Context, cmd UpdateMethodologyCommand) error {
	methodology, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	return methodology.UpdateInfo(cmd.Name, cmd.Description, cmd.UserID)
}

// DeleteMethodology 删除方法学
func (s *MethodologyAppService) DeleteMethodology(ctx context.Context, id int64, userID int64) error {
	methodology, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 执行领域方法
	if err := methodology.Delete(userID); err != nil {
		return err
	}

	// 持久化
	return s.repo.Update(ctx, methodology)
}

// GetMethodologyByID 根据 ID 获取方法学
func (s *MethodologyAppService) GetMethodologyByID(ctx context.Context, id int64) (*domain.Methodology, error) {
	return s.repo.FindByID(ctx, id)
}

// GetMethodologyByCode 根据编码获取方法学
func (s *MethodologyAppService) GetMethodologyByCode(ctx context.Context, code string) (*domain.Methodology, error) {
	return s.repo.FindByCode(ctx, code)
}

// GetMethodologyPage 分页查询方法学
func (s *MethodologyAppService) GetMethodologyPage(ctx context.Context, query domain.MethodologyPageQuery) ([]*domain.Methodology, int64, error) {
	return s.repo.FindPage(ctx, query)
}

// ChangeMethodologyStatus 变更方法学状态
func (s *MethodologyAppService) ChangeMethodologyStatus(ctx context.Context, cmd ChangeMethodologyStatusCommand) error {
	methodology, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	return methodology.ChangeStatus(cmd.Status, cmd.UserID)
}

// ===== 实现 Service Ports =====

func (s *MethodologyAppService) GetMethodology(ctx context.Context, id int64) (*domain.Methodology, error) {
	return s.GetMethodologyByID(ctx, id)
}

func (s *MethodologyAppService) GetMethodologyByCodeService(ctx context.Context, code string) (*domain.Methodology, error) {
	return s.GetMethodologyByCode(ctx, code)
}
