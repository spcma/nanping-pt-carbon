package application

import (
	"app/internal/module/carbonreportday/domain"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"
	"sync"
)

// _defaultService 全局服务实例，供定时任务和其他模块使用
var _defaultService *CarbonReportDayService
var _once sync.Once

// Service 获取默认的碳日报应用服务实例
func Service() *CarbonReportDayService {
	return _defaultService
}

// CarbonReportDayService 碳日报应用服务
type CarbonReportDayService struct {
	repo domain.CarbonReportDayRepository
}

// NewCarbonReportDayService 创建碳日报应用服务
func NewCarbonReportDayService(repo domain.CarbonReportDayRepository) *CarbonReportDayService {
	_once.Do(func() {
		_defaultService = &CarbonReportDayService{repo: repo}
	})
	return _defaultService
}

// CreateCarbonReportDayCommand 创建碳日报命令
type CreateCarbonReportDayCommand struct {
	Hash            string        `json:"hash"`
	UserID          int64         `json:"userId"`
	Turnover        float64       `json:"turnover"`
	Baseline        float64       `json:"baseline"`
	CarbonReduction float64       `json:"carbonReduction"`
	CollectionDate  timeutil.Time `json:"collection_date"`
}

// Create 创建碳日报
func (s *CarbonReportDayService) Create(ctx context.Context, cmd CreateCarbonReportDayCommand) (int64, error) {
	// 调用领域层创建聚合根
	report, err := domain.NewCarbonReportDay(
		cmd.Turnover,
		cmd.Baseline,
		cmd.CollectionDate,
		cmd.UserID,
	)
	if err != nil {
		return 0, err
	}

	// 设置哈希值
	report.SetHash(cmd.Hash)

	// 持久化
	err = s.repo.Create(ctx, report)
	if err != nil {
		return 0, err
	}

	return report.Id, nil
}

// UpdateCarbonReportDayCommand 更新碳日报命令
type UpdateCarbonReportDayCommand struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"userId"`
}

// Update 更新碳日报
func (s *CarbonReportDayService) Update(ctx context.Context, cmd UpdateCarbonReportDayCommand) error {
	// 获取聚合根
	report, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if report == nil {
		return errors.New("碳日报不存在")
	}

	// 调用领域行为
	return report.UpdateInfo(cmd.UserID)
}

// Delete 删除碳日报
func (s *CarbonReportDayService) Delete(ctx context.Context, id int64, userID int64) error {
	return s.repo.Delete(ctx, id, userID)
}

// GetByID 根据 ID 获取碳日报
func (s *CarbonReportDayService) GetByID(ctx context.Context, id int64) (*domain.CarbonReportDay, error) {
	return s.repo.FindByID(ctx, id)
}

// GetPage 分页查询碳日报
func (s *CarbonReportDayService) GetPage(ctx context.Context, query *domain.CarbonReportDayPageQuery) (*entity.PaginationResult[*domain.CarbonReportDay], error) {
	return s.repo.FindPage(ctx, query)
}

// FindByMonth 根据年月查询该月的所有日报数据（供其他模块调用）
func (s *CarbonReportDayService) FindByMonth(ctx context.Context, year int, month int) ([]*domain.CarbonReportDay, error) {
	return s.repo.FindByMonth(ctx, year, month)
}
