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
	MaxItem              int64
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

func (c *Cache) InvalidateSearchCache(ctx context.Context, tableName string) error {
	return c.DeleteKeysWithPrefix(ctx, c.genCachePrefix(tableName))
}

func (c *Cache) InvalidatePrimaryCache(ctx context.Context, tableName string, primaryKey string) error {
	return c.DeleteKey(ctx, c.genPrimaryCacheKey(tableName, primaryKey))
}

func (c *Cache) BatchInvalidatePrimaryCache(ctx context.Context, tableName string, primaryKeys []string) error {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, c.genPrimaryCacheKey(tableName, primaryKey))
	}
	return c.BatchDeleteKeys(ctx, cacheKeys)
}

func (c *Cache) InvalidateAllPrimaryCache(ctx context.Context, tableName string) error {
	return c.DeleteKeysWithPrefix(ctx, c.genCachePrefix(tableName))
}

func (c *Cache) BatchPrimaryKeyExists(ctx context.Context, tableName string, primaryKeys []string) (bool, error) {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, c.genPrimaryCacheKey(tableName, primaryKey))
	}
	return c.BatchKeyExist(ctx, cacheKeys)
}

func (c *Cache) SearchKeyExists(ctx context.Context, tableName string, sql string, vars ...interface{}) (bool, error) {
	return c.KeyExists(ctx, c.genSearchCacheKey(tableName, sql, vars...))
}

func (c *Cache) BatchSetPrimaryKeyCache(ctx context.Context, tableName string, kvs []Kv) error {
	for idx, kv := range kvs {
		kvs[idx].Key = c.genPrimaryCacheKey(tableName, kv.Key)
	}
	return c.BatchSetKeys(ctx, kvs)
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

func (c *Cache) BatchGetPrimaryCache(ctx context.Context, tableName string, primaryKeys []string) ([]string, error) {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, c.genPrimaryCacheKey(tableName, primaryKey))
	}
	return c.BatchGetValues(ctx, cacheKeys)
}
