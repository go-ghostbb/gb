package dbcache

import (
	"context"
	"ghostbb.io/gb/contrib/dbcache/cache"
	"ghostbb.io/gb/contrib/dbcache/crud"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	"gorm.io/gorm"
)

const PluginName = "gb:gorm:cache"

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
	return nil
}

func WithLevelAllCtx() context.Context {
	return context.WithValue(gbctx.New(), crud.CacheCtxLevelKey, crud.CacheAll)
}