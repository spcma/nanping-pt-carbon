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
	Type        string     `json:"type" gorm:"column:type"` // 用户类型 1 系统用户 2 项目子用户
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

// UpdateInfo updates user info
func (u *SysUser) UpdateInfo(nickname, phone, email, description, avatar string, userID int64) error {
	u.Nickname = nickname
	u.Phone = phone
	u.Email = email
	u.Description = description
	u.Avatar = avatar
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
