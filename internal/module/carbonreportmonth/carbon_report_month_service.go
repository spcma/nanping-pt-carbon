package carbonreportmonth

import (
	"app/internal/module/carbonreportday"
	"app/internal/shared/entity"
	"app/internal/shared/logger"
	"app/internal/shared/timeutil"
	"context"
	"time"

	"go.uber.org/zap"
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
	repo    CarbonReportMonthRepo
	dayRepo carbonreportday.CarbonReportDayRepo
}

// NewCarbonReportMonthAppService 创建碳报告月报应用服务
func NewCarbonReportMonthAppService(repo CarbonReportMonthRepo, dayRepo carbonreportday.CarbonReportDayRepo) *CarbonReportMonthAppService {
	return &CarbonReportMonthAppService{
		repo:    repo,
		dayRepo: dayRepo,
	}
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

// AggregateMonthlyReport 汇总月报：从日报数据汇总生成月报
func (s *CarbonReportMonthAppService) AggregateMonthlyReport(ctx context.Context, year int, month int) error {
	logger.SchedulerL.Info("开始汇总月报",
		zap.Int("year", year),
		zap.Int("month", month),
	)

	// 1. 检查该月是否已有月报
	existingMonthReport, err := s.repo.FindByMonth(ctx, year, month)
	if err != nil {
		return err
	}
	if existingMonthReport != nil {
		logger.SchedulerL.Info("该月已存在月报，跳过汇总",
			zap.Int("year", year),
			zap.Int("month", month),
		)
		return nil
	}

	// 2. 查询该月的所有日报数据
	dayReports, err := s.dayRepo.FindByMonth(ctx, year, month)
	if err != nil {
		return err
	}

	if len(dayReports) == 0 {
		logger.SchedulerL.Info("该月没有日报数据，跳过汇总",
			zap.Int("year", year),
			zap.Int("month", month),
		)
		return nil
	}

	// 3. 汇总数据
	var totalTurnover, totalBaseline, totalCarbonReduction float64
	for _, dayReport := range dayReports {
		totalTurnover += dayReport.Turnover
		totalBaseline += dayReport.Baseline
		totalCarbonReduction += dayReport.CarbonReduction
	}

	// 4. 创建月报记录
	// 使用当月1号作为采集日期
	collectionDate := timeutil.FromGoTime(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local))
	monthReport, err := NewCarbonReportMonth(
		totalTurnover,
		totalBaseline,
		0, // EnergyConsumption 需要人工填写
		collectionDate,
		1, // 系统自动汇总，使用默认用户ID
	)
	if err != nil {
		return err
	}

	// 设置碳减排量为汇总值
	monthReport.CarbonReduction = totalCarbonReduction

	// 5. 保存月报
	err = s.repo.Create(ctx, monthReport)
	if err != nil {
		return err
	}

	logger.SchedulerL.Info("月报汇总完成",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Float64("totalTurnover", totalTurnover),
		zap.Float64("totalBaseline", totalBaseline),
		zap.Float64("totalCarbonReduction", totalCarbonReduction),
		zap.Int("dayReportCount", len(dayReports)),
	)

	return nil
}
