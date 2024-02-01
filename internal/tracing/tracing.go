// Package tracing provides some utility functions for tracing functionality.
package tracing

import (
	gbtype "github.com/Ghostbb-io/gb/container/gb_type"
	gbbinary "github.com/Ghostbb-io/gb/encoding/gb_binary"
	gbrand "github.com/Ghostbb-io/gb/util/gb_rand"
	"math"
	"time"

	"go.opentelemetry.io/otel/trace"
)

var (
	randomInitSequence = int32(gbrand.Intn(math.MaxInt32))
	sequence           = gbtype.NewInt32(randomInitSequence)
)

// NewIDs creates and returns a new trace and span ID.
func NewIDs() (traceID trace.TraceID, spanID trace.SpanID) {
	return NewTraceID(), NewSpanID()
}

// NewTraceID creates and returns a trace ID.
func NewTraceID() (traceID trace.TraceID) {
	var (
		timestampNanoBytes = gbbinary.EncodeInt64(time.Now().UnixNano())
		sequenceBytes      = gbbinary.EncodeInt32(sequence.Add(1))
		randomBytes        = gbrand.B(4)
	)
	copy(traceID[:], timestampNanoBytes)
	copy(traceID[8:], sequenceBytes)
	copy(traceID[12:], randomBytes)
	return
}

// NewSpanID creates and returns a span ID.
func NewSpanID() (spanID trace.SpanID) {
	copy(spanID[:], gbbinary.EncodeInt64(time.Now().UnixNano()/1e3))
	copy(spanID[4:], gbrand.B(4))
	return
}
