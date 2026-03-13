package persistence

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type Config struct {
	DSN                  string        `mapstructure:"dsn"`
	MaxIdleConns         int           `mapstructure:"max_idle_conns"`
	MaxOpenConns         int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime      time.Duration `mapstructure:"conn_max_lifetime"`
	SlowThreshold        time.Duration `mapstructure:"slow_threshold"`
	LogLevel             string        `mapstructure:"log_level"`
	Colorful             bool          `mapstructure:"colorful"`
	IgnoreRecordNotFound bool          `mapstructure:"ignore_record_not_found"`
}

func NewGORM(config Config) (*gorm.DB, error) {
	// 验证配置
	if config.DSN == "" {
		return nil, errors.New("DSN cannot be empty")
	}

	// 设置默认值
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 10
	}
	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 100
	}
	if config.ConnMaxLifetime == 0 {
		config.ConnMaxLifetime = time.Hour
	}
	if config.SlowThreshold == 0 {
		config.SlowThreshold = 200 * time.Millisecond
	}
	if config.LogLevel == "" {
		config.LogLevel = "warn"
	}

	var logLevel logger.LogLevel
	switch config.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Warn
	}

	newLogger := logger.Default.LogMode(logLevel)
	if config.Colorful {
		newLogger = newLogger.LogMode(logger.Info) // 启用彩色输出
	}

	gormConfig := &gorm.Config{
		Logger: newLogger,
	}

	if config.IgnoreRecordNotFound {
		gormConfig.SkipDefaultTransaction = true
	}

	db, err := gorm.Open(postgres.Open(config.DSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying DB object: %w", err)
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}
