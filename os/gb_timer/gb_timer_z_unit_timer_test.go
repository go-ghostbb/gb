// Timer Operations

package gbtimer_test

import (
	"context"
	gbarray "ghostbb.io/gb/container/gb_array"
	gbtimer "ghostbb.io/gb/os/gb_timer"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
	"time"
)

func TestTimer_Add_Close(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		//fmt.Println("start", time.Now())
		timer.Add(ctx, 200*time.Millisecond, func(ctx context.Context) {
			//fmt.Println("job1", time.Now())
			array.Append(1)
		})
		timer.Add(ctx, 200*time.Millisecond, func(ctx context.Context) {
			//fmt.Println("job2", time.Now())
			array.Append(1)
		})
		timer.Add(ctx, 400*time.Millisecond, func(ctx context.Context) {
			//fmt.Println("job3", time.Now())
			array.Append(1)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 2)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 5)
		timer.Close()
		time.Sleep(250 * time.Millisecond)
		fixedLength := array.Len()
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), fixedLength)
	})
}

func TestTimer_Start_Stop_Close(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		timer.Add(ctx, 1000*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		})
		t.Assert(array.Len(), 0)
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 1)
		timer.Stop()
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 1)
		timer.Start()
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 2)
		timer.Close()
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestJob_Reset(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		job := timer.AddSingleton(ctx, 500*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(300 * time.Millisecond)
		job.Reset()
		time.Sleep(300 * time.Millisecond)
		job.Reset()
		time.Sleep(300 * time.Millisecond)
		job.Reset()
		time.Sleep(600 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_AddSingleton(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		timer.AddSingleton(ctx, 200*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
			time.Sleep(10 * time.Second)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(500 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_AddSingletonWithQuick(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New(gbtimer.TimerOptions{
			Interval: 100 * time.Millisecond,
			Quick:    true,
		})
		array := gbarray.New(true)
		timer.AddSingleton(ctx, 5*time.Second, func(ctx context.Context) {
			array.Append(1)
			time.Sleep(10 * time.Second)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(500 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_AddSingletonWithoutQuick(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New(gbtimer.TimerOptions{
			Interval: 100 * time.Millisecond,
			Quick:    false,
		})
		array := gbarray.New(true)
		timer.AddSingleton(ctx, 5*time.Second, func(ctx context.Context) {
			array.Append(1)
			time.Sleep(10 * time.Second)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 0)

		time.Sleep(500 * time.Millisecond)
		t.Assert(array.Len(), 0)
	})
}

func TestTimer_AddOnce(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		timer.AddOnce(ctx, 200*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		})
		timer.AddOnce(ctx, 200*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 2)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 2)
		timer.Close()
		time.Sleep(250 * time.Millisecond)
		fixedLength := array.Len()
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), fixedLength)
	})
}

func TestTimer_AddTimes(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		timer.AddTimes(ctx, 200*time.Millisecond, 2, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestTimer_DelayAdd(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		timer.DelayAdd(ctx, 200*time.Millisecond, 200*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_DelayAddJob(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		timer.DelayAddEntry(ctx, 200*time.Millisecond, 200*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		}, false, 100, gbtimer.StatusReady)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_DelayAddSingleton(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		timer.DelayAddSingleton(ctx, 200*time.Millisecond, 200*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
			time.Sleep(10 * time.Second)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 0)

		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_DelayAddOnce(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		timer.DelayAddOnce(ctx, 200*time.Millisecond, 200*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 0)

		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(500 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_DelayAddTimes(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		timer.DelayAddTimes(ctx, 200*time.Millisecond, 500*time.Millisecond, 2, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 0)

		time.Sleep(600 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(600 * time.Millisecond)
		t.Assert(array.Len(), 2)

		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestTimer_AddLessThanInterval(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New(gbtimer.TimerOptions{
			Interval: 100 * time.Millisecond,
		})
		array := gbarray.New(true)
		timer.Add(ctx, 20*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(50 * time.Millisecond)
		t.Assert(array.Len(), 0)

		time.Sleep(110 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(110 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestTimer_AddLeveledJob1(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		timer.DelayAdd(ctx, 1000*time.Millisecond, 1000*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(1500 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(1300 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestTimer_Exit(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		timer.Add(ctx, 200*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
			gbtimer.Exit()
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}
