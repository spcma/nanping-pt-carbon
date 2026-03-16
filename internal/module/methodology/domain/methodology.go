package domain

import (
	"app/internal/shared/entity"
	idgen "app/internal/shared/idgen"
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
			Id:         idgen.NumID(),
			CreateBy:   createUser,
			CreateTime: timeutil.Now(),
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
	m.UpdateTime = timeutil.Now()
	return nil
}

// ChangeStatus 变更方法学状态
func (m *Methodology) ChangeStatus(status MethodologyStatus, userID int64) error {
	m.Status = status
	m.UpdateBy = userID
	m.UpdateTime = timeutil.Now()
	return nil
}

// Delete 逻辑删除方法学
func (m *Methodology) Delete(userID int64) error {
	m.DeleteBy = userID
	m.DeleteTime = timeutil.Now()
	return nil
}

// MethodologyPageQuery 方法学分页查询对象
type MethodologyPageQuery struct {
	entity.Pagination
	Name      string `json:"name"`   // 方法学名模糊匹配
	Code      string `json:"code"`   // 方法学编码精确匹配
	Status    string `json:"status"` // 状态
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder" binding:"oneof=asc desc"` // "asc" or "desc"
}
