package gbstr_test

import (
	"fmt"
	gbstr "ghostbb.io/gb/text/gb_str"
)

func ExampleCount() {
	var (
		str     = `ghostbb is very, very easy to use`
		substr1 = "ghostbb"
		substr2 = "very"
		result1 = gbstr.Count(str, substr1)
		result2 = gbstr.Count(str, substr2)
	)
	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 1
	// 2
}

func ExampleCountI() {
	var (
		str     = `ghostbb is very, very easy to use`
		substr1 = "GHOSTBB"
		substr2 = "VERY"
		result1 = gbstr.CountI(str, substr1)
		result2 = gbstr.CountI(str, substr2)
	)
	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// 1
	// 2
}

func ExampleToLower() {
	var (
		s      = `GHOSTBB`
		result = gbstr.ToLower(s)
	)
	fmt.Println(result)

	// Output:
	// ghostbb
}

func ExampleToUpper() {
	var (
		s      = `ghostbb`
		result = gbstr.ToUpper(s)
	)
	fmt.Println(result)

	// Output:
	// GHOSTBB
}

func ExampleUcFirst() {
	var (
		s      = `hello`
		result = gbstr.UcFirst(s)
	)
	fmt.Println(result)

	// Output:
	// Hello
}

func ExampleLcFirst() {
	var (
		str    = `Ghostbb`
		result = gbstr.LcFirst(str)
	)
	fmt.Println(result)

	// Output:
	// ghostbb
}

func ExampleUcWords() {
	var (
		str    = `hello world`
		result = gbstr.UcWords(str)
	)
	fmt.Println(result)

	// Output:
	// Hello World
}

func ExampleIsLetterLower() {
	fmt.Println(gbstr.IsLetterLower('a'))
	fmt.Println(gbstr.IsLetterLower('A'))

	// Output:
	// true
	// false
}

func ExampleIsLetterUpper() {
	fmt.Println(gbstr.IsLetterUpper('A'))
	fmt.Println(gbstr.IsLetterUpper('a'))

	// Output:
	// true
	// false
}

func ExampleIsNumeric() {
	fmt.Println(gbstr.IsNumeric("88"))
	fmt.Println(gbstr.IsNumeric("3.1415926"))
	fmt.Println(gbstr.IsNumeric("abc"))
	// Output:
	// true
	// true
	// false
}

func ExampleReverse() {
	var (
		str    = `123456`
		result = gbstr.Reverse(str)
	)
	fmt.Println(result)

	// Output:
	// 654321
}

func ExampleNumberFormat() {
	var (
		number       float64 = 123456
		decimals             = 2
		decPoint             = "."
		thousandsSep         = ","
		result               = gbstr.NumberFormat(number, decimals, decPoint, thousandsSep)
	)
	fmt.Println(result)

	// Output:
	// 123,456.00
}

func ExampleChunkSplit() {
	var (
		body     = `1234567890`
		chunkLen = 2
		end      = "#"
		result   = gbstr.ChunkSplit(body, chunkLen, end)
	)
	fmt.Println(result)

	// Output:
	// 12#34#56#78#90#
}

func ExampleCompare() {
	fmt.Println(gbstr.Compare("c", "c"))
	fmt.Println(gbstr.Compare("a", "b"))
	fmt.Println(gbstr.Compare("c", "b"))

	// Output:
	// 0
	// -1
	// 1
}

func ExampleEqual() {
	fmt.Println(gbstr.Equal(`A`, `a`))
	fmt.Println(gbstr.Equal(`A`, `A`))
	fmt.Println(gbstr.Equal(`A`, `B`))

	// Output:
	// true
	// true
	// false
}

func ExampleFields() {
	var (
		str    = `Hello World`
		result = gbstr.Fields(str)
	)
	fmt.Printf(`%#v`, result)

	// Output:
	// []string{"Hello", "World"}
}

func ExampleHasPrefix() {
	var (
		s      = `Hello World`
		prefix = "Hello"
		result = gbstr.HasPrefix(s, prefix)
	)
	fmt.Println(result)

	// Output:
	// true
}

func ExampleHasSuffix() {
	var (
		s      = `my best love is ghostbb`
		prefix = "ghostbb"
		result = gbstr.HasSuffix(s, prefix)
	)
	fmt.Println(result)

	// Output:
	// true
}

func ExampleCountWords() {
	var (
		str    = `ghostbb is very, very easy to use!`
		result = gbstr.CountWords(str)
	)
	fmt.Printf(`%#v`, result)

	// Output:
	// map[string]int{"easy":1, "ghostbb":1, "is":1, "to":1, "use!":1, "very":1, "very,":1}
}

func ExampleCountChars() {
	var (
		str    = `ghostbb`
		result = gbstr.CountChars(str)
	)
	fmt.Println(result)

	// May Output:
	// map[a:1 e:1 f:1 g:1 m:1 o:1 r:1]
}

func ExampleWordWrap() {
	{
		var (
			str    = `A very long woooooooooooooooooord. and something`
			width  = 8
			br     = "\n"
			result = gbstr.WordWrap(str, width, br)
		)
		fmt.Println(result)
	}
	{
		var (
			str    = `The quick brown fox jumped over the lazy dog.`
			width  = 20
			br     = "<br />\n"
			result = gbstr.WordWrap(str, width, br)
		)
		fmt.Printf("%v", result)
	}

	// Output:
	// A very
	// long
	// woooooooooooooooooord.
	// and
	// something
	// The quick brown fox<br />
	// jumped over the lazy<br />
	// dog.
}

func ExampleLenRune() {
	var (
		str    = `Ghostbb框架`
		result = gbstr.LenRune(str)
	)
	fmt.Println(result)

	// Output:
	// 9
}

func ExampleRepeat() {
	var (
		input      = `ghostbb `
		multiplier = 3
		result     = gbstr.Repeat(input, multiplier)
	)
	fmt.Println(result)

	// Output:
	// ghostbb ghostbb ghostbb
}

func ExampleShuffle() {
	var (
		str    = `123456`
		result = gbstr.Shuffle(str)
	)
	fmt.Println(result)

	// May Output:
	// 563214
}

func ExampleSplit() {
	var (
		str       = `a|b|c|d`
		delimiter = `|`
		result    = gbstr.Split(str, delimiter)
	)
	fmt.Printf(`%#v`, result)

	// Output:
	// []string{"a", "b", "c", "d"}
}

func ExampleSplitAndTrim() {
	var (
		str       = `a|b|||||c|d`
		delimiter = `|`
		result    = gbstr.SplitAndTrim(str, delimiter)
	)
	fmt.Printf(`%#v`, result)

	// Output:
	// []string{"a", "b", "c", "d"}
}

func ExampleJoin() {
	var (
		array  = []string{"ghostbb", "is", "very", "easy", "to", "use"}
		sep    = ` `
		result = gbstr.Join(array, sep)
	)
	fmt.Println(result)

	// Output:
	// ghostbb is very easy to use
}

func ExampleJoinAny() {
	var (
		sep    = `,`
		arr2   = []int{99, 73, 85, 66}
		result = gbstr.JoinAny(arr2, sep)
	)
	fmt.Println(result)

	// Output:
	// 99,73,85,66
}

func ExampleExplode() {
	var (
		str       = `Hello World`
		delimiter = " "
		result    = gbstr.Explode(delimiter, str)
	)
	fmt.Printf(`%#v`, result)

	// Output:
	// []string{"Hello", "World"}
}

func ExampleImplode() {
	var (
		pieces = []string{"ghostbb", "is", "very", "easy", "to", "use"}
		glue   = " "
		result = gbstr.Implode(glue, pieces)
	)
	fmt.Println(result)

	// Output:
	// ghostbb is very easy to use
}

func ExampleChr() {
	var (
		ascii  = 65 // A
		result = gbstr.Chr(ascii)
	)
	fmt.Println(result)

	// Output:
	// A
}

// '103' is the 'g' in ASCII
func ExampleOrd() {
	var (
		str    = `ghostbb`
		result = gbstr.Ord(str)
	)

	fmt.Println(result)

	// Output:
	// 103
}

func ExampleHideStr() {
	var (
		str     = `13800138000`
		percent = 40
		hide    = `*`
		result  = gbstr.HideStr(str, percent, hide)
	)
	fmt.Println(result)

	// Output:
	// 138****8000
}

func ExampleNl2Br() {
	var (
		str = `ghostbb
is
very
easy
to
use`
		result = gbstr.Nl2Br(str)
	)

	fmt.Println(result)

	// Output:
	// ghostbb<br>is<br>very<br>easy<br>to<br>use
}

func ExampleAddSlashes() {
	var (
		str    = `'aa'"bb"cc\r\n\d\t`
		result = gbstr.AddSlashes(str)
	)

	fmt.Println(result)

	// Output:
	// \'aa\'\"bb\"cc\\r\\n\\d\\t
}

func ExampleStripSlashes() {
	var (
		str    = `C:\\windows\\Ghostbb\\test`
		result = gbstr.StripSlashes(str)
	)
	fmt.Println(result)

	// Output:
	// C:\windows\Ghostbb\test
}

func ExampleQuoteMeta() {
	{
		var (
			str    = `.\+?[^]()`
			result = gbstr.QuoteMeta(str)
		)
		fmt.Println(result)
	}
	{
		var (
			str    = `https://ghostbb.io/pages/viewpage.action?pageId=1114327`
			result = gbstr.QuoteMeta(str)
		)
		fmt.Println(result)
	}

	// Output:
	// \.\\\+\?\[\^\]\(\)
	// https://ghostbb\.io/pages/viewpage\.action\?pageId=1114327

}

// array
func ExampleSearchArray() {
	var (
		array  = []string{"ghostbb", "is", "very", "nice"}
		str    = `ghostbb`
		result = gbstr.SearchArray(array, str)
	)
	fmt.Println(result)

	// Output:
	// 0
}

func ExampleInArray() {
	var (
		a      = []string{"ghostbb", "is", "very", "easy", "to", "use"}
		s      = "ghostbb"
		result = gbstr.InArray(a, s)
	)
	fmt.Println(result)

	// Output:
	// true
}

func ExamplePrefixArray() {
	var (
		strArray = []string{"tom", "lily", "john"}
	)

	gbstr.PrefixArray(strArray, "classA_")

	fmt.Println(strArray)

	// Output:
	// [classA_tom classA_lily classA_john]
}

// case
func ExampleCaseCamel() {
	var (
		str    = `hello world`
		result = gbstr.CaseCamel(str)
	)
	fmt.Println(result)

	// Output:
	// HelloWorld
}

func ExampleCaseCamelLower() {
	var (
		str    = `hello world`
		result = gbstr.CaseCamelLower(str)
	)
	fmt.Println(result)

	// Output:
	// helloWorld
}

func ExampleCaseSnake() {
	var (
		str    = `hello world`
		result = gbstr.CaseSnake(str)
	)
	fmt.Println(result)

	// Output:
	// hello_world
}

func ExampleCaseSnakeScreaming() {
	var (
		str    = `hello world`
		result = gbstr.CaseSnakeScreaming(str)
	)
	fmt.Println(result)

	// Output:
	// HELLO_WORLD
}

func ExampleCaseSnakeFirstUpper() {
	var (
		str    = `RGBCodeMd5`
		result = gbstr.CaseSnakeFirstUpper(str)
	)
	fmt.Println(result)

	// Output:
	// rgb_code_md5
}

func ExampleCaseKebab() {
	var (
		str    = `hello world`
		result = gbstr.CaseKebab(str)
	)
	fmt.Println(result)

	// Output:
	// hello-world
}

func ExampleCaseKebabScreaming() {
	var (
		str    = `hello world`
		result = gbstr.CaseKebabScreaming(str)
	)
	fmt.Println(result)

	// Output:
	// HELLO-WORLD
}

func ExampleCaseDelimited() {
	var (
		str    = `hello world`
		del    = byte('-')
		result = gbstr.CaseDelimited(str, del)
	)
	fmt.Println(result)

	// Output:
	// hello-world
}

func ExampleCaseDelimitedScreaming() {
	{
		var (
			str    = `hello world`
			del    = byte('-')
			result = gbstr.CaseDelimitedScreaming(str, del, true)
		)
		fmt.Println(result)
	}
	{
		var (
			str    = `hello world`
			del    = byte('-')
			result = gbstr.CaseDelimitedScreaming(str, del, false)
		)
		fmt.Println(result)
	}

	// Output:
	// HELLO-WORLD
	// hello-world
}

// contain
func ExampleContains() {
	{
		var (
			str    = `Hello World`
			substr = `Hello`
			result = gbstr.Contains(str, substr)
		)
		fmt.Println(result)
	}
	{
		var (
			str    = `Hello World`
			substr = `hello`
			result = gbstr.Contains(str, substr)
		)
		fmt.Println(result)
	}

	// Output:
	// true
	// false
}

func ExampleContainsI() {
	var (
		str     = `Hello World`
		substr  = "hello"
		result1 = gbstr.Contains(str, substr)
		result2 = gbstr.ContainsI(str, substr)
	)
	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// false
	// true
}

func ExampleContainsAny() {
	{
		var (
			s      = `ghostbb`
			chars  = "g"
			result = gbstr.ContainsAny(s, chars)
		)
		fmt.Println(result)
	}
	{
		var (
			s      = `ghostbb`
			chars  = "G"
			result = gbstr.ContainsAny(s, chars)
		)
		fmt.Println(result)
	}

	// Output:
	// true
	// false
}

// convert
func ExampleOctStr() {
	var (
		str    = `\346\200\241`
		result = gbstr.OctStr(str)
	)
	fmt.Println(result)

	// Output:
	// 怡
}

// domain
func ExampleIsSubDomain() {
	var (
		subDomain  = `s.ghostbb.io`
		mainDomain = `ghostbb.io`
		result     = gbstr.IsSubDomain(subDomain, mainDomain)
	)
	fmt.Println(result)

	// Output:
	// true
}

// levenshtein
func ExampleLevenshtein() {
	var (
		str1    = "Hello World"
		str2    = "hallo World"
		costIns = 1
		costRep = 1
		costDel = 1
		result  = gbstr.Levenshtein(str1, str2, costIns, costRep, costDel)
	)
	fmt.Println(result)

	// Output:
	// 2
}

// parse
func ExampleParse() {
	{
		var (
			str       = `v1=m&v2=n`
			result, _ = gbstr.Parse(str)
		)
		fmt.Println(result)
	}
	{
		var (
			str       = `v[a][a]=m&v[a][b]=n`
			result, _ = gbstr.Parse(str)
		)
		fmt.Println(result)
	}
	{
		// The form of nested Slice is not yet supported.
		var str = `v[][]=m&v[][]=n`
		result, err := gbstr.Parse(str)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	}
	{
		// This will produce an error.
		var str = `v=m&v[a]=n`
		result, err := gbstr.Parse(str)
		if err != nil {
			println(err)
		}
		fmt.Println(result)
	}
	{
		var (
			str       = `a .[[b=c`
			result, _ = gbstr.Parse(str)
		)
		fmt.Println(result)
	}

	// May Output:
	// map[v1:m v2:n]
	// map[v:map[a:map[a:m b:n]]]
	// map[v:map[]]
	// Error: expected type 'map[string]interface{}' for key 'v', but got 'string'
	// map[]
	// map[a___[b:c]
}

// pos
func ExamplePos() {
	var (
		haystack = `Hello World`
		needle   = `World`
		result   = gbstr.Pos(haystack, needle)
	)
	fmt.Println(result)

	// Output:
	// 6
}

func ExamplePosRune() {
	var (
		haystack = `Ghostbb是一款模塊化、高性能、企業級的Go基礎開發框架`
		needle   = `Gh`
		rNeedle  = `Go`
		posI     = gbstr.PosRune(haystack, needle)
		posR     = gbstr.PosRRune(haystack, rNeedle)
	)
	fmt.Println(posI)
	fmt.Println(posR)

	// Output:
	// 0
	// 22
}

func ExamplePosI() {
	var (
		haystack = `ghostbb is very, very easy to use`
		needle   = `very`
		posI     = gbstr.PosI(haystack, needle)
		posR     = gbstr.PosR(haystack, needle)
	)
	fmt.Println(posI)
	fmt.Println(posR)

	// Output:
	// 11
	// 17
}

func ExamplePosIRune() {
	{
		var (
			haystack    = `Ghostbb是一款模塊化、高性能、企業級的Go基礎開發框架`
			needle      = `高性能`
			startOffset = 10
			result      = gbstr.PosIRune(haystack, needle, startOffset)
		)
		fmt.Println(result)
	}
	{
		var (
			haystack    = `Ghostbb是一款模塊化、高性能、企業級的Go基礎開發框架`
			needle      = `高性能`
			startOffset = 30
			result      = gbstr.PosIRune(haystack, needle, startOffset)
		)
		fmt.Println(result)
	}

	// Output:
	// 14
	// -1
}

func ExamplePosR() {
	var (
		haystack = `ghostbb is very, very easy to use`
		needle   = `very`
		posI     = gbstr.PosI(haystack, needle)
		posR     = gbstr.PosR(haystack, needle)
	)
	fmt.Println(posI)
	fmt.Println(posR)

	// Output:
	// 11
	// 17
}

func ExamplePosRRune() {
	var (
		haystack = `Ghostbb是一款模塊化、高性能、企業級的Go基礎開發框架`
		needle   = `Gh`
		rNeedle  = `Go`
		posI     = gbstr.PosIRune(haystack, needle)
		posR     = gbstr.PosRRune(haystack, rNeedle)
	)
	fmt.Println(posI)
	fmt.Println(posR)

	// Output:
	// 0
	// 22
}

func ExamplePosRI() {
	var (
		haystack = `ghostbb is very, very easy to use`
		needle   = `VERY`
		posI     = gbstr.PosI(haystack, needle)
		posR     = gbstr.PosRI(haystack, needle)
	)
	fmt.Println(posI)
	fmt.Println(posR)

	// Output:
	// 11
	// 17
}

func ExamplePosRIRune() {
	var (
		haystack = `Ghostbb是一款模塊化、高性能、企業級的Go基礎開發框架`
		needle   = `Gh`
		rNeedle  = `GO`
		posI     = gbstr.PosIRune(haystack, needle)
		posR     = gbstr.PosRIRune(haystack, rNeedle)
	)
	fmt.Println(posI)
	fmt.Println(posR)

	// Output:
	// 0
	// 22
}

// replace
func ExampleReplace() {
	var (
		origin  = `golang is very nice!`
		search  = `golang`
		replace = `ghostbb`
		result  = gbstr.Replace(origin, search, replace)
	)
	fmt.Println(result)

	// Output:
	// ghostbb is very nice!
}

func ExampleReplaceI() {
	var (
		origin  = `golang is very nice!`
		search  = `GOLANG`
		replace = `ghostbb`
		result  = gbstr.ReplaceI(origin, search, replace)
	)
	fmt.Println(result)

	// Output:
	// ghostbb is very nice!
}

func ExampleReplaceByArray() {
	{
		var (
			origin = `golang is very nice`
			array  = []string{"olang", "hostbb"}
			result = gbstr.ReplaceByArray(origin, array)
		)
		fmt.Println(result)
	}
	{
		var (
			origin = `golang is very good`
			array  = []string{"golang", "ghostbb", "good", "nice"}
			result = gbstr.ReplaceByArray(origin, array)
		)
		fmt.Println(result)
	}

	// Output:
	// ghostbb is very nice
	// ghostbb is very nice
}

func ExampleReplaceIByArray() {
	var (
		origin = `golang is very Good`
		array  = []string{"Golang", "ghostbb", "GOOD", "nice"}
		result = gbstr.ReplaceIByArray(origin, array)
	)

	fmt.Println(result)

	// Output:
	// ghostbb is very nice
}

func ExampleReplaceByMap() {
	{
		var (
			origin   = `golang is very nice`
			replaces = map[string]string{
				"olang": "hostbb",
			}
			result = gbstr.ReplaceByMap(origin, replaces)
		)
		fmt.Println(result)
	}
	{
		var (
			origin   = `golang is very good`
			replaces = map[string]string{
				"golang": "ghostbb",
				"good":   "nice",
			}
			result = gbstr.ReplaceByMap(origin, replaces)
		)
		fmt.Println(result)
	}

	// Output:
	// ghostbb is very nice
	// ghostbb is very nice
}

func ExampleReplaceIByMap() {
	var (
		origin   = `golang is very nice`
		replaces = map[string]string{
			"Olang": "hostbb",
		}
		result = gbstr.ReplaceIByMap(origin, replaces)
	)
	fmt.Println(result)

	// Output:
	// ghostbb is very nice
}

// similartext
func ExampleSimilarText() {
	var (
		first   = `AaBbCcDd`
		second  = `ad`
		percent = 0.80
		result  = gbstr.SimilarText(first, second, &percent)
	)
	fmt.Println(result)

	// Output:
	// 2
}

// soundex
func ExampleSoundex() {
	var (
		str1    = `Hello`
		str2    = `Hallo`
		result1 = gbstr.Soundex(str1)
		result2 = gbstr.Soundex(str2)
	)
	fmt.Println(result1, result2)

	// Output:
	// H400 H400
}

// str
func ExampleStr() {
	var (
		haystack = `xxx.jpg`
		needle   = `.`
		result   = gbstr.Str(haystack, needle)
	)
	fmt.Println(result)

	// Output:
	// .jpg
}

func ExampleStrEx() {
	var (
		haystack = `https://ghostbb.io/index.html?a=1&b=2`
		needle   = `?`
		result   = gbstr.StrEx(haystack, needle)
	)
	fmt.Println(result)

	// Output:
	// a=1&b=2
}

func ExampleStrTill() {
	var (
		haystack = `https://ghostbb.io/index.html?test=123456`
		needle   = `?`
		result   = gbstr.StrTill(haystack, needle)
	)
	fmt.Println(result)

	// Output:
	// https://ghostbb.io/index.html?
}

func ExampleStrTillEx() {
	var (
		haystack = `https://ghostbb.io/index.html?test=123456`
		needle   = `?`
		result   = gbstr.StrTillEx(haystack, needle)
	)
	fmt.Println(result)

	// Output:
	// https://ghostbb.io/index.html
}

// substr
func ExampleSubStr() {
	var (
		str    = `1234567890`
		start  = 0
		length = 4
		subStr = gbstr.SubStr(str, start, length)
	)
	fmt.Println(subStr)

	// Output:
	// 1234
}

func ExampleSubStrRune() {
	var (
		str    = `Ghostbb是一款模塊化、高性能、企業級的Go基礎開發框架。`
		start  = 14
		length = 3
		subStr = gbstr.SubStrRune(str, start, length)
	)
	fmt.Println(subStr)

	// Output:
	// 高性能
}

func ExampleStrLimit() {
	var (
		str    = `123456789`
		length = 3
		suffix = `...`
		result = gbstr.StrLimit(str, length, suffix)
	)
	fmt.Println(result)

	// Output:
	// 123...
}

func ExampleStrLimitRune() {
	var (
		str    = `Ghostbb是一款模塊化、高性能、企業級的Go基礎開發框架。`
		length = 17
		suffix = "..."
		result = gbstr.StrLimitRune(str, length, suffix)
	)
	fmt.Println(result)

	// Output:
	// Ghostbb是一款模塊化、高性能...
}

func ExampleSubStrFrom() {
	var (
		str  = "我愛GhostbbGood"
		need = `愛`
	)

	fmt.Println(gbstr.SubStrFrom(str, need))

	// Output:
	// 愛GhostbbGood
}

func ExampleSubStrFromEx() {
	var (
		str  = "我愛GhostbbGood"
		need = `愛`
	)

	fmt.Println(gbstr.SubStrFromEx(str, need))

	// Output:
	// GhostbbGood
}

func ExampleSubStrFromR() {
	var (
		str  = "我愛GhostbbGood"
		need = `Go`
	)

	fmt.Println(gbstr.SubStrFromR(str, need))

	// Output:
	// Good
}

func ExampleSubStrFromREx() {
	var (
		str  = "我愛GhostbbGood"
		need = `Go`
	)

	fmt.Println(gbstr.SubStrFromREx(str, need))

	// Output:
	// od
}

// trim
func ExampleTrim() {
	var (
		str           = `*Hello World*`
		characterMask = "*"
		result        = gbstr.Trim(str, characterMask)
	)
	fmt.Println(result)

	// Output:
	// Hello World
}

func ExampleTrimStr() {
	var (
		str    = `Hello World`
		cut    = "World"
		count  = -1
		result = gbstr.TrimStr(str, cut, count)
	)
	fmt.Println(result)

	// Output:
	// Hello
}

func ExampleTrimLeft() {
	var (
		str           = `*Hello World*`
		characterMask = "*"
		result        = gbstr.TrimLeft(str, characterMask)
	)
	fmt.Println(result)

	// Output:
	// Hello World*
}

func ExampleTrimLeftStr() {
	var (
		str    = `**Hello World**`
		cut    = "*"
		count  = 1
		result = gbstr.TrimLeftStr(str, cut, count)
	)
	fmt.Println(result)

	// Output:
	// *Hello World**
}

func ExampleTrimRight() {
	var (
		str           = `**Hello World**`
		characterMask = "*def" // []byte{"*", "d", "e", "f"}
		result        = gbstr.TrimRight(str, characterMask)
	)
	fmt.Println(result)

	// Output:
	// **Hello Worl
}

func ExampleTrimRightStr() {
	var (
		str    = `Hello World!`
		cut    = "!"
		count  = -1
		result = gbstr.TrimRightStr(str, cut, count)
	)
	fmt.Println(result)

	// Output:
	// Hello World
}

func ExampleTrimAll() {
	var (
		str           = `*Hello World*`
		characterMask = "*"
		result        = gbstr.TrimAll(str, characterMask)
	)
	fmt.Println(result)

	// Output:
	// HelloWorld
}

// version
func ExampleCompareVersion() {
	fmt.Println(gbstr.CompareVersion("v2.11.9", "v2.10.8"))
	fmt.Println(gbstr.CompareVersion("1.10.8", "1.19.7"))
	fmt.Println(gbstr.CompareVersion("2.8.beta", "2.8"))

	// Output:
	// 1
	// -1
	// 0
}

func ExampleCompareVersionGo() {
	fmt.Println(gbstr.CompareVersionGo("v2.11.9", "v2.10.8"))
	fmt.Println(gbstr.CompareVersionGo("v4.20.1", "v4.20.1+incompatible"))
	fmt.Println(gbstr.CompareVersionGo(
		"v0.0.2-20180626092158-b2ccc119800e",
		"v1.0.1-20190626092158-b2ccc519800e",
	))

	// Output:
	// 1
	// 1
	// -1
}
