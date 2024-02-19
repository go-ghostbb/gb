package gbconv_test

import (
	gbvar "ghostbb.io/gb/container/gb_var"
	"ghostbb.io/gb/frame/g"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
	"time"
)

func Test_Time(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Duration(""), time.Duration(int64(0)))
		t.AssertEQ(gbconv.GBTime(""), gbtime.New())
		t.AssertEQ(gbconv.GBTime(nil), nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		s := "2011-10-10 01:02:03.456"
		t.AssertEQ(gbconv.GBTime(s), gbtime.NewFromStr(s))
		t.AssertEQ(gbconv.Time(nil), time.Time{})
		t.AssertEQ(gbconv.Time(s), gbtime.NewFromStr(s).Time)
		t.AssertEQ(gbconv.Duration(100), 100*time.Nanosecond)
	})
	gbtest.C(t, func(t *gbtest.T) {
		s := "01:02:03.456"
		t.AssertEQ(gbconv.GBTime(s).Hour(), 1)
		t.AssertEQ(gbconv.GBTime(s).Minute(), 2)
		t.AssertEQ(gbconv.GBTime(s).Second(), 3)
		t.AssertEQ(gbconv.GBTime(s), gbtime.NewFromStr(s))
		t.AssertEQ(gbconv.Time(s), gbtime.NewFromStr(s).Time)
	})
	gbtest.C(t, func(t *gbtest.T) {
		s := "0000-01-01 01:02:03"
		t.AssertEQ(gbconv.GBTime(s).Year(), 0)
		t.AssertEQ(gbconv.GBTime(s).Month(), 1)
		t.AssertEQ(gbconv.GBTime(s).Day(), 1)
		t.AssertEQ(gbconv.GBTime(s).Hour(), 1)
		t.AssertEQ(gbconv.GBTime(s).Minute(), 2)
		t.AssertEQ(gbconv.GBTime(s).Second(), 3)
		t.AssertEQ(gbconv.GBTime(s), gbtime.NewFromStr(s))
		t.AssertEQ(gbconv.Time(s), gbtime.NewFromStr(s).Time)
	})
	gbtest.C(t, func(t *gbtest.T) {
		t1 := gbtime.NewFromStr("2021-05-21 05:04:51.206547+00")
		t2 := gbconv.GBTime(gbvar.New(t1))
		t3 := gbvar.New(t1).GBTime()
		t.AssertEQ(t1, t2)
		t.AssertEQ(t1.Local(), t2.Local())
		t.AssertEQ(t1, t3)
		t.AssertEQ(t1.Local(), t3.Local())
	})
}

func Test_Time_Slice_Attribute(t *testing.T) {
	type SelectReq struct {
		Arr []*gbtime.Time
		One *gbtime.Time
	}
	gbtest.C(t, func(t *gbtest.T) {
		var s *SelectReq
		err := gbconv.Struct(g.Map{
			"arr": g.Slice{"2021-01-12 12:34:56", "2021-01-12 12:34:57"},
			"one": "2021-01-12 12:34:58",
		}, &s)
		t.AssertNil(err)
		t.Assert(s.One, "2021-01-12 12:34:58")
		t.Assert(s.Arr[0], "2021-01-12 12:34:56")
		t.Assert(s.Arr[1], "2021-01-12 12:34:57")
	})
}

func Test_Issue2901(t *testing.T) {
	type GameApp2 struct {
		ForceUpdateTime *time.Time
	}
	gbtest.C(t, func(t *gbtest.T) {
		src := map[string]interface{}{
			"FORCE_UPDATE_TIME": time.Now(),
		}
		m := GameApp2{}
		err := gbconv.Scan(src, &m)
		t.AssertNil(err)
	})
}
