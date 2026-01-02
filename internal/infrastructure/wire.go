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
	"member-pre/internal/infrastructure/repository"
)

// 类型别名，用于 Wire 依赖注入时区分不同的 string 和 int 类型
type ConfigPath string
type JWTSecret string
type TokenExpires int

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
		NewLogger,
		ProvideLogger,           // 将ZapLogger转换为pkg/logger.Logger接口
		ProvideJWTSecret,        // 提供JWT密钥（类型别名）
		ProvideAuthJWTSecret,    // 将JWTSecret转换为string（供domain层使用）
		ProvideTokenExpires,     // 提供Token过期时间（类型别名）
		ProvideAuthTokenExpires, // 将TokenExpires转换为int（供domain层使用）
		database.NewDatabase,
		persistence.NewClient,
		repository.WireRepoSet, // Repository需要logger
		domain.WireDoMainSet,   // Domain需要logger和配置值
		http.WireHttpSet,       // Handler需要logger和RouteRegistrar
		// 提供RouteRegistrar列表
		ProvideRouteRegistrars,
		// HTTP服务器（需要logger用于中间件）
		httpInfra.NewServer,
		// 应用
		NewApp,
	)
	return &App{}, nil
}

// NewLogger 创建日志实例
func NewLogger(cfg *config.Config) (*infraLogger.ZapLogger, error) {
	return infraLogger.NewZapLogger(
		cfg.Log.Level,
		cfg.Log.Format,
		cfg.Log.Output,
		cfg.Log.FilePath,
	)
}

// ProvideLogger 将ZapLogger转换为pkg/logger.Logger接口
// 由于ZapLogger已经实现了pkg/logger.Logger接口，可以直接返回
func ProvideLogger(zapLogger *infraLogger.ZapLogger) pkgLogger.Logger {
	return zapLogger
}

// ProvideJWTSecret 提供JWT密钥（类型别名，避免与cfgPath string冲突）
func ProvideJWTSecret(cfg *config.Config) JWTSecret {
	return JWTSecret(cfg.Auth.JWTSecret)
}

// ProvideTokenExpires 提供Token过期时间（类型别名）
func ProvideTokenExpires(cfg *config.Config) TokenExpires {
	return TokenExpires(cfg.Auth.TokenExpires)
}

// ProvideAuthJWTSecret 将JWTSecret转换为string，供domain层使用
// 使用明确的函数名避免与cfgPath参数冲突
func ProvideAuthJWTSecret(secret JWTSecret) string {
	return string(secret)
}

// ProvideAuthTokenExpires 将TokenExpires转换为int，供domain层使用
func ProvideAuthTokenExpires(expires TokenExpires) int {
	return int(expires)
}

// ProvideRouteRegistrars 提供路由注册器列表
func ProvideRouteRegistrars(authRegistrar httpInfra.RouteRegistrar) []httpInfra.RouteRegistrar {
	return []httpInfra.RouteRegistrar{authRegistrar}
}

// App 应用结构
type App struct {
	Config *config.Config
	Logger *infraLogger.ZapLogger
	DB     database.Database
	Redis  *persistence.Client
	Server *httpInfra.Server
}

// NewApp 创建应用实例
func NewApp(
	cfg *config.Config,
	log *infraLogger.ZapLogger,
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
