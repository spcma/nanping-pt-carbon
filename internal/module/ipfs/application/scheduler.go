package application

import (
	"app/internal/shared/logger"
	"context"

	"go.uber.org/zap"
)

//	定时任务

func (s *Service) AggregateDailyReport(ctx context.Context, year int, month int) {
	logger.SchedulerL.Info("开始汇总日报",
		zap.Int("year", year),
		zap.Int("month", month),
	)

	//	调用ipfs

	//ctx := context.Background()
	//dir, err := s.CalcDir(ctx, "", "")
	//if err != nil {
	//	logger.SchedulerL.Error("日报汇总错误", zap.Error(err))
	//	return
	//}

}
