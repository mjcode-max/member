package http

import (
	"github.com/gin-gonic/gin"
	"member-pre/internal/infrastructure/config"
	"member-pre/internal/infrastructure/persistence/mysql"
	"member-pre/internal/infrastructure/persistence/redis"
	"member-pre/internal/interfaces/http/router"
)

// SetupRoutes 设置路由
func SetupRoutes(engine *gin.Engine, cfg *config.Config, db *mysql.DB, rdb *redis.Client) {
	// 初始化各个模块的路由
	routers := []Router{
		router.InitAuthRouter(cfg, db, rdb),
		// 添加新模块时，只需在这里添加一行：
		// router.Init{Module}Router(cfg, db, rdb),
	}

	// 使用统一的注册函数
	RegisterRoutes(engine, routers...)
}
