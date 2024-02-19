package gbutil_test

import (
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbutil "ghostbb.io/gb/util/gb_util"
	"testing"
)

func Test_MapCopy(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		m1 := g.Map{
			"k1": "v1",
		}
		m2 := gbutil.MapCopy(m1)
		m2["k2"] = "v2"

		t.Assert(m1["k1"], "v1")
		t.Assert(m1["k2"], nil)
		t.Assert(m2["k1"], "v1")
		t.Assert(m2["k2"], "v2")
	})
}

func Test_MapContains(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		m1 := g.Map{
			"k1": "v1",
		}
		t.Assert(gbutil.MapContains(m1, "k1"), true)
		t.Assert(gbutil.MapContains(m1, "K1"), false)
		t.Assert(gbutil.MapContains(m1, "k2"), false)
		m2 := g.Map{}
		t.Assert(gbutil.MapContains(m2, "k1"), false)
	})
}

func Test_MapDelete(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		m1 := g.Map{
			"k1": "v1",
		}
		gbutil.MapDelete(m1, "k1")
		gbutil.MapDelete(m1, "K1")
		m2 := g.Map{}
		gbutil.MapDelete(m2, "k1")
	})
}

func Test_MapMerge(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		m1 := g.Map{
			"k1": "v1",
		}
		m2 := g.Map{
			"k2": "v2",
		}
		m3 := g.Map{
			"k3": "v3",
		}
		gbutil.MapMerge(m1, m2, m3, nil)
		t.Assert(m1["k1"], "v1")
		t.Assert(m1["k2"], "v2")
		t.Assert(m1["k3"], "v3")
		t.Assert(m2["k1"], nil)
		t.Assert(m3["k1"], nil)
		gbutil.MapMerge(nil)
	})
}

func Test_MapMergeCopy(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		m1 := g.Map{
			"k1": "v1",
		}
		m2 := g.Map{
			"k2": "v2",
		}
		m3 := g.Map{
			"k3": "v3",
		}
		m := gbutil.MapMergeCopy(m1, m2, m3, nil)
		t.Assert(m["k1"], "v1")
		t.Assert(m["k2"], "v2")
		t.Assert(m["k3"], "v3")
		t.Assert(m1["k1"], "v1")
		t.Assert(m1["k2"], nil)
		t.Assert(m2["k1"], nil)
		t.Assert(m3["k1"], nil)
	})
}

func Test_MapPossibleItemByKey(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		m := g.Map{
			"name":     "guo",
			"NickName": "john",
		}
		k, v := gbutil.MapPossibleItemByKey(m, "NAME")
		t.Assert(k, "name")
		t.Assert(v, "guo")

		k, v = gbutil.MapPossibleItemByKey(m, "nick name")
		t.Assert(k, "NickName")
		t.Assert(v, "john")

		k, v = gbutil.MapPossibleItemByKey(m, "none")
		t.Assert(k, "")
		t.Assert(v, nil)
	})
}

func Test_MapContainsPossibleKey(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		m := g.Map{
			"name":     "guo",
			"NickName": "john",
		}
		t.Assert(gbutil.MapContainsPossibleKey(m, "name"), true)
		t.Assert(gbutil.MapContainsPossibleKey(m, "NAME"), true)
		t.Assert(gbutil.MapContainsPossibleKey(m, "nickname"), true)
		t.Assert(gbutil.MapContainsPossibleKey(m, "nick name"), true)
		t.Assert(gbutil.MapContainsPossibleKey(m, "nick_name"), true)
		t.Assert(gbutil.MapContainsPossibleKey(m, "nick-name"), true)
		t.Assert(gbutil.MapContainsPossibleKey(m, "nick.name"), true)
		t.Assert(gbutil.MapContainsPossibleKey(m, "none"), false)
	})
}

func Test_MapOmitEmpty(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		m := g.Map{
			"k1": "john",
			"e1": "",
			"e2": 0,
			"e3": nil,
			"k2": "smith",
		}
		gbutil.MapOmitEmpty(m)
		t.Assert(len(m), 2)
		t.AssertNE(m["k1"], nil)
		t.AssertNE(m["k2"], nil)
		m1 := g.Map{}
		gbutil.MapOmitEmpty(m1)
		t.Assert(len(m1), 0)
	})
}

func Test_MapToSlice(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		m := g.Map{
			"k1": "v1",
			"k2": "v2",
		}
		s := gbutil.MapToSlice(m)
		t.Assert(len(s), 4)
		t.AssertIN(s[0], g.Slice{"k1", "k2", "v1", "v2"})
		t.AssertIN(s[1], g.Slice{"k1", "k2", "v1", "v2"})
		t.AssertIN(s[2], g.Slice{"k1", "k2", "v1", "v2"})
		t.AssertIN(s[3], g.Slice{"k1", "k2", "v1", "v2"})
		s1 := gbutil.MapToSlice(&m)
		t.Assert(len(s1), 4)
		t.AssertIN(s1[0], g.Slice{"k1", "k2", "v1", "v2"})
		t.AssertIN(s1[1], g.Slice{"k1", "k2", "v1", "v2"})
		t.AssertIN(s1[2], g.Slice{"k1", "k2", "v1", "v2"})
		t.AssertIN(s1[3], g.Slice{"k1", "k2", "v1", "v2"})
	})
	gbtest.C(t, func(t *gbtest.T) {
		m := g.MapStrStr{
			"k1": "v1",
			"k2": "v2",
		}
		s := gbutil.MapToSlice(m)
		t.Assert(len(s), 4)
		t.AssertIN(s[0], g.Slice{"k1", "k2", "v1", "v2"})
		t.AssertIN(s[1], g.Slice{"k1", "k2", "v1", "v2"})
		t.AssertIN(s[2], g.Slice{"k1", "k2", "v1", "v2"})
		t.AssertIN(s[3], g.Slice{"k1", "k2", "v1", "v2"})
	})
	gbtest.C(t, func(t *gbtest.T) {
		m := g.MapStrStr{}
		s := gbutil.MapToSlice(m)
		t.Assert(len(s), 0)
	})
	gbtest.C(t, func(t *gbtest.T) {
		s := gbutil.MapToSlice(1)
		t.Assert(s, nil)
	})
}
