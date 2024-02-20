package redis_test

import (
	gbredis "ghostbb.io/gb/database/gb_redis"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"strings"
	"testing"
)

func Test_GroupList_LPush(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1)
		t.AssertNil(err)
		_, err = redis.GroupList().LPush(ctx, k1, v2)
		t.AssertNil(err)

		r1, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r1, []string{v2, v1})
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2)
		t.AssertNil(err)

		r1, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r1, []string{v2, v1})
	})
}

func Test_GroupList_LPushX(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
		)
		_, err := redis.GroupList().LPushX(ctx, k1, v1)
		t.AssertNil(err)

		_, err = redis.GroupList().LPush(ctx, k1, v2)
		t.AssertNil(err)
		_, err = redis.GroupList().LPushX(ctx, k1, v1)
		t.AssertNil(err)

		r1, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r1, []string{v1, v2})
	})
}

func Test_GroupList_RPush(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
		)
		_, err := redis.GroupList().RPush(ctx, k1, v1)
		t.AssertNil(err)
		_, err = redis.GroupList().RPush(ctx, k1, v2)
		t.AssertNil(err)

		r1, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r1, []string{v1, v2})
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
		)
		_, err := redis.GroupList().RPush(ctx, k1, v1, v2)
		t.AssertNil(err)

		r1, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r1, []string{v1, v2})
	})
}

func Test_GroupList_RPushX(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
		)
		_, err := redis.GroupList().RPushX(ctx, k1, v1)
		t.AssertNil(err)

		_, err = redis.GroupList().RPush(ctx, k1, v2)
		t.AssertNil(err)
		_, err = redis.GroupList().RPushX(ctx, k1, v1)
		t.AssertNil(err)

		r1, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r1, []string{v2, v1})
	})
}

func InfoServerMap() map[string]string {
	v, err := redis.Do(ctx, "INFO", "server")
	if err != nil {
		return nil
	}
	server := make(map[string]string)
	list := strings.Split(v.String(), "\r\n")
	for _, v := range list {
		if strings.Contains(v, ":") {
			kv := strings.Split(v, ":")
			if len(kv) == 2 {
				server[kv[0]] = kv[1]
			}
		}
	}
	return server
}

func GetRedisVersion() string {
	svr := InfoServerMap()
	if svr != nil {
		return svr["redis_version"]
	}
	return ""
}

func Test_GroupList_LPop(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
		t.AssertNil(err)

		r1, err := redis.GroupList().LPop(ctx, k1)
		t.AssertNil(err)
		t.Assert(r1, v3)

		r3, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r3, []string{v2, v1})
	})

	// redis version check
	if gbstr.CompareVersion(GetRedisVersion(), "6.2.0") > 0 {
		gbtest.C(t, func(t *gbtest.T) {
			defer redis.FlushDB(ctx)
			var (
				k1 = "k1"
				v1 = "v1"
				v2 = "v2"
				v3 = "v3"
			)
			_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
			t.AssertNil(err)

			r1, err := redis.GroupList().LPop(ctx, k1, 2)
			t.AssertNil(err)
			t.Assert(r1, []string{v3, v2})

			r3, err := redis.GroupList().LRange(ctx, k1, 0, -1)
			t.AssertNil(err)
			t.Assert(r3, []string{v1})
		})
	}
}

func Test_GroupList_RPop(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
		t.AssertNil(err)

		r1, err := redis.GroupList().RPop(ctx, k1)
		t.AssertNil(err)
		t.Assert(r1, v1)

		r2, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r2, []string{v3, v2})
	})

	// redis version check
	if gbstr.CompareVersion(GetRedisVersion(), "6.2.0") > 0 {
		gbtest.C(t, func(t *gbtest.T) {
			defer redis.FlushDB(ctx)
			var (
				k1 = "k1"
				v1 = "v1"
				v2 = "v2"
				v3 = "v3"
			)
			_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
			t.AssertNil(err)

			r1, err := redis.GroupList().RPop(ctx, k1, 2)
			t.AssertNil(err)
			t.Assert(r1, []string{v1, v2})

			r3, err := redis.GroupList().LRange(ctx, k1, 0, -1)
			t.AssertNil(err)
			t.Assert(r3, []string{v3})
		})
	}
}

func Test_GroupList_LRem(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v1)
		t.AssertNil(err)

		r1, err := redis.GroupList().LRem(ctx, k1, 1, v1)
		t.AssertNil(err)
		t.Assert(r1, int64(1))

		r2, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r2, []string{v2, v1})
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v1)
		t.AssertNil(err)

		r1, err := redis.GroupList().LRem(ctx, k1, -1, v1)
		t.AssertNil(err)
		t.Assert(r1, int64(1))

		r2, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r2, []string{v1, v2})
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v1)
		t.AssertNil(err)

		r1, err := redis.GroupList().LRem(ctx, k1, 0, v1)
		t.AssertNil(err)
		t.Assert(r1, int64(2))

		r2, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r2, []string{v2})
	})
}

func Test_GroupList_LLen(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
		t.AssertNil(err)

		r1, err := redis.GroupList().LLen(ctx, k1)
		t.AssertNil(err)
		t.Assert(r1, int64(3))
	})
}

func Test_GroupList_LIndex(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
		t.AssertNil(err)

		r1, err := redis.GroupList().LIndex(ctx, k1, 1)
		t.AssertNil(err)
		t.Assert(r1, v2)
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
		t.AssertNil(err)

		r1, err := redis.GroupList().LIndex(ctx, k1, -2)
		t.AssertNil(err)
		t.Assert(r1, v2)
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
		t.AssertNil(err)

		r1, err := redis.GroupList().LIndex(ctx, k1, 3)
		t.AssertNil(err)
		t.AssertNil(r1)
	})
}

func Test_GroupList_LInsert(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
		t.AssertNil(err)

		r1, err := redis.GroupList().LInsert(ctx, k1, gbredis.LInsertBefore, v2, v1)
		t.AssertNil(err)
		t.Assert(r1, int64(4))

		r2, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r2, []string{v3, v1, v2, v1})
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
		t.AssertNil(err)

		r1, err := redis.GroupList().LInsert(ctx, k1, gbredis.LInsertAfter, v2, v1)
		t.AssertNil(err)
		t.Assert(r1, int64(4))

		r2, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r2, []string{v3, v2, v1, v1})
	})
}

func Test_GroupList_LSet(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
		t.AssertNil(err)

		r1, err := redis.GroupList().LSet(ctx, k1, 1, v1)
		t.AssertNil(err)
		t.Assert(r1, "OK")

		r2, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r2, []string{v3, v1, v1})
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
		t.AssertNil(err)

		r1, err := redis.GroupList().LSet(ctx, k1, -2, v1)
		t.AssertNil(err)
		t.Assert(r1, "OK")

		r2, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r2, []string{v3, v1, v1})
	})
}

func Test_GroupList_LRange(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
		t.AssertNil(err)

		r1, err := redis.GroupList().LRange(ctx, k1, 0, 1)
		t.AssertNil(err)
		t.Assert(r1, []string{v3, v2})
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
		t.AssertNil(err)

		r1, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r1, []string{v3, v2, v1})
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
		t.AssertNil(err)

		r1, err := redis.GroupList().LRange(ctx, k1, 0, 100)
		t.AssertNil(err)
		t.Assert(r1, []string{v3, v2, v1})
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
		t.AssertNil(err)

		r1, err := redis.GroupList().LRange(ctx, k1, 10, 100)
		t.AssertNil(err)
		t.AssertNil(r1)
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3)
		t.AssertNil(err)

		r1, err := redis.GroupList().LRange(ctx, k1, -3, -2)
		t.AssertNil(err)
		t.Assert(r1, []string{v3, v2})
	})
}

func Test_GroupList_LTrim(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
			v4 = "v4"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3, v4)
		t.AssertNil(err)

		err = redis.GroupList().LTrim(ctx, k1, 1, 2)
		t.AssertNil(err)

		r2, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r2, []string{v3, v2})
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
			v4 = "v4"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3, v4)
		t.AssertNil(err)

		err = redis.GroupList().LTrim(ctx, k1, 5, 10)
		t.AssertNil(err)

		r2, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.AssertNil(r2)
	})

	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
			v4 = "v4"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3, v4)
		t.AssertNil(err)

		err = redis.GroupList().LTrim(ctx, k1, -3, -2)
		t.AssertNil(err)

		r2, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r2, []string{v3, v2})
	})
}

func Test_GroupList_BLPop(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			k2 = "k2"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
			v4 = "v4"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3, v4)
		t.AssertNil(err)

		r1, err := redis.GroupList().BLPop(ctx, 1, k1, k2)
		t.AssertNil(err)
		t.Assert(r1, []string{k1, v4})
	})
}

func Test_GroupList_BRPop(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			k2 = "k2"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
			v4 = "v4"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3, v4)
		t.AssertNil(err)

		r1, err := redis.GroupList().BRPop(ctx, 1, k1, k2)
		t.AssertNil(err)
		t.Assert(r1, []string{k1, v1})
	})
}

func Test_GroupList_RPopLPush(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			k2 = "k2"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
			v4 = "v4"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3, v4)
		t.AssertNil(err)

		r1, err := redis.GroupList().RPopLPush(ctx, k1, k2)
		t.AssertNil(err)
		t.Assert(r1, v1)

		r2, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r2, []string{v4, v3, v2})

		r3, err := redis.GroupList().LRange(ctx, k2, 0, -1)
		t.AssertNil(err)
		t.Assert(r3, []string{v1})
	})
}

func Test_GroupList_BRPopLPush(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			k2 = "k2"
			v1 = "v1"
			v2 = "v2"
			v3 = "v3"
			v4 = "v4"
		)
		_, err := redis.GroupList().LPush(ctx, k1, v1, v2, v3, v4)
		t.AssertNil(err)

		r1, err := redis.GroupList().BRPopLPush(ctx, k1, k2, 1)
		t.AssertNil(err)
		t.Assert(r1, v1)

		r2, err := redis.GroupList().LRange(ctx, k1, 0, -1)
		t.AssertNil(err)
		t.Assert(r2, []string{v4, v3, v2})

		r3, err := redis.GroupList().LRange(ctx, k2, 0, -1)
		t.AssertNil(err)
		t.Assert(r3, []string{v1})
	})
}
