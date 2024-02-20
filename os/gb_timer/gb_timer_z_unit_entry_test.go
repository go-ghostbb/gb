// Job Operations

package gbtimer_test

import (
	"context"
	gbarray "ghostbb.io/gb/container/gb_array"
	gbtimer "ghostbb.io/gb/os/gb_timer"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
	"time"
)

func TestJob_Start_Stop_Close(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		job := timer.Add(ctx, 200*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
		job.Stop()
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
		job.Start()
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 2)
		job.Close()
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 2)

		t.Assert(job.Status(), gbtimer.StatusClosed)
	})
}

func TestJob_Singleton(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		job := timer.Add(ctx, 200*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
			time.Sleep(10 * time.Second)
		})
		t.Assert(job.IsSingleton(), false)
		job.SetSingleton(true)
		t.Assert(job.IsSingleton(), true)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestJob_SingletonQuick(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New(gbtimer.TimerOptions{
			Quick: true,
		})
		array := gbarray.New(true)
		job := timer.Add(ctx, 5*time.Second, func(ctx context.Context) {
			array.Append(1)
			time.Sleep(10 * time.Second)
		})
		t.Assert(job.IsSingleton(), false)
		job.SetSingleton(true)
		t.Assert(job.IsSingleton(), true)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestJob_SetTimes(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		job := timer.Add(ctx, 200*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		})
		job.SetTimes(2)
		//job.IsSingleton()
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestJob_Run(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timer := gbtimer.New()
		array := gbarray.New(true)
		job := timer.Add(ctx, 1000*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		})
		job.Job()(ctx)
		t.Assert(array.Len(), 1)
	})
}
