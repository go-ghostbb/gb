package gbrpool

import (
	"context"
	gbtimer "ghostbb.io/gb/os/gb_timer"
)

// supervisor checks the job list and fork new worker goroutine to handle the job
// if there are jobs but no workers in pool.
func (p *Pool) supervisor(_ context.Context) {
	if p.IsClosed() {
		gbtimer.Exit()
	}
	if p.list.Size() > 0 && p.count.Val() == 0 {
		var number = p.list.Size()
		if p.limit > 0 {
			number = p.limit
		}
		for i := 0; i < number; i++ {
			p.checkAndForkNewGoroutineWorker()
		}
	}
}
