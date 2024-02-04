package cache

import (
	"context"
	"fmt"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"time"
)

func (c *Cache) ClearCache(ctx context.Context) error {
	return c.cache.Clear(ctx)
}

func (c *Cache) BatchKeyExist(ctx context.Context, keys []string) (bool, error) {
	for _, key := range keys {
		item, err := c.cache.Get(ctx, key)
		if err != nil {
			return false, err
		}
		if item.IsNil() {
			return false, nil
		}
	}
	return true, nil
}

func (c *Cache) KeyExists(ctx context.Context, key string) (bool, error) {
	item, err := c.cache.Get(ctx, key)
	if err != nil {
		return false, err
	}
	if item == nil || item.IsNil() {
		return false, nil
	}
	return true, nil
}

func (c *Cache) GetValue(ctx context.Context, key string) (string, error) {
	item, err := c.cache.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if item == nil || item.IsNil() {
		return "", ErrCacheNotFound
	}
	return item.String(), nil
}

func (c *Cache) BatchGetValues(ctx context.Context, keys []string) ([]string, error) {
	values := make([]string, 0, len(keys))
	for _, key := range keys {
		item, err := c.cache.Get(ctx, key)
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

func (c *Cache) DeleteKeysWithPrefix(ctx context.Context, keyPrefix string) error {
	keys, err := c.cache.KeyStrings(ctx)
	if err != nil {
		return err
	}

	list := make([]interface{}, 0)
	for _, k := range keys {
		if gbstr.Contains(k, keyPrefix) {
			list = append(list, k)
		}
	}
	err = c.cache.Removes(ctx, list)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) DeleteKey(ctx context.Context, key string) error {
	_, err := c.cache.Remove(ctx, key)
	return err
}

func (c *Cache) BatchDeleteKeys(ctx context.Context, keys []string) error {
	return c.cache.Removes(ctx, gbconv.Interfaces(keys))
}

func (c *Cache) BatchSetKeys(ctx context.Context, kvs []Kv) error {
	for _, kv := range kvs {
		return c.cache.Set(ctx, kv.Key, kv.Value, time.Duration(c.TTL())*time.Millisecond)
	}
	return nil
}

func (c *Cache) SetKey(ctx context.Context, kv Kv) error {
	return c.cache.Set(ctx, kv.Key, kv.Value, time.Duration(c.TTL())*time.Millisecond)
}
