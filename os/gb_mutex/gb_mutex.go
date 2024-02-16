// Package gbmutex inherits and extends sync.Mutex and sync.RWMutex with more futures.
//
// Note that, it is refracted using stdlib mutex of package sync from GoFrame version v2.5.2.
package gbmutex

// New creates and returns a new mutex.
// Deprecated: use Mutex or RWMutex instead.
func New() *RWMutex {
	return &RWMutex{}
}
