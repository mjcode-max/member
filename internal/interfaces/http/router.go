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
	authHandler       *handler.AuthHandler
	userHandler       *handler.UserHandler
	storeHandler      *handler.StoreHandler
	templateHandler   *handler.SlotTemplateHandler
	slotHandler       *handler.SlotHandler
	memberHandler     *handler.MemberHandler
	appointmentHandler *handler.AppointmentHandler
	authService       *auth.AuthService
	logger            logger.Logger
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

		// 通过微信code换取openid并保存（公开接口，无需认证）
		publicGroup.POST("/customer/wechat/login", r.userHandler.WechatLoginByCode)
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

	// ==================== 时段模板管理相关路由 ====================
	templateGroup := api.Group("/slot-templates")
	{
		// 需要认证的路由组
		templateProtected := templateGroup.Group("")
		templateProtected.Use(httpInfra.AuthMiddleware(r.authService, r.logger))
		{
			// 获取时段模板列表（后台、店长）
			templateProtected.GET("", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager), r.templateHandler.GetTemplateList)

			// 获取时段模板详情（后台、店长）
			templateProtected.GET("/:id", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager), r.templateHandler.GetTemplate)

			// 创建时段模板（后台、店长）
			templateProtected.POST("", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager), r.templateHandler.CreateTemplate)

			// 更新时段模板（后台、店长）
			templateProtected.PUT("/:id", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager), r.templateHandler.UpdateTemplate)

			// 删除时段模板（后台、店长）
			templateProtected.DELETE("/:id", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager), r.templateHandler.DeleteTemplate)
		}
	}

	// ==================== 时段管理相关路由 ====================
	slotsGroup := api.Group("/slots")
	{
		// 获取可用时段列表（无需认证，顾客端查询可预约时间）
		slotsGroup.GET("/available", r.slotHandler.GetAvailableSlotList)

		// 需要认证的路由组
		slotsProtected := slotsGroup.Group("")
		slotsProtected.Use(httpInfra.AuthMiddleware(r.authService, r.logger))
		{
			// 获取时段列表（所有已认证用户）
			slotsProtected.GET("", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician, user.RoleCustomer), r.slotHandler.GetSlotList)

			// 获取时段详情（所有已认证用户）
			slotsProtected.GET("/:id", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician, user.RoleCustomer), r.slotHandler.GetSlot)

			// 生成时段（后台、店长）
			slotsProtected.POST("/generate", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager), r.slotHandler.GenerateSlots)

			// 锁定时段（所有已认证用户，预约时调用）
			slotsProtected.POST("/lock", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician, user.RoleCustomer), r.slotHandler.LockSlot)

			// 解锁时段（所有已认证用户，取消预约时调用）
			slotsProtected.POST("/unlock", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician, user.RoleCustomer), r.slotHandler.UnlockSlot)

			// 预约时段（所有已认证用户）
			slotsProtected.POST("/book", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician, user.RoleCustomer), r.slotHandler.BookSlot)

			// 释放时段（所有已认证用户，取消或完成时调用）
			slotsProtected.POST("/release", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician, user.RoleCustomer), r.slotHandler.ReleaseSlot)

			// 重新计算时段容量（后台、店长）
			slotsProtected.POST("/recalculate-capacity", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager), r.slotHandler.RecalculateCapacity)
		}
	}

	// ==================== 会员管理相关路由 ====================
	membersGroup := api.Group("/members")
	{
		// 需要认证的路由组
		membersProtected := membersGroup.Group("")
		membersProtected.Use(httpInfra.AuthMiddleware(r.authService, r.logger))
		{
			// 创建会员（管理员、店长）
			membersProtected.POST("", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager), r.memberHandler.CreateMember)

			// 获取会员列表（所有已认证用户）
			// 后台可以查看所有会员，店长只能查看自己门店的会员
			membersProtected.GET("", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician), r.memberHandler.GetMemberList)

			// 获取会员详情（所有已认证用户）
			// 后台可以查看所有会员，店长只能查看自己门店的会员
			membersProtected.GET("/:id", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician, user.RoleCustomer), r.memberHandler.GetMember)

			// 根据手机号查询会员列表（所有已认证用户）
			membersProtected.GET("/phone/:phone", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician), r.memberHandler.GetMemberByPhone)

			// 更新会员信息（管理员、店长）
			// 后台可以更新所有会员，店长只能更新自己门店的会员
			membersProtected.PUT("/:id", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager), r.memberHandler.UpdateMember)

			// 更新会员状态（管理员、店长）
			membersProtected.PUT("/:id/status", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager), r.memberHandler.UpdateMemberStatus)

			// 创建使用记录（管理员、店长、美甲师）
			membersProtected.POST("/:id/usages", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician), r.memberHandler.CreateUsage)

			// 获取会员使用记录列表（不分页，管理员、店长、美甲师）
			membersProtected.GET("/:id/usages", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician), r.memberHandler.GetUsageByMemberID)
		}
	}

	// ==================== 使用记录相关路由 ====================
	usagesGroup := api.Group("/usages")
	{
		// 需要认证的路由组
		usagesProtected := usagesGroup.Group("")
		usagesProtected.Use(httpInfra.AuthMiddleware(r.authService, r.logger))
		{
			// 获取使用记录列表（不分页，管理员、店长、美甲师）
			usagesProtected.GET("", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician), r.memberHandler.GetUsageList)

			// 获取使用记录详情（管理员、店长、美甲师）
			usagesProtected.GET("/:id", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician), r.memberHandler.GetUsage)

			// 删除使用记录（仅管理员）
			usagesProtected.DELETE("/:id", httpInfra.RoleMiddleware(user.RoleAdmin), r.memberHandler.DeleteUsage)
		}
	}

	// ==================== 预约管理相关路由 ====================
	appointmentsGroup := api.Group("/appointments")
	{
		// 需要认证的路由组
		appointmentsProtected := appointmentsGroup.Group("")
		appointmentsProtected.Use(httpInfra.AuthMiddleware(r.authService, r.logger))
		{
			// 获取我的预约列表
			appointmentsProtected.GET("/my", r.appointmentHandler.GetMyAppointments)

			// 获取我即将到来的预约
			appointmentsProtected.GET("/my/upcoming", r.appointmentHandler.GetMyUpcomingAppointments)

			// 获取预约详情
			appointmentsProtected.GET("/:id", r.appointmentHandler.GetAppointment)

			// 创建预约
			appointmentsProtected.POST("", r.appointmentHandler.CreateAppointment)

			// 支付押金
			appointmentsProtected.POST("/pay-deposit", r.appointmentHandler.PayDeposit)

			// 美甲师确认到店
			appointmentsProtected.POST("/confirm-arrival", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician), r.appointmentHandler.ConfirmArrival)

			// 完成预约
			appointmentsProtected.POST("/complete", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician), r.appointmentHandler.Complete)

			// 顾客取消预约（需提前3小时）
			appointmentsProtected.POST("/cancel/customer", r.appointmentHandler.CancelByCustomer)

			// 美甲师取消预约（随时可以取消）
			appointmentsProtected.POST("/cancel/technician", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician), r.appointmentHandler.CancelByTechnician)

			// 获取门店预约列表（后台、店长、美甲师）
			appointmentsProtected.GET("/store", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician), r.appointmentHandler.GetStoreAppointments)

			// 获取美甲师预约列表
			appointmentsProtected.GET("/technician", httpInfra.RoleMiddleware(user.RoleAdmin, user.RoleStoreManager, user.RoleTechnician), r.appointmentHandler.GetTechnicianAppointments)
		}
	}
}

// NewAppRouteRegistrar 创建应用路由注册器
func NewAppRouteRegistrar(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	storeHandler *handler.StoreHandler,
	templateHandler *handler.SlotTemplateHandler,
	slotHandler *handler.SlotHandler,
	memberHandler *handler.MemberHandler,
	appointmentHandler *handler.AppointmentHandler,
	authService *auth.AuthService,
	log logger.Logger,
) httpInfra.RouteRegistrar {
	return &appRouteRegistrar{
		authHandler:        authHandler,
		userHandler:        userHandler,
		storeHandler:       storeHandler,
		templateHandler:    templateHandler,
		slotHandler:        slotHandler,
		memberHandler:      memberHandler,
		appointmentHandler: appointmentHandler,
		authService:        authService,
		logger:             log,
	}
}
