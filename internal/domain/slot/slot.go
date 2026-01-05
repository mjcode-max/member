package slot

import (
	"context"
	"time"

	"member-pre/pkg/errors"
	"member-pre/pkg/logger"
)

// 时段状态常量
const (
	SlotStatusAvailable = "available" // 可用
	SlotStatusLocked    = "locked"    // 已锁定（预约中）
	SlotStatusBooked    = "booked"    // 已预约
	SlotStatusCompleted = "completed" // 已完成
	SlotStatusCancelled = "cancelled" // 已取消
)

// Slot 时段实体
type Slot struct {
	ID           uint      `json:"id"`
	StoreID      uint      `json:"store_id"`       // 门店ID
	TechnicianID *uint     `json:"technician_id"`  // 美甲师ID（可选，如果为空表示任意美甲师）
	Date         time.Time `json:"date"`            // 日期
	StartTime    time.Time `json:"start_time"`      // 开始时间
	EndTime      time.Time `json:"end_time"`        // 结束时间
	Capacity     int       `json:"capacity"`        // 可用容量（在岗美甲师数量）
	LockedCount  int       `json:"locked_count"`    // 已锁定数量
	BookedCount  int       `json:"booked_count"`    // 已预约数量
	Status       string    `json:"status"`          // 状态: available, locked, booked, completed, cancelled
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ISlotRepository 时段仓储接口
type ISlotRepository interface {
	// FindByID 根据ID查找时段
	FindByID(ctx context.Context, id uint) (*Slot, error)
	// FindByStoreIDAndDate 根据门店ID和日期查找时段列表
	FindByStoreIDAndDate(ctx context.Context, storeID uint, date time.Time) ([]*Slot, error)
	// FindByStoreIDAndDateRange 根据门店ID和日期范围查找时段列表
	FindByStoreIDAndDateRange(ctx context.Context, storeID uint, startDate, endDate time.Time) ([]*Slot, error)
	// FindAvailableByStoreIDAndDate 根据门店ID和日期查找可用时段列表
	FindAvailableByStoreIDAndDate(ctx context.Context, storeID uint, date time.Time) ([]*Slot, error)
	// Create 创建时段
	Create(ctx context.Context, slot *Slot) error
	// CreateBatch 批量创建时段
	CreateBatch(ctx context.Context, slots []*Slot) error
	// Update 更新时段
	Update(ctx context.Context, slot *Slot) error
	// UpdateBatch 批量更新时段
	UpdateBatch(ctx context.Context, slots []*Slot) error
	// LockSlot 锁定时段（原子操作，增加锁定数量）
	LockSlot(ctx context.Context, slotID uint, count int) error
	// UnlockSlot 解锁时段（原子操作，减少锁定数量）
	UnlockSlot(ctx context.Context, slotID uint, count int) error
	// BookSlot 预约时段（原子操作，从锁定转为已预约）
	BookSlot(ctx context.Context, slotID uint, count int) error
	// ReleaseSlot 释放时段（原子操作，取消或完成时释放）
	ReleaseSlot(ctx context.Context, slotID uint, count int) error
	// DeleteByTechnicianIDAndFuture 删除指定美甲师的所有未来时段
	DeleteByTechnicianIDAndFuture(ctx context.Context, technicianID uint, fromDate time.Time) error
	// FindByTechnicianIDAndDateRange 根据美甲师ID和日期范围查找时段列表
	FindByTechnicianIDAndDateRange(ctx context.Context, technicianID uint, startDate, endDate time.Time) ([]*Slot, error)
}

// 领域错误定义
var (
	ErrSlotNotFound        = errors.ErrNotFound("时段不存在")
	ErrSlotNotAvailable   = errors.ErrInvalidParams("时段不可用")
	ErrInsufficientCapacity = errors.ErrInvalidParams("时段容量不足")
	ErrInvalidSlotDate     = errors.ErrInvalidParams("无效的时段日期")
	ErrInvalidSlotTime     = errors.ErrInvalidParams("无效的时段时间")
)

// IUserRepositoryForSlot 用于时段服务获取美甲师列表的接口（避免循环依赖）
type IUserRepositoryForSlot interface {
	FindByStoreID(ctx context.Context, storeID uint, role string) ([]interface{}, error)
}

// SlotService 时段服务
type SlotService struct {
	slotRepo     ISlotRepository
	templateRepo ITemplateRepository
	logger       logger.Logger
}

// NewSlotService 创建时段服务
func NewSlotService(
	slotRepo ISlotRepository,
	templateRepo ITemplateRepository,
	log logger.Logger,
) *SlotService {
	return &SlotService{
		slotRepo:     slotRepo,
		templateRepo: templateRepo,
		logger:       log,
	}
}

// GetByID 根据ID获取时段
func (s *SlotService) GetByID(ctx context.Context, id uint) (*Slot, error) {
	s.logger.Debug("获取时段", logger.NewField("slot_id", id))

	slot, err := s.slotRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("查找时段失败", logger.NewField("slot_id", id), logger.NewField("error", err.Error()))
		return nil, err
	}

	if slot == nil {
		s.logger.Warn("时段不存在", logger.NewField("slot_id", id))
		return nil, ErrSlotNotFound
	}

	return slot, nil
}

// GetByStoreIDAndDate 根据门店ID和日期获取时段列表
func (s *SlotService) GetByStoreIDAndDate(ctx context.Context, storeID uint, date time.Time) ([]*Slot, error) {
	s.logger.Debug("获取门店时段列表",
		logger.NewField("store_id", storeID),
		logger.NewField("date", date.Format("2006-01-02")),
	)

	slots, err := s.slotRepo.FindByStoreIDAndDate(ctx, storeID, date)
	if err != nil {
		s.logger.Error("查找时段列表失败", logger.NewField("store_id", storeID), logger.NewField("error", err.Error()))
		return nil, err
	}

	return slots, nil
}

// GetAvailableByStoreIDAndDate 根据门店ID和日期获取可用时段列表
func (s *SlotService) GetAvailableByStoreIDAndDate(ctx context.Context, storeID uint, date time.Time) ([]*Slot, error) {
	s.logger.Debug("获取门店可用时段列表",
		logger.NewField("store_id", storeID),
		logger.NewField("date", date.Format("2006-01-02")),
	)

	slots, err := s.slotRepo.FindAvailableByStoreIDAndDate(ctx, storeID, date)
	if err != nil {
		s.logger.Error("查找可用时段列表失败", logger.NewField("store_id", storeID), logger.NewField("error", err.Error()))
		return nil, err
	}

	return slots, nil
}

// GenerateSlots 根据模板生成时段（用于初始化或重新生成）
func (s *SlotService) GenerateSlots(ctx context.Context, storeID uint, startDate, endDate time.Time, technicianCount int) error {
	s.logger.Info("生成时段",
		logger.NewField("store_id", storeID),
		logger.NewField("start_date", startDate.Format("2006-01-02")),
		logger.NewField("end_date", endDate.Format("2006-01-02")),
		logger.NewField("technician_count", technicianCount),
	)

	// 获取启用的模板
	template, err := s.templateRepo.FindActiveByStoreID(ctx, storeID)
	if err != nil {
		return err
	}
	if template == nil {
		return ErrTemplateNotFound
	}

	// 生成时段列表
	slots := make([]*Slot, 0)
	currentDate := startDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		weekday := int(currentDate.Weekday())
		// 查找该天的规则
		for _, rule := range template.WeekdayRules {
			if rule.Weekday == weekday {
				// 为该天的每个时段规则生成时段
				for _, timeSlot := range rule.Slots {
					startTime, err := parseTime(currentDate, timeSlot.StartTime)
					if err != nil {
						s.logger.Warn("解析开始时间失败", logger.NewField("time", timeSlot.StartTime), logger.NewField("error", err.Error()))
						continue
					}
					endTime, err := parseTime(currentDate, timeSlot.EndTime)
					if err != nil {
						s.logger.Warn("解析结束时间失败", logger.NewField("time", timeSlot.EndTime), logger.NewField("error", err.Error()))
						continue
					}

					slot := &Slot{
						StoreID:     storeID,
						TechnicianID: nil, // 不指定具体美甲师
						Date:        currentDate,
						StartTime:   startTime,
						EndTime:     endTime,
						Capacity:    technicianCount, // 基于在岗美甲师数量
						LockedCount: 0,
						BookedCount: 0,
						Status:      SlotStatusAvailable,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}
					slots = append(slots, slot)
				}
				break
			}
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	if len(slots) == 0 {
		s.logger.Warn("没有生成任何时段")
		return nil
	}

	// 批量创建时段
	if err := s.slotRepo.CreateBatch(ctx, slots); err != nil {
		s.logger.Error("批量创建时段失败", logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("生成时段成功", logger.NewField("count", len(slots)))
	return nil
}

// RecalculateCapacity 重新计算时段容量（当美甲师状态变更时调用）
func (s *SlotService) RecalculateCapacity(ctx context.Context, storeID uint, fromDate time.Time, technicianCount int) error {
	s.logger.Info("重新计算时段容量",
		logger.NewField("store_id", storeID),
		logger.NewField("from_date", fromDate.Format("2006-01-02")),
		logger.NewField("technician_count", technicianCount),
	)

	// 查找所有未来的时段
	endDate := fromDate.AddDate(0, 0, 90) // 未来90天
	slots, err := s.slotRepo.FindByStoreIDAndDateRange(ctx, storeID, fromDate, endDate)
	if err != nil {
		return err
	}

	// 更新容量
	slotsToUpdate := make([]*Slot, 0)
	for _, slot := range slots {
		// 只更新可用状态的时段，已预约的时段不改变容量
		if slot.Status == SlotStatusAvailable {
			oldCapacity := slot.Capacity
			slot.Capacity = technicianCount
			// 如果新容量小于已锁定数量，需要调整
			if slot.Capacity < slot.LockedCount {
				slot.LockedCount = slot.Capacity
			}
			if oldCapacity != slot.Capacity {
				slot.UpdatedAt = time.Now()
				slotsToUpdate = append(slotsToUpdate, slot)
			}
		}
	}

	if len(slotsToUpdate) > 0 {
		if err := s.slotRepo.UpdateBatch(ctx, slotsToUpdate); err != nil {
			s.logger.Error("批量更新时段容量失败", logger.NewField("error", err.Error()))
			return err
		}
		s.logger.Info("重新计算时段容量成功", logger.NewField("updated_count", len(slotsToUpdate)))
	}

	return nil
}

// LockSlot 锁定时段（预约时调用，原子操作）
func (s *SlotService) LockSlot(ctx context.Context, slotID uint, count int) error {
	s.logger.Info("锁定时段",
		logger.NewField("slot_id", slotID),
		logger.NewField("count", count),
	)

	if count <= 0 {
		return errors.ErrInvalidParams("锁定数量必须大于0")
	}

	// 获取时段信息
	slot, err := s.slotRepo.FindByID(ctx, slotID)
	if err != nil {
		return err
	}
	if slot == nil {
		return ErrSlotNotFound
	}

	// 检查可用容量
	available := slot.Capacity - slot.LockedCount - slot.BookedCount
	if available < count {
		return ErrInsufficientCapacity
	}

	// 原子操作：增加锁定数量
	if err := s.slotRepo.LockSlot(ctx, slotID, count); err != nil {
		s.logger.Error("锁定时段失败", logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("锁定时段成功", logger.NewField("slot_id", slotID))
	return nil
}

// UnlockSlot 解锁时段（取消预约时调用，原子操作）
func (s *SlotService) UnlockSlot(ctx context.Context, slotID uint, count int) error {
	s.logger.Info("解锁时段",
		logger.NewField("slot_id", slotID),
		logger.NewField("count", count),
	)

	if count <= 0 {
		return errors.ErrInvalidParams("解锁数量必须大于0")
	}

	// 原子操作：减少锁定数量
	if err := s.slotRepo.UnlockSlot(ctx, slotID, count); err != nil {
		s.logger.Error("解锁时段失败", logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("解锁时段成功", logger.NewField("slot_id", slotID))
	return nil
}

// BookSlot 预约时段（从锁定转为已预约，原子操作）
func (s *SlotService) BookSlot(ctx context.Context, slotID uint, count int) error {
	s.logger.Info("预约时段",
		logger.NewField("slot_id", slotID),
		logger.NewField("count", count),
	)

	if count <= 0 {
		return errors.ErrInvalidParams("预约数量必须大于0")
	}

	// 原子操作：从锁定转为已预约
	if err := s.slotRepo.BookSlot(ctx, slotID, count); err != nil {
		s.logger.Error("预约时段失败", logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("预约时段成功", logger.NewField("slot_id", slotID))
	return nil
}

// ReleaseSlot 释放时段（取消或完成时调用，原子操作）
func (s *SlotService) ReleaseSlot(ctx context.Context, slotID uint, count int) error {
	s.logger.Info("释放时段",
		logger.NewField("slot_id", slotID),
		logger.NewField("count", count),
	)

	if count <= 0 {
		return errors.ErrInvalidParams("释放数量必须大于0")
	}

	// 原子操作：释放时段
	if err := s.slotRepo.ReleaseSlot(ctx, slotID, count); err != nil {
		s.logger.Error("释放时段失败", logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("释放时段成功", logger.NewField("slot_id", slotID))
	return nil
}

// ReleaseTechnicianSlots 释放美甲师的所有未来时段（当美甲师状态变更为休息时调用）
func (s *SlotService) ReleaseTechnicianSlots(ctx context.Context, technicianID uint, fromDate time.Time) error {
	s.logger.Info("释放美甲师未来时段",
		logger.NewField("technician_id", technicianID),
		logger.NewField("from_date", fromDate.Format("2006-01-02")),
	)

	// 查找美甲师的所有未来时段
	endDate := fromDate.AddDate(0, 0, 90) // 未来90天
	slots, err := s.slotRepo.FindByTechnicianIDAndDateRange(ctx, technicianID, fromDate, endDate)
	if err != nil {
		return err
	}

	// 释放所有已锁定和已预约的时段
	for _, slot := range slots {
		if slot.LockedCount > 0 {
			if err := s.slotRepo.UnlockSlot(ctx, slot.ID, slot.LockedCount); err != nil {
				s.logger.Warn("解锁时段失败", logger.NewField("slot_id", slot.ID), logger.NewField("error", err.Error()))
			}
		}
		if slot.BookedCount > 0 {
			if err := s.slotRepo.ReleaseSlot(ctx, slot.ID, slot.BookedCount); err != nil {
				s.logger.Warn("释放时段失败", logger.NewField("slot_id", slot.ID), logger.NewField("error", err.Error()))
			}
		}
	}

	// 删除美甲师的未来时段记录
	if err := s.slotRepo.DeleteByTechnicianIDAndFuture(ctx, technicianID, fromDate); err != nil {
		s.logger.Error("删除美甲师未来时段失败", logger.NewField("error", err.Error()))
		return err
	}

	s.logger.Info("释放美甲师未来时段成功", logger.NewField("technician_id", technicianID))
	return nil
}

// parseTime 解析时间字符串并组合日期
func parseTime(date time.Time, timeStr string) (time.Time, error) {
	// 解析时间字符串 (HH:MM)
	layout := "15:04"
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return time.Time{}, err
	}

	// 组合日期和时间
	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		parsedTime.Hour(),
		parsedTime.Minute(),
		0,
		0,
		date.Location(),
	), nil
}

