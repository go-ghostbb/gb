package gbstr_test

import (
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"testing"
)

func Test_SearchArray(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		a := g.SliceStr{"a", "b", "c"}
		t.AssertEQ(gbstr.SearchArray(a, "a"), 0)
		t.AssertEQ(gbstr.SearchArray(a, "b"), 1)
		t.AssertEQ(gbstr.SearchArray(a, "c"), 2)
		t.AssertEQ(gbstr.SearchArray(a, "d"), -1)
	})
}

func Test_InArray(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		a := g.SliceStr{"a", "b", "c"}
		t.AssertEQ(gbstr.InArray(a, "a"), true)
		t.AssertEQ(gbstr.InArray(a, "b"), true)
		t.AssertEQ(gbstr.InArray(a, "c"), true)
		t.AssertEQ(gbstr.InArray(a, "d"), false)
	})
}

func Test_PrefixArray(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		a := g.SliceStr{"a", "b", "c"}
		gbstr.PrefixArray(a, "1-")
		t.AssertEQ(a, g.SliceStr{"1-a", "1-b", "1-c"})
	})
}
