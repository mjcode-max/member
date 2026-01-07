package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"member-pre/internal/domain/user"
	"member-pre/internal/infrastructure/persistence"
	"member-pre/internal/infrastructure/persistence/database"
	"member-pre/pkg/errors"
	"member-pre/pkg/logger"
)

var _ user.IUserRepository = (*UserRepository)(nil)

// 在包初始化时注册模型，确保迁移时能检测到
func init() {
	persistence.Register(&UserModel{})
	persistence.Register(&CustomerOpenIDModel{})
}

// UserRepository 用户仓储实现
type UserRepository struct {
	db     database.Database
	redis  *persistence.Client
	logger logger.Logger
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db database.Database, rdb *persistence.Client, log logger.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		redis:  rdb,
		logger: log,
	}
}

// UserModel 用户数据库模型
type UserModel struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Username   string         `gorm:"uniqueIndex;size:50" json:"username"`
	Email      string         `gorm:"size:100" json:"email"`
	Phone      string         `gorm:"uniqueIndex;size:20" json:"phone"`
	Password   string         `gorm:"size:255" json:"-"`
	Role       string         `gorm:"size:20;default:'customer';not null" json:"role"`
	Status     string         `gorm:"size:20;default:'active';not null" json:"status"`
	StoreID    *uint          `gorm:"index" json:"store_id"`      // 店长和美甲师必须关联门店
	WorkStatus *string        `gorm:"size:20" json:"work_status"` // 美甲师工作状态
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (UserModel) TableName() string {
	return "users"
}

// CustomerOpenIDModel 顾客微信OpenID数据库模型
type CustomerOpenIDModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OpenID    string    `gorm:"uniqueIndex;size:100;not null" json:"openid"` // 微信OpenID
	Phone     string    `gorm:"size:20" json:"phone"`                        // 手机号
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (CustomerOpenIDModel) TableName() string {
	return "customer_openids"
}

// ToEntity 转换为领域实体
func (m *UserModel) ToEntity() *user.User {
	if m == nil {
		return nil
	}
	return &user.User{
		ID:         m.ID,
		Username:   m.Username,
		Email:      m.Email,
		Phone:      m.Phone,
		Password:   m.Password,
		Role:       m.Role,
		Status:     m.Status,
		StoreID:    m.StoreID,
		WorkStatus: m.WorkStatus,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

// FromEntity 从领域实体创建模型
func (m *UserModel) FromEntity(u *user.User) {
	if u == nil {
		return
	}
	m.ID = u.ID
	m.Username = u.Username
	m.Email = u.Email
	m.Phone = u.Phone
	m.Password = u.Password
	m.Role = u.Role
	m.Status = u.Status
	m.StoreID = u.StoreID
	m.WorkStatus = u.WorkStatus
	m.CreatedAt = u.CreatedAt
	m.UpdatedAt = u.UpdatedAt
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
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
func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*user.User, error) {
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
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*user.User, error) {
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
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	r.logger.Info("创建用户",
		logger.NewField("username", u.Username),
		logger.NewField("phone", u.Phone),
		logger.NewField("role", u.Role),
	)

	model := &UserModel{}
	model.FromEntity(u)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Error("创建用户失败",
			logger.NewField("username", u.Username),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	// 更新实体的ID和时间戳
	u.ID = model.ID
	u.CreatedAt = model.CreatedAt
	u.UpdatedAt = model.UpdatedAt

	r.logger.Info("创建用户成功",
		logger.NewField("user_id", u.ID),
		logger.NewField("username", u.Username),
	)

	return nil
}

// Update 更新用户
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	r.logger.Info("更新用户",
		logger.NewField("user_id", u.ID),
	)

	model := &UserModel{}
	model.FromEntity(u)
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Model(&UserModel{}).Where("id = ?", u.ID).Updates(model).Error; err != nil {
		r.logger.Error("更新用户失败",
			logger.NewField("user_id", u.ID),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	u.UpdatedAt = model.UpdatedAt

	r.logger.Info("更新用户成功",
		logger.NewField("user_id", u.ID),
	)

	return nil
}

// UpdateWorkStatus 更新美甲师工作状态
func (r *UserRepository) UpdateWorkStatus(ctx context.Context, userID uint, workStatus string) error {
	r.logger.Info("更新美甲师工作状态",
		logger.NewField("user_id", userID),
		logger.NewField("work_status", workStatus),
	)

	if err := r.db.DB().WithContext(ctx).Model(&UserModel{}).
		Where("id = ? AND role = ?", userID, user.RoleTechnician).
		Update("work_status", workStatus).Error; err != nil {
		r.logger.Error("更新工作状态失败",
			logger.NewField("user_id", userID),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	r.logger.Info("更新工作状态成功",
		logger.NewField("user_id", userID),
		logger.NewField("work_status", workStatus),
	)

	return nil
}

// FindList 获取用户列表（支持筛选和分页）
func (r *UserRepository) FindList(ctx context.Context, role, status string, storeID *uint, username, phone string, page, pageSize int) ([]*user.User, int64, error) {
	r.logger.Debug("查找用户列表",
		logger.NewField("role", role),
		logger.NewField("status", status),
		logger.NewField("store_id", storeID),
		logger.NewField("username", username),
		logger.NewField("phone", phone),
		logger.NewField("page", page),
		logger.NewField("page_size", pageSize),
	)

	var models []UserModel
	var total int64

	query := r.db.DB().WithContext(ctx).Model(&UserModel{})

	// 筛选条件
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("获取用户总数失败",
			logger.NewField("error", err.Error()),
		)
		return nil, 0, errors.ErrDatabase(err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&models).Error; err != nil {
		r.logger.Error("查找用户列表失败",
			logger.NewField("error", err.Error()),
		)
		return nil, 0, errors.ErrDatabase(err)
	}

	users := make([]*user.User, 0, len(models))
	for _, model := range models {
		users = append(users, model.ToEntity())
	}

	r.logger.Debug("查找用户列表成功",
		logger.NewField("count", len(users)),
		logger.NewField("total", total),
	)

	return users, total, nil
}

// FindByStoreID 根据门店ID查找用户（店长和美甲师）
func (r *UserRepository) FindByStoreID(ctx context.Context, storeID uint, role string) ([]*user.User, error) {
	r.logger.Debug("查找用户：根据门店ID",
		logger.NewField("store_id", storeID),
		logger.NewField("role", role),
	)

	var models []UserModel
	query := r.db.DB().WithContext(ctx).Where("store_id = ?", storeID)
	if role != "" {
		query = query.Where("role = ?", role)
	}

	if err := query.Find(&models).Error; err != nil {
		r.logger.Error("查找用户失败：根据门店ID",
			logger.NewField("store_id", storeID),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	users := make([]*user.User, 0, len(models))
	for _, model := range models {
		users = append(users, model.ToEntity())
	}

	r.logger.Debug("查找用户成功：根据门店ID",
		logger.NewField("store_id", storeID),
		logger.NewField("count", len(users)),
	)

	return users, nil
}

// SaveCustomerOpenID 保存顾客微信OpenID和手机号
func (r *UserRepository) SaveCustomerOpenID(ctx context.Context, openID, phone string) error {
	r.logger.Info("保存顾客微信OpenID和手机号",
		logger.NewField("openid", openID),
		logger.NewField("phone", phone),
	)

	// 检查是否已存在
	var existing CustomerOpenIDModel
	err := r.db.DB().WithContext(ctx).Where("open_id = ?", openID).First(&existing).Error
	if err == nil {
		// 已存在，更新手机号和更新时间
		existing.Phone = phone
		existing.UpdatedAt = time.Now()
		if err := r.db.DB().WithContext(ctx).Save(&existing).Error; err != nil {
			r.logger.Error("更新顾客OpenID失败",
				logger.NewField("openid", openID),
				logger.NewField("phone", phone),
				logger.NewField("error", err.Error()),
			)
			return errors.ErrDatabase(err)
		}
		r.logger.Info("更新顾客OpenID成功",
			logger.NewField("openid", openID),
			logger.NewField("phone", phone),
		)
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		r.logger.Error("查询顾客OpenID失败",
			logger.NewField("openid", openID),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	// 不存在则创建
	model := &CustomerOpenIDModel{
		OpenID:    openID,
		Phone:     phone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.db.DB().WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Error("保存顾客OpenID失败",
			logger.NewField("openid", openID),
			logger.NewField("phone", phone),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	r.logger.Info("保存顾客OpenID成功",
		logger.NewField("openid", openID),
		logger.NewField("phone", phone),
		logger.NewField("id", model.ID),
	)
	return nil
}
