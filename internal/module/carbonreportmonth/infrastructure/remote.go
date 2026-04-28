package infrastructure

import (
	carbonreportday_app "app/internal/module/carbonreportday/application"
	"app/internal/module/carbonreportmonth/application"
	"context"
)

// CarbonReportDayServiceAdapter 碳日报服务适配器
type CarbonReportDayServiceAdapter struct {
	dayService *carbonreportday_app.CarbonReportDayService
}

func NewAdapter(dayService *carbonreportday_app.CarbonReportDayService) *CarbonReportDayServiceAdapter {
	return &CarbonReportDayServiceAdapter{
		dayService: dayService,
	}
}

// FindByMonth 实现 CarbonReportDayService 接口
func (a *CarbonReportDayServiceAdapter) FindByMonth(ctx context.Context, year int, month int) ([]*application.CarbonReportDaySummary, error) {
	// 调用碳日报服务获取数据
	dayReports, err := a.dayService.FindByMonth(ctx, year, month)
	if err != nil {
		return nil, err
	}

	// 转换为汇总数据类型
	var summaries []*application.CarbonReportDaySummary
	for _, report := range dayReports {
		summaries = append(summaries, &application.CarbonReportDaySummary{
			Turnover:        report.Turnover,
			Baseline:        report.Baseline,
			CarbonReduction: report.CarbonReduction,
		})
	}

	return summaries, nil
}
