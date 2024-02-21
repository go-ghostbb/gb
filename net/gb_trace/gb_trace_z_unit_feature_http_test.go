package gbtrace_test

import (
	"context"
	"fmt"
	"ghostbb.io/gb/frame/g"
	gbhttp "ghostbb.io/gb/net/gb_http"
	gbtrace "ghostbb.io/gb/net/gb_trace"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbtest "ghostbb.io/gb/test/gb_test"
	"github.com/gin-gonic/gin"
	"testing"
	"time"
)

func Test_Client_Server_Tracing(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		p := 8888
		s := g.Server(p)
		s.SetDumpRouterMap(false)
		s.SetTerminal(false)
		s.SetPort(p)
		s.GET("/", func(c *gin.Context) {
			ctx := gbhttp.Ctx(c)
			g.Log().Print(ctx, "GetTraceID:", gbtrace.GetTraceID(ctx))
			c.String(200, gbtrace.GetTraceID(ctx))
		})
		t.AssertNil(s.Start())
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)
		ctx := gbctx.New()
		prefix := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := g.Client()
		client.SetPrefix(prefix)
		t.Assert(gbtrace.IsUsingDefaultProvider(), true)
		t.Assert(client.GetContent(ctx, "/"), gbtrace.GetTraceID(ctx))
		t.Assert(client.GetContent(ctx, "/"), gbctx.CtxId(ctx))
	})
}

func Test_WithTraceID(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		p := 8889
		s := g.Server(p)
		s.SetDumpRouterMap(false)
		s.SetTerminal(false)
		s.SetPort(p)
		s.GET("/", func(c *gin.Context) {
			ctx := gbhttp.Ctx(c)
			g.Log().Print(ctx, "GetTraceID:", gbtrace.GetTraceID(ctx))
			c.String(200, gbtrace.GetTraceID(ctx))
		})
		t.AssertNil(s.Start())
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		ctx, err := gbtrace.WithTraceID(context.TODO(), traceID.String())
		t.AssertNil(err)

		prefix := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := g.Client()
		client.SetPrefix(prefix)
		t.Assert(gbtrace.IsUsingDefaultProvider(), true)
		t.Assert(client.GetContent(ctx, "/"), gbtrace.GetTraceID(ctx))
		t.Assert(client.GetContent(ctx, "/"), traceIDStr)
	})
}
