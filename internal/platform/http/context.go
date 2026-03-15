package http

import (
	"app/internal/platform/security"
	"context"

	"github.com/gin-gonic/gin"
)

// contextKey 用于上下文字符串 key
type contextKey string

const (
	// UserContextKey 用户信息在上下文中的 key
	UserContextKey contextKey = "user"
)

// GetUserFromContext 从 Gin 上下文中获取用户信息
func GetUserFromContext(c *gin.Context) *security.User {
	if val, exists := c.Get(string(UserContextKey)); exists {
		if user, ok := val.(*security.User); ok {
			return user
		}
	}
	return nil
}

// GetUserFromGoContext 从 Go context 中获取用户信息
func GetUserFromGoContext(ctx context.Context) *security.User {
	if val := ctx.Value(UserContextKey); val != nil {
		if user, ok := val.(*security.User); ok {
			return user
		}
	}
	return nil
}

// Ctx 获取请求上下文
func Ctx(c *gin.Context) context.Context {
	return c.Request.Context()
}

// GetCurrentUser 从上下文中获取当前用户
func GetCurrentUser(c *gin.Context) *security.User {
	// 从 Gin 上下文中获取用户信息
	user := GetUserFromContext(c)
	if user == nil {
		// 如果未找到，尝试从 Go context 中获取
		user = GetUserFromGoContext(c.Request.Context())
	}
	return user
}

// CheckPermission 检查权限
func CheckPermission(permissionChecker *security.DefaultPermissionChecker, c *gin.Context, action string, resource interface{}) bool {
	user := GetCurrentUser(c)
	return permissionChecker.Can(c.Request.Context(), user, action, resource)
}

// RequirePermission 需要权限 (如果没有权限则返回错误)
func RequirePermission(permissionChecker *security.DefaultPermissionChecker, c *gin.Context, action string, resource interface{}) bool {
	if !CheckPermission(permissionChecker, c, action, resource) {
		// TODO: 返回禁止访问错误
		return false
	}
	return true
}
