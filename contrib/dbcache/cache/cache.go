package cache

import (
	gbcache "ghostbb.io/gb/os/gb_cache"
	gbrand "ghostbb.io/gb/util/gb_rand"
	"time"
)

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
