package application

import (
	"app/internal/config"
	"app/internal/shared/logger"
	"context"
	"time"

	"go.uber.org/zap"
)

func (s *Service) AggregateDailyReport(ctx context.Context, year int, month int) {
	logger.SchedulerL.Info("开始汇总日报",
		zap.Int("year", year),
		zap.Int("month", month),
	)

	now := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)

	dir, err := parseDirByPort(config.GlobalConfig.Ipfs.Port, now)
	if err != nil {
		logger.SchedulerL.Error("获取端口错误", zap.Error(err))
		return
	}

	totalTurnover, err := s.CalcDir(ctx, dir, now.Format(time.DateTime))
	if err != nil {
		logger.SchedulerL.Error("日报汇总错误", zap.Error(err))
		return
	}

	logger.SchedulerL.Info("日报汇总完成",
		zap.Float64("total_turnover", totalTurnover),
	)
}
