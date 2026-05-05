package http

import (
	"app/internal/platform/security"
	"app/internal/shared/logger"
	"context"
	"fmt"

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

// MustGetCurrentUser 从上下文中获取当前用户（必须存在）
// 适用于已使用 AuthMiddleware 的路由，如果用户不存在会 panic
// 这样可以避免在每个 handler 中都写 nil 检查
func MustGetCurrentUser(c *gin.Context) *security.User {
	currentUser := GetCurrentUser(c)
	if currentUser == nil {
		// 理论上不应该发生，因为 AuthMiddleware 已经保证了用户存在
		// 如果发生了，说明路由配置有问题（应该用 AuthMiddleware 但没用）
		panic("user not found in context, please check if AuthMiddleware is applied")
	}
	return currentUser
}

// GetCurrentUserOrError 从上下文中获取当前用户，如果不存在则返回错误
// 这是防御性编程的最佳实践，既安全又优雅
func GetCurrentUserOrError(c *gin.Context) (*security.User, error) {
	currentUser := GetCurrentUser(c)
	if currentUser == nil {
		return nil, fmt.Errorf("user not found in context")
	}
	return currentUser, nil
}

// CheckPermission 检查权限
func CheckPermission(permissionChecker *security.DefaultPermissionChecker, c *gin.Context, action string, resource interface{}) bool {
	currentUser := GetCurrentUser(c)
	return permissionChecker.Can(c.Request.Context(), currentUser, action, resource)
}

// RequirePermission 需要权限 (如果没有权限则返回错误)
func RequirePermission(permissionChecker *security.DefaultPermissionChecker, c *gin.Context, action string, resource interface{}) bool {
	if !CheckPermission(permissionChecker, c, action, resource) {
		// TODO: 返回禁止访问错误
		return false
	}
	return true
}

// GetTraceID 获取请求的追踪 ID
func GetTraceID(c *gin.Context) string {
	return logger.GetTraceID(c)
}

// GetTraceIDFromContext 从 Go context 获取追踪 ID
func GetTraceIDFromContext(ctx context.Context) string {
	return logger.TraceID(ctx)
}
