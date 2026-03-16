package token

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"

	"app/internal/shared/cache"
	"github.com/bwmarrin/snowflake"
	"github.com/redis/go-redis/v9"
)

// SnowflakeTokenManager 基于雪花 ID 的 Token 管理器
type SnowflakeTokenManager struct {
	redisClient *redis.Client
	snowflake   *snowflake.Node
	userCache   *cache.UserCache
	keyPrefix   string
	userPrefix  string
	expireTime  time.Duration
}

// NewSnowflakeTokenManager 创建雪花 Token 管理器
func NewSnowflakeTokenManager(redisClient *redis.Client, userCache *cache.UserCache, expireTime time.Duration) (*SnowflakeTokenManager, error) {
	if expireTime == 0 {
		expireTime = 24 * time.Hour // 默认 24 小时
	}

	// 初始化雪花 ID 生成器
	node, err := snowflake.NewNode(1) // 节点 ID 可以配置化
	if err != nil {
		return nil, err
	}

	return &SnowflakeTokenManager{
		redisClient: redisClient,
		snowflake:   node,
		userCache:   userCache,
		keyPrefix:   "token:",
		userPrefix:  "user_tokens:",
		expireTime:  expireTime,
	}, nil
}

// GenerateTokenWithRoles 生成带角色的 Token
func (m *SnowflakeTokenManager) GenerateToken(userID int64, username string, roles []string) (string, error) {
	if len(roles) == 0 {
		roles = []string{"USER"}
	}

	// 生成雪花 ID 作为 Token
	token := m.snowflake.Generate().String()

	// 存储 Token -> UserID 映射（不存储完整用户信息）
	tokenKey := m.keyPrefix + token
	tokenData, err := json.Marshal(map[string]int64{"user_id": userID})
	if err != nil {
		return "", err
	}

	if err := m.redisClient.Set(context.Background(), tokenKey, tokenData, m.expireTime).Err(); err != nil {
		return "", err
	}

	// 建立用户索引：user_tokens:{user_id} -> [token_ids]
	userKey := m.userPrefix + strconv.FormatInt(userID, 10)
	if err := m.redisClient.SAdd(context.Background(), userKey, token).Err(); err != nil {
		return "", err
	}
	m.redisClient.Expire(context.Background(), userKey, m.expireTime)

	// 更新用户缓存（如果不存在则创建）
	_, err = m.userCache.GetUserInfo(context.Background(), userID)
	if errors.Is(err, cache.ErrCacheMiss) {
		// 缓存未命中，写入基本信息
		m.userCache.SetUserInfo(context.Background(), &cache.UserInfo{
			UserID:   userID,
			Username: username,
			Roles:    roles,
			Status:   "normal",
		})
	}

	return token, nil
}

// ValidateToken 验证 Token（返回 Claims 以符合 Manager 接口）
func (m *SnowflakeTokenManager) ValidateToken(token string) (*Claims, error) {
	// 从 Redis 获取 UserID
	tokenKey := m.keyPrefix + token
	tokenData, err := m.redisClient.Get(context.Background(), tokenKey).Result()
	if err == redis.Nil {
		return nil, ErrTokenExpired
	}
	if err != nil {
		return nil, err
	}

	// 解析 UserID
	var tokenInfo map[string]int64
	if err := json.Unmarshal([]byte(tokenData), &tokenInfo); err != nil {
		return nil, err
	}
	userID := tokenInfo["user_id"]

	// 通过用户缓存管理器获取用户信息（共用缓存！）
	userInfo, err := m.userCache.GetUserInfo(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	// 检查用户状态
	if userInfo.Status != "normal" {
		return nil, errors.New("user is frozen or canceled")
	}

	// 续期 Token TTL
	m.redisClient.Expire(context.Background(), tokenKey, m.expireTime)
	userKey := m.userPrefix + strconv.FormatInt(userID, 10)
	m.redisClient.Expire(context.Background(), userKey, m.expireTime)

	// 转换为 Claims 格式返回
	return &Claims{
		UserID:   userInfo.UserID,
		Username: userInfo.Username,
		Roles:    userInfo.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userInfo.Username,
		},
	}, nil
}

// RefreshToken 刷新 Token
func (m *SnowflakeTokenManager) RefreshToken(refreshToken string) (string, error) {
	claims, err := m.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	// 生成新 Token
	newToken, err := m.GenerateToken(claims.UserID, claims.Username, claims.Roles)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

// AddToBlacklist 添加 Token 到黑名单
func (m *SnowflakeTokenManager) AddToBlacklist(tokenString string, expireTime time.Time) error {
	// 直接撤销 Token
	return m.RevokeToken(tokenString)
}

// IsInBlacklist 检查 Token 是否在黑名单中
func (m *SnowflakeTokenManager) IsInBlacklist(tokenString string) (bool, error) {
	// 检查 Token 是否存在于 Redis 中
	tokenKey := m.keyPrefix + tokenString
	_, err := m.redisClient.Get(context.Background(), tokenKey).Result()
	if err == redis.Nil {
		return true, nil // 不存在即已在黑名单（被撤销）
	}
	if err != nil {
		return false, err
	}
	return false, nil // 存在说明有效
}

// GetTokenCount 获取用户的活跃 Token 数量
func (m *SnowflakeTokenManager) GetTokenCount(userID int64) (int64, error) {
	userKey := m.userPrefix + strconv.FormatInt(userID, 10)
	return m.redisClient.SCard(context.Background(), userKey).Result()
}

// RevokeToken 撤销指定 Token
func (m *SnowflakeTokenManager) RevokeToken(token string) error {
	tokenKey := m.keyPrefix + token

	// 获取用户 ID
	tokenData, err := m.redisClient.Get(context.Background(), tokenKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil // 已不存在
	}
	if err != nil {
		return err
	}

	var tokenInfo map[string]int64
	if err := json.Unmarshal([]byte(tokenData), &tokenInfo); err != nil {
		return err
	}

	// 删除 Token
	m.redisClient.Del(context.Background(), tokenKey)

	// 从用户索引中移除
	userKey := m.userPrefix + strconv.FormatInt(tokenInfo["user_id"], 10)
	m.redisClient.SRem(context.Background(), userKey, token)

	return nil
}

// RevokeAllUserTokens 撤销用户的所有 Token
func (m *SnowflakeTokenManager) RevokeAllUserTokens(userID int64) error {
	userKey := m.userPrefix + strconv.FormatInt(userID, 10)

	// 获取用户的所有 Token
	tokens, err := m.redisClient.SMembers(context.Background(), userKey).Result()
	if err != nil {
		return err
	}

	// 批量删除
	for _, token := range tokens {
		tokenKey := m.keyPrefix + token
		m.redisClient.Del(context.Background(), tokenKey)
	}

	// 删除用户索引
	m.redisClient.Del(context.Background(), userKey)

	return nil
}

// KickoutUser 踢用户下线
func (m *SnowflakeTokenManager) KickoutUser(userID int64) error {
	return m.RevokeAllUserTokens(userID)
}

// GetUserInfo 获取用户信息（用于 Manager 接口）
func (m *SnowflakeTokenManager) GetUserInfo(tokenString string) (*UserInfo, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		UserID:   claims.UserID,
		Username: claims.Username,
		Roles:    claims.Roles,
		Status:   "normal",
	}, nil
}

// 错误定义
var (
	ErrTokenExpired = errors.New("token expired or not found")
)
