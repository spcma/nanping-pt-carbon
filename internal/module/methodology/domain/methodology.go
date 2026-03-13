package domain

import (
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
)

// MethodologyStatus 方法学状态
type MethodologyStatus string

const (
	MethodologyStatusActive   MethodologyStatus = "1" // 启用
	MethodologyStatusInactive MethodologyStatus = "0" // 禁用
)

// Methodology 方法学聚合根
type Methodology struct {
	entity.BaseEntity
	Name        string            `json:"name" gorm:"column:name"`
	Code        string            `json:"code" gorm:"column:code"`
	Status      MethodologyStatus `json:"status" gorm:"column:status"`
	Description string            `json:"description" gorm:"column:description"`
}

// TableName 表名
func (*Methodology) TableName() string {
	return "methodology"
}

// NewMethodology 创建新方法学
func NewMethodology(name, code, description string, createUser int64) (*Methodology, error) {
	methodology := &Methodology{
		BaseEntity: entity.BaseEntity{
			CreateBy:   createUser,
			CreateTime: timeutil.New(),
		},
		Name:        name,
		Code:        code,
		Description: description,
		Status:      MethodologyStatusActive,
	}
	return methodology, nil
}

// UpdateInfo 更新方法学信息
func (m *Methodology) UpdateInfo(name, description string, userID int64) error {
	m.Name = name
	m.Description = description
	m.UpdateBy = userID
	m.UpdateTime = timeutil.New()
	return nil
}

// ChangeStatus 变更方法学状态
func (m *Methodology) ChangeStatus(status MethodologyStatus, userID int64) error {
	m.Status = status
	m.UpdateBy = userID
	m.UpdateTime = timeutil.New()
	return nil
}
