package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"member-pre/internal/infrastructure/config"
	"member-pre/pkg/cache"
)

// Client Redis客户端
// 实现了pkg/cache.Cache接口
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

	// 测试连接（使用较短的超时时间）
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		// 关闭连接以避免资源泄漏
		rdb.Close()
		return nil, fmt.Errorf("Redis连接失败: %w", err)
	}

	return &Client{Client: rdb}, nil
}

// NewClientOptional 创建Redis客户端（可选，连接失败时返回nil和nil error）
// 用于迁移等不需要Redis的场景
func NewClientOptional(cfg *config.Config) (*Client, error) {
	client, err := NewClient(cfg)
	if err != nil {
		// 连接失败时返回nil和nil error，允许继续执行
		// Wire 可以处理这种情况
		return nil, nil
	}
	return client, nil
}

// 实现pkg/cache.Cache接口

// Get 获取缓存值
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	result := c.Client.Get(ctx, key)
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return "", cache.ErrNotFound
		}
		return "", result.Err()
	}
	return result.Val(), nil
}

// Set 设置缓存值
func (c *Client) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

// Delete 删除缓存
func (c *Client) Delete(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}

// Exists 检查key是否存在
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	result := c.Client.Exists(ctx, key)
	if result.Err() != nil {
		return false, result.Err()
	}
	return result.Val() > 0, nil
}

// SetNX 当key不存在时设置值
func (c *Client) SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	result := c.Client.SetNX(ctx, key, value, expiration)
	if result.Err() != nil {
		return false, result.Err()
	}
	return result.Val(), nil
}

// Close 关闭Redis连接
func (c *Client) Close() error {
	return c.Client.Close()
}
