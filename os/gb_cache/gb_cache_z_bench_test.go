package gbcache_test

import (
	"context"
	gbcache "ghostbb.io/gb/os/gb_cache"
	"testing"
)

var (
	localCache    = gbcache.New()
	localCacheLru = gbcache.New(10000)
)

func Benchmark_CacheSet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCache.Set(ctx, i, i, 0)
			i++
		}
	})
}

func Benchmark_CacheGet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCache.Get(ctx, i)
			i++
		}
	})
}

func Benchmark_CacheRemove(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCache.Remove(ctx, i)
			i++
		}
	})
}

func Benchmark_CacheLruSet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCacheLru.Set(ctx, i, i, 0)
			i++
		}
	})
}

func Benchmark_CacheLruGet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCacheLru.Get(ctx, i)
			i++
		}
	})
}

func Benchmark_CacheLruRemove(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCacheLru.Remove(context.TODO(), i)
			i++
		}
	})
}
