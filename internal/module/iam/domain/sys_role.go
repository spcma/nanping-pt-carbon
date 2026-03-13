package domain

import (
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
)

// RoleStatus role status
type RoleStatus string

const (
	RoleStatusNormal   RoleStatus = "normal"   // normal
	RoleStatusFrozen   RoleStatus = "frozen"   // frozen
	RoleStatusCanceled RoleStatus = "canceled" // canceled
)

// SysRole role aggregate root
type SysRole struct {
	entity.BaseEntity
	Name        string     `json:"name" gorm:"column:name"`
	Code        string     `json:"code" gorm:"column:code"`
	Status      RoleStatus `json:"status" gorm:"column:status"`
	Description string     `json:"description" gorm:"column:description"`
}

// TableName table name
func (*SysRole) TableName() string {
	return "sys_role"
}

// NewSysRole creates a new role
func NewSysRole(name, code, description string, createUser int64) (*SysRole, error) {
	role := &SysRole{
		BaseEntity: entity.BaseEntity{
			CreateBy:   createUser,
			CreateTime: timeutil.New(),
		},
		Name:        name,
		Code:        code,
		Description: description,
		Status:      RoleStatusNormal,
	}
	return role, nil
}

// UpdateInfo updates role info
func (r *SysRole) UpdateInfo(name, description string, userID int64) error {
	r.Name = name
	r.Description = description
	r.UpdateBy = userID
	r.UpdateTime = timeutil.New()
	return nil
}

// ChangeStatus changes role status
func (r *SysRole) ChangeStatus(status RoleStatus, userID int64) error {
	r.Status = status
	r.UpdateBy = userID
	r.UpdateTime = timeutil.New()
	return nil
}
