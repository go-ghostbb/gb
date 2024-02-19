package gbutil_test

import (
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbutil "ghostbb.io/gb/util/gb_util"
	"testing"
)

func Test_SliceCopy(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := g.Slice{
			"K1", "v1", "K2", "v2",
		}
		s1 := gbutil.SliceCopy(s)
		t.Assert(s, s1)
	})
}

func Test_SliceDelete(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := g.Slice{
			"K1", "v1", "K2", "v2",
		}
		t.Assert(gbutil.SliceDelete(s, 0), g.Slice{
			"v1", "K2", "v2",
		})
		t.Assert(gbutil.SliceDelete(s, 5), s)
	})
}

func Test_SliceToMap(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := g.Slice{
			"K1", "v1", "K2", "v2",
		}
		m := gbutil.SliceToMap(s)
		t.Assert(len(m), 2)
		t.Assert(m, g.Map{
			"K1": "v1",
			"K2": "v2",
		})

		m1 := gbutil.SliceToMap(&s)
		t.Assert(len(m1), 2)
		t.Assert(m1, g.Map{
			"K1": "v1",
			"K2": "v2",
		})
	})
	gbtest.C(t, func(t *gbtest.T) {
		s := g.Slice{
			"K1", "v1", "K2",
		}
		m := gbutil.SliceToMap(s)
		t.Assert(len(m), 0)
		t.Assert(m, nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		m := gbutil.SliceToMap(1)
		t.Assert(len(m), 0)
		t.Assert(m, nil)
	})
}

func Test_SliceToMapWithColumnAsKey(t *testing.T) {
	m1 := g.Map{"K1": "v1", "K2": 1}
	m2 := g.Map{"K1": "v2", "K2": 2}
	s := g.Slice{m1, m2}
	gbtest.C(t, func(t *gbtest.T) {
		m := gbutil.SliceToMapWithColumnAsKey(s, "K1")
		t.Assert(m, g.MapAnyAny{
			"v1": m1,
			"v2": m2,
		})

		n := gbutil.SliceToMapWithColumnAsKey(&s, "K1")
		t.Assert(n, g.MapAnyAny{
			"v1": m1,
			"v2": m2,
		})
	})
	gbtest.C(t, func(t *gbtest.T) {
		m := gbutil.SliceToMapWithColumnAsKey(s, "K2")
		t.Assert(m, g.MapAnyAny{
			1: m1,
			2: m2,
		})
	})
}
