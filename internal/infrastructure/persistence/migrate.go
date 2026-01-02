package persistence

import (
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
	"member-pre/pkg/logger"
)

// models 所有需要迁移的模型
var models []interface{}

// Register 注册需要迁移的模型
func Register(model interface{}) {
	models = append(models, model)
}

// RegisterModels 批量注册需要迁移的模型
func RegisterModels(modelList ...interface{}) {
	models = append(models, modelList...)
}

// getModelName 获取模型名称
func getModelName(model interface{}) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

// Up 执行所有迁移
func Up(db *gorm.DB, log logger.Logger) error {
	// 创建迁移记录表
	if err := createMigrationTable(db); err != nil {
		return fmt.Errorf("创建迁移记录表失败: %w", err)
	}

	for _, model := range models {
		modelName := getModelName(model)

		// 检查是否已执行
		executed, err := isMigrationExecuted(db, modelName)
		if err != nil {
			return fmt.Errorf("检查迁移状态失败: %w", err)
		}

		if executed {
			log.Info("迁移已执行，跳过", logger.NewField("model", modelName))
			continue
		}

		// 执行迁移（使用 GORM 的 AutoMigrate）
		log.Info("开始执行迁移", logger.NewField("model", modelName))
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("执行迁移失败 [%s]: %w", modelName, err)
		}

		// 记录迁移
		if err := recordMigration(db, modelName); err != nil {
			return fmt.Errorf("记录迁移失败: %w", err)
		}

		log.Info("迁移执行成功", logger.NewField("model", modelName))
	}

	return nil
}

// Down 回滚最后一次迁移
func Down(db *gorm.DB, log logger.Logger) error {
	if len(models) == 0 {
		log.Info("没有可回滚的迁移")
		return nil
	}

	// 获取最后一次迁移的模型
	lastModel := models[len(models)-1]
	modelName := getModelName(lastModel)

	// 检查是否已执行
	executed, err := isMigrationExecuted(db, modelName)
	if err != nil {
		return fmt.Errorf("检查迁移状态失败: %w", err)
	}

	if !executed {
		log.Info("迁移未执行，无需回滚", logger.NewField("model", modelName))
		return nil
	}

	// 执行回滚（删除表）
	log.Info("开始回滚迁移", logger.NewField("model", modelName))
	if err := db.Migrator().DropTable(lastModel); err != nil {
		return fmt.Errorf("回滚迁移失败 [%s]: %w", modelName, err)
	}

	// 删除迁移记录
	if err := removeMigrationRecord(db, modelName); err != nil {
		return fmt.Errorf("删除迁移记录失败: %w", err)
	}

	log.Info("迁移回滚成功", logger.NewField("model", modelName))
	return nil
}

// Status 查看迁移状态
func Status(db *gorm.DB, log logger.Logger) error {
	// 创建迁移记录表（如果不存在）
	if err := createMigrationTable(db); err != nil {
		return fmt.Errorf("创建迁移记录表失败: %w", err)
	}

	log.Info("迁移状态:")
	for _, model := range models {
		modelName := getModelName(model)
		executed, err := isMigrationExecuted(db, modelName)
		if err != nil {
			log.Error("检查迁移状态失败", logger.NewField("model", modelName), logger.NewField("error", err.Error()))
			continue
		}

		status := "未执行"
		if executed {
			status = "已执行"
		}
		log.Info(fmt.Sprintf("  %s: %s", modelName, status))
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
