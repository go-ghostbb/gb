// go test *.go -bench=".*"

package gbregex_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbregex "ghostbb.io/gb/text/gb_regex"
	"strings"
	"testing"
)

var (
	PatternErr = `([\d+`
)

func Test_Quote(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := `[foo]` //`\[foo\]`
		t.Assert(gbregex.Quote(s1), `\[foo\]`)
	})
}

func Test_Validate(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var s1 = `(.+):(\d+)`
		t.Assert(gbregex.Validate(s1), nil)
		s1 = `((.+):(\d+)`
		t.Assert(gbregex.Validate(s1) == nil, false)
	})
}

func Test_IsMatch(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var pattern = `(.+):(\d+)`
		s1 := []byte(`sfs:2323`)
		t.Assert(gbregex.IsMatch(pattern, s1), true)
		s1 = []byte(`sfs2323`)
		t.Assert(gbregex.IsMatch(pattern, s1), false)
		s1 = []byte(`sfs:`)
		t.Assert(gbregex.IsMatch(pattern, s1), false)
		// error pattern
		t.Assert(gbregex.IsMatch(PatternErr, s1), false)
	})
}

func Test_IsMatchString(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var pattern = `(.+):(\d+)`
		s1 := `sfs:2323`
		t.Assert(gbregex.IsMatchString(pattern, s1), true)
		s1 = `sfs2323`
		t.Assert(gbregex.IsMatchString(pattern, s1), false)
		s1 = `sfs:`
		t.Assert(gbregex.IsMatchString(pattern, s1), false)
		// error pattern
		t.Assert(gbregex.IsMatchString(PatternErr, s1), false)
	})
}

func Test_Match(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		s := "acbb" + wantSubs + "dd"
		subs, err := gbregex.Match(re, []byte(s))
		t.AssertNil(err)
		if string(subs[0]) != wantSubs {
			t.Fatalf("regex:%s,Match(%q)[0] = %q; want %q", re, s, subs[0], wantSubs)
		}
		if string(subs[1]) != "aab" {
			t.Fatalf("Match(%q)[1] = %q; want %q", s, subs[1], "aab")
		}
		// error pattern
		_, err = gbregex.Match(PatternErr, []byte(s))
		t.AssertNE(err, nil)
	})
}

func Test_MatchString(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		s := "acbb" + wantSubs + "dd"
		subs, err := gbregex.MatchString(re, s)
		t.AssertNil(err)
		if string(subs[0]) != wantSubs {
			t.Fatalf("regex:%s,Match(%q)[0] = %q; want %q", re, s, subs[0], wantSubs)
		}
		if string(subs[1]) != "aab" {
			t.Fatalf("Match(%q)[1] = %q; want %q", s, subs[1], "aab")
		}
		// error pattern
		_, err = gbregex.MatchString(PatternErr, s)
		t.AssertNE(err, nil)
	})
}

func Test_MatchAll(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		s := "acbb" + wantSubs + "dd"
		s = s + `其他的` + s
		subs, err := gbregex.MatchAll(re, []byte(s))
		t.AssertNil(err)
		if string(subs[0][0]) != wantSubs {
			t.Fatalf("regex:%s,Match(%q)[0] = %q; want %q", re, s, subs[0][0], wantSubs)
		}
		if string(subs[0][1]) != "aab" {
			t.Fatalf("Match(%q)[1] = %q; want %q", s, subs[0][1], "aab")
		}

		if string(subs[1][0]) != wantSubs {
			t.Fatalf("regex:%s,Match(%q)[0] = %q; want %q", re, s, subs[1][0], wantSubs)
		}
		if string(subs[1][1]) != "aab" {
			t.Fatalf("Match(%q)[1] = %q; want %q", s, subs[1][1], "aab")
		}
		// error pattern
		_, err = gbregex.MatchAll(PatternErr, []byte(s))
		t.AssertNE(err, nil)
	})
}

func Test_MatchAllString(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		s := "acbb" + wantSubs + "dd"
		subs, err := gbregex.MatchAllString(re, s+`其他的`+s)
		t.AssertNil(err)
		if string(subs[0][0]) != wantSubs {
			t.Fatalf("regex:%s,Match(%q)[0] = %q; want %q", re, s, subs[0][0], wantSubs)
		}
		if string(subs[0][1]) != "aab" {
			t.Fatalf("Match(%q)[1] = %q; want %q", s, subs[0][1], "aab")
		}

		if string(subs[1][0]) != wantSubs {
			t.Fatalf("regex:%s,Match(%q)[0] = %q; want %q", re, s, subs[1][0], wantSubs)
		}
		if string(subs[1][1]) != "aab" {
			t.Fatalf("Match(%q)[1] = %q; want %q", s, subs[1][1], "aab")
		}
		// error pattern
		_, err = gbregex.MatchAllString(PatternErr, s)
		t.AssertNE(err, nil)
	})
}

func Test_Replace(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		replace := "12345"
		s := "acbb" + wantSubs + "dd"
		wanted := "acbb" + replace + "dd"
		replacedStr, err := gbregex.Replace(re, []byte(replace), []byte(s))
		t.AssertNil(err)
		if string(replacedStr) != wanted {
			t.Fatalf("regex:%s,old:%s; want %q", re, s, wanted)
		}
		// error pattern
		_, err = gbregex.Replace(PatternErr, []byte(replace), []byte(s))
		t.AssertNE(err, nil)
	})
}

func Test_ReplaceString(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		replace := "12345"
		s := "acbb" + wantSubs + "dd"
		wanted := "acbb" + replace + "dd"
		replacedStr, err := gbregex.ReplaceString(re, replace, s)
		t.AssertNil(err)
		if replacedStr != wanted {
			t.Fatalf("regex:%s,old:%s; want %q", re, s, wanted)
		}
		// error pattern
		_, err = gbregex.ReplaceString(PatternErr, replace, s)
		t.AssertNE(err, nil)
	})
}

func Test_ReplaceFun(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		//replace :="12345"
		s := "acbb" + wantSubs + "dd"
		wanted := "acbb[x" + wantSubs + "y]dd"
		wanted = "acbb" + "3个a" + "dd"
		replacedStr, err := gbregex.ReplaceFunc(re, []byte(s), func(s []byte) []byte {
			if strings.Contains(string(s), "aaa") {
				return []byte("3个a")
			}
			return []byte("[x" + string(s) + "y]")
		})
		t.AssertNil(err)
		if string(replacedStr) != wanted {
			t.Fatalf("regex:%s,old:%s; want %q", re, s, wanted)
		}
		// error pattern
		_, err = gbregex.ReplaceFunc(PatternErr, []byte(s), func(s []byte) []byte {
			return []byte("")
		})
		t.AssertNE(err, nil)
	})
}

func Test_ReplaceFuncMatch(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := []byte("1234567890")
		p := `(\d{3})(\d{3})(.+)`
		s0, e0 := gbregex.ReplaceFuncMatch(p, s, func(match [][]byte) []byte {
			return match[0]
		})
		t.Assert(e0, nil)
		t.Assert(s0, s)
		s1, e1 := gbregex.ReplaceFuncMatch(p, s, func(match [][]byte) []byte {
			return match[1]
		})
		t.Assert(e1, nil)
		t.Assert(s1, []byte("123"))
		s2, e2 := gbregex.ReplaceFuncMatch(p, s, func(match [][]byte) []byte {
			return match[2]
		})
		t.Assert(e2, nil)
		t.Assert(s2, []byte("456"))
		s3, e3 := gbregex.ReplaceFuncMatch(p, s, func(match [][]byte) []byte {
			return match[3]
		})
		t.Assert(e3, nil)
		t.Assert(s3, []byte("7890"))
		// error pattern
		_, err := gbregex.ReplaceFuncMatch(PatternErr, s, func(match [][]byte) []byte {
			return match[3]
		})
		t.AssertNE(err, nil)
	})
}

func Test_ReplaceStringFunc(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		re := "a(a+b+)b"
		wantSubs := "aaabb"
		//replace :="12345"
		s := "acbb" + wantSubs + "dd"
		wanted := "acbb[x" + wantSubs + "y]dd"
		wanted = "acbb" + "3个a" + "dd"
		replacedStr, err := gbregex.ReplaceStringFunc(re, s, func(s string) string {
			if strings.Contains(s, "aaa") {
				return "3个a"
			}
			return "[x" + s + "y]"
		})
		t.AssertNil(err)
		if replacedStr != wanted {
			t.Fatalf("regex:%s,old:%s; want %q", re, s, wanted)
		}
		// error pattern
		_, err = gbregex.ReplaceStringFunc(PatternErr, s, func(s string) string {
			return ""
		})
		t.AssertNE(err, nil)
	})
}

func Test_ReplaceStringFuncMatch(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := "1234567890"
		p := `(\d{3})(\d{3})(.+)`
		s0, e0 := gbregex.ReplaceStringFuncMatch(p, s, func(match []string) string {
			return match[0]
		})
		t.Assert(e0, nil)
		t.Assert(s0, s)
		s1, e1 := gbregex.ReplaceStringFuncMatch(p, s, func(match []string) string {
			return match[1]
		})
		t.Assert(e1, nil)
		t.Assert(s1, "123")
		s2, e2 := gbregex.ReplaceStringFuncMatch(p, s, func(match []string) string {
			return match[2]
		})
		t.Assert(e2, nil)
		t.Assert(s2, "456")
		s3, e3 := gbregex.ReplaceStringFuncMatch(p, s, func(match []string) string {
			return match[3]
		})
		t.Assert(e3, nil)
		t.Assert(s3, "7890")
		// error pattern
		_, err := gbregex.ReplaceStringFuncMatch(PatternErr, s, func(match []string) string {
			return ""
		})
		t.AssertNE(err, nil)
	})
}

func Test_Split(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		re := "a(a+b+)b"
		matched := "aaabb"
		item0 := "acbb"
		item1 := "dd"
		s := item0 + matched + item1
		t.Assert(gbregex.IsMatchString(re, matched), true)
		items := gbregex.Split(re, s) //split string with matched
		if items[0] != item0 {
			t.Fatalf("regex:%s,Split(%q) want %q", re, s, item0)
		}
		if items[1] != item1 {
			t.Fatalf("regex:%s,Split(%q) want %q", re, s, item0)
		}
	})

	gbtest.C(t, func(t *gbtest.T) {
		re := "a(a+b+)b"
		notmatched := "aaxbb"
		item0 := "acbb"
		item1 := "dd"
		s := item0 + notmatched + item1
		t.Assert(gbregex.IsMatchString(re, notmatched), false)
		items := gbregex.Split(re, s) //split string with notmatched then nosplitting
		if items[0] != s {
			t.Fatalf("regex:%s,Split(%q) want %q", re, s, item0)
		}
		// error pattern
		items = gbregex.Split(PatternErr, s)
		t.AssertEQ(items, nil)

	})
}
