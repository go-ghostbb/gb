package redis_test

import (
	gbvar "ghostbb.io/gb/container/gb_var"
	gbredis "ghostbb.io/gb/database/gb_redis"
	"ghostbb.io/gb/frame/g"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbuid "ghostbb.io/gb/util/gb_uid"
	gbutil "ghostbb.io/gb/util/gb_util"
	"testing"
	"time"
)

func Test_NewClose(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		redis, err := gbredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)

		err = redis.Close(ctx)
		t.AssertNil(err)
	})
}

func Test_Do(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		_, err := redis.Do(ctx, "SET", "k", "v")
		t.AssertNil(err)

		r, err := redis.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(r, []byte("v"))

		_, err = redis.Do(ctx, "DEL", "k")
		t.AssertNil(err)
		r, err = redis.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(r, nil)
	})
}

func Test_Conn(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := redis.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)

		key := gbconv.String(gbtime.TimestampNano())
		value := []byte("v")
		r, err := conn.Do(ctx, "SET", key, value)
		t.AssertNil(err)

		r, err = conn.Do(ctx, "GET", key)
		t.AssertNil(err)
		t.Assert(r, value)

		_, err = conn.Do(ctx, "DEL", key)
		t.AssertNil(err)
		r, err = conn.Do(ctx, "GET", key)
		t.AssertNil(err)
		t.Assert(r, nil)
	})
}

func Test_Instance(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		group := "my-test"
		gbredis.SetConfig(config, group)
		defer gbredis.RemoveConfig(group)

		redis := gbredis.Instance(group)
		defer redis.Close(ctx)

		conn, err := redis.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)

		_, err = conn.Do(ctx, "SET", "k", "v")
		t.AssertNil(err)

		r, err := conn.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(r, []byte("v"))

		_, err = conn.Do(ctx, "DEL", "k")
		t.AssertNil(err)
		r, err = conn.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(r, nil)
	})
}

func Test_Error(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		config1 := &gbredis.Config{
			Address:     "192.111.0.2:6379",
			Db:          1,
			DialTimeout: time.Second,
		}
		r, err := gbredis.New(config1)
		t.AssertNil(err)
		t.AssertNE(r, nil)
		defer r.Close(ctx)

		_, err = r.Do(ctx, "info")
		t.AssertNE(err, nil)

		config1 = &gbredis.Config{
			Address: "127.0.0.1:6379",
			Db:      100,
		}
		r, err = gbredis.New(config1)
		t.AssertNil(err)
		t.AssertNE(r, nil)
		defer r.Close(ctx)

		_, err = r.Do(ctx, "info")
		t.AssertNE(err, nil)

		r = gbredis.Instance("gb")
		t.Assert(r == nil, true)
		gbredis.ClearConfig()

		r, err = gbredis.New(config)
		t.AssertNil(err)
		t.AssertNE(r, nil)
		defer r.Close(ctx)

		_, err = r.Do(ctx, "SET", "k", "v")
		t.AssertNil(err)

		v, err := r.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(v.String(), "v")

		conn, err := r.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)
		_, err = conn.Do(ctx, "SET", "k", "v")
		t.AssertNil(err)
	})
}

func Test_Bool(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer func() {
			redis.Do(ctx, "DEL", "key-true")
			redis.Do(ctx, "DEL", "key-false")
		}()

		_, err := redis.Do(ctx, "SET", "key-true", true)
		t.AssertNil(err)

		_, err = redis.Do(ctx, "SET", "key-false", false)
		t.AssertNil(err)

		r, err := redis.Do(ctx, "GET", "key-true")
		t.AssertNil(err)
		t.Assert(r.Bool(), true)

		r, err = redis.Do(ctx, "GET", "key-false")
		t.AssertNil(err)
		t.Assert(r.Bool(), false)
	})
}

func Test_Int(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		redis, err := gbredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		key := gbuid.S()
		defer redis.Do(ctx, "DEL", key)

		_, err = redis.Do(ctx, "SET", key, 1)
		t.AssertNil(err)

		r, err := redis.Do(ctx, "GET", key)
		t.AssertNil(err)
		t.Assert(r.Int(), 1)
	})
}

func Test_HSet(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		redis, err := gbredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		key := gbuid.S()
		defer redis.Do(ctx, "DEL", key)

		_, err = redis.Do(ctx, "HSET", key, "name", "john")
		t.AssertNil(err)

		r, err := redis.Do(ctx, "HGETALL", key)
		t.AssertNil(err)
		t.Assert(r.MapStrStr(), g.MapStrStr{"name": "john"})
	})
}

func Test_HGetAll1(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			key = gbuid.S()
		)
		redis, err := gbredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)
		defer redis.Do(ctx, "DEL", key)

		_, err = redis.Do(ctx, "HSET", key, "id", 100)
		t.AssertNil(err)
		_, err = redis.Do(ctx, "HSET", key, "name", "john")
		t.AssertNil(err)

		r, err := redis.Do(ctx, "HGETALL", key)
		t.AssertNil(err)
		t.Assert(r.Map(), g.MapStrAny{
			"id":   100,
			"name": "john",
		})
	})
}

func Test_HGetAll2(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			key = gbuid.S()
		)
		redis, err := gbredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)
		defer redis.Do(ctx, "DEL", key)

		_, err = redis.Do(ctx, "HSET", key, "id", 100)
		t.AssertNil(err)
		_, err = redis.Do(ctx, "HSET", key, "name", "john")
		t.AssertNil(err)

		result, err := redis.Do(ctx, "HGETALL", key)
		t.AssertNil(err)

		t.Assert(gbconv.Uint(result.MapStrVar()["id"]), 100)
		t.Assert(result.MapStrVar()["id"].Uint(), 100)
	})
}

func Test_HMSet(t *testing.T) {
	// map
	gbtest.C(t, func(t *gbtest.T) {
		var (
			key  = gbuid.S()
			data = g.Map{
				"name":  "gb",
				"sex":   0,
				"score": 100,
			}
		)
		redis, err := gbredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)
		defer redis.Do(ctx, "DEL", key)

		_, err = redis.Do(ctx, "HMSET", append(g.Slice{key}, gbutil.MapToSlice(data)...)...)
		t.AssertNil(err)
		v, err := redis.Do(ctx, "HMGET", key, "name")
		t.AssertNil(err)
		t.Assert(v.Slice(), g.Slice{data["name"]})
	})
	// struct
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Name  string `json:"name"`
			Sex   int    `json:"sex"`
			Score int    `json:"score"`
		}
		var (
			key  = gbuid.S()
			data = &User{
				Name:  "gb",
				Sex:   0,
				Score: 100,
			}
		)
		redis, err := gbredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)
		defer redis.Do(ctx, "DEL", key)

		_, err = redis.Do(ctx, "HMSET", append(g.Slice{key}, gbutil.StructToSlice(data)...)...)
		t.AssertNil(err)
		v, err := redis.Do(ctx, "HMGET", key, "name")
		t.AssertNil(err)
		t.Assert(v.Slice(), g.Slice{data.Name})
	})
}

func Test_Auto_Marshal(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			key = gbuid.S()
		)
		redis, err := gbredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		defer redis.Do(ctx, "DEL", key)

		type User struct {
			Id   int
			Name string
		}

		user := &User{
			Id:   10000,
			Name: "john",
		}

		_, err = redis.Do(ctx, "SET", key, user)
		t.AssertNil(err)

		r, err := redis.Do(ctx, "GET", key)
		t.AssertNil(err)
		t.Assert(r.Map(), g.MapStrAny{
			"Id":   user.Id,
			"Name": user.Name,
		})

		var user2 *User
		t.Assert(r.Struct(&user2), nil)
		t.Assert(user2.Id, user.Id)
		t.Assert(user2.Name, user.Name)
	})
}

func Test_Auto_MarshalSlice(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			key = "user-slice"
		)
		redis, err := gbredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Do(ctx, "DEL", key)
		type User struct {
			Id   int
			Name string
		}
		var (
			result *gbvar.Var
			users1 = []User{
				{
					Id:   1,
					Name: "john1",
				},
				{
					Id:   2,
					Name: "john2",
				},
			}
		)

		_, err = redis.Do(ctx, "SET", key, users1)
		t.AssertNil(err)

		result, err = redis.Do(ctx, "GET", key)
		t.AssertNil(err)

		var users2 []User
		err = result.Structs(&users2)
		t.AssertNil(err)
		t.Assert(users2, users1)
	})
}
