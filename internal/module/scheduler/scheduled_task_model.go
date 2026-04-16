package scheduler

import (
	"app/internal/shared/entity"
)

// ScheduledTask 定时任务配置实体
type ScheduledTask struct {
	entity.BaseEntity
	Name        string `gorm:"column:name;uniqueIndex;not null" json:"name"`
	CronSpec    string `gorm:"column:cron_spec;not null" json:"cronSpec"`
	Description string `gorm:"column:description" json:"description"`
	Enabled     bool   `gorm:"column:enabled;default:true" json:"enabled"`
	TaskType    string `gorm:"column:task_type;not null" json:"taskType"` // 任务类型标识,用于映射到具体的TaskFunc
}

// TableName 表名
func (*ScheduledTask) TableName() string {
	return "sys_scheduled_task"
}
