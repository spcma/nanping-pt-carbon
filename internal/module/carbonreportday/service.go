package carbonreportday

import (
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
)

// CarbonReportDayService 碳报告日报应用服务
type CarbonReportDayService struct {
	repo CarbonReportDayRepo
}

// NewCarbonReportDayService 创建碳报告日报应用服务
func NewCarbonReportDayService(repo CarbonReportDayRepo) *CarbonReportDayService {
	return &CarbonReportDayService{repo: repo}
}

// CreateCarbonReportDayCommand 创建碳报告日报命令
type CreateCarbonReportDayCommand struct {
	Hash            string        `json:"hash"`
	UserID          int64         `json:"userId"`
	Turnover        float64       `json:"turnover"`
	Baseline        float64       `json:"baseline" gorm:"column:baseline"`                // 基准值
	CarbonReduction float64       `json:"carbonReduction" gorm:"column:carbon_reduction"` // 碳减排量
	CollectionDate  timeutil.Time `json:"collection_date" gorm:"column:collection_date"`  // 数据采集日期
}

// CreateCarbonReportDay 创建碳报告日报
func (s *CarbonReportDayService) CreateCarbonReportDay(ctx context.Context, cmd CreateCarbonReportDayCommand) (int64, error) {
	report, err := NewCarbonReportDay(cmd.Turnover, cmd.Baseline, cmd.CollectionDate, cmd.UserID)
	if err != nil {
		return 0, err
	}

	report.Hash = cmd.Hash
	err = s.repo.Create(ctx, report)
	if err != nil {
		return 0, err
	}
	return report.Id, nil
}

// UpdateCarbonReportDayCommand 更新碳报告日报命令
type UpdateCarbonReportDayCommand struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"userId"`
}

// UpdateCarbonReportDay 更新碳报告日报
func (s *CarbonReportDayService) UpdateCarbonReportDay(ctx context.Context, cmd UpdateCarbonReportDayCommand) error {
	report, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	return report.UpdateInfo(cmd.UserID)
}

// DeleteCarbonReportDay 删除碳报告日报
func (s *CarbonReportDayService) DeleteCarbonReportDay(ctx context.Context, id int64, userID int64) error {
	return s.repo.Delete(ctx, id, userID)
}

// GetCarbonReportDayByID 根据 ID 获取碳报告日报
func (s *CarbonReportDayService) GetCarbonReportDayByID(ctx context.Context, id int64) (*CarbonReportDay, error) {
	return s.repo.FindByID(ctx, id)
}

// GetCarbonReportDayPage 分页查询碳报告日报
func (s *CarbonReportDayService) GetCarbonReportDayPage(ctx context.Context, query *CarbonReportDayPageQuery) (*entity.PaginationResult[CarbonReportDay], error) {
	res, err := s.repo.FindPage(ctx, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}
