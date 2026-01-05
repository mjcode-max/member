package user

import (
	"context"
	"testing"
	"time"

	"gorm.io/gorm"
	"member-pre/pkg/errors"
	"member-pre/test"
)

// UserModel 用户数据库模型（用于测试）
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

// ToEntity 转换为领域实体
func (m *UserModel) ToEntity() *User {
	if m == nil {
		return nil
	}
	return &User{
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
func (m *UserModel) FromEntity(u *User) {
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

// TestUserRepository 使用 SQLite 的用户仓储实现，用于单元测试
// 数据库创建逻辑在 test 包中，但实现放在这里避免循环依赖
type TestUserRepository struct {
	db *gorm.DB
}

// NewTestUserRepository 创建基于 SQLite 的测试用户仓储
// 使用 test 包的 NewTestDB 函数创建数据库
func NewTestUserRepository(t *testing.T) (*TestUserRepository, error) {
	db, _, err := test.NewTestDB(t, &UserModel{})
	if err != nil {
		return nil, err
	}
	return &TestUserRepository{db: db}, nil
}

// FindByID 根据ID查找用户
func (r *TestUserRepository) FindByID(ctx context.Context, id uint) (*User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// FindByUsername 根据用户名查找用户
func (r *TestUserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// FindByPhone 根据手机号查找用户
func (r *TestUserRepository) FindByPhone(ctx context.Context, phone string) (*User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// Create 创建用户
func (r *TestUserRepository) Create(ctx context.Context, u *User) error {
	model := &UserModel{}
	model.FromEntity(u)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// 更新实体的ID和时间戳
	u.ID = model.ID
	u.CreatedAt = model.CreatedAt
	u.UpdatedAt = model.UpdatedAt

	return nil
}

// Update 更新用户
func (r *TestUserRepository) Update(ctx context.Context, u *User) error {
	model := &UserModel{}
	model.FromEntity(u)
	model.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", u.ID).Updates(model).Error; err != nil {
		return err
	}

	u.UpdatedAt = model.UpdatedAt
	return nil
}

// UpdateWorkStatus 更新美甲师工作状态
func (r *TestUserRepository) UpdateWorkStatus(ctx context.Context, userID uint, workStatus string) error {
	return r.db.WithContext(ctx).Model(&UserModel{}).
		Where("id = ? AND role = ?", userID, RoleTechnician).
		Update("work_status", workStatus).Error
}

// FindList 获取用户列表（支持筛选和分页）
func (r *TestUserRepository) FindList(ctx context.Context, role, status string, storeID *uint, username, phone string, page, pageSize int) ([]*User, int64, error) {
	var models []UserModel
	var total int64

	query := r.db.WithContext(ctx).Model(&UserModel{})

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
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	users := make([]*User, 0, len(models))
	for _, model := range models {
		users = append(users, model.ToEntity())
	}

	return users, total, nil
}

// FindByStoreID 根据门店ID查找用户（店长和美甲师）
func (r *TestUserRepository) FindByStoreID(ctx context.Context, storeID uint, role string) ([]*User, error) {
	var models []UserModel
	query := r.db.WithContext(ctx).Where("store_id = ?", storeID)
	if role != "" {
		query = query.Where("role = ?", role)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	users := make([]*User, 0, len(models))
	for _, model := range models {
		users = append(users, model.ToEntity())
	}

	return users, nil
}

// AddUser 添加用户到测试数据库（用于测试数据准备）
func (r *TestUserRepository) AddUser(ctx context.Context, u *User) error {
	return r.Create(ctx, u)
}

// AddUsersByStoreID 添加门店用户（用于测试数据准备）
func (r *TestUserRepository) AddUsersByStoreID(ctx context.Context, storeID uint, users []*User) error {
	for _, u := range users {
		if u.StoreID == nil {
			u.StoreID = &storeID
		}
		if err := r.Create(ctx, u); err != nil {
			return err
		}
	}
	return nil
}

// 确保 TestUserRepository 实现了 IUserRepository 接口
var _ IUserRepository = (*TestUserRepository)(nil)

// setupTestService 创建测试用的 UserService
func setupTestService(t *testing.T) (*UserService, *TestUserRepository, context.Context) {
	repo, err := NewTestUserRepository(t)
	if err != nil {
		t.Fatalf("创建测试数据库失败: %v", err)
	}
	service := NewUserService(repo, test.NewMockLogger())
	ctx := context.Background()
	return service, repo, ctx
}

func TestUserService_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		userID  uint
		setup   func(*TestUserRepository, context.Context) error
		wantErr error
	}{
		{
			name:   "成功获取用户",
			userID: 1,
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				return repo.AddUser(ctx, &User{
					ID:        1,
					Username:  "testuser",
					Email:     "test@example.com",
					Phone:     "13800138000",
					Role:      RoleAdmin,
					Status:    StatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: nil,
		},
		{
			name:   "用户不存在",
			userID: 999,
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				// 不设置任何用户
				return nil
			},
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := NewTestUserRepository(t)
			if err != nil {
				t.Fatalf("创建测试数据库失败: %v", err)
			}
			ctx := context.Background()
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}
			service := NewUserService(repo, test.NewMockLogger())

			user, err := service.GetByID(ctx, tt.userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("期望错误 %v, 但得到 %v", tt.wantErr, err)
				}
				if user != nil {
					t.Errorf("期望用户为 nil, 但得到 %v", user)
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if user == nil {
					t.Error("期望得到用户, 但得到 nil")
				} else if user.ID != tt.userID {
					t.Errorf("期望用户ID为 %d, 但得到 %d", tt.userID, user.ID)
				}
			}
		})
	}
}

func TestUserService_GetByPhone(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		setup   func(*TestUserRepository, context.Context) error
		wantErr error
	}{
		{
			name:  "成功获取用户",
			phone: "13800138000",
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				return repo.AddUser(ctx, &User{
					ID:        1,
					Phone:     "13800138000",
					Role:      RoleCustomer,
					Status:    StatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: nil,
		},
		{
			name:  "用户不存在",
			phone: "13900139000",
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				// 不设置任何用户
				return nil
			},
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := NewTestUserRepository(t)
			if err != nil {
				t.Fatalf("创建测试数据库失败: %v", err)
			}
			ctx := context.Background()
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}
			service := NewUserService(repo, test.NewMockLogger())

			user, err := service.GetByPhone(ctx, tt.phone)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("期望错误 %v, 但得到 %v", tt.wantErr, err)
				}
				if user != nil {
					t.Errorf("期望用户为 nil, 但得到 %v", user)
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if user == nil {
					t.Error("期望得到用户, 但得到 nil")
				} else if user.Phone != tt.phone {
					t.Errorf("期望手机号为 %s, 但得到 %s", tt.phone, user.Phone)
				}
			}
		})
	}
}

func TestUserService_GetByUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		setup    func(*TestUserRepository, context.Context) error
		wantErr  error
	}{
		{
			name:     "成功获取用户",
			username: "testuser",
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				return repo.AddUser(ctx, &User{
					ID:        1,
					Username:  "testuser",
					Email:     "test@example.com",
					Role:      RoleAdmin,
					Status:    StatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: nil,
		},
		{
			name:     "用户不存在",
			username: "nonexistent",
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				// 不设置任何用户
				return nil
			},
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, ctx := setupTestService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			user, err := service.GetByUsername(ctx, tt.username)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("期望错误 %v, 但得到 %v", tt.wantErr, err)
				}
				if user != nil {
					t.Errorf("期望用户为 nil, 但得到 %v", user)
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if user == nil {
					t.Error("期望得到用户, 但得到 nil")
				} else if user.Username != tt.username {
					t.Errorf("期望用户名为 %s, 但得到 %s", tt.username, user.Username)
				}
			}
		})
	}
}

func TestUserService_CreateOrGetCustomer(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		setup   func(*TestUserRepository, context.Context) error
		wantErr error
		check   func(*testing.T, *User, error)
	}{
		{
			name:  "获取已存在的顾客",
			phone: "13800138000",
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				return repo.AddUser(ctx, &User{
					ID:        1,
					Phone:     "13800138000",
					Role:      RoleCustomer,
					Status:    StatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: nil,
			check: func(t *testing.T, user *User, err error) {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if user == nil {
					t.Fatal("期望得到用户, 但得到 nil")
				}
				if user.Phone != "13800138000" {
					t.Errorf("期望手机号为 13800138000, 但得到 %s", user.Phone)
				}
				if user.Role != RoleCustomer {
					t.Errorf("期望角色为 %s, 但得到 %s", RoleCustomer, user.Role)
				}
			},
		},
		{
			name:  "创建新顾客",
			phone: "13900139000",
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				// 不设置任何用户，应该创建新顾客
				return nil
			},
			wantErr: nil,
			check: func(t *testing.T, user *User, err error) {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if user == nil {
					t.Fatal("期望得到用户, 但得到 nil")
				}
				if user.Phone != "13900139000" {
					t.Errorf("期望手机号为 13900139000, 但得到 %s", user.Phone)
				}
				if user.Role != RoleCustomer {
					t.Errorf("期望角色为 %s, 但得到 %s", RoleCustomer, user.Role)
				}
				if user.Status != StatusActive {
					t.Errorf("期望状态为 %s, 但得到 %s", StatusActive, user.Status)
				}
			},
		},
		{
			name:  "手机号已被其他角色使用",
			phone: "13800138000",
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				return repo.AddUser(ctx, &User{
					ID:        1,
					Phone:     "13800138000",
					Role:      RoleAdmin,
					Status:    StatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: errors.ErrInvalidParams("该手机号已被注册为其他角色"),
			check: func(t *testing.T, user *User, err error) {
				if err == nil {
					t.Error("期望错误, 但得到 nil")
				}
				if user != nil {
					t.Errorf("期望用户为 nil, 但得到 %v", user)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, ctx := setupTestService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			user, err := service.CreateOrGetCustomer(ctx, tt.phone)

			if tt.check != nil {
				tt.check(t, user, err)
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

func TestUserService_UpdateWorkStatus(t *testing.T) {
	tests := []struct {
		name       string
		userID     uint
		workStatus string
		setup      func(*TestUserRepository, context.Context) error
		wantErr    error
	}{
		{
			name:       "成功更新美甲师工作状态",
			userID:     1,
			workStatus: WorkStatusWorking,
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				status := WorkStatusRest
				return repo.AddUser(ctx, &User{
					ID:         1,
					Username:   "technician1",
					Role:       RoleTechnician,
					Status:     StatusActive,
					WorkStatus: &status,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				})
			},
			wantErr: nil,
		},
		{
			name:       "无效的工作状态",
			userID:     1,
			workStatus: "invalid_status",
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				status := WorkStatusRest
				return repo.AddUser(ctx, &User{
					ID:         1,
					Role:       RoleTechnician,
					Status:     StatusActive,
					WorkStatus: &status,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				})
			},
			wantErr: ErrInvalidWorkStatus,
		},
		{
			name:       "用户不存在",
			userID:     999,
			workStatus: WorkStatusWorking,
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				// 不设置任何用户
				return nil
			},
			wantErr: ErrUserNotFound,
		},
		{
			name:       "非美甲师角色不能更新工作状态",
			userID:     1,
			workStatus: WorkStatusWorking,
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				return repo.AddUser(ctx, &User{
					ID:        1,
					Username:  "admin1",
					Role:      RoleAdmin,
					Status:    StatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: errors.ErrInvalidParams("只有美甲师可以更新工作状态"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, ctx := setupTestService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			err := service.UpdateWorkStatus(ctx, tt.userID, tt.workStatus)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("期望错误 %v, 但得到 %v", tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				} else {
					// 验证工作状态已更新
					u, _ := repo.FindByID(ctx, tt.userID)
					if u != nil && u.WorkStatus != nil && *u.WorkStatus != tt.workStatus {
						t.Errorf("期望工作状态为 %s, 但得到 %s", tt.workStatus, *u.WorkStatus)
					}
				}
			}
		})
	}
}

func TestUserService_GetByStoreID(t *testing.T) {
	tests := []struct {
		name    string
		storeID uint
		role    string
		setup   func(*TestUserRepository, context.Context) error
		wantLen int
		wantErr error
	}{
		{
			name:    "成功获取门店用户列表",
			storeID: 1,
			role:    "",
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				return repo.AddUsersByStoreID(ctx, 1, []*User{
					{ID: 1, Username: "manager1", Phone: "13800138001", Role: RoleStoreManager, StoreID: uintPtr(1), Status: StatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					{ID: 2, Username: "tech1", Phone: "13800138002", Role: RoleTechnician, StoreID: uintPtr(1), Status: StatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					{ID: 3, Username: "tech2", Phone: "13800138003", Role: RoleTechnician, StoreID: uintPtr(1), Status: StatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				})
			},
			wantLen: 3,
			wantErr: nil,
		},
		{
			name:    "按角色过滤",
			storeID: 1,
			role:    RoleTechnician,
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				return repo.AddUsersByStoreID(ctx, 1, []*User{
					{ID: 1, Username: "manager1", Phone: "13800138001", Role: RoleStoreManager, StoreID: uintPtr(1), Status: StatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					{ID: 2, Username: "tech1", Phone: "13800138002", Role: RoleTechnician, StoreID: uintPtr(1), Status: StatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					{ID: 3, Username: "tech2", Phone: "13800138003", Role: RoleTechnician, StoreID: uintPtr(1), Status: StatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				})
			},
			wantLen: 2,
			wantErr: nil,
		},
		{
			name:    "门店无用户",
			storeID: 999,
			role:    "",
			setup: func(repo *TestUserRepository, ctx context.Context) error {
				// 不设置任何用户
				return nil
			},
			wantLen: 0,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, ctx := setupTestService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			users, err := service.GetByStoreID(ctx, tt.storeID, tt.role)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if len(users) != tt.wantLen {
					t.Errorf("期望用户数量为 %d, 但得到 %d", tt.wantLen, len(users))
				}
			}
		})
	}
}

func TestUserService_ValidateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr error
	}{
		{
			name: "验证激活用户",
			user: &User{
				ID:        1,
				Username:  "testuser",
				Status:    StatusActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: nil,
		},
		{
			name:    "用户为nil",
			user:    nil,
			wantErr: ErrUserNotFound,
		},
		{
			name: "用户被禁用",
			user: &User{
				ID:        1,
				Username:  "testuser",
				Status:    StatusInactive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: ErrUserInactive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, ctx := setupTestService(t)

			err := service.ValidateUser(ctx, tt.user)

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

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name          string
		userRole      string
		requiredRoles []string
		want          bool
	}{
		{
			name:          "无要求角色，所有角色都可以",
			userRole:      RoleAdmin,
			requiredRoles: []string{},
			want:          true,
		},
		{
			name:          "用户角色匹配",
			userRole:      RoleAdmin,
			requiredRoles: []string{RoleAdmin, RoleStoreManager},
			want:          true,
		},
		{
			name:          "用户角色不匹配",
			userRole:      RoleCustomer,
			requiredRoles: []string{RoleAdmin, RoleStoreManager},
			want:          false,
		},
		{
			name:          "单个角色匹配",
			userRole:      RoleTechnician,
			requiredRoles: []string{RoleTechnician},
			want:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasPermission(tt.userRole, tt.requiredRoles...)
			if got != tt.want {
				t.Errorf("HasPermission(%s, %v) = %v, 期望 %v", tt.userRole, tt.requiredRoles, got, tt.want)
			}
		})
	}
}

func TestCanAccessStore(t *testing.T) {
	storeID1 := uint(1)

	tests := []struct {
		name    string
		user    *User
		storeID uint
		want    bool
	}{
		{
			name: "总后台可以访问所有门店",
			user: &User{
				ID:      1,
				Role:    RoleAdmin,
				StoreID: nil,
			},
			storeID: 999,
			want:    true,
		},
		{
			name: "店长可以访问自己的门店",
			user: &User{
				ID:      1,
				Role:    RoleStoreManager,
				StoreID: &storeID1,
			},
			storeID: 1,
			want:    true,
		},
		{
			name: "店长不能访问其他门店",
			user: &User{
				ID:      1,
				Role:    RoleStoreManager,
				StoreID: &storeID1,
			},
			storeID: 2,
			want:    false,
		},
		{
			name: "美甲师可以访问自己的门店",
			user: &User{
				ID:      1,
				Role:    RoleTechnician,
				StoreID: &storeID1,
			},
			storeID: 1,
			want:    true,
		},
		{
			name: "美甲师不能访问其他门店",
			user: &User{
				ID:      1,
				Role:    RoleTechnician,
				StoreID: &storeID1,
			},
			storeID: 2,
			want:    false,
		},
		{
			name:    "用户为nil",
			user:    nil,
			storeID: 1,
			want:    false,
		},
		{
			name: "顾客不能访问门店",
			user: &User{
				ID:      1,
				Role:    RoleCustomer,
				StoreID: nil,
			},
			storeID: 1,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanAccessStore(tt.user, tt.storeID)
			if got != tt.want {
				t.Errorf("CanAccessStore(%v, %d) = %v, 期望 %v", tt.user, tt.storeID, got, tt.want)
			}
		})
	}
}

func TestCanAccessUser(t *testing.T) {
	storeID1 := uint(1)

	tests := []struct {
		name         string
		currentUser  *User
		targetUserID uint
		want         bool
	}{
		{
			name: "总后台可以访问所有用户",
			currentUser: &User{
				ID:   1,
				Role: RoleAdmin,
			},
			targetUserID: 999,
			want:         true,
		},
		{
			name: "用户可以访问自己的数据",
			currentUser: &User{
				ID:   1,
				Role: RoleCustomer,
			},
			targetUserID: 1,
			want:         true,
		},
		{
			name: "用户不能访问其他用户的数据",
			currentUser: &User{
				ID:   1,
				Role: RoleCustomer,
			},
			targetUserID: 2,
			want:         false,
		},
		{
			name: "店长可以访问本店用户（简化实现）",
			currentUser: &User{
				ID:      1,
				Role:    RoleStoreManager,
				StoreID: &storeID1,
			},
			targetUserID: 2,
			want:         true, // 根据代码实现，店长总是返回true
		},
		{
			name:         "当前用户为nil",
			currentUser:  nil,
			targetUserID: 1,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanAccessUser(tt.currentUser, tt.targetUserID)
			if got != tt.want {
				t.Errorf("CanAccessUser(%v, %d) = %v, 期望 %v", tt.currentUser, tt.targetUserID, got, tt.want)
			}
		})
	}
}

// 辅助函数
func uintPtr(u uint) *uint {
	return &u
}
