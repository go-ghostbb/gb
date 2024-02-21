package gbhttp

import (
	"compress/gzip"
	"context"
	"fmt"
	"ghostbb.io/gb"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/httputil"
	"ghostbb.io/gb/internal/utils"
	gbtrace "ghostbb.io/gb/net/gb_trace"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
	"strings"
)

const (
	tracingInstrumentName                        = "ghostbb.io/gb/net/gbhttp.Server"
	tracingEventHttpRequest                      = "http.request"
	tracingEventHttpRequestHeaders               = "http.request.headers"
	tracingEventHttpRequestBaggage               = "http.request.baggage"
	tracingEventHttpRequestBody                  = "http.request.body"
	tracingEventHttpResponse                     = "http.response"
	tracingEventHttpResponseHeaders              = "http.response.headers"
	tracingEventHttpResponseBody                 = "http.response.body"
	tracingEventHttpRequestUrl                   = "http.request.url"
	tracingMiddlewareHandled        gbctx.StrKey = `MiddlewareServerTracingHandled`
)

func (s *Server) traceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			ctx = c.Request.Context()
		)
		// Mark this request is handled by server tracing middleware,
		// to avoid repeated handling by the same middleware.
		if ctx.Value(tracingMiddlewareHandled) != nil {
			c.Next()
			return
		}

		ctx = context.WithValue(ctx, tracingMiddlewareHandled, 1)
		var (
			span trace.Span
			tr   = otel.GetTracerProvider().Tracer(
				tracingInstrumentName,
				trace.WithInstrumentationVersion(gb.VERSION),
			)
		)

		ctx, span = tr.Start(
			otel.GetTextMapPropagator().Extract(
				ctx,
				propagation.HeaderCarrier(c.Request.Header),
			),
			c.Request.URL.Path,
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		span.SetAttributes(gbtrace.CommonLabels()...)
		*c.Request = *c.Request.WithContext(ctx)

		// If it is now using a default trace provider, it then does no complex tracing jobs.
		if gbtrace.IsUsingDefaultProvider() {
			c.Next()
			return
		}

		// Request content logging.
		reqBodyContentBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			_ = c.Error(gberror.Wrap(err, `read request body failed`))
			span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
			return
		}
		c.Request.Body = utils.NewReadCloser(reqBodyContentBytes, false)

		span.AddEvent(tracingEventHttpRequest, trace.WithAttributes(
			attribute.String(tracingEventHttpRequestUrl, c.Request.URL.String()),
			attribute.String(tracingEventHttpRequestHeaders, gbconv.String(httputil.HeaderToMap(c.Request.Header))),
			attribute.String(tracingEventHttpRequestBaggage, gbtrace.GetBaggageMap(ctx).String()),
			attribute.String(tracingEventHttpRequestBody, gbstr.StrLimit(
				string(reqBodyContentBytes),
				gbtrace.MaxContentLogSize(),
				"...",
			)),
		))

		// Continue executing.
		c.Next()

		// Error logging.
		if err = c.Err(); err != nil {
			span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
		}
		// Response content logging.
		bodyByte := make([]byte, 0)
		_, _ = c.Request.Response.Body.Read(bodyByte)
		var resBodyContent = gbstr.StrLimit(string(bodyByte), gbtrace.MaxContentLogSize(), "...")
		if gzipAccepted(c.Request.Response.Header) {
			reader, err := gzip.NewReader(c.Request.Response.Body)
			if err != nil {
				span.SetStatus(codes.Error, fmt.Sprintf(`read gzip response err:%+v`, err))
			}
			defer reader.Close()
			uncompressed, err := io.ReadAll(reader)
			if err != nil {
				span.SetStatus(codes.Error, fmt.Sprintf(`get uncompress value err:%+v`, err))
			}
			resBodyContent = gbstr.StrLimit(string(uncompressed), gbtrace.MaxContentLogSize(), "...")
		}

		span.AddEvent(tracingEventHttpResponse, trace.WithAttributes(
			attribute.String(tracingEventHttpResponseHeaders, gbconv.String(httputil.HeaderToMap(c.Request.Response.Header))),
			attribute.String(tracingEventHttpResponseBody, resBodyContent),
		))
	}
}

// gzipAccepted returns whether the client will accept gzip-encoded content.
func gzipAccepted(header http.Header) bool {
	a := header.Get("Content-Encoding")
	parts := strings.Split(a, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "gzip" || strings.HasPrefix(part, "gzip;") {
			return true
		}
	}
	return false
}
