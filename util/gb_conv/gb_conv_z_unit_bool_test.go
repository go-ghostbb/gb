package gbconv_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
)

type boolStruct struct{}

func Test_Bool(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var any interface{} = nil
		t.AssertEQ(gbconv.Bool(any), false)
		t.AssertEQ(gbconv.Bool(false), false)
		t.AssertEQ(gbconv.Bool(nil), false)
		t.AssertEQ(gbconv.Bool(0), false)
		t.AssertEQ(gbconv.Bool("0"), false)
		t.AssertEQ(gbconv.Bool(""), false)
		t.AssertEQ(gbconv.Bool("false"), false)
		t.AssertEQ(gbconv.Bool("off"), false)
		t.AssertEQ(gbconv.Bool([]byte{}), false)
		t.AssertEQ(gbconv.Bool([]string{}), false)
		t.AssertEQ(gbconv.Bool([]interface{}{}), false)
		t.AssertEQ(gbconv.Bool([]map[int]int{}), false)

		t.AssertEQ(gbconv.Bool("1"), true)
		t.AssertEQ(gbconv.Bool("on"), true)
		t.AssertEQ(gbconv.Bool(1), true)
		t.AssertEQ(gbconv.Bool(123.456), true)
		t.AssertEQ(gbconv.Bool(boolStruct{}), true)
		t.AssertEQ(gbconv.Bool(&boolStruct{}), true)
	})
}
