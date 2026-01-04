package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// NewTestDB 创建用于测试的 SQLite 数据库
func NewTestDB(t *testing.T, models ...interface{}) (*gorm.DB, string, error) {
	// 为每个测试创建独立的数据库文件
	dbPath := filepath.Join(os.TempDir(), fmt.Sprintf("test_%d_%d.db", time.Now().UnixNano(), os.Getpid()))

	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent), // 测试时静默日志
	})
	if err != nil {
		return nil, "", fmt.Errorf("连接测试数据库失败: %w", err)
	}

	// 执行迁移
	if len(models) > 0 {
		if err := db.AutoMigrate(models...); err != nil {
			return nil, "", fmt.Errorf("数据库迁移失败: %w", err)
		}
	}

	// 注册清理函数
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		os.Remove(dbPath)
	})

	return db, dbPath, nil
}
