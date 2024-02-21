package nacos

import (
	"context"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbsvc "ghostbb.io/gb/net/gb_svc"

	"github.com/joy999/nacos-sdk-go/model"
)

// Watcher used to manage service event such as update.
type Watcher struct {
	ctx   context.Context
	event chan *watchEvent
	close func() error
}

// watchEvent
type watchEvent struct {
	Services []model.Instance
	Err      error
}

// newWatcher new a Watcher's instance
func newWatcher(ctx context.Context) *Watcher {
	w := &Watcher{
		ctx:   ctx,
		event: make(chan *watchEvent, 10),
	}
	return w
}

// Proceed proceeds watch in blocking way.
// It returns all complete services that watched by `key` if any change.
func (w *Watcher) Proceed() (services []gbsvc.Service, err error) {
	e, ok := <-w.event
	if !ok || e == nil {
		err = gberror.NewCode(gbcode.CodeNil)
		return
	}
	if e.Err != nil {
		err = e.Err
		return
	}
	services = NewServicesFromInstances(e.Services)
	return
}

// Close closes the watcher.
func (w *Watcher) Close() (err error) {
	if w.close != nil {
		err = w.close()
	}
	return
}

// SetCloseFunc set the close callback function
func (w *Watcher) SetCloseFunc(close func() error) {
	w.close = close
}

// Push add the services watchevent to event queue
func (w *Watcher) Push(services []model.Instance, err error) {
	w.event <- &watchEvent{
		Services: services,
		Err:      err,
	}
}
