package gbcache_test

import (
	"context"
	gbset "ghostbb.io/gb/container/gb_set"
	"ghostbb.io/gb/frame/g"
	gbcache "ghostbb.io/gb/os/gb_cache"
	gbrpool "ghostbb.io/gb/os/gb_rpool"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbuid "ghostbb.io/gb/util/gb_uid"
	"math"
	"testing"
	"time"
)

var (
	ctx = context.Background()
)

func TestCache_GCache_Set(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertNil(gbcache.Set(ctx, 1, 11, 0))
		defer gbcache.Remove(ctx, g.Slice{1, 2, 3}...)
		v, _ := gbcache.Get(ctx, 1)
		t.Assert(v, 11)
		b, _ := gbcache.Contains(ctx, 1)
		t.Assert(b, true)
	})
}

func TestCache_Set(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		c := gbcache.New()
		defer c.Close(ctx)
		t.Assert(c.Set(ctx, 1, 11, 0), nil)
		v, _ := c.Get(ctx, 1)
		t.Assert(v, 11)
		b, _ := c.Contains(ctx, 1)
		t.Assert(b, true)
	})
}

func TestCache_Set_Expire(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New()
		t.Assert(cache.Set(ctx, 2, 22, 100*time.Millisecond), nil)
		v, _ := cache.Get(ctx, 2)
		t.Assert(v, 22)
		time.Sleep(200 * time.Millisecond)
		v, _ = cache.Get(ctx, 2)
		t.Assert(v, nil)
		time.Sleep(3 * time.Second)
		n, _ := cache.Size(ctx)
		t.Assert(n, 0)
		t.Assert(cache.Close(ctx), nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New()
		t.Assert(cache.Set(ctx, 1, 11, 100*time.Millisecond), nil)
		v, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)
		time.Sleep(200 * time.Millisecond)
		v, _ = cache.Get(ctx, 1)
		t.Assert(v, nil)
	})
}

func TestCache_Update(t *testing.T) {
	// gbcache
	gbtest.C(t, func(t *gbtest.T) {
		key := gbuid.S()
		t.AssertNil(gbcache.Set(ctx, key, 11, 3*time.Second))
		expire1, _ := gbcache.GetExpire(ctx, key)
		oldValue, exist, err := gbcache.Update(ctx, key, 12)
		t.AssertNil(err)
		t.Assert(oldValue, 11)
		t.Assert(exist, true)

		expire2, _ := gbcache.GetExpire(ctx, key)
		v, _ := gbcache.Get(ctx, key)
		t.Assert(v, 12)
		t.Assert(math.Ceil(expire1.Seconds()), math.Ceil(expire2.Seconds()))
	})
	// gbcache.Cache
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New()
		t.AssertNil(cache.Set(ctx, 1, 11, 3*time.Second))

		oldValue, exist, err := cache.Update(ctx, 1, 12)
		t.AssertNil(err)
		t.Assert(oldValue, 11)
		t.Assert(exist, true)

		expire1, _ := cache.GetExpire(ctx, 1)
		expire2, _ := cache.GetExpire(ctx, 1)
		v, _ := cache.Get(ctx, 1)
		t.Assert(v, 12)
		t.Assert(math.Ceil(expire1.Seconds()), math.Ceil(expire2.Seconds()))
	})
}

func TestCache_UpdateExpire(t *testing.T) {
	// gbcache
	gbtest.C(t, func(t *gbtest.T) {
		key := gbuid.S()
		t.AssertNil(gbcache.Set(ctx, key, 11, 3*time.Second))
		defer gbcache.Remove(ctx, key)
		oldExpire, _ := gbcache.GetExpire(ctx, key)
		newExpire := 10 * time.Second
		oldExpire2, err := gbcache.UpdateExpire(ctx, key, newExpire)
		t.AssertNil(err)
		t.AssertIN(oldExpire2, g.Slice{oldExpire, `2.999s`})

		e, _ := gbcache.GetExpire(ctx, key)
		t.AssertNE(e, oldExpire)
		e, _ = gbcache.GetExpire(ctx, key)
		t.Assert(math.Ceil(e.Seconds()), 10)
	})
	// gbcache.Cache
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New()
		t.AssertNil(cache.Set(ctx, 1, 11, 3*time.Second))
		oldExpire, _ := cache.GetExpire(ctx, 1)
		newExpire := 10 * time.Second
		oldExpire2, err := cache.UpdateExpire(ctx, 1, newExpire)
		t.AssertNil(err)
		t.AssertIN(oldExpire2, g.Slice{oldExpire, `2.999s`})

		e, _ := cache.GetExpire(ctx, 1)
		t.AssertNE(e, oldExpire)

		e, _ = cache.GetExpire(ctx, 1)
		t.Assert(math.Ceil(e.Seconds()), 10)
	})
}

func TestCache_Keys_Values(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		c := gbcache.New()
		for i := 0; i < 10; i++ {
			t.Assert(c.Set(ctx, i, i*10, 0), nil)
		}
		var (
			keys, _   = c.Keys(ctx)
			values, _ = c.Values(ctx)
		)
		t.Assert(len(keys), 10)
		t.Assert(len(values), 10)
		t.AssertIN(0, keys)
		t.AssertIN(90, values)
	})
}

func TestCache_LRU(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New(2)
		for i := 0; i < 10; i++ {
			t.AssertNil(cache.Set(ctx, i, i, 0))
		}
		n, _ := cache.Size(ctx)
		t.Assert(n, 10)
		v, _ := cache.Get(ctx, 6)
		t.Assert(v, 6)
		time.Sleep(4 * time.Second)
		g.Log().Debugf(ctx, `items after lru: %+v`, cache.MustData(ctx))
		n, _ = cache.Size(ctx)
		t.Assert(n, 2)
		v, _ = cache.Get(ctx, 6)
		t.Assert(v, 6)
		v, _ = cache.Get(ctx, 1)
		t.Assert(v, nil)
		t.Assert(cache.Close(ctx), nil)
	})
}

func TestCache_LRU_expire(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New(2)
		t.Assert(cache.Set(ctx, 1, nil, 1000), nil)
		n, _ := cache.Size(ctx)
		t.Assert(n, 1)
		v, _ := cache.Get(ctx, 1)

		t.Assert(v, nil)
	})
}

func TestCache_SetIfNotExist(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New()
		ok, err := cache.SetIfNotExist(ctx, 1, 11, 0)
		t.AssertNil(err)
		t.Assert(ok, true)

		v, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)

		ok, err = cache.SetIfNotExist(ctx, 1, 22, 0)
		t.AssertNil(err)
		t.Assert(ok, false)

		v, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)

		ok, err = cache.SetIfNotExist(ctx, 2, 22, 0)
		t.AssertNil(err)
		t.Assert(ok, true)

		v, _ = cache.Get(ctx, 2)
		t.Assert(v, 22)

		gbcache.Remove(ctx, g.Slice{1, 2, 3}...)
		ok, err = gbcache.SetIfNotExist(ctx, 1, 11, 0)
		t.AssertNil(err)
		t.Assert(ok, true)

		v, _ = gbcache.Get(ctx, 1)
		t.Assert(v, 11)

		ok, err = gbcache.SetIfNotExist(ctx, 1, 22, 0)
		t.AssertNil(err)
		t.Assert(ok, false)

		v, _ = gbcache.Get(ctx, 1)
		t.Assert(v, 11)
	})
}

func TestCache_SetIfNotExistFunc(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New()
		exist, err := cache.SetIfNotExistFunc(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 11, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, true)

		v, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)

		exist, err = cache.SetIfNotExistFunc(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 22, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, false)

		v, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
	})
	gbtest.C(t, func(t *gbtest.T) {
		gbcache.Remove(ctx, g.Slice{1, 2, 3}...)

		ok, err := gbcache.SetIfNotExistFunc(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 11, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(ok, true)

		v, _ := gbcache.Get(ctx, 1)
		t.Assert(v, 11)

		ok, err = gbcache.SetIfNotExistFunc(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 22, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(ok, false)

		v, _ = gbcache.Get(ctx, 1)
		t.Assert(v, 11)
	})
}

func TestCache_SetIfNotExistFuncLock(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New()
		exist, err := cache.SetIfNotExistFuncLock(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 11, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, true)

		v, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)

		exist, err = cache.SetIfNotExistFuncLock(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 22, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, false)

		v, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
	})
	gbtest.C(t, func(t *gbtest.T) {
		gbcache.Remove(ctx, g.Slice{1, 2, 3}...)

		exist, err := gbcache.SetIfNotExistFuncLock(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 11, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, true)

		v, _ := gbcache.Get(ctx, 1)
		t.Assert(v, 11)

		exist, err = gbcache.SetIfNotExistFuncLock(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 22, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, false)

		v, _ = gbcache.Get(ctx, 1)
		t.Assert(v, 11)
	})
}

func TestCache_SetMap(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New()
		t.AssertNil(cache.SetMap(ctx, g.MapAnyAny{1: 11, 2: 22}, 0))
		v, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)

		gbcache.Remove(ctx, g.Slice{1, 2, 3}...)
		t.AssertNil(gbcache.SetMap(ctx, g.MapAnyAny{1: 11, 2: 22}, 0))
		v, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
	})
}

func TestCache_GetOrSet(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New()
		value, err := cache.GetOrSet(ctx, 1, 11, 0)
		t.AssertNil(err)
		t.Assert(value, 11)

		v, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)
		value, err = cache.GetOrSet(ctx, 1, 111, 0)
		t.AssertNil(err)
		t.Assert(value, 11)

		v, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
	})

	gbtest.C(t, func(t *gbtest.T) {
		gbcache.Remove(ctx, g.Slice{1, 2, 3}...)
		value, err := gbcache.GetOrSet(ctx, 1, 11, 0)
		t.AssertNil(err)
		t.Assert(value, 11)

		v, err := gbcache.Get(ctx, 1)
		t.AssertNil(err)
		t.Assert(v, 11)

		value, err = gbcache.GetOrSet(ctx, 1, 111, 0)
		t.AssertNil(err)
		t.Assert(value, 11)

		v, err = gbcache.Get(ctx, 1)
		t.AssertNil(err)
		t.Assert(v, 11)
	})
}

func TestCache_GetOrSetFunc(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New()
		cache.GetOrSetFunc(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 11, nil
		}, 0)
		v, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)

		cache.GetOrSetFunc(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)

		gbcache.Remove(ctx, g.Slice{1, 2, 3}...)

		gbcache.GetOrSetFunc(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 11, nil
		}, 0)
		v, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)

		gbcache.GetOrSetFunc(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
	})
}

func TestCache_GetOrSetFuncLock(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New()
		cache.GetOrSetFuncLock(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 11, nil
		}, 0)
		v, _ := cache.Get(ctx, 1)
		t.Assert(v, 11)

		cache.GetOrSetFuncLock(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)

		gbcache.Remove(ctx, g.Slice{1, 2, 3}...)
		gbcache.GetOrSetFuncLock(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 11, nil
		}, 0)
		v, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)

		gbcache.GetOrSetFuncLock(ctx, 1, func(ctx context.Context) (value interface{}, err error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(ctx, 1)
		t.Assert(v, 11)
	})
}

func TestCache_Clear(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New()
		cache.SetMap(ctx, g.MapAnyAny{1: 11, 2: 22}, 0)
		cache.Clear(ctx)
		n, _ := cache.Size(ctx)
		t.Assert(n, 0)
	})
}

func TestCache_SetConcurrency(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New()
		pool := gbrpool.New(4)
		go func() {
			for {
				pool.Add(ctx, func(ctx context.Context) {
					cache.SetIfNotExist(ctx, 1, 11, 10)
				})
			}
		}()
		select {
		case <-time.After(2 * time.Second):
			// t.Log("first part end")
		}

		go func() {
			for {
				pool.Add(ctx, func(ctx context.Context) {
					cache.SetIfNotExist(ctx, 1, nil, 10)
				})
			}
		}()
		select {
		case <-time.After(2 * time.Second):
			// t.Log("second part end")
		}
	})
}

func TestCache_Basic(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		{
			cache := gbcache.New()
			cache.SetMap(ctx, g.MapAnyAny{1: 11, 2: 22}, 0)
			b, _ := cache.Contains(ctx, 1)
			t.Assert(b, true)
			v, _ := cache.Get(ctx, 1)
			t.Assert(v, 11)
			data, _ := cache.Data(ctx)
			t.Assert(data[1], 11)
			t.Assert(data[2], 22)
			t.Assert(data[3], nil)
			n, _ := cache.Size(ctx)
			t.Assert(n, 2)
			keys, _ := cache.Keys(ctx)
			t.Assert(gbset.NewFrom(g.Slice{1, 2}).Equal(gbset.NewFrom(keys)), true)
			keyStrs, _ := cache.KeyStrings(ctx)
			t.Assert(gbset.NewFrom(g.Slice{"1", "2"}).Equal(gbset.NewFrom(keyStrs)), true)
			values, _ := cache.Values(ctx)
			t.Assert(gbset.NewFrom(g.Slice{11, 22}).Equal(gbset.NewFrom(values)), true)
			removeData1, _ := cache.Remove(ctx, 1)
			t.Assert(removeData1, 11)
			n, _ = cache.Size(ctx)
			t.Assert(n, 1)

			cache.Remove(ctx, 2)
			n, _ = cache.Size(ctx)
			t.Assert(n, 0)
		}

		gbcache.Remove(ctx, g.Slice{1, 2, 3}...)
		{
			gbcache.SetMap(ctx, g.MapAnyAny{1: 11, 2: 22}, 0)
			b, _ := gbcache.Contains(ctx, 1)
			t.Assert(b, true)
			v, _ := gbcache.Get(ctx, 1)
			t.Assert(v, 11)
			data, _ := gbcache.Data(ctx)
			t.Assert(data[1], 11)
			t.Assert(data[2], 22)
			t.Assert(data[3], nil)
			n, _ := gbcache.Size(ctx)
			t.Assert(n, 2)
			keys, _ := gbcache.Keys(ctx)
			t.Assert(gbset.NewFrom(g.Slice{1, 2}).Equal(gbset.NewFrom(keys)), true)
			keyStrs, _ := gbcache.KeyStrings(ctx)
			t.Assert(gbset.NewFrom(g.Slice{"1", "2"}).Equal(gbset.NewFrom(keyStrs)), true)
			values, _ := gbcache.Values(ctx)
			t.Assert(gbset.NewFrom(g.Slice{11, 22}).Equal(gbset.NewFrom(values)), true)
			removeData1, _ := gbcache.Remove(ctx, 1)
			t.Assert(removeData1, 11)
			n, _ = gbcache.Size(ctx)
			t.Assert(n, 1)
			gbcache.Remove(ctx, 2)
			n, _ = gbcache.Size(ctx)
			t.Assert(n, 0)
		}
	})
}

func TestCache_Removes(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.New()
		t.AssertNil(cache.Set(ctx, 1, 11, 0))
		t.AssertNil(cache.Set(ctx, 2, 22, 0))
		t.AssertNil(cache.Set(ctx, 3, 33, 0))
		t.AssertNil(cache.Removes(ctx, g.Slice{2, 3}))

		ok, err := cache.Contains(ctx, 1)
		t.AssertNil(err)
		t.Assert(ok, true)

		ok, err = cache.Contains(ctx, 2)
		t.AssertNil(err)
		t.Assert(ok, false)
	})

	gbtest.C(t, func(t *gbtest.T) {
		t.AssertNil(gbcache.Set(ctx, 1, 11, 0))
		t.AssertNil(gbcache.Set(ctx, 2, 22, 0))
		t.AssertNil(gbcache.Set(ctx, 3, 33, 0))
		t.AssertNil(gbcache.Removes(ctx, g.Slice{2, 3}))

		ok, err := gbcache.Contains(ctx, 1)
		t.AssertNil(err)
		t.Assert(ok, true)

		ok, err = gbcache.Contains(ctx, 2)
		t.AssertNil(err)
		t.Assert(ok, false)
	})
}

func TestCache_Basic_Must(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer gbcache.Remove(ctx, g.Slice{1, 2, 3, 4}...)

		t.AssertNil(gbcache.Set(ctx, 1, 11, 0))
		v := gbcache.MustGet(ctx, 1)
		t.Assert(v, 11)
		gbcache.MustGetOrSet(ctx, 2, 22, 0)
		v = gbcache.MustGet(ctx, 2)
		t.Assert(v, 22)

		gbcache.MustGetOrSetFunc(ctx, 3, func(ctx context.Context) (value interface{}, err error) {
			return 33, nil
		}, 0)
		v = gbcache.MustGet(ctx, 3)
		t.Assert(v, 33)

		gbcache.GetOrSetFuncLock(ctx, 4, func(ctx context.Context) (value interface{}, err error) {
			return 44, nil
		}, 0)
		v = gbcache.MustGet(ctx, 4)
		t.Assert(v, 44)

		t.Assert(gbcache.MustContains(ctx, 1), true)

		t.AssertNil(gbcache.Set(ctx, 1, 11, 3*time.Second))
		expire := gbcache.MustGetExpire(ctx, 1)
		t.AssertGE(expire, 0)

		n := gbcache.MustSize(ctx)
		t.Assert(n, 4)

		data := gbcache.MustData(ctx)
		t.Assert(len(data), 4)

		keys := gbcache.MustKeys(ctx)
		t.Assert(len(keys), 4)

		keyStrings := gbcache.MustKeyStrings(ctx)
		t.Assert(len(keyStrings), 4)

		values := gbcache.MustValues(ctx)
		t.Assert(len(values), 4)
	})
}

func TestCache_NewWithAdapter(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		cache := gbcache.NewWithAdapter(gbcache.NewAdapterMemory())
		t.AssertNE(cache, nil)
	})
}
