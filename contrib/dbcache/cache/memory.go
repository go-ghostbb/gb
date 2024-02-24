package cache

import (
	"context"
	gbcache "ghostbb.io/gb/os/gb_cache"
	gbstr "ghostbb.io/gb/text/gb_str"
)

type Memory struct {
	*gbcache.Cache
	config *Config
}

func (m *Memory) setConfig(config *Config) {
	m.config = config
}

func (m *Memory) Clear(ctx context.Context) error {
	return m.Cache.Clear(ctx)
}

func (m *Memory) Get(ctx context.Context, key string) (string, error) {
	item, err := m.Cache.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if item == nil || item.IsNil() {
		return "", ErrCacheNotFound
	}
	return item.String(), nil
}

func (m *Memory) Set(ctx context.Context, key string, value string) error {
	return m.Cache.Set(ctx, key, value, randTTL(m.config.TTL))
}

func (m *Memory) DeleteKeysWithPrefix(ctx context.Context, keyPrefix string) error {
	keys, err := m.Cache.KeyStrings(ctx)
	if err != nil {
		return err
	}

	list := make([]interface{}, 0)
	for _, k := range keys {
		if gbstr.Contains(k, keyPrefix) {
			list = append(list, k)
		}
	}
	err = m.Cache.Removes(ctx, list)
	if err != nil {
		return err
	}
	return nil
}
