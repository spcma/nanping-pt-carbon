package cache

import (
	"sync"
)

var (
	mu sync.RWMutex
)

// Default 返回默认的 Redis 客户端实例
// 如果未初始化，返回 nil
func Default() *RedisClient {
	mu.RLock()
	defer mu.RUnlock()
	return RDS
}

// MustDefault 返回默认的 Redis 客户端实例
// 如果未初始化，会 panic
func MustDefault() *RedisClient {
	mu.RLock()
	defer mu.RUnlock()
	if RDS == nil {
		panic("cache: default redis client not initialized, call SetDefault first")
	}
	return RDS
}
