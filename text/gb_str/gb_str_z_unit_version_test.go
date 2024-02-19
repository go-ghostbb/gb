package gbstr_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"testing"
)

func Test_IsGNUVersion(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbstr.IsGNUVersion(""), false)
		t.AssertEQ(gbstr.IsGNUVersion("v"), false)
		t.AssertEQ(gbstr.IsGNUVersion("v0"), true)
		t.AssertEQ(gbstr.IsGNUVersion("v0."), false)
		t.AssertEQ(gbstr.IsGNUVersion("v1."), false)
		t.AssertEQ(gbstr.IsGNUVersion("v1.1"), true)
		t.AssertEQ(gbstr.IsGNUVersion("v1.1.0"), true)
		t.AssertEQ(gbstr.IsGNUVersion("v1.1."), false)
		t.AssertEQ(gbstr.IsGNUVersion("v1.1.0.0"), false)
		t.AssertEQ(gbstr.IsGNUVersion("v0.0.0"), true)
		t.AssertEQ(gbstr.IsGNUVersion("v1.1.-1"), false)
		t.AssertEQ(gbstr.IsGNUVersion("v1.1.+1"), false)
	})
}

func Test_CompareVersion(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbstr.CompareVersion("1", ""), 1)
		t.AssertEQ(gbstr.CompareVersion("", ""), 0)
		t.AssertEQ(gbstr.CompareVersion("", "v0.1"), -1)
		t.AssertEQ(gbstr.CompareVersion("1", "v0.99"), 1)
		t.AssertEQ(gbstr.CompareVersion("v1.0", "v0.99"), 1)
		t.AssertEQ(gbstr.CompareVersion("v1.0.1", "v1.1.0"), -1)
		t.AssertEQ(gbstr.CompareVersion("1.0.1", "v1.1.0"), -1)
		t.AssertEQ(gbstr.CompareVersion("1.0.0", "v0.1.0"), 1)
		t.AssertEQ(gbstr.CompareVersion("1.0.0", "v1.0.0"), 0)
	})
}

func Test_CompareVersionGo(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbstr.CompareVersionGo("1", ""), 1)
		t.AssertEQ(gbstr.CompareVersionGo("", ""), 0)
		t.AssertEQ(gbstr.CompareVersionGo("", "v0.1"), -1)
		t.AssertEQ(gbstr.CompareVersionGo("v1.0.1", "v1.1.0"), -1)
		t.AssertEQ(gbstr.CompareVersionGo("1.0.1", "v1.1.0"), -1)
		t.AssertEQ(gbstr.CompareVersionGo("1.0.0", "v0.1.0"), 1)
		t.AssertEQ(gbstr.CompareVersionGo("1.0.0", "v1.0.0"), 0)
		t.AssertEQ(gbstr.CompareVersionGo("1.0.0", "v1.0"), 0)
		t.AssertEQ(gbstr.CompareVersionGo("v0.0.0-20190626092158-b2ccc519800e", "0.0.0-20190626092158"), 0)
		t.AssertEQ(gbstr.CompareVersionGo("v0.0.0-20190626092159-b2ccc519800e", "0.0.0-20190626092158"), 1)

		// Specially in Golang:
		// "v1.12.2-0.20200413154443-b17e3a6804fa" < "v1.12.2"
		// "v1.12.3-0.20200413154443-b17e3a6804fa" > "v1.12.2"
		t.AssertEQ(gbstr.CompareVersionGo("v1.12.2-0.20200413154443-b17e3a6804fa", "v1.12.2"), -1)
		t.AssertEQ(gbstr.CompareVersionGo("v1.12.2", "v1.12.2-0.20200413154443-b17e3a6804fa"), 1)
		t.AssertEQ(gbstr.CompareVersionGo("v1.12.3-0.20200413154443-b17e3a6804fa", "v1.12.2"), 1)
		t.AssertEQ(gbstr.CompareVersionGo("v1.12.2", "v1.12.3-0.20200413154443-b17e3a6804fa"), -1)
		t.AssertEQ(gbstr.CompareVersionGo("v1.12.2-0.20200413154443-b17e3a6804fa", "v0.0.0-20190626092158-b2ccc519800e"), 1)
		t.AssertEQ(gbstr.CompareVersionGo("v1.12.2-0.20200413154443-b17e3a6804fa", "v1.12.2-0.20200413154444-b2ccc519800e"), -1)

		// Specially in Golang:
		// "v4.20.1+incompatible" < "v4.20.1"
		t.AssertEQ(gbstr.CompareVersionGo("v4.20.0+incompatible", "4.20.0"), -1)
		t.AssertEQ(gbstr.CompareVersionGo("4.20.0", "v4.20.0+incompatible"), 1)
		t.AssertEQ(gbstr.CompareVersionGo("v4.20.0+incompatible", "4.20.1"), -1)
		t.AssertEQ(gbstr.CompareVersionGo("v4.20.0+incompatible", "v4.20.0+incompatible"), 0)

	})
}
