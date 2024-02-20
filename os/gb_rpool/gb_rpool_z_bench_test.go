// go test *.go -bench=".*"

package gbrpool_test

import (
	"context"
	gbrpool "ghostbb.io/gb/os/gb_rpool"
	"testing"
)

var (
	ctx = context.TODO()
	n   = 500000
)

func increment(ctx context.Context) {
	for i := 0; i < 1000000; i++ {
	}
}

func BenchmarkGBrpool_1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbrpool.Add(ctx, increment)
	}
}

func BenchmarkGoroutine_1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go increment(ctx)
	}
}

func BenchmarkGBrpool2(b *testing.B) {
	b.N = n
	for i := 0; i < b.N; i++ {
		gbrpool.Add(ctx, increment)
	}
}

func BenchmarkGoroutine2(b *testing.B) {
	b.N = n
	for i := 0; i < b.N; i++ {
		go increment(ctx)
	}
}
