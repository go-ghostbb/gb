package gbcron_test

import (
	"context"
	gbarray "ghostbb.io/gb/container/gb_array"
	gbcron "ghostbb.io/gb/os/gb_cron"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
	"time"
)

func TestCron_Entry_Operations(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			cron  = gbcron.New()
			array = gbarray.New(true)
		)
		cron.DelayAddTimes(ctx, 500*time.Millisecond, "* * * * * *", 2, func(ctx context.Context) {
			array.Append(1)
		})
		t.Assert(cron.Size(), 0)
		time.Sleep(800 * time.Millisecond)
		t.Assert(array.Len(), 0)
		t.Assert(cron.Size(), 1)
		time.Sleep(3000 * time.Millisecond)
		t.Assert(array.Len(), 2)
		t.Assert(cron.Size(), 0)
	})

	gbtest.C(t, func(t *gbtest.T) {
		var (
			cron  = gbcron.New()
			array = gbarray.New(true)
		)
		entry, err1 := cron.Add(ctx, "* * * * * *", func(ctx context.Context) {
			array.Append(1)
		})
		t.Assert(err1, nil)
		t.Assert(array.Len(), 0)
		t.Assert(cron.Size(), 1)
		time.Sleep(1300 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 1)
		entry.Stop()
		time.Sleep(5000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 1)
		entry.Start()
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
		t.Assert(cron.Size(), 1)
		entry.Close()
		time.Sleep(1200 * time.Millisecond)
		t.Assert(cron.Size(), 0)
	})
}
