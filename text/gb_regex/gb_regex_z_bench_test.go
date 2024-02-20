// go test *.go -bench=".*"

package gbregex_test

import (
	gbregex "ghostbb.io/gb/text/gb_regex"
	"regexp"
	"testing"
)

var pattern = `(\w+).+\-\-\s*(.+)`

var src = `GF is best! -- John`

func Benchmark_GF_IsMatchString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbregex.IsMatchString(pattern, src)
	}
}

func Benchmark_GF_MatchString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbregex.MatchString(pattern, src)
	}
}

func Benchmark_Compile(b *testing.B) {
	var wcdRegexp = regexp.MustCompile(pattern)
	for i := 0; i < b.N; i++ {
		wcdRegexp.MatchString(src)
	}
}

func Benchmark_Compile_Actual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wcdRegexp := regexp.MustCompile(pattern)
		wcdRegexp.MatchString(src)
	}
}
