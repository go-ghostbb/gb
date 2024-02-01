package gbcache

import (
	"context"
	gbconv "github.com/Ghostbb-io/gb/util/gb_conv"
)

// Cache struct.
type Cache struct {
	localAdapter
}

// localAdapter is alias of Adapter, for embedded attribute purpose only.
type localAdapter = Adapter

// New creates and returns a new cache object using default memory adapter.
// Note that the LRU feature is only available using memory adapter.
func New(lruCap ...int) *Cache {
	memAdapter := NewAdapterMemory(lruCap...)
	c := &Cache{
		localAdapter: memAdapter,
	}
	return c
}

// NewWithAdapter creates and returns a Cache object with given Adapter implements.
func NewWithAdapter(adapter Adapter) *Cache {
	return &Cache{
		localAdapter: adapter,
	}
}

// SetAdapter changes the adapter for this cache.
// Be very note that, this setting function is not concurrent-safe, which means you should not call
// this setting function concurrently in multiple goroutines.
func (c *Cache) SetAdapter(adapter Adapter) {
	c.localAdapter = adapter
}

// GetAdapter returns the adapter that is set in current Cache.
func (c *Cache) GetAdapter() Adapter {
	return c.localAdapter
}

// Removes deletes `keys` in the cache.
func (c *Cache) Removes(ctx context.Context, keys []interface{}) error {
	_, err := c.Remove(ctx, keys...)
	return err
}

// KeyStrings returns all keys in the cache as string slice.
func (c *Cache) KeyStrings(ctx context.Context) ([]string, error) {
	keys, err := c.Keys(ctx)
	if err != nil {
		return nil, err
	}
	return gbconv.Strings(keys), nil
}
