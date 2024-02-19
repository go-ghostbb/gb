package gbutil_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbutil "ghostbb.io/gb/util/gb_util"
	"testing"
)

func Test_ComparatorString(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorString(1, 1), 0)
		t.Assert(gbutil.ComparatorString(1, 2), -1)
		t.Assert(gbutil.ComparatorString(2, 1), 1)
	})
}

func Test_ComparatorInt(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorInt(1, 1), 0)
		t.Assert(gbutil.ComparatorInt(1, 2), -1)
		t.Assert(gbutil.ComparatorInt(2, 1), 1)
	})
}

func Test_ComparatorInt8(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorInt8(1, 1), 0)
		t.Assert(gbutil.ComparatorInt8(1, 2), -1)
		t.Assert(gbutil.ComparatorInt8(2, 1), 1)
	})
}

func Test_ComparatorInt16(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorInt16(1, 1), 0)
		t.Assert(gbutil.ComparatorInt16(1, 2), -1)
		t.Assert(gbutil.ComparatorInt16(2, 1), 1)
	})
}

func Test_ComparatorInt32(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorInt32(1, 1), 0)
		t.Assert(gbutil.ComparatorInt32(1, 2), -1)
		t.Assert(gbutil.ComparatorInt32(2, 1), 1)
	})
}

func Test_ComparatorInt64(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorInt64(1, 1), 0)
		t.Assert(gbutil.ComparatorInt64(1, 2), -1)
		t.Assert(gbutil.ComparatorInt64(2, 1), 1)
	})
}

func Test_ComparatorUint(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorUint(1, 1), 0)
		t.Assert(gbutil.ComparatorUint(1, 2), -1)
		t.Assert(gbutil.ComparatorUint(2, 1), 1)
	})
}

func Test_ComparatorUint8(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorUint8(1, 1), 0)
		t.Assert(gbutil.ComparatorUint8(2, 6), 252)
		t.Assert(gbutil.ComparatorUint8(2, 1), 1)
	})
}

func Test_ComparatorUint16(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorUint16(1, 1), 0)
		t.Assert(gbutil.ComparatorUint16(1, 2), 65535)
		t.Assert(gbutil.ComparatorUint16(2, 1), 1)
	})
}

func Test_ComparatorUint32(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorUint32(1, 1), 0)
		t.Assert(gbutil.ComparatorUint32(-1000, 2147483640), 2147482656)
		t.Assert(gbutil.ComparatorUint32(2, 1), 1)
	})
}

func Test_ComparatorUint64(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorUint64(1, 1), 0)
		t.Assert(gbutil.ComparatorUint64(1, 2), -1)
		t.Assert(gbutil.ComparatorUint64(2, 1), 1)
	})
}

func Test_ComparatorFloat32(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorFloat32(1, 1), 0)
		t.Assert(gbutil.ComparatorFloat32(1, 2), -1)
		t.Assert(gbutil.ComparatorFloat32(2, 1), 1)
	})
}

func Test_ComparatorFloat64(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorFloat64(1, 1), 0)
		t.Assert(gbutil.ComparatorFloat64(1, 2), -1)
		t.Assert(gbutil.ComparatorFloat64(2, 1), 1)
	})
}

func Test_ComparatorByte(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorByte(1, 1), 0)
		t.Assert(gbutil.ComparatorByte(1, 2), 255)
		t.Assert(gbutil.ComparatorByte(2, 1), 1)
	})
}

func Test_ComparatorRune(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorRune(1, 1), 0)
		t.Assert(gbutil.ComparatorRune(1, 2), -1)
		t.Assert(gbutil.ComparatorRune(2, 1), 1)
	})
}

func Test_ComparatorTime(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		j := gbutil.ComparatorTime("2019-06-14", "2019-06-14")
		t.Assert(j, 0)

		k := gbutil.ComparatorTime("2019-06-15", "2019-06-14")
		t.Assert(k, 1)

		l := gbutil.ComparatorTime("2019-06-13", "2019-06-14")
		t.Assert(l, -1)
	})
}

func Test_ComparatorFloat32OfFixed(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorFloat32(0.1, 0.1), 0)
		t.Assert(gbutil.ComparatorFloat32(1.1, 2.1), -1)
		t.Assert(gbutil.ComparatorFloat32(2.1, 1.1), 1)
	})
}

func Test_ComparatorFloat64OfFixed(t *testing.T) {

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.ComparatorFloat64(0.1, 0.1), 0)
		t.Assert(gbutil.ComparatorFloat64(1.1, 2.1), -1)
		t.Assert(gbutil.ComparatorFloat64(2.1, 1.1), 1)
	})
}
