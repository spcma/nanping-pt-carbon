package db

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Driver     string
	Host       string
	Port       int
	User       string
	Password   string
	DbName     string
	SearchPath string
	// Name 数据源名称，为空时使用 "default"
	Name string
}

func NewGormDB(cfg Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d search_path=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host, cfg.User, cfg.Password, cfg.DbName, cfg.Port, cfg.SearchPath,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect db:", err)
		return nil
	}

	// 确定数据源名称
	name := cfg.Name
	if name == "" {
		name = DefaultDataSourceName
	}

	// 注册到多数据源管理器
	RegisterDB(name, db)

	return db
}

// GetDB 从 context 中获取数据库实例
// 如果 context 中有事务 DB，则返回事务 DB
// 否则返回默认数据源
func GetDB(ctx context.Context) *gorm.DB {
	// 检查context中是否有事务DB
	if txDB, ok := ctx.Value("tx_db").(*gorm.DB); ok {
		return txDB
	}
	// 否则使用默认数据源
	return Default()
}

// GetDBWithContext 从 context 中获取指定数据源的数据库实例
// 优先检查 context 中的事务 DB
// 如果未指定数据源名称，使用默认数据源
func GetDBWithContext(ctx context.Context, dbName string) *gorm.DB {
	// 检查context中是否有事务DB
	if txDB, ok := ctx.Value("tx_db").(*gorm.DB); ok {
		return txDB
	}
	// 否则使用指定数据源
	return GetDBByName(dbName)
}
