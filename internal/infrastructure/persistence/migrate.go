package persistence

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
	"member-pre/internal/domain/auth"
	"member-pre/internal/domain/user"
	"member-pre/internal/infrastructure/config"
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
func Up(db *gorm.DB, log logger.Logger, cfg *config.Config) error {
	// 创建迁移记录表
	if err := createMigrationTable(db); err != nil {
		return fmt.Errorf("创建迁移记录表失败: %w", err)
	}

	// 首先，对所有模型执行 AutoMigrate，确保表结构和字段都是最新的
	// 这样可以检测并添加新字段，即使迁移记录已存在
	log.Info("执行 AutoMigrate 更新所有表结构")
	for _, model := range models {
		modelName := getModelName(model)
		log.Debug("执行 AutoMigrate", logger.NewField("model", modelName))

		if err := db.AutoMigrate(model); err != nil {
			// 检查是否是删除不存在索引的错误，如果是则忽略（这是 GORM 的已知问题）
			if isIndexDropError(err) {
				log.Info("AutoMigrate 更新表结构成功（忽略索引删除错误）", logger.NewField("model", modelName))
			} else {
				log.Warn("AutoMigrate 更新表结构失败", logger.NewField("model", modelName), logger.NewField("error", err.Error()))
				// 不返回错误，继续执行其他迁移
			}
		} else {
			log.Debug("AutoMigrate 更新表结构成功", logger.NewField("model", modelName))
		}
	}

	// 然后，检查并记录首次迁移
	for _, model := range models {
		modelName := getModelName(model)

		// 检查是否已记录首次迁移
		executed, err := isMigrationExecuted(db, modelName)
		if err != nil {
			return fmt.Errorf("检查迁移状态失败: %w", err)
		}

		if !executed {
			// 首次迁移，记录迁移记录
			log.Info("记录首次迁移", logger.NewField("model", modelName))
			if err := recordMigration(db, modelName); err != nil {
				return fmt.Errorf("记录迁移失败: %w", err)
			}
			log.Info("首次迁移记录成功", logger.NewField("model", modelName))
		}
	}

	// 迁移完成后，初始化默认 admin 账户
	if err := initAdminUser(db, log, cfg); err != nil {
		log.Warn("初始化 admin 账户失败", logger.NewField("error", err.Error()))
		// 不返回错误，因为这不是迁移失败
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
	Name       string `gorm:"uniqueIndex;not null;size:255"` // 指定size避免MySQL TEXT类型索引错误
	ExecutedAt int64  `gorm:"not null"`
}

// createMigrationTable 创建迁移记录表
func createMigrationTable(db *gorm.DB) error {
	migrator := db.Migrator()

	// 如果表不存在，直接创建
	if !migrator.HasTable(&MigrationRecord{}) {
		if err := migrator.CreateTable(&MigrationRecord{}); err != nil {
			return fmt.Errorf("创建迁移记录表失败: %w", err)
		}
		return nil
	}

	// 如果表已存在，使用 AutoMigrate 更新结构
	// 如果遇到删除不存在索引的错误，可以忽略（这是 GORM 的已知问题）
	if err := db.AutoMigrate(&MigrationRecord{}); err != nil {
		// 检查是否是删除不存在索引的错误，如果是则忽略
		if isIndexDropError(err) {
			// 忽略此错误，表结构已经正确
			return nil
		}
		return fmt.Errorf("更新迁移记录表失败: %w", err)
	}

	return nil
}

// isIndexDropError 检查错误是否是删除不存在的索引
func isIndexDropError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "Can't DROP") && strings.Contains(errStr, "check that column/key exists")
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

// UserModel 用户数据库模型（用于迁移时创建 admin 账户，避免循环导入）
type UserModel struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Username   string         `gorm:"uniqueIndex;size:50" json:"username"`
	Email      string         `gorm:"size:100" json:"email"`
	Phone      string         `gorm:"uniqueIndex;size:20" json:"phone"`
	Password   string         `gorm:"size:255" json:"-"`
	Role       string         `gorm:"size:20;default:'customer';not null" json:"role"`
	Status     string         `gorm:"size:20;default:'active';not null" json:"status"`
	StoreID    *uint          `gorm:"index" json:"store_id"`
	WorkStatus *string        `gorm:"size:20" json:"work_status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (UserModel) TableName() string {
	return "users"
}

// initAdminUser 初始化默认 admin 账户
func initAdminUser(db *gorm.DB, log logger.Logger, cfg *config.Config) error {
	if cfg == nil || cfg.Admin.Username == "" {
		log.Info("未配置 admin 账户，跳过初始化")
		return nil
	}

	// 检查 admin 账户是否已存在
	var count int64
	if err := db.Model(&UserModel{}).Where("username = ? AND role = ?", cfg.Admin.Username, user.RoleAdmin).Count(&count).Error; err != nil {
		return fmt.Errorf("检查 admin 账户失败: %w", err)
	}

	if count > 0 {
		log.Info("admin 账户已存在，跳过创建", logger.NewField("username", cfg.Admin.Username))
		return nil
	}

	// 加密密码
	hashedPassword, err := auth.HashPassword(cfg.Admin.Password)
	if err != nil {
		return fmt.Errorf("加密密码失败: %w", err)
	}

	// 创建 admin 账户
	adminUser := &UserModel{
		Username:  cfg.Admin.Username,
		Password:  hashedPassword,
		Role:      user.RoleAdmin,
		Status:    user.StatusActive,
		StoreID:   nil, // admin 不关联门店
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(adminUser).Error; err != nil {
		return fmt.Errorf("创建 admin 账户失败: %w", err)
	}

	log.Info("admin 账户创建成功", logger.NewField("username", cfg.Admin.Username))
	return nil
}
