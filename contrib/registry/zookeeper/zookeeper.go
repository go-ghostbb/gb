// Package zookeeper implements service Registry and Discovery using zookeeper.
package zookeeper

import (
	gberror "ghostbb.io/gb/errors/gb_error"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	"github.com/go-zookeeper/zk"
	"golang.org/x/sync/singleflight"
	"time"
)

var _ gbsvc.Registry = &Registry{}

// Content for custom service Marshal/Unmarshal.
type Content struct {
	Key   string
	Value string
}

// Option is etcd registry option.
type Option func(o *options)

type options struct {
	namespace string
	user      string
	password  string
}

// WithRootPath with registry root path.
func WithRootPath(path string) Option {
	return func(o *options) { o.namespace = path }
}

// WithDigestACL with registry password.
func WithDigestACL(user string, password string) Option {
	return func(o *options) {
		o.user = user
		o.password = password
	}
}

// Registry is consul registry
type Registry struct {
	opts  *options
	conn  *zk.Conn
	group singleflight.Group
}

func New(address []string, opts ...Option) *Registry {
	conn, _, err := zk.Connect(address, time.Second*120)
	if err != nil {
		panic(gberror.Wrapf(err,
			"Error with connect to zookeeper"),
		)
	}
	options := &options{
		namespace: "/microservices",
	}
	for _, o := range opts {
		o(options)
	}
	return &Registry{
		opts: options,
		conn: conn,
	}
}
