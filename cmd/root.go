package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configPath string
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "member-pre",
	Short: "会员预约系统",
	Long:  "基于 Go 语言开发的会员预约系统，采用领域驱动设计（DDD）架构",
}

// Execute 执行命令
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "configs/config.yaml", "配置文件路径")

	// 注册子命令
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(migrateCmd)
}

// getConfigPath 获取配置文件路径
func getConfigPath() string {
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("配置文件不存在: %s\n", configPath)
		os.Exit(1)
	}
	return configPath
}
