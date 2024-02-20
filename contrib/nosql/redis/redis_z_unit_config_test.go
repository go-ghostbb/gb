package redis_test

import (
	gbvar "ghostbb.io/gb/container/gb_var"
	gbredis "ghostbb.io/gb/database/gb_redis"
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
	"time"
)

func Test_ConfigFromMap(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		c, err := gbredis.ConfigFromMap(g.Map{
			`address`:     `127.0.0.1:6379`,
			`db`:          `10`,
			`pass`:        `&*^%$#65Gv`,
			`minIdle`:     `10`,
			`MaxIdle`:     `100`,
			`ReadTimeout`: `10s`,
		})
		t.AssertNil(err)
		t.Assert(c.Address, `127.0.0.1:6379`)
		t.Assert(c.Db, `10`)
		t.Assert(c.Pass, `&*^%$#65Gv`)
		t.Assert(c.MinIdle, 10)
		t.Assert(c.MaxIdle, 100)
		t.Assert(c.ReadTimeout, 10*time.Second)
	})
}

func Test_ConfigAddUser(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			c   *gbredis.Redis
			err error
			r   *gbvar.Var
		)

		c, err = gbredis.New(&gbredis.Config{
			Address: `127.0.0.1`,
			Db:      1,
			User:    "root",
			Pass:    "",
		})
		t.AssertNil(err)

		_, err = c.Conn(ctx)
		t.AssertNil(err)

		_, err = redis.Do(ctx, "SET", "k", "v")
		t.AssertNil(err)

		r, err = redis.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(r, []byte("v"))

		_, err = redis.Do(ctx, "DEL", "k")
		t.AssertNil(err)

		r, err = redis.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(r, nil)
	})
}
