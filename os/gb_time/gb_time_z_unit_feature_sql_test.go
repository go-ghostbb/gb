package gbtime_test

import (
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func TestTime_Scan(t1 *testing.T) {
	gbtest.C(t1, func(t *gbtest.T) {
		tt := gbtime.Time{}
		// test string
		s := gbtime.Now().String()
		t.Assert(tt.Scan(s), nil)
		t.Assert(tt.String(), s)
		// test nano
		n := gbtime.TimestampNano()
		t.Assert(tt.Scan(n), nil)
		t.Assert(tt.TimestampNano(), n)
		// test nil
		none := (*gbtime.Time)(nil)
		t.Assert(none.Scan(nil), nil)
		t.Assert(none, nil)
	})

}

func TestTime_Value(t1 *testing.T) {
	gbtest.C(t1, func(t *gbtest.T) {
		tt := gbtime.Now()
		s, err := tt.Value()
		t.AssertNil(err)
		t.Assert(s, tt.Time)
		// test nil
		none := (*gbtime.Time)(nil)
		s, err = none.Value()
		t.AssertNil(err)
		t.Assert(s, nil)

	})
}
