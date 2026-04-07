package token

import (
	"app/internal/shared/cache"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenType Token 类型
type TokenType int

const (
	// TokenType_JWT 传统 JWT 模式（所有信息在 token 中）
	TokenType_JWT TokenType = iota

	// TokenType_Snowflake 雪花 ID + 缓存模式（token 只是 ID，用户信息在 Redis 缓存）
	TokenType_Snowflake
)

// ConfigEx 扩展配置（用于创建 Token 管理器）
type ConfigEx struct {
	Type        TokenType        // Token 类型
	JWTConfig   Config           // JWT 配置（当 Type=TokenTypeJWT 时使用）
	RedisClient *redis.Client    // Redis 客户端（当 Type=TokenTypeSnowflake 时使用）
	UserCache   *cache.UserCache // 用户缓存（当 Type=TokenTypeSnowflake 时使用）
	ExpireTime  time.Duration    // 过期时间
}

// NewManager 创建 Token 管理器（工厂方法）
func NewManager(config ConfigEx) (Manager, error) {
	switch config.Type {
	case TokenType_JWT:
		// 使用传统 JWT 模式
		return NewJWTManager(config.JWTConfig), nil

	case TokenType_Snowflake:
		// 使用雪花 ID + 缓存模式
		if config.RedisClient == nil {
			return nil, ErrRedisClientRequired
		}
		if config.UserCache == nil {
			config.UserCache = cache.NewUserCache(config.RedisClient, config.ExpireTime)
		}
		return NewSnowflakeTokenManager(config.RedisClient, config.UserCache, config.ExpireTime)

	default:
		return nil, ErrUnknownTokenType
	}
}

// 错误定义
var (
	ErrRedisClientRequired = &Error{"redis client is required for Snowflake token mode"}
	ErrUnknownTokenType    = &Error{"unknown token type"}
)

// Error Token 错误
type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}
