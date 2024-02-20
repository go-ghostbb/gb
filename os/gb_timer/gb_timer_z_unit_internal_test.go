package gbtimer

import (
	"context"
	gbarray "ghostbb.io/gb/container/gb_array"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
	"time"
)

func TestTimer_Proceed(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		array := gbarray.New(true)
		timer := New(TimerOptions{
			Interval: time.Hour,
		})
		timer.Add(ctx, 10000*time.Hour, func(ctx context.Context) {
			array.Append(1)
		})
		timer.proceed(10001)
		time.Sleep(10 * time.Millisecond)
		t.Assert(array.Len(), 1)
		timer.proceed(20001)
		time.Sleep(10 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
	gbtest.C(t, func(t *gbtest.T) {
		array := gbarray.New(true)
		timer := New(TimerOptions{
			Interval: time.Millisecond * 100,
		})
		timer.Add(ctx, 10000*time.Hour, func(ctx context.Context) {
			array.Append(1)
		})
		ticks := int64((10000 * time.Hour) / (time.Millisecond * 100))
		timer.proceed(ticks + 1)
		time.Sleep(10 * time.Millisecond)
		t.Assert(array.Len(), 1)
		timer.proceed(2*ticks + 1)
		time.Sleep(10 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestTimer_PriorityQueue(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		queue := newPriorityQueue()
		queue.Push(1, 1)
		queue.Push(4, 4)
		queue.Push(5, 5)
		queue.Push(2, 2)
		queue.Push(3, 3)
		t.Assert(queue.Pop(), 1)
		t.Assert(queue.Pop(), 2)
		t.Assert(queue.Pop(), 3)
		t.Assert(queue.Pop(), 4)
		t.Assert(queue.Pop(), 5)
	})
}

func TestTimer_PriorityQueue_FirstOneInArrayIsTheLeast(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			size  = 1000000
			array = gbarray.NewIntArrayRange(0, size, 1)
		)
		array.Shuffle()
		queue := newPriorityQueue()
		array.Iterator(func(k int, v int) bool {
			queue.Push(v, int64(v))
			return true
		})
		for i := 0; i < size; i++ {
			t.Assert(queue.Pop(), i)
			t.Assert(queue.heap.array[0].priority, i+1)
		}
	})
}
