package domain

import (
	"app/internal/shared/entity"
	idgen "app/internal/shared/idgen"
	"app/internal/shared/timeutil"
	"errors"
)

// ProjectStatus 项目状态（值对象）
type ProjectStatus string

const (
	ProjectStatusActive    ProjectStatus = "1" // 进行中
	ProjectStatusPending   ProjectStatus = "0" // 待启动
	ProjectStatusCompleted ProjectStatus = "2" // 已完成
	ProjectStatusCancelled ProjectStatus = "3" // 已取消
)

// IsValid 验证状态是否有效
func (s ProjectStatus) IsValid() bool {
	switch s {
	case ProjectStatusActive, ProjectStatusPending, ProjectStatusCompleted, ProjectStatusCancelled:
		return true
	default:
		return false
	}
}

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

// NewProject 创建新项目（工厂方法）
func NewProject(name, code, icon, description string, createUser int64, startDate, endDate timeutil.Time) (*Project, error) {
	// 领域规则验证
	if name == "" {
		return nil, errors.New("项目名称不能为空")
	}
	if code == "" {
		return nil, errors.New("项目编码不能为空")
	}
	if !startDate.Before(endDate.ToTime()) {
		return nil, errors.New("开始日期必须早于结束日期")
	}

	project := &Project{
		BaseEntity: entity.BaseEntity{
			Id:         idgen.NumID(),
			CreateBy:   createUser,
			CreateTime: timeutil.Now(),
		},
		Name:        name,
		Code:        code,
		Status:      ProjectStatusPending, // 新项目默认为待启动
		Description: description,
		Icon:        icon,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	return project, nil
}

// UpdateInfo 更新项目信息（领域行为）
func (p *Project) UpdateInfo(name *string, description *string, userID int64) error {
	// 领域规则：只有进行中的项目可以更新
	if p.Status != ProjectStatusActive && p.Status != ProjectStatusPending {
		return errors.New("只有进行中或待启动的项目可以更新")
	}

	// 只有当指针非空时才更新对应字段
	if name != nil {
		if *name == "" {
			return errors.New("项目名称不能为空")
		}
		p.Name = *name
	}
	if description != nil {
		p.Description = *description
	}

	p.UpdateBy = userID
	p.UpdateTime = timeutil.Now()
	return nil
}

// ChangeStatus 变更项目状态（领域行为）
func (p *Project) ChangeStatus(status ProjectStatus, userID int64) error {
	// 领域规则：验证状态转换的合法性
	if !status.IsValid() {
		return errors.New("无效的项目状态")
	}

	// 状态转换规则
	if p.Status == ProjectStatusCompleted || p.Status == ProjectStatusCancelled {
		return errors.New("已完成或已取消的项目不能变更状态")
	}

	// 特定状态转换规则
	if p.Status == ProjectStatusPending && status == ProjectStatusCompleted {
		return errors.New("待启动的项目不能直接变更为已完成")
	}

	p.Status = status
	p.UpdateBy = userID
	p.UpdateTime = timeutil.Now()
	return nil
}

// Delete 逻辑删除项目（领域行为）
func (p *Project) Delete(userID int64) error {
	// 领域规则：只有待启动的项目可以删除
	if p.Status != ProjectStatusPending {
		return errors.New("只有待启动的项目可以删除")
	}

	p.DeleteBy = userID
	p.DeleteTime = timeutil.Now()
	return nil
}

// Activate 激活项目（启动项目）
func (p *Project) Activate(userID int64) error {
	if p.Status != ProjectStatusPending {
		return errors.New("只有待启动的项目可以激活")
	}

	p.Status = ProjectStatusActive
	p.UpdateBy = userID
	p.UpdateTime = timeutil.Now()
	return nil
}

// Complete 完成项目
func (p *Project) Complete(userID int64) error {
	if p.Status != ProjectStatusActive {
		return errors.New("只有进行中的项目可以完成")
	}

	p.Status = ProjectStatusCompleted
	p.UpdateBy = userID
	p.UpdateTime = timeutil.Now()
	return nil
}

// Cancel 取消项目
func (p *Project) Cancel(userID int64) error {
	if p.Status == ProjectStatusCompleted || p.Status == ProjectStatusCancelled {
		return errors.New("已完成或已取消的项目不能再次取消")
	}

	p.Status = ProjectStatusCancelled
	p.UpdateBy = userID
	p.UpdateTime = timeutil.Now()
	return nil
}

// IsOverdue 检查项目是否逾期
func (p *Project) IsOverdue() bool {
	if p.Status == ProjectStatusCompleted || p.Status == ProjectStatusCancelled {
		return false
	}
	return timeutil.Now().After(p.EndDate.ToTime())
}
