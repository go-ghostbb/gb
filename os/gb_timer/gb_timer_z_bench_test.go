package gbtimer

import (
	"context"
	"testing"
	"time"
)

var (
	ctx   = context.TODO()
	timer = New()
)

func Benchmark_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		timer.Add(ctx, time.Hour, func(ctx context.Context) {

		})
	}
}

func Benchmark_PriorityQueue_Pop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		timer.queue.Pop()
	}
}

func Benchmark_StartStop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		timer.Start()
		timer.Stop()
	}
}
