package gblog_test

import (
	"bytes"
	"context"
	gbarray "ghostbb.io/gb/container/gb_array"
	gblog "ghostbb.io/gb/os/gb_log"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"testing"
)

var arrayForHandlerTest1 = gbarray.NewStrArray()

func customHandler1(ctx context.Context, input *gblog.HandlerInput) {
	arrayForHandlerTest1.Append(input.String(false))
	input.Next(ctx)
}

func TestLogger_SetHandlers1(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		w := bytes.NewBuffer(nil)
		l := gblog.NewWithWriter(w)
		l.SetHandlers(customHandler1)
		l.SetCtxKeys("Trace-Id", "Span-Id", "Test")
		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Print(ctx, 1, 2, 3)
		t.Assert(gbstr.Count(w.String(), "1234567890"), 1)
		t.Assert(gbstr.Count(w.String(), "abcdefg"), 1)
		t.Assert(gbstr.Count(w.String(), "1 2 3"), 1)

		t.Assert(arrayForHandlerTest1.Len(), 1)
		t.Assert(gbstr.Count(arrayForHandlerTest1.At(0), "1234567890"), 1)
		t.Assert(gbstr.Count(arrayForHandlerTest1.At(0), "abcdefg"), 1)
		t.Assert(gbstr.Count(arrayForHandlerTest1.At(0), "1 2 3"), 1)
	})
}

var arrayForHandlerTest2 = gbarray.NewStrArray()

func customHandler2(ctx context.Context, input *gblog.HandlerInput) {
	arrayForHandlerTest2.Append(input.String(false))
}

func TestLogger_SetHandlers2(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		w := bytes.NewBuffer(nil)
		l := gblog.NewWithWriter(w)
		l.SetHandlers(customHandler2)
		l.SetCtxKeys("Trace-Id", "Span-Id", "Test")
		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Print(ctx, 1, 2, 3)
		t.Assert(gbstr.Count(w.String(), "1234567890"), 0)
		t.Assert(gbstr.Count(w.String(), "abcdefg"), 0)
		t.Assert(gbstr.Count(w.String(), "1 2 3"), 0)

		t.Assert(arrayForHandlerTest2.Len(), 1)
		t.Assert(gbstr.Count(arrayForHandlerTest2.At(0), "1234567890"), 1)
		t.Assert(gbstr.Count(arrayForHandlerTest2.At(0), "abcdefg"), 1)
		t.Assert(gbstr.Count(arrayForHandlerTest2.At(0), "1 2 3"), 1)
	})
}

func TestLogger_SetHandlers_HandlerJson(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		w := bytes.NewBuffer(nil)
		l := gblog.NewWithWriter(w)
		l.SetHandlers(gblog.HandlerJson)
		l.SetCtxKeys("Trace-Id", "Span-Id", "Test")
		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Debug(ctx, 1, 2, 3)
		t.Assert(gbstr.Count(w.String(), "1234567890"), 1)
		t.Assert(gbstr.Count(w.String(), "abcdefg"), 1)
		t.Assert(gbstr.Count(w.String(), `"1 2 3"`), 1)
		t.Assert(gbstr.Count(w.String(), `"DEBU"`), 1)
	})
}

func TestLogger_SetHandlers_HandlerStructure(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		w := bytes.NewBuffer(nil)
		l := gblog.NewWithWriter(w)
		l.SetHandlers(gblog.HandlerStructure)
		l.SetCtxKeys("Trace-Id", "Span-Id", "Test")
		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Debug(ctx, "debug", "uid", 1000)
		l.Info(ctx, "info", "' '", `"\n`)

		t.Assert(gbstr.Count(w.String(), "uid=1000"), 1)
		t.Assert(gbstr.Count(w.String(), "Content=debug"), 1)
		t.Assert(gbstr.Count(w.String(), `"' '"="\"\\n"`), 1)
		t.Assert(gbstr.Count(w.String(), `CtxStr="1234567890, abcdefg"`), 2)
	})
}

func Test_SetDefaultHandler(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		oldHandler := gblog.GetDefaultHandler()
		gblog.SetDefaultHandler(func(ctx context.Context, in *gblog.HandlerInput) {
			gblog.HandlerJson(ctx, in)
		})
		defer gblog.SetDefaultHandler(oldHandler)

		w := bytes.NewBuffer(nil)
		l := gblog.NewWithWriter(w)
		l.SetCtxKeys("Trace-Id", "Span-Id", "Test")
		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Debug(ctx, 1, 2, 3)
		t.Assert(gbstr.Count(w.String(), "1234567890"), 1)
		t.Assert(gbstr.Count(w.String(), "abcdefg"), 1)
		t.Assert(gbstr.Count(w.String(), `"1 2 3"`), 1)
		t.Assert(gbstr.Count(w.String(), `"DEBU"`), 1)
	})
}
