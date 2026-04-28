package application

import (
	"app/internal/module/carbonreportmonth/domain"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"
	"time"
)

// CarbonReportDayService 碳日报服务接口（用于跨模块调用）
type CarbonReportDayService interface {
	FindByMonth(ctx context.Context, year int, month int) ([]*CarbonReportDaySummary, error)
}

// CarbonReportDaySummary 碳日报汇总数据
type CarbonReportDaySummary struct {
	Turnover        float64
	Baseline        float64
	CarbonReduction float64
}

// CreateCarbonReportMonthCommand 创建碳月报命令
type CreateCarbonReportMonthCommand struct {
	UserID            int64   `json:"userId"`
	CollectionDate    string  `json:"collection_date"`   // 数据采集日期字符串
	Turnover          float64 `json:"turnover"`          // 周转量
	Baseline          float64 `json:"baseline"`          // 基准值
	EnergyConsumption float64 `json:"energyConsumption"` // 能耗
}

// UpdateCarbonReportMonthCommand 更新碳月报命令
type UpdateCarbonReportMonthCommand struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"userId"`
}

// CarbonReportMonthAppService 碳月报应用服务
type CarbonReportMonthAppService struct {
	repo       domain.CarbonReportMonthRepository
	dayService CarbonReportDayService // 通过接口依赖，避免循环依赖
}

// NewCarbonReportMonthAppService 创建碳月报应用服务
func NewCarbonReportMonthAppService(repo domain.CarbonReportMonthRepository, dayService CarbonReportDayService) *CarbonReportMonthAppService {
	return &CarbonReportMonthAppService{
		repo:       repo,
		dayService: dayService,
	}
}

// Create 创建碳月报
func (s *CarbonReportMonthAppService) Create(ctx context.Context, cmd CreateCarbonReportMonthCommand) (int64, error) {
	// 解析日期
	collectionDate, err := parseDate(cmd.CollectionDate)
	if err != nil {
		return 0, errors.New("无效的日期格式")
	}

	// 调用领域层创建聚合根
	report, err := domain.NewCarbonReportMonth(
		cmd.Turnover,
		cmd.Baseline,
		cmd.EnergyConsumption,
		collectionDate,
		cmd.UserID,
	)
	if err != nil {
		return 0, err
	}

	// 持久化
	err = s.repo.Create(ctx, report)
	if err != nil {
		return 0, err
	}

	return report.Id, nil
}

// Update 更新碳月报
func (s *CarbonReportMonthAppService) Update(ctx context.Context, cmd UpdateCarbonReportMonthCommand) error {
	// 获取聚合根
	report, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if report == nil {
		return errors.New("碳月报不存在")
	}

	// 调用领域行为
	return report.UpdateInfo(cmd.UserID)
}

// Delete 删除碳月报
func (s *CarbonReportMonthAppService) Delete(ctx context.Context, id int64, userID int64) error {
	return s.repo.Delete(ctx, id, userID)
}

// GetByID 根据 ID 获取碳月报
func (s *CarbonReportMonthAppService) GetByID(ctx context.Context, id int64) (*domain.CarbonReportMonth, error) {
	return s.repo.FindByID(ctx, id)
}

// GetPage 分页查询碳月报
func (s *CarbonReportMonthAppService) GetPage(ctx context.Context, query *domain.CarbonReportMonthPageQuery) (*entity.PaginationResult[*domain.CarbonReportMonth], error) {
	return s.repo.FindPage(ctx, query)
}

// AggregateMonthlyReport 汇总月报：从日报数据汇总生成月报
func (s *CarbonReportMonthAppService) AggregateMonthlyReport(ctx context.Context, year int, month int) error {
	// 1. 检查该月是否已有月报
	existingMonthReport, err := s.repo.FindByMonth(ctx, year, month)
	if err != nil {
		return err
	}
	if existingMonthReport != nil {
		// 已存在月报，跳过
		return nil
	}

	// 2. 查询该月的所有日报数据（通过接口调用，避免循环依赖）
	if s.dayService == nil {
		return errors.New("碳日报服务未配置")
	}

	dayReports, err := s.dayService.FindByMonth(ctx, year, month)
	if err != nil {
		return err
	}

	if len(dayReports) == 0 {
		// 没有日报数据，跳过
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
	collectionDate, _ := parseDateFromYearMonth(year, month)
	monthReport, err := domain.NewCarbonReportMonth(
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

	return nil
}

// Helper functions for date parsing
func parseDate(dateStr string) (timeutil.Time, error) {
	if dateStr == "" {
		return timeutil.Now(), nil
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		// 尝试另一种格式
		t, err = time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			return timeutil.Time{}, errors.New("无效的日期格式")
		}
	}

	return timeutil.FromGoTime(t), nil
}

func parseDateFromYearMonth(year int, month int) (timeutil.Time, error) {
	t := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	return timeutil.FromGoTime(t), nil
}
