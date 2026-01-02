package http

import (
	"github.com/gin-gonic/gin"
)

// RouteRegistrar 路由注册器接口
// 每个模块的 handler 文件需要实现此接口来注册路由
type RouteRegistrar interface {
	RegisterRoutes(api *gin.RouterGroup)
}

// SetupRoutes 设置路由
// 接收路由注册器列表，通过 Wire 注入
func SetupRoutes(engine *gin.Engine, registrars []RouteRegistrar) {
	// API路由组
	api := engine.Group("/api/v1")
	{
		// 健康检查
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "服务运行正常",
			})
		})

		// 注册所有路由注册器
		for _, registrar := range registrars {
			registrar.RegisterRoutes(api)
		}
	}
}
