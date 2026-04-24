package config

import (
	"app/internal/shared/logger"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var GlobalConfig *Config

type Config struct {
	Server         ServerConfig     `mapstructure:"server"`
	Logger         logger.LogConfig `mapstructure:"logger"`
	Database       DatabaseConfig   `mapstructure:"database"`
	RemoteDatabase DatabaseConfig   `mapstructure:"remote_database"`
	Redis          RedisConfig      `mapstructure:"redis"`
	Token          TokenConfig      `mapstructure:"token"`
	Ipfs           IpfsConfig       `mapstructure:"ipfs"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release, test
}

type IpfsConfig struct {
	//Host   string `mapstructure:"host"`
	Port   int  `mapstructure:"port"`
	Status bool `mapstructure:"status"` //  是否开启ipfs服务, 运行环境中没有ipfs会导致初始化失败
}

type DatabaseConfig struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
	DBName     string `mapstructure:"dbname"`
	SearchPath string `mapstructure:"searchpath"`
	Name       string `mapstructure:"name"`
}

type RedisConfig struct {
	Addr         string        `mapstructure:"addr"`
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
