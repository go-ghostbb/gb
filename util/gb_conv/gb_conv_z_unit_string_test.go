package gbconv_test

import (
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
	"time"
)

type stringStruct1 struct {
	Name string
}

type stringStruct2 struct {
	Name string
}

func (s *stringStruct1) String() string {
	return s.Name
}

func Test_String(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.String(int(123)), "123")
		t.AssertEQ(gbconv.String(int(-123)), "-123")
		t.AssertEQ(gbconv.String(int8(123)), "123")
		t.AssertEQ(gbconv.String(int8(-123)), "-123")
		t.AssertEQ(gbconv.String(int16(123)), "123")
		t.AssertEQ(gbconv.String(int16(-123)), "-123")
		t.AssertEQ(gbconv.String(int32(123)), "123")
		t.AssertEQ(gbconv.String(int32(-123)), "-123")
		t.AssertEQ(gbconv.String(int64(123)), "123")
		t.AssertEQ(gbconv.String(int64(-123)), "-123")
		t.AssertEQ(gbconv.String(int64(1552578474888)), "1552578474888")
		t.AssertEQ(gbconv.String(int64(-1552578474888)), "-1552578474888")

		t.AssertEQ(gbconv.String(uint(123)), "123")
		t.AssertEQ(gbconv.String(uint8(123)), "123")
		t.AssertEQ(gbconv.String(uint16(123)), "123")
		t.AssertEQ(gbconv.String(uint32(123)), "123")
		t.AssertEQ(gbconv.String(uint64(155257847488898765)), "155257847488898765")

		t.AssertEQ(gbconv.String(float32(123.456)), "123.456")
		t.AssertEQ(gbconv.String(float32(-123.456)), "-123.456")
		t.AssertEQ(gbconv.String(float64(1552578474888.456)), "1552578474888.456")
		t.AssertEQ(gbconv.String(float64(-1552578474888.456)), "-1552578474888.456")

		t.AssertEQ(gbconv.String(true), "true")
		t.AssertEQ(gbconv.String(false), "false")

		t.AssertEQ(gbconv.String([]byte("bytes")), "bytes")

		t.AssertEQ(gbconv.String(stringStruct1{"john"}), `{"Name":"john"}`)
		t.AssertEQ(gbconv.String(&stringStruct1{"john"}), "john")

		t.AssertEQ(gbconv.String(stringStruct2{"john"}), `{"Name":"john"}`)
		t.AssertEQ(gbconv.String(&stringStruct2{"john"}), `{"Name":"john"}`)

		var nilTime *time.Time = nil
		t.AssertEQ(gbconv.String(nilTime), "")
		var nilGTime *gbtime.Time = nil
		t.AssertEQ(gbconv.String(nilGTime), "")
	})
}
