package gbstr_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"testing"
)

func Test_IsSubDomain(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		main := "ghostbb.io"
		t.Assert(gbstr.IsSubDomain("ghostbb.io", main), true)
		t.Assert(gbstr.IsSubDomain("s.ghostbb.io", main), true)
		t.Assert(gbstr.IsSubDomain("s.s.ghostbb.io", main), true)
		t.Assert(gbstr.IsSubDomain("s.s.ghostbb.io:8080", main), true)
		t.Assert(gbstr.IsSubDomain("test.tw", main), false)
		t.Assert(gbstr.IsSubDomain("s.test.tw", main), false)
		t.Assert(gbstr.IsSubDomain("s.s.test.tw", main), false)
	})
	gbtest.C(t, func(t *gbtest.T) {
		main := "*.ghostbb.io"
		t.Assert(gbstr.IsSubDomain("ghostbb.io", main), true)
		t.Assert(gbstr.IsSubDomain("s.ghostbb.io", main), true)
		t.Assert(gbstr.IsSubDomain("s.ghostbb.io:80", main), true)
		t.Assert(gbstr.IsSubDomain("s.s.ghostbb.io", main), false)
		t.Assert(gbstr.IsSubDomain("test.tw", main), false)
		t.Assert(gbstr.IsSubDomain("s.test.tw", main), false)
		t.Assert(gbstr.IsSubDomain("s.s.test.tw", main), false)
	})
	gbtest.C(t, func(t *gbtest.T) {
		main := "*.*.ghostbb.io"
		t.Assert(gbstr.IsSubDomain("ghostbb.io", main), true)
		t.Assert(gbstr.IsSubDomain("s.ghostbb.io", main), true)
		t.Assert(gbstr.IsSubDomain("s.s.ghostbb.io", main), true)
		t.Assert(gbstr.IsSubDomain("s.s.ghostbb.io:8000", main), true)
		t.Assert(gbstr.IsSubDomain("s.s.s.ghostbb.io", main), false)
		t.Assert(gbstr.IsSubDomain("test.tw", main), false)
		t.Assert(gbstr.IsSubDomain("s.test.tw", main), false)
		t.Assert(gbstr.IsSubDomain("s.s.test.tw", main), false)
	})
	gbtest.C(t, func(t *gbtest.T) {
		main := "*.*.ghostbb.io:8080"
		t.Assert(gbstr.IsSubDomain("ghostbb.io", main), true)
		t.Assert(gbstr.IsSubDomain("s.ghostbb.io", main), true)
		t.Assert(gbstr.IsSubDomain("s.s.ghostbb.io", main), true)
		t.Assert(gbstr.IsSubDomain("s.s.ghostbb.io:8000", main), true)
		t.Assert(gbstr.IsSubDomain("s.s.s.ghostbb.io", main), false)
		t.Assert(gbstr.IsSubDomain("test.tw", main), false)
		t.Assert(gbstr.IsSubDomain("s.test.tw", main), false)
		t.Assert(gbstr.IsSubDomain("s.s.test.tw", main), false)
	})

	gbtest.C(t, func(t *gbtest.T) {
		main := "*.*.ghostbb.io:8080"
		t.Assert(gbstr.IsSubDomain("ghostbb.io", main), true)
		t.Assert(gbstr.IsSubDomain("s.ghostbb.io", main), true)
		t.Assert(gbstr.IsSubDomain("s.s.ghostbb.io", main), true)
		t.Assert(gbstr.IsSubDomain("s.s.ghostbb.io:8000", main), true)
		t.Assert(gbstr.IsSubDomain("s.s.s.ghostbb.io", main), false)
		t.Assert(gbstr.IsSubDomain("test.tw", main), false)
		t.Assert(gbstr.IsSubDomain("s.test.tw", main), false)
		t.Assert(gbstr.IsSubDomain("s.s.test.tw", main), false)
	})
	gbtest.C(t, func(t *gbtest.T) {
		main := "s.ghostbb.io"
		t.Assert(gbstr.IsSubDomain("ghostbb.io", main), false)
	})
}
