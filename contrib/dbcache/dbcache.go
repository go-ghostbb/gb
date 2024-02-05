package dbcache

import (
	"context"
	"ghostbb.io/gb/contrib/dbcache/cache"
	"ghostbb.io/gb/contrib/dbcache/crud"
	gbcache "ghostbb.io/gb/os/gb_cache"
	"gorm.io/gorm"
)

const PluginName = "gb:gorm:cache"

var Context = context.WithValue(context.Background(), crud.CacheCtxKey, crud.CacheSearch)

func WithCtx(parent context.Context) context.Context {
	return context.WithValue(parent, crud.CacheCtxKey, crud.CacheSearch)
}

func New() *Plugin {
	return &Plugin{
		cache: cache.New(),
	}
}

type Plugin struct {
	cache *cache.Cache
}

func (c *Plugin) Name() string {
	return PluginName
}

func (c *Plugin) Initialize(db *gorm.DB) (err error) {
	if err = crud.New(c.cache).Bind(db); err != nil {
		return err
	}
	c.cache.ClearCache(context.TODO())
	return nil
}

func (c *Plugin) SetAdapter(adapter gbcache.Adapter) {
	c.cache.SetAdapter(adapter)
	c.cache.ClearCache(context.TODO())
}
