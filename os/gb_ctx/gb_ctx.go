// Package gbctx wraps context.Context and provides extra context features.
package gbctx

import (
	"context"
	gbtrace "github.com/Ghostbb-io/gb/net/gb_trace"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type (
	Ctx    = context.Context // Ctx is short name alias for context.Context.
	StrKey string            // StrKey is a type for warps basic type string as context key.
)

var (
	// initCtx is the context initialized from process environment.
	initCtx context.Context
)

func init() {
	// All environment key-value pairs.
	m := make(map[string]string)
	i := 0
	for _, s := range os.Environ() {
		i = strings.IndexByte(s, '=')
		if i == -1 {
			continue
		}
		m[s[0:i]] = s[i+1:]
	}
	// OpenTelemetry from environments.
	initCtx = otel.GetTextMapPropagator().Extract(
		context.Background(),
		propagation.MapCarrier(m),
	)
	initCtx = WithCtx(initCtx)
}

// New creates and returns a context which contains context id.
func New() context.Context {
	return WithCtx(context.Background())
}

// WithCtx creates and returns a context containing context id upon given parent context `ctx`.
func WithCtx(ctx context.Context) context.Context {
	if CtxId(ctx) != "" {
		return ctx
	}
	var span *gbtrace.Span
	ctx, span = gbtrace.NewSpan(ctx, "gbctx.WithCtx")
	defer span.End()
	return ctx
}

// CtxId retrieves and returns the context id from context.
func CtxId(ctx context.Context) string {
	return gbtrace.GetTraceID(ctx)
}

// SetInitCtx sets custom initialization context.
// Note that this function cannot be called in multiple goroutines.
func SetInitCtx(ctx context.Context) {
	initCtx = ctx
}

// GetInitCtx returns the initialization context.
// Initialization context is used in `main` or `init` functions.
func GetInitCtx() context.Context {
	return initCtx
}
