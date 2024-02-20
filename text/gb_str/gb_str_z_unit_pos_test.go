package gbstr_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"testing"
)

func Test_Pos(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "abcdEFGabcdefg"
		t.Assert(gbstr.Pos(s1, "ab"), 0)
		t.Assert(gbstr.Pos(s1, "ab", 2), 7)
		t.Assert(gbstr.Pos(s1, "abd", 0), -1)
		t.Assert(gbstr.Pos(s1, "e", -4), 11)
	})
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "我愛Taiwan very much"
		t.Assert(gbstr.Pos(s1, "愛"), 3)
		t.Assert(gbstr.Pos(s1, "T"), 6)
		t.Assert(gbstr.Pos(s1, "Taiwan"), 6)
	})
}

func Test_PosRune(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "abcdEFGabcdefg"
		t.Assert(gbstr.PosRune(s1, "ab"), 0)
		t.Assert(gbstr.PosRune(s1, "ab", 2), 7)
		t.Assert(gbstr.PosRune(s1, "abd", 0), -1)
		t.Assert(gbstr.PosRune(s1, "e", -4), 11)
	})
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "我愛Taiwan very much"
		t.Assert(gbstr.PosRune(s1, "愛"), 1)
		t.Assert(gbstr.PosRune(s1, "T"), 2)
		t.Assert(gbstr.PosRune(s1, "Taiwan"), 2)
	})
}

func Test_PosI(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "abcdEFGabcdefg"
		t.Assert(gbstr.PosI(s1, "zz"), -1)
		t.Assert(gbstr.PosI(s1, "ab"), 0)
		t.Assert(gbstr.PosI(s1, "ef", 2), 4)
		t.Assert(gbstr.PosI(s1, "abd", 0), -1)
		t.Assert(gbstr.PosI(s1, "E", -4), 11)
	})
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "我愛Taiwan very much"
		t.Assert(gbstr.PosI(s1, "愛"), 3)
		t.Assert(gbstr.PosI(s1, "t"), 6)
		t.Assert(gbstr.PosI(s1, "taiwan"), 6)
	})
}

func Test_PosIRune(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "abcdEFGabcdefg"
		t.Assert(gbstr.PosIRune(s1, "zz"), -1)
		t.Assert(gbstr.PosIRune(s1, "ab"), 0)
		t.Assert(gbstr.PosIRune(s1, "ef", 2), 4)
		t.Assert(gbstr.PosIRune(s1, "abd", 0), -1)
		t.Assert(gbstr.PosIRune(s1, "E", -4), 11)
	})
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "我愛Taiwan very much"
		t.Assert(gbstr.PosIRune(s1, "愛"), 1)
		t.Assert(gbstr.PosIRune(s1, "t"), 2)
		t.Assert(gbstr.PosIRune(s1, "taiwan"), 2)
	})
}

func Test_PosR(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		t.Assert(gbstr.PosR(s1, "zz"), -1)
		t.Assert(gbstr.PosR(s1, "ab"), 7)
		t.Assert(gbstr.PosR(s2, "ab", -2), 0)
		t.Assert(gbstr.PosR(s1, "ef"), 11)
		t.Assert(gbstr.PosR(s1, "abd", 0), -1)
		t.Assert(gbstr.PosR(s1, "e", -4), -1)
	})
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "我愛Taiwan very much"
		t.Assert(gbstr.PosR(s1, "愛"), 3)
		t.Assert(gbstr.PosR(s1, "T"), 6)
		t.Assert(gbstr.PosR(s1, "Taiwan"), 6)
	})
}

func Test_PosRRune(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		t.Assert(gbstr.PosRRune(s1, "zz"), -1)
		t.Assert(gbstr.PosRRune(s1, "ab"), 7)
		t.Assert(gbstr.PosRRune(s2, "ab", -2), 0)
		t.Assert(gbstr.PosRRune(s1, "ef"), 11)
		t.Assert(gbstr.PosRRune(s1, "abd", 0), -1)
		t.Assert(gbstr.PosRRune(s1, "e", -4), -1)
	})
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "我愛Taiwan very much"
		t.Assert(gbstr.PosRRune(s1, "愛"), 1)
		t.Assert(gbstr.PosRRune(s1, "T"), 2)
		t.Assert(gbstr.PosRRune(s1, "Taiwan"), 2)
	})
}

func Test_PosRI(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		t.Assert(gbstr.PosRI(s1, "zz"), -1)
		t.Assert(gbstr.PosRI(s1, "AB"), 7)
		t.Assert(gbstr.PosRI(s2, "AB", -2), 0)
		t.Assert(gbstr.PosRI(s1, "EF"), 11)
		t.Assert(gbstr.PosRI(s1, "abd", 0), -1)
		t.Assert(gbstr.PosRI(s1, "e", -5), 4)
	})
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "我愛Taiwan very much"
		t.Assert(gbstr.PosRI(s1, "愛"), 3)
		t.Assert(gbstr.PosRI(s1, "C"), 20)
		t.Assert(gbstr.PosRI(s1, "Taiwan"), 6)
	})
}

func Test_PosRIRune(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		t.Assert(gbstr.PosRIRune(s1, "zz"), -1)
		t.Assert(gbstr.PosRIRune(s1, "AB"), 7)
		t.Assert(gbstr.PosRIRune(s2, "AB", -2), 0)
		t.Assert(gbstr.PosRIRune(s1, "EF"), 11)
		t.Assert(gbstr.PosRIRune(s1, "abd", 0), -1)
		t.Assert(gbstr.PosRIRune(s1, "e", -5), 4)
	})
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "我愛Taiwan very much"
		t.Assert(gbstr.PosRIRune(s1, "愛"), 1)
		t.Assert(gbstr.PosRIRune(s1, "C"), 16)
		t.Assert(gbstr.PosRIRune(s1, "Taiwan"), 2)
	})
}
