package user

import (
	"context"
	"time"

	"member-pre/pkg/errors"
	"member-pre/pkg/logger"
)

// 用户角色常量
const (
	RoleAdmin        = "admin"         // 总后台
	RoleStoreManager = "store_manager" // 店长
	RoleTechnician   = "technician"    // 美甲师
	RoleCustomer     = "customer"      // 顾客
)

// 用户状态常量
const (
	StatusActive   = "active"   // 激活
	StatusInactive = "inactive" // 禁用
)

// 美甲师工作状态常量
const (
	WorkStatusWorking = "working" // 在岗
	WorkStatusRest    = "rest"    // 休息
	WorkStatusOffline = "offline" // 离岗
)

// User 用户实体
type User struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`     // 用户名（员工使用）
	Email       string    `json:"email"`        // 邮箱
	Phone       string    `json:"phone"`        // 手机号（顾客和员工都可能有）
	Password    string    `json:"-"`            // 密码（不序列化，顾客无密码）
	Role        string    `json:"role"`         // 角色: admin, store_manager, technician, customer
	Status      string    `json:"status"`      // 状态: active, inactive
	StoreID     *uint     `json:"store_id"`     // 门店ID（店长和美甲师必须关联，顾客和总后台为nil）
	WorkStatus  *string   `json:"work_status"` // 工作状态（仅美甲师有效）: working, rest, offline
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// IUserRepository 用户仓储接口
type IUserRepository interface {
	// FindByUsername 根据用户名查找用户
	FindByUsername(ctx context.Context, username string) (*User, error)
	// FindByPhone 根据手机号查找用户
	FindByPhone(ctx context.Context, phone string) (*User, error)
	// FindByID 根据ID查找用户
	FindByID(ctx context.Context, id uint) (*User, error)
	// Create 创建用户
	Create(ctx context.Context, user *User) error
	// Update 更新用户
	Update(ctx context.Context, user *User) error
	// UpdateWorkStatus 更新美甲师工作状态
	UpdateWorkStatus(ctx context.Context, userID uint, workStatus string) error
	// FindByStoreID 根据门店ID查找用户（店长和美甲师）
	FindByStoreID(ctx context.Context, storeID uint, role string) ([]*User, error)
}

// 领域错误定义
var (
	ErrUserNotFound      = errors.ErrNotFound("用户不存在")
	ErrUserInactive      = errors.ErrForbidden("用户已被禁用")
	ErrInvalidRole       = errors.ErrInvalidParams("无效的用户角色")
	ErrStoreRequired     = errors.ErrInvalidParams("店长和美甲师必须关联门店")
	ErrInvalidWorkStatus = errors.ErrInvalidParams("无效的工作状态")
	ErrPhoneRequired     = errors.ErrInvalidParams("手机号不能为空")
)

// UserService 用户服务
type UserService struct {
	repo   IUserRepository
	logger logger.Logger
}

// NewUserService 创建用户服务
func NewUserService(repo IUserRepository, log logger.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: log,
	}
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(ctx context.Context, id uint) (*User, error) {
	s.logger.Debug("获取用户", logger.NewField("user_id", id))

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("查找用户失败", logger.NewField("user_id", id), logger.NewField("error", err.Error()))
		return nil, err
	}

	if user == nil {
		s.logger.Warn("用户不存在", logger.NewField("user_id", id))
		return nil, ErrUserNotFound
	}

	return user, nil
}

// GetByPhone 根据手机号获取用户（用于顾客登录）
func (s *UserService) GetByPhone(ctx context.Context, phone string) (*User, error) {
	s.logger.Debug("根据手机号获取用户", logger.NewField("phone", phone))

	user, err := s.repo.FindByPhone(ctx, phone)
	if err != nil {
		s.logger.Error("查找用户失败", logger.NewField("phone", phone), logger.NewField("error", err.Error()))
		return nil, err
	}

	if user == nil {
		s.logger.Warn("用户不存在", logger.NewField("phone", phone))
		return nil, ErrUserNotFound
	}

	return user, nil
}

// GetByUsername 根据用户名获取用户（用于员工登录）
func (s *UserService) GetByUsername(ctx context.Context, username string) (*User, error) {
	s.logger.Debug("根据用户名获取用户", logger.NewField("username", username))

	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		s.logger.Error("查找用户失败", logger.NewField("username", username), logger.NewField("error", err.Error()))
		return nil, err
	}

	if user == nil {
		s.logger.Warn("用户不存在", logger.NewField("username", username))
		return nil, ErrUserNotFound
	}

	return user, nil
}

// CreateOrGetCustomer 创建或获取顾客（手机号登录时自动创建）
func (s *UserService) CreateOrGetCustomer(ctx context.Context, phone string) (*User, error) {
	s.logger.Debug("创建或获取顾客", logger.NewField("phone", phone))

	// 先尝试查找
	user, err := s.repo.FindByPhone(ctx, phone)
	if err != nil {
		s.logger.Error("查找用户失败", logger.NewField("phone", phone), logger.NewField("error", err.Error()))
		return nil, err
	}

	// 如果存在，检查是否是顾客
	if user != nil {
		if user.Role != RoleCustomer {
			s.logger.Warn("手机号已被其他角色使用", logger.NewField("phone", phone), logger.NewField("role", user.Role))
			return nil, errors.ErrInvalidParams("该手机号已被注册为其他角色")
		}
		return user, nil
	}

	// 不存在则创建新顾客
	newUser := &User{
		Phone:    phone,
		Role:     RoleCustomer,
		Status:   StatusActive,
		StoreID:  nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		s.logger.Error("创建顾客失败", logger.NewField("phone", phone), logger.NewField("error", err.Error()))
		return nil, err
	}

	s.logger.Info("创建顾客成功", logger.NewField("user_id", newUser.ID), logger.NewField("phone", phone))
	return newUser, nil
}

// UpdateWorkStatus 更新美甲师工作状态
func (s *UserService) UpdateWorkStatus(ctx context.Context, userID uint, workStatus string) error {
	s.logger.Info("更新美甲师工作状态", logger.NewField("user_id", userID), logger.NewField("work_status", workStatus))

	// 验证工作状态
	if workStatus != WorkStatusWorking && workStatus != WorkStatusRest && workStatus != WorkStatusOffline {
		return ErrInvalidWorkStatus
	}

	// 获取用户，验证是否为美甲师
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	if user.Role != RoleTechnician {
		return errors.ErrInvalidParams("只有美甲师可以更新工作状态")
	}

	// 更新工作状态
	if err := s.repo.UpdateWorkStatus(ctx, userID, workStatus); err != nil {
		s.logger.Error("更新工作状态失败", logger.NewField("user_id", userID), logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("更新工作状态成功", logger.NewField("user_id", userID), logger.NewField("work_status", workStatus))
	return nil
}

// GetByStoreID 根据门店ID获取用户列表
func (s *UserService) GetByStoreID(ctx context.Context, storeID uint, role string) ([]*User, error) {
	s.logger.Debug("根据门店ID获取用户", logger.NewField("store_id", storeID), logger.NewField("role", role))

	users, err := s.repo.FindByStoreID(ctx, storeID, role)
	if err != nil {
		s.logger.Error("查找用户失败", logger.NewField("store_id", storeID), logger.NewField("error", err.Error()))
		return nil, err
	}

	return users, nil
}

// ValidateUser 验证用户状态
func (s *UserService) ValidateUser(ctx context.Context, user *User) error {
	if user == nil {
		return ErrUserNotFound
	}

	if user.Status != StatusActive {
		return ErrUserInactive
	}

	return nil
}

// HasPermission 检查用户是否有权限访问接口
// role: 用户角色
// requiredRole: 需要的角色（空字符串表示所有角色都可以）
func HasPermission(userRole string, requiredRoles ...string) bool {
	if len(requiredRoles) == 0 {
		return true
	}

	for _, role := range requiredRoles {
		if userRole == role {
			return true
		}
	}

	return false
}

// CanAccessStore 检查用户是否可以访问指定门店的数据
// 总后台可以访问所有门店，店长只能访问自己的门店，美甲师只能访问自己的门店
func CanAccessStore(user *User, storeID uint) bool {
	if user == nil {
		return false
	}

	// 总后台可以访问所有门店
	if user.Role == RoleAdmin {
		return true
	}

	// 店长和美甲师只能访问自己的门店
	if user.StoreID != nil && *user.StoreID == storeID {
		return true
	}

	return false
}

// CanAccessUser 检查用户是否可以访问指定用户的数据
// 总后台可以访问所有用户，店长可以访问本店用户，美甲师只能访问自己，顾客只能访问自己
func CanAccessUser(currentUser *User, targetUserID uint) bool {
	if currentUser == nil {
		return false
	}

	// 总后台可以访问所有用户
	if currentUser.Role == RoleAdmin {
		return true
	}

	// 用户只能访问自己的数据
	if currentUser.ID == targetUserID {
		return true
	}

	// 店长可以访问本店的美甲师和顾客（通过预约关联）
	// 这里简化处理，实际应该通过预约等关联关系判断
	if currentUser.Role == RoleStoreManager {
		// 需要查询目标用户是否属于当前门店，这里简化处理
		return true // 实际应该查询数据库
	}

	return false
}
