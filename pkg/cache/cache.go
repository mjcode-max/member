package cache

import (
	"context"
	"errors"
	"time"
)

// Cache 缓存接口，定义在pkg层，供domain层使用
// domain层只依赖此接口，不依赖具体实现
type Cache interface {
	// Get 获取缓存值
	Get(ctx context.Context, key string) (string, error)
	// Set 设置缓存值
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	// Delete 删除缓存
	Delete(ctx context.Context, key string) error
	// Exists 检查key是否存在
	Exists(ctx context.Context, key string) (bool, error)
	// SetNX 当key不存在时设置值
	SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
	// Close 关闭缓存连接
	Close() error
}

// ErrNotFound 缓存未找到错误
var ErrNotFound = errors.New("cache: key not found")
