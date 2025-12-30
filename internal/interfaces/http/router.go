package http

import (
	"github.com/gin-gonic/gin"
)

// Router 路由接口
type Router interface {
	RegisterRoutes(engine *gin.Engine)
}

// RegisterRoutes 注册所有路由
func RegisterRoutes(engine *gin.Engine, routers ...Router) {
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

		// 注册各个模块的路由
		for _, router := range routers {
			router.RegisterRoutes(engine)
		}
	}
}
