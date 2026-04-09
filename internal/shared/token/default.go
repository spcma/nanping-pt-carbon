package token

import "sync"

var (
	defaultManager Manager
	once           sync.Once
	mu             sync.RWMutex
)

// Default 返回默认的 Token 管理器实例
// 如果未初始化，返回 nil
func Default() Manager {
	mu.RLock()
	defer mu.RUnlock()
	return defaultManager
}

// SetDefault 设置默认的 Token 管理器实例
func SetDefault(m Manager) {
	mu.Lock()
	defer mu.Unlock()
	defaultManager = m
}

// MustDefault 返回默认的 Token 管理器实例
// 如果未初始化，会 panic
func MustDefault() Manager {
	mu.RLock()
	defer mu.RUnlock()
	if defaultManager == nil {
		panic("token: default manager not initialized, call SetDefault first")
	}
	return defaultManager
}
