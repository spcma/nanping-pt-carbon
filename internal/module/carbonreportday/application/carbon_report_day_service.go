package application

import (
	"app/internal/module/carbonreportday/domain"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
)

// CreateCarbonReportDayCommand 创建碳报告日报命令
type CreateCarbonReportDayCommand struct {
	// TODO: 根据实际业务需求添加字段
	// ReportDate string `json:"report_date"`
	UserID          int64         `json:"userId"`
	Turnover        float64       `json:"turnover"`
	Baseline        float64       `json:"baseline" gorm:"column:baseline"`                // 基准值
	CarbonReduction float64       `json:"carbonReduction" gorm:"column:carbon_reduction"` // 碳减排量
	CollectionDate  timeutil.Time `json:"collection_date" gorm:"column:collection_date"`  // 数据采集日期
}

// UpdateCarbonReportDayCommand 更新碳报告日报命令
type UpdateCarbonReportDayCommand struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"userId"`
	// TODO: 根据实际业务需求添加字段
}

// CarbonReportDayAppService 碳报告日报应用服务
type CarbonReportDayAppService struct {
	repo CarbonReportDayRepo
}

// NewCarbonReportDayAppService 创建碳报告日报应用服务
func NewCarbonReportDayAppService(repo CarbonReportDayRepo) *CarbonReportDayAppService {
	return &CarbonReportDayAppService{repo: repo}
}

// CreateCarbonReportDay 创建碳报告日报
func (s *CarbonReportDayAppService) CreateCarbonReportDay(ctx context.Context, cmd CreateCarbonReportDayCommand) (int64, error) {
	report, err := domain.NewCarbonReportDay(cmd.Turnover, cmd.Baseline, cmd.CollectionDate, cmd.UserID)
	if err != nil {
		return 0, err
	}
	err = s.repo.Create(ctx, report)
	if err != nil {
		return 0, err
	}
	return report.Id, nil
}

// UpdateCarbonReportDay 更新碳报告日报
func (s *CarbonReportDayAppService) UpdateCarbonReportDay(ctx context.Context, cmd UpdateCarbonReportDayCommand) error {
	report, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	return report.UpdateInfo(cmd.UserID)
}

// DeleteCarbonReportDay 删除碳报告日报
func (s *CarbonReportDayAppService) DeleteCarbonReportDay(ctx context.Context, id int64, userID int64) error {
	return s.repo.Delete(ctx, id, userID)
}

// GetCarbonReportDayByID 根据 ID 获取碳报告日报
func (s *CarbonReportDayAppService) GetCarbonReportDayByID(ctx context.Context, id int64) (*domain.CarbonReportDay, error) {
	return s.repo.FindByID(ctx, id)
}

// GetCarbonReportDayPage 分页查询碳报告日报
func (s *CarbonReportDayAppService) GetCarbonReportDayPage(ctx context.Context, query *domain.CarbonReportDayPageQuery) (*entity.PaginationResult[domain.CarbonReportDay], error) {
	res, err := s.repo.FindPage(ctx, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}
