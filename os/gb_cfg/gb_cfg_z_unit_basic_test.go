package gbcfg_test

import (
	gbcfg "ghostbb.io/gb/os/gb_cfg"
	gbcmd "ghostbb.io/gb/os/gb_cmd"
	gbenv "ghostbb.io/gb/os/gb_env"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func Test_Basic1(t *testing.T) {
	config := `
v1    = 1
v2    = "true"
v3    = "off"
v4    = "1.23"
array = [1,2,3]
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`
	gbtest.C(t, func(t *gbtest.T) {
		var (
			path = gbcfg.DefaultConfigFileName
			err  = gbfile.PutContents(path, config)
		)
		t.AssertNil(err)
		defer gbfile.Remove(path)

		c, err := gbcfg.New()
		t.AssertNil(err)
		t.Assert(c.MustGet(ctx, "v1"), 1)
		filepath, _ := c.GetAdapter().(*gbcfg.AdapterFile).GetFilePath()
		t.AssertEQ(filepath, gbfile.Pwd()+gbfile.Separator+path)
	})
}

func Test_Basic2(t *testing.T) {
	config := `log-path = "logs"`
	gbtest.C(t, func(t *gbtest.T) {
		var (
			path = gbcfg.DefaultConfigFileName
			err  = gbfile.PutContents(path, config)
		)
		t.AssertNil(err)
		defer func() {
			_ = gbfile.Remove(path)
		}()

		c, err := gbcfg.New()
		t.AssertNil(err)
		t.Assert(c.MustGet(ctx, "log-path"), "logs")
	})
}

func Test_Content(t *testing.T) {
	content := `
v1    = 1
v2    = "true"
v3    = "off"
v4    = "1.23"
array = [1,2,3]
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`
	gbtest.C(t, func(t *gbtest.T) {
		c, err := gbcfg.New()
		t.AssertNil(err)
		c.GetAdapter().(*gbcfg.AdapterFile).SetContent(content)
		defer c.GetAdapter().(*gbcfg.AdapterFile).ClearContent()
		t.Assert(c.MustGet(ctx, "v1"), 1)
	})
}

func Test_SetFileName(t *testing.T) {
	config := `
{
	"array": [
		1,
		2,
		3
	],
	"redis": {
		"cache": "127.0.0.1:6379,1",
		"disk": "127.0.0.1:6379,0"
	},
	"v1": 1,
	"v2": "true",
	"v3": "off",
	"v4": "1.234"
}
`
	gbtest.C(t, func(t *gbtest.T) {
		path := "config.json"
		err := gbfile.PutContents(path, config)
		t.AssertNil(err)
		defer func() {
			_ = gbfile.Remove(path)
		}()

		config, err := gbcfg.New()
		t.AssertNil(err)
		c := config.GetAdapter().(*gbcfg.AdapterFile)
		c.SetFileName(path)
		t.Assert(c.MustGet(ctx, "v1"), 1)
		t.AssertEQ(c.MustGet(ctx, "v1").Int(), 1)
		t.AssertEQ(c.MustGet(ctx, "v1").Int8(), int8(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Int16(), int16(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Int32(), int32(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Int64(), int64(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Uint(), uint(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Uint8(), uint8(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Uint16(), uint16(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Uint32(), uint32(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Uint64(), uint64(1))

		t.AssertEQ(c.MustGet(ctx, "v1").String(), "1")
		t.AssertEQ(c.MustGet(ctx, "v1").Bool(), true)
		t.AssertEQ(c.MustGet(ctx, "v2").String(), "true")
		t.AssertEQ(c.MustGet(ctx, "v2").Bool(), true)

		t.AssertEQ(c.MustGet(ctx, "v1").String(), "1")
		t.AssertEQ(c.MustGet(ctx, "v4").Float32(), float32(1.234))
		t.AssertEQ(c.MustGet(ctx, "v4").Float64(), float64(1.234))
		t.AssertEQ(c.MustGet(ctx, "v2").String(), "true")
		t.AssertEQ(c.MustGet(ctx, "v2").Bool(), true)
		t.AssertEQ(c.MustGet(ctx, "v3").Bool(), false)

		t.AssertEQ(c.MustGet(ctx, "array").Ints(), []int{1, 2, 3})
		t.AssertEQ(c.MustGet(ctx, "array").Strings(), []string{"1", "2", "3"})
		t.AssertEQ(c.MustGet(ctx, "array").Interfaces(), []interface{}{1, 2, 3})
		t.AssertEQ(c.MustGet(ctx, "redis").Map(), map[string]interface{}{
			"disk":  "127.0.0.1:6379,0",
			"cache": "127.0.0.1:6379,1",
		})
		filepath, _ := c.GetFilePath()
		t.AssertEQ(filepath, gbfile.Pwd()+gbfile.Separator+path)
	})
}

func TestCfg_Get_WrongConfigFile(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var err error
		configPath := gbfile.Temp(gbtime.TimestampNanoStr())
		err = gbfile.Mkdir(configPath)
		t.AssertNil(err)
		defer gbfile.Remove(configPath)

		defer gbfile.Chdir(gbfile.Pwd())
		err = gbfile.Chdir(configPath)
		t.AssertNil(err)

		err = gbfile.PutContents(
			gbfile.Join(configPath, "config.yml"),
			"wrong config",
		)
		t.AssertNil(err)
		adapterFile, err := gbcfg.NewAdapterFile("config.yml")
		t.AssertNil(err)

		c := gbcfg.NewWithAdapter(adapterFile)
		v, err := c.Get(ctx, "name")
		t.AssertNE(err, nil)
		t.Assert(v, nil)
		adapterFile.Clear()
	})
}

func Test_GetWithEnv(t *testing.T) {
	content := `
v1    = 1
v2    = "true"
v3    = "off"
v4    = "1.23"
array = [1,2,3]
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`
	gbtest.C(t, func(t *gbtest.T) {
		c, err := gbcfg.New()
		t.AssertNil(err)
		c.GetAdapter().(*gbcfg.AdapterFile).SetContent(content)
		defer c.GetAdapter().(*gbcfg.AdapterFile).ClearContent()
		t.Assert(c.MustGet(ctx, "v1"), 1)
		t.Assert(c.MustGetWithEnv(ctx, `redis.user`), nil)
		t.Assert(gbenv.Set("REDIS_USER", `1`), nil)
		defer gbenv.Remove(`REDIS_USER`)
		t.Assert(c.MustGetWithEnv(ctx, `redis.user`), `1`)
	})
}

func Test_GetWithCmd(t *testing.T) {
	content := `
v1    = 1
v2    = "true"
v3    = "off"
v4    = "1.23"
array = [1,2,3]
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`
	gbtest.C(t, func(t *gbtest.T) {

		c, err := gbcfg.New()
		t.AssertNil(err)
		c.GetAdapter().(*gbcfg.AdapterFile).SetContent(content)
		defer c.GetAdapter().(*gbcfg.AdapterFile).ClearContent()
		t.Assert(c.MustGet(ctx, "v1"), 1)
		t.Assert(c.MustGetWithCmd(ctx, `redis.user`), nil)

		gbcmd.Init([]string{"gb", "--redis.user=2"}...)
		t.Assert(c.MustGetWithCmd(ctx, `redis.user`), `2`)
	})
}
