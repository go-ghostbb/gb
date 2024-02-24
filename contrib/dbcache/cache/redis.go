package cache

import (
	"context"
	gbredis "ghostbb.io/gb/database/gb_redis"
	gbconv "ghostbb.io/gb/util/gb_conv"
)

type Redis struct {
	*gbredis.Redis
	config *Config
}

func (r *Redis) setConfig(config *Config) {
	r.config = config
}

func (r *Redis) Clear(ctx context.Context) error {
	return r.DeleteKeysWithPrefix(ctx, cacheName)
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	item, err := r.Redis.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if item == nil || item.IsNil() {
		return "", ErrCacheNotFound
	}
	return item.String(), nil
}

func (r *Redis) Set(ctx context.Context, key string, value string) error {
	_, err := r.Redis.Set(ctx, key, value, gbredis.SetOption{TTLOption: gbredis.TTLOption{
		PX: gbconv.PtrInt64(randTTL(r.config.TTL).Milliseconds()),
	}})
	return err
}

func (r *Redis) DeleteKeysWithPrefix(ctx context.Context, keyPrefix string) error {
	keys, err := r.Redis.Keys(ctx, keyPrefix+":*")
	if err != nil {
		return err
	}
	for _, key := range keys {
		if _, err = r.Redis.Del(ctx, key); err != nil {
			return err
		}
	}
	return nil
}
