//go:build wireinject
// +build wireinject

package infrastructure

import (
	"member-pre/internal/infrastructure/config"
	"member-pre/internal/infrastructure/http"
	"member-pre/internal/infrastructure/logger"
	"member-pre/internal/infrastructure/persistence/mysql"
	"member-pre/internal/infrastructure/persistence/redis"

	"github.com/google/wire"
)

// InitializeApp 初始化应用
func InitializeApp(cfgPath string) (*App, error) {
	wire.Build(
		// 配置
		config.Load,
		// 日志
		NewLogger,
		// 数据库
		mysql.NewDB,
		// Redis
		redis.NewClient,
		// HTTP服务器（需要完整配置、数据库和Redis）
		http.NewServer,
		// 应用
		NewApp,
	)
	return &App{}, nil
}

// NewLogger 创建日志实例
func NewLogger(cfg *config.Config) (*logger.ZapLogger, error) {
	return logger.NewZapLogger(
		cfg.Log.Level,
		cfg.Log.Format,
		cfg.Log.Output,
		cfg.Log.FilePath,
	)
}

// App 应用结构
type App struct {
	Config *config.Config
	Logger *logger.ZapLogger
	DB     *mysql.DB
	Redis  *redis.Client
	Server *http.Server
}

// NewApp 创建应用实例
func NewApp(
	cfg *config.Config,
	log *logger.ZapLogger,
	db *mysql.DB,
	rdb *redis.Client,
	srv *http.Server,
) *App {
	return &App{
		Config: cfg,
		Logger: log,
		DB:     db,
		Redis:  rdb,
		Server: srv,
	}
}
