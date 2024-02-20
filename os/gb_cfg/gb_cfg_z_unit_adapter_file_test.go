package gbcfg_test

import (
	gbcfg "ghostbb.io/gb/os/gb_cfg"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func TestAdapterFile_Dump(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		c, err := gbcfg.NewAdapterFile("config.yml")
		t.AssertNil(err)

		t.Assert(c.GetFileName(), "config.yml")

		c.Dump()
		c.Data(ctx)
	})

	gbtest.C(t, func(t *gbtest.T) {
		c, err := gbcfg.NewAdapterFile("testdata/default/config.toml")
		t.AssertNil(err)

		c.Dump()
		c.Data(ctx)
		c.GetPaths()
	})

}
func TestAdapterFile_Available(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		c, err := gbcfg.NewAdapterFile("testdata/default/config.toml")
		t.AssertNil(err)
		c.Available(ctx)
	})
}

func TestAdapterFile_SetPath(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		c, err := gbcfg.NewAdapterFile("config.yml")
		t.AssertNil(err)

		err = c.SetPath("/tmp")
		t.AssertNil(err)

		err = c.SetPath("notexist")
		t.AssertNE(err, nil)

		err = c.SetPath("testdata/c1.toml")
		t.AssertNE(err, nil)

		err = c.SetPath("")
		t.AssertNil(err)

		err = c.SetPath("gbcfg.go")
		t.AssertNE(err, nil)

		v, err := c.Get(ctx, "name")
		t.AssertNE(err, nil)
		t.Assert(v, nil)
	})
}

func TestAdapterFile_AddPath(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		c, err := gbcfg.NewAdapterFile("config.yml")
		t.AssertNil(err)

		err = c.AddPath("/tmp")
		t.AssertNil(err)

		err = c.AddPath("notexist")
		t.AssertNE(err, nil)

		err = c.SetPath("testdata/c1.toml")
		t.AssertNE(err, nil)

		err = c.SetPath("")
		t.AssertNil(err)

		err = c.AddPath("gbcfg.go")
		t.AssertNE(err, nil)

		v, err := c.Get(ctx, "name")
		t.AssertNE(err, nil)
		t.Assert(v, nil)
	})
}

func TestAdapterFile_SetViolenceCheck(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		c, err := gbcfg.NewAdapterFile("config.yml")
		t.AssertNil(err)
		c.SetViolenceCheck(true)
		v, err := c.Get(ctx, "name")
		t.AssertNE(err, nil)
		t.Assert(v, nil)
	})
}

func TestAdapterFile_FilePath(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		c, err := gbcfg.NewAdapterFile("config.yml")
		t.AssertNil(err)

		path, _ := c.GetFilePath("tmp")
		t.Assert(path, "")

		path, _ = c.GetFilePath("tmp")
		t.Assert(path, "")
	})
}

func TestAdapterFile_Content(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		c, err := gbcfg.NewAdapterFile()
		t.AssertNil(err)

		c.SetContent("gb", "config.yml")
		t.Assert(c.GetContent("config.yml"), "gb")
		c.SetContent("gb1", "config.yml")
		t.Assert(c.GetContent("config.yml"), "gb1")
		c.RemoveContent("config.yml")
		c.ClearContent()
		t.Assert(c.GetContent("name"), "")
	})
}

func TestAdapterFile_With_UTF8_BOM(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		c, err := gbcfg.NewAdapterFile("test-cfg-with-utf8-bom")
		t.AssertNil(err)

		t.Assert(c.SetPath("testdata"), nil)
		c.SetFileName("cfg-with-utf8-bom.toml")
		t.Assert(c.MustGet(ctx, "test.testInt"), 1)
		t.Assert(c.MustGet(ctx, "test.testStr"), "test")
	})
}

func TestAdapterFile_Set(t *testing.T) {
	config := `log-path = "logs"`
	gbtest.C(t, func(t *gbtest.T) {
		var (
			path = gbcfg.DefaultConfigFileName
			err  = gbfile.PutContents(path, config)
		)
		t.Assert(err, nil)
		defer gbfile.Remove(path)

		c, err := gbcfg.New()
		t.Assert(c.MustGet(ctx, "log-path").String(), "logs")

		err = c.GetAdapter().(*gbcfg.AdapterFile).Set("log-path", "custom-logs")
		t.Assert(err, nil)
		t.Assert(c.MustGet(ctx, "log-path").String(), "custom-logs")
	})
}
