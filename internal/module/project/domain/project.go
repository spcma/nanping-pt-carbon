package domain

import (
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
)

// ProjectStatus 项目状态
type ProjectStatus string

const (
	ProjectStatusActive    ProjectStatus = "1" // 进行中
	ProjectStatusPending   ProjectStatus = "0" // 待启动
	ProjectStatusCompleted ProjectStatus = "2" // 已完成
	ProjectStatusCancelled ProjectStatus = "3" // 已取消
)

// Project 项目聚合根
type Project struct {
	entity.BaseEntity
	Name        string        `json:"name" gorm:"column:name"`
	Code        string        `json:"code" gorm:"column:code"`
	Status      ProjectStatus `json:"status" gorm:"column:status"`
	Description string        `json:"description" gorm:"column:description"`
}

// TableName 表名
func (*Project) TableName() string {
	return "project"
}

// NewProject 创建新项目
func NewProject(name, code, description string, createUser int64) (*Project, error) {
	project := &Project{
		BaseEntity: entity.BaseEntity{
			CreateBy:   createUser,
			CreateTime: timeutil.New(),
		},
		Name:        name,
		Code:        code,
		Description: description,
		Status:      ProjectStatusActive,
	}
	return project, nil
}

// UpdateInfo 更新项目信息
func (p *Project) UpdateInfo(name, description string, userID int64) error {
	p.Name = name
	p.Description = description
	p.UpdateBy = userID
	p.UpdateTime = timeutil.New()
	return nil
}

// ChangeStatus 变更项目状态
func (p *Project) ChangeStatus(status ProjectStatus, userID int64) error {
	p.Status = status
	p.UpdateBy = userID
	p.UpdateTime = timeutil.New()
	return nil
}
