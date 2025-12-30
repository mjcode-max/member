package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"member-pre/internal/domain/auth"
	"member-pre/internal/infrastructure/config"
	"member-pre/internal/infrastructure/persistence/mysql"
	"member-pre/internal/infrastructure/persistence/redis"
	repository "member-pre/internal/infrastructure/persistence/repository"
	"member-pre/internal/interfaces/http/handler"
	"member-pre/internal/interfaces/http/middleware"
)

// AuthRouter 认证路由
type AuthRouter struct {
	handler    *handler.AuthHandler
	middleware *middleware.AuthMiddleware
}

// NewAuthRouter 创建认证路由
func NewAuthRouter(authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) *AuthRouter {
	return &AuthRouter{
		handler:    authHandler,
		middleware: authMiddleware,
	}
}

// InitAuthRouter 初始化认证路由（包含依赖创建）
func InitAuthRouter(cfg *config.Config, db *mysql.DB, rdb *redis.Client) *AuthRouter {
	// 创建仓储
	userRepo := repository.NewUserRepository(db, rdb)

	// 创建认证服务
	authService := auth.NewService(
		userRepo,
		cfg.Auth.JWTSecret,
		time.Duration(cfg.Auth.TokenExpires)*time.Second,
	)

	// 创建处理器和中间件
	authHandler := handler.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// 创建路由实例
	return NewAuthRouter(authHandler, authMiddleware)
}

// RegisterRoutes 注册路由
// engine 参数应该是 /api/v1 路由组
func (r *AuthRouter) RegisterRoutes(engine *gin.Engine) {
	auth := engine.Group("/auth")
	{
		// 登录（不需要认证）
		auth.POST("/login", r.handler.Login)

		// 登出（需要认证）
		auth.POST("/logout", r.middleware.RequireAuth(), r.handler.Logout)

		// 获取当前用户信息（需要认证）
		auth.GET("/me", r.middleware.RequireAuth(), r.handler.GetCurrentUser)
	}
}
