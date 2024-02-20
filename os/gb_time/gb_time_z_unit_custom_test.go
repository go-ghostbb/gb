package gbtime_test

import (
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func Test_Custom1(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbtime.New("2022-03-08T03:01:14-07:00").Local().Time, gbtime.New("2022-03-08T10:01:14Z").Local().Time)
		t.Assert(gbtime.New("2022-03-08T03:01:14-08:00").Local().Time, gbtime.New("2022-03-08T11:01:14Z").Local().Time)
		t.Assert(gbtime.New("2022-03-08T03:01:14-09:00").Local().Time, gbtime.New("2022-03-08T12:01:14Z").Local().Time)
		t.Assert(gbtime.New("2022-03-08T03:01:14+08:00").Local().Time, gbtime.New("2022-03-07T19:01:14Z").Local().Time)
	})
}

func Test_Custom2(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		newTime := gbtime.New("2023-07-26").LayoutTo("2006-01")
		t.Assert(newTime.Year(), 2023)
		t.Assert(newTime.Month(), 7)
		t.Assert(newTime.Day(), 1)
		t.Assert(newTime.Hour(), 0)
		t.Assert(newTime.Minute(), 0)
		t.Assert(newTime.Second(), 0)
	})
}
