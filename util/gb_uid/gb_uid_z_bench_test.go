package gbuid_test

import (
	gbuid "ghostbb.io/gb/util/gb_uid"
	"testing"
)

func Benchmark_S(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gbuid.S()
		}
	})
}

func Benchmark_S_Data_1(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gbuid.S([]byte("123"))
		}
	})
}

func Benchmark_S_Data_2(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gbuid.S([]byte("123"), []byte("456"))
		}
	})
}
