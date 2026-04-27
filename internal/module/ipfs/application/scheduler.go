package application

import (
	"app/internal/module/scheduler"
	"app/internal/shared/logger"
	"context"
	"time"

	"go.uber.org/zap"
)

func RegisterTask() {
	// 注册每日碳计算任务
	scheduler.RegisterTask("carbon.report.day", func(ctx context.Context, params map[string]interface{}) error {
		logger.SchedulerL.Info("执行每日碳计算任务")
		cst := time.Now()

		// 从参数中获取年份和月份，如果没有则使用上个月
		var year, month, day int
		if yearParam, ok := params["year"]; ok {
			if y, valid := yearParam.(float64); valid {
				year = int(y)
			}
		}
		if monthParam, ok := params["month"]; ok {
			if m, valid := monthParam.(float64); valid {
				month = int(m)
			}
		}
		if dayParam, ok := params["day"]; ok {
			if d, valid := dayParam.(float64); valid {
				day = int(d)
			}
		}

		// 如果参数中没有指定，则默认使用昨日
		if year == 0 || month == 0 {
			now := time.Now()
			lastDay := now.AddDate(0, 0, -1)
			year = lastDay.Year()
			month = int(lastDay.Month())
			day = lastDay.Day()
		}

		defaultService.AggregateDailyReport(ctx, year, month, day)

		logger.SchedulerL.Info("碳日报汇总任务执行成功",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Int("day", day),
			zap.Duration("duration", time.Since(cst)),
		)
		return nil
	})
}

func (s *Service) AggregateDailyReport(ctx context.Context, year int, month int, day int) {
	logger.SchedulerL.Info("开始汇总日报",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Int("day", day),
	)

	now := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)

	dir, err := parseDirByPort(4800, now)
	if err != nil {
		logger.SchedulerL.Error("获取端口错误", zap.Error(err))
		return
	}

	//	调度任务使用4800客户端
	totalTurnover, err := s.CalcDir(ctx, "4800", dir, now.Format(time.DateTime))
	if err != nil {
		logger.SchedulerL.Error("日报汇总错误", zap.Error(err))
		return
	}

	logger.SchedulerL.Info("日报汇总完成",
		zap.Float64("total_turnover", totalTurnover),
	)
}
