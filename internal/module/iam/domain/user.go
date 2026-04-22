package domain

import (
	"app/internal/shared/crypto"
	"app/internal/shared/entity"
	idgen "app/internal/shared/idgen"
	"app/internal/shared/timeutil"
	"fmt"
)

// UserStatus user status
type UserStatus string

const (
	UserStatusNormal   UserStatus = "1" // 可用
	UserStatusFrozen   UserStatus = "2" // 冻结
	UserStatusCanceled UserStatus = "9" // 注销
)

// User system user aggregate root
type User struct {
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
}

// TableName table name
func (*User) TableName() string {
	return "sys_user"
}

// NewUser creates a new user
//
//	统一使用构造方法创建用户，方便统一进行参数校验
func NewUser(username, nickname, password, salt string, createUser int64) (*User, error) {
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	if nickname == "" {
		return nil, fmt.Errorf("nickname cannot be empty")
	}

	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	// 如果没有提供盐值，生成一个新的
	if salt == "" {
		var err error
		salt, err = crypto.GenerateSalt()
		if err != nil {
			return nil, err
		}
	}

	if createUser == 0 {
		createUser = 1
	}

	// 加密密码
	encryptedPassword := crypto.HashPassword(password, salt)

	user := &User{
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
	Nickname    *string `json:"nickname"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email"`
	Description *string `json:"description"`
	Avatar      *string `json:"avatar"`
	Status      *string `json:"status"`
}

// UpdateInfo 更新用户信息（部分更新）
func (u *User) UpdateInfo(cmd UpdateUserCommand, userID int64) error {
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
	if cmd.Status != nil {
		u.Status = UserStatus(*cmd.Status)
	}

	u.UpdateBy = userID
	u.UpdateTime = timeutil.Now()
	return nil
}

// ChangeStatus changes user status
func (u *User) ChangeStatus(status UserStatus, userID int64) error {
	u.Status = status
	u.UpdateBy = userID
	u.UpdateTime = timeutil.Now()
	return nil
}

// ChangePassword changes password
func (u *User) ChangePassword(password string, userID int64) error {
	u.Password = password
	u.UpdateBy = userID
	u.UpdateTime = timeutil.Now()
	return nil
}

// Delete 逻辑删除用户
func (u *User) Delete(userID int64) error {
	u.DeleteBy = userID
	u.DeleteTime = timeutil.Now()

	return nil
}

// UsersPageQuery system user page query object
type UsersPageQuery struct {
	entity.PaginationQuery
	Username  string `json:"username" form:"username"`
	Nickname  string `json:"nickname" form:"nickname"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	UserType  string `json:"userType"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"` // "asc" or "desc"
}

type UserQuery struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}
