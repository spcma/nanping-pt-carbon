package domain

import (
	"app/internal/shared/crypto"
	"app/internal/shared/entity"
	idgen "app/internal/shared/idgen"
	"app/internal/shared/timeutil"
)

// UserStatus user status
type UserStatus string

const (
	UserStatusNormal   UserStatus = "normal"   // normal
	UserStatusFrozen   UserStatus = "frozen"   // frozen
	UserStatusCanceled UserStatus = "canceled" // canceled
)

// SysUser system user aggregate root
type SysUser struct {
	entity.BaseEntity
	Username    string     `json:"username" gorm:"column:username"`
	Nickname    string     `json:"nickname" gorm:"column:nickname"`
	Password    string     `json:"password" gorm:"column:password"`
	Salt        string     `json:"salt" gorm:"column:salt"`
	Status      UserStatus `json:"status" gorm:"column:status"`
	Phone       string     `json:"phone" gorm:"column:phone"`
	Avatar      string     `json:"avatar" gorm:"column:avatar"`
	Email       string     `json:"email" gorm:"column:email"`
	Description string     `json:"description" gorm:"column:description"`
	Type        string     `json:"type" gorm:"column:type"`
	ParentId    int64      `json:"parent_id" gorm:"column:parent_id"`
}

// TableName table name
func (*SysUser) TableName() string {
	return "sys_user"
}

// NewSysUser creates a new user
func NewSysUser(username, nickname, password, salt string, createUser int64) (*SysUser, error) {
	// 如果没有提供盐值，生成一个新的
	if salt == "" {
		var err error
		salt, err = crypto.GenerateSalt()
		if err != nil {
			return nil, err
		}
	}

	// 加密密码
	encryptedPassword := crypto.HashPassword(password, salt)

	user := &SysUser{
		BaseEntity: entity.BaseEntity{
			Id:         idgen.NumID(),
			CreateBy:   createUser,
			CreateTime: timeutil.Now(),
		},
		Username: username,
		Nickname: nickname,
		Password: encryptedPassword,
		Salt:     salt,
		Status:   UserStatusNormal,
	}
	return user, nil
}

// UpdateUserCommand 用户更新命令
type UpdateUserCommand struct {
	Nickname    *string `json:"nickname" validate:"omitempty,max=100"`
	Phone       *string `json:"phone" validate:"omitempty,max=20"`
	Email       *string `json:"email" validate:"omitempty,email,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	Avatar      *string `json:"avatar" validate:"omitempty,url,max=500"`
}

// UpdateInfo 更新用户信息（部分更新）
func (u *SysUser) UpdateInfo(cmd UpdateUserCommand, userID int64) error {
	if cmd.Nickname != nil {
		u.Nickname = *cmd.Nickname
	}
	if cmd.Phone != nil {
		u.Phone = *cmd.Phone
	}
	if cmd.Email != nil {
		u.Email = *cmd.Email
	}
	if cmd.Description != nil {
		u.Description = *cmd.Description
	}
	if cmd.Avatar != nil {
		u.Avatar = *cmd.Avatar
	}
	u.UpdateBy = userID
	u.UpdateTime = timeutil.Now()
	return nil
}

// ChangeStatus changes user status
func (u *SysUser) ChangeStatus(status UserStatus, userID int64) error {
	u.Status = status
	u.UpdateBy = userID
	u.UpdateTime = timeutil.Now()
	return nil
}

// ChangePassword changes password
func (u *SysUser) ChangePassword(password string, userID int64) error {
	u.Password = password
	u.UpdateBy = userID
	u.UpdateTime = timeutil.Now()
	return nil
}

// Delete 逻辑删除用户
func (u *SysUser) Delete(userID int64) error {
	u.DeleteBy = userID
	u.DeleteTime = timeutil.Now()

	return nil
}

// StringPtr creates a string pointer from string value
// This is a helper function for UpdateUserCommand
func StringPtr(s string) *string {
	return &s
}

// SysUserPageQuery system user page query object
type SysUserPageQuery struct {
	PageNum   int64  `json:"pageNum" binding:"required,min=1"`
	PageSize  int64  `json:"pageSize" binding:"required,min=1,max=100"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	UserType  string `json:"userType"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder" binding:"oneof=asc desc"` // "asc" or "desc"
}
