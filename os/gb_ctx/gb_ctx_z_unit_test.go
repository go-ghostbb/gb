package gbctx_test

import (
	"context"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func Test_New(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		ctx := gbctx.New()
		t.AssertNE(ctx, nil)
		t.AssertNE(gbctx.CtxId(ctx), "")
	})
}

func Test_WithCtx(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		ctx := context.WithValue(context.TODO(), "TEST", 1)
		ctx = gbctx.WithCtx(ctx)
		t.AssertNE(gbctx.CtxId(ctx), "")
		t.Assert(ctx.Value("TEST"), 1)
	})
}

func Test_SetInitCtx(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		ctx := context.WithValue(context.TODO(), "TEST", 1)
		gbctx.SetInitCtx(ctx)
		t.AssertNE(gbctx.GetInitCtx(), "")
		t.Assert(gbctx.GetInitCtx().Value("TEST"), 1)
	})
}
