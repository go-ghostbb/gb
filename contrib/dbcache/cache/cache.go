package cache

import (
	"context"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbcache "ghostbb.io/gb/os/gb_cache"
	gbrand "ghostbb.io/gb/util/gb_rand"
	"time"
)

var cacheName = "gb:gorm:cache"

var (
	ErrCacheNotFound = gberror.New("cache not found")
)

type Kv struct {
	Key   string
	Value string
}

type Cache struct {
	cache      *gbcache.Cache
	InstanceId string
	config     *Config
}

type Config struct {
	InvalidateWhenUpdate bool
	AsyncWrite           bool
	TTL                  time.Duration
	MaxItem              int
}

func New() *Cache {
	return &Cache{
		cache:      gbcache.New(),
		InstanceId: gbrand.S(5),
		config: &Config{
			InvalidateWhenUpdate: true,
			AsyncWrite:           false,
			TTL:                  5 * time.Minute,
			MaxItem:              100,
		},
	}
}

func (c *Cache) Config() *Config {
	return c.config
}

func (c *Cache) SetAdapter(adapter gbcache.Adapter) {
	c.cache.SetAdapter(adapter)
}

func (c *Cache) InvalidateSearchCache(ctx context.Context, tableName string) error {
	return c.DeleteKeysWithPrefix(ctx, c.genCachePrefix(tableName))
}

func (c *Cache) SetSearchCache(ctx context.Context, cacheValue string, tableName string, sql string, vars ...interface{}) error {
	return c.SetKey(ctx, Kv{
		Key:   c.genSearchCacheKey(tableName, sql, vars...),
		Value: cacheValue,
	})
}

func (c *Cache) GetSearchCache(ctx context.Context, tableName string, sql string, vars ...interface{}) (string, error) {
	return c.GetValue(ctx, c.genSearchCacheKey(tableName, sql, vars...))
}
