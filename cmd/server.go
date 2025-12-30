package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"member-pre/internal/infrastructure"
)

// serverCmd 启动服务器命令
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动HTTP服务器",
	Long:  "启动HTTP服务器，提供API服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := getConfigPath()

		// 初始化应用
		app, err := infrastructure.InitializeApp(cfgPath)
		if err != nil {
			return err
		}

		// 优雅关闭
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// 启动服务器
		go func() {
			app.Logger.Info("启动HTTP服务器",
				zap.Int("port", app.Config.Server.Port),
				zap.String("mode", app.Config.Server.Mode),
			)

			if err := app.Server.Start(); err != nil {
				app.Logger.Fatal("HTTP服务器启动失败", zap.Error(err))
			}
		}()

		// 等待中断信号
		<-quit
		app.Logger.Info("正在关闭服务器...")

		// 优雅关闭
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// 关闭HTTP服务器
		if err := app.Server.Stop(ctx); err != nil {
			app.Logger.Error("关闭HTTP服务器失败", zap.Error(err))
		}

		// 关闭数据库连接
		if err := app.DB.Close(); err != nil {
			app.Logger.Error("关闭数据库连接失败", zap.Error(err))
		}

		// 关闭Redis连接
		if err := app.Redis.Close(); err != nil {
			app.Logger.Error("关闭Redis连接失败", zap.Error(err))
		}

		// 同步日志
		if err := app.Logger.Sync(); err != nil {
			app.Logger.Error("同步日志失败", zap.Error(err))
		}

		app.Logger.Info("服务器已关闭")
		return nil
	},
}
