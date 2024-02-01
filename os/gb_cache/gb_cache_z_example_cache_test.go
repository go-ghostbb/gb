// go test *.go -bench=".*" -benchmem

package gbcache_test

import (
	"context"
	"fmt"
	gbredis "ghostbb.io/database/gb_redis"
	gbcache "ghostbb.io/os/gb_cache"
	gbctx "ghostbb.io/os/gb_ctx"
	"time"
)

var (
	ctx = context.Background()
)

func ExampleCache_MustContains() {

	// Create a cache object,
	// Of course, you can also easily use the gbcache package method directly
	c := gbcache.New()

	// Set Cache
	c.Set(ctx, "k", "v", 0)

	// MustContains returns true if `key` exists in the cache, or else returns false.
	// return true
	data := c.MustContains(ctx, "k")
	fmt.Println(data)

	// return false
	data1 := c.MustContains(ctx, "k1")
	fmt.Println(data1)

	// Output:
	// true
	// false
}

func ExampleCache_SetAdapter() {
	var (
		err         error
		ctx         = gbctx.New()
		cache       = gbcache.New()
		redisConfig = &gbredis.Config{
			Address: "127.0.0.1:6379",
			Db:      0,
		}
		cacheKey   = `ping`
		cacheValue = `pong`
	)
	// Create redis client object.
	redis, err := gbredis.New(redisConfig)
	if err != nil {
		panic(err)
	}
	// Create redis cache adapter and set it to cache object.
	cache.SetAdapter(gbcache.NewAdapterRedis(redis))

	// Set and Get using cache object.
	err = cache.Set(ctx, cacheKey, cacheValue, time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println(cache.MustGet(ctx, cacheKey).String())

	// Get using redis client.
	fmt.Println(redis.MustDo(ctx, "GET", cacheKey).String())

	// May Output:
	// value
	// value
}
