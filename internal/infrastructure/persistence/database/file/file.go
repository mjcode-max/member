package file

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"member-pre/internal/infrastructure/config"
)

// DB File数据库实例（使用SQLite）
type DB struct {
	db       *gorm.DB
	filePath string
}

// NewDB 创建File数据库连接（用于测试）
// 返回类型使用接口，避免循环依赖
func NewDB(cfg *config.Config) (*DB, error) {
	dbCfg := &cfg.Database

	// 如果未指定文件路径，使用默认路径
	filePath := dbCfg.FilePath
	if filePath == "" {
		filePath = "test.db"
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("创建数据库目录失败: %w", err)
		}
	}

	// 连接SQLite数据库
	db, err := gorm.Open(sqlite.Open(filePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 测试时静默日志
	})
	if err != nil {
		return nil, fmt.Errorf("连接File数据库失败: %w", err)
	}

	// 测试连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("File数据库连接测试失败: %w", err)
	}

	return &DB{
		db:       db,
		filePath: filePath,
	}, nil
}

// DB 实现Database接口
func (d *DB) DB() *gorm.DB {
	return d.db
}

// Close 关闭数据库连接
func (d *DB) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Clear 清除数据库文件（用于测试）
func (d *DB) Clear() error {
	if err := d.Close(); err != nil {
		return err
	}
	if d.filePath != "" {
		return os.Remove(d.filePath)
	}
	return nil
}
