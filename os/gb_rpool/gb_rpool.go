// Package gbrpool implements a goroutine reusable pool.
package gbrpool

import (
	"context"
	gblist "ghostbb.io/gb/container/gb_list"
	gbtype "ghostbb.io/gb/container/gb_type"
	gbtimer "ghostbb.io/gb/os/gb_timer"
	gbrand "ghostbb.io/gb/util/gb_rand"
	"time"
)

// Func is the pool function which contains context parameter.
type Func func(ctx context.Context)

// RecoverFunc is the pool runtime panic recover function which contains context parameter.
type RecoverFunc func(ctx context.Context, exception error)

// Pool manages the goroutines using pool.
type Pool struct {
	limit  int          // Max goroutine count limit.
	count  *gbtype.Int  // Current running goroutine count.
	list   *gblist.List // List for asynchronous job adding purpose.
	closed *gbtype.Bool // Is pool closed or not.
}

// localPoolItem is the job item storing in job list.
type localPoolItem struct {
	Ctx  context.Context // Context.
	Func Func            // Job function.
}

const (
	minSupervisorTimerDuration = 500 * time.Millisecond
	maxSupervisorTimerDuration = 1500 * time.Millisecond
)

// Default goroutine pool.
var (
	defaultPool = New()
)

// New creates and returns a new goroutine pool object.
// The parameter `limit` is used to limit the max goroutine count,
// which is not limited in default.
func New(limit ...int) *Pool {
	var (
		pool = &Pool{
			limit:  -1,
			count:  gbtype.NewInt(),
			list:   gblist.New(true),
			closed: gbtype.NewBool(),
		}
		timerDuration = gbrand.D(
			minSupervisorTimerDuration,
			maxSupervisorTimerDuration,
		)
	)
	if len(limit) > 0 && limit[0] > 0 {
		pool.limit = limit[0]
	}
	gbtimer.Add(context.Background(), timerDuration, pool.supervisor)
	return pool
}

// Add pushes a new job to the default goroutine pool.
// The job will be executed asynchronously.
func Add(ctx context.Context, f Func) error {
	return defaultPool.Add(ctx, f)
}

// AddWithRecover pushes a new job to the default pool with specified recover function.
//
// The optional `recoverFunc` is called when any panic during executing of `userFunc`.
// If `recoverFunc` is not passed or given nil, it ignores the panic from `userFunc`.
// The job will be executed asynchronously.
func AddWithRecover(ctx context.Context, userFunc Func, recoverFunc RecoverFunc) error {
	return defaultPool.AddWithRecover(ctx, userFunc, recoverFunc)
}

// Size returns current goroutine count of default goroutine pool.
func Size() int {
	return defaultPool.Size()
}

// Jobs returns current job count of default goroutine pool.
func Jobs() int {
	return defaultPool.Jobs()
}
