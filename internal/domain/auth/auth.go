package auth

import (
	"context"
	"time"

	"member-pre/pkg/errors"
	"member-pre/pkg/logger"
)

type IAuthRepository interface {
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
}

// User 用户实体
type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"-"`      // 不序列化密码
	Role      string    `json:"role"`   // 角色: admin, staff, customer, store
	Status    string    `json:"status"` // 状态: active, inactive
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名或手机号
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string    `json:"token"`
	User      *User     `json:"user"`
	ExpiresAt time.Time `json:"expires_at"`
}

// GetCurrentUserResponse 获取当前用户响应
type GetCurrentUserResponse struct {
	User *User `json:"user"`
}

// 领域错误定义
var (
	ErrUserNotFound    = errors.ErrNotFound("用户不存在")
	ErrInvalidPassword = errors.ErrUnauthorized("密码错误")
	ErrUserInactive    = errors.ErrForbidden("用户已被禁用")
)

// AuthService 认证服务
type AuthService struct {
	repo         IAuthRepository
	logger       logger.Logger
	jwtSecret    string
	tokenExpires int // 秒
}

// NewAuthService 创建认证服务
// jwtSecret: JWT密钥
// tokenExpires: token过期时间（秒）
func NewAuthService(repo IAuthRepository, log logger.Logger, jwtSecret string, tokenExpires int) *AuthService {
	return &AuthService{
		repo:         repo,
		logger:       log,
		jwtSecret:    jwtSecret,
		tokenExpires: tokenExpires,
	}
}

// Login 登录
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	s.logger.Info("用户登录尝试", logger.NewField("username", req.Username))

	// 根据用户名或手机号查找用户
	var user *User
	var err error

	// 判断是用户名还是手机号（简化处理，实际可以根据格式判断）
	if len(req.Username) == 11 && isNumeric(req.Username) {
		user, err = s.repo.FindByPhone(ctx, req.Username)
	} else {
		user, err = s.repo.FindByUsername(ctx, req.Username)
	}

	if err != nil {
		s.logger.Error("查找用户失败", logger.NewField("username", req.Username), logger.NewField("error", err.Error()))
		return nil, err
	}

	if user == nil {
		s.logger.Warn("用户不存在", logger.NewField("username", req.Username))
		return nil, ErrUserNotFound
	}

	// 验证密码（简化处理，实际应该使用 bcrypt 等）
	if user.Password != req.Password {
		s.logger.Warn("密码错误", logger.NewField("user_id", user.ID))
		return nil, ErrInvalidPassword
	}

	// 检查用户状态
	if user.Status != "active" {
		s.logger.Warn("用户已被禁用", logger.NewField("user_id", user.ID))
		return nil, ErrUserInactive
	}

	// 生成 Token（简化处理，实际应该使用 JWT）
	token := generateToken(user.ID, s.jwtSecret)
	expiresAt := time.Now().Add(time.Duration(s.tokenExpires) * time.Second)

	s.logger.Info("用户登录成功", logger.NewField("user_id", user.ID), logger.NewField("role", user.Role))

	return &LoginResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

// GetCurrentUser 获取当前用户
func (s *AuthService) GetCurrentUser(ctx context.Context, userID uint) (*User, error) {
	s.logger.Debug("获取当前用户", logger.NewField("user_id", userID))

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error("查找用户失败", logger.NewField("user_id", userID), logger.NewField("error", err.Error()))
		return nil, err
	}

	if user == nil {
		s.logger.Warn("用户不存在", logger.NewField("user_id", userID))
		return nil, ErrUserNotFound
	}

	return user, nil
}

// Logout 登出（简化处理，实际应该使 Token 失效）
func (s *AuthService) Logout(ctx context.Context, token string) error {
	// 只记录token前20个字符（如果token长度足够）
	tokenPreview := token
	if len(token) > 20 {
		tokenPreview = token[:20]
	}
	s.logger.Info("用户登出", logger.NewField("token_preview", tokenPreview))
	// 实际实现中应该将 token 加入黑名单或从 Redis 中删除
	return nil
}

// 辅助函数

// isNumeric 判断字符串是否为纯数字
func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// generateToken 生成 Token（简化实现，实际应该使用 JWT）
func generateToken(userID uint, secret string) string {
	// 简化实现，实际应该使用 JWT 库生成
	return "token_" + string(rune(userID)) + "_" + secret
}
