// Package polaris implements service Registry and Discovery using polaris.
package polaris

import (
	"ghostbb.io/gb/frame/g"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gblog "ghostbb.io/gb/os/gb_log"
	"time"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/config"
)

var (
	_ gbsvc.Registry = &Registry{}
)

const (
	instanceIDSeparator = "-"
	metadataKeyKind     = "kind"
	metadataKeyVersion  = "version"
)

type options struct {
	// required, namespace in polaris
	Namespace string

	// required, service access token
	ServiceToken string

	// optional, protocol in polaris. Default value is nil, it means use protocol config in service
	Protocol *string

	// service weight in polaris. Default value is 100, 0 <= weight <= 10000
	Weight int

	// service priority. Default value is 0. The smaller the value, the lower the priority
	Priority int

	// To show service is healthy or not. Default value is True.
	Healthy bool

	// To show service is isolate or not. Default value is False.
	Isolate bool

	// TTL timeout. if the node needs to use a heartbeat to report, required. If not set,server will throw ErrorCode-400141
	TTL int

	// Timeout for the single query. Default value is global config
	// Total is (1+RetryCount) * Timeout
	Timeout time.Duration

	// optional, retry count. Default value is global config
	RetryCount int

	// optional, logger for polaris
	Logger gblog.ILogger
}

// Option The option is a polaris option.
type Option func(o *options)

// Registry is polaris registry.
type Registry struct {
	opt      options
	provider polaris.ProviderAPI
	consumer polaris.ConsumerAPI
}

// WithNamespace with the Namespace option.
func WithNamespace(namespace string) Option {
	return func(o *options) { o.Namespace = namespace }
}

// WithServiceToken with ServiceToken option.
func WithServiceToken(serviceToken string) Option {
	return func(o *options) { o.ServiceToken = serviceToken }
}

// WithProtocol with the Protocol option.
func WithProtocol(protocol string) Option {
	return func(o *options) { o.Protocol = &protocol }
}

// WithWeight with the Weight option.
func WithWeight(weight int) Option {
	return func(o *options) { o.Weight = weight }
}

// WithHealthy with the Healthy option.
func WithHealthy(healthy bool) Option {
	return func(o *options) { o.Healthy = healthy }
}

// WithIsolate with the Isolate option.
func WithIsolate(isolate bool) Option {
	return func(o *options) { o.Isolate = isolate }
}

// WithTTL with the TTL option.
func WithTTL(TTL int) Option {
	return func(o *options) { o.TTL = TTL }
}

// WithTimeout the Timeout option.
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.Timeout = timeout }
}

// WithRetryCount with RetryCount option.
func WithRetryCount(retryCount int) Option {
	return func(o *options) { o.RetryCount = retryCount }
}

// WithLogger with the Logger option.
func WithLogger(logger gblog.ILogger) Option {
	return func(o *options) { o.Logger = logger }
}

// New create a new registry.
func New(provider polaris.ProviderAPI, consumer polaris.ConsumerAPI, opts ...Option) gbsvc.Registry {
	op := options{
		Namespace:    gbsvc.DefaultNamespace,
		ServiceToken: "",
		Protocol:     nil,
		Weight:       100,
		Priority:     0,
		Healthy:      true,
		Isolate:      false,
		TTL:          0,
		Timeout:      0,
		RetryCount:   0,
		Logger:       g.Log(),
	}
	for _, option := range opts {
		option(&op)
	}
	return &Registry{
		opt:      op,
		provider: provider,
		consumer: consumer,
	}
}

// NewWithConfig new a registry with config.
func NewWithConfig(conf config.Configuration, opts ...Option) gbsvc.Registry {
	provider, err := polaris.NewProviderAPIByConfig(conf)
	if err != nil {
		panic(err)
	}
	consumer, err := polaris.NewConsumerAPIByConfig(conf)
	if err != nil {
		panic(err)
	}
	return New(provider, consumer, opts...)
}
