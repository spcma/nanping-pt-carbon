package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// AuthType 认证类型
type AuthType int

const (
	// AuthTypeNone 不需要认证（完全公开）
	AuthTypeNone AuthType = iota

	// AuthTypeOptional 可选认证（有 token 返回增强信息，无 token 返回基础信息）
	AuthTypeOptional

	// AuthTypeRequired 强制认证（必须有有效 token）
	AuthTypeRequired
)

// RouteHandler 定义单个路由处理器
type RouteHandler struct {
	Method  string          // HTTP 方法：GET, POST, PUT, DELETE
	Path    string          // 路由路径（相对于组路径）
	Handler gin.HandlerFunc // 处理器函数
}

// RouteGroupConfig 路由组配置（保留用于向后兼容）
type RouteGroupConfig struct {
	Prefix   string         // 路由组前缀，如 "/bus", "/sys/user"
	Handlers []RouteHandler // 该组下的所有路由处理器
	AuthType AuthType       // 认证类型
}

// RouteRegistry 路由注册器接口（新标准）
// 各模块应实现此接口来注册自己的路由
type RouteRegistry interface {
	// RegisterRoutes 注册模块路由
	// group: 模块级别的 RouterGroup（已包含 /api 前缀）
	// middlewares: 根据认证类型提供的中间件映射
	RegisterRoutes(group *gin.RouterGroup, middlewares map[AuthType]gin.HandlerFunc)
}

// ValidateRouteConfigs 验证路由配置的有效性（保留用于向后兼容）
func ValidateRouteConfigs(configs []RouteGroupConfig) error {
	seen := make(map[string]string) // method:path -> route

	for i, config := range configs {
		for _, handler := range config.Handlers {
			fullPath := config.Prefix + handler.Path
			key := fmt.Sprintf("%s:%s", handler.Method, fullPath)

			if _, exists := seen[key]; exists {
				return fmt.Errorf("route conflict detected: %s is registered by multiple routes", key)
			}
			seen[key] = fmt.Sprintf("config[%d]%s", i, config.Prefix)
		}
	}

	return nil
}

// GetAuthMiddleware 根据认证类型获取对应的中间件
func GetAuthMiddleware(authType AuthType, middlewares map[AuthType]gin.HandlerFunc) gin.HandlerFunc {
	if middleware, exists := middlewares[authType]; exists {
		return middleware
	}
	return nil
}
