package application

import (
	"app/internal/module/iam/domain"
	"app/internal/shared/crypto"
	"app/internal/shared/entity"
	"app/internal/shared/token"
	"context"
)

// UsersService system user application service
type UsersService struct {
	repo UserRepo
}

// NewUsersService creates system user application service
func NewUsersService(repo UserRepo) *UsersService {
	return &UsersService{repo: repo}
}

// CreateUserCommand create system user command
type CreateUserCommand struct {
	Username *string `json:"username"`
	Nickname *string `json:"nickname"`
	Password *string `json:"password"`
	UserID   int64   `json:"userId"`
}

// Create creates a system user
func (u *UsersService) Create(ctx context.Context, cmd CreateUserCommand) (int64, error) {
	user, err := domain.NewUser(*cmd.Username, *cmd.Nickname, *cmd.Password, "", cmd.UserID)
	if err != nil {
		return 0, err
	}
	err = u.repo.Create(ctx, user)
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

// UpdateUserCommand update system user command
type UpdateUserCommand struct {
	ID          int64   `json:"id"`
	Nickname    *string `json:"nickname"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email"`
	Description *string `json:"description"`
	Avatar      *string `json:"avatar"`
	Status      *string `json:"status"`
	UserID      int64   `json:"userId"`
}

// Update updates a system user
func (u *UsersService) Update(ctx context.Context, updateParam UpdateUserCommand) error {
	user, err := u.repo.FindByID(ctx, updateParam.ID)
	if err != nil {
		return err
	}

	if user == nil {
		return domain.ErrUserNotFound
	}

	// 将 application 层的 Command 转换为 domain 层的 Command
	domainCmd := domain.UpdateUserCommand{
		Nickname:    updateParam.Nickname,
		Phone:       updateParam.Phone,
		Email:       updateParam.Email,
		Description: updateParam.Description,
		Avatar:      updateParam.Avatar,
		Status:      updateParam.Status,
	}

	err = user.UpdateInfo(domainCmd, updateParam.UserID)
	if err != nil {
		return err
	}

	return u.repo.Update(ctx, user)
}

type DeleteUserCommand struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"userId"`
}

// Delete deletes a system user (logical delete)
func (u *UsersService) Delete(ctx context.Context, param *DeleteUserCommand) error {
	user, err := u.repo.FindByID(ctx, param.ID)
	if err != nil {
		return err
	}

	if user == nil {
		return domain.ErrUserNotFound
	}

	err = user.Delete(param.UserID)
	if err != nil {
		return err
	}

	// 持久化
	return u.repo.Update(ctx, user)
}

// GetByID gets system user by ID
func (u *UsersService) GetByID(ctx context.Context, id int64) (*domain.Users, error) {
	return u.repo.FindByID(ctx, id)
}

// GetByUsername gets system user by username
func (u *UsersService) GetByUsername(ctx context.Context, username string) (*domain.Users, error) {
	return u.repo.FindByUsername(ctx, username)
}

func (u *UsersService) GetByQuery(ctx context.Context, query *domain.UserQuery) (*domain.Users, error) {
	user, err := u.repo.FindByQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UsersService) GetList(ctx context.Context) ([]*domain.Users, error) {
	list, err := u.repo.FindList(ctx)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetPage queries system users with pagination
func (u *UsersService) GetPage(ctx context.Context, query *domain.UsersPageQuery) (*entity.PaginationResult[domain.Users], error) {
	result, err := u.repo.FindPage(ctx, query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ChangePasswordCommand change password command
type ChangePasswordCommand struct {
	Id          int64  `json:"id"` // userid
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
	Force       bool   `json:"force"`
}

// ChangePassword changes password
func (u *UsersService) ChangePassword(ctx context.Context, cmd *ChangePasswordCommand) error {
	user, err := u.repo.FindByID(ctx, cmd.Id)
	if err != nil {
		return err
	}

	if user.Id == 0 {
		return domain.ErrUserNotFound
	}

	//	非强制修改密码，需要验证旧密码
	if cmd.Force {
		//	验证旧密码
		validPassword := crypto.HashPassword(cmd.Password, user.Salt)
		if validPassword == user.Password {
			return domain.ErrNewPasswordSameAsOldPassword
		}
	}

	//	加密新密码
	newPassword := crypto.HashPassword(cmd.NewPassword, user.Salt)

	err = user.ChangePassword(newPassword, user.Id)
	if err != nil {
		return err
	}

	err = u.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// ChangeUserStatus changes user status
func (u *UsersService) ChangeUserStatus(ctx context.Context, id int64, status domain.UserStatus, userID int64) error {
	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return user.ChangeStatus(status, userID)
}

// RegisterCommand 注册命令
type RegisterCommand struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

// Register 用户注册
func (u *UsersService) Register(ctx context.Context, cmd RegisterCommand) (int64, error) {
	// 检查用户名是否已存在
	existingUser, err := u.repo.FindByUsername(ctx, cmd.Username)
	if err == nil && existingUser != nil && !existingUser.IsDeleted() {
		return 0, domain.ErrUserAlreadyExists
	}

	// 创建用户（密码加密在 domain 层处理）
	return u.Create(ctx, CreateUserCommand{
		Username: &cmd.Username,
		Nickname: &cmd.Nickname,
		Password: &cmd.Password,
		UserID:   0, // 注册用户没有创建者 ID
	})
}

// LoginCommand 登录命令
type LoginCommand struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string         `json:"token"`
	UserInfo *LoginUserInfo `json:"userInfo"`
}

// LoginUserInfo 登录用户信息
type LoginUserInfo struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	Roles    []string `json:"roles"`
}

// Login 用户登录
func (u *UsersService) Login(ctx context.Context, cmd LoginCommand, jwtManager token.Manager) (*LoginResponse, error) {
	// 查找用户
	user, err := u.repo.FindByUsername(ctx, cmd.Username)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// 检查用户是否存在
	if user.Id == 0 {
		return nil, domain.ErrUserNotFound
	}

	// 检查用户状态
	if user.Status != domain.UserStatusNormal {
		return nil, domain.ErrUserFrozen
	}

	// 验证密码
	if !crypto.VerifyPassword(cmd.Password, user.Password, user.Salt) {
		return nil, domain.ErrInvalidPassword
	}

	// 生成 JWT Token（默认角色为 USER）
	tokenString, err := jwtManager.GenerateToken(user.Id, user.Username, []string{"USER"})
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: tokenString,
		UserInfo: &LoginUserInfo{
			ID:       user.Id,
			Username: user.Username,
			Nickname: user.Nickname,
			Roles:    []string{"USER"},
		},
	}, nil
}
