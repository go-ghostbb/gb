package gbstr_test

import (
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"testing"
)

func Test_Replace(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "abcdEFG亂入的中文abcdefg"
		t.Assert(gbstr.Replace(s1, "ab", "AB"), "ABcdEFG亂入的中文ABcdefg")
		t.Assert(gbstr.Replace(s1, "EF", "ef"), "abcdefG亂入的中文abcdefg")
		t.Assert(gbstr.Replace(s1, "MN", "mn"), s1)

		t.Assert(gbstr.ReplaceByArray(s1, g.ArrayStr{
			"a", "A",
			"A", "-",
			"a",
		}), "-bcdEFG亂入的中文-bcdefg")

		t.Assert(gbstr.ReplaceByMap(s1, g.MapStrStr{
			"a": "A",
			"G": "g",
		}), "AbcdEFg亂入的中文Abcdefg")
	})
}

func Test_ReplaceI_1(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "abcd亂入的中文ABCD"
		s2 := "a"
		t.Assert(gbstr.ReplaceI(s1, "ab", "aa"), "aacd亂入的中文aaCD")
		t.Assert(gbstr.ReplaceI(s1, "ab", "aa", 0), "abcd亂入的中文ABCD")
		t.Assert(gbstr.ReplaceI(s1, "ab", "aa", 1), "aacd亂入的中文ABCD")

		t.Assert(gbstr.ReplaceI(s1, "abcd", "-"), "-亂入的中文-")
		t.Assert(gbstr.ReplaceI(s1, "abcd", "-", 1), "-亂入的中文ABCD")

		t.Assert(gbstr.ReplaceI(s1, "abcd亂入的", ""), "中文ABCD")
		t.Assert(gbstr.ReplaceI(s1, "ABCD亂入的", ""), "中文ABCD")

		t.Assert(gbstr.ReplaceI(s2, "A", "-"), "-")
		t.Assert(gbstr.ReplaceI(s2, "a", "-"), "-")

		t.Assert(gbstr.ReplaceIByArray(s1, g.ArrayStr{
			"abcd亂入的", "-",
			"-", "=",
			"a",
		}), "=中文ABCD")

		t.Assert(gbstr.ReplaceIByMap(s1, g.MapStrStr{
			"ab": "-",
			"CD": "=",
		}), "-=亂入的中文-=")
	})
}

func Test_ReplaceI_2(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.ReplaceI("aaa", "A", "-a-"), `-a--a--a-`)
		t.Assert(gbstr.ReplaceI("aaaa", "AA", "-"), `--`)
		t.Assert(gbstr.ReplaceI("a a a", "A", "b"), `b b b`)
		t.Assert(gbstr.ReplaceI("aaaaaa", "aa", "a"), `aaa`)
		t.Assert(gbstr.ReplaceI("aaaaaa", "AA", "A"), `AAA`)
		t.Assert(gbstr.ReplaceI("aaa", "A", "AA"), `AAAAAA`)
		t.Assert(gbstr.ReplaceI("aaa", "A", "AA"), `AAAAAA`)
		t.Assert(gbstr.ReplaceI("a duration", "duration", "recordduration"), `a recordduration`)
	})
	// With count parameter.
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.ReplaceI("aaaaaa", "aa", "a", 2), `aaaa`)
		t.Assert(gbstr.ReplaceI("aaaaaa", "AA", "A", 1), `Aaaaa`)
		t.Assert(gbstr.ReplaceI("aaaaaa", "AA", "A", 3), `AAA`)
		t.Assert(gbstr.ReplaceI("aaaaaa", "AA", "A", 4), `AAA`)
		t.Assert(gbstr.ReplaceI("aaa", "A", "AA", 2), `AAAAa`)
		t.Assert(gbstr.ReplaceI("aaa", "A", "AA", 3), `AAAAAA`)
		t.Assert(gbstr.ReplaceI("aaa", "A", "AA", 4), `AAAAAA`)
	})
}
