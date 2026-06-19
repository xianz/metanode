package config

/*
声明配置结构体、加载配置文件
*/

import (
	"fmt"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
)

var configFilePath = "config.yaml"

// // 绑定配置文件的结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Jwt      JwtConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	GinMode string `mapstructure:"gin_mode"`
}

type DatabaseConfig struct {
	Sqlite SqliteConfig `mapstructure:"sqlite"`
}

type SqliteConfig struct {
	Path        string          `mapstructure:"path"`
	LogMode     logger.LogLevel `mapstructure:"log_mode"`
	TablePrefix string          `mapstructure:"table_prefix"`
}

type JwtConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// 加载配置文件
func LoadConfig(config *Config) error {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		fmt.Println("未找到 .env 文件")
	}

	v := viper.New()
	// 读取config.yaml 文件
	v.SetConfigFile(configFilePath)
	v.SetConfigType("yaml")

	v.SetEnvPrefix("BLOG")                             // 环境变量前缀 BLOG_
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 支持嵌套，如 blog.jwt.secret -> BLOG_JWT_SECRET
	v.AutomaticEnv()                                   // 自动读取环境变量

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("配置文件不存在，使用默认配置: %v", err)
			return nil
		} else {
			log.Printf("加载配置文件失败: %v", err)
			return err
		}
	}
	// 反序列化配置
	if err := v.Unmarshal(config); err != nil {
		return err
	}

	// 热更新
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		if err := v.Unmarshal(config); err != nil {
			log.Printf("重新加载配置失败: %v", err)
		}
		log.Printf("配置已更新: %v", e.Name)
	})

	return nil
}
