package gbconv_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
)

func Test_Unsafe(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := "I love 鬼徹"
		t.AssertEQ(gbconv.UnsafeStrToBytes(s), []byte(s))
	})

	gbtest.C(t, func(t *gbtest.T) {
		b := []byte("I love 鬼徹")
		t.AssertEQ(gbconv.UnsafeBytesToStr(b), string(b))
	})
}
