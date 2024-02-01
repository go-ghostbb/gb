// Package redis provides gbredis.Adapter implements using go-redis.
package redis

import (
	"crypto/tls"
	gbredis "ghostbb.io/database/gb_redis"
	gbstr "ghostbb.io/text/gb_str"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis is an implement of Adapter using go-redis.
type Redis struct {
	gbredis.AdapterOperation

	client redis.UniversalClient
	config *gbredis.Config
}

const (
	defaultPoolMaxIdle     = 10
	defaultPoolMaxActive   = 100
	defaultPoolIdleTimeout = 10 * time.Second
	defaultPoolWaitTimeout = 10 * time.Second
	defaultPoolMaxLifeTime = 30 * time.Second
	defaultMaxRetries      = -1
)

func init() {
	gbredis.RegisterAdapterFunc(func(config *gbredis.Config) gbredis.Adapter {
		return New(config)
	})
}

// New creates and returns a redis adapter using go-redis.
func New(config *gbredis.Config) *Redis {
	fillWithDefaultConfiguration(config)
	opts := &redis.UniversalOptions{
		Addrs:           gbstr.SplitAndTrim(config.Address, ","),
		Username:        config.User,
		Password:        config.Pass,
		DB:              config.Db,
		MaxRetries:      defaultMaxRetries,
		PoolSize:        config.MaxActive,
		MinIdleConns:    config.MinIdle,
		MaxIdleConns:    config.MaxIdle,
		ConnMaxLifetime: config.MaxConnLifetime,
		ConnMaxIdleTime: config.IdleTimeout,
		PoolTimeout:     config.WaitTimeout,
		DialTimeout:     config.DialTimeout,
		ReadTimeout:     config.ReadTimeout,
		WriteTimeout:    config.WriteTimeout,
		MasterName:      config.MasterName,
		TLSConfig:       config.TLSConfig,
		Protocol:        config.Protocol,
	}

	var client redis.UniversalClient
	if opts.MasterName != "" {
		redisSentinel := opts.Failover()
		redisSentinel.ReplicaOnly = config.SlaveOnly
		client = redis.NewFailoverClient(redisSentinel)
	} else if len(opts.Addrs) > 1 || config.Cluster {
		client = redis.NewClusterClient(opts.Cluster())
	} else {
		client = redis.NewClient(opts.Simple())
	}

	r := &Redis{
		client: client,
		config: config,
	}
	r.AdapterOperation = r
	return r
}

func fillWithDefaultConfiguration(config *gbredis.Config) {
	// The MaxIdle is the most important attribute of the connection pool.
	// Only if this attribute is set, the created connections from client
	// can not exceed the limit of the server.
	if config.MaxIdle == 0 {
		config.MaxIdle = defaultPoolMaxIdle
	}
	// This value SHOULD NOT exceed the connection limit of redis server.
	if config.MaxActive == 0 {
		config.MaxActive = defaultPoolMaxActive
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = defaultPoolIdleTimeout
	}
	if config.WaitTimeout == 0 {
		config.WaitTimeout = defaultPoolWaitTimeout
	}
	if config.MaxConnLifetime == 0 {
		config.MaxConnLifetime = defaultPoolMaxLifeTime
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = -1
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = -1
	}
	if config.TLSConfig == nil && config.TLS {
		config.TLSConfig = &tls.Config{
			InsecureSkipVerify: config.TLSSkipVerify,
		}
	}
}
