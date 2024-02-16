package gbmutex_test

import (
	"context"
	gbarray "ghostbb.io/gb/container/gb_array"
	gblog "ghostbb.io/gb/os/gb_log"
	gbmutex "ghostbb.io/gb/os/gb_mutex"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
	"time"
)

func Test_RWMutex_RUnlock(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		mu := gbmutex.RWMutex{}
		mu.RLockFunc(func() {
			time.Sleep(200 * time.Millisecond)
		})
	})

	// RLock before Lock
	gbtest.C(t, func(t *gbtest.T) {
		mu := gbmutex.RWMutex{}
		mu.RLock()
		go func() {
			mu.Lock()
			time.Sleep(300 * time.Millisecond)
			mu.Unlock()
		}()
		time.Sleep(100 * time.Millisecond)
		mu.RUnlock()
	})
}

func Test_RWMutex_IsLocked(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		mu := gbmutex.RWMutex{}
		go func() {
			mu.LockFunc(func() {
				time.Sleep(200 * time.Millisecond)
			})
		}()
		time.Sleep(100 * time.Millisecond)

		go func() {
			mu.RLockFunc(func() {
				time.Sleep(200 * time.Millisecond)
			})
		}()
	})
}

func Test_RWMutex_Unlock(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		mu := gbmutex.RWMutex{}
		array := gbarray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.LockFunc(func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.LockFunc(func() {
				array.Append(1)
			})
		}()

		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(400 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func Test_RWMutex_LockFunc(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		mu := gbmutex.RWMutex{}
		array := gbarray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.LockFunc(func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func Test_RWMutex_TryLockFunc(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		mu := gbmutex.RWMutex{}
		array := gbarray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.TryLockFunc(func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(400 * time.Millisecond)
			mu.TryLockFunc(func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func Test_RWMutex_RLockFunc(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		mu := gbmutex.RWMutex{}
		array := gbarray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

	gbtest.C(t, func(t *gbtest.T) {
		mu := gbmutex.RWMutex{}
		array := gbarray.New(true)
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.RLockFunc(func() {
				array.Append(1)
				time.Sleep(100 * time.Millisecond)
			})
		}()
		t.Assert(array.Len(), 0)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func Test_RWMutex_TryRLockFunc(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			mu    = gbmutex.RWMutex{}
			array = gbarray.New(true)
		)
		// First writing lock
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				gblog.Print(context.TODO(), "lock1 done")
				time.Sleep(2000 * time.Millisecond)
			})
		}()
		// This goroutine never gets the lock.
		go func() {
			time.Sleep(1000 * time.Millisecond)
			mu.TryRLockFunc(func() {
				array.Append(1)
			})
		}()
		for index := 0; index < 1000; index++ {
			go func() {
				time.Sleep(4000 * time.Millisecond)
				mu.TryRLockFunc(func() {
					array.Append(1)
				})
			}()
		}
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(2000 * time.Millisecond)
		t.Assert(array.Len(), 1001)
	})
}
