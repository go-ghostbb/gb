package gbstr_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"testing"
)

func Test_OctStr(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.OctStr(`\346\200\241`), "怡")
	})
}

func Test_WordWrap(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.WordWrap("12 34", 2, "<br>"), "12<br>34")
		t.Assert(gbstr.WordWrap("12 34", 2, "\n"), "12\n34")
		t.Assert(gbstr.WordWrap("我愛 GF", 2, "\n"), "我愛\nGF")
		t.Assert(gbstr.WordWrap("A very long woooooooooooooooooord. and something", 7, "<br>"),
			"A very<br>long<br>woooooooooooooooooord.<br>and<br>something")
	})
	// Chinese Punctuations.
	gbtest.C(t, func(t *gbtest.T) {
		var (
			br      = "                       "
			content = "    DelRouteKeyIPv6    刪除VPC内的服務的Route信息;和DelRouteIPv6接口相比，這個接口可以刪除滿足條件的多條RS\n"
			length  = 120
		)
		wrappedContent := gbstr.WordWrap(content, length, "\n"+br)
		t.Assert(wrappedContent, `    DelRouteKeyIPv6    刪除VPC内的服務的Route信息;和DelRouteIPv6接口相比，
                       這個接口可以刪除滿足條件的多條RS
`)
	})
}
