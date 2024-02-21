package file

import (
	"context"
	gbsvc "ghostbb.io/gb/net/gb_svc"
)

// Watcher for file changes watch.
type Watcher struct {
	prefix    string             // Watched prefix key, not file name prefix.
	discovery gbsvc.Discovery    // Service discovery.
	ch        chan gbsvc.Service // Changes that caused by inotify.
}

// Proceed proceeds watch in blocking way.
// It returns all complete services that watched by `key` if any change.
func (w *Watcher) Proceed() (services []gbsvc.Service, err error) {
	<-w.ch
	return w.discovery.Search(context.Background(), gbsvc.SearchInput{
		Prefix: w.prefix,
	})
}

// Close closes the watcher.
func (w *Watcher) Close() error {
	close(w.ch)
	return nil
}
