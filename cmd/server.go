package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"member-pre/internal/infrastructure"
	"member-pre/pkg/logger"
)

// serverCmd 启动服务器命令
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动HTTP服务器",
	Long:  "启动HTTP服务器，提供API服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := getConfigPath()

		// 初始化应用
		app, err := infrastructure.InitializeApp(infrastructure.ConfigPath(cfgPath))
		if err != nil {
			return err
		}

		// 优雅关闭
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// 启动定时任务
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			app.CronScheduler.Start(ctx)
		}()

		// 启动服务器
		go func() {
			app.Logger.Info("启动HTTP服务器",
				logger.NewField("port", app.Config.Server.Port),
				logger.NewField("mode", app.Config.Server.Mode),
			)

			if err := app.Server.Start(); err != nil {
				app.Logger.Fatal("HTTP服务器启动失败", logger.NewField("error", err.Error()))
			}
		}()

		// 等待中断信号
		<-quit
		app.Logger.Info("正在关闭服务器...")

		// 停止定时任务
		app.CronScheduler.Stop()

		// 优雅关闭
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// 关闭HTTP服务器
		if err := app.Server.Stop(shutdownCtx); err != nil {
			app.Logger.Error("关闭HTTP服务器失败", logger.NewField("error", err.Error()))
		}

		// 关闭数据库连接
		if err := app.DB.Close(); err != nil {
			app.Logger.Error("关闭数据库连接失败", logger.NewField("error", err.Error()))
		}

		// 关闭Redis连接
		if err := app.Redis.Close(); err != nil {
			app.Logger.Error("关闭Redis连接失败", logger.NewField("error", err.Error()))
		}

		// 同步日志
		if err := app.Logger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			app.Logger.Error("同步日志失败", logger.NewField("error", err.Error()))
		}

		app.Logger.Info("服务器已关闭")
		return nil
	},
}
