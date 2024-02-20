package gbcron_test

import (
	"context"
	gbcron "ghostbb.io/gb/os/gb_cron"
	"testing"
)

func Benchmark_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbcron.Add(ctx, "1 1 1 1 1 1", func(ctx context.Context) {

		})
	}
}
