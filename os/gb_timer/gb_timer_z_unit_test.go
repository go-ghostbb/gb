// Package functions

package gbtimer_test

import (
	"context"
	gbarray "ghostbb.io/gb/container/gb_array"
	gbtimer "ghostbb.io/gb/os/gb_timer"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
	"time"
)

var (
	ctx = context.TODO()
)

func TestSetTimeout(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		array := gbarray.New(true)
		gbtimer.SetTimeout(ctx, 200*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestSetInterval(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		array := gbarray.New(true)
		gbtimer.SetInterval(ctx, 300*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func TestAddEntry(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		array := gbarray.New(true)
		gbtimer.AddEntry(ctx, 200*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		}, false, 2, gbtimer.StatusReady)
		time.Sleep(1100 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestAddSingleton(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		array := gbarray.New(true)
		gbtimer.AddSingleton(ctx, 200*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
			time.Sleep(10000 * time.Millisecond)
		})
		time.Sleep(1100 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestAddTimes(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		array := gbarray.New(true)
		gbtimer.AddTimes(ctx, 200*time.Millisecond, 2, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestDelayAdd(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		array := gbarray.New(true)
		gbtimer.DelayAdd(ctx, 500*time.Millisecond, 500*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(600 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(600 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestDelayAddEntry(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		array := gbarray.New(true)
		gbtimer.DelayAddEntry(ctx, 500*time.Millisecond, 500*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		}, false, 2, gbtimer.StatusReady)
		time.Sleep(500 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(2000 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestDelayAddSingleton(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		array := gbarray.New(true)
		gbtimer.DelayAddSingleton(ctx, 500*time.Millisecond, 500*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
			time.Sleep(10000 * time.Millisecond)
		})
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestDelayAddOnce(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		array := gbarray.New(true)
		gbtimer.DelayAddOnce(ctx, 1000*time.Millisecond, 2000*time.Millisecond, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(2000 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(2000 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestDelayAddTimes(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		array := gbarray.New(true)
		gbtimer.DelayAddTimes(ctx, 500*time.Millisecond, 500*time.Millisecond, 2, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(1500 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}
