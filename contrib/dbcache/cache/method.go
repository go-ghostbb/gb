package cache

import (
	"context"
	gbstr "ghostbb.io/gb/text/gb_str"
	"time"
)

func (c *Cache) ClearCache(ctx context.Context) error {
	return c.cache.Clear(ctx)
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

func (c *Cache) SetKey(ctx context.Context, kv Kv) error {
	return c.cache.Set(ctx, kv.Key, kv.Value, time.Duration(c.TTL())*time.Millisecond)
}
