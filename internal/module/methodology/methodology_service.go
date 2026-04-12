package methodology

import (
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
)

// MethodologyAppService 方法学应用服务
type MethodologyAppService struct {
	repo MethodologyRepo
}

// NewMethodologyAppService 创建方法学应用服务
func NewMethodologyAppService(repo MethodologyRepo) *MethodologyAppService {
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
	methodology, err := NewMethodology(cmd.Name, cmd.Code, "", cmd.Description, cmd.UserID, timeutil.Now(), timeutil.Now())
	if err != nil {
		return 0, err
	}
	err = s.repo.Create(ctx, methodology)
	if err != nil {
		return 0, err
	}
	return methodology.Id, nil
}

// UpdateMethodologyCommand 更新方法学命令（使用指针字段支持部分更新）
type UpdateMethodologyCommand struct {
	ID          int64   `json:"id"`
	Name        *string `json:"name"`        // 指针表示可选，nil 表示不更新
	Description *string `json:"description"` // 指针表示可选
	UserID      int64   `json:"userId"`
}

// Update 更新方法学（支持部分字段更新）
func (s *MethodologyAppService) Update(ctx context.Context, cmd UpdateMethodologyCommand) error {
	methodology, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	// 调用领域层的 UpdateInfo，传入指针字段
	return methodology.UpdateInfo(cmd.Name, cmd.Description, cmd.UserID)
}

// Delete 删除方法学
func (s *MethodologyAppService) Delete(ctx context.Context, id int64, userID int64) error {
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

// GetByID 根据 ID 获取方法学
func (s *MethodologyAppService) GetByID(ctx context.Context, id int64) (*Methodology, error) {
	return s.repo.FindByID(ctx, id)
}

type MethodologyQuery struct {
	ID     int64  `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// GetByQuery 综合查询
func (s *MethodologyAppService) GetByQuery(ctx context.Context, query *MethodologyQuery) (*Methodology, error) {

	domainQuery := &MethodologyQuery{
		ID:     query.ID,
		Code:   query.Code,
		Name:   query.Name,
		Status: MethodologyStatus(query.Status),
	}

	return s.repo.FindByQuery(ctx, domainQuery)
}

// GetList 获取方法学列表
func (s *MethodologyAppService) GetList(ctx context.Context) ([]*Methodology, error) {
	return s.repo.FindList(ctx)
}

// GetPage 分页查询方法学
func (s *MethodologyAppService) GetPage(ctx context.Context, query *MethodologyPageQuery) (*entity.PaginationResult[Methodology], error) {
	return s.repo.FindPage(ctx, query)
}

// ChangeMethodologyStatusCommand 变更方法学状态命令
type ChangeMethodologyStatusCommand struct {
	ID     int64             `json:"id"`
	Status MethodologyStatus `json:"status"`
	UserID int64             `json:"userId"`
}

// ChangeStatus 变更方法学状态
func (s *MethodologyAppService) ChangeStatus(ctx context.Context, cmd ChangeMethodologyStatusCommand) error {
	methodology, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	return methodology.ChangeStatus(cmd.Status, cmd.UserID)
}
