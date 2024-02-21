// Package etcd implements service Registry and Discovery using etcd.
package etcd

import (
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/frame/g"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gblog "ghostbb.io/gb/os/gb_log"
	gbstr "ghostbb.io/gb/text/gb_str"
	"time"

	etcd3 "go.etcd.io/etcd/client/v3"
)

var (
	_ gbsvc.Registry = &Registry{}
)

// Registry implements gbsvc.Registry interface.
type Registry struct {
	client       *etcd3.Client
	kv           etcd3.KV
	lease        etcd3.Lease
	keepaliveTTL time.Duration
	logger       gblog.ILogger
}

// Option is the option for the etcd registry.
type Option struct {
	Logger       gblog.ILogger
	KeepaliveTTL time.Duration
}

const (
	// DefaultKeepAliveTTL is the default keepalive TTL.
	DefaultKeepAliveTTL = 10 * time.Second
)

// New creates and returns a new etcd registry.
func New(address string, option ...Option) gbsvc.Registry {
	endpoints := gbstr.SplitAndTrim(address, ",")
	if len(endpoints) == 0 {
		panic(gberror.NewCodef(gbcode.CodeInvalidParameter, `invalid etcd address "%s"`, address))
	}
	client, err := etcd3.New(etcd3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		panic(gberror.Wrap(err, `create etcd client failed`))
	}
	return NewWithClient(client, option...)
}

// NewWithClient creates and returns a new etcd registry with the given client.
func NewWithClient(client *etcd3.Client, option ...Option) *Registry {
	r := &Registry{
		client: client,
		kv:     etcd3.NewKV(client),
	}
	if len(option) > 0 {
		r.logger = option[0].Logger
		r.keepaliveTTL = option[0].KeepaliveTTL
	}
	if r.logger == nil {
		r.logger = g.Log()
	}
	if r.keepaliveTTL == 0 {
		r.keepaliveTTL = DefaultKeepAliveTTL
	}
	return r
}

// extractResponseToServices extracts etcd watch response context to service list.
func extractResponseToServices(res *etcd3.GetResponse) ([]gbsvc.Service, error) {
	if res == nil || res.Kvs == nil {
		return nil, nil
	}
	var (
		services         []gbsvc.Service
		servicePrefixMap = make(map[string]*Service)
	)
	for _, kv := range res.Kvs {
		service, err := gbsvc.NewServiceWithKV(
			string(kv.Key), string(kv.Value),
		)
		if err != nil {
			return services, err
		}
		s := NewService(service)
		if v, ok := servicePrefixMap[service.GetPrefix()]; ok {
			v.Endpoints = append(v.Endpoints, service.GetEndpoints()...)
		} else {
			servicePrefixMap[s.GetPrefix()] = s
			services = append(services, s)
		}
	}
	return services, nil
}
