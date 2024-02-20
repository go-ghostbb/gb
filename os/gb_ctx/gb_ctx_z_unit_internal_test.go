package gbctx_test

import (
	"context"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
	"time"
)

func Test_NeverDone(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		ctx, _ := context.WithDeadline(gbctx.New(), time.Now().Add(time.Hour))
		t.AssertNE(ctx, nil)
		t.AssertNE(ctx.Done(), nil)
		t.Assert(ctx.Err(), nil)

		tm, ok := ctx.Deadline()
		t.AssertNE(tm, time.Time{})
		t.Assert(ok, true)

		ctx = gbctx.NeverDone(ctx)
		t.AssertNE(ctx, nil)
		t.Assert(ctx.Done(), nil)
		t.Assert(ctx.Err(), nil)

		tm, ok = ctx.Deadline()
		t.Assert(tm, time.Time{})
		t.Assert(ok, false)
	})
}
