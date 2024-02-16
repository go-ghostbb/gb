package gbmlock_test

import (
	gbarray "ghostbb.io/gb/container/gb_array"
	gbmlock "ghostbb.io/gb/os/gb_mlock"
	gbtest "ghostbb.io/gb/test/gb_test"
	"sync"
	"testing"
	"time"
)

func Test_Locker_Lock(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		key := "testLock"
		array := gbarray.New(true)
		go func() {
			gbmlock.Lock(key)
			array.Append(1)
			time.Sleep(300 * time.Millisecond)
			gbmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gbmlock.Lock(key)
			array.Append(1)
			gbmlock.Unlock(key)
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
		gbmlock.Remove(key)
	})

	gbtest.C(t, func(t *gbtest.T) {
		key := "testLock"
		array := gbarray.New(true)
		lock := gbmlock.New()
		go func() {
			lock.Lock(key)
			array.Append(1)
			time.Sleep(300 * time.Millisecond)
			lock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			lock.Lock(key)
			array.Append(1)
			lock.Unlock(key)
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
		lock.Clear()
	})

}

func Test_Locker_TryLock(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		key := "testTryLock"
		array := gbarray.New(true)
		go func() {
			gbmlock.Lock(key)
			array.Append(1)
			time.Sleep(300 * time.Millisecond)
			gbmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(150 * time.Millisecond)
			if gbmlock.TryLock(key) {
				array.Append(1)
				gbmlock.Unlock(key)
			}
		}()
		go func() {
			time.Sleep(400 * time.Millisecond)
			if gbmlock.TryLock(key) {
				array.Append(1)
				gbmlock.Unlock(key)
			}
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

}

func Test_Locker_LockFunc(t *testing.T) {
	//no expire
	gbtest.C(t, func(t *gbtest.T) {
		key := "testLockFunc"
		array := gbarray.New(true)
		go func() {
			gbmlock.LockFunc(key, func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			}) //
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gbmlock.LockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1) //
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func Test_Locker_TryLockFunc(t *testing.T) {
	//no expire
	gbtest.C(t, func(t *gbtest.T) {
		key := "testTryLockFunc"
		array := gbarray.New(true)
		go func() {
			gbmlock.TryLockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gbmlock.TryLockFunc(key, func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(300 * time.Millisecond)
			gbmlock.TryLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(150 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(400 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func Test_Multiple_Goroutine(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		ch := make(chan struct{})
		num := 1000
		wait := sync.WaitGroup{}
		wait.Add(num)
		for i := 0; i < num; i++ {
			go func() {
				defer wait.Done()
				<-ch
				gbmlock.Lock("test")
				defer gbmlock.Unlock("test")
				time.Sleep(time.Millisecond)
			}()
		}
		close(ch)
		wait.Wait()
	})

	gbtest.C(t, func(t *gbtest.T) {
		ch := make(chan struct{})
		num := 100
		wait := sync.WaitGroup{}
		wait.Add(num * 2)
		for i := 0; i < num; i++ {
			go func() {
				defer wait.Done()
				<-ch
				gbmlock.Lock("test")
				defer gbmlock.Unlock("test")
				time.Sleep(time.Millisecond)
			}()
		}
		for i := 0; i < num; i++ {
			go func() {
				defer wait.Done()
				<-ch
				gbmlock.RLock("test")
				defer gbmlock.RUnlock("test")
				time.Sleep(time.Millisecond)
			}()
		}
		close(ch)
		wait.Wait()
	})
}

func Test_Locker_RLock(t *testing.T) {
	// RLock before Lock
	gbtest.C(t, func(t *gbtest.T) {
		key := "testRLockBeforeLock"
		array := gbarray.New(true)
		go func() {
			gbmlock.RLock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gbmlock.RUnlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gbmlock.Lock(key)
			array.Append(1)
			gbmlock.Unlock(key)
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

	// Lock before RLock
	gbtest.C(t, func(t *gbtest.T) {
		key := "testLockBeforeRLock"
		array := gbarray.New(true)
		go func() {
			gbmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gbmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gbmlock.RLock(key)
			array.Append(1)
			gbmlock.RUnlock(key)
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

	// Lock before RLocks
	gbtest.C(t, func(t *gbtest.T) {
		key := "testLockBeforeRLocks"
		array := gbarray.New(true)
		go func() {
			gbmlock.Lock(key)
			array.Append(1)
			time.Sleep(300 * time.Millisecond)
			gbmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gbmlock.RLock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gbmlock.RUnlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gbmlock.RLock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gbmlock.RUnlock(key)
		}()
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func Test_Locker_TryRLock(t *testing.T) {
	// Lock before TryRLock
	gbtest.C(t, func(t *gbtest.T) {
		key := "testLockBeforeTryRLock"
		array := gbarray.New(true)
		go func() {
			gbmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gbmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			if gbmlock.TryRLock(key) {
				array.Append(1)
				gbmlock.RUnlock(key)
			}
		}()
		time.Sleep(150 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})

	// Lock before TryRLocks
	gbtest.C(t, func(t *gbtest.T) {
		key := "testLockBeforeTryRLocks"
		array := gbarray.New(true)
		go func() {
			gbmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gbmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			if gbmlock.TryRLock(key) {
				array.Append(1)
				gbmlock.RUnlock(key)
			}
		}()
		go func() {
			time.Sleep(300 * time.Millisecond)
			if gbmlock.TryRLock(key) {
				array.Append(1)
				gbmlock.RUnlock(key)
			}
		}()
		time.Sleep(150 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func Test_Locker_RLockFunc(t *testing.T) {
	// RLockFunc before Lock
	gbtest.C(t, func(t *gbtest.T) {
		key := "testRLockFuncBeforeLock"
		array := gbarray.New(true)
		go func() {
			gbmlock.RLockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gbmlock.Lock(key)
			array.Append(1)
			gbmlock.Unlock(key)
		}()
		time.Sleep(150 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

	// Lock before RLockFunc
	gbtest.C(t, func(t *gbtest.T) {
		key := "testLockBeforeRLockFunc"
		array := gbarray.New(true)
		go func() {
			gbmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gbmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gbmlock.RLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})

	// Lock before RLockFuncs
	gbtest.C(t, func(t *gbtest.T) {
		key := "testLockBeforeRLockFuncs"
		array := gbarray.New(true)
		go func() {
			gbmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gbmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gbmlock.RLockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gbmlock.RLockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func Test_Locker_TryRLockFunc(t *testing.T) {
	// Lock before TryRLockFunc
	gbtest.C(t, func(t *gbtest.T) {
		key := "testLockBeforeTryRLockFunc"
		array := gbarray.New(true)
		go func() {
			gbmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gbmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gbmlock.TryRLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})

	// Lock before TryRLockFuncs
	gbtest.C(t, func(t *gbtest.T) {
		key := "testLockBeforeTryRLockFuncs"
		array := gbarray.New(true)
		go func() {
			gbmlock.Lock(key)
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			gbmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gbmlock.TryRLockFunc(key, func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(300 * time.Millisecond)
			gbmlock.TryRLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}
