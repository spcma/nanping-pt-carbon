package scheduler

import (
	"app/internal/shared/db"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"

	"gorm.io/gorm"
)

// ScheduledTaskRepository 定时任务配置仓储接口
type ScheduledTaskRepository interface {
	Create(ctx context.Context, task *ScheduledTask) error
	Update(ctx context.Context, task *ScheduledTask) error
	Delete(ctx context.Context, id int64) error
	FindByName(ctx context.Context, name string) (*ScheduledTask, error)
	FindAll(ctx context.Context) ([]*ScheduledTask, error)
	FindByEnabled(ctx context.Context) ([]*ScheduledTask, error)
}

// scheduledTaskRepository 定时任务配置仓储实现
type scheduledTaskRepository struct {
}

// NewScheduledTaskRepository 创建定时任务配置仓储
func NewScheduledTaskRepository() ScheduledTaskRepository {
	return &scheduledTaskRepository{}
}

// Create 创建任务配置
func (r *scheduledTaskRepository) Create(ctx context.Context, task *ScheduledTask) error {
	return db.GetDB(ctx).WithContext(ctx).Create(task).Error
}

// Update 更新任务配置
func (r *scheduledTaskRepository) Update(ctx context.Context, task *ScheduledTask) error {
	return db.GetDB(ctx).WithContext(ctx).Updates(task).Error
}

// Delete 删除任务配置(逻辑删除)
func (r *scheduledTaskRepository) Delete(ctx context.Context, id int64) error {
	updates := map[string]interface{}{
		entity.FieldDeleteBy:   1, // TODO: 从上下文获取用户ID
		entity.FieldDeleteTime: timeutil.Now(),
	}
	return db.GetDB(ctx).WithContext(ctx).Model(&ScheduledTask{}).Where("id = ?", id).Updates(updates).Error
}

// FindByName 根据名称查找任务配置
func (r *scheduledTaskRepository) FindByName(ctx context.Context, name string) (*ScheduledTask, error) {
	var task ScheduledTask
	err := db.GetDB(ctx).WithContext(ctx).
		Where("name = ? AND "+entity.FieldDeleteBy+" = 0", name).
		First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

// FindAll 查找所有任务配置
func (r *scheduledTaskRepository) FindAll(ctx context.Context) ([]*ScheduledTask, error) {
	var tasks []*ScheduledTask
	err := db.GetDB(ctx).WithContext(ctx).
		Where(entity.FieldDeleteBy + " = 0").
		Find(&tasks).Error
	return tasks, err
}

// FindByEnabled 查找所有启用的任务配置
func (r *scheduledTaskRepository) FindByEnabled(ctx context.Context) ([]*ScheduledTask, error) {
	var tasks []*ScheduledTask
	err := db.GetDB(ctx).WithContext(ctx).
		Where("enabled = ? AND "+entity.FieldDeleteBy+" = 0", true).
		Find(&tasks).Error
	return tasks, err
}
