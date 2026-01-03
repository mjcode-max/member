package auth

import (
	"context"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	"member-pre/internal/domain/user"
	"member-pre/test"
)

// MockAuthConfig 模拟认证配置
type MockAuthConfig struct {
	jwtSecret    string
	tokenExpires int
}

func NewMockAuthConfig() *MockAuthConfig {
	return &MockAuthConfig{
		jwtSecret:    "test-secret-key-for-unit-testing",
		tokenExpires: 7200, // 2小时
	}
}

func (m *MockAuthConfig) GetJwtSecret() string {
	return m.jwtSecret
}

func (m *MockAuthConfig) GetTokenExpires() int {
	return m.tokenExpires
}

// MockUserService 模拟用户服务，用于单元测试
type MockUserService struct {
	usersByID       map[uint]*user.User
	usersByPhone    map[string]*user.User
	usersByUsername map[string]*user.User
	createError     error
	validateError   error
}

func NewMockUserService() *MockUserService {
	return &MockUserService{
		usersByID:       make(map[uint]*user.User),
		usersByPhone:    make(map[string]*user.User),
		usersByUsername: make(map[string]*user.User),
	}
}

func (m *MockUserService) GetByID(ctx context.Context, id uint) (*user.User, error) {
	u, ok := m.usersByID[id]
	if !ok {
		return nil, user.ErrUserNotFound
	}
	return u, nil
}

func (m *MockUserService) GetByPhone(ctx context.Context, phone string) (*user.User, error) {
	u, ok := m.usersByPhone[phone]
	if !ok {
		return nil, user.ErrUserNotFound
	}
	return u, nil
}

func (m *MockUserService) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	u, ok := m.usersByUsername[username]
	if !ok {
		return nil, user.ErrUserNotFound
	}
	return u, nil
}

func (m *MockUserService) CreateOrGetCustomer(ctx context.Context, phone string) (*user.User, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	// 如果已存在，返回现有用户
	if u, ok := m.usersByPhone[phone]; ok {
		return u, nil
	}
	// 创建新用户
	u := &user.User{
		ID:        uint(len(m.usersByID) + 1),
		Phone:     phone,
		Role:      user.RoleCustomer,
		Status:    user.StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.usersByID[u.ID] = u
	m.usersByPhone[phone] = u
	return u, nil
}

func (m *MockUserService) ValidateUser(ctx context.Context, u *user.User) error {
	if m.validateError != nil {
		return m.validateError
	}
	if u == nil {
		return user.ErrUserNotFound
	}
	if u.Status != user.StatusActive {
		return user.ErrUserInactive
	}
	return nil
}

// AddUser 添加用户到mock（用于测试数据准备）
func (m *MockUserService) AddUser(u *user.User) {
	m.usersByID[u.ID] = u
	if u.Phone != "" {
		m.usersByPhone[u.Phone] = u
	}
	if u.Username != "" {
		m.usersByUsername[u.Username] = u
	}
}

// setupTestService 创建测试用的 AuthService
func setupTestAuthService(t *testing.T) (*AuthService, *MockUserService, context.Context) {
	mockUserService := NewMockUserService()
	authService := NewAuthService(mockUserService, test.NewMockLogger(), NewMockAuthConfig())
	ctx := context.Background()
	return authService, mockUserService, ctx
}

func TestAuthService_Login_Customer(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr error
		check   func(*testing.T, *LoginResponse, error)
	}{
		{
			name:    "顾客手机号登录-新用户自动创建",
			phone:   "13800138000",
			wantErr: nil,
			check: func(t *testing.T, resp *LoginResponse, err error) {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if resp == nil {
					t.Fatal("期望得到登录响应, 但得到 nil")
				}
				if resp.User == nil {
					t.Fatal("期望得到用户, 但得到 nil")
				}
				if resp.User.Phone != "13800138000" {
					t.Errorf("期望手机号为 13800138000, 但得到 %s", resp.User.Phone)
				}
				if resp.User.Role != user.RoleCustomer {
					t.Errorf("期望角色为 %s, 但得到 %s", user.RoleCustomer, resp.User.Role)
				}
				if resp.Token == "" {
					t.Error("期望得到token, 但得到空字符串")
				}
				if resp.ExpiresAt.IsZero() {
					t.Error("期望得到过期时间, 但得到零值")
				}
			},
		},
		{
			name:    "顾客手机号登录-已存在用户",
			phone:   "13900139000",
			wantErr: nil,
			check: func(t *testing.T, resp *LoginResponse, err error) {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if resp == nil {
					t.Fatal("期望得到登录响应, 但得到 nil")
				}
				if resp.User.Phone != "13900139000" {
					t.Errorf("期望手机号为 13900139000, 但得到 %s", resp.User.Phone)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService, mockUserService, ctx := setupTestAuthService(t)

			// 如果是已存在用户测试，先创建用户
			if tt.name == "顾客手机号登录-已存在用户" {
				mockUserService.AddUser(&user.User{
					ID:        1,
					Phone:     tt.phone,
					Role:      user.RoleCustomer,
					Status:    user.StatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			}

			req := &LoginRequest{
				Username: tt.phone,
				Password: "", // 顾客不需要密码
			}

			resp, err := authService.Login(ctx, req)

			if tt.check != nil {
				tt.check(t, resp, err)
			} else {
				if tt.wantErr != nil {
					if err == nil {
						t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
					} else if err.Error() != tt.wantErr.Error() {
						t.Errorf("期望错误 %v, 但得到 %v", tt.wantErr, err)
					}
				} else {
					if err != nil {
						t.Errorf("不期望错误, 但得到 %v", err)
					}
				}
			}
		})
	}
}

func TestAuthService_Login_Staff(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		setup    func(*MockUserService) error
		wantErr  error
		check    func(*testing.T, *LoginResponse, error)
	}{
		{
			name:     "员工登录-成功",
			username: "admin1",
			password: "password123",
			setup: func(mock *MockUserService) error {
				hashedPassword, err := HashPassword("password123")
				if err != nil {
					return err
				}
				mock.AddUser(&user.User{
					ID:        1,
					Username:  "admin1",
					Password:  hashedPassword,
					Role:      user.RoleAdmin,
					Status:    user.StatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				return nil
			},
			wantErr: nil,
			check: func(t *testing.T, resp *LoginResponse, err error) {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if resp == nil {
					t.Fatal("期望得到登录响应, 但得到 nil")
				}
				if resp.User.Username != "admin1" {
					t.Errorf("期望用户名为 admin1, 但得到 %s", resp.User.Username)
				}
				if resp.Token == "" {
					t.Error("期望得到token, 但得到空字符串")
				}
			},
		},
		{
			name:     "员工登录-密码为空",
			username: "admin1",
			password: "",
			setup: func(mock *MockUserService) error {
				mock.AddUser(&user.User{
					ID:        1,
					Username:  "admin1",
					Password:  "hashed_password",
					Role:      user.RoleAdmin,
					Status:    user.StatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				return nil
			},
			wantErr: ErrInvalidPassword,
		},
		{
			name:     "员工登录-用户不存在",
			username: "nonexistent",
			password: "password123",
			setup: func(mock *MockUserService) error {
				// 不创建用户
				return nil
			},
			wantErr: ErrUserNotFound,
		},
		{
			name:     "员工登录-密码错误",
			username: "admin1",
			password: "wrong_password",
			setup: func(mock *MockUserService) error {
				hashedPassword, err := HashPassword("correct_password")
				if err != nil {
					return err
				}
				mock.AddUser(&user.User{
					ID:        1,
					Username:  "admin1",
					Password:  hashedPassword,
					Role:      user.RoleAdmin,
					Status:    user.StatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				return nil
			},
			wantErr: ErrInvalidPassword,
		},
		{
			name:     "员工登录-用户被禁用",
			username: "admin1",
			password: "password123",
			setup: func(mock *MockUserService) error {
				hashedPassword, err := HashPassword("password123")
				if err != nil {
					return err
				}
				mock.AddUser(&user.User{
					ID:        1,
					Username:  "admin1",
					Password:  hashedPassword,
					Role:      user.RoleAdmin,
					Status:    user.StatusInactive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				return nil
			},
			wantErr: user.ErrUserInactive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService, mockUserService, ctx := setupTestAuthService(t)

			if err := tt.setup(mockUserService); err != nil {
				t.Fatalf("准备测试数据失败: %v", err)
			}

			req := &LoginRequest{
				Username: tt.username,
				Password: tt.password,
			}

			resp, err := authService.Login(ctx, req)

			if tt.check != nil {
				tt.check(t, resp, err)
			} else {
				if tt.wantErr != nil {
					if err == nil {
						t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
					} else if err.Error() != tt.wantErr.Error() {
						t.Errorf("期望错误 %v, 但得到 %v", tt.wantErr, err)
					}
				} else {
					if err != nil {
						t.Errorf("不期望错误, 但得到 %v", err)
					}
				}
			}
		})
	}
}

func TestAuthService_GetCurrentUser(t *testing.T) {
	tests := []struct {
		name    string
		userID  uint
		setup   func(*MockUserService) error
		wantErr error
		check   func(*testing.T, *user.User, error)
	}{
		{
			name:   "成功获取当前用户",
			userID: 1,
			setup: func(mock *MockUserService) error {
				mock.AddUser(&user.User{
					ID:        1,
					Username:  "testuser",
					Role:      user.RoleAdmin,
					Status:    user.StatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				return nil
			},
			wantErr: nil,
			check: func(t *testing.T, u *user.User, err error) {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if u == nil {
					t.Fatal("期望得到用户, 但得到 nil")
				}
				if u.ID != 1 {
					t.Errorf("期望用户ID为 1, 但得到 %d", u.ID)
				}
			},
		},
		{
			name:   "用户不存在",
			userID: 999,
			setup: func(mock *MockUserService) error {
				// 不创建用户
				return nil
			},
			wantErr: user.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService, mockUserService, ctx := setupTestAuthService(t)

			if err := tt.setup(mockUserService); err != nil {
				t.Fatalf("准备测试数据失败: %v", err)
			}

			u, err := authService.GetCurrentUser(ctx, tt.userID)

			if tt.check != nil {
				tt.check(t, u, err)
			} else {
				if tt.wantErr != nil {
					if err == nil {
						t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
					} else if err.Error() != tt.wantErr.Error() {
						t.Errorf("期望错误 %v, 但得到 %v", tt.wantErr, err)
					}
				} else {
					if err != nil {
						t.Errorf("不期望错误, 但得到 %v", err)
					}
				}
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	authService, mockUserService, _ := setupTestAuthService(t)

	// 创建测试用户并生成token
	testUser := &user.User{
		ID:        1,
		Username:  "testuser",
		Role:      user.RoleAdmin,
		Status:    user.StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockUserService.AddUser(testUser)

	validToken, _, err := authService.GenerateTokenForUser(testUser)
	if err != nil {
		t.Fatalf("生成token失败: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr error
		check   func(*testing.T, *TokenClaims, error)
	}{
		{
			name:    "验证有效token",
			token:   validToken,
			wantErr: nil,
			check: func(t *testing.T, claims *TokenClaims, err error) {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if claims == nil {
					t.Fatal("期望得到claims, 但得到 nil")
				}
				if claims.UserID != 1 {
					t.Errorf("期望用户ID为 1, 但得到 %d", claims.UserID)
				}
				if claims.Role != user.RoleAdmin {
					t.Errorf("期望角色为 %s, 但得到 %s", user.RoleAdmin, claims.Role)
				}
			},
		},
		{
			name:    "验证无效token",
			token:   "invalid.token.string",
			wantErr: ErrInvalidToken,
		},
		{
			name:    "验证空token",
			token:   "",
			wantErr: ErrInvalidToken,
		},
		{
			name:    "验证错误签名的token",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoiYWRtaW4ifQ.invalid_signature",
			wantErr: ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := authService.ValidateToken(tt.token)

			if tt.check != nil {
				tt.check(t, claims, err)
			} else {
				if tt.wantErr != nil {
					if err == nil {
						t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
					} else if err.Error() != tt.wantErr.Error() {
						t.Errorf("期望错误 %v, 但得到 %v", tt.wantErr, err)
					}
				} else {
					if err != nil {
						t.Errorf("不期望错误, 但得到 %v", err)
					}
				}
			}
		})
	}
}

func TestAuthService_GenerateTokenForUser(t *testing.T) {
	authService, mockUserService, _ := setupTestAuthService(t)

	tests := []struct {
		name  string
		user  *user.User
		setup func(*MockUserService) error
		check func(*testing.T, *AuthService, string, time.Time, error)
	}{
		{
			name: "为管理员生成token",
			user: &user.User{
				ID:        1,
				Username:  "admin1",
				Role:      user.RoleAdmin,
				Status:    user.StatusActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setup: func(mock *MockUserService) error {
				mock.AddUser(&user.User{
					ID:        1,
					Username:  "admin1",
					Role:      user.RoleAdmin,
					Status:    user.StatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				return nil
			},
			check: func(t *testing.T, authService *AuthService, token string, expiresAt time.Time, err error) {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if token == "" {
					t.Error("期望得到token, 但得到空字符串")
				}
				if expiresAt.IsZero() {
					t.Error("期望得到过期时间, 但得到零值")
				}
				// 验证token可以解析
				claims, err := authService.ValidateToken(token)
				if err != nil {
					t.Errorf("生成的token无法验证: %v", err)
				}
				if claims.UserID != 1 {
					t.Errorf("期望用户ID为 1, 但得到 %d", claims.UserID)
				}
			},
		},
		{
			name: "为店长生成token（带门店ID）",
			user: &user.User{
				ID:        2,
				Username:  "manager1",
				Role:      user.RoleStoreManager,
				StoreID:   uintPtr(1),
				Status:    user.StatusActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setup: func(mock *MockUserService) error {
				storeID := uint(1)
				mock.AddUser(&user.User{
					ID:        2,
					Username:  "manager1",
					Role:      user.RoleStoreManager,
					StoreID:   &storeID,
					Status:    user.StatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				return nil
			},
			check: func(t *testing.T, authService *AuthService, token string, expiresAt time.Time, err error) {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				claims, err := authService.ValidateToken(token)
				if err != nil {
					t.Errorf("生成的token无法验证: %v", err)
				}
				if claims.StoreID == nil || *claims.StoreID != 1 {
					t.Errorf("期望门店ID为 1, 但得到 %v", claims.StoreID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.setup(mockUserService); err != nil {
				t.Fatalf("准备测试数据失败: %v", err)
			}

			token, expiresAt, err := authService.GenerateTokenForUser(tt.user)

			if tt.check != nil {
				tt.check(t, authService, token, expiresAt, err)
			} else {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	authService, _, ctx := setupTestAuthService(t)

	tests := []struct {
		name    string
		token   string
		wantErr error
	}{
		{
			name:    "成功登出",
			token:   "test_token_string",
			wantErr: nil,
		},
		{
			name:    "空token登出",
			token:   "",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.Logout(ctx, tt.token)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("期望错误 %v, 但得到 %v", tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
		check    func(*testing.T, string, error)
	}{
		{
			name:     "成功加密密码",
			password: "test_password_123",
			wantErr:  nil,
			check: func(t *testing.T, hashed string, err error) {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if hashed == "" {
					t.Error("期望得到加密后的密码, 但得到空字符串")
				}
				if hashed == "test_password_123" {
					t.Error("加密后的密码不应与原文相同")
				}
				// 验证可以正确验证密码
				err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte("test_password_123"))
				if err != nil {
					t.Errorf("加密后的密码无法验证: %v", err)
				}
			},
		},
		{
			name:     "空密码",
			password: "",
			wantErr:  nil,
			check: func(t *testing.T, hashed string, err error) {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if hashed == "" {
					t.Error("期望得到加密后的密码, 但得到空字符串")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := HashPassword(tt.password)

			if tt.check != nil {
				tt.check(t, hashed, err)
			} else {
				if tt.wantErr != nil {
					if err == nil {
						t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
					}
				} else {
					if err != nil {
						t.Errorf("不期望错误, 但得到 %v", err)
					}
				}
			}
		})
	}
}

func TestIsPhoneNumber(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "有效手机号",
			input: "13800138000",
			want:  true,
		},
		{
			name:  "10位数字",
			input: "1380013800",
			want:  false,
		},
		{
			name:  "12位数字",
			input: "138001380001",
			want:  false,
		},
		{
			name:  "包含字母",
			input: "1380013800a",
			want:  false,
		},
		{
			name:  "包含特殊字符",
			input: "138-0013-8000",
			want:  false,
		},
		{
			name:  "空字符串",
			input: "",
			want:  false,
		},
		{
			name:  "用户名",
			input: "admin",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isPhoneNumber(tt.input)
			if got != tt.want {
				t.Errorf("isPhoneNumber(%s) = %v, 期望 %v", tt.input, got, tt.want)
			}
		})
	}
}

// 辅助函数
func uintPtr(u uint) *uint {
	return &u
}
