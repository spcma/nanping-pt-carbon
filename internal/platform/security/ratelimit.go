package security

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiter 限流器接口
type RateLimiter interface {
	// Allow 检查是否允许请求
	Allow(ctx context.Context, key string) (bool, error)
	// GetRate 获取当前速率信息
	GetRate(ctx context.Context, key string) (*RateInfo, error)
}

// RateInfo 速率信息
type RateInfo struct {
	Allowed     bool  `json:"allowed"`
	Remaining   int64 `json:"remaining"`
	ResetTime   int64 `json:"reset_time"`
	Limit       int64 `json:"limit"`
	RequestTime int64 `json:"request_time"`
}

// RedisRateLimiter Redis限流器实现
type RedisRateLimiter struct {
	client *redis.Client
	limit  int64
	window time.Duration
	prefix string
}

// NewRedisRateLimiter 创建Redis限流器
func NewRedisRateLimiter(client *redis.Client, limit int64, window time.Duration, prefix string) *RedisRateLimiter {
	return &RedisRateLimiter{
		client: client,
		limit:  limit,
		window: window,
		prefix: prefix,
	}
}

// Allow 检查是否允许请求
func (r *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	fullKey := fmt.Sprintf("%s:%s", r.prefix, key)
	now := time.Now().Unix()
	windowStart := now - int64(r.window.Seconds())

	// 清理过期的计数
	err := r.client.ZRemRangeByScore(ctx, fullKey, "0", fmt.Sprintf("%d", windowStart)).Err()
	if err != nil {
		return false, err
	}

	// 获取当前窗口内的请求数
	count, err := r.client.ZCard(ctx, fullKey).Result()
	if err != nil {
		return false, err
	}

	// 检查是否超过限制
	if count >= r.limit {
		return false, nil
	}

	// 添加当前请求
	err = r.client.ZAdd(ctx, fullKey, redis.Z{
		Score:  float64(now),
		Member: now,
	}).Err()
	if err != nil {
		return false, err
	}

	// 设置过期时间
	err = r.client.Expire(ctx, fullKey, r.window).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetRate 获取当前速率信息
func (r *RedisRateLimiter) GetRate(ctx context.Context, key string) (*RateInfo, error) {
	fullKey := fmt.Sprintf("%s:%s", r.prefix, key)
	now := time.Now().Unix()
	windowStart := now - int64(r.window.Seconds())

	// 清理过期的计数
	err := r.client.ZRemRangeByScore(ctx, fullKey, "0", fmt.Sprintf("%d", windowStart)).Err()
	if err != nil {
		return nil, err
	}

	// 获取当前计数
	count, err := r.client.ZCard(ctx, fullKey).Result()
	if err != nil {
		return nil, err
	}

	return &RateInfo{
		Allowed:     count < r.limit,
		Remaining:   r.limit - count,
		ResetTime:   now + int64(r.window.Seconds()),
		Limit:       r.limit,
		RequestTime: now,
	}, nil
}

// MemoryRateLimiter 内存限流器实现（备用方案）
type MemoryRateLimiter struct {
	limit    int64
	window   time.Duration
	requests map[string][]int64
	prefix   string
}

// NewMemoryRateLimiter 创建内存限流器
func NewMemoryRateLimiter(limit int64, window time.Duration, prefix string) *MemoryRateLimiter {
	return &MemoryRateLimiter{
		limit:    limit,
		window:   window,
		requests: make(map[string][]int64),
		prefix:   prefix,
	}
}

// Allow 检查是否允许请求
func (m *MemoryRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	fullKey := fmt.Sprintf("%s:%s", m.prefix, key)
	now := time.Now().Unix()
	windowStart := now - int64(m.window.Seconds())

	// 清理过期的请求
	requests := m.requests[fullKey]
	filtered := make([]int64, 0)
	for _, reqTime := range requests {
		if reqTime > windowStart {
			filtered = append(filtered, reqTime)
		}
	}
	m.requests[fullKey] = filtered

	// 检查是否超过限制
	if int64(len(filtered)) >= m.limit {
		return false, nil
	}

	// 添加当前请求
	m.requests[fullKey] = append(m.requests[fullKey], now)
	return true, nil
}

// GetRate 获取当前速率信息
func (m *MemoryRateLimiter) GetRate(ctx context.Context, key string) (*RateInfo, error) {
	fullKey := fmt.Sprintf("%s:%s", m.prefix, key)
	now := time.Now().Unix()
	windowStart := now - int64(m.window.Seconds())

	// 清理过期的请求
	requests := m.requests[fullKey]
	filtered := make([]int64, 0)
	for _, reqTime := range requests {
		if reqTime > windowStart {
			filtered = append(filtered, reqTime)
		}
	}
	m.requests[fullKey] = filtered

	count := int64(len(filtered))

	return &RateInfo{
		Allowed:     count < m.limit,
		Remaining:   m.limit - count,
		ResetTime:   now + int64(m.window.Seconds()),
		Limit:       m.limit,
		RequestTime: now,
	}, nil
}
