package gbrpool_test

import (
	"context"
	gbarray "ghostbb.io/gb/container/gb_array"
	gbrpool "ghostbb.io/gb/os/gb_rpool"
	gbtest "ghostbb.io/gb/test/gb_test"
	"sync"
	"testing"
	"time"
)

func Test_Basic(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			err   error
			wg    = sync.WaitGroup{}
			array = gbarray.NewArray(true)
			size  = 100
		)
		wg.Add(size)
		for i := 0; i < size; i++ {
			err = gbrpool.Add(ctx, func(ctx context.Context) {
				array.Append(1)
				wg.Done()
			})
			t.AssertNil(err)
		}
		wg.Wait()

		time.Sleep(100 * time.Millisecond)

		t.Assert(array.Len(), size)
		t.Assert(gbrpool.Jobs(), 0)
		t.Assert(gbrpool.Size(), 0)
	})
}

func Test_Limit1(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			wg    = sync.WaitGroup{}
			array = gbarray.NewArray(true)
			size  = 100
			pool  = gbrpool.New(10)
		)
		wg.Add(size)
		for i := 0; i < size; i++ {
			pool.Add(ctx, func(ctx context.Context) {
				array.Append(1)
				wg.Done()
			})
		}
		wg.Wait()
		t.Assert(array.Len(), size)
	})
}

func Test_Limit2(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			err   error
			wg    = sync.WaitGroup{}
			array = gbarray.NewArray(true)
			size  = 100
			pool  = gbrpool.New(1)
		)
		wg.Add(size)
		for i := 0; i < size; i++ {
			err = pool.Add(ctx, func(ctx context.Context) {
				defer wg.Done()
				array.Append(1)
			})
			t.AssertNil(err)
		}
		wg.Wait()
		t.Assert(array.Len(), size)
	})
}

func Test_Limit3(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			array = gbarray.NewArray(true)
			size  = 1000
			pool  = gbrpool.New(100)
		)
		t.Assert(pool.Cap(), 100)
		for i := 0; i < size; i++ {
			pool.Add(ctx, func(ctx context.Context) {
				array.Append(1)
				time.Sleep(2 * time.Second)
			})
		}
		time.Sleep(time.Second)
		t.Assert(pool.Size(), 100)
		t.Assert(pool.Jobs(), 900)
		t.Assert(array.Len(), 100)
		pool.Close()
		time.Sleep(2 * time.Second)
		t.Assert(pool.Size(), 0)
		t.Assert(pool.Jobs(), 900)
		t.Assert(array.Len(), 100)
		t.Assert(pool.IsClosed(), true)
		t.AssertNE(pool.Add(ctx, func(ctx context.Context) {}), nil)
	})
}

func Test_AddWithRecover(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			err   error
			array = gbarray.NewArray(true)
		)
		err = gbrpool.AddWithRecover(ctx, func(ctx context.Context) {
			array.Append(1)
			panic(1)
		}, func(ctx context.Context, err error) {
			array.Append(1)
		})
		t.AssertNil(err)
		err = gbrpool.AddWithRecover(ctx, func(ctx context.Context) {
			panic(1)
			array.Append(1)
		}, nil)
		t.AssertNil(err)

		time.Sleep(500 * time.Millisecond)

		t.Assert(array.Len(), 2)
	})
}
