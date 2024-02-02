package gbdb

import (
	"context"
	"fmt"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbcache "ghostbb.io/gb/os/gb_cache"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"time"
)

var (
	ErrCacheNotFound = gberror.New("cache not found")
)

type IDataCache interface {
	Init() error
	ClearCache(ctx context.Context) error

	// read
	BatchKeyExist(ctx context.Context, keys []string) (bool, error)
	KeyExists(ctx context.Context, key string) (bool, error)
	GetValue(ctx context.Context, key string) (string, error)
	BatchGetValues(ctx context.Context, keys []string) ([]string, error)

	// write
	DeleteKeysWithPrefix(ctx context.Context, keyPrefix string) error
	DeleteKey(ctx context.Context, key string) error
	BatchDeleteKeys(ctx context.Context, keys []string) error
	BatchSetKeys(ctx context.Context, kvs []Kv) error
	SetKey(ctx context.Context, kv Kv) error
}

func newDataCache(config *CacheConfig) IDataCache {
	return &DataCache{
		cache: config.Cache,
		ttl:   config.TTL,
	}
}

type DataCache struct {
	cache *gbcache.Cache
	ttl   int64
}

func (d *DataCache) Init() error {
	// nothing here
	return nil
}

func (d *DataCache) ClearCache(ctx context.Context) error {
	return d.cache.Clear(ctx)
}

func (d *DataCache) BatchKeyExist(ctx context.Context, keys []string) (bool, error) {
	for _, key := range keys {
		item, err := d.cache.Get(ctx, key)
		if err != nil {
			return false, err
		}
		if item.IsNil() {
			return false, nil
		}
	}
	return true, nil
}

func (d *DataCache) KeyExists(ctx context.Context, key string) (bool, error) {
	item, err := d.cache.Get(ctx, key)
	if err != nil {
		return false, err
	}
	if item == nil || item.IsNil() {
		return false, nil
	}
	return true, nil
}

func (d *DataCache) GetValue(ctx context.Context, key string) (string, error) {
	item, err := d.cache.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if item == nil || item.IsNil() {
		return "", ErrCacheNotFound
	}
	return item.String(), nil
}

func (d *DataCache) BatchGetValues(ctx context.Context, keys []string) ([]string, error) {
	values := make([]string, 0, len(keys))
	for _, key := range keys {
		item, err := d.cache.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		if item != nil {
			values = append(values, item.String())
		}
	}
	if len(values) != len(keys) {
		return nil, fmt.Errorf("cannot get items")
	}
	return values, nil
}

func (d *DataCache) DeleteKeysWithPrefix(ctx context.Context, keyPrefix string) error {
	keys, err := d.cache.KeyStrings(ctx)
	if err != nil {
		return err
	}

	list := make([]interface{}, 0)
	for _, k := range keys {
		if gbstr.Contains(k, keyPrefix) {
			list = append(list, k)
		}
	}
	err = d.cache.Removes(ctx, list)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataCache) DeleteKey(ctx context.Context, key string) error {
	_, err := d.cache.Remove(ctx, key)
	return err
}

func (d *DataCache) BatchDeleteKeys(ctx context.Context, keys []string) error {
	return d.cache.Removes(ctx, gbconv.Interfaces(keys))
}

func (d *DataCache) BatchSetKeys(ctx context.Context, kvs []Kv) error {
	for _, kv := range kvs {
		if d.ttl > 0 {
			return d.cache.Set(ctx, kv.Key, kv.Value, time.Duration(randFloatingInt64(d.ttl))*time.Millisecond)
		} else {
			return d.cache.Set(ctx, kv.Key, kv.Value, time.Duration(randFloatingInt64(24))*time.Hour)
		}
	}
	return nil
}

func (d *DataCache) SetKey(ctx context.Context, kv Kv) error {
	if d.ttl > 0 {
		return d.cache.Set(ctx, kv.Key, kv.Value, time.Duration(randFloatingInt64(d.ttl))*time.Millisecond)
	} else {
		return d.cache.Set(ctx, kv.Key, kv.Value, time.Duration(randFloatingInt64(24))*time.Hour)
	}
}
