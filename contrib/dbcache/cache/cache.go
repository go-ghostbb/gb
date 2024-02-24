package cache

import (
	"context"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbrand "ghostbb.io/gb/util/gb_rand"
	"time"
)

var cacheName = "gb:gorm:cache"

var (
	ErrCacheNotFound = gberror.New("cache not found")
)

type (
	Cache struct {
		cache  ICache
		config *Config
	}

	Config struct {
		instance             string
		InvalidateWhenUpdate bool
		AsyncWrite           bool
		TTL                  time.Duration
		MaxItem              int
	}

	ICache interface {
		setConfig(config *Config)
		Clear(ctx context.Context) error
		Get(ctx context.Context, key string) (string, error)
		Set(ctx context.Context, key string, value string) error
		DeleteKeysWithPrefix(ctx context.Context, keyPrefix string) error
	}
)

func New(cache ICache) *Cache {
	config := &Config{
		instance:             gbrand.S(5),
		InvalidateWhenUpdate: true,
		AsyncWrite:           false,
		TTL:                  5 * time.Minute,
		MaxItem:              100,
	}
	cache.setConfig(config)
	return &Cache{
		cache:  cache,
		config: config,
	}
}

func (c *Cache) Config() *Config {
	return c.config
}

func (c *Cache) SetName(name string) {
	c.config.instance = name
}

func (c *Cache) Clear(ctx context.Context) error {
	return c.cache.Clear(ctx)
}

func (c *Cache) InvalidateSearchCache(ctx context.Context, tableName string) error {
	return c.cache.DeleteKeysWithPrefix(ctx, c.genCachePrefix(tableName))
}

func (c *Cache) SetSearchCache(ctx context.Context, cacheValue string, tableName string, sql string, vars ...interface{}) error {
	return c.cache.Set(ctx, c.genSearchCacheKey(tableName, sql, vars...), cacheValue)
}

func (c *Cache) GetSearchCache(ctx context.Context, tableName string, sql string, vars ...interface{}) (string, error) {
	return c.cache.Get(ctx, c.genSearchCacheKey(tableName, sql, vars...))
}
