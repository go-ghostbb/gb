package gbcfg

import (
	"context"
	gbmap "ghostbb.io/gb/container/gb_map"
	gbenv "ghostbb.io/gb/os/gb_env"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

var (
	ctx = context.TODO()
)

func Test_Instance_Basic(t *testing.T) {
	config := `
array = [1.0, 2.0, 3.0]
v1 = 1.0
v2 = "true"
v3 = "off"
v4 = "1.234"

[redis]
  cache = "127.0.0.1:6379,1"
  disk = "127.0.0.1:6379,0"

`
	gbtest.C(t, func(t *gbtest.T) {
		var (
			path = DefaultConfigFileName
			err  = gbfile.PutContents(path, config)
		)
		t.AssertNil(err)
		defer func() {
			t.AssertNil(gbfile.Remove(path))
		}()

		c := Instance()
		t.Assert(c.MustGet(ctx, "v1"), 1)
		filepath, _ := c.GetAdapter().(*AdapterFile).GetFilePath()
		t.AssertEQ(filepath, gbfile.Pwd()+gbfile.Separator+path)
	})
}

func Test_Instance_AutoLocateConfigFile(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(Instance("gb") != nil, true)
	})
	// Automatically locate the configuration file with supported file extensions.
	gbtest.C(t, func(t *gbtest.T) {
		pwd := gbfile.Pwd()
		t.AssertNil(gbfile.Chdir(gbtest.DataPath()))
		defer gbfile.Chdir(pwd)
		t.Assert(Instance("c1") != nil, true)
		t.Assert(Instance("c1").MustGet(ctx, "my-config"), "1")
		t.Assert(Instance("folder1/c1").MustGet(ctx, "my-config"), "2")
	})
	// Automatically locate the configuration file with supported file extensions.
	gbtest.C(t, func(t *gbtest.T) {
		pwd := gbfile.Pwd()
		t.AssertNil(gbfile.Chdir(gbtest.DataPath("folder1")))
		defer gbfile.Chdir(pwd)
		t.Assert(Instance("c2").MustGet(ctx, "my-config"), 2)
	})
	// Default configuration file.
	gbtest.C(t, func(t *gbtest.T) {
		localInstances.Clear()
		pwd := gbfile.Pwd()
		t.AssertNil(gbfile.Chdir(gbtest.DataPath("default")))
		defer gbfile.Chdir(pwd)
		t.Assert(Instance().MustGet(ctx, "my-config"), 1)

		localInstances.Clear()
		t.AssertNil(gbenv.Set("GB_CFG_FILE", "config.json"))
		defer gbenv.Set("GB_CFG_FILE", "")
		t.Assert(Instance().MustGet(ctx, "my-config"), 2)
	})
}

func Test_Instance_EnvPath(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		gbenv.Set("GB_CFG_PATH", gbtest.DataPath("envpath"))
		defer gbenv.Set("GB_CFG_PATH", "")
		t.Assert(Instance("c3") != nil, true)
		t.Assert(Instance("c3").MustGet(ctx, "my-config"), "3")
		t.Assert(Instance("c4").MustGet(ctx, "my-config"), "4")
		localInstances = gbmap.NewStrAnyMap(true)
	})
}

func Test_Instance_EnvFile(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		gbenv.Set("GB_CFG_PATH", gbtest.DataPath("envfile"))
		defer gbenv.Set("GB_CFG_PATH", "")
		gbenv.Set("GB_CFG_FILE", "c6.json")
		defer gbenv.Set("GB_CFG_FILE", "")
		t.Assert(Instance().MustGet(ctx, "my-config"), "6")
		localInstances = gbmap.NewStrAnyMap(true)
	})
}
