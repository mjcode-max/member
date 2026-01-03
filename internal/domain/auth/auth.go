package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"member-pre/internal/domain/user"
	pkgErrors "member-pre/pkg/errors"
	"member-pre/pkg/logger"
)

// IAuthConfig 认证配置接口
type IAuthConfig interface {
	GetJwtSecret() string
	GetTokenExpires() int
}

// IUserService 用户服务接口（用于认证模块调用用户服务）
type IUserService interface {
	GetByID(ctx context.Context, id uint) (*user.User, error)
	GetByPhone(ctx context.Context, phone string) (*user.User, error)
	GetByUsername(ctx context.Context, username string) (*user.User, error)
	CreateOrGetCustomer(ctx context.Context, phone string) (*user.User, error)
	ValidateUser(ctx context.Context, u *user.User) error
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名（员工）或手机号（顾客）
	Password string `json:"password"`                    // 密码（员工必填，顾客不需要）
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string         `json:"token"`
	User      *user.User     `json:"user"`
	ExpiresAt time.Time      `json:"expires_at"`
}

// GetCurrentUserResponse 获取当前用户响应
type GetCurrentUserResponse struct {
	User *user.User `json:"user"`
}

// TokenClaims JWT令牌声明
type TokenClaims struct {
	UserID   uint   `json:"user_id"`
	Role     string `json:"role"`
	StoreID  *uint  `json:"store_id,omitempty"` // 店长和美甲师的门店ID
	jwt.RegisteredClaims
}

// 领域错误定义
var (
	ErrUserNotFound    = pkgErrors.ErrNotFound("用户不存在")
	ErrInvalidPassword = pkgErrors.ErrUnauthorized("密码错误")
	ErrUserInactive    = pkgErrors.ErrForbidden("用户已被禁用")
	ErrInvalidToken    = pkgErrors.ErrUnauthorized("无效的令牌")
	ErrTokenExpired    = pkgErrors.ErrUnauthorized("令牌已过期")
)

// AuthService 认证服务
type AuthService struct {
	userService IUserService
	logger      logger.Logger
	authConfig  IAuthConfig
}

// NewAuthService 创建认证服务
func NewAuthService(userService IUserService, log logger.Logger, authConfig IAuthConfig) *AuthService {
	return &AuthService{
		userService: userService,
		logger:      log,
		authConfig:  authConfig,
	}
}

// Login 登录
// 顾客：使用手机号，无需密码，自动创建账户
// 员工：使用用户名+密码
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	s.logger.Info("用户登录尝试", logger.NewField("username", req.Username))

	var u *user.User
	var err error

	// 判断是手机号还是用户名
	if isPhoneNumber(req.Username) {
		// 手机号登录（顾客）
		if req.Password != "" {
			s.logger.Warn("顾客登录不应提供密码", logger.NewField("phone", req.Username))
			// 不返回错误，允许继续处理
		}

		// 顾客使用手机号直接登录，自动创建账户
		u, err = s.userService.CreateOrGetCustomer(ctx, req.Username)
		if err != nil {
			s.logger.Error("创建或获取顾客失败", logger.NewField("phone", req.Username), logger.NewField("error", err.Error()))
			return nil, err
		}
	} else {
		// 用户名登录（员工：店长、美甲师、总后台）
		if req.Password == "" {
			s.logger.Warn("员工登录必须提供密码", logger.NewField("username", req.Username))
			return nil, ErrInvalidPassword
		}

		u, err = s.userService.GetByUsername(ctx, req.Username)
		if err != nil {
			s.logger.Error("查找用户失败", logger.NewField("username", req.Username), logger.NewField("error", err.Error()))
			return nil, err
		}

		if u == nil {
			s.logger.Warn("用户不存在", logger.NewField("username", req.Username))
			return nil, ErrUserNotFound
		}

		// 验证密码
		if err := s.verifyPassword(u.Password, req.Password); err != nil {
			s.logger.Warn("密码错误", logger.NewField("user_id", u.ID))
			return nil, ErrInvalidPassword
		}
	}

	// 验证用户状态
	if err := s.userService.ValidateUser(ctx, u); err != nil {
		s.logger.Warn("用户状态验证失败", logger.NewField("user_id", u.ID), logger.NewField("error", err.Error()))
		return nil, err
	}

	// 生成JWT令牌
	token, expiresAt, err := s.generateToken(u)
	if err != nil {
		s.logger.Error("生成令牌失败", logger.NewField("user_id", u.ID), logger.NewField("error", err.Error()))
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	s.logger.Info("用户登录成功", logger.NewField("user_id", u.ID), logger.NewField("role", u.Role))

	return &LoginResponse{
		Token:     token,
		User:      u,
		ExpiresAt: expiresAt,
	}, nil
}

// GetCurrentUser 获取当前用户
func (s *AuthService) GetCurrentUser(ctx context.Context, userID uint) (*user.User, error) {
	s.logger.Debug("获取当前用户", logger.NewField("user_id", userID))

	u, err := s.userService.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error("查找用户失败", logger.NewField("user_id", userID), logger.NewField("error", err.Error()))
		return nil, err
	}

	return u, nil
}

// ValidateToken 验证JWT令牌
func (s *AuthService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(s.authConfig.GetJwtSecret()), nil
	})

	if err != nil {
		// 检查是否是JWT库的过期错误
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// GenerateTokenForUser 为用户生成令牌（用于刷新令牌等场景）
func (s *AuthService) GenerateTokenForUser(u *user.User) (string, time.Time, error) {
	return s.generateToken(u)
}

// Logout 登出（将token加入黑名单，这里简化处理）
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

// generateToken 生成JWT令牌
func (s *AuthService) generateToken(u *user.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(s.authConfig.GetTokenExpires()) * time.Second)

	claims := &TokenClaims{
		UserID:  u.ID,
		Role:     u.Role,
		StoreID:  u.StoreID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.authConfig.GetJwtSecret()))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// verifyPassword 验证密码
func (s *AuthService) verifyPassword(hashedPassword, password string) error {
	// 如果密码为空，说明是旧数据或未加密，直接比较（兼容性处理）
	if hashedPassword == "" {
		return ErrInvalidPassword
	}

	// 尝试bcrypt验证
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == nil {
		return nil
	}

	// 如果bcrypt验证失败，可能是明文密码（兼容旧数据）
	if hashedPassword == password {
		return nil
	}

	return ErrInvalidPassword
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// isPhoneNumber 判断是否为手机号（简单判断：11位数字）
func isPhoneNumber(s string) bool {
	if len(s) != 11 {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
