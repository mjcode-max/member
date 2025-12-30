package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"member-pre/internal/infrastructure"
	"member-pre/internal/infrastructure/persistence/migrate"
)

// migrateCmd 数据库迁移命令
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "数据库迁移",
	Long:  "执行数据库迁移操作",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "执行迁移",
	Long:  "执行所有待执行的数据库迁移",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := getConfigPath()

		// 初始化应用（只需要配置和数据库）
		app, err := infrastructure.InitializeApp(cfgPath)
		if err != nil {
			return fmt.Errorf("初始化应用失败: %w", err)
		}
		defer app.DB.Close()

		app.Logger.Info("开始执行数据库迁移")

		if err := migrate.Up(app.DB.DB, app.Logger); err != nil {
			app.Logger.Error("数据库迁移失败", zap.Error(err))
			return err
		}

		app.Logger.Info("数据库迁移完成")
		return nil
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "回滚迁移",
	Long:  "回滚最后一次数据库迁移",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := getConfigPath()

		// 初始化应用（只需要配置和数据库）
		app, err := infrastructure.InitializeApp(cfgPath)
		if err != nil {
			return fmt.Errorf("初始化应用失败: %w", err)
		}
		defer app.DB.Close()

		app.Logger.Info("开始回滚数据库迁移")

		if err := migrate.Down(app.DB.DB, app.Logger); err != nil {
			app.Logger.Error("数据库迁移回滚失败", zap.Error(err))
			return err
		}

		app.Logger.Info("数据库迁移回滚完成")
		return nil
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看迁移状态",
	Long:  "查看数据库迁移的当前状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := getConfigPath()

		// 初始化应用（只需要配置和数据库）
		app, err := infrastructure.InitializeApp(cfgPath)
		if err != nil {
			return fmt.Errorf("初始化应用失败: %w", err)
		}
		defer app.DB.Close()

		if err := migrate.Status(app.DB.DB, app.Logger); err != nil {
			app.Logger.Error("查看迁移状态失败", zap.Error(err))
			return err
		}

		return nil
	},
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
}
