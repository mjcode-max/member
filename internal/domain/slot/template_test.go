package slot

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"testing"
	"time"

	"gorm.io/gorm"
	"member-pre/pkg/errors"
	"member-pre/test"
)

// TemplateModel 时段模板数据库模型（用于测试）
type TemplateModel struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	StoreID      uint            `gorm:"index;not null" json:"store_id"`
	Name         string          `gorm:"size:100;not null" json:"name"`
	Status       string          `gorm:"size:20;default:'active';not null" json:"status"`
	WeekdayRules WeekdayRuleJSON `gorm:"type:text" json:"weekday_rules"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// WeekdayRuleJSON 星期规则JSON类型（用于数据库存储）
type WeekdayRuleJSON []WeekdayRule

// Value 实现driver.Valuer接口
func (w WeekdayRuleJSON) Value() (driver.Value, error) {
	if len(w) == 0 {
		return "[]", nil
	}
	return json.Marshal(w)
}

// Scan 实现sql.Scanner接口
func (w *WeekdayRuleJSON) Scan(value interface{}) error {
	if value == nil {
		*w = WeekdayRuleJSON{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.ErrInvalidParams("无法扫描WeekdayRuleJSON")
	}
	if len(bytes) == 0 {
		*w = WeekdayRuleJSON{}
		return nil
	}
	return json.Unmarshal(bytes, w)
}

// TableName 指定表名
func (TemplateModel) TableName() string {
	return "slot_templates"
}

// ToEntity 转换为领域实体
func (m *TemplateModel) ToEntity() *Template {
	if m == nil {
		return nil
	}
	return &Template{
		ID:           m.ID,
		StoreID:      m.StoreID,
		Name:         m.Name,
		Status:       m.Status,
		WeekdayRules: []WeekdayRule(m.WeekdayRules),
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// FromEntity 从领域实体创建模型
func (m *TemplateModel) FromEntity(t *Template) {
	if t == nil {
		return
	}
	m.ID = t.ID
	m.StoreID = t.StoreID
	m.Name = t.Name
	m.Status = t.Status
	m.WeekdayRules = WeekdayRuleJSON(t.WeekdayRules)
	m.CreatedAt = t.CreatedAt
	m.UpdatedAt = t.UpdatedAt
}

// TestTemplateRepository 使用 SQLite 的时段模板仓储实现，用于单元测试
type TestTemplateRepository struct {
	db *gorm.DB
}

// NewTestTemplateRepository 创建基于 SQLite 的测试时段模板仓储
func NewTestTemplateRepository(t *testing.T) (*TestTemplateRepository, error) {
	db, _, err := test.NewTestDB(t, &TemplateModel{})
	if err != nil {
		return nil, err
	}
	return &TestTemplateRepository{db: db}, nil
}

// FindByID 根据ID查找模板
func (r *TestTemplateRepository) FindByID(ctx context.Context, id uint) (*Template, error) {
	var model TemplateModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

// FindByStoreID 根据门店ID查找模板列表
func (r *TestTemplateRepository) FindByStoreID(ctx context.Context, storeID uint) ([]*Template, error) {
	var models []TemplateModel
	if err := r.db.WithContext(ctx).Where("store_id = ?", storeID).Find(&models).Error; err != nil {
		return nil, err
	}

	templates := make([]*Template, 0, len(models))
	for _, model := range models {
		templates = append(templates, model.ToEntity())
	}

	return templates, nil
}

// FindActiveByStoreID 根据门店ID查找启用的模板
func (r *TestTemplateRepository) FindActiveByStoreID(ctx context.Context, storeID uint) (*Template, error) {
	var model TemplateModel
	if err := r.db.WithContext(ctx).
		Where("store_id = ? AND status = ?", storeID, TemplateStatusActive).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

// Create 创建模板
func (r *TestTemplateRepository) Create(ctx context.Context, t *Template) error {
	model := &TemplateModel{}
	model.FromEntity(t)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	t.ID = model.ID
	t.CreatedAt = model.CreatedAt
	t.UpdatedAt = model.UpdatedAt

	return nil
}

// Update 更新模板
func (r *TestTemplateRepository) Update(ctx context.Context, t *Template) error {
	model := &TemplateModel{}
	model.FromEntity(t)
	model.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Model(&TemplateModel{}).Where("id = ?", t.ID).Updates(model).Error; err != nil {
		return err
	}

	t.UpdatedAt = model.UpdatedAt
	return nil
}

// Delete 删除模板
func (r *TestTemplateRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&TemplateModel{}, id).Error
}

// 确保 TestTemplateRepository 实现了 ITemplateRepository 接口
var _ ITemplateRepository = (*TestTemplateRepository)(nil)

// setupTestTemplateService 创建测试用的 TemplateService
func setupTestTemplateService(t *testing.T) (*TemplateService, *TestTemplateRepository, context.Context) {
	repo, err := NewTestTemplateRepository(t)
	if err != nil {
		t.Fatalf("创建测试数据库失败: %v", err)
	}
	service := NewTemplateService(repo, test.NewMockLogger())
	ctx := context.Background()
	return service, repo, ctx
}

// addTestTemplate 添加模板到测试数据库（用于测试数据准备）
func (r *TestTemplateRepository) addTestTemplate(ctx context.Context, template *Template) error {
	return r.Create(ctx, template)
}

func TestTemplateService_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		templateID uint
		setup     func(*TestTemplateRepository, context.Context) error
		wantErr   error
	}{
		{
			name:      "成功获取模板",
			templateID: 1,
			setup: func(repo *TestTemplateRepository, ctx context.Context) error {
				return repo.addTestTemplate(ctx, &Template{
					ID:        1,
					StoreID:   1,
					Name:      "测试模板",
					Status:    TemplateStatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: nil,
		},
		{
			name:      "模板不存在",
			templateID: 999,
			setup: func(repo *TestTemplateRepository, ctx context.Context) error {
				return nil
			},
			wantErr: ErrTemplateNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, ctx := setupTestTemplateService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			template, err := service.GetByID(ctx, tt.templateID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("期望错误 %v, 但得到 %v", tt.wantErr, err)
				}
				if template != nil {
					t.Errorf("期望模板为 nil, 但得到 %v", template)
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if template == nil {
					t.Error("期望得到模板, 但得到 nil")
				} else if template.ID != tt.templateID {
					t.Errorf("期望模板ID为 %d, 但得到 %d", tt.templateID, template.ID)
				}
			}
		})
	}
}

func TestTemplateService_GetByStoreID(t *testing.T) {
	tests := []struct {
		name    string
		storeID uint
		setup   func(*TestTemplateRepository, context.Context) error
		wantLen int
		wantErr error
	}{
		{
			name:    "成功获取门店模板列表",
			storeID: 1,
			setup: func(repo *TestTemplateRepository, ctx context.Context) error {
				repo.addTestTemplate(ctx, &Template{
					StoreID:   1,
					Name:      "模板1",
					Status:    TemplateStatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				repo.addTestTemplate(ctx, &Template{
					StoreID:   1,
					Name:      "模板2",
					Status:    TemplateStatusInactive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				return repo.addTestTemplate(ctx, &Template{
					StoreID:   2,
					Name:      "模板3",
					Status:    TemplateStatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantLen: 2,
			wantErr: nil,
		},
		{
			name:    "门店无模板",
			storeID: 999,
			setup: func(repo *TestTemplateRepository, ctx context.Context) error {
				return nil
			},
			wantLen: 0,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, ctx := setupTestTemplateService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			templates, err := service.GetByStoreID(ctx, tt.storeID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if len(templates) != tt.wantLen {
					t.Errorf("期望模板数量为 %d, 但得到 %d", tt.wantLen, len(templates))
				}
			}
		})
	}
}

func TestTemplateService_GetActiveByStoreID(t *testing.T) {
	tests := []struct {
		name    string
		storeID uint
		setup   func(*TestTemplateRepository, context.Context) error
		wantErr error
	}{
		{
			name:    "成功获取启用的模板",
			storeID: 1,
			setup: func(repo *TestTemplateRepository, ctx context.Context) error {
				repo.addTestTemplate(ctx, &Template{
					StoreID:   1,
					Name:      "禁用模板",
					Status:    TemplateStatusInactive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
				return repo.addTestTemplate(ctx, &Template{
					StoreID:   1,
					Name:      "启用模板",
					Status:    TemplateStatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: nil,
		},
		{
			name:    "门店没有启用的模板",
			storeID: 1,
			setup: func(repo *TestTemplateRepository, ctx context.Context) error {
				return repo.addTestTemplate(ctx, &Template{
					StoreID:   1,
					Name:      "禁用模板",
					Status:    TemplateStatusInactive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: ErrTemplateNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, ctx := setupTestTemplateService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			template, err := service.GetActiveByStoreID(ctx, tt.storeID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("期望错误 %v, 但得到 nil", tt.wantErr)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("期望错误 %v, 但得到 %v", tt.wantErr, err)
				}
				if template != nil {
					t.Errorf("期望模板为 nil, 但得到 %v", template)
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误, 但得到 %v", err)
				}
				if template == nil {
					t.Error("期望得到模板, 但得到 nil")
				} else if template.Status != TemplateStatusActive {
					t.Errorf("期望模板状态为 %s, 但得到 %s", TemplateStatusActive, template.Status)
				}
			}
		})
	}
}

func TestTemplateService_Create(t *testing.T) {
	tests := []struct {
		name    string
		template *Template
		wantErr error
	}{
		{
			name: "成功创建模板",
			template: &Template{
				StoreID: 1,
				Name:    "新模板",
				Status:  TemplateStatusActive,
				WeekdayRules: []WeekdayRule{
					{
						Weekday: 1, // 周一
						Slots: []TimeSlotRule{
							{StartTime: "09:00", EndTime: "18:00", Duration: 60},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "模板名称为空",
			template: &Template{
				StoreID: 1,
				Name:    "",
			},
			wantErr: ErrInvalidTemplateName,
		},
		{
			name: "无效的星期规则（星期超出范围）",
			template: &Template{
				StoreID: 1,
				Name:    "测试模板",
				WeekdayRules: []WeekdayRule{
					{
						Weekday: 7, // 无效的星期
						Slots:   []TimeSlotRule{},
					},
				},
			},
			wantErr: ErrInvalidWeekdayRule,
		},
		{
			name: "无效的时间格式",
			template: &Template{
				StoreID: 1,
				Name:    "测试模板",
				WeekdayRules: []WeekdayRule{
					{
						Weekday: 1,
						Slots: []TimeSlotRule{
							{StartTime: "9:00", EndTime: "18:00", Duration: 60}, // 格式不正确，缺少前导0
						},
					},
				},
			},
			wantErr: ErrInvalidTimeSlot,
		},
		{
			name: "开始时间晚于结束时间",
			template: &Template{
				StoreID: 1,
				Name:    "测试模板",
				WeekdayRules: []WeekdayRule{
					{
						Weekday: 1,
						Slots: []TimeSlotRule{
							{StartTime: "18:00", EndTime: "09:00", Duration: 60},
						},
					},
				},
			},
			wantErr: errors.ErrInvalidParams("开始时间必须早于结束时间"),
		},
		{
			name: "默认状态为启用",
			template: &Template{
				StoreID: 1,
				Name:    "测试模板",
				WeekdayRules: []WeekdayRule{
					{
						Weekday: 1,
						Slots: []TimeSlotRule{
							{StartTime: "09:00", EndTime: "18:00", Duration: 60},
						},
					},
				},
				// 不设置 Status，应该默认为 active
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, ctx := setupTestTemplateService(t)

			err := service.Create(ctx, tt.template)

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
				if tt.template.ID == 0 {
					t.Error("期望模板ID不为0")
				}
				if tt.template.Status == "" {
					t.Error("期望模板状态不为空")
				}
			}
		})
	}
}

func TestTemplateService_Update(t *testing.T) {
	tests := []struct {
		name    string
		templateID uint
		update  *Template
		setup   func(*TestTemplateRepository, context.Context) error
		wantErr error
	}{
		{
			name:      "成功更新模板",
			templateID: 1,
			update: &Template{
				ID:      1,
				Name:    "更新后的模板名",
				Status:  TemplateStatusInactive,
			},
			setup: func(repo *TestTemplateRepository, ctx context.Context) error {
				return repo.addTestTemplate(ctx, &Template{
					ID:        1,
					StoreID:   1,
					Name:      "原模板名",
					Status:    TemplateStatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: nil,
		},
		{
			name:      "模板不存在",
			templateID: 999,
			update: &Template{
				ID:   999,
				Name: "不存在的模板",
			},
			setup: func(repo *TestTemplateRepository, ctx context.Context) error {
				return nil
			},
			wantErr: ErrTemplateNotFound,
		},
		{
			name:      "无效的星期规则",
			templateID: 1,
			update: &Template{
				ID:   1,
				Name: "测试模板",
				WeekdayRules: []WeekdayRule{
					{
						Weekday: 7, // 无效的星期
						Slots:   []TimeSlotRule{},
					},
				},
			},
			setup: func(repo *TestTemplateRepository, ctx context.Context) error {
				return repo.addTestTemplate(ctx, &Template{
					ID:        1,
					StoreID:   1,
					Name:      "测试模板",
					Status:    TemplateStatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: ErrInvalidWeekdayRule,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, ctx := setupTestTemplateService(t)
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

func TestTemplateService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		templateID uint
		setup     func(*TestTemplateRepository, context.Context) error
		wantErr   error
	}{
		{
			name:      "成功删除模板",
			templateID: 1,
			setup: func(repo *TestTemplateRepository, ctx context.Context) error {
				return repo.addTestTemplate(ctx, &Template{
					ID:        1,
					StoreID:   1,
					Name:      "测试模板",
					Status:    TemplateStatusActive,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				})
			},
			wantErr: nil,
		},
		{
			name:      "模板不存在",
			templateID: 999,
			setup: func(repo *TestTemplateRepository, ctx context.Context) error {
				return nil
			},
			wantErr: ErrTemplateNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, repo, ctx := setupTestTemplateService(t)
			if err := tt.setup(repo, ctx); err != nil {
				t.Fatalf("测试数据准备失败: %v", err)
			}

			err := service.Delete(ctx, tt.templateID)

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

