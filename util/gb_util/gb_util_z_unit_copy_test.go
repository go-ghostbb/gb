package gbutil_test

import (
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbutil "ghostbb.io/gb/util/gb_util"
	"testing"
)

func Test_Copy(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.Copy(0), 0)
		t.Assert(gbutil.Copy(1), 1)
		t.Assert(gbutil.Copy("a"), "a")
		t.Assert(gbutil.Copy(nil), nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		src := g.Map{
			"k1": "v1",
			"k2": "v2",
		}
		dst := gbutil.Copy(src)
		t.Assert(dst, src)

		dst.(g.Map)["k3"] = "v3"
		t.Assert(src, g.Map{
			"k1": "v1",
			"k2": "v2",
		})
		t.Assert(dst, g.Map{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		})
	})
}
