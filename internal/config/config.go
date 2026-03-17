package config

import (
	"app/internal/shared/logger"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var GlobalConfig *Config

type Config struct {
	Server   ServerConfig     `mapstructure:"server"`
	Logger   logger.LogConfig `mapstructure:"logger"`
	Database DatabaseConfig   `mapstructure:"database"`
	Redis    RedisConfig      `mapstructure:"redis"`
	Token    TokenConfig      `mapstructure:"token"`
	Idgen    IdgenConfig      `mapstructure:"idgen"`
	Ipfs     IpfsConfig       `mapstructure:"ipfs"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release, test
}

type IpfsConfig struct {
	Status bool `mapstructure:"status"` //  是否开启ipfs服务, 运行环境中没有ipfs会导致初始化失败
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
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	Timeout      time.Duration `mapstructure:"timeout"`
}

type TokenConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
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
