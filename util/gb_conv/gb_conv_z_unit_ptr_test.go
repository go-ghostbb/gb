package gbconv_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
)

func Test_Ptr_Functions(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var v interface{} = 1
		t.AssertEQ(gbconv.PtrAny(v), &v)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var v string = "1"
		t.AssertEQ(gbconv.PtrString(v), &v)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var v bool = true
		t.AssertEQ(gbconv.PtrBool(v), &v)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var v int = 1
		t.AssertEQ(gbconv.PtrInt(v), &v)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var v int8 = 1
		t.AssertEQ(gbconv.PtrInt8(v), &v)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var v int16 = 1
		t.AssertEQ(gbconv.PtrInt16(v), &v)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var v int32 = 1
		t.AssertEQ(gbconv.PtrInt32(v), &v)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var v int64 = 1
		t.AssertEQ(gbconv.PtrInt64(v), &v)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var v uint = 1
		t.AssertEQ(gbconv.PtrUint(v), &v)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var v uint8 = 1
		t.AssertEQ(gbconv.PtrUint8(v), &v)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var v uint16 = 1
		t.AssertEQ(gbconv.PtrUint16(v), &v)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var v uint32 = 1
		t.AssertEQ(gbconv.PtrUint32(v), &v)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var v uint64 = 1
		t.AssertEQ(gbconv.PtrUint64(v), &v)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var v float32 = 1.01
		t.AssertEQ(gbconv.PtrFloat32(v), &v)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var v float64 = 1.01
		t.AssertEQ(gbconv.PtrFloat64(v), &v)
	})
}
