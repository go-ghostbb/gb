package gbmlock_test

import (
	gbmlock "ghostbb.io/gb/os/gb_mlock"
	"testing"
)

var (
	lockKey = "This is the lock key for gbmlock."
)

func Benchmark_GMLock_Lock_Unlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbmlock.Lock(lockKey)
		gbmlock.Unlock(lockKey)
	}
}

func Benchmark_GMLock_RLock_RUnlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbmlock.RLock(lockKey)
		gbmlock.RUnlock(lockKey)
	}
}

func Benchmark_GMLock_TryLock_Unlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if gbmlock.TryLock(lockKey) {
			gbmlock.Unlock(lockKey)
		}
	}
}

func Benchmark_GMLock_TryRLock_RUnlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if gbmlock.TryRLock(lockKey) {
			gbmlock.RUnlock(lockKey)
		}
	}
}
