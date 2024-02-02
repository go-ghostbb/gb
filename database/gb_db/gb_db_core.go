package gbdb

import (
	"context"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbcache "ghostbb.io/gb/os/gb_cache"
	gblog "ghostbb.io/gb/os/gb_log"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbutil "ghostbb.io/gb/util/gb_util"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strings"
	"time"
)

func (c *Core) GetLogger() *gblog.Logger {
	return c.logger.(*Logger).logger
}

func (c *Core) GetConfig() *ConfigNode {
	return c.config
}

func (c *Core) GormConfig() *gorm.Config {
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.config.TablePrefix,
			SingularTable: c.config.SingularTable,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	config.Logger = c.logger.LogMode(logger.Info)
	return config
}

func (c *Core) SetSlowThreshold(t time.Duration) {
	c.logger.SetSlowThreshold(t)
}

func (c *Core) SetIgnoreRecordNotFoundError(b bool) {
	c.logger.SetIgnoreRecordNotFoundError(b)
}

func (c *Core) UseCache(config *CacheConfig) error {
	cache, err := NewCache(config)
	if err != nil {
		return err
	}
	return c.DB.Use(cache)
}

func (c *Core) UseCacheWithMap(m map[string]interface{}) error {
	if len(m) == 0 {
		return gberror.NewCode(gbcode.CodeInvalidParameter, "configuration cannot be empty")
	}

	m = gbutil.MapCopy(m)

	// Change string configuration to int value for level.
	levelKey, levelValue := gbutil.MapPossibleItemByKey(m, "Level")
	if levelValue != nil {
		if level, ok := levelStringMap[strings.ToUpper(gbconv.String(levelValue))]; ok {
			m[levelKey] = level
		} else {
			return gberror.NewCodef(gbcode.CodeInvalidParameter, `invalid level string: %v`, levelValue)
		}
	}
	if c.config.CacheLevel != "" {
		if level, ok := levelStringMap[strings.ToUpper(c.config.CacheLevel)]; ok {
			m[levelKey] = level
		} else {
			return gberror.NewCodef(gbcode.CodeInvalidParameter, `invalid level string: %v`, levelValue)
		}
	}
	// ttl
	ttlKey, ttlValue := gbutil.MapPossibleItemByKey(m, "ttl")
	if ttlValue != nil {
		d, err := gbtime.ParseDuration(gbconv.String(ttlValue))
		if err != nil {
			return gberror.NewCodef(gbcode.CodeInvalidParameter, `invalid ttl string: %v`, ttlValue)
		}
		m[ttlKey] = d.Milliseconds()
	}

	config := new(CacheConfig)
	if err := gbconv.Struct(m, config); err != nil {
		return err
	}

	if c.config.CacheTTL.Milliseconds() != 0 {
		config.TTL = c.config.CacheTTL.Milliseconds()
	}
	if c.config.CacheMaxItemCnt != 0 {
		config.MaxItemCnt = c.config.CacheMaxItemCnt
	}

	c.cache = gbcache.New()
	config.Cache = c.cache

	return c.UseCache(config)
}

func (c *Core) SetCacheAdapter(adapter gbcache.Adapter) {
	if c.cache == nil {
		return
	}
	c.cache.SetAdapter(adapter)
	// clear cache
	if err := c.cache.Clear(context.TODO()); err != nil {
		panic(err)
	}
}
