package utils

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
)

type Config struct {
	Db struct {
		Host    string          `mapstructure:"host"`
		Port    int             `mapstructure:"port"`
		Address string          `mapstructure:"address"`
		LogMode logger.LogLevel `mapstructure:"log_mode"`
	} `mapstructure:"db"`
}

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		fmt.Println("未找到 .env 文件")
	}

	// 初始化 viper
	viper.SetConfigFile("config.yaml") // 直接指定配置文件

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("读取配置文件失败: %v\n", err)
	}

	// 设置环境变量支持
	viper.SetEnvPrefix("BOX")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 设置默认值
	viper.SetDefault("db.host", "default-host")
	viper.SetDefault("db.port", 9999)

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("解析配置失败: %v\n", err)
	}
	fmt.Println(config)

	// 输出结果
	fmt.Printf("db.host = %v\n", viper.Get("db.host"))
	fmt.Printf("db.port = %v\n", viper.GetInt("db.port"))
	fmt.Printf("db.address = %v\n", viper.Get("db.address"))
	fmt.Printf("db.log_mode = %v\n", viper.GetInt("db.log_mode"))

	// // 检查环境变量原始值
	// fmt.Printf("\n实际环境变量:\n")
	// fmt.Printf("APP_DB_HOST = %v\n", os.Getenv("APP_DB_HOST"))
	// fmt.Printf("APP_DB_PORT = %v\n", os.Getenv("APP_DB_PORT"))

}
