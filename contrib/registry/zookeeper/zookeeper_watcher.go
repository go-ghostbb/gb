package zookeeper

import (
	"context"
	"errors"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	"github.com/go-zookeeper/zk"
	"golang.org/x/sync/singleflight"
	"path"
	"strings"
)

var _ gbsvc.Watcher = (*watcher)(nil)

// ErrWatcherStopped is the certain error for watcher closed.
var ErrWatcherStopped = errors.New("watcher stopped")

type watcher struct {
	ctx       context.Context
	event     chan zk.Event
	conn      *zk.Conn
	cancel    context.CancelFunc
	prefix    string
	nameSpace string
	group     singleflight.Group
}

func newWatcher(ctx context.Context, nameSpace, prefix string, conn *zk.Conn) (*watcher, error) {
	w := &watcher{
		conn:      conn,
		event:     make(chan zk.Event, 1),
		nameSpace: nameSpace,
		prefix:    prefix,
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	go w.watch(w.ctx)
	return w, nil
}

// Proceed proceeds watch in blocking way.
// It returns all complete services that watched by `key` if any change.
func (w *watcher) Proceed() ([]gbsvc.Service, error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case e := <-w.event:
		if e.State == zk.StateDisconnected {
			return nil, gberror.Wrapf(
				ErrWatcherStopped,
				"watcher stopped",
			)
		}
		if e.Err != nil {
			return nil, e.Err
		}
		return w.getServicesByPrefix()
	}
}

func (w *watcher) getServicesByPrefix() ([]gbsvc.Service, error) {
	prefix := strings.TrimPrefix(strings.ReplaceAll(w.prefix, "/", "-"), "-")
	serviceNamePath := path.Join(w.nameSpace, prefix)
	instances, err, _ := w.group.Do(serviceNamePath, func() (interface{}, error) {
		servicesID, _, err := w.conn.Children(serviceNamePath)
		if err != nil {
			return nil, gberror.Wrapf(
				err,
				"Error with search the children node under %s",
				serviceNamePath,
			)
		}
		items := make([]gbsvc.Service, 0, len(servicesID))
		for _, service := range servicesID {
			servicePath := path.Join(serviceNamePath, service)
			byteData, _, err := w.conn.Get(servicePath)
			if err != nil {
				return nil, gberror.Wrapf(
					err,
					"Error with node data which name is %s",
					servicePath,
				)
			}
			item, err := unmarshal(byteData)
			if err != nil {
				return nil, gberror.Wrapf(
					err,
					"Error with unmarshal node data to Content",
				)
			}
			svc, err := gbsvc.NewServiceWithKV(item.Key, item.Value)
			if err != nil {
				return nil, gberror.Wrapf(
					err,
					"Error with new service with KV in Content",
				)
			}
			items = append(items, svc)
		}
		return items, nil
	})
	if err != nil {
		return nil, gberror.Wrapf(
			err,
			"Error with group do",
		)
	}
	return instances.([]gbsvc.Service), nil
}

// Close closes the watcher.
func (w *watcher) Close() error {
	w.cancel()
	return nil
}

func (w *watcher) watch(ctx context.Context) {
	prefix := strings.TrimPrefix(strings.ReplaceAll(w.prefix, "/", "-"), "-")
	serviceNamePath := path.Join(w.nameSpace, prefix)
	for {

		if w.conn.State() == zk.StateConnected || w.conn.State() == zk.StateHasSession {
			// each watch action is only valid once
			_, _, ch, err := w.conn.ChildrenW(serviceNamePath)
			if err != nil {
				w.event <- zk.Event{Err: err}
			}
			select {
			case <-ctx.Done():
				return
			default:
				w.event <- <-ch
			}
		}
	}
}
