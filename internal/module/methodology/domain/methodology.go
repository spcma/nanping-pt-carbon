package domain

import (
	"app/internal/shared/entity"
	idgen "app/internal/shared/idgen"
	"app/internal/shared/timeutil"
	"errors"
)

// MethodologyStatus 方法学状态（值对象）
type MethodologyStatus string

const (
	MethodologyStatusActive   MethodologyStatus = "1" // 启用
	MethodologyStatusInactive MethodologyStatus = "0" // 禁用
)

// IsValid 验证状态是否有效
func (s MethodologyStatus) IsValid() bool {
	switch s {
	case MethodologyStatusActive, MethodologyStatusInactive:
		return true
	default:
		return false
	}
}

// Methodology 方法学聚合根
type Methodology struct {
	entity.BaseEntity
	Name        string            `json:"name" gorm:"column:name"`
	Code        string            `json:"code" gorm:"column:code"`
	Status      MethodologyStatus `json:"status" gorm:"column:status"`
	Description string            `json:"description" gorm:"column:description"`
	Icon        string            `json:"icon" gorm:"column:icon"`
	StartDate   timeutil.Time     `json:"startDate" gorm:"column:start_date"`
	EndDate     timeutil.Time     `json:"endDate" gorm:"column:end_date"`
}

// TableName 表名
func (*Methodology) TableName() string {
	return "methodology"
}

// NewMethodology 创建新方法学（工厂方法）
func NewMethodology(name, code, icon, description string, createUser int64, startDate, endDate timeutil.Time) (*Methodology, error) {
	// 领域规则验证
	if name == "" {
		return nil, errors.New("方法学名称不能为空")
	}
	if code == "" {
		return nil, errors.New("方法学编码不能为空")
	}

	methodology := &Methodology{
		BaseEntity: entity.BaseEntity{
			Id:         idgen.NumID(),
			CreateBy:   createUser,
			CreateTime: timeutil.Now(),
		},
		Name:        name,
		Code:        code,
		Status:      MethodologyStatusActive, // 默认启用
		Description: description,
		Icon:        icon,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	return methodology, nil
}

// UpdateInfo 更新方法学信息（领域行为）
func (m *Methodology) UpdateInfo(name *string, description *string, userID int64) error {
	// 只有当指针非空时才更新对应字段
	if name != nil {
		if *name == "" {
			return errors.New("方法学名称不能为空")
		}
		m.Name = *name
	}
	if description != nil {
		m.Description = *description
	}

	m.UpdateBy = userID
	m.UpdateTime = timeutil.Now()
	return nil
}

// ChangeStatus 变更方法学状态（领域行为）
func (m *Methodology) ChangeStatus(status MethodologyStatus, userID int64) error {
	// 验证状态有效性
	if !status.IsValid() {
		return errors.New("无效的方法学状态")
	}

	m.Status = status
	m.UpdateBy = userID
	m.UpdateTime = timeutil.Now()
	return nil
}

// Delete 逻辑删除方法学（领域行为）
func (m *Methodology) Delete(userID int64) error {
	m.DeleteBy = userID
	m.DeleteTime = timeutil.Now()
	return nil
}

// Activate 启用方法学
func (m *Methodology) Activate(userID int64) error {
	return m.ChangeStatus(MethodologyStatusActive, userID)
}

// Deactivate 禁用方法学
func (m *Methodology) Deactivate(userID int64) error {
	return m.ChangeStatus(MethodologyStatusInactive, userID)
}

// IsActive 检查是否启用
func (m *Methodology) IsActive() bool {
	return m.Status == MethodologyStatusActive
}
