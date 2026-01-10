package scheduler

import (
	"context"
	"time"

	"member-pre/internal/domain/slot"
	"member-pre/internal/domain/store"
	"member-pre/internal/domain/user"
	"member-pre/pkg/logger"
)

// SlotSchedulerService 时段调度服务（用于定时任务）
type SlotSchedulerService struct {
	slotService  *slot.SlotService
	storeService *store.StoreService
	userService  *user.UserService
	logger       logger.Logger
}

// NewSlotSchedulerService 创建时段调度服务
func NewSlotSchedulerService(
	slotService *slot.SlotService,
	storeService *store.StoreService,
	userService *user.UserService,
	log logger.Logger,
) *SlotSchedulerService {
	return &SlotSchedulerService{
		slotService:  slotService,
		storeService: storeService,
		userService:  userService,
		logger:       log,
	}
}

// UpdateSlotsForAllStores 为所有营业中的门店更新可预约时间
// 这个方法应该在每天晚上定时执行
func (s *SlotSchedulerService) UpdateSlotsForAllStores(ctx context.Context) error {
	s.logger.Info("开始为所有门店更新可预约时间")

	// 获取所有营业中的门店
	stores, err := s.storeService.GetByStatus(ctx, store.StatusOperating)
	if err != nil {
		s.logger.Error("获取营业中门店列表失败", logger.NewField("error", err.Error()))
		return err
	}

	if len(stores) == 0 {
		s.logger.Info("没有营业中的门店，跳过更新")
		return nil
	}

	s.logger.Info("找到营业中门店", logger.NewField("count", len(stores)))

	// 计算日期：明天
	tomorrow := time.Now().AddDate(0, 0, 1)
	// 将时间设置为00:00:00
	tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())

	successCount := 0
	failCount := 0
	skipCount := 0

	// 为每个门店更新时段
	for _, storeEntity := range stores {
		storeID := storeEntity.ID

		// 检查门店是否配置了模板
		if storeEntity.TemplateID == nil {
			s.logger.Warn("门店未配置模板，跳过",
				logger.NewField("store_id", storeID),
				logger.NewField("store_name", storeEntity.Name),
			)
			failCount++
			continue
		}

		// 检查明天是否已经有时段（且有已预约或已锁定的）
		existingSlots, err := s.slotService.GetByStoreIDAndDate(ctx, storeID, tomorrow)
		if err != nil {
			s.logger.Warn("检查现有时段失败",
				logger.NewField("store_id", storeID),
				logger.NewField("error", err.Error()),
			)
			failCount++
			continue
		}

		// 如果明天已有时段且有已预约或已锁定的记录，跳过同步
		hasActiveBookings := false
		for _, slot := range existingSlots {
			if slot.BookedCount > 0 || slot.LockedCount > 0 {
				hasActiveBookings = true
				break
			}
		}
		if hasActiveBookings {
			s.logger.Info("明天已有预约，跳过同步",
				logger.NewField("store_id", storeID),
				logger.NewField("store_name", storeEntity.Name),
				logger.NewField("date", tomorrow.Format("2006-01-02")),
			)
			skipCount++
			continue
		}

		// 获取门店在岗美甲师数量
		technicians, err := s.userService.GetByStoreID(ctx, storeID, user.RoleTechnician)
		if err != nil {
			s.logger.Warn("获取门店美甲师列表失败",
				logger.NewField("store_id", storeID),
				logger.NewField("error", err.Error()),
			)
			failCount++
			continue
		}

		// 统计在岗美甲师数量
		workingCount := 0
		for _, tech := range technicians {
			if tech.WorkStatus != nil && *tech.WorkStatus == user.WorkStatusWorking {
				workingCount++
			}
		}

		// 根据门店配置的模板ID生成明天的新时段（内部会检测是否已存在时段）
		if err := s.slotService.GenerateSlotsByTemplateID(ctx, storeID, *storeEntity.TemplateID, tomorrow, tomorrow, workingCount); err != nil {
			s.logger.Error("生成时段失败",
				logger.NewField("store_id", storeID),
				logger.NewField("template_id", *storeEntity.TemplateID),
				logger.NewField("error", err.Error()),
			)
			failCount++
			continue
		}

		s.logger.Info("门店时段更新成功",
			logger.NewField("store_id", storeID),
			logger.NewField("store_name", storeEntity.Name),
			logger.NewField("template_id", *storeEntity.TemplateID),
			logger.NewField("technician_count", workingCount),
			logger.NewField("date", tomorrow.Format("2006-01-02")),
		)
		successCount++
	}

	s.logger.Info("门店时段更新完成",
		logger.NewField("total", len(stores)),
		logger.NewField("success", successCount),
		logger.NewField("skipped", skipCount),
		logger.NewField("failed", failCount),
	)

	return nil
}

// deleteOldSlots 删除旧时段（只删除可用状态的时段）
func (s *SlotSchedulerService) deleteOldSlots(ctx context.Context, storeID uint, startDate, endDate time.Time) error {
	// 查找该日期范围内的所有时段
	slots, err := s.slotService.GetByStoreIDAndDateRange(ctx, storeID, startDate, endDate)
	if err != nil {
		return err
	}

	// 只删除可用状态的时段（且没有已预约和已锁定的）
	slotsToDelete := make([]uint, 0)
	for _, slotEntity := range slots {
		if slotEntity.Status == slot.SlotStatusAvailable && slotEntity.BookedCount == 0 && slotEntity.LockedCount == 0 {
			slotsToDelete = append(slotsToDelete, slotEntity.ID)
		}
	}

	if len(slotsToDelete) == 0 {
		return nil
	}

	s.logger.Info("删除旧时段",
		logger.NewField("store_id", storeID),
		logger.NewField("count", len(slotsToDelete)),
	)

	// 批量删除旧时段
	if err := s.slotService.DeleteByIDs(ctx, slotsToDelete); err != nil {
		s.logger.Error("删除旧时段失败",
			logger.NewField("store_id", storeID),
			logger.NewField("error", err.Error()),
		)
		return err
	}

	return nil
}
