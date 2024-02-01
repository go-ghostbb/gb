package gbdb

import gbcache "github.com/Ghostbb-io/gb/os/gb_cache"

type CacheConfig struct {
	// CacheLevel there are 2 types of cache and 4 kinds of cache option
	Level int `json:"level"`

	// CacheStorage choose proper storage medium
	Cache *gbcache.Cache `json:"cache"`

	// Tables only cache data within given data tables (cache all if empty)
	Tables []string `json:"tables"`

	// InvalidateWhenUpdate
	// if user update/delete/create something in DB, we invalidate all cached data to ensure consistency,
	// else we do nothing to outdated cache.
	InvalidateWhenUpdate bool `json:"invalidateWhenUpdate"`

	// AsyncWrite if true, then we will write cache in async mode
	AsyncWrite bool `json:"asyncWrite"`

	// CacheTTL cache ttl in ms, where 0 represents forever
	TTL int64 `json:"ttl"`

	// CacheMaxItemCnt for given query, if objects retrieved are more than this cnt,
	// then we choose not to cache for this query. 0 represents caching all queries.
	MaxItemCnt int64 `json:"maxItemCnt"`

	// DisableCachePenetration if true, then we will not cache nil result
	DisableCachePenetrationProtect bool `json:"disableCachePenetrationProtect"`
}

const (
	CacheLevelOff         int = 0
	CacheLevelOnlyPrimary int = 1
	CacheLevelOnlySearch  int = 2
	CacheLevelAll         int = 3
)

var levelStringMap = map[string]int{
	"OFF":     CacheLevelOff,
	"PRIMARY": CacheLevelOnlyPrimary,
	"SEARCH":  CacheLevelOnlySearch,
	"ALL":     CacheLevelAll,
}
