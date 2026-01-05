package store

import (
	"context"
	"testing"
	"time"

	"gorm.io/gorm"
	"member-pre/test"
)

// StoreModel 门店数据库模型（用于测试）
type StoreModel struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Name               string         `gorm:"size:100;not null" json:"name"`
	Address            string         `gorm:"size:255" json:"address"`
	Phone              string         `gorm:"size:20" json:"phone"`
	ContactPerson      string         `gorm:"size:50" json:"contact_person"`
	Status             string         `gorm:"size:20;default:'operating';not null" json:"status"`
	BusinessHoursStart string         `gorm:"size:5" json:"business_hours_start"`
	BusinessHoursEnd   string         `gorm:"size:5" json:"business_hours_end"`
	DepositAmount      float64        `json:"deposit_amount"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (StoreModel) TableName() string {
	return "stores"
}

// ToEntity 转换为领域实体
func (m *StoreModel) ToEntity() *Store {
	if m == nil {
		return nil
	}
	return &Store{
		ID:                 m.ID,
		Name:               m.Name,
		Address:            m.Address,
		Phone:              m.Phone,
		ContactPerson:      m.ContactPerson,
		Status:             m.Status,
		BusinessHoursStart: m.BusinessHoursStart,
		BusinessHoursEnd:   m.BusinessHoursEnd,
		DepositAmount:      m.DepositAmount,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

// FromEntity 从领域实体创建模型
func (m *StoreModel) FromEntity(s *Store) {
	if s == nil {
		return
	}
	m.ID = s.ID
	m.Name = s.Name
	m.Address = s.Address
	m.Phone = s.Phone
	m.ContactPerson = s.ContactPerson
	m.Status = s.Status
	m.BusinessHoursStart = s.BusinessHoursStart
	m.BusinessHoursEnd = s.BusinessHoursEnd
	m.DepositAmount = s.DepositAmount
	m.CreatedAt = s.CreatedAt
	m.UpdatedAt = s.UpdatedAt
}

// TestStoreRepository 使用 SQLite 的门店仓储实现，用于单元测试
type TestStoreRepository struct {
	db *gorm.DB
}

// NewTestStoreRepository 创建基于 SQLite 的测试门店仓储
func NewTestStoreRepository(t *testing.T) (*TestStoreRepository, error) {
	db, _, err := test.NewTestDB(t, &StoreModel{})
	if err != nil {
		return nil, err
	}
	return &TestStoreRepository{db: db}, nil
}

// FindByID 根据ID查找门店
func (r *TestStoreRepository) FindByID(ctx context.Context, id uint) (*Store, error) {
	var model StoreModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// FindList 获取门店列表（支持筛选和分页）
func (r *TestStoreRepository) FindList(ctx context.Context, status, name string, page, pageSize int) ([]*Store, int64, error) {
	var models []StoreModel
	var total int64

	query := r.db.WithContext(ctx).Model(&StoreModel{})

	// 筛选条件
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
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

	stores := make([]*Store, 0, len(models))
	for _, model := range models {
		stores = append(stores, model.ToEntity())
	}

	return stores, total, nil
}

// Create 创建门店
func (r *TestStoreRepository) Create(ctx context.Context, store *Store) error {
	model := &StoreModel{}
	model.FromEntity(store)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// 更新实体的ID和时间戳
	store.ID = model.ID
	store.CreatedAt = model.CreatedAt
	store.UpdatedAt = model.UpdatedAt

	return nil
}

// Update 更新门店
func (r *TestStoreRepository) Update(ctx context.Context, store *Store) error {
	model := &StoreModel{}
	model.FromEntity(store)
	model.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Model(&StoreModel{}).Where("id = ?", store.ID).Updates(model).Error; err != nil {
		return err
	}

	store.UpdatedAt = model.UpdatedAt
	return nil
}

// Delete 删除门店（软删除）
func (r *TestStoreRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&StoreModel{}, id).Error
}

// FindByManagerID 根据店长ID查找门店
func (r *TestStoreRepository) FindByManagerID(ctx context.Context, managerID uint) (*Store, error) {
	// 简化实现：通过 manager_id 关联查找（实际实现可能需要关联 users 表）
	var model StoreModel
	// 这里假设有一个 manager_id 字段，实际可能需要通过 users 表关联
	// 为了测试，我们简化处理
	if err := r.db.WithContext(ctx).Where("id = ?", managerID).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// FindByStatus 根据状态查找门店列表
func (r *TestStoreRepository) FindByStatus(ctx context.Context, status string) ([]*Store, error) {
	var models []StoreModel
	if err := r.db.WithContext(ctx).Where("status = ?", status).Find(&models).Error; err != nil {
		return nil, err
	}

	stores := make([]*Store, 0, len(models))
	for _, model := range models {
		stores = append(stores, model.ToEntity())
	}

	return stores, nil
}

// 确保 TestStoreRepository 实现了 IStoreRepository 接口
var _ IStoreRepository = (*TestStoreRepository)(nil)

// setupTestService 创建测试用的 StoreService
func setupTestStoreService(t *testing.T) (*StoreService, *TestStoreRepository, context.Context) {
	repo, err := NewTestStoreRepository(t)
	if err != nil {
		t.Fatalf("创建测试数据库失败: %v", err)
	}
	service := NewStoreService(repo, test.NewMockLogger())
	ctx := context.Background()
	return service, repo, ctx
}

// addTestStore 添加门店到测试数据库（用于测试数据准备）
func (r *TestStoreRepository) addTestStore(ctx context.Context, s *Store) error {
	return r.Create(ctx, s)
}

func TestStoreService_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		storeID uint
		setup   func(*TestStoreRepository, context.Context) error
		wantErr error
	}{
		{
			name:    "成功获取门店",
			storeID: 1,
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				return repo.addTestStore(ctx, &Store{
					ID:        1,
					Name:      "测试门店1",
					Address:   "测试地址1",
					Status:    StatusOperating,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: nil,
		},
		{
			name:    "门店不存在",
			storeID: 999,
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				// 不设置任何门店
				return nil
			},
			wantErr: ErrStoreNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, ctx := setupTestStoreService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			store, err := service.GetByID(ctx, tt.storeID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("期望错误 %v, 但得到 %v", tt.wantErr, err)
				}
				if store != nil {
					t.Errorf("期望门店为 nil, 但得到 %v", store)
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if store == nil {
					t.Error("期望得到门店, 但得到 nil")
				} else if store.ID != tt.storeID {
					t.Errorf("期望门店ID为 %d, 但得到 %d", tt.storeID, store.ID)
				}
			}
		})
	}
}

func TestStoreService_GetList(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		nameFilter string
		page       int
		pageSize   int
		setup      func(*TestStoreRepository, context.Context) error
		wantLen    int
		wantTotal  int64
		wantErr    error
	}{
		{
			name:       "获取所有门店",
			status:     "",
			nameFilter: "",
			page:       1,
			pageSize:   10,
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				repo.addTestStore(ctx, &Store{Name: "门店1", Status: StatusOperating, CreatedAt: time.Now(), UpdatedAt: time.Now()})
				repo.addTestStore(ctx, &Store{Name: "门店2", Status: StatusClosed, CreatedAt: time.Now(), UpdatedAt: time.Now()})
				return repo.addTestStore(ctx, &Store{Name: "门店3", Status: StatusOperating, CreatedAt: time.Now(), UpdatedAt: time.Now()})
			},
			wantLen:   3,
			wantTotal: 3,
			wantErr:   nil,
		},
		{
			name:       "按状态筛选",
			status:     StatusOperating,
			nameFilter: "",
			page:       1,
			pageSize:   10,
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				repo.addTestStore(ctx, &Store{Name: "门店1", Status: StatusOperating, CreatedAt: time.Now(), UpdatedAt: time.Now()})
				repo.addTestStore(ctx, &Store{Name: "门店2", Status: StatusClosed, CreatedAt: time.Now(), UpdatedAt: time.Now()})
				return repo.addTestStore(ctx, &Store{Name: "门店3", Status: StatusOperating, CreatedAt: time.Now(), UpdatedAt: time.Now()})
			},
			wantLen:   2,
			wantTotal: 2,
			wantErr:   nil,
		},
		{
			name:       "按名称模糊搜索",
			status:     "",
			nameFilter: "测试",
			page:       1,
			pageSize:   10,
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				repo.addTestStore(ctx, &Store{Name: "测试门店1", Status: StatusOperating, CreatedAt: time.Now(), UpdatedAt: time.Now()})
				repo.addTestStore(ctx, &Store{Name: "测试门店2", Status: StatusOperating, CreatedAt: time.Now(), UpdatedAt: time.Now()})
				return repo.addTestStore(ctx, &Store{Name: "其他门店", Status: StatusOperating, CreatedAt: time.Now(), UpdatedAt: time.Now()})
			},
			wantLen:   2,
			wantTotal: 2,
			wantErr:   nil,
		},
		{
			name:       "分页测试",
			status:     "",
			nameFilter: "",
			page:       1,
			pageSize:   2,
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				repo.addTestStore(ctx, &Store{Name: "门店1", Status: StatusOperating, CreatedAt: time.Now(), UpdatedAt: time.Now()})
				repo.addTestStore(ctx, &Store{Name: "门店2", Status: StatusOperating, CreatedAt: time.Now(), UpdatedAt: time.Now()})
				return repo.addTestStore(ctx, &Store{Name: "门店3", Status: StatusOperating, CreatedAt: time.Now(), UpdatedAt: time.Now()})
			},
			wantLen:   2,
			wantTotal: 3,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, ctx := setupTestStoreService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			stores, total, err := service.GetList(ctx, tt.status, tt.nameFilter, tt.page, tt.pageSize)

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
				if len(stores) != tt.wantLen {
					t.Errorf("期望门店数量为 %d, 但得到 %d", tt.wantLen, len(stores))
				}
				if total != tt.wantTotal {
					t.Errorf("期望总数为 %d, 但得到 %d", tt.wantTotal, total)
				}
			}
		})
	}
}

func TestStoreService_Create(t *testing.T) {
	tests := []struct {
		name    string
		store   *Store
		wantErr error
	}{
		{
			name: "成功创建门店",
			store: &Store{
				Name:               "新门店",
				Address:            "新地址",
				Phone:              "13800138000",
				ContactPerson:      "联系人",
				Status:             StatusOperating,
				BusinessHoursStart: "09:00",
				BusinessHoursEnd:   "18:00",
				DepositAmount:      100.0,
			},
			wantErr: nil,
		},
		{
			name: "门店名称为空",
			store: &Store{
				Name:   "",
				Status: StatusOperating,
			},
			wantErr: ErrNameRequired,
		},
		{
			name: "无效的状态",
			store: &Store{
				Name:   "门店1",
				Status: "invalid_status",
			},
			wantErr: ErrInvalidStatus,
		},
		{
			name: "无效的营业时间格式",
			store: &Store{
				Name:               "门店1",
				Status:             StatusOperating,
				BusinessHoursStart: "25:00",
			},
			wantErr: ErrInvalidTimeFormat,
		},
		{
			name: "押金金额为负数",
			store: &Store{
				Name:          "门店1",
				Status:        StatusOperating,
				DepositAmount: -100.0,
			},
			wantErr: ErrInvalidDeposit,
		},
		{
			name: "默认状态为营业中",
			store: &Store{
				Name: "门店1",
				// 不设置 Status，应该默认为 operating
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, ctx := setupTestStoreService(t)

			err := service.Create(ctx, tt.store)

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
				if tt.store.ID == 0 {
					t.Error("期望门店ID不为0")
				}
				if tt.store.Status == "" {
					t.Error("期望门店状态不为空")
				}
			}
		})
	}
}

func TestStoreService_Update(t *testing.T) {
	tests := []struct {
		name    string
		storeID uint
		update  *Store
		setup   func(*TestStoreRepository, context.Context) error
		wantErr error
	}{
		{
			name:    "成功更新门店",
			storeID: 1,
			update: &Store{
				ID:      1,
				Name:    "更新后的门店名",
				Address: "更新后的地址",
			},
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				return repo.addTestStore(ctx, &Store{
					ID:        1,
					Name:      "原门店名",
					Address:   "原地址",
					Status:    StatusOperating,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: nil,
		},
		{
			name:    "门店不存在",
			storeID: 999,
			update: &Store{
				ID:   999,
				Name: "不存在的门店",
			},
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				return nil
			},
			wantErr: ErrStoreNotFound,
		},
		{
			name:    "无效的状态",
			storeID: 1,
			update: &Store{
				ID:     1,
				Status: "invalid_status",
			},
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				return repo.addTestStore(ctx, &Store{
					ID:        1,
					Name:      "门店1",
					Status:    StatusOperating,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, ctx := setupTestStoreService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			err := service.Update(ctx, tt.update)

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

func TestStoreService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		storeID uint
		setup   func(*TestStoreRepository, context.Context) error
		wantErr error
	}{
		{
			name:    "成功删除门店",
			storeID: 1,
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				return repo.addTestStore(ctx, &Store{
					ID:        1,
					Name:      "门店1",
					Status:    StatusOperating,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: nil,
		},
		{
			name:    "门店不存在",
			storeID: 999,
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				return nil
			},
			wantErr: ErrStoreNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, ctx := setupTestStoreService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			err := service.Delete(ctx, tt.storeID)

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

func TestStoreService_UpdateStatus(t *testing.T) {
	tests := []struct {
		name    string
		storeID uint
		status  string
		setup   func(*TestStoreRepository, context.Context) error
		wantErr error
	}{
		{
			name:    "成功更新门店状态",
			storeID: 1,
			status:  StatusClosed,
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				return repo.addTestStore(ctx, &Store{
					ID:        1,
					Name:      "门店1",
					Status:    StatusOperating,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: nil,
		},
		{
			name:    "无效的状态",
			storeID: 1,
			status:  "invalid_status",
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				return repo.addTestStore(ctx, &Store{
					ID:        1,
					Name:      "门店1",
					Status:    StatusOperating,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: ErrInvalidStatus,
		},
		{
			name:    "门店不存在",
			storeID: 999,
			status:  StatusOperating,
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				return nil
			},
			wantErr: ErrStoreNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, ctx := setupTestStoreService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			err := service.UpdateStatus(ctx, tt.storeID, tt.status)

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

func TestStoreService_GetByStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		setup   func(*TestStoreRepository, context.Context) error
		wantLen int
		wantErr error
	}{
		{
			name:   "成功获取营业中的门店",
			status: StatusOperating,
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				repo.addTestStore(ctx, &Store{Name: "门店1", Status: StatusOperating, CreatedAt: time.Now(), UpdatedAt: time.Now()})
				repo.addTestStore(ctx, &Store{Name: "门店2", Status: StatusClosed, CreatedAt: time.Now(), UpdatedAt: time.Now()})
				return repo.addTestStore(ctx, &Store{Name: "门店3", Status: StatusOperating, CreatedAt: time.Now(), UpdatedAt: time.Now()})
			},
			wantLen: 2,
			wantErr: nil,
		},
		{
			name:   "无效的状态",
			status: "invalid_status",
			setup: func(repo *TestStoreRepository, ctx context.Context) error {
				return nil
			},
			wantLen: 0,
			wantErr: ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, ctx := setupTestStoreService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			stores, err := service.GetByStatus(ctx, tt.status)

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
				if len(stores) != tt.wantLen {
					t.Errorf("期望门店数量为 %d, 但得到 %d", tt.wantLen, len(stores))
				}
			}
		})
	}
}

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"营业中", StatusOperating, true},
		{"停业", StatusClosed, true},
		{"关闭", StatusShutdown, true},
		{"无效状态", "invalid", false},
		{"空字符串", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidStatus(tt.status)
			if got != tt.want {
				t.Errorf("isValidStatus(%s) = %v, 期望 %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestIsValidTimeFormat(t *testing.T) {
	tests := []struct {
		name    string
		timeStr string
		want    bool
	}{
		{"有效时间 09:00", "09:00", true},
		{"有效时间 23:59", "23:59", true},
		{"有效时间 00:00", "00:00", true},
		{"无效时间 25:00", "25:00", false},
		{"无效时间 12:60", "12:60", false},
		{"无效格式 9:00", "9:00", false},
		{"无效格式 09:0", "09:0", false},
		{"空字符串", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidTimeFormat(tt.timeStr)
			if got != tt.want {
				t.Errorf("isValidTimeFormat(%s) = %v, 期望 %v", tt.timeStr, got, tt.want)
			}
		})
	}
}
