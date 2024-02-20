package gbstr_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"testing"
)

func Test_List2(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2 := gbstr.List2("1:2", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2 := gbstr.List2("1:", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2 := gbstr.List2("1", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2 := gbstr.List2("", ":")
		t.Assert(p1, "")
		t.Assert(p2, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2 := gbstr.List2("1:2:3", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2:3")
	})
}

func Test_ListAndTrim2(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2 := gbstr.ListAndTrim2("1::2", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2 := gbstr.ListAndTrim2("1::", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2 := gbstr.ListAndTrim2("1:", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2 := gbstr.ListAndTrim2("", ":")
		t.Assert(p1, "")
		t.Assert(p2, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2 := gbstr.ListAndTrim2("1::2::3", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2:3")
	})
}

func Test_List3(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2, p3 := gbstr.List3("1:2:3", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
		t.Assert(p3, "3")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2, p3 := gbstr.List3("1:2:", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
		t.Assert(p3, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2, p3 := gbstr.List3("1:2", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
		t.Assert(p3, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2, p3 := gbstr.List3("1:", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
		t.Assert(p3, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2, p3 := gbstr.List3("1", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
		t.Assert(p3, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2, p3 := gbstr.List3("", ":")
		t.Assert(p1, "")
		t.Assert(p2, "")
		t.Assert(p3, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2, p3 := gbstr.List3("1:2:3:4", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
		t.Assert(p3, "3:4")
	})
}

func Test_ListAndTrim3(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2, p3 := gbstr.ListAndTrim3("1::2:3", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
		t.Assert(p3, "3")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2, p3 := gbstr.ListAndTrim3("1::2:", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
		t.Assert(p3, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2, p3 := gbstr.ListAndTrim3("1::2", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
		t.Assert(p3, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2, p3 := gbstr.ListAndTrim3("1::", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
		t.Assert(p3, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2, p3 := gbstr.ListAndTrim3("1::", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
		t.Assert(p3, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		p1, p2, p3 := gbstr.ListAndTrim3("", ":")
		t.Assert(p1, "")
		t.Assert(p2, "")
		t.Assert(p3, "")
	})
}
