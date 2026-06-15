package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

func main() {
	// New一个最好，避免全局变量的污染
	v := viper.New()

	// 基础设置
	v.SetConfigName("config") // 配置文件名
	v.SetConfigType("yaml")   // 配置文件扩展名
	v.AddConfigPath(".")      // 配置文件所在路径

	// 设置默认值（优先级最低）
	v.SetDefault("server", map[string]interface{}{
		"port": 8080,
		"mode": "debug",
	})
	v.SetDefault("database", map[string]interface{}{
		"host":     "127.0.0.1",
		"port":     3306,
		"user":     "dev",
		"password": "dev",
		"name":     "test",
	})
	// v.SetDefault("server.port", 8080)
	// v.SetDefault("server.mode", "debug")
	// v.SetDefault("database.host", "localhost")
	// v.SetDefault("database.port", 3306)
	// v.SetDefault("database.user", "root")
	// v.SetDefault("database.password", "password")
	// v.SetDefault("database.name", "test")

	// 读取文件
	err := v.ReadInConfig()
	if err != nil {
		// if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Printf("使用默认值（读取配置文件出错：%v）", err)
	}

	// 热更新
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("配置文件已更新，正在重新加载... %s", e.Name)
		if err := reloadConfig(v); err != nil {
			log.Printf("重新加载配置失败: %v", err)
		} else {
			log.Printf("配置热加载成功")
		}
	})

	// 支持环境变量
	v.SetEnvPrefix("MYAPP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 解析配置
	var config Config
	err = v.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	// 使用配置
	runApp(config)
}

func reloadConfig(v *viper.Viper) error {
	var config Config
	err := v.Unmarshal(&config)
	if err != nil {
		return err
	}
	// 这里可以重新初始化一些依赖配置的组件
	println("配置已更新")
	printConfig(config)
	return nil
}

func runApp(cfg Config) {
	// 打印配置验证
	printConfig(cfg)

	// 实际运行逻辑
	fmt.Println("应用正在运行... (修改 config/config.yaml 文件查看热更新效果)")
	select {} // 阻塞主线程，保持热更新监听
}

func printConfig(cfg Config) {
	fmt.Println("========== 应用配置信息 ==========")
	fmt.Printf("Server Port: %d\n", cfg.Server.Port)
	fmt.Printf("Server Mode: %s\n", cfg.Server.Mode)
	fmt.Printf("DB Host: %s\n", cfg.Database.Host)
	fmt.Printf("DB Port: %d\n", cfg.Database.Port)
	fmt.Println("==================================")
}
