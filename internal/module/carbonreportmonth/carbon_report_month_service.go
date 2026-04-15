package carbonreportmonth

import (
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
)

// CreateCarbonReportMonthCommand 创建碳报告月报命令
type CreateCarbonReportMonthCommand struct {
	UserID            int64         `json:"userId"`
	CollectionDate    timeutil.Time `json:"collection_date"`   // 数据采集日期
	Turnover          float64       `json:"turnover"`          // 周转量
	Baseline          float64       `json:"baseline"`          // 基准值
	EnergyConsumption float64       `json:"energyConsumption"` // 能耗, 人工填入
	CarbonReduction   float64       `json:"carbonReduction"`   // 碳减排量
}

// UpdateCarbonReportMonthCommand 更新碳报告月报命令
type UpdateCarbonReportMonthCommand struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"userId"`
}

// CarbonReportMonthAppService 碳报告月报应用服务
type CarbonReportMonthAppService struct {
	repo CarbonReportMonthRepo
}

// NewCarbonReportMonthAppService 创建碳报告月报应用服务
func NewCarbonReportMonthAppService(repo CarbonReportMonthRepo) *CarbonReportMonthAppService {
	return &CarbonReportMonthAppService{repo: repo}
}

// CreateCarbonReportMonth 创建碳报告月报
func (s *CarbonReportMonthAppService) CreateCarbonReportMonth(ctx context.Context, cmd CreateCarbonReportMonthCommand) (int64, error) {
	report, err := NewCarbonReportMonth(cmd.Turnover, cmd.Baseline, cmd.EnergyConsumption, cmd.CollectionDate, cmd.UserID)
	if err != nil {
		return 0, err
	}

	err = s.repo.Create(ctx, report)
	if err != nil {
		return 0, err
	}
	return report.Id, nil
}

// UpdateCarbonReportMonth 更新碳报告月报
func (s *CarbonReportMonthAppService) UpdateCarbonReportMonth(ctx context.Context, cmd UpdateCarbonReportMonthCommand) error {
	report, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	return report.UpdateInfo(cmd.UserID)
}

// DeleteCarbonReportMonth 删除碳报告月报
func (s *CarbonReportMonthAppService) DeleteCarbonReportMonth(ctx context.Context, id int64, userID int64) error {
	return s.repo.Delete(ctx, id, userID)
}

// GetCarbonReportMonthByID 根据 ID 获取碳报告月报
func (s *CarbonReportMonthAppService) GetCarbonReportMonthByID(ctx context.Context, id int64) (*CarbonReportMonth, error) {
	return s.repo.FindByID(ctx, id)
}

// GetCarbonReportMonthPage 分页查询碳报告月报
func (s *CarbonReportMonthAppService) GetCarbonReportMonthPage(ctx context.Context, query *CarbonReportMonthPageQuery) (*entity.PaginationResult[*CarbonReportMonth], error) {
	res, err := s.repo.FindPage(ctx, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}
