package security

import "context"

// User 用户信息
type User struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

// PermissionChecker 权限检查器接口
type PermissionChecker interface {
	// Can 检查用户是否有权限执行某个操作
	Can(ctx context.Context, user *User, action string, resource interface{}) bool

	// HasRole 检查用户是否有某个角色
	HasRole(ctx context.Context, user *User, role string) bool

	// HasAnyRole 检查用户是否有任何角色
	HasAnyRole(ctx context.Context, user *User, roles []string) bool
}

// DefaultPermissionChecker 默认权限检查器
type DefaultPermissionChecker struct {
	// TODO: 添加权限规则配置
}

// NewDefaultPermissionChecker 创建默认权限检查器
func NewDefaultPermissionChecker() *DefaultPermissionChecker {
	return &DefaultPermissionChecker{}
}

// Can 检查用户是否有权限执行某个操作
func (p *DefaultPermissionChecker) Can(ctx context.Context, user *User, action string, resource interface{}) bool {
	// TODO: 实现权限检查逻辑
	// 示例：管理员可以执行所有操作
	return p.HasRole(ctx, user, "admin")
}

// HasRole 检查用户是否有某个角色
func (p *DefaultPermissionChecker) HasRole(ctx context.Context, user *User, role string) bool {
	if user == nil {
		return false
	}
	for _, r := range user.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole 检查用户是否有任何角色
func (p *DefaultPermissionChecker) HasAnyRole(ctx context.Context, user *User, roles []string) bool {
	for _, role := range roles {
		if p.HasRole(ctx, user, role) {
			return true
		}
	}
	return false
}
