package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func GenerateSQLiteFilePath(filename string) (string, error) {
	dbDir, err := GetDBDir()
	if err != nil {
		return "", err
	}

	// Ensure filename contains "sqlite"
	if filename == "" {
		filename = "test.sqlite.db"
	} else {
		// Extract base name and extension
		ext := filepath.Ext(filename)
		if ext == "" {
			ext = ".db"
		}
		// fmt.Println("ext:", ext)
		base := filename[:len(filename)-len(ext)]
		// fmt.Println("base:", base)
		if base == "" {
			base = "test"
		}
		// Check if "sqlite" is already in the filename (case-insensitive check)
		baseLower := strings.ToLower(base)
		if !strings.Contains(baseLower, "sqlite") {
			filename = base + "_sqlite" + ext
		} else {
			filename = base + ext
		}
	}

	dbPath := filepath.Join(dbDir, filename)
	return dbPath, nil
}

func OpenSqlite(dbPath string) *gorm.DB {
	// Configure GORM with:
	// 1. Logger: Control SQL logging level
	//    - Silent: No logs
	//    - Error: Only errors
	//    - Warn: Errors and warnings
	//    - Info: All SQL queries (default)
	// 2. NamingStrategy: Customize table and column naming
	//    - TableName: How struct names map to table names
	//    - ColumnName: How field names map to column names
	//    - JoinTableName: How join table names are generated
	//    - SchemaName: Schema name for databases that support it
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		// Logger configuration
		Logger: logger.Default.LogMode(logger.Info), // Silent for tests, use logger.Info for development

		// NamingStrategy: Customize how GORM names tables and columns
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",    // Prefix for all table names (e.g., "app_")
			SingularTable: false, // Use singular table names (User -> user instead of users)
			NoLowerCase:   false, // Disable automatic lowercasing
			NameReplacer:  nil,   // Custom name replacer function
		},
	})
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("get generic db: %v", err)
	}

	// Connection pool configuration:
	// - SetMaxIdleConns: Maximum number of idle connections in the pool
	//   (connections that are open but not currently in use)
	// - SetMaxOpenConns: Maximum number of open connections to the database
	//   (total connections, including idle and in-use)
	// - SetConnMaxLifetime: Maximum amount of time a connection may be reused
	//   (prevents using stale connections)
	sqlDB.SetMaxIdleConns(2)                   // Keep 2 idle connections ready
	sqlDB.SetMaxOpenConns(5)                   // Allow up to 5 concurrent connections
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Reuse connections for up to 30 minutes

	return db
}

func GetDBDir() (string, error) {
	// Get the directory where this file (helpers.go) is located
	// This file is in examples/testutil directory
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", os.ErrInvalid
	}

	// Get the examples directory (parent of testutil)
	testutilDir := filepath.Dir(currentFile)
	examplesDir := filepath.Dir(testutilDir)

	// The db directory is examples/db
	dbDir := filepath.Join(examplesDir, "db")

	// Ensure db directory exists
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return "", err
	}

	return dbDir, nil
}
