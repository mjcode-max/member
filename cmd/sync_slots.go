package main

import (
	"context"
	"errors"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"member-pre/internal/infrastructure"
	infraLogger "member-pre/internal/infrastructure/logger"
	"member-pre/pkg/logger"
)

// syncSlotsCmd 同步时段命令
var syncSlotsCmd = &cobra.Command{
	Use:   "sync-slots",
	Short: "同步门店可预约时段",
	Long:  `为所有营业中的门店更新可预约时间，根据门店配置的模板生成未来30天的时段`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSyncSlots()
	},
}

func runSyncSlots() error {
	cfgPath := getConfigPath()

	// 初始化日志
	log, err := infraLogger.NewZapLogger("info", "console", "stdout", "")
	if err != nil {
		return err
	}
	defer log.Sync()

	log.Info("开始同步门店时段...")

	// 初始化应用
	app, err := infrastructure.InitializeApp(infrastructure.ConfigPath(cfgPath))
	if err != nil {
		log.Error("初始化应用失败", logger.NewField("error", err.Error()))
		return err
	}
	defer func() {
		// 关闭数据库连接
		if err := app.DB.Close(); err != nil {
			log.Error("关闭数据库连接失败", logger.NewField("error", err.Error()))
		}
		// 关闭Redis连接
		if err := app.Redis.Close(); err != nil {
			log.Error("关闭Redis连接失败", logger.NewField("error", err.Error()))
		}
		// 同步日志
		if err := app.Logger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			log.Error("同步日志失败", logger.NewField("error", err.Error()))
		}
	}()

	// 执行同步
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := app.SlotSchedulerService.UpdateSlotsForAllStores(ctx); err != nil {
		log.Error("同步门店时段失败", logger.NewField("error", err.Error()))
		return err
	}

	log.Info("同步门店时段完成")
	return nil
}

func init() {
	rootCmd.AddCommand(syncSlotsCmd)
}
