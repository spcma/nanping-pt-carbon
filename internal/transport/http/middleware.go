package http

import (
	"app/internal/platform/http/response"
	"app/internal/platform/security"
	"app/internal/shared/logger"
	"app/internal/shared/token"
	"context"

	"github.com/gin-gonic/gin"
)

// contextKey 用于上下文字符串 key
type contextKey string

const (
	// UserContextKey 用户信息在上下文中的 key
	UserContextKey contextKey = "user"
)

// LogMiddleware 日志中间件（使用新的 logger middleware）
func LogMiddleware() gin.HandlerFunc {
	return logger.LoggingMiddleware()
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Recovery 恢复中间件（使用新的 logger middleware）
func Recovery() gin.HandlerFunc {
	return logger.RecoveryMiddleware()
}

// parseToken 解析 JWT token 并返回 claims
func parseToken(tokenString string, jwtManager token.Manager) (*token.Claims, error) {
	return jwtManager.ValidateToken(tokenString)
}

// AuthMiddleware 认证中间件（强制要求 token）
func AuthMiddleware(jwtManager token.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			response.Unauthorized(c, "未授权")
			c.Abort()
			return
		}

		claims, err := parseToken(authToken, jwtManager)
		if err != nil {
			logger.Error("auth", "Token validation failed: "+err.Error())
			response.Unauthorized(c, "令牌无效")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		user := &security.User{
			ID:       claims.UserID,
			Username: claims.Username,
			Roles:    claims.GetRoles(),
		}

		// 存储到 Gin 上下文
		c.Set(string(UserContextKey), user)

		// 同时存储到 Go context
		ctx := context.WithValue(c.Request.Context(), UserContextKey, user)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（有 token 则解析，没有也放行）
// 适用于：公开接口但希望获取当前用户信息的场景
func OptionalAuthMiddleware(jwtManager token.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			// 没有 token，继续执行（不拦截）
			c.Next()
			return
		}

		claims, err := parseToken(authToken, jwtManager)
		if err != nil {
			// token 无效，忽略错误继续执行（不拦截）
			logger.Debug("auth", "Optional auth - invalid token ignored")
			c.Next()
			return
		}

		// token 有效，将用户信息存入上下文
		user := &security.User{
			ID:       claims.UserID,
			Username: claims.Username,
			Roles:    claims.GetRoles(),
		}

		c.Set(string(UserContextKey), user)
		ctx := context.WithValue(c.Request.Context(), UserContextKey, user)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
