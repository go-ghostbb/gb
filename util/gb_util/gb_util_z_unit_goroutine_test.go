package gbutil_test

import (
	"context"
	gbarray "ghostbb.io/gb/container/gb_array"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbutil "ghostbb.io/gb/util/gb_util"
	"sync"
	"testing"
)

func Test_Go(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			wg    = sync.WaitGroup{}
			array = gbarray.NewArray(true)
		)
		wg.Add(1)
		gbutil.Go(ctx, func(ctx context.Context) {
			defer wg.Done()
			array.Append(1)
		}, nil)
		wg.Wait()
		t.Assert(array.Len(), 1)
	})
	// recover
	gbtest.C(t, func(t *gbtest.T) {
		var (
			wg    = sync.WaitGroup{}
			array = gbarray.NewArray(true)
		)
		wg.Add(1)
		gbutil.Go(ctx, func(ctx context.Context) {
			defer wg.Done()
			panic("error")
			array.Append(1)
		}, nil)
		wg.Wait()
		t.Assert(array.Len(), 0)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var (
			wg    = sync.WaitGroup{}
			array = gbarray.NewArray(true)
		)
		wg.Add(1)
		gbutil.Go(ctx, func(ctx context.Context) {
			panic("error")
		}, func(ctx context.Context, exception error) {
			defer wg.Done()
			array.Append(exception)
		})
		wg.Wait()
		t.Assert(array.Len(), 1)
		t.Assert(array.At(0), "error")
	})
}
