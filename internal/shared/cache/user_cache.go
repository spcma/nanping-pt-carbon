package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// UserCache 用户缓存管理器
type UserCache struct {
	redisClient *redis.Client
	keyPrefix   string
	expireTime  time.Duration
}

// NewUserCache 创建用户缓存管理器
func NewUserCache(redisClient *redis.Client, expireTime time.Duration) *UserCache {
	if expireTime == 0 {
		expireTime = 1 * time.Hour // 默认 1 小时
	}

	return &UserCache{
		redisClient: redisClient,
		keyPrefix:   "user:",
		expireTime:  expireTime,
	}
}

// GetUserInfo 获取用户信息（带缓存）
func (c *UserCache) GetUserInfo(ctx context.Context, userID int64) (*UserInfo, error) {
	key := c.keyPrefix + strconv.FormatInt(userID, 10)

	// 1. 尝试从缓存获取
	data, err := c.redisClient.Get(ctx, key).Result()
	if err == nil {
		var info UserInfo
		if err := json.Unmarshal([]byte(data), &info); err == nil {
			// 续期
			c.redisClient.Expire(ctx, key, c.expireTime)
			return &info, nil
		}
	}

	// 2. 缓存未命中，返回错误（由调用方从数据库加载）
	return nil, ErrCacheMiss
}

// SetUserInfo 设置用户信息到缓存
func (c *UserCache) SetUserInfo(ctx context.Context, info *UserInfo) error {
	key := c.keyPrefix + strconv.FormatInt(info.UserID, 10)
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return c.redisClient.Set(ctx, key, data, c.expireTime).Err()
}

// DeleteUserInfo 删除用户缓存
func (c *UserCache) DeleteUserInfo(ctx context.Context, userID int64) error {
	key := c.keyPrefix + strconv.FormatInt(userID, 10)
	return c.redisClient.Del(ctx, key).Err()
}

// UpdateUserInfo 更新用户信息（原子操作）
func (c *UserCache) UpdateUserInfo(ctx context.Context, userID int64, updater func(*UserInfo) *UserInfo) error {
	key := c.keyPrefix + strconv.FormatInt(userID, 10)

	// 使用 Redis 事务
	return c.redisClient.Watch(ctx, func(tx *redis.Tx) error {
		// 获取当前值
		data, err := tx.Get(ctx, key).Result()
		if err != nil && err != redis.Nil {
			return err
		}

		var info *UserInfo
		if err == nil {
			if err := json.Unmarshal([]byte(data), &info); err != nil {
				return err
			}
		}

		// 调用更新函数
		info = updater(info)

		// 写回
		dataBytes, err := json.Marshal(info)
		if err != nil {
			return err
		}
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, dataBytes, c.expireTime)
			return nil
		})
		return err
	}, key)
}

// UserInfo 用户信息（用于缓存）
type UserInfo struct {
	UserID   int64    `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Status   string   `json:"status"` // normal, frozen, canceled
}

// 错误定义
var (
	ErrCacheMiss = fmt.Errorf("cache miss")
)
