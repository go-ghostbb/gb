package gbmutex

import "sync"

// Mutex is a high level Mutex, which implements more rich features for mutex.
type Mutex struct {
	sync.Mutex
}

// LockFunc locks the mutex for writing with given callback function `f`.
// If there's a write/reading lock the mutex, it will block until the lock is released.
//
// It releases the lock after `f` is executed.
func (m *Mutex) LockFunc(f func()) {
	m.Lock()
	defer m.Unlock()
	f()
}

// TryLockFunc tries locking the mutex for writing with given callback function `f`.
// it returns true immediately if success, or if there's a write/reading lock on the mutex,
// it returns false immediately.
//
// It releases the lock after `f` is executed.
func (m *Mutex) TryLockFunc(f func()) (result bool) {
	if m.TryLock() {
		result = true
		defer m.Unlock()
		f()
	}
	return
}
