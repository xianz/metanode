package utils

import (
	"blog/config"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func GetSqliteDB(cfg *config.SqliteConfig) *gorm.DB {
	dbDir := filepath.Dir(cfg.Path)
	if dbDir != "." && dbDir != "/" {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			log.Fatal("创建数据库目录失败:", err)
		}
	}
	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogMode),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.TablePrefix,
			SingularTable: false,
			NoLowerCase:   false,
			NameReplacer:  nil,
		},
	})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("返回db了")
	return db
}

func SqliteFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
