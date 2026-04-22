package db

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
)

var (
	// defaultDB 默认数据库实例
	defaultDB *gorm.DB
	// databases 多数据源管理器，key 为数据源名称
	databases = make(map[string]*gorm.DB)
	mu        sync.RWMutex
)

const (
	// DefaultDataSourceName 默认数据源名称
	DefaultDataSourceName = "default"
	RemoteDataSourceName  = "remote"
)

var datasources = map[string]string{
	"default": DefaultDataSourceName,
	"remote":  RemoteDataSourceName,
}

// Default 返回默认的数据库实例
func Default() *gorm.DB {
	return GetDBByName(DefaultDataSourceName)
}

// SetDefault 设置默认的数据库实例
func SetDefault(db *gorm.DB) {
	RegisterDB(DefaultDataSourceName, db)
}

// MustDefault 返回默认的数据库实例
// 如果未初始化，会 panic
func MustDefault() *gorm.DB {
	db := GetDBByName(DefaultDataSourceName)
	if db == nil {
		panic("db: default database not initialized, call SetDefault or RegisterDB first")
	}
	return db
}

// RegisterDB 注册一个命名数据源
func RegisterDB(name string, db *gorm.DB) {
	if name == "" {
		panic("db: data source name cannot be empty")
	}
	if db == nil {
		panic("db: database instance cannot be nil")
	}

	if _, ok := datasources[name]; !ok {
		//	该数据源不允许被注册
		panicf("db: data source '%s' is not allowed to be registered", name)
	}

	mu.Lock()
	defer mu.Unlock()
	databases[name] = db

	// 如果是默认数据源，同时更新 defaultDB
	if name == DefaultDataSourceName {
		defaultDB = db
	}
}

// GetDBByName 根据名称获取数据库实例
func GetDBByName(name string) *gorm.DB {
	if name == "" {
		name = DefaultDataSourceName
	}

	if _, ok := datasources[name]; !ok {
		panicf("db: data source '%s' not registered", name)
	}

	mu.RLock()
	defer mu.RUnlock()
	return databases[name]
}

// MustGetDBByName 根据名称获取数据库实例，如果不存在则 panic
func MustGetDBByName(name string) *gorm.DB {
	if name == "" {
		name = DefaultDataSourceName
	}

	if _, ok := datasources[name]; !ok {
		panicf("db: data source '%s' not registered", name)
	}

	mu.RLock()
	defer mu.RUnlock()
	db := databases[name]
	if db == nil {
		panicf("db: database '%s' not registered", name)
	}
	return db
}

// RemoveDB 移除一个命名数据源
func RemoveDB(name string) {
	if name == DefaultDataSourceName {
		panic("db: cannot remove default database")
	}

	mu.Lock()
	defer mu.Unlock()
	delete(databases, name)
}

// GetAllDBNames 获取所有已注册的数据源名称
func GetAllDBNames() []string {
	mu.RLock()
	defer mu.RUnlock()

	names := make([]string, 0, len(databases))
	for name := range databases {
		names = append(names, name)
	}
	return names
}

// HasDB 检查是否存在指定名称的数据源
func HasDB(name string) bool {
	if name == "" {
		name = DefaultDataSourceName
	}

	mu.RLock()
	defer mu.RUnlock()
	_, exists := databases[name]
	return exists
}

// panicf 格式化 panic 消息
func panicf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}
