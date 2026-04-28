package server

import (
	carbonreportday_scheduler "app/internal/module/carbonreportday/scheduled"
	"app/internal/module/scheduler"
	"app/internal/shared/logger"

	"go.uber.org/zap"
)

func startSchedulers() *scheduler.Scheduler {
	//	初始化任务注册表
	scheduler.RegistryStore()

	//	初始化调度器
	sched := scheduler.Default()

	// 注册默认任务函数
	scheduler.RegisterTaskFunctions()

	//	注册各个模块的调度任务
	carbonreportday_scheduler.RegisterTask()

	// 从数据库加载已启用的任务
	if err := sched.LoadTasksFromDatabase(); err != nil {
		logger.SchedulerL.Error("Failed to load tasks from database", zap.Error(err))
		panic("Failed to load tasks from database" + err.Error())
	}

	//	启动调度器
	sched.Start()

	return sched
}
