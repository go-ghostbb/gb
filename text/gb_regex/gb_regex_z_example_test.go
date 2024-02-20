package gbregex_test

import (
	"bytes"
	"fmt"
	"ghostbb.io/gb/frame/g"
	gbregex "ghostbb.io/gb/text/gb_regex"
	"strings"
)

func ExampleIsMatch() {
	patternStr := `\d+`
	g.Dump(gbregex.IsMatch(patternStr, []byte("hello 2022! hello gb!")))
	g.Dump(gbregex.IsMatch(patternStr, nil))
	g.Dump(gbregex.IsMatch(patternStr, []byte("hello gb!")))

	// Output:
	// true
	// false
	// false
}

func ExampleIsMatchString() {
	patternStr := `\d+`
	g.Dump(gbregex.IsMatchString(patternStr, "hello 2022! hello gb!"))
	g.Dump(gbregex.IsMatchString(patternStr, "hello gb!"))
	g.Dump(gbregex.IsMatchString(patternStr, ""))

	// Output:
	// true
	// false
	// false
}

func ExampleMatch() {
	patternStr := `(\w+)=(\w+)`
	matchStr := "https://ghostbb.io/pages/viewpage.action?pageId=1114219&searchId=8QC5D1D2E!"
	// This method looks for the first match index
	result, err := gbregex.Match(patternStr, []byte(matchStr))
	g.Dump(result)
	g.Dump(err)

	// Output:
	// [
	//     "pageId=1114219",
	//     "pageId",
	//     "1114219",
	// ]
	// <nil>
}

func ExampleMatchString() {
	patternStr := `(\w+)=(\w+)`
	matchStr := "https://ghostbb.io/pages/viewpage.action?pageId=1114219&searchId=8QC5D1D2E!"
	// This method looks for the first match index
	result, err := gbregex.MatchString(patternStr, matchStr)
	g.Dump(result)
	g.Dump(err)

	// Output:
	// [
	//     "pageId=1114219",
	//     "pageId",
	//     "1114219",
	// ]
	// <nil>
}

func ExampleMatchAll() {
	patternStr := `(\w+)=(\w+)`
	matchStr := "https://ghostbb.io/pages/viewpage.action?pageId=1114219&searchId=8QC5D1D2E!"
	result, err := gbregex.MatchAll(patternStr, []byte(matchStr))
	g.Dump(result)
	g.Dump(err)

	// Output:
	//  [
	//     [
	//         "pageId=1114219",
	//         "pageId",
	//         "1114219",
	//     ],
	//     [
	//         "searchId=8QC5D1D2E",
	//         "searchId",
	//         "8QC5D1D2E",
	//     ],
	// ]
	// <nil>
}

func ExampleMatchAllString() {
	patternStr := `(\w+)=(\w+)`
	matchStr := "https://ghostbb.io/pages/viewpage.action?pageId=1114219&searchId=8QC5D1D2E!"
	result, err := gbregex.MatchAllString(patternStr, matchStr)
	g.Dump(result)
	g.Dump(err)

	// Output:
	// [
	//     [
	//         "pageId=1114219",
	//         "pageId",
	//         "1114219",
	//     ],
	//     [
	//         "searchId=8QC5D1D2E",
	//         "searchId",
	//         "8QC5D1D2E",
	//     ],
	// ]
	// <nil>
}

func ExampleQuote() {
	result := gbregex.Quote(`[1-9]\d+`)
	fmt.Println(result)

	// Output:
	// \[1-9\]\\d\+
}

func ExampleReplace() {
	var (
		patternStr  = `\d+`
		str         = "hello gb 2020!"
		repStr      = "2021"
		result, err = gbregex.Replace(patternStr, []byte(repStr), []byte(str))
	)
	g.Dump(err)
	g.Dump(result)

	// Output:
	// <nil>
	// "hello gb 2021!"
}

func ExampleReplaceFunc() {
	// In contrast to [ExampleReplaceFunc]
	// the result contains the `pattern' of all subpattern that use the matching function
	result, err := gbregex.ReplaceFuncMatch(`(\d+)~(\d+)`, []byte("hello gb 2018~2020!"), func(match [][]byte) []byte {
		g.Dump(match)
		match[2] = []byte("2021")
		return bytes.Join(match[1:], []byte("~"))
	})
	g.Dump(result)
	g.Dump(err)

	// Output:
	// [
	//     "2018~2020",
	//     "2018",
	//     "2020",
	// ]
	// "hello gb 2018~2021!"
	// <nil>
}

func ExampleReplaceFuncMatch() {
	var (
		patternStr = `(\d+)~(\d+)`
		str        = "hello gb 2018~2020!"
	)
	// In contrast to [ExampleReplaceFunc]
	// the result contains the `pattern' of all subpatterns that use the matching function
	result, err := gbregex.ReplaceFuncMatch(patternStr, []byte(str), func(match [][]byte) []byte {
		g.Dump(match)
		match[2] = []byte("2021")
		return bytes.Join(match[1:], []byte("-"))
	})
	g.Dump(result)
	g.Dump(err)

	// Output:
	// [
	//     "2018~2020",
	//     "2018",
	//     "2020",
	// ]
	// "hello gb 2018-2021!"
	// <nil>
}

func ExampleReplaceString() {
	patternStr := `\d+`
	str := "hello gb 2020!"
	replaceStr := "2021"
	result, err := gbregex.ReplaceString(patternStr, replaceStr, str)

	g.Dump(result)
	g.Dump(err)

	// Output:
	// "hello gb 2021!"
	// <nil>
}

func ExampleReplaceStringFunc() {
	replaceStrMap := map[string]string{
		"2020": "2021",
	}
	// When the regular statement can match multiple results
	// func can be used to further control the value that needs to be modified
	result, err := gbregex.ReplaceStringFunc(`\d+`, `hello gb 2018~2020!`, func(b string) string {
		g.Dump(b)
		if replaceStr, ok := replaceStrMap[b]; ok {
			return replaceStr
		}
		return b
	})
	g.Dump(result)
	g.Dump(err)

	result, err = gbregex.ReplaceStringFunc(`[a-z]*`, "gb@ghostbb.io", strings.ToUpper)
	g.Dump(result)
	g.Dump(err)

	// Output:
	// "2018"
	// "2020"
	// "hello gb 2018~2021!"
	// <nil>
	// "GB@GHOSTBB.IO"
	// <nil>
}

func ExampleReplaceStringFuncMatch() {
	var (
		patternStr = `([A-Z])\w+`
		str        = "hello Golang 2018~2021!"
	)
	// In contrast to [ExampleReplaceFunc]
	// the result contains the `pattern' of all subpatterns that use the matching function
	result, err := gbregex.ReplaceStringFuncMatch(patternStr, str, func(match []string) string {
		g.Dump(match)
		match[0] = "Gb"
		return match[0]
	})
	g.Dump(result)
	g.Dump(err)

	// Output:
	// [
	//     "Golang",
	//     "G",
	// ]
	// "hello Gb 2018~2021!"
	// <nil>
}

func ExampleSplit() {
	patternStr := `\d+`
	str := "hello2020gb"
	result := gbregex.Split(patternStr, str)
	g.Dump(result)

	// Output:
	// [
	//     "hello",
	//     "gb",
	// ]
}

func ExampleValidate() {
	// Valid match statement
	fmt.Println(gbregex.Validate(`\d+`))
	// Mismatched statement
	fmt.Println(gbregex.Validate(`[a-9]\d+`))

	// Output:
	// <nil>
	// regexp.Compile failed for pattern "[a-9]\d+": error parsing regexp: invalid character class range: `a-9`
}
