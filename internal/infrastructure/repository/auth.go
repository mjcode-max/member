package repository

import (
	"context"
	"gorm.io/gorm"
	"member-pre/internal/domain/auth"
	"member-pre/internal/infrastructure/persistence"
	"member-pre/internal/infrastructure/persistence/database"
	"member-pre/pkg/errors"
	"member-pre/pkg/logger"
	"time"
)

// AuthRepository 用户仓储实现
type AuthRepository struct {
	db     database.Database
	redis  *persistence.Client
	logger logger.Logger
}

// NewAuthRepository 创建用户仓储实例
// 返回 auth 层定义的 Repository 接口
func NewAuthRepository(db database.Database, rdb *persistence.Client, log logger.Logger) *AuthRepository {
	return &AuthRepository{
		db:     db,
		redis:  rdb,
		logger: log,
	}
}

// UserModel 用户数据库模型
type UserModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email     string         `gorm:"size:100" json:"email"`
	Phone     string         `gorm:"uniqueIndex;size:20" json:"phone"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Role      string         `gorm:"size:20;default:'customer'" json:"role"`
	Status    string         `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (UserModel) TableName() string {
	return "users"
}

// ToEntity 转换为领域实体
func (m *UserModel) ToEntity() *auth.User {
	if m == nil {
		return nil
	}
	return &auth.User{
		ID:        m.ID,
		Username:  m.Username,
		Email:     m.Email,
		Phone:     m.Phone,
		Password:  m.Password,
		Role:      m.Role,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromEntity 从领域实体创建模型
func (m *UserModel) FromEntity(user *auth.User) {
	if user == nil {
		return
	}
	m.ID = user.ID
	m.Username = user.Username
	m.Email = user.Email
	m.Phone = user.Phone
	m.Password = user.Password
	m.Role = user.Role
	m.Status = user.Status
	m.CreatedAt = user.CreatedAt
	m.UpdatedAt = user.UpdatedAt
}

// FindByUsername 根据用户名查找用户
func (r *AuthRepository) FindByUsername(ctx context.Context, username string) (*auth.User, error) {
	r.logger.Debug("查找用户：根据用户名",
		logger.NewField("username", username),
	)

	var model UserModel
	if err := r.db.DB().WithContext(ctx).Where("username = ?", username).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("用户不存在：根据用户名",
				logger.NewField("username", username),
			)
			return nil, nil
		}
		r.logger.Error("查找用户失败：根据用户名",
			logger.NewField("username", username),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	r.logger.Debug("查找用户成功：根据用户名",
		logger.NewField("username", username),
		logger.NewField("user_id", model.ID),
	)

	return model.ToEntity(), nil
}

// FindByPhone 根据手机号查找用户
func (r *AuthRepository) FindByPhone(ctx context.Context, phone string) (*auth.User, error) {
	r.logger.Debug("查找用户：根据手机号",
		logger.NewField("phone", phone),
	)

	var model UserModel
	if err := r.db.DB().WithContext(ctx).Where("phone = ?", phone).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("用户不存在：根据手机号",
				logger.NewField("phone", phone),
			)
			return nil, nil
		}
		r.logger.Error("查找用户失败：根据手机号",
			logger.NewField("phone", phone),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	r.logger.Debug("查找用户成功：根据手机号",
		logger.NewField("phone", phone),
		logger.NewField("user_id", model.ID),
	)

	return model.ToEntity(), nil
}

// FindByID 根据ID查找用户
func (r *AuthRepository) FindByID(ctx context.Context, id uint) (*auth.User, error) {
	r.logger.Debug("查找用户：根据ID",
		logger.NewField("user_id", id),
	)

	var model UserModel
	if err := r.db.DB().WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("用户不存在：根据ID",
				logger.NewField("user_id", id),
			)
			return nil, nil
		}
		r.logger.Error("查找用户失败：根据ID",
			logger.NewField("user_id", id),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	r.logger.Debug("查找用户成功：根据ID",
		logger.NewField("user_id", id),
	)

	return model.ToEntity(), nil
}

// Create 创建用户
func (r *AuthRepository) Create(ctx context.Context, user *auth.User) error {
	r.logger.Info("创建用户",
		logger.NewField("username", user.Username),
		logger.NewField("phone", user.Phone),
	)

	model := &UserModel{}
	model.FromEntity(user)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Error("创建用户失败",
			logger.NewField("username", user.Username),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	// 更新实体的ID和时间戳
	user.ID = model.ID
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt

	r.logger.Info("创建用户成功",
		logger.NewField("user_id", user.ID),
		logger.NewField("username", user.Username),
	)

	return nil
}

// Update 更新用户
func (r *AuthRepository) Update(ctx context.Context, user *auth.User) error {
	r.logger.Info("更新用户",
		logger.NewField("user_id", user.ID),
	)

	model := &UserModel{}
	model.FromEntity(user)
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Model(&UserModel{}).Where("id = ?", user.ID).Updates(model).Error; err != nil {
		r.logger.Error("更新用户失败",
			logger.NewField("user_id", user.ID),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	user.UpdatedAt = model.UpdatedAt

	r.logger.Info("更新用户成功",
		logger.NewField("user_id", user.ID),
	)

	return nil
}
