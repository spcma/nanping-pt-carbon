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
	prefixSet := make(map[string]bool)

	for i, config := range configs {
		// 检查空的前缀
		if config.Prefix == "" {
			return fmt.Errorf("route config[%d]: empty prefix", i)
		}

		// 检查重复的前缀
		if prefixSet[config.Prefix] {
			return fmt.Errorf("route config[%d]: duplicate prefix '%s'", i, config.Prefix)
		}
		prefixSet[config.Prefix] = true

		// 检查空的 Handlers
		if len(config.Handlers) == 0 {
			return fmt.Errorf("route config[%d]: no handlers for prefix '%s'", i, config.Prefix)
		}

		// 检查每个 Handler 的有效性
		for j, handler := range config.Handlers {
			if handler.Method == "" {
				return fmt.Errorf("route config[%d].handler[%d]: empty HTTP method for prefix '%s'", i, j, config.Prefix)
			}

			if handler.Handler == nil {
				return fmt.Errorf("route config[%d].handler[%d]: nil handler for prefix '%s'", i, j, config.Prefix)
			}

			// Path 可以为空（表示组路径本身，如 GET /users）
			// 不需要检查 Path 是否为空
		}
	}

	return nil
}
