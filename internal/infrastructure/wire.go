//go:build wireinject
// +build wireinject

package infrastructure

import (
	"github.com/google/wire"
	"member-pre/internal/infrastructure/config"
	httpInfra "member-pre/internal/infrastructure/http"
	infraLogger "member-pre/internal/infrastructure/logger"
	"member-pre/internal/infrastructure/persistence"
	"member-pre/internal/infrastructure/persistence/database"
	"member-pre/internal/interfaces/http"
	pkgLogger "member-pre/pkg/logger"

	"member-pre/internal/domain"
	"member-pre/internal/domain/auth"
	"member-pre/internal/infrastructure/repository"
)

// 类型别名，用于 Wire 依赖注入时区分不同的 string 类型
type ConfigPath string

// LoadConfig 加载配置（包装config.Load，避免string类型冲突）
func LoadConfig(cfgPath ConfigPath) (*config.Config, error) {
	return config.Load(string(cfgPath))
}

// InitializeApp 初始化应用
// 所有模块的 Provider Set 都在 handler 包中定义，这里只需要引用
func InitializeApp(cfgPath ConfigPath) (*App, error) {
	wire.Build(
		// 基础设施层
		LoadConfig, // 加载配置（接受ConfigPath类型）
		// 绑定配置到接口
		wire.Bind(new(auth.IAuthConfig), new(*config.Config)),
		NewLogger, // 创建日志实例，直接返回logger.Logger接口
		database.NewDatabase,
		persistence.NewClient,
		repository.WireRepoSet, // Repository需要logger
		domain.WireDoMainSet,   // Domain需要logger和配置值
		http.WireHttpSet,       // Handler需要logger和RouteRegistrar
		// 提供RouteRegistrar切片（直接收集所有RouteRegistrar）
		ProvideRouteRegistrars,
		// HTTP服务器（需要logger和RouteRegistrar切片）
		httpInfra.NewServer,
		// 应用
		NewApp,
	)
	return &App{}, nil
}

// NewLogger 创建日志实例，直接返回logger.Logger接口
func NewLogger(cfg *config.Config) (pkgLogger.Logger, error) {
	zapLogger, err := infraLogger.NewZapLogger(
		cfg.Log.Level,
		cfg.Log.Format,
		cfg.Log.Output,
		cfg.Log.FilePath,
	)
	if err != nil {
		return nil, err
	}
	// ZapLogger已经实现了pkgLogger.Logger接口，可以直接返回
	return zapLogger, nil
}

// ProvideRouteRegistrars 提供路由注册器切片
// 直接收集所有RouteRegistrar，这是Wire要求的实现方式
func ProvideRouteRegistrars(appRegistrar httpInfra.RouteRegistrar) []httpInfra.RouteRegistrar {
	return []httpInfra.RouteRegistrar{appRegistrar}
}

// App 应用结构
type App struct {
	Config *config.Config
	Logger pkgLogger.Logger
	DB     database.Database
	Redis  *persistence.Client
	Server *httpInfra.Server
}

// NewApp 创建应用实例
func NewApp(
	cfg *config.Config,
	log pkgLogger.Logger,
	db database.Database,
	rdb *persistence.Client,
	srv *httpInfra.Server,
) *App {
	return &App{
		Config: cfg,
		Logger: log,
		DB:     db,
		Redis:  rdb,
		Server: srv,
	}
}
