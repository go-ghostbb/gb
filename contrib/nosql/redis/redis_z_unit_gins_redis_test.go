package redis_test

import (
	"ghostbb.io/gb/frame/gins"
	gbcfg "ghostbb.io/gb/os/gb_cfg"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
	"time"
)

func Test_GINS_Redis(t *testing.T) {
	redisContent := gbfile.GetContents(
		gbtest.DataPath("redis", "config.toml"),
	)

	gbtest.C(t, func(t *gbtest.T) {
		var err error
		dirPath := gbfile.Temp(gbtime.TimestampNanoStr())
		err = gbfile.Mkdir(dirPath)
		t.AssertNil(err)
		defer gbfile.Remove(dirPath)

		name := "config.toml"
		err = gbfile.PutContents(gbfile.Join(dirPath, name), redisContent)
		t.AssertNil(err)

		err = gins.Config().GetAdapter().(*gbcfg.AdapterFile).AddPath(dirPath)
		t.AssertNil(err)

		defer gins.Config().GetAdapter().(*gbcfg.AdapterFile).Clear()

		// for gbfsnotify callbacks to refresh cache of config file
		time.Sleep(500 * time.Millisecond)

		// fmt.Println("gins Test_Redis", Config().Get("test"))

		var (
			redisDefault = gins.Redis()
			redisCache   = gins.Redis("cache")
			redisDisk    = gins.Redis("disk")
		)
		t.AssertNE(redisDefault, nil)
		t.AssertNE(redisCache, nil)
		t.AssertNE(redisDisk, nil)

		r, err := redisDefault.Do(ctx, "PING")
		t.AssertNil(err)
		t.Assert(r, "PONG")

		r, err = redisCache.Do(ctx, "PING")
		t.AssertNil(err)
		t.Assert(r, "PONG")

		_, err = redisDisk.Do(ctx, "SET", "k", "v")
		t.AssertNil(err)
		r, err = redisDisk.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(r, []byte("v"))
	})
}
