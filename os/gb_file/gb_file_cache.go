package gbfile

import (
	"context"
	gbcode "github.com/Ghostbb-io/gb/errors/gb_code"
	gberror "github.com/Ghostbb-io/gb/errors/gb_error"
	"github.com/Ghostbb-io/gb/internal/command"
	"github.com/Ghostbb-io/gb/internal/intlog"
	gbcache "github.com/Ghostbb-io/gb/os/gb_cache"
	gbfsnotify "github.com/Ghostbb-io/gb/os/gb_fsnotify"
	"time"
)

const (
	defaultCacheDuration  = "1m"            // defaultCacheExpire is the expire time for file content caching in seconds.
	commandEnvKeyForCache = "gb.file.cache" // commandEnvKeyForCache is the configuration key for command argument or environment configuring cache expire duration.
)

var (
	// Default expire time for file content caching.
	cacheDuration = getCacheDuration()

	// internalCache is the memory cache for internal usage.
	internalCache = gbcache.New()
)

func getCacheDuration() time.Duration {
	cacheDurationConfigured := command.GetOptWithEnv(commandEnvKeyForCache, defaultCacheDuration)
	d, err := time.ParseDuration(cacheDurationConfigured)
	if err != nil {
		panic(gberror.WrapCodef(
			gbcode.CodeInvalidConfiguration,
			err,
			`error parsing string "%s" to time duration`,
			cacheDurationConfigured,
		))
	}
	return d
}

// GetContentsWithCache returns string content of given file by `path` from cache.
// If there's no content in the cache, it will read it from disk file specified by `path`.
// The parameter `expire` specifies the caching time for this file content in seconds.
func GetContentsWithCache(path string, duration ...time.Duration) string {
	return string(GetBytesWithCache(path, duration...))
}

// GetBytesWithCache returns []byte content of given file by `path` from cache.
// If there's no content in the cache, it will read it from disk file specified by `path`.
// The parameter `expire` specifies the caching time for this file content in seconds.
func GetBytesWithCache(path string, duration ...time.Duration) []byte {
	var (
		ctx      = context.Background()
		expire   = cacheDuration
		cacheKey = commandEnvKeyForCache + path
	)

	if len(duration) > 0 {
		expire = duration[0]
	}
	r, _ := internalCache.GetOrSetFuncLock(ctx, cacheKey, func(ctx context.Context) (interface{}, error) {
		b := GetBytes(path)
		if b != nil {
			// Adding this `path` to gfsnotify,
			// it will clear its cache if there's any changes of the file.
			_, _ = gbfsnotify.Add(path, func(event *gbfsnotify.Event) {
				_, err := internalCache.Remove(ctx, cacheKey)
				if err != nil {
					intlog.Errorf(ctx, `%+v`, err)
				}
				gbfsnotify.Exit()
			})
		}
		return b, nil
	}, expire)
	if r != nil {
		return r.Bytes()
	}
	return nil
}
