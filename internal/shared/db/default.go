package db

import (
	"sync"

	"gorm.io/gorm"
)

var (
	defaultDB     *gorm.DB
	defaultRemote *gorm.DB
	mu            sync.RWMutex
)

// Default 返回默认的数据库实例
func Default() *gorm.DB {
	mu.RLock()
	defer mu.RUnlock()
	return defaultDB
}

// SetDefault 设置默认的数据库实例
func SetDefault(db *gorm.DB) {
	mu.Lock()
	defer mu.Unlock()
	defaultDB = db
}

// MustDefault 返回默认的数据库实例
// 如果未初始化，会 panic
func MustDefault() *gorm.DB {
	mu.RLock()
	defer mu.RUnlock()
	if defaultDB == nil {
		panic("db: default database not initialized, call SetDefault first")
	}
	return defaultDB
}

// Remote 返回远程数据库实例
func Remote() *gorm.DB {
	mu.RLock()
	defer mu.RUnlock()
	return defaultRemote
}

// SetRemote 设置远程数据库实例
func SetRemote(db *gorm.DB) {
	mu.Lock()
	defer mu.Unlock()
	defaultRemote = db
}

// MustRemote 返回远程数据库实例
// 如果未初始化，会 panic
func MustRemote() *gorm.DB {
	mu.RLock()
	defer mu.RUnlock()
	if defaultRemote == nil {
		panic("db: remote database not initialized, call SetRemote first")
	}
	return defaultRemote
}
