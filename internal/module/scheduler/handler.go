package scheduler

import (
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SchedulerHandler 定时任务处理器
type SchedulerHandler struct {
	scheduler *Scheduler
}

// NewSchedulerHandler 创建定时任务处理器
func NewSchedulerHandler(scheduler *Scheduler) *SchedulerHandler {
	return &SchedulerHandler{
		scheduler: scheduler,
	}
}

// AddTaskRequest 添加任务请求
type AddTaskRequest struct {
	Name        string `json:"name" binding:"required"`
	CronSpec    string `json:"cron_spec" binding:"required"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Name        string `json:"name" binding:"required"`
	CronSpec    string `json:"cron_spec" binding:"required"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// AddTask 添加定时任务
func (h *SchedulerHandler) AddTask(c *gin.Context) {
	var req AddTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request "+err.Error())
		return
	}

	// 注意：这里需要注册具体的任务函数
	// 实际使用时，应该通过任务注册表来获取任务函数
	response.InternalError(c, "task function not registered")
}

// RemoveTask 移除定时任务
func (h *SchedulerHandler) RemoveTask(c *gin.Context) {
	taskName := c.Param("name")
	if taskName == "" {
		response.BadRequest(c, "task name is required")
		return
	}

	if err := h.scheduler.RemoveTask(taskName); err != nil {
		logger.SchedulerL.Error("Failed to remove task",
			zap.String("task_name", taskName),
			zap.Error(err),
		)

		response.InternalError(c, fmt.Sprintf("Failed to remove task: %v", err))
		return
	}

	logger.SchedulerL.Info("Task removed via API", zap.String("task_name", taskName))

	response.Success(c, nil)
}

// GetTaskStatus 获取任务状态
func (h *SchedulerHandler) GetTaskStatus(c *gin.Context) {
	taskName := c.Param("name")
	if taskName == "" {
		response.BadRequest(c, "task name is required")
		return
	}

	status, err := h.scheduler.GetStatus(taskName)
	if err != nil {
		response.InternalError(c, "get task status failed")
		return
	}

	response.Success(c, status)
}

// ListTasks 列出所有任务
func (h *SchedulerHandler) ListTasks(c *gin.Context) {
	tasks := h.scheduler.ListTasks()
	response.Success(c, tasks)
}

// EnableTask 启用任务
func (h *SchedulerHandler) EnableTask(c *gin.Context) {
	taskName := c.Param("name")
	if taskName == "" {
		response.BadRequest(c, "task name is required")
		return
	}

	if err := h.scheduler.EnableTask(taskName); err != nil {
		logger.SchedulerL.Error("Failed to enable task",
			zap.String("task_name", taskName),
			zap.Error(err),
		)
		response.InternalError(c, "enable task failed")
		return
	}

	response.Success(c, nil)
}

// DisableTask 禁用任务
func (h *SchedulerHandler) DisableTask(c *gin.Context) {
	taskName := c.Param("name")
	if taskName == "" {
		response.Error(c, http.StatusBadRequest, "task name is required")
		return
	}

	if err := h.scheduler.DisableTask(taskName); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
