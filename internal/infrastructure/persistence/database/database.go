package database

import (
	"fmt"
	"gorm.io/gorm"
	"member-pre/internal/infrastructure/config"
	"member-pre/internal/infrastructure/persistence/database/file"
	"member-pre/internal/infrastructure/persistence/database/mysql"
)

// Database 数据库抽象接口
type Database interface {
	// DB 获取GORM数据库实例
	DB() *gorm.DB

	// Close 关闭数据库连接
	Close() error
}

// NewDatabase 根据配置创建数据库实例
func NewDatabase(cfg *config.Config) (Database, error) {
	dbType := cfg.Database.Type
	if dbType == "" {
		dbType = "mysql" // 默认使用mysql
	}

	switch dbType {
	case "mysql":
		return mysql.NewDB(cfg)
	case "file":
		return file.NewDB(cfg)
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s，支持的类型: mysql, postgres, file", dbType)
	}
}
