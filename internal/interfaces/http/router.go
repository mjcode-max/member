package http

import (
	"github.com/gin-gonic/gin"
	"member-pre/internal/domain/auth"
	"member-pre/internal/domain/user"
	httpInfra "member-pre/internal/infrastructure/http"
	"member-pre/internal/interfaces/http/handler"
	"member-pre/pkg/logger"
)

// appRouteRegistrar 应用路由注册器
// 所有HTTP路由必须在此注册器中定义
type appRouteRegistrar struct {
	authHandler  *handler.AuthHandler
	userHandler  *handler.UserHandler
	storeHandler *handler.StoreHandler
	authService  *auth.AuthService
	logger       logger.Logger
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
			// 获取用户列表（后台、店长）
			// 后台可以查看所有用户，店长只能查看自己门店的美甲师
			usersProtected.GET("", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager), r.userHandler.GetUserList)

			// 创建用户（后台、店长）
			// 后台可以创建所有角色，店长只能创建美甲师
			usersProtected.POST("", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager), r.userHandler.CreateUser)

			// 获取用户详情（所有已认证用户）
			// 后台可以查看所有用户，店长可以查看自己门店的用户，美甲师和顾客只能查看自己的信息
			usersProtected.GET("/:id", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician, user.RoleCustomer), r.userHandler.GetUser)

			// 更新用户（后台、店长、美甲师、顾客）
			// 后台可以更新所有用户，店长可以更新自己门店的用户，美甲师和顾客只能更新自己的信息
			usersProtected.PUT("/:id", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician, user.RoleCustomer), r.userHandler.UpdateUser)

			// 更新用户状态（仅后台）
			usersProtected.PUT("/:id/status", httpInfra.RoleMiddleware(user.RoleAdmin), r.userHandler.UpdateUserStatus)

			// 更新美甲师工作状态（后台、店长、美甲师）
			// 后台可以更新任何美甲师，店长可以更新自己门店的美甲师，美甲师只能更新自己的状态
			usersProtected.PUT("/:id/work-status", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician), r.userHandler.UpdateWorkStatus)
		}
	}

	// ==================== 公开接口（无需认证） ====================
	publicGroup := api.Group("/public")
	{
		// 获取公开门店列表（供客户使用，默认只返回营业中的门店）
		publicGroup.GET("/stores", r.storeHandler.GetPublicStoreList)
	}

	// ==================== 门店管理相关路由 ====================
	storesGroup := api.Group("/stores")
	{

		// 需要认证的路由组
		storesProtected := storesGroup.Group("")
		storesProtected.Use(httpInfra.AuthMiddleware(r.authService, r.logger))
		{
			// 获取门店列表（所有已认证用户）
			// 后台可以查看所有门店，店长、美甲师、顾客可以查看门店列表
			storesProtected.GET("", httpInfra.RoleMiddleware(user.RoleAdmin), r.storeHandler.GetStoreList)

			// 获取门店详情（所有已认证用户）
			// 后台可以查看所有门店，店长只能查看自己的门店，美甲师和顾客可以查看门店详情
			storesProtected.GET("/:id", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician, user.RoleCustomer), r.storeHandler.GetStore)

			// 创建门店（仅后台）
			storesProtected.POST("", httpInfra.RoleMiddleware(user.RoleAdmin), r.storeHandler.CreateStore)

			// 更新门店（后台、店长）
			// 后台可以更新所有门店，店长只能更新自己的门店
			storesProtected.PUT("/:id", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager), r.storeHandler.UpdateStore)

			// 删除门店（仅后台）
			storesProtected.DELETE("/:id", httpInfra.RoleMiddleware(user.RoleAdmin), r.storeHandler.DeleteStore)

			// 更新门店状态（仅后台）
			storesProtected.PUT("/:id/status", httpInfra.RoleMiddleware(user.RoleAdmin), r.storeHandler.UpdateStoreStatus)
		}
	}
}

// NewAppRouteRegistrar 创建应用路由注册器
func NewAppRouteRegistrar(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	storeHandler *handler.StoreHandler,
	authService *auth.AuthService,
	log logger.Logger,
) httpInfra.RouteRegistrar {
	return &appRouteRegistrar{
		authHandler:  authHandler,
		userHandler:  userHandler,
		storeHandler: storeHandler,
		authService:  authService,
		logger:       log,
	}
}
