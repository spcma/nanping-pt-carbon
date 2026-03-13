package config

import (
	"app/internal/shared/logger"
	"fmt"

	"github.com/spf13/viper"
)

var GlobalConfig *Config

type Config struct {
	Server   ServerConfig     `mapstructure:"server"`
	Logger   logger.LogConfig `mapstructure:"logger"`
	Database DatabaseConfig   `mapstructure:"database"`
	Redis    RedisConfig      `mapstructure:"redis"`
	JWT      JWTConfig        `mapstructure:"jwt"`
	Idgen    IdgenConfig      `mapstructure:"idgen"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release, test
}

type DatabaseConfig struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
	DBName     string `mapstructure:"dbname"`
	SearchPath string `mapstructure:"searchpath"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"` // seconds
}

type IdgenConfig struct {
	WorkerID int `mapstructure:"worker_id"`
}

// Init 初始化配置
func Init(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
