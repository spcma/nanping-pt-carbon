package http

import (
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

// RouteGroupConfig 路由组配置
type RouteGroupConfig struct {
	Prefix   string         // 路由组前缀，如 "/bus", "/sys/user"
	Handlers []RouteHandler // 该组下的所有路由处理器
	AuthType AuthType       // 认证类型
}

type RegistryRoute interface {
	RegisterRoutes() []RouteGroupConfig
}

// ValidateRouteConfigs 验证路由配置的有效性
func ValidateRouteConfigs(configs []RouteGroupConfig) error {
	return nil
}
