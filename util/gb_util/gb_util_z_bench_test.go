// go test *.go -bench=".*" -benchmem

package gbutil

import (
	"context"
	"testing"
)

var (
	m1 = map[string]interface{}{
		"k1": "v1",
	}
	m2 = map[string]interface{}{
		"k2": "v2",
	}
)

func Benchmark_TryCatch(b *testing.B) {
	ctx := context.TODO()
	for i := 0; i < b.N; i++ {
		TryCatch(ctx, func(ctx context.Context) {

		}, func(ctx context.Context, err error) {

		})
	}
}

func Benchmark_MapMergeCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MapMergeCopy(m1, m2)
	}
}
