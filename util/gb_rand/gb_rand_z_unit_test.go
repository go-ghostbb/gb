package gbrand_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbrand "ghostbb.io/gb/util/gb_rand"
	"strings"
	"testing"
	"time"
)

func Test_Intn(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 1000000; i++ {
			n := gbrand.Intn(100)
			t.AssertLT(n, 100)
			t.AssertGE(n, 0)
		}
		for i := 0; i < 1000000; i++ {
			n := gbrand.Intn(-100)
			t.AssertLE(n, 0)
			t.Assert(n, -100)
		}
	})
}

func Test_Meet(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(gbrand.Meet(100, 100), true)
		}
		for i := 0; i < 100; i++ {
			t.Assert(gbrand.Meet(0, 100), false)
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(gbrand.Meet(50, 100), []bool{true, false})
		}
	})
}

func Test_MeetProb(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(gbrand.MeetProb(1), true)
		}
		for i := 0; i < 100; i++ {
			t.Assert(gbrand.MeetProb(0), false)
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(gbrand.MeetProb(0.5), []bool{true, false})
		}
	})
}

func Test_N(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(gbrand.N(1, 1), 1)
		}
		for i := 0; i < 100; i++ {
			t.Assert(gbrand.N(0, 0), 0)
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(gbrand.N(1, 2), []int{1, 2})
		}
	})
}

func Test_D(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(gbrand.D(time.Second, time.Second), time.Second)
		}
		for i := 0; i < 100; i++ {
			t.Assert(gbrand.D(0, 0), time.Duration(0))
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(
				gbrand.D(1*time.Second, 3*time.Second),
				[]time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second},
			)
		}
	})
}

func Test_Rand(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(gbrand.N(1, 1), 1)
		}
		for i := 0; i < 100; i++ {
			t.Assert(gbrand.N(0, 0), 0)
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(gbrand.N(1, 2), []int{1, 2})
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(gbrand.N(-1, 2), []int{-1, 0, 1, 2})
		}
	})
}

func Test_S(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(gbrand.S(5)), 5)
		}
	})
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(gbrand.S(5, true)), 5)
		}
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(len(gbrand.S(0)), 0)
	})
}

func Test_B(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			b := gbrand.B(5)
			t.Assert(len(b), 5)
			t.AssertNE(b, make([]byte, 5))
		}
	})
	gbtest.C(t, func(t *gbtest.T) {
		b := gbrand.B(0)
		t.AssertNil(b)
	})
}

func Test_Str(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(gbrand.S(5)), 5)
		}
	})
}

func Test_RandStr(t *testing.T) {
	str := "我愛Ghostbb"
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 10; i++ {
			s := gbrand.Str(str, 100000)
			t.Assert(gbstr.Contains(s, "我"), true)
			t.Assert(gbstr.Contains(s, "愛"), true)
			t.Assert(gbstr.Contains(s, "G"), true)
			t.Assert(gbstr.Contains(s, "h"), true)
			t.Assert(gbstr.Contains(s, "o"), true)
			t.Assert(gbstr.Contains(s, "s"), true)
			t.Assert(gbstr.Contains(s, "t"), true)
			t.Assert(gbstr.Contains(s, "b"), true)
			t.Assert(gbstr.Contains(s, "b"), true)
			t.Assert(gbstr.Contains(s, "w"), false)
		}
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbrand.Str(str, 0), "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		list := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
		str := ""
		for _, s := range list {
			tmp := ""
			for i := 0; i < 15; i++ {
				tmp += tmp + s
			}
			str += tmp
		}
		t.Assert(len(gbrand.Str(str, 300)), 300)
	})
}

func Test_Digits(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(gbrand.Digits(5)), 5)
		}
	})
}

func Test_RandDigits(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(gbrand.Digits(5)), 5)
		}
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(len(gbrand.Digits(0)), 0)
	})
}

func Test_Letters(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(gbrand.Letters(5)), 5)
		}
	})
}

func Test_RandLetters(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(gbrand.Letters(5)), 5)
		}
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(len(gbrand.Letters(0)), 0)
	})
}

func Test_Perm(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			t.AssertIN(gbrand.Perm(5), []int{0, 1, 2, 3, 4})
		}
	})
}

func Test_Symbols(t *testing.T) {
	symbols := "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			syms := []byte(gbrand.Symbols(5))
			for _, sym := range syms {
				t.AssertNE(strings.Index(symbols, string(sym)), -1)
			}
		}
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbrand.Symbols(0), "")
	})
}
