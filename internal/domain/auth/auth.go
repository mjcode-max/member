package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"member-pre/internal/infrastructure/persistence/model"
	"member-pre/internal/infrastructure/persistence/repository"
)

// TokenInfo Token信息
type TokenInfo struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// Claims JWT Claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Service 认证服务
type Service struct {
	repo         repository.UserRepository
	jwtSecret    string
	tokenExpires time.Duration
}

// NewService 创建认证服务
func NewService(repo repository.UserRepository, jwtSecret string, tokenExpires time.Duration) *Service {
	return &Service{
		repo:         repo,
		jwtSecret:    jwtSecret,
		tokenExpires: tokenExpires,
	}
}

// Login 用户登录
func (s *Service) Login(req *LoginRequest) (*LoginResponse, error) {
	// 查找用户
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成token
	token, _, err := s.generateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("生成token失败")
	}

	// 保存token到Redis
	expiresIn := int64(s.tokenExpires.Seconds())
	if err := s.repo.SaveToken(user.ID, token, expiresIn); err != nil {
		return nil, errors.New("保存token失败")
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// Logout 用户登出
func (s *Service) Logout(token string) error {
	return s.repo.DeleteToken(token)
}

// ValidateToken 验证token
func (s *Service) ValidateToken(token string) (*model.User, error) {
	// 先验证JWT token
	claims, err := s.ParseToken(token)
	if err != nil {
		return nil, errors.New("token无效或已过期")
	}

	// 再验证Redis中的token
	userID, err := s.repo.ValidateToken(token)
	if err != nil {
		return nil, errors.New("token无效或已过期")
	}

	// 确保JWT中的userID和Redis中的一致
	if claims.UserID != userID {
		return nil, errors.New("token不匹配")
	}

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	return user, nil
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// generateToken 生成JWT token
func (s *Service) generateToken(userID uint, username, role string) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.tokenExpires)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ParseToken 解析JWT token（公开方法，供中间件使用）
func (s *Service) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的token")
}
