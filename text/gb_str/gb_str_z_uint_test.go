// go test *.go -bench=".*"

package gbstr_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"

	"testing"
)

func Test_ToLower(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "abcdEFG亂入的中文abcdefg"
		e1 := "abcdefg亂入的中文abcdefg"
		t.Assert(gbstr.ToLower(s1), e1)
	})
}

func Test_ToUpper(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "abcdEFG亂入的中文abcdefg"
		e1 := "ABCDEFG亂入的中文ABCDEFG"
		t.Assert(gbstr.ToUpper(s1), e1)
	})
}

func Test_UcFirst(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "abcdEFG亂入的中文abcdefg"
		e1 := "AbcdEFG亂入的中文abcdefg"
		t.Assert(gbstr.UcFirst(""), "")
		t.Assert(gbstr.UcFirst(s1), e1)
		t.Assert(gbstr.UcFirst(e1), e1)
	})
}

func Test_LcFirst(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "AbcdEFG亂入的中文abcdefg"
		e1 := "abcdEFG亂入的中文abcdefg"
		t.Assert(gbstr.LcFirst(""), "")
		t.Assert(gbstr.LcFirst(s1), e1)
		t.Assert(gbstr.LcFirst(e1), e1)
	})
}

func Test_UcWords(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s1 := "我愛GB: i love go ghostbb"
		e1 := "我愛GB: I Love Go Ghostbb"
		t.Assert(gbstr.UcWords(s1), e1)
	})
}

func Test_IsLetterLower(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.IsLetterLower('a'), true)
		t.Assert(gbstr.IsLetterLower('A'), false)
		t.Assert(gbstr.IsLetterLower('1'), false)
	})
}

func Test_IsLetterUpper(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.IsLetterUpper('a'), false)
		t.Assert(gbstr.IsLetterUpper('A'), true)
		t.Assert(gbstr.IsLetterUpper('1'), false)
	})
}

func Test_IsNumeric(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.IsNumeric("1a我"), false)
		t.Assert(gbstr.IsNumeric("0123"), true)
		t.Assert(gbstr.IsNumeric("我是台灣人"), false)
		t.Assert(gbstr.IsNumeric("1.2.3.4"), false)
	})
}

func Test_SubStr(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.SubStr("我愛Ghostbb", 0), "我愛Ghostbb")
		t.Assert(gbstr.SubStr("我愛Ghostbb", 6), "Ghostbb")
		t.Assert(gbstr.SubStr("我愛Ghostbb", 6, 2), "Gh")
		t.Assert(gbstr.SubStr("我愛Ghostbb", -1, 30), "b")
		t.Assert(gbstr.SubStr("我愛Ghostbb", 30, 30), "")
		t.Assert(gbstr.SubStr("abcdef", 0, -1), "abcde")
		t.Assert(gbstr.SubStr("abcdef", 2, -1), "cde")
		t.Assert(gbstr.SubStr("abcdef", 4, -4), "")
		t.Assert(gbstr.SubStr("abcdef", -3, -1), "de")
	})
}

func Test_SubStrRune(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.SubStrRune("我愛Ghostbb", 0), "我愛Ghostbb")
		t.Assert(gbstr.SubStrRune("我愛Ghostbb", 2), "Ghostbb")
		t.Assert(gbstr.SubStrRune("我愛Ghostbb", 2, 2), "Gh")
		t.Assert(gbstr.SubStrRune("我愛Ghostbb", -1, 30), "b")
		t.Assert(gbstr.SubStrRune("我愛Ghostbb", 30, 30), "")
		t.Assert(gbstr.SubStrRune("abcdef", 0, -1), "abcde")
		t.Assert(gbstr.SubStrRune("abcdef", 2, -1), "cde")
		t.Assert(gbstr.SubStrRune("abcdef", 4, -4), "")
		t.Assert(gbstr.SubStrRune("abcdef", -3, -1), "de")
		t.Assert(gbstr.SubStrRune("我愛Ghostbb呵呵", -3, 100), "b呵呵")
		t.Assert(gbstr.SubStrRune("abcdef哈哈", -3, -1), "f哈")
		t.Assert(gbstr.SubStrRune("ab我愛Ghostbbcdef哈哈", -3, -1), "f哈")
		t.Assert(gbstr.SubStrRune("我愛Ghostbb", 0, 3), "我愛G")
	})
}

func Test_StrLimit(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.StrLimit("我愛Ghostbb", 6), "我愛...")
		t.Assert(gbstr.StrLimit("我愛Ghostbb", 6, ""), "我愛")
		t.Assert(gbstr.StrLimit("我愛Ghostbb", 6, "**"), "我愛**")
		t.Assert(gbstr.StrLimit("我愛Ghostbb", 8, ""), "我愛Gh")
		t.Assert(gbstr.StrLimit("*", 4, ""), "*")
	})
}

func Test_StrLimitRune(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.StrLimitRune("我愛Ghostbb", 2), "我愛...")
		t.Assert(gbstr.StrLimitRune("我愛Ghostbb", 2, ""), "我愛")
		t.Assert(gbstr.StrLimitRune("我愛Ghostbb", 2, "**"), "我愛**")
		t.Assert(gbstr.StrLimitRune("我愛Ghostbb", 4, ""), "我愛Gh")
		t.Assert(gbstr.StrLimitRune("*", 4, ""), "*")
	})
}

func Test_HasPrefix(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.HasPrefix("我愛Ghostbb", "我愛"), true)
		t.Assert(gbstr.HasPrefix("en我愛Ghostbb", "我愛"), false)
		t.Assert(gbstr.HasPrefix("en我愛Ghostbb", "en"), true)
	})
}

func Test_HasSuffix(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.HasSuffix("我愛Ghostbb", "Ghostbb"), true)
		t.Assert(gbstr.HasSuffix("en我愛Ghostbb", "a"), false)
		t.Assert(gbstr.HasSuffix("Ghostbb很棒", "棒"), true)
	})
}

func Test_Reverse(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.Reverse("我愛123"), "321愛我")
	})
}

func Test_NumberFormat(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.NumberFormat(1234567.8910, 2, ".", ","), "1,234,567.89")
		t.Assert(gbstr.NumberFormat(1234567.8910, 2, "#", "/"), "1/234/567#89")
		t.Assert(gbstr.NumberFormat(-1234567.8910, 2, "#", "/"), "-1/234/567#89")
	})
}

func Test_ChunkSplit(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.ChunkSplit("1234", 1, "#"), "1#2#3#4#")
		t.Assert(gbstr.ChunkSplit("我愛123", 1, "#"), "我#愛#1#2#3#")
		t.Assert(gbstr.ChunkSplit("1234", 1, ""), "1\r\n2\r\n3\r\n4\r\n")
	})
}

func Test_SplitAndTrim(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := `

010    

020  

`
		a := gbstr.SplitAndTrim(s, "\n", "0")
		t.Assert(len(a), 2)
		t.Assert(a[0], "1")
		t.Assert(a[1], "2")
	})
}

func Test_Fields(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.Fields("我愛 Go Frame"), []string{
			"我愛", "Go", "Frame",
		})
	})
}

func Test_CountWords(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.CountWords("我愛 Go Go Go"), map[string]int{
			"Go": 3,
			"我愛": 1,
		})
	})
}

func Test_CountChars(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.CountChars("我愛 Go Go Go"), map[string]int{
			" ": 3,
			"G": 3,
			"o": 3,
			"我": 1,
			"愛": 1,
		})
		t.Assert(gbstr.CountChars("我愛 Go Go Go", true), map[string]int{
			"G": 3,
			"o": 3,
			"我": 1,
			"愛": 1,
		})
	})
}

func Test_LenRune(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.LenRune("1234"), 4)
		t.Assert(gbstr.LenRune("我愛Ghostbb"), 9)
	})
}

func Test_Repeat(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.Repeat("go", 3), "gogogo")
		t.Assert(gbstr.Repeat("好的", 3), "好的好的好的")
	})
}

func Test_Str(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.Str("name@example.com", "@"), "@example.com")
		t.Assert(gbstr.Str("name@example.com", ""), "")
		t.Assert(gbstr.Str("name@example.com", "z"), "")
	})
}

func Test_StrEx(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.StrEx("name@example.com", "@"), "example.com")
		t.Assert(gbstr.StrEx("name@example.com", ""), "")
		t.Assert(gbstr.StrEx("name@example.com", "z"), "")
	})
}

func Test_StrTill(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.StrTill("name@example.com", "@"), "name@")
		t.Assert(gbstr.StrTill("name@example.com", ""), "")
		t.Assert(gbstr.StrTill("name@example.com", "z"), "")
	})
}

func Test_StrTillEx(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.StrTillEx("name@example.com", "@"), "name")
		t.Assert(gbstr.StrTillEx("name@example.com", ""), "")
		t.Assert(gbstr.StrTillEx("name@example.com", "z"), "")
	})
}

func Test_Shuffle(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(len(gbstr.Shuffle("123456")), 6)
	})
}

func Test_Split(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.Split("1.2", "."), []string{"1", "2"})
		t.Assert(gbstr.Split("我愛 - Ghostbb", " - "), []string{"我愛", "Ghostbb"})
	})
}

func Test_Join(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.Join([]string{"我愛", "Ghostbb"}, " - "), "我愛 - Ghostbb")
	})
}

func Test_Explode(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.Explode(" - ", "我愛 - Ghostbb"), []string{"我愛", "Ghostbb"})
	})
}

func Test_Implode(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.Implode(" - ", []string{"我愛", "Ghostbb"}), "我愛 - Ghostbb")
	})
}

func Test_Chr(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.Chr(65), "A")
	})
}

func Test_Ord(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.Ord("A"), 65)
	})
}

func Test_HideStr(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.HideStr("15928008611", 40, "*"), "159****8611")
		t.Assert(gbstr.HideStr("john@ghostbb.io", 40, "*"), "jo*n@ghostbb.io")
		t.Assert(gbstr.HideStr("張三", 50, "*"), "張*")
		t.Assert(gbstr.HideStr("張小三", 50, "*"), "張*三")
		t.Assert(gbstr.HideStr("歐陽小三", 50, "*"), "歐**三")
	})
}

func Test_Nl2Br(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.Nl2Br("1\n2"), "1<br>2")
		t.Assert(gbstr.Nl2Br("1\r\n2"), "1<br>2")
		t.Assert(gbstr.Nl2Br("1\r\n2", true), "1<br />2")
	})
}

func Test_AddSlashes(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.AddSlashes(`1'2"3\`), `1\'2\"3\\`)
	})
}

func Test_StripSlashes(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.StripSlashes(`1\'2\"3\\`), `1'2"3\`)
	})
}

func Test_QuoteMeta(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.QuoteMeta(`.\+*?[^]($)`), `\.\\\+\*\?\[\^\]\(\$\)`)
		t.Assert(gbstr.QuoteMeta(`.\+*台灣?[^]($)`), `\.\\\+\*台灣\?\[\^\]\(\$\)`)
		t.Assert(gbstr.QuoteMeta(`.''`, `'`), `.\'\'`)
		t.Assert(gbstr.QuoteMeta(`台灣.''`, `'`), `台灣.\'\'`)
	})
}

func Test_Count(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := "abcdaAD"
		t.Assert(gbstr.Count(s, "0"), 0)
		t.Assert(gbstr.Count(s, "a"), 2)
		t.Assert(gbstr.Count(s, "b"), 1)
		t.Assert(gbstr.Count(s, "d"), 1)
	})
}

func Test_CountI(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := "abcdaAD"
		t.Assert(gbstr.CountI(s, "0"), 0)
		t.Assert(gbstr.CountI(s, "a"), 3)
		t.Assert(gbstr.CountI(s, "b"), 1)
		t.Assert(gbstr.CountI(s, "d"), 2)
	})
}

func Test_Compare(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.Compare("a", "b"), -1)
		t.Assert(gbstr.Compare("a", "a"), 0)
		t.Assert(gbstr.Compare("b", "a"), 1)
	})
}

func Test_Equal(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.Equal("a", "A"), true)
		t.Assert(gbstr.Equal("a", "a"), true)
		t.Assert(gbstr.Equal("b", "a"), false)
	})
}

func Test_Contains(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.Contains("abc", "a"), true)
		t.Assert(gbstr.Contains("abc", "A"), false)
		t.Assert(gbstr.Contains("abc", "ab"), true)
		t.Assert(gbstr.Contains("abc", "abc"), true)
	})
}

func Test_ContainsI(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.ContainsI("abc", "a"), true)
		t.Assert(gbstr.ContainsI("abc", "A"), true)
		t.Assert(gbstr.ContainsI("abc", "Ab"), true)
		t.Assert(gbstr.ContainsI("abc", "ABC"), true)
		t.Assert(gbstr.ContainsI("abc", "ABCD"), false)
		t.Assert(gbstr.ContainsI("abc", "D"), false)
	})
}

func Test_ContainsAny(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.ContainsAny("abc", "a"), true)
		t.Assert(gbstr.ContainsAny("abc", "cd"), true)
		t.Assert(gbstr.ContainsAny("abc", "de"), false)
		t.Assert(gbstr.ContainsAny("abc", "A"), false)
	})
}

func Test_SubStrFrom(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.SubStrFrom("我愛GhostbbGood", `G`), "GhostbbGood")
		t.Assert(gbstr.SubStrFrom("我愛GhostbbGood", `GG`), "")
		t.Assert(gbstr.SubStrFrom("我愛GhostbbGood", `我`), "我愛GhostbbGood")
		t.Assert(gbstr.SubStrFrom("我愛GhostbbGood", `bb`), "bbGood")
	})
}

func Test_SubStrFromEx(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.SubStrFromEx("我愛GhostbbGood", `Ghost`), "bbGood")
		t.Assert(gbstr.SubStrFromEx("我愛GhostbbGood", `GG`), "")
		t.Assert(gbstr.SubStrFromEx("我愛GhostbbGood", `我`), "愛GhostbbGood")
		t.Assert(gbstr.SubStrFromEx("我愛GhostbbGood", `bb`), `Good`)
	})
}

func Test_SubStrFromR(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.SubStrFromR("我愛GhostbbGood", `G`), "Good")
		t.Assert(gbstr.SubStrFromR("我愛GhostbbGood", `GG`), "")
		t.Assert(gbstr.SubStrFromR("我愛GhostbbGood", `我`), "我愛GhostbbGood")
		t.Assert(gbstr.SubStrFromR("我愛GhostbbGood", `bb`), "bbGood")
	})
}

func Test_SubStrFromREx(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbstr.SubStrFromREx("我愛GhostbbGood", `G`), "ood")
		t.Assert(gbstr.SubStrFromREx("我愛GhostbbGood", `GG`), "")
		t.Assert(gbstr.SubStrFromREx("我愛GhostbbGood", `我`), "愛GhostbbGood")
		t.Assert(gbstr.SubStrFromREx("我愛GhostbbGood", `bb`), `Good`)
	})
}
