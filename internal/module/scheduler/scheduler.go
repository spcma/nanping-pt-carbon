package scheduler

import (
	"app/internal/shared/logger"
	"context"
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// TaskFunc 定时任务函数类型
type TaskFunc func(ctx context.Context) error

// TaskConfig 定时任务配置
type TaskConfig struct {
	Name        string   // 任务名称
	CronSpec    string   // Cron 表达式
	TaskFunc    TaskFunc // 任务执行函数
	Description string   // 任务描述
	Enabled     bool     // 是否启用
}

// TaskStatus 任务状态
type TaskStatus struct {
	Name        string `json:"name"`
	CronSpec    string `json:"cron_spec"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	LastRun     string `json:"last_run"`
	NextRun     string `json:"next_run"`
	TotalRuns   int    `json:"total_runs"`
}

// Scheduler 定时任务调度器
type Scheduler struct {
	cron      *cron.Cron
	tasks     map[string]*TaskConfig
	taskStats map[string]*TaskStatus
	mu        sync.RWMutex
	entryIDs  map[string]cron.EntryID
}

var (
	defaultScheduler *Scheduler
	once             sync.Once
)

// Default 获取默认调度器实例
func Default() *Scheduler {
	once.Do(func() {
		defaultScheduler = NewScheduler()
	})
	return defaultScheduler
}

// NewScheduler 创建新的调度器
func NewScheduler() *Scheduler {
	c := cron.New(
		cron.WithSeconds(), // 支持秒级精度
		cron.WithChain(
			cron.Recover(cron.DefaultLogger), // 自动恢复 panic
		),
	)

	s := &Scheduler{
		cron:      c,
		tasks:     make(map[string]*TaskConfig),
		taskStats: make(map[string]*TaskStatus),
		entryIDs:  make(map[string]cron.EntryID),
	}

	return s
}

// AddTask 添加定时任务
func (s *Scheduler) AddTask(config *TaskConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !config.Enabled {
		logger.SchedulerL.Info("Task disabled, skipping",
			zap.String("task_name", config.Name),
		)
		return nil
	}

	// 检查任务是否已存在
	if _, exists := s.tasks[config.Name]; exists {
		return fmt.Errorf("task %s already exists", config.Name)
	}

	// 初始化任务状态
	s.taskStats[config.Name] = &TaskStatus{
		Name:        config.Name,
		CronSpec:    config.CronSpec,
		Description: config.Description,
		Enabled:     config.Enabled,
		TotalRuns:   0,
	}

	// 包装任务函数，添加统计和日志
	wrappedFunc := func() {
		s.executeTask(config)
	}

	// 添加到 cron 调度器
	entryID, err := s.cron.AddFunc(config.CronSpec, wrappedFunc)
	if err != nil {
		return fmt.Errorf("failed to add task %s: %w", config.Name, err)
	}

	s.tasks[config.Name] = config
	s.entryIDs[config.Name] = entryID

	logger.SchedulerL.Info("Task added successfully",
		zap.String("task_name", config.Name),
		zap.String("cron_spec", config.CronSpec),
		zap.String("description", config.Description),
	)

	return nil
}

// RemoveTask 移除定时任务
func (s *Scheduler) RemoveTask(taskName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.tasks[taskName]
	if !exists {
		return fmt.Errorf("task %s not found", taskName)
	}

	// 从 cron 中移除
	entryID, ok := s.entryIDs[taskName]
	if ok {
		s.cron.Remove(entryID)
		delete(s.entryIDs, taskName)
	}

	delete(s.tasks, taskName)
	delete(s.taskStats, taskName)

	logger.SchedulerL.Info("Task removed",
		zap.String("task_name", taskName),
	)

	return nil
}

// UpdateTask 更新定时任务
func (s *Scheduler) UpdateTask(config *TaskConfig) error {
	// 先移除旧任务
	if err := s.RemoveTask(config.Name); err != nil && config.Enabled {
		// 如果任务不存在且需要启用，则直接添加
		if err.Error() != fmt.Sprintf("task %s not found", config.Name) {
			return err
		}
	}

	// 如果启用，则添加新任务
	if config.Enabled {
		return s.AddTask(config)
	}

	return nil
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.cron.Start()
	logger.SchedulerL.Info("Scheduler started")
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cron.Stop()
	logger.SchedulerL.Info("Scheduler stopped")
}

// GetStatus 获取任务状态
func (s *Scheduler) GetStatus(taskName string) (*TaskStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status, exists := s.taskStats[taskName]
	if !exists {
		return nil, fmt.Errorf("task %s not found", taskName)
	}

	// 获取下次执行时间
	if entryID, ok := s.entryIDs[taskName]; ok {
		entry := s.cron.Entry(entryID)
		status.NextRun = entry.Next.Format("2006-01-02 15:04:05")
		if !entry.Prev.IsZero() {
			status.LastRun = entry.Prev.Format("2006-01-02 15:04:05")
		}
	}

	return status, nil
}

// ListTasks 列出所有任务状态
func (s *Scheduler) ListTasks() []*TaskStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	statuses := make([]*TaskStatus, 0, len(s.taskStats))
	for _, status := range s.taskStats {
		// 复制状态以避免并发问题
		statusCopy := *status
		statuses = append(statuses, &statusCopy)
	}

	return statuses
}

// executeTask 执行任务（内部方法）
func (s *Scheduler) executeTask(config *TaskConfig) {
	s.mu.Lock()
	if stats, exists := s.taskStats[config.Name]; exists {
		stats.TotalRuns++
		stats.LastRun = config.Name // 临时标记，实际时间在任务执行后更新
	}
	s.mu.Unlock()

	logger.SchedulerL.Info("Executing task",
		zap.String("task_name", config.Name),
	)

	// 执行任务
	ctx := context.Background()
	err := config.TaskFunc(ctx)

	if err != nil {
		logger.SchedulerL.Error("Task execution failed",
			zap.String("task_name", config.Name),
			zap.Error(err),
		)
	} else {
		logger.SchedulerL.Info("Task executed successfully",
			zap.String("task_name", config.Name),
		)
	}

	// 更新最后执行时间
	s.mu.Lock()
	if stats, exists := s.taskStats[config.Name]; exists {
		stats.LastRun = config.Name
	}
	s.mu.Unlock()
}

// EnableTask 启用任务
func (s *Scheduler) EnableTask(taskName string) error {
	s.mu.RLock()
	config, exists := s.tasks[taskName]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("task %s not found", taskName)
	}

	config.Enabled = true
	return s.UpdateTask(config)
}

// DisableTask 禁用任务
func (s *Scheduler) DisableTask(taskName string) error {
	s.mu.RLock()
	config, exists := s.tasks[taskName]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("task %s not found", taskName)
	}

	config.Enabled = false
	return s.RemoveTask(taskName)
}
