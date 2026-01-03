package http

import (
	"github.com/gin-gonic/gin"
	"member-pre/internal/domain/auth"
	httpInfra "member-pre/internal/infrastructure/http"
	"member-pre/internal/interfaces/http/handler"
	"member-pre/pkg/logger"
)

// appRouteRegistrar 应用路由注册器
// 所有HTTP路由必须在此注册器中定义
type appRouteRegistrar struct {
	authHandler *handler.AuthHandler
	userHandler *handler.UserHandler
	authService *auth.AuthService
	logger      logger.Logger
}

// RegisterRoutes 注册所有应用路由
//
// 重要提示：所有HTTP路由必须且只能在此函数中定义！
//
// 路由注册规则：
// 1. 所有路由必须在此函数内通过 api.Group() 创建路由组
// 2. 需要认证的路由必须使用 AuthMiddleware 中间件
// 3. 路由路径使用 RESTful 风格，例如：/users/:id
// 4. 每个路由组添加清晰的注释说明用途
// 5. 不要在其他地方定义路由，包括 handler 文件、setup.go 等
//
// 示例：
//
//	authGroup := api.Group("/auth")
//	{
//	    // 公开接口
//	    authGroup.POST("/login", r.authHandler.Login)
//
//	    // 需要认证的接口
//	    authProtected := authGroup.Group("")
//	    authProtected.Use(httpInfra.AuthMiddleware(r.authService, r.logger))
//	    {
//	        authProtected.GET("/me", r.authHandler.GetCurrentUser)
//	    }
//	}
func (r *appRouteRegistrar) RegisterRoutes(api *gin.RouterGroup) {
	// ==================== 认证相关路由 ====================
	authGroup := api.Group("/auth")
	{
		// 登录接口（无需认证）
		authGroup.POST("/login", r.authHandler.Login)

		// 需要认证的路由组
		authProtected := authGroup.Group("")
		authProtected.Use(httpInfra.AuthMiddleware(r.authService, r.logger))
		{
			// 获取当前用户（需要认证）
			authProtected.GET("/me", r.authHandler.GetCurrentUser)
			// 登出接口（需要认证）
			authProtected.POST("/logout", r.authHandler.Logout)
		}
	}

	// ==================== 用户管理相关路由 ====================
	usersGroup := api.Group("/users")
	{
		// 需要认证的路由组
		usersProtected := usersGroup.Group("")
		usersProtected.Use(httpInfra.AuthMiddleware(r.authService, r.logger))
		{
			// 获取用户详情
			usersProtected.GET("/:id", r.userHandler.GetUser)
			// 更新美甲师工作状态
			usersProtected.PUT("/:id/work-status", r.userHandler.UpdateWorkStatus)
		}
	}

	// ==================== 门店相关路由 ====================
	storesGroup := api.Group("/stores")
	{
		// 需要认证的路由组
		storesProtected := storesGroup.Group("")
		storesProtected.Use(httpInfra.AuthMiddleware(r.authService, r.logger))
		{
			// 根据门店获取用户列表
			storesProtected.GET("/:store_id/users", r.userHandler.GetUsersByStore)
		}
	}
}

// NewAppRouteRegistrar 创建应用路由注册器
func NewAppRouteRegistrar(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	authService *auth.AuthService,
	log logger.Logger,
) httpInfra.RouteRegistrar {
	return &appRouteRegistrar{
		authHandler: authHandler,
		userHandler: userHandler,
		authService: authService,
		logger:      log,
	}
}
