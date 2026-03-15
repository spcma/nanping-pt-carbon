package domain

import (
	"time"
)

// SysUserRole user role association entity
type SysUserRole struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	SysUserId  int64     `json:"sysUserId" gorm:"column:sys_user_id"`
	SysRoleId  int64     `json:"sysRoleId" gorm:"column:sys_role_id"`
	CreateUser int64     `json:"createUser" gorm:"column:create_user"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`
}

// TableName table name
func (SysUserRole) TableName() string {
	return "sys_user_role"
}

// NewSysUserRole creates a new user role association
func NewSysUserRole(userId, roleId, createUser int64) (*SysUserRole, error) {
	userRole := &SysUserRole{
		SysUserId:  userId,
		SysRoleId:  roleId,
		CreateUser: createUser,
		CreateTime: time.Now(),
	}
	return userRole, nil
}
