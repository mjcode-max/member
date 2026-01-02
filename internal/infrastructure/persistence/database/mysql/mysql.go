package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"member-pre/internal/infrastructure/config"
)

// DB MySQL数据库实例
type DB struct {
	db *gorm.DB
}

// NewDB 创建MySQL数据库连接
// 返回类型使用具体类型，避免循环依赖
func NewDB(cfg *config.Config) (*DB, error) {
	dbCfg := &cfg.Database
	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.DBName,
		dbCfg.Charset,
		dbCfg.ParseTime,
		dbCfg.Loc,
	)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接MySQL数据库失败: %w", err)
	}

	// 获取底层sql.DB设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbCfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(dbCfg.ConnMaxLifetime * time.Second)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("MySQL数据库连接测试失败: %w", err)
	}

	return &DB{db: db}, nil
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
