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

// CreateSysUserCommand create system user command
type CreateSysUserCommand struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	UserID   int64  `json:"userId"`
}

// Create creates a system user
func (s *UsersService) Create(ctx context.Context, cmd CreateSysUserCommand) (int64, error) {
	user, err := domain.NewUser(cmd.Username, cmd.Nickname, cmd.Password, "", cmd.UserID)
	if err != nil {
		return 0, err
	}
	err = s.repo.Create(ctx, user)
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

// UpdateSysUserParam update system user command
type UpdateSysUserParam struct {
	ID          int64   `json:"id"`
	Nickname    *string `json:"nickname"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email"`
	Description *string `json:"description"`
	Avatar      *string `json:"avatar"`
	UserID      int64   `json:"userId"`
}

// UpdateSysUser updates a system user
func (s *UsersService) UpdateSysUser(ctx context.Context, updateParam UpdateSysUserParam) error {
	user, err := s.repo.FindByID(ctx, updateParam.ID)
	if err != nil {
		return err
	}

	// 将 application 层的 Command 转换为 domain 层的 Command
	domainCmd := domain.UpdateUserCommand{
		Nickname:    updateParam.Nickname,
		Phone:       updateParam.Phone,
		Email:       updateParam.Email,
		Description: updateParam.Description,
		Avatar:      updateParam.Avatar,
	}

	err = user.UpdateInfo(domainCmd, updateParam.UserID)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, user)
}

type DeleteSysUserParam struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"userId"`
}

// DeleteSysUser deletes a system user (logical delete)
func (s *UsersService) DeleteSysUser(ctx context.Context, param *DeleteSysUserParam) error {
	user, err := s.repo.FindByID(ctx, param.ID)
	if err != nil {
		return err
	}

	// 检查是否已被删除
	if user.IsDeleted() {
		return domain.ErrUserAlreadyDeleted
	}

	// 执行领域方法
	if err := user.Delete(param.UserID); err != nil {
		return err
	}

	// 持久化
	return s.repo.Update(ctx, user)
}

type UsersQuery struct {
	Username string `json:"username"`
}

// GetSysUserByID gets system user by ID
func (s *UsersService) GetSysUserByID(ctx context.Context, id int64) (*domain.Users, error) {
	return s.repo.FindByID(ctx, id)
}

// GetSysUserByUsername gets system user by username
func (s *UsersService) GetSysUserByUsername(ctx context.Context, username string) (*domain.Users, error) {
	return s.repo.FindByUsername(ctx, username)
}

func (s *UsersService) GetList(ctx context.Context) ([]*domain.Users, error) {
	list, err := s.repo.FindList(ctx)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetSysUserPage queries system users with pagination
func (s *UsersService) GetSysUserPage(ctx context.Context, query *domain.SysUserPageQuery) (*entity.PaginationResult[*domain.Users], error) {
	result, err := s.repo.FindPage(ctx, query)
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
}

// ChangePassword changes password
func (s *UsersService) ChangePassword(ctx context.Context, cmd ChangePasswordCommand) error {
	user, err := s.repo.FindByID(ctx, cmd.Id)
	if err != nil {
		return err
	}

	if user.Id == 0 {
		return domain.ErrUserNotFound
	}

	//	验证旧密码
	newPassword := crypto.HashPassword(cmd.NewPassword, user.Salt)

	if newPassword == user.Password {
		return domain.ErrNewPasswordSameAsOldPassword
	}

	return user.ChangePassword(cmd.Password, user.Id)
}

// ChangeUserStatus changes user status
func (s *UsersService) ChangeUserStatus(ctx context.Context, id int64, status domain.UserStatus, userID int64) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return user.ChangeStatus(status, userID)
}

// RegisterCommand 注册命令
type RegisterCommand struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Nickname string `json:"nickname"`
}

// Register 用户注册
func (s *UsersService) Register(ctx context.Context, cmd RegisterCommand) (int64, error) {
	// 检查用户名是否已存在
	existingUser, err := s.repo.FindByUsername(ctx, cmd.Username)
	if err == nil && existingUser != nil && !existingUser.IsDeleted() {
		return 0, domain.ErrUserAlreadyExists
	}

	// 创建用户（密码加密在 domain 层处理）
	return s.Create(ctx, CreateSysUserCommand{
		Username: cmd.Username,
		Nickname: cmd.Nickname,
		Password: cmd.Password,
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
func (s *UsersService) Login(ctx context.Context, cmd LoginCommand, jwtManager token.Manager) (*LoginResponse, error) {
	// 查找用户
	user, err := s.repo.FindByUsername(ctx, cmd.Username)
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
