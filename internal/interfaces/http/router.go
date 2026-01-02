package http

import (
	"github.com/gin-gonic/gin"
	httpInfra "member-pre/internal/infrastructure/http"
	"member-pre/internal/interfaces/http/handler"
)

// authRouteRegistrar 认证路由注册器
type authRouteRegistrar struct {
	handler *handler.AuthHandler
}

// RegisterRoutes 注册路由（实现 RouteRegistrar 接口）
func (r *authRouteRegistrar) RegisterRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		// 登录接口（无需认证）
		auth.POST("/login", r.handler.Login)
		// 获取当前用户（需要认证）
		auth.GET("/me", r.handler.GetCurrentUser)
		// 登出接口（需要认证）
		auth.POST("/logout", r.handler.Logout)
	}
}

// NewAuthRouteRegistrar 创建认证路由注册器
func NewAuthRouteRegistrar(handler *handler.AuthHandler) httpInfra.RouteRegistrar {
	return &authRouteRegistrar{
		handler: handler,
	}
}
