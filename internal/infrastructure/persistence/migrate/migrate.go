package migrate

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"member-pre/internal/infrastructure/logger"
)

// Migration 迁移接口
type Migration interface {
	Up(db *gorm.DB) error
	Down(db *gorm.DB) error
	Name() string
}

// migrations 所有迁移
var migrations []Migration

// Register 注册迁移
func Register(m Migration) {
	migrations = append(migrations, m)
}

// Up 执行所有迁移
func Up(db *gorm.DB, log *logger.ZapLogger) error {
	// 创建迁移记录表
	if err := createMigrationTable(db); err != nil {
		return fmt.Errorf("创建迁移记录表失败: %w", err)
	}

	for _, m := range migrations {
		// 检查是否已执行
		executed, err := isMigrationExecuted(db, m.Name())
		if err != nil {
			return fmt.Errorf("检查迁移状态失败: %w", err)
		}

		if executed {
			log.Info("迁移已执行，跳过", zap.String("migration", m.Name()))
			continue
		}

		// 执行迁移
		log.Info("开始执行迁移", zap.String("migration", m.Name()))
		if err := m.Up(db); err != nil {
			return fmt.Errorf("执行迁移失败 [%s]: %w", m.Name(), err)
		}

		// 记录迁移
		if err := recordMigration(db, m.Name()); err != nil {
			return fmt.Errorf("记录迁移失败: %w", err)
		}

		log.Info("迁移执行成功", zap.String("migration", m.Name()))
	}

	return nil
}

// Down 回滚最后一次迁移
func Down(db *gorm.DB, log *logger.ZapLogger) error {
	if len(migrations) == 0 {
		log.Info("没有可回滚的迁移")
		return nil
	}

	// 获取最后一次迁移
	lastMigration := migrations[len(migrations)-1]

	// 检查是否已执行
	executed, err := isMigrationExecuted(db, lastMigration.Name())
	if err != nil {
		return fmt.Errorf("检查迁移状态失败: %w", err)
	}

	if !executed {
		log.Info("迁移未执行，无需回滚", zap.String("migration", lastMigration.Name()))
		return nil
	}

	// 执行回滚
	log.Info("开始回滚迁移", zap.String("migration", lastMigration.Name()))
	if err := lastMigration.Down(db); err != nil {
		return fmt.Errorf("回滚迁移失败 [%s]: %w", lastMigration.Name(), err)
	}

	// 删除迁移记录
	if err := removeMigrationRecord(db, lastMigration.Name()); err != nil {
		return fmt.Errorf("删除迁移记录失败: %w", err)
	}

	log.Info("迁移回滚成功", zap.String("migration", lastMigration.Name()))
	return nil
}

// Status 查看迁移状态
func Status(db *gorm.DB, log *logger.ZapLogger) error {
	// 创建迁移记录表（如果不存在）
	if err := createMigrationTable(db); err != nil {
		return fmt.Errorf("创建迁移记录表失败: %w", err)
	}

	log.Info("迁移状态:")
	for _, m := range migrations {
		executed, err := isMigrationExecuted(db, m.Name())
		if err != nil {
			log.Error("检查迁移状态失败", zap.String("migration", m.Name()), zap.Error(err))
			continue
		}

		status := "未执行"
		if executed {
			status = "已执行"
		}
		log.Info(fmt.Sprintf("  %s: %s", m.Name(), status))
	}

	return nil
}

// MigrationRecord 迁移记录
type MigrationRecord struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"uniqueIndex;not null"`
	ExecutedAt int64  `gorm:"not null"`
}

// createMigrationTable 创建迁移记录表
func createMigrationTable(db *gorm.DB) error {
	return db.AutoMigrate(&MigrationRecord{})
}

// isMigrationExecuted 检查迁移是否已执行
func isMigrationExecuted(db *gorm.DB, name string) (bool, error) {
	var count int64
	err := db.Model(&MigrationRecord{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// recordMigration 记录迁移
func recordMigration(db *gorm.DB, name string) error {
	record := &MigrationRecord{
		Name:       name,
		ExecutedAt: time.Now().Unix(),
	}
	return db.Create(record).Error
}

// removeMigrationRecord 删除迁移记录
func removeMigrationRecord(db *gorm.DB, name string) error {
	return db.Where("name = ?", name).Delete(&MigrationRecord{}).Error
}
