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
