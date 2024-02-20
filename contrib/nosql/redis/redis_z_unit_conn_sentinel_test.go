package redis_test

import (
	"context"
	gbredis "ghostbb.io/gb/database/gb_redis"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

var (
	sentinelCtx    = context.TODO()
	sentinelConfig = &gbredis.Config{
		Address:    `127.0.0.1:26379,127.0.0.1:26380,127.0.0.1:26381`,
		MasterName: `mymaster`,
		Pass:       "111111",
	}
)

func TestConn_sentinel_master(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		sentinelConfig.SlaveOnly = false
		redis, err := gbredis.New(sentinelConfig)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(sentinelCtx)

		conn, err := redis.Conn(sentinelCtx)
		t.AssertNil(err)
		defer conn.Close(sentinelCtx)

		_, err = conn.Do(sentinelCtx, "set", "test", "123")
		t.AssertNil(err)
		defer conn.Do(sentinelCtx, "del", "test")

		r, err := conn.Do(sentinelCtx, "get", "test")
		t.AssertNil(err)
		t.Assert(r.String(), "123")
	})
}

func TestConn_sentinel_slave(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		sentinelConfig.SlaveOnly = true
		redis, err := gbredis.New(sentinelConfig)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(sentinelCtx)

		conn, err := redis.Conn(sentinelCtx)
		t.AssertNil(err)
		defer conn.Close(sentinelCtx)

		_, err = conn.Do(sentinelCtx, "set", "test", "123")
		t.AssertNQ(err, nil)
	})
}
