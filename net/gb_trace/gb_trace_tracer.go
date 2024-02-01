package gbtrace

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Tracer warps trace.Tracer for compatibility and extension.
type Tracer struct {
	trace.Tracer
}

// NewTracer Tracer is a short function for retrieving Tracer.
func NewTracer(name ...string) *Tracer {
	tracerName := ""
	if len(name) > 0 {
		tracerName = name[0]
	}
	return &Tracer{
		Tracer: otel.Tracer(tracerName),
	}
}
