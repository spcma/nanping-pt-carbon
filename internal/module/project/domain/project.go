package domain

import (
	"app/internal/shared/entity"
	idgen "app/internal/shared/idgen"
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
	Icon        string        `json:"icon" gorm:"column:icon"`
	StartDate   timeutil.Time `json:"startDate" gorm:"column:start_date"`
	EndDate     timeutil.Time `json:"endDate" gorm:"column:end_date"`
}

// TableName 表名
func (*Project) TableName() string {
	return "project"
}

// NewProject 创建新项目
func NewProject(name, code, icon, description string, createUser int64, startDate, endDate timeutil.Time) (*Project, error) {
	project := &Project{
		BaseEntity: entity.BaseEntity{
			Id:         idgen.NumID(),
			CreateBy:   createUser,
			CreateTime: timeutil.Now(),
		},
		Name:        name,
		Code:        code,
		Status:      ProjectStatusActive,
		Description: description,
		Icon:        icon,
		StartDate:   startDate,
		EndDate:     endDate,
	}
	return project, nil
}

// UpdateInfo 更新项目信息
func (p *Project) UpdateInfo(name, description string, userID int64) error {
	p.Name = name
	p.Description = description
	p.UpdateBy = userID
	p.UpdateTime = timeutil.Now()
	return nil
}

// ChangeStatus 变更项目状态
func (p *Project) ChangeStatus(status ProjectStatus, userID int64) error {
	p.Status = status
	p.UpdateBy = userID
	p.UpdateTime = timeutil.Now()
	return nil
}

// Delete 逻辑删除项目
func (p *Project) Delete(userID int64) error {
	p.DeleteBy = userID
	p.DeleteTime = timeutil.Now()
	return nil
}

// ProjectPageQuery system user page query object
type ProjectPageQuery struct {
	entity.PaginationQuery
	Name      string `json:"name"`
	Code      string `json:"code"`
	Status    string `json:"status"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"` // "asc" or "desc"
}
