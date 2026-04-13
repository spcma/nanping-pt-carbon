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
	NotDefault bool
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

	if !cfg.NotDefault {
		defaultDB = db
	}

	return db
}

func GetDB(ctx context.Context) *gorm.DB {
	// 检查context中是否有事务DB
	if txDB, ok := ctx.Value("tx_db").(*gorm.DB); ok {
		return txDB
	}
	// 否则使用普通DB
	return defaultDB
}
