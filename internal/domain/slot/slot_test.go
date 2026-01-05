package slot

import (
	"context"
	"testing"
	"time"

	"gorm.io/gorm"
	"member-pre/pkg/errors"
	"member-pre/test"
)

// SlotModel 时段数据库模型（用于测试）
type SlotModel struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	StoreID      uint           `gorm:"index;not null" json:"store_id"`
	TechnicianID *uint          `gorm:"index" json:"technician_id"`
	Date         time.Time      `gorm:"type:date;index" json:"date"`
	StartTime    time.Time      `gorm:"not null" json:"start_time"`
	EndTime      time.Time      `gorm:"not null" json:"end_time"`
	Capacity     int            `gorm:"default:0;not null" json:"capacity"`
	LockedCount  int            `gorm:"default:0;not null" json:"locked_count"`
	BookedCount  int            `gorm:"default:0;not null" json:"booked_count"`
	Status       string         `gorm:"size:20;default:'available';not null;index" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SlotModel) TableName() string {
	return "slots"
}

// ToEntity 转换为领域实体
func (m *SlotModel) ToEntity() *Slot {
	if m == nil {
		return nil
	}
	return &Slot{
		ID:           m.ID,
		StoreID:      m.StoreID,
		TechnicianID: m.TechnicianID,
		Date:         m.Date,
		StartTime:    m.StartTime,
		EndTime:      m.EndTime,
		Capacity:     m.Capacity,
		LockedCount:  m.LockedCount,
		BookedCount:  m.BookedCount,
		Status:       m.Status,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// FromEntity 从领域实体创建模型
func (m *SlotModel) FromEntity(s *Slot) {
	if s == nil {
		return
	}
	m.ID = s.ID
	m.StoreID = s.StoreID
	m.TechnicianID = s.TechnicianID
	m.Date = s.Date
	m.StartTime = s.StartTime
	m.EndTime = s.EndTime
	m.Capacity = s.Capacity
	m.LockedCount = s.LockedCount
	m.BookedCount = s.BookedCount
	m.Status = s.Status
	m.CreatedAt = s.CreatedAt
	m.UpdatedAt = s.UpdatedAt
}

// TestSlotRepository 使用 SQLite 的时段仓储实现，用于单元测试
type TestSlotRepository struct {
	db *gorm.DB
}

// NewTestSlotRepository 创建基于 SQLite 的测试时段仓储
func NewTestSlotRepository(t *testing.T) (*TestSlotRepository, error) {
	db, _, err := test.NewTestDB(t, &SlotModel{})
	if err != nil {
		return nil, err
	}
	return &TestSlotRepository{db: db}, nil
}

// FindByID 根据ID查找时段
func (r *TestSlotRepository) FindByID(ctx context.Context, id uint) (*Slot, error) {
	var model SlotModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// FindByStoreIDAndDate 根据门店ID和日期查找时段列表
func (r *TestSlotRepository) FindByStoreIDAndDate(ctx context.Context, storeID uint, date time.Time) ([]*Slot, error) {
	var models []SlotModel
	dateStr := date.Format("2006-01-02")
	if err := r.db.WithContext(ctx).
		Where("store_id = ? AND DATE(date) = ?", storeID, dateStr).
		Order("start_time ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	slots := make([]*Slot, 0, len(models))
	for _, model := range models {
		slots = append(slots, model.ToEntity())
	}

	return slots, nil
}

// FindByStoreIDAndDateRange 根据门店ID和日期范围查找时段列表
func (r *TestSlotRepository) FindByStoreIDAndDateRange(ctx context.Context, storeID uint, startDate, endDate time.Time) ([]*Slot, error) {
	var models []SlotModel
	if err := r.db.WithContext(ctx).
		Where("store_id = ? AND date >= ? AND date <= ?", storeID, startDate, endDate).
		Order("date ASC, start_time ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	slots := make([]*Slot, 0, len(models))
	for _, model := range models {
		slots = append(slots, model.ToEntity())
	}

	return slots, nil
}

// FindAvailableByStoreIDAndDate 根据门店ID和日期查找可用时段列表
func (r *TestSlotRepository) FindAvailableByStoreIDAndDate(ctx context.Context, storeID uint, date time.Time) ([]*Slot, error) {
	var models []SlotModel
	dateStr := date.Format("2006-01-02")
	if err := r.db.WithContext(ctx).
		Where("store_id = ? AND DATE(date) = ? AND status = ? AND (capacity - locked_count - booked_count) > 0",
			storeID, dateStr, SlotStatusAvailable).
		Order("start_time ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	slots := make([]*Slot, 0, len(models))
	for _, model := range models {
		slots = append(slots, model.ToEntity())
	}

	return slots, nil
}

// Create 创建时段
func (r *TestSlotRepository) Create(ctx context.Context, s *Slot) error {
	model := &SlotModel{}
	model.FromEntity(s)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	s.ID = model.ID
	s.CreatedAt = model.CreatedAt
	s.UpdatedAt = model.UpdatedAt

	return nil
}

// CreateBatch 批量创建时段
func (r *TestSlotRepository) CreateBatch(ctx context.Context, slots []*Slot) error {
	if len(slots) == 0 {
		return nil
	}

	models := make([]SlotModel, 0, len(slots))
	now := time.Now()
	for _, s := range slots {
		model := &SlotModel{}
		model.FromEntity(s)
		model.CreatedAt = now
		model.UpdatedAt = now
		models = append(models, *model)
	}

	if err := r.db.WithContext(ctx).CreateInBatches(models, 100).Error; err != nil {
		return err
	}

	// 更新实体的ID和时间戳
	for i, model := range models {
		slots[i].ID = model.ID
		slots[i].CreatedAt = model.CreatedAt
		slots[i].UpdatedAt = model.UpdatedAt
	}

	return nil
}

// Update 更新时段
func (r *TestSlotRepository) Update(ctx context.Context, s *Slot) error {
	model := &SlotModel{}
	model.FromEntity(s)
	model.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Model(&SlotModel{}).Where("id = ?", s.ID).Updates(model).Error; err != nil {
		return err
	}

	s.UpdatedAt = model.UpdatedAt
	return nil
}

// UpdateBatch 批量更新时段
func (r *TestSlotRepository) UpdateBatch(ctx context.Context, slots []*Slot) error {
	if len(slots) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, s := range slots {
			model := &SlotModel{}
			model.FromEntity(s)
			model.UpdatedAt = time.Now()

			if err := tx.Model(&SlotModel{}).Where("id = ?", s.ID).Updates(model).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// LockSlot 锁定时段（原子操作，增加锁定数量）
func (r *TestSlotRepository) LockSlot(ctx context.Context, slotID uint, count int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model SlotModel
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("id = ?", slotID).
			First(&model).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return ErrSlotNotFound
			}
			return err
		}

		available := model.Capacity - model.LockedCount - model.BookedCount
		if available < count {
			return ErrInsufficientCapacity
		}

		if err := tx.Model(&SlotModel{}).
			Where("id = ?", slotID).
			Update("locked_count", gorm.Expr("locked_count + ?", count)).Error; err != nil {
			return err
		}

		return nil
	})
}

// UnlockSlot 解锁时段（原子操作，减少锁定数量）
func (r *TestSlotRepository) UnlockSlot(ctx context.Context, slotID uint, count int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&SlotModel{}).
			Where("id = ? AND locked_count >= ?", slotID, count).
			Update("locked_count", gorm.Expr("locked_count - ?", count))

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.ErrInvalidParams("解锁数量超过当前锁定数量")
		}

		return nil
	})
}

// BookSlot 预约时段（原子操作，从锁定转为已预约）
func (r *TestSlotRepository) BookSlot(ctx context.Context, slotID uint, count int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model SlotModel
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("id = ?", slotID).
			First(&model).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return ErrSlotNotFound
			}
			return err
		}

		if model.LockedCount < count {
			return errors.ErrInvalidParams("预约数量超过当前锁定数量")
		}

		if err := tx.Model(&SlotModel{}).
			Where("id = ?", slotID).
			Updates(map[string]interface{}{
				"locked_count": gorm.Expr("locked_count - ?", count),
				"booked_count": gorm.Expr("booked_count + ?", count),
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

// ReleaseSlot 释放时段（原子操作，取消或完成时释放）
func (r *TestSlotRepository) ReleaseSlot(ctx context.Context, slotID uint, count int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model SlotModel
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("id = ?", slotID).
			First(&model).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return ErrSlotNotFound
			}
			return err
		}

		if model.BookedCount < count {
			return errors.ErrInvalidParams("释放数量超过当前已预约数量")
		}

		if err := tx.Model(&SlotModel{}).
			Where("id = ?", slotID).
			Update("booked_count", gorm.Expr("booked_count - ?", count)).Error; err != nil {
			return err
		}

		return nil
	})
}

// DeleteByTechnicianIDAndFuture 删除指定美甲师的所有未来时段
func (r *TestSlotRepository) DeleteByTechnicianIDAndFuture(ctx context.Context, technicianID uint, fromDate time.Time) error {
	return r.db.WithContext(ctx).
		Where("technician_id = ? AND date >= ?", technicianID, fromDate).
		Delete(&SlotModel{}).Error
}

// FindByTechnicianIDAndDateRange 根据美甲师ID和日期范围查找时段列表
func (r *TestSlotRepository) FindByTechnicianIDAndDateRange(ctx context.Context, technicianID uint, startDate, endDate time.Time) ([]*Slot, error) {
	var models []SlotModel
	if err := r.db.WithContext(ctx).
		Where("technician_id = ? AND date >= ? AND date <= ?", technicianID, startDate, endDate).
		Order("date ASC, start_time ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	slots := make([]*Slot, 0, len(models))
	for _, model := range models {
		slots = append(slots, model.ToEntity())
	}

	return slots, nil
}

// 确保 TestSlotRepository 实现了 ISlotRepository 接口
var _ ISlotRepository = (*TestSlotRepository)(nil)

// TestTemplateRepositoryForSlot 用于 SlotService 测试的模板仓储
type TestTemplateRepositoryForSlot struct {
	templates map[uint]*Template
}

func NewTestTemplateRepositoryForSlot() *TestTemplateRepositoryForSlot {
	return &TestTemplateRepositoryForSlot{
		templates: make(map[uint]*Template),
	}
}

func (r *TestTemplateRepositoryForSlot) FindByID(ctx context.Context, id uint) (*Template, error) {
	return r.templates[id], nil
}

func (r *TestTemplateRepositoryForSlot) FindByStoreID(ctx context.Context, storeID uint) ([]*Template, error) {
	var result []*Template
	for _, t := range r.templates {
		if t.StoreID == storeID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *TestTemplateRepositoryForSlot) FindActiveByStoreID(ctx context.Context, storeID uint) (*Template, error) {
	for _, t := range r.templates {
		if t.StoreID == storeID && t.Status == TemplateStatusActive {
			return t, nil
		}
	}
	return nil, nil
}

func (r *TestTemplateRepositoryForSlot) Create(ctx context.Context, template *Template) error {
	r.templates[template.ID] = template
	return nil
}

func (r *TestTemplateRepositoryForSlot) Update(ctx context.Context, template *Template) error {
	r.templates[template.ID] = template
	return nil
}

func (r *TestTemplateRepositoryForSlot) Delete(ctx context.Context, id uint) error {
	delete(r.templates, id)
	return nil
}

// setupTestSlotService 创建测试用的 SlotService
func setupTestSlotService(t *testing.T) (*SlotService, *TestSlotRepository, *TestTemplateRepositoryForSlot, context.Context) {
	slotRepo, err := NewTestSlotRepository(t)
	if err != nil {
		t.Fatalf("创建测试数据库失败: %v", err)
	}
	templateRepo := NewTestTemplateRepositoryForSlot()
	service := NewSlotService(slotRepo, templateRepo, test.NewMockLogger())
	ctx := context.Background()
	return service, slotRepo, templateRepo, ctx
}

// addTestSlot 添加时段到测试数据库（用于测试数据准备）
func (r *TestSlotRepository) addTestSlot(ctx context.Context, slot *Slot) error {
	return r.Create(ctx, slot)
}

func TestSlotService_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		slotID  uint
		setup   func(*TestSlotRepository, context.Context) error
		wantErr error
	}{
		{
			name:   "成功获取时段",
			slotID: 1,
			setup: func(repo *TestSlotRepository, ctx context.Context) error {
				now := time.Now()
				return repo.addTestSlot(ctx, &Slot{
					ID:          1,
					StoreID:     1,
					Date:        now,
					StartTime:   now,
					EndTime:     now.Add(time.Hour),
					Capacity:    5,
					LockedCount: 0,
					BookedCount: 0,
					Status:      SlotStatusAvailable,
					CreatedAt:   now,
					UpdatedAt:   now,
				})
			},
			wantErr: nil,
		},
		{
			name:   "时段不存在",
			slotID: 999,
			setup: func(repo *TestSlotRepository, ctx context.Context) error {
				return nil
			},
			wantErr: ErrSlotNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, _, ctx := setupTestSlotService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			slot, err := service.GetByID(ctx, tt.slotID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("期望错误 %v, 但得到 %v", tt.wantErr, err)
				}
				if slot != nil {
					t.Errorf("期望时段为 nil, 但得到 %v", slot)
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if slot == nil {
					t.Error("期望得到时段, 但得到 nil")
				} else if slot.ID != tt.slotID {
					t.Errorf("期望时段ID为 %d, 但得到 %d", tt.slotID, slot.ID)
				}
			}
		})
	}
}

func TestSlotService_GetByStoreIDAndDate(t *testing.T) {
	tests := []struct {
		name    string
		storeID uint
		date    time.Time
		setup   func(*TestSlotRepository, context.Context) error
		wantLen int
		wantErr error
	}{
		{
			name:    "成功获取时段列表",
			storeID: 1,
			date:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			setup: func(repo *TestSlotRepository, ctx context.Context) error {
				date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				repo.addTestSlot(ctx, &Slot{
					StoreID:     1,
					Date:        date,
					StartTime:   date.Add(9 * time.Hour),
					EndTime:     date.Add(10 * time.Hour),
					Capacity:    5,
					LockedCount: 0,
					BookedCount: 0,
					Status:      SlotStatusAvailable,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
				return repo.addTestSlot(ctx, &Slot{
					StoreID:     1,
					Date:        date,
					StartTime:   date.Add(10 * time.Hour),
					EndTime:     date.Add(11 * time.Hour),
					Capacity:    5,
					LockedCount: 0,
					BookedCount: 0,
					Status:      SlotStatusAvailable,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				})
			},
			wantLen: 2,
			wantErr: nil,
		},
		{
			name:    "无时段",
			storeID: 1,
			date:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			setup: func(repo *TestSlotRepository, ctx context.Context) error {
				return nil
			},
			wantLen: 0,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, _, ctx := setupTestSlotService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			slots, err := service.GetByStoreIDAndDate(ctx, tt.storeID, tt.date)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if len(slots) != tt.wantLen {
					t.Errorf("期望时段数量为 %d, 但得到 %d", tt.wantLen, len(slots))
				}
			}
		})
	}
}

func TestSlotService_LockSlot(t *testing.T) {
	tests := []struct {
		name    string
		slotID  uint
		count   int
		setup   func(*TestSlotRepository, context.Context) error
		wantErr error
	}{
		{
			name:   "成功锁定时段",
			slotID: 1,
			count:  2,
			setup: func(repo *TestSlotRepository, ctx context.Context) error {
				now := time.Now()
				return repo.addTestSlot(ctx, &Slot{
					ID:          1,
					StoreID:     1,
					Date:        now,
					StartTime:   now,
					EndTime:     now.Add(time.Hour),
					Capacity:    5,
					LockedCount: 0,
					BookedCount: 0,
					Status:      SlotStatusAvailable,
					CreatedAt:   now,
					UpdatedAt:   now,
				})
			},
			wantErr: nil,
		},
		{
			name:   "容量不足",
			slotID: 1,
			count:  10,
			setup: func(repo *TestSlotRepository, ctx context.Context) error {
				now := time.Now()
				return repo.addTestSlot(ctx, &Slot{
					ID:          1,
					StoreID:     1,
					Date:        now,
					StartTime:   now,
					EndTime:     now.Add(time.Hour),
					Capacity:    5,
					LockedCount: 0,
					BookedCount: 0,
					Status:      SlotStatusAvailable,
					CreatedAt:   now,
					UpdatedAt:   now,
				})
			},
			wantErr: ErrInsufficientCapacity,
		},
		{
			name:   "时段不存在",
			slotID: 999,
			count:  1,
			setup: func(repo *TestSlotRepository, ctx context.Context) error {
				return nil
			},
			wantErr: ErrSlotNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, _, ctx := setupTestSlotService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			err := service.LockSlot(ctx, tt.slotID, tt.count)

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

func TestSlotService_UnlockSlot(t *testing.T) {
	tests := []struct {
		name    string
		slotID  uint
		count   int
		setup   func(*TestSlotRepository, context.Context) error
		wantErr error
	}{
		{
			name:   "成功解锁时段",
			slotID: 1,
			count:  2,
			setup: func(repo *TestSlotRepository, ctx context.Context) error {
				now := time.Now()
				return repo.addTestSlot(ctx, &Slot{
					ID:          1,
					StoreID:     1,
					Date:        now,
					StartTime:   now,
					EndTime:     now.Add(time.Hour),
					Capacity:    5,
					LockedCount: 3,
					BookedCount: 0,
					Status:      SlotStatusAvailable,
					CreatedAt:   now,
					UpdatedAt:   now,
				})
			},
			wantErr: nil,
		},
		{
			name:   "解锁数量超过锁定数量",
			slotID: 1,
			count:  10,
			setup: func(repo *TestSlotRepository, ctx context.Context) error {
				now := time.Now()
				return repo.addTestSlot(ctx, &Slot{
					ID:          1,
					StoreID:     1,
					Date:        now,
					StartTime:   now,
					EndTime:     now.Add(time.Hour),
					Capacity:    5,
					LockedCount: 3,
					BookedCount: 0,
					Status:      SlotStatusAvailable,
					CreatedAt:   now,
					UpdatedAt:   now,
				})
			},
			wantErr: errors.ErrInvalidParams("解锁数量超过当前锁定数量"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, _, ctx := setupTestSlotService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			err := service.UnlockSlot(ctx, tt.slotID, tt.count)

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

func TestSlotService_BookSlot(t *testing.T) {
	tests := []struct {
		name    string
		slotID  uint
		count   int
		setup   func(*TestSlotRepository, context.Context) error
		wantErr error
	}{
		{
			name:   "成功预约时段",
			slotID: 1,
			count:  2,
			setup: func(repo *TestSlotRepository, ctx context.Context) error {
				now := time.Now()
				return repo.addTestSlot(ctx, &Slot{
					ID:          1,
					StoreID:     1,
					Date:        now,
					StartTime:   now,
					EndTime:     now.Add(time.Hour),
					Capacity:    5,
					LockedCount: 3,
					BookedCount: 0,
					Status:      SlotStatusAvailable,
					CreatedAt:   now,
					UpdatedAt:   now,
				})
			},
			wantErr: nil,
		},
		{
			name:   "预约数量超过锁定数量",
			slotID: 1,
			count:  10,
			setup: func(repo *TestSlotRepository, ctx context.Context) error {
				now := time.Now()
				return repo.addTestSlot(ctx, &Slot{
					ID:          1,
					StoreID:     1,
					Date:        now,
					StartTime:   now,
					EndTime:     now.Add(time.Hour),
					Capacity:    5,
					LockedCount: 3,
					BookedCount: 0,
					Status:      SlotStatusAvailable,
					CreatedAt:   now,
					UpdatedAt:   now,
				})
			},
			wantErr: errors.ErrInvalidParams("预约数量超过当前锁定数量"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, _, ctx := setupTestSlotService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			err := service.BookSlot(ctx, tt.slotID, tt.count)

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

func TestSlotService_ReleaseSlot(t *testing.T) {
	tests := []struct {
		name    string
		slotID  uint
		count   int
		setup   func(*TestSlotRepository, context.Context) error
		wantErr error
	}{
		{
			name:   "成功释放时段",
			slotID: 1,
			count:  2,
			setup: func(repo *TestSlotRepository, ctx context.Context) error {
				now := time.Now()
				return repo.addTestSlot(ctx, &Slot{
					ID:          1,
					StoreID:     1,
					Date:        now,
					StartTime:   now,
					EndTime:     now.Add(time.Hour),
					Capacity:    5,
					LockedCount: 0,
					BookedCount: 3,
					Status:      SlotStatusBooked,
					CreatedAt:   now,
					UpdatedAt:   now,
				})
			},
			wantErr: nil,
		},
		{
			name:   "释放数量超过已预约数量",
			slotID: 1,
			count:  10,
			setup: func(repo *TestSlotRepository, ctx context.Context) error {
				now := time.Now()
				return repo.addTestSlot(ctx, &Slot{
					ID:          1,
					StoreID:     1,
					Date:        now,
					StartTime:   now,
					EndTime:     now.Add(time.Hour),
					Capacity:    5,
					LockedCount: 0,
					BookedCount: 3,
					Status:      SlotStatusBooked,
					CreatedAt:   now,
					UpdatedAt:   now,
				})
			},
			wantErr: errors.ErrInvalidParams("释放数量超过当前已预约数量"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, _, ctx := setupTestSlotService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			err := service.ReleaseSlot(ctx, tt.slotID, tt.count)

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

