package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"member-pre/internal/infrastructure/config"
)

// Client Redis客户端
type Client struct {
	*redis.Client
}

// NewClient 创建Redis客户端
func NewClient(cfg *config.Config) (*Client, error) {
	redisCfg := &cfg.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password:     redisCfg.Password,
		DB:           redisCfg.DB,
		PoolSize:     redisCfg.PoolSize,
		MinIdleConns: redisCfg.MinIdleConns,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis连接失败: %w", err)
	}

	return &Client{Client: rdb}, nil
}

// Close 关闭Redis连接
func (c *Client) Close() error {
	return c.Client.Close()
}
