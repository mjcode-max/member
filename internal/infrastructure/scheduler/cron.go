package scheduler

import (
	"context"
	"time"

	"member-pre/pkg/logger"
)

// CronScheduler 定时任务调度器
type CronScheduler struct {
	slotScheduler *SlotSchedulerService
	logger        logger.Logger
	stopChan      chan struct{}
}

// NewCronScheduler 创建定时任务调度器
func NewCronScheduler(slotScheduler *SlotSchedulerService, log logger.Logger) *CronScheduler {
	return &CronScheduler{
		slotScheduler: slotScheduler,
		logger:        log,
		stopChan:      make(chan struct{}),
	}
}

// Start 启动定时任务
func (c *CronScheduler) Start(ctx context.Context) {
	c.logger.Info("启动定时任务调度器")

	// 计算到明天凌晨的时间
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	duration := tomorrow.Sub(now)

	c.logger.Info("将在明天凌晨执行第一次时段更新任务",
		logger.NewField("duration_hours", duration.Hours()),
	)

	// 等待到明天凌晨
	timer := time.NewTimer(duration)
	select {
	case <-timer.C:
		// 立即执行一次
		c.runSlotUpdateTask(ctx)
	case <-c.stopChan:
		timer.Stop()
		return
	}

	// 之后每天凌晨执行
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.runSlotUpdateTask(ctx)
		case <-c.stopChan:
			c.logger.Info("定时任务调度器已停止")
			return
		}
	}
}

// Stop 停止定时任务
func (c *CronScheduler) Stop() {
	close(c.stopChan)
}

// runSlotUpdateTask 执行时段更新任务
func (c *CronScheduler) runSlotUpdateTask(ctx context.Context) {
	c.logger.Info("开始执行时段更新任务")

	// 创建带超时的上下文（5分钟超时）
	taskCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	if err := c.slotScheduler.UpdateSlotsForAllStores(taskCtx); err != nil {
		c.logger.Error("时段更新任务执行失败", logger.NewField("error", err.Error()))
	} else {
		c.logger.Info("时段更新任务执行成功")
	}
}
