package gbconv_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
)

func Test_Basic(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		f32 := float32(123.456)
		i64 := int64(1552578474888)
		t.AssertEQ(gbconv.Int(f32), int(123))
		t.AssertEQ(gbconv.Int8(f32), int8(123))
		t.AssertEQ(gbconv.Int16(f32), int16(123))
		t.AssertEQ(gbconv.Int32(f32), int32(123))
		t.AssertEQ(gbconv.Int64(f32), int64(123))
		t.AssertEQ(gbconv.Int64(f32), int64(123))
		t.AssertEQ(gbconv.Uint(f32), uint(123))
		t.AssertEQ(gbconv.Uint8(f32), uint8(123))
		t.AssertEQ(gbconv.Uint16(f32), uint16(123))
		t.AssertEQ(gbconv.Uint32(f32), uint32(123))
		t.AssertEQ(gbconv.Uint64(f32), uint64(123))
		t.AssertEQ(gbconv.Float32(f32), float32(123.456))
		t.AssertEQ(gbconv.Float64(i64), float64(i64))
		t.AssertEQ(gbconv.Bool(f32), true)
		t.AssertEQ(gbconv.String(f32), "123.456")
		t.AssertEQ(gbconv.String(i64), "1552578474888")
	})

	gbtest.C(t, func(t *gbtest.T) {
		s := "-0xFF"
		t.Assert(gbconv.Int(s), int64(-0xFF))
	})
}

func Test_Duration(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		d := gbconv.Duration("1s")
		t.Assert(d.String(), "1s")
		t.Assert(d.Nanoseconds(), 1000000000)
	})
}

func Test_ConvertWithRefer(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.ConvertWithRefer("1", 100), 1)
		t.AssertEQ(gbconv.ConvertWithRefer("1.01", 1.111), 1.01)
		t.AssertEQ(gbconv.ConvertWithRefer("1.01", "1.111"), "1.01")
		t.AssertEQ(gbconv.ConvertWithRefer("1.01", false), true)
		t.AssertNE(gbconv.ConvertWithRefer("1.01", false), false)
	})
}
