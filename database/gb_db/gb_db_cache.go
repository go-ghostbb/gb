package gbdb

import (
	"context"
	"ghostbb.io/gb/internal/intlog"
	gbcache "ghostbb.io/gb/os/gb_cache"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbrand "ghostbb.io/gb/util/gb_rand"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	DBCacheName = "gb.database.cache"
)

func NewCache(config *CacheConfig) (ICache, error) {
	cache := &Cache{
		Config: config,
		stats:  &stats{},
	}
	if err := cache.Init(); err != nil {
		return nil, err
	}
	return cache, nil
}

type ICache interface {
	Name() string
	Initialize(db *gorm.DB) error
	AttachToDB(db *gorm.DB)

	ResetCache() error
	StatsAccessor
}

type Cache struct {
	Config     *CacheConfig
	InstanceId string

	cache  IDataCache
	logger logger.Interface

	*stats
}

func (c *Cache) Name() string {
	return DBCacheName
}

func (c *Cache) Initialize(db *gorm.DB) (err error) {
	c.logger = db.Logger

	err = db.Callback().Create().After("gorm:create").Register("gorm:cache:after_create", AfterCreate(c))
	if err != nil {
		return err
	}

	err = db.Callback().Delete().After("gorm:delete").Register("gorm:cache:after_delete", AfterDelete(c))
	if err != nil {
		return err
	}

	err = db.Callback().Update().After("gorm:update").Register("gorm:cache:after_update", AfterUpdate(c))
	if err != nil {
		return err
	}

	err = newQueryHandler(c).Bind(db)
	if err != nil {
		return err
	}

	return
}

func (c *Cache) AttachToDB(db *gorm.DB) {
	_ = c.Initialize(db)
}

func (c *Cache) Init() error {
	c.InstanceId = gbrand.S(5)

	if c.Config.Cache == nil {
		c.Config.Cache = gbcache.New()
	}
	c.cache = newDataCache(c.Config)

	err := c.cache.Init()
	if err != nil {
		intlog.Errorf(gbctx.New(), "[Init] cache init error: %v", err)
		return err
	}
	return nil
}

func (c *Cache) ResetCache() error {
	c.stats.ResetHitCount()
	ctx := gbctx.New()
	err := c.cache.ClearCache(ctx)
	if err != nil {
		c.logger.Error(context.TODO(), "[ResetCache] reset cache error: %v", err)
		return err
	}
	return nil
}

func (c *Cache) InvalidateSearchCache(ctx context.Context, tableName string) error {
	return c.cache.DeleteKeysWithPrefix(ctx, genSearchCachePrefix(c.InstanceId, tableName))
}

func (c *Cache) InvalidatePrimaryCache(ctx context.Context, tableName string, primaryKey string) error {
	return c.cache.DeleteKey(ctx, genPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
}

func (c *Cache) BatchInvalidatePrimaryCache(ctx context.Context, tableName string, primaryKeys []string) error {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, genPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
	}
	return c.cache.BatchDeleteKeys(ctx, cacheKeys)
}

func (c *Cache) InvalidateAllPrimaryCache(ctx context.Context, tableName string) error {
	return c.cache.DeleteKeysWithPrefix(ctx, genPrimaryCachePrefix(c.InstanceId, tableName))
}

func (c *Cache) BatchPrimaryKeyExists(ctx context.Context, tableName string, primaryKeys []string) (bool, error) {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, genPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
	}
	return c.cache.BatchKeyExist(ctx, cacheKeys)
}

func (c *Cache) SearchKeyExists(ctx context.Context, tableName string, sql string, vars ...interface{}) (bool, error) {
	cacheKey := genSearchCacheKey(c.InstanceId, tableName, sql, vars...)
	return c.cache.KeyExists(ctx, cacheKey)
}

func (c *Cache) BatchSetPrimaryKeyCache(ctx context.Context, tableName string, kvs []Kv) error {
	for idx, kv := range kvs {
		kvs[idx].Key = genPrimaryCacheKey(c.InstanceId, tableName, kv.Key)
	}
	return c.cache.BatchSetKeys(ctx, kvs)
}

func (c *Cache) SetSearchCache(ctx context.Context, cacheValue string, tableName string,
	sql string, vars ...interface{}) error {
	key := genSearchCacheKey(c.InstanceId, tableName, sql, vars...)
	return c.cache.SetKey(ctx, Kv{
		Key:   key,
		Value: cacheValue,
	})
}

func (c *Cache) GetSearchCache(ctx context.Context, tableName string, sql string, vars ...interface{}) (string, error) {
	key := genSearchCacheKey(c.InstanceId, tableName, sql, vars...)
	return c.cache.GetValue(ctx, key)
}

func (c *Cache) BatchGetPrimaryCache(ctx context.Context, tableName string, primaryKeys []string) ([]string, error) {
	cacheKeys := make([]string, 0, len(primaryKeys))
	for _, primaryKey := range primaryKeys {
		cacheKeys = append(cacheKeys, genPrimaryCacheKey(c.InstanceId, tableName, primaryKey))
	}
	return c.cache.BatchGetValues(ctx, cacheKeys)
}
