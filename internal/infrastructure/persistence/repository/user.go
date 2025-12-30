package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"member-pre/internal/infrastructure/persistence/model"
	"member-pre/internal/infrastructure/persistence/mysql"
	"member-pre/internal/infrastructure/persistence/redis"
)

// UserRepository 用户仓储接口（定义在 repository 层，由 domain 层引用）
type UserRepository interface {
	// FindByUsername 根据用户名查找用户
	FindByUsername(username string) (*model.User, error)

	// FindByID 根据ID查找用户
	FindByID(id uint) (*model.User, error)

	// Create 创建用户
	Create(user *model.User) error

	// Update 更新用户
	Update(user *model.User) error

	// SaveToken 保存token到Redis
	SaveToken(userID uint, token string, expiresIn int64) error

	// DeleteToken 删除token
	DeleteToken(token string) error

	// ValidateToken 验证token是否有效
	ValidateToken(token string) (uint, error)
}

// UserModel 用户模型（数据库表结构）
type UserModel struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;not null;size:50"`
	Password  string `gorm:"not null;size:255"`
	Email     string `gorm:"size:100"`
	Phone     string `gorm:"size:20"`
	Role      string `gorm:"size:20;not null;default:'customer'"` // admin, staff, store, customer
	Status    int    `gorm:"not null;default:1"`                  // 1-正常，0-禁用
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 指定表名
func (UserModel) TableName() string {
	return "users"
}

// ToEntity 转换为领域实体
func (m *UserModel) ToEntity() *model.User {
	return &model.User{
		ID:        m.ID,
		Username:  m.Username,
		Password:  m.Password,
		Email:     m.Email,
		Phone:     m.Phone,
		Role:      m.Role,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromEntity 从领域实体创建模型
func (m *UserModel) FromEntity(u *model.User) {
	m.ID = u.ID
	m.Username = u.Username
	m.Password = u.Password
	m.Email = u.Email
	m.Phone = u.Phone
	m.Role = u.Role
	m.Status = u.Status
	m.CreatedAt = u.CreatedAt
	m.UpdatedAt = u.UpdatedAt
}

// userRepository 用户仓储实现
type userRepository struct {
	db    *mysql.DB
	redis *redis.Client
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *mysql.DB, rdb *redis.Client) UserRepository {
	return &userRepository{
		db:    db,
		redis: rdb,
	}
}

// FindByUsername 根据用户名查找用户
func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var model UserModel
	if err := r.db.Where("username = ?", username).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// FindByID 根据ID查找用户
func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var model UserModel
	if err := r.db.First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) error {
	model := &UserModel{}
	model.FromEntity(user)
	return r.db.Create(model).Error
}

// Update 更新用户
func (r *userRepository) Update(user *model.User) error {
	model := &UserModel{}
	model.FromEntity(user)
	return r.db.Save(model).Error
}

// SaveToken 保存token到Redis
func (r *userRepository) SaveToken(userID uint, token string, expiresIn int64) error {
	ctx := context.Background()
	key := fmt.Sprintf("token:%s", token)
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	// 保存token，key为token，value为userID
	if err := r.redis.Set(ctx, key, userIDStr, time.Duration(expiresIn)*time.Second).Err(); err != nil {
		return err
	}

	// 同时保存用户的所有token列表（用于登出时清理）
	userTokensKey := fmt.Sprintf("user:tokens:%d", userID)
	return r.redis.SAdd(ctx, userTokensKey, token).Err()
}

// DeleteToken 删除token
func (r *userRepository) DeleteToken(token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("token:%s", token)

	// 获取userID
	userIDStr, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		// token可能已经过期，直接返回成功
		return nil
	}

	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	// 删除token
	r.redis.Del(ctx, key)

	// 从用户token列表中移除
	userTokensKey := fmt.Sprintf("user:tokens:%d", userID)
	r.redis.SRem(ctx, userTokensKey, token)

	return nil
}

// ValidateToken 验证token是否有效
func (r *userRepository) ValidateToken(token string) (uint, error) {
	ctx := context.Background()
	key := fmt.Sprintf("token:%s", token)

	userIDStr, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("token无效或已过期")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("无效的token")
	}

	return uint(userID), nil
}
