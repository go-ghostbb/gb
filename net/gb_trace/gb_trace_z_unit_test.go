package gbtrace_test

import (
	"context"
	gbtrace "ghostbb.io/gb/net/gb_trace"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"testing"
)

func TestWithTraceID(t *testing.T) {
	var (
		ctx  = context.Background()
		uuid = `a323f910-f690-11ec-963d-79c0b7fcf119`
	)
	gbtest.C(t, func(t *gbtest.T) {
		newCtx, err := gbtrace.WithTraceID(ctx, uuid)
		t.AssertNE(err, nil)
		t.Assert(newCtx, ctx)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var traceId = gbstr.Replace(uuid, "-", "")
		newCtx, err := gbtrace.WithTraceID(ctx, traceId)
		t.AssertNil(err)
		t.AssertNE(newCtx, ctx)
		t.Assert(gbtrace.GetTraceID(ctx), "")
		t.Assert(gbtrace.GetTraceID(newCtx), traceId)
	})
}

func TestWithUUID(t *testing.T) {
	var (
		ctx  = context.Background()
		uuid = `a323f910-f690-11ec-963d-79c0b7fcf119`
	)
	gbtest.C(t, func(t *gbtest.T) {
		newCtx, err := gbtrace.WithTraceID(ctx, uuid)
		t.AssertNE(err, nil)
		t.Assert(newCtx, ctx)
	})
	gbtest.C(t, func(t *gbtest.T) {
		newCtx, err := gbtrace.WithUUID(ctx, uuid)
		t.AssertNil(err)
		t.AssertNE(newCtx, ctx)
		t.Assert(gbtrace.GetTraceID(ctx), "")
		t.Assert(gbtrace.GetTraceID(newCtx), gbstr.Replace(uuid, "-", ""))
	})
}
