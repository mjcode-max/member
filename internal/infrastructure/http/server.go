package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"member-pre/internal/infrastructure/config"
	"member-pre/internal/infrastructure/persistence/mysql"
	"member-pre/internal/infrastructure/persistence/redis"
	httpInterface "member-pre/internal/interfaces/http"
)

// Server HTTP服务器
type Server struct {
	engine *gin.Engine
	server *http.Server
	config *config.ServerConfig
}

// NewServer 创建HTTP服务器
func NewServer(appCfg *config.Config, db *mysql.DB, rdb *redis.Client) *Server {
	cfg := &appCfg.Server
	// 设置Gin模式
	gin.SetMode(cfg.Mode)

	engine := gin.New()

	// 添加中间件
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// 设置路由
	httpInterface.SetupRoutes(engine, appCfg, db, rdb)

	return &Server{
		engine: engine,
		config: cfg,
	}
}

// GetEngine 获取Gin引擎
func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

// Start 启动服务器
func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		Handler:      s.engine,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("启动HTTP服务器失败: %w", err)
	}

	return nil
}

// Stop 停止服务器
func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}
