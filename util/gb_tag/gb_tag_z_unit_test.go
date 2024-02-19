package gbtag_test

import (
	"fmt"
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbtag "ghostbb.io/gb/util/gb_tag"
	gbuid "ghostbb.io/gb/util/gb_uid"
	"testing"
)

func Test_Set_Get(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		k := gbuid.S()
		v := gbuid.S()
		gbtag.Set(k, v)
		t.Assert(gbtag.Get(k), v)
	})
}

func Test_SetOver_Get(t *testing.T) {
	// panic by Set
	gbtest.C(t, func(t *gbtest.T) {
		var (
			k  = gbuid.S()
			v1 = gbuid.S()
			v2 = gbuid.S()
		)
		gbtag.Set(k, v1)
		t.Assert(gbtag.Get(k), v1)
		defer func() {
			t.AssertNE(recover(), nil)
		}()
		gbtag.Set(k, v2)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var (
			k  = gbuid.S()
			v1 = gbuid.S()
			v2 = gbuid.S()
		)
		gbtag.SetOver(k, v1)
		t.Assert(gbtag.Get(k), v1)
		gbtag.SetOver(k, v2)
		t.Assert(gbtag.Get(k), v2)
	})
}

func Test_Sets_Get(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			k1 = gbuid.S()
			k2 = gbuid.S()
			v1 = gbuid.S()
			v2 = gbuid.S()
		)
		gbtag.Sets(g.MapStrStr{
			k1: v1,
			k2: v2,
		})
		t.Assert(gbtag.Get(k1), v1)
		t.Assert(gbtag.Get(k2), v2)
	})
}

func Test_SetsOver_Get(t *testing.T) {
	// panic by Sets
	gbtest.C(t, func(t *gbtest.T) {
		var (
			k1 = gbuid.S()
			k2 = gbuid.S()
			v1 = gbuid.S()
			v2 = gbuid.S()
			v3 = gbuid.S()
		)
		gbtag.Sets(g.MapStrStr{
			k1: v1,
			k2: v2,
		})
		t.Assert(gbtag.Get(k1), v1)
		t.Assert(gbtag.Get(k2), v2)
		defer func() {
			t.AssertNE(recover(), nil)
		}()
		gbtag.Sets(g.MapStrStr{
			k1: v3,
			k2: v3,
		})
	})
	gbtest.C(t, func(t *gbtest.T) {
		var (
			k1 = gbuid.S()
			k2 = gbuid.S()
			v1 = gbuid.S()
			v2 = gbuid.S()
			v3 = gbuid.S()
		)
		gbtag.SetsOver(g.MapStrStr{
			k1: v1,
			k2: v2,
		})
		t.Assert(gbtag.Get(k1), v1)
		t.Assert(gbtag.Get(k2), v2)
		gbtag.SetsOver(g.MapStrStr{
			k1: v3,
			k2: v3,
		})
		t.Assert(gbtag.Get(k1), v3)
		t.Assert(gbtag.Get(k2), v3)
	})
}

func Test_Parse(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			k1      = gbuid.S()
			k2      = gbuid.S()
			v1      = gbuid.S()
			v2      = gbuid.S()
			content = fmt.Sprintf(`this is {%s} and {%s}`, k1, k2)
			expect  = fmt.Sprintf(`this is %s and %s`, v1, v2)
		)
		gbtag.Sets(g.MapStrStr{
			k1: v1,
			k2: v2,
		})
		t.Assert(gbtag.Parse(content), expect)
	})
}

func Test_SetGlobalEnums(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		oldEnumsJson, err := gbtag.GetGlobalEnums()
		t.AssertNil(err)

		err = gbtag.SetGlobalEnums(`{"k8s.io/apimachinery/pkg/api/resource.Format": [
        "BinarySI",
        "DecimalExponent",
        "DecimalSI"
    ]}`)
		t.AssertNil(err)
		t.Assert(gbtag.GetEnumsByType("k8s.io/apimachinery/pkg/api/resource.Format"), `[
        "BinarySI",
        "DecimalExponent",
        "DecimalSI"
    ]`)
		t.AssertNil(gbtag.SetGlobalEnums(oldEnumsJson))
	})
}
