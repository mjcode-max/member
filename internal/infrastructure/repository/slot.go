package repository

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
	"member-pre/internal/domain/slot"
	"member-pre/internal/infrastructure/persistence"
	"member-pre/internal/infrastructure/persistence/database"
	"member-pre/pkg/errors"
	"member-pre/pkg/logger"
)

var _ slot.ITemplateRepository = (*TemplateRepository)(nil)
var _ slot.ISlotRepository = (*SlotRepository)(nil)

// 在包初始化时注册模型，确保迁移时能检测到
func init() {
	persistence.Register(&TemplateModel{})
	persistence.Register(&SlotModel{})
}

// ==================== 时段模板仓储 ====================

// TemplateRepository 时段模板仓储实现
type TemplateRepository struct {
	db     database.Database
	redis  *persistence.Client
	logger logger.Logger
}

// NewTemplateRepository 创建时段模板仓储实例
func NewTemplateRepository(db database.Database, rdb *persistence.Client, log logger.Logger) *TemplateRepository {
	return &TemplateRepository{
		db:     db,
		redis:  rdb,
		logger: log,
	}
}

// TimeSlotRuleJSON 时间段规则JSON类型（用于数据库存储）
type TimeSlotRuleJSON []slot.TimeSlotRule

// Value 实现driver.Valuer接口
func (t TimeSlotRuleJSON) Value() (driver.Value, error) {
	if len(t) == 0 {
		return "[]", nil
	}
	return json.Marshal(t)
}

// Scan 实现sql.Scanner接口
func (t *TimeSlotRuleJSON) Scan(value interface{}) error {
	if value == nil {
		*t = TimeSlotRuleJSON{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.ErrInvalidParams("无法扫描TimeSlotRuleJSON")
	}
	return json.Unmarshal(bytes, t)
}

// TemplateModel 时段模板数据库模型
type TemplateModel struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	Name      string           `gorm:"size:100;not null" json:"name"`
	Status    string           `gorm:"size:20;default:'active';not null" json:"status"`
	TimeSlots TimeSlotRuleJSON `gorm:"column:time_slots;type:json" json:"time_slots"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// TableName 指定表名
func (TemplateModel) TableName() string {
	return "slot_templates"
}

// ToEntity 转换为领域实体
func (m *TemplateModel) ToEntity() *slot.Template {
	if m == nil {
		return nil
	}
	return &slot.Template{
		ID:        m.ID,
		Name:      m.Name,
		Status:    m.Status,
		TimeSlots: []slot.TimeSlotRule(m.TimeSlots),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromEntity 从领域实体创建模型
func (m *TemplateModel) FromEntity(t *slot.Template) {
	if t == nil {
		return
	}
	m.ID = t.ID
	m.Name = t.Name
	m.Status = t.Status
	m.TimeSlots = TimeSlotRuleJSON(t.TimeSlots)
	m.CreatedAt = t.CreatedAt
	m.UpdatedAt = t.UpdatedAt
}

// FindByID 根据ID查找模板
func (r *TemplateRepository) FindByID(ctx context.Context, id uint) (*slot.Template, error) {
	r.logger.Debug("查找时段模板：根据ID", logger.NewField("template_id", id))

	var model TemplateModel
	if err := r.db.DB().WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("时段模板不存在：根据ID", logger.NewField("template_id", id))
			return nil, nil
		}
		r.logger.Error("查找时段模板失败：根据ID",
			logger.NewField("template_id", id),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	return model.ToEntity(), nil
}

// FindAll 查找所有模板列表
func (r *TemplateRepository) FindAll(ctx context.Context) ([]*slot.Template, error) {
	r.logger.Debug("查找所有时段模板列表")

	var models []TemplateModel
	if err := r.db.DB().WithContext(ctx).Find(&models).Error; err != nil {
		r.logger.Error("查找所有时段模板列表失败",
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	templates := make([]*slot.Template, 0, len(models))
	for _, model := range models {
		templates = append(templates, model.ToEntity())
	}

	return templates, nil
}

// FindByStatus 根据状态查找模板列表
func (r *TemplateRepository) FindByStatus(ctx context.Context, status string) ([]*slot.Template, error) {
	r.logger.Debug("查找时段模板列表：根据状态", logger.NewField("status", status))

	var models []TemplateModel
	if err := r.db.DB().WithContext(ctx).Where("status = ?", status).Find(&models).Error; err != nil {
		r.logger.Error("查找时段模板列表失败：根据状态",
			logger.NewField("status", status),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	templates := make([]*slot.Template, 0, len(models))
	for _, model := range models {
		templates = append(templates, model.ToEntity())
	}

	return templates, nil
}

// Create 创建模板
func (r *TemplateRepository) Create(ctx context.Context, t *slot.Template) error {
	r.logger.Info("创建时段模板",
		logger.NewField("name", t.Name),
	)

	model := &TemplateModel{}
	model.FromEntity(t)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Error("创建时段模板失败",
			logger.NewField("template_id", model.ID),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	t.ID = model.ID
	t.CreatedAt = model.CreatedAt
	t.UpdatedAt = model.UpdatedAt

	r.logger.Info("创建时段模板成功",
		logger.NewField("template_id", t.ID),
		logger.NewField("name", t.Name),
	)

	return nil
}

// Update 更新模板
func (r *TemplateRepository) Update(ctx context.Context, t *slot.Template) error {
	r.logger.Info("更新时段模板", logger.NewField("template_id", t.ID))

	model := &TemplateModel{}
	model.FromEntity(t)
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Model(&TemplateModel{}).Where("id = ?", t.ID).Updates(model).Error; err != nil {
		r.logger.Error("更新时段模板失败",
			logger.NewField("template_id", t.ID),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	t.UpdatedAt = model.UpdatedAt

	r.logger.Info("更新时段模板成功", logger.NewField("template_id", t.ID))
	return nil
}

// Delete 删除模板
func (r *TemplateRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("删除时段模板", logger.NewField("template_id", id))

	if err := r.db.DB().WithContext(ctx).Delete(&TemplateModel{}, id).Error; err != nil {
		r.logger.Error("删除时段模板失败",
			logger.NewField("template_id", id),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	r.logger.Info("删除时段模板成功", logger.NewField("template_id", id))
	return nil
}

// ==================== 时段仓储 ====================

// SlotRepository 时段仓储实现
type SlotRepository struct {
	db     database.Database
	redis  *persistence.Client
	logger logger.Logger
}

// NewSlotRepository 创建时段仓储实例
func NewSlotRepository(db database.Database, rdb *persistence.Client, log logger.Logger) *SlotRepository {
	return &SlotRepository{
		db:     db,
		redis:  rdb,
		logger: log,
	}
}

// SlotModel 时段数据库模型
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
func (m *SlotModel) ToEntity() *slot.Slot {
	if m == nil {
		return nil
	}
	return &slot.Slot{
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
func (m *SlotModel) FromEntity(s *slot.Slot) {
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

// FindByID 根据ID查找时段
func (r *SlotRepository) FindByID(ctx context.Context, id uint) (*slot.Slot, error) {
	r.logger.Debug("查找时段：根据ID", logger.NewField("slot_id", id))

	var model SlotModel
	if err := r.db.DB().WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("时段不存在：根据ID", logger.NewField("slot_id", id))
			return nil, nil
		}
		r.logger.Error("查找时段失败：根据ID",
			logger.NewField("slot_id", id),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	return model.ToEntity(), nil
}

// FindByStoreIDAndDate 根据门店ID和日期查找时段列表
func (r *SlotRepository) FindByStoreIDAndDate(ctx context.Context, storeID uint, date time.Time) ([]*slot.Slot, error) {
	r.logger.Debug("查找时段列表：根据门店ID和日期",
		logger.NewField("store_id", storeID),
		logger.NewField("date", date.Format("2006-01-02")),
	)

	var models []SlotModel
	dateStr := date.Format("2006-01-02")
	if err := r.db.DB().WithContext(ctx).
		Where("store_id = ? AND DATE(date) = ?", storeID, dateStr).
		Order("start_time ASC").
		Find(&models).Error; err != nil {
		r.logger.Error("查找时段列表失败：根据门店ID和日期",
			logger.NewField("store_id", storeID),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	slots := make([]*slot.Slot, 0, len(models))
	for _, model := range models {
		slots = append(slots, model.ToEntity())
	}

	return slots, nil
}

// FindByStoreIDAndDateRange 根据门店ID和日期范围查找时段列表
func (r *SlotRepository) FindByStoreIDAndDateRange(ctx context.Context, storeID uint, startDate, endDate time.Time) ([]*slot.Slot, error) {
	r.logger.Debug("查找时段列表：根据门店ID和日期范围",
		logger.NewField("store_id", storeID),
		logger.NewField("start_date", startDate.Format("2006-01-02")),
		logger.NewField("end_date", endDate.Format("2006-01-02")),
	)

	var models []SlotModel
	if err := r.db.DB().WithContext(ctx).
		Where("store_id = ? AND date >= ? AND date <= ?", storeID, startDate, endDate).
		Order("date ASC, start_time ASC").
		Find(&models).Error; err != nil {
		r.logger.Error("查找时段列表失败：根据门店ID和日期范围",
			logger.NewField("store_id", storeID),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	slots := make([]*slot.Slot, 0, len(models))
	for _, model := range models {
		slots = append(slots, model.ToEntity())
	}

	return slots, nil
}

// FindAvailableByStoreIDAndDate 根据门店ID和日期查找可用时段列表
func (r *SlotRepository) FindAvailableByStoreIDAndDate(ctx context.Context, storeID uint, date time.Time) ([]*slot.Slot, error) {
	r.logger.Debug("查找可用时段列表：根据门店ID和日期",
		logger.NewField("store_id", storeID),
		logger.NewField("date", date.Format("2006-01-02")),
	)

	var models []SlotModel
	dateStr := date.Format("2006-01-02")
	if err := r.db.DB().WithContext(ctx).
		Where("store_id = ? AND DATE(date) = ? AND status = ? AND (capacity - locked_count - booked_count) > 0",
			storeID, dateStr, slot.SlotStatusAvailable).
		Order("start_time ASC").
		Find(&models).Error; err != nil {
		r.logger.Error("查找可用时段列表失败：根据门店ID和日期",
			logger.NewField("store_id", storeID),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	slots := make([]*slot.Slot, 0, len(models))
	for _, model := range models {
		slots = append(slots, model.ToEntity())
	}

	return slots, nil
}

// Create 创建时段
func (r *SlotRepository) Create(ctx context.Context, s *slot.Slot) error {
	r.logger.Info("创建时段",
		logger.NewField("store_id", s.StoreID),
		logger.NewField("date", s.Date.Format("2006-01-02")),
	)

	model := &SlotModel{}
	model.FromEntity(s)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Error("创建时段失败",
			logger.NewField("store_id", s.StoreID),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	s.ID = model.ID
	s.CreatedAt = model.CreatedAt
	s.UpdatedAt = model.UpdatedAt

	return nil
}

// CreateBatch 批量创建时段
func (r *SlotRepository) CreateBatch(ctx context.Context, slots []*slot.Slot) error {
	if len(slots) == 0 {
		return nil
	}

	r.logger.Info("批量创建时段", logger.NewField("count", len(slots)))

	models := make([]SlotModel, 0, len(slots))
	now := time.Now()
	for _, s := range slots {
		model := &SlotModel{}
		model.FromEntity(s)
		model.CreatedAt = now
		model.UpdatedAt = now
		models = append(models, *model)
	}

	if err := r.db.DB().WithContext(ctx).CreateInBatches(models, 100).Error; err != nil {
		r.logger.Error("批量创建时段失败", logger.NewField("error", err.Error()))
		return errors.ErrDatabase(err)
	}

	// 更新实体的ID和时间戳
	for i, model := range models {
		slots[i].ID = model.ID
		slots[i].CreatedAt = model.CreatedAt
		slots[i].UpdatedAt = model.UpdatedAt
	}

	r.logger.Info("批量创建时段成功", logger.NewField("count", len(slots)))
	return nil
}

// Update 更新时段
func (r *SlotRepository) Update(ctx context.Context, s *slot.Slot) error {
	r.logger.Info("更新时段", logger.NewField("slot_id", s.ID))

	model := &SlotModel{}
	model.FromEntity(s)
	model.UpdatedAt = time.Now()

	if err := r.db.DB().WithContext(ctx).Model(&SlotModel{}).Where("id = ?", s.ID).Updates(model).Error; err != nil {
		r.logger.Error("更新时段失败",
			logger.NewField("slot_id", s.ID),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	s.UpdatedAt = model.UpdatedAt
	return nil
}

// UpdateBatch 批量更新时段
func (r *SlotRepository) UpdateBatch(ctx context.Context, slots []*slot.Slot) error {
	if len(slots) == 0 {
		return nil
	}

	r.logger.Info("批量更新时段", logger.NewField("count", len(slots)))

	// 使用事务批量更新
	return r.db.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, s := range slots {
			model := &SlotModel{}
			model.FromEntity(s)
			model.UpdatedAt = time.Now()

			if err := tx.Model(&SlotModel{}).Where("id = ?", s.ID).Updates(model).Error; err != nil {
				r.logger.Error("批量更新时段失败",
					logger.NewField("slot_id", s.ID),
					logger.NewField("error", err.Error()),
				)
				return errors.ErrDatabase(err)
			}
		}
		return nil
	})
}

// LockSlot 锁定时段（原子操作，增加锁定数量）
func (r *SlotRepository) LockSlot(ctx context.Context, slotID uint, count int) error {
	r.logger.Info("锁定时段（原子操作）",
		logger.NewField("slot_id", slotID),
		logger.NewField("count", count),
	)

	// 使用事务和行锁确保原子性
	return r.db.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model SlotModel
		// 使用SELECT FOR UPDATE获取行锁
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("id = ?", slotID).
			First(&model).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return slot.ErrSlotNotFound
			}
			return errors.ErrDatabase(err)
		}

		// 检查可用容量
		available := model.Capacity - model.LockedCount - model.BookedCount
		if available < count {
			return slot.ErrInsufficientCapacity
		}

		// 原子更新锁定数量
		if err := tx.Model(&SlotModel{}).
			Where("id = ?", slotID).
			Update("locked_count", gorm.Expr("locked_count + ?", count)).Error; err != nil {
			return errors.ErrDatabase(err)
		}

		return nil
	})
}

// UnlockSlot 解锁时段（原子操作，减少锁定数量）
func (r *SlotRepository) UnlockSlot(ctx context.Context, slotID uint, count int) error {
	r.logger.Info("解锁时段（原子操作）",
		logger.NewField("slot_id", slotID),
		logger.NewField("count", count),
	)

	// 使用事务和行锁确保原子性
	return r.db.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 原子更新锁定数量（确保不会小于0）
		result := tx.Model(&SlotModel{}).
			Where("id = ? AND locked_count >= ?", slotID, count).
			Update("locked_count", gorm.Expr("locked_count - ?", count))

		if result.Error != nil {
			return errors.ErrDatabase(result.Error)
		}

		if result.RowsAffected == 0 {
			return errors.ErrInvalidParams("解锁数量超过当前锁定数量")
		}

		return nil
	})
}

// BookSlot 预约时段（原子操作，从锁定转为已预约）
func (r *SlotRepository) BookSlot(ctx context.Context, slotID uint, count int) error {
	r.logger.Info("预约时段（原子操作）",
		logger.NewField("slot_id", slotID),
		logger.NewField("count", count),
	)

	// 使用事务和行锁确保原子性
	return r.db.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model SlotModel
		// 使用SELECT FOR UPDATE获取行锁
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("id = ?", slotID).
			First(&model).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return slot.ErrSlotNotFound
			}
			return errors.ErrDatabase(err)
		}

		// 检查锁定数量
		if model.LockedCount < count {
			return errors.ErrInvalidParams("预约数量超过当前锁定数量")
		}

		// 原子更新：减少锁定数量，增加已预约数量
		if err := tx.Model(&SlotModel{}).
			Where("id = ?", slotID).
			Updates(map[string]interface{}{
				"locked_count": gorm.Expr("locked_count - ?", count),
				"booked_count": gorm.Expr("booked_count + ?", count),
			}).Error; err != nil {
			return errors.ErrDatabase(err)
		}

		return nil
	})
}

// ReleaseSlot 释放时段（原子操作，取消或完成时释放）
func (r *SlotRepository) ReleaseSlot(ctx context.Context, slotID uint, count int) error {
	r.logger.Info("释放时段（原子操作）",
		logger.NewField("slot_id", slotID),
		logger.NewField("count", count),
	)

	// 使用事务和行锁确保原子性
	return r.db.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model SlotModel
		// 使用SELECT FOR UPDATE获取行锁
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("id = ?", slotID).
			First(&model).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return slot.ErrSlotNotFound
			}
			return errors.ErrDatabase(err)
		}

		// 检查已预约数量
		if model.BookedCount < count {
			return errors.ErrInvalidParams("释放数量超过当前已预约数量")
		}

		// 原子更新：减少已预约数量
		if err := tx.Model(&SlotModel{}).
			Where("id = ?", slotID).
			Update("booked_count", gorm.Expr("booked_count - ?", count)).Error; err != nil {
			return errors.ErrDatabase(err)
		}

		return nil
	})
}

// DeleteByTechnicianIDAndFuture 删除指定美甲师的所有未来时段
func (r *SlotRepository) DeleteByTechnicianIDAndFuture(ctx context.Context, technicianID uint, fromDate time.Time) error {
	r.logger.Info("删除美甲师未来时段",
		logger.NewField("technician_id", technicianID),
		logger.NewField("from_date", fromDate.Format("2006-01-02")),
	)

	if err := r.db.DB().WithContext(ctx).
		Where("technician_id = ? AND date >= ?", technicianID, fromDate).
		Delete(&SlotModel{}).Error; err != nil {
		r.logger.Error("删除美甲师未来时段失败",
			logger.NewField("technician_id", technicianID),
			logger.NewField("error", err.Error()),
		)
		return errors.ErrDatabase(err)
	}

	r.logger.Info("删除美甲师未来时段成功", logger.NewField("technician_id", technicianID))
	return nil
}

// DeleteByIDs 批量删除时段
func (r *SlotRepository) DeleteByIDs(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	r.logger.Info("批量删除时段", logger.NewField("count", len(ids)))

	if err := r.db.DB().WithContext(ctx).Where("id IN ?", ids).Delete(&SlotModel{}).Error; err != nil {
		r.logger.Error("批量删除时段失败", logger.NewField("error", err.Error()))
		return errors.ErrDatabase(err)
	}

	r.logger.Info("批量删除时段成功", logger.NewField("count", len(ids)))
	return nil
}

// FindByTechnicianIDAndDateRange 根据美甲师ID和日期范围查找时段列表
func (r *SlotRepository) FindByTechnicianIDAndDateRange(ctx context.Context, technicianID uint, startDate, endDate time.Time) ([]*slot.Slot, error) {
	r.logger.Debug("查找时段列表：根据美甲师ID和日期范围",
		logger.NewField("technician_id", technicianID),
		logger.NewField("start_date", startDate.Format("2006-01-02")),
		logger.NewField("end_date", endDate.Format("2006-01-02")),
	)

	var models []SlotModel
	if err := r.db.DB().WithContext(ctx).
		Where("technician_id = ? AND date >= ? AND date <= ?", technicianID, startDate, endDate).
		Order("date ASC, start_time ASC").
		Find(&models).Error; err != nil {
		r.logger.Error("查找时段列表失败：根据美甲师ID和日期范围",
			logger.NewField("technician_id", technicianID),
			logger.NewField("error", err.Error()),
		)
		return nil, errors.ErrDatabase(err)
	}

	slots := make([]*slot.Slot, 0, len(models))
	for _, model := range models {
		slots = append(slots, model.ToEntity())
	}

	return slots, nil
}
