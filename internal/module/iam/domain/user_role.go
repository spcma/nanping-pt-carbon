package domain

import (
	"time"
)

// UserRole user role association entity
type UserRole struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	SysUserId  int64     `json:"sysUserId" gorm:"column:sys_user_id"`
	SysRoleId  int64     `json:"sysRoleId" gorm:"column:sys_role_id"`
	CreateUser int64     `json:"createUser" gorm:"column:create_user"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`
}

// TableName table name
func (UserRole) TableName() string {
	return "sys_user_role"
}

// NewUserRole creates a new user role association
func NewUserRole(userId, roleId, createUser int64) (*UserRole, error) {
	userRole := &UserRole{
		SysUserId:  userId,
		SysRoleId:  roleId,
		CreateUser: createUser,
		CreateTime: time.Now(),
	}
	return userRole, nil
}
