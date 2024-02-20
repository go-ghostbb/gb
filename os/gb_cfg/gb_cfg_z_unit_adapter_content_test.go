package gbcfg_test

import (
	"ghostbb.io/gb/frame/g"
	gbcfg "ghostbb.io/gb/os/gb_cfg"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func TestAdapterContent_Available_Get_Data(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		adapter, err := gbcfg.NewAdapterContent()
		t.AssertNil(err)
		t.Assert(adapter.Available(ctx), false)
	})
	gbtest.C(t, func(t *gbtest.T) {
		content := `{"a": 1, "b": 2, "c": {"d": 3}}`
		adapter, err := gbcfg.NewAdapterContent(content)
		t.AssertNil(err)

		c := gbcfg.NewWithAdapter(adapter)
		t.Assert(c.Available(ctx), true)
		t.Assert(c.MustGet(ctx, "a"), 1)
		t.Assert(c.MustGet(ctx, "b"), 2)
		t.Assert(c.MustGet(ctx, "c.d"), 3)
		t.Assert(c.MustGet(ctx, "d"), nil)
		t.Assert(c.MustData(ctx), g.Map{
			"a": 1,
			"b": 2,
			"c": g.Map{
				"d": 3,
			},
		})
	})
}

func TestAdapterContent_SetContent(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		adapter, err := gbcfg.NewAdapterContent()
		t.AssertNil(err)
		t.Assert(adapter.Available(ctx), false)

		content := `{"a": 1, "b": 2, "c": {"d": 3}}`
		err = adapter.SetContent(content)
		t.AssertNil(err)
		c := gbcfg.NewWithAdapter(adapter)
		t.Assert(c.Available(ctx), true)
		t.Assert(c.MustGet(ctx, "a"), 1)
		t.Assert(c.MustGet(ctx, "b"), 2)
		t.Assert(c.MustGet(ctx, "c.d"), 3)
		t.Assert(c.MustGet(ctx, "d"), nil)
		t.Assert(c.MustData(ctx), g.Map{
			"a": 1,
			"b": 2,
			"c": g.Map{
				"d": 3,
			},
		})
	})

}
