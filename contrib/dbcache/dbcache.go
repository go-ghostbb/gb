package dbcache

import (
	"context"
	"ghostbb.io/gb/contrib/dbcache/cache"
	"ghostbb.io/gb/contrib/dbcache/crud"
	gbredis "ghostbb.io/gb/database/gb_redis"
	gbcache "ghostbb.io/gb/os/gb_cache"
	"gorm.io/gorm"
)

const PluginName = "gb:gorm:cache"

var Context = context.WithValue(context.Background(), crud.CacheCtxKey, crud.CacheSearch)

func WithCtx(parent context.Context) context.Context {
	return context.WithValue(parent, crud.CacheCtxKey, crud.CacheSearch)
}

func New(c cache.ICache) *Plugin {
	return &Plugin{
		cache: cache.New(c),
	}
}

type Plugin struct {
	cache *cache.Cache
}

func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Initialize(db *gorm.DB) (err error) {
	if err = crud.New(p.cache).Bind(db); err != nil {
		return err
	}

	if name, ok := db.Get("gb:database:name"); ok {
		p.cache.SetName(name.(string))
	}
	return p.cache.Clear(context.TODO())
}

func NewMemory(c *gbcache.Cache) *cache.Memory {
	return &cache.Memory{
		Cache: c,
	}
}

func NewRedis(r *gbredis.Redis) *cache.Redis {
	return &cache.Redis{
		Redis: r,
	}
}
