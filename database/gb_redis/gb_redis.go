// Package gbredis provides convenient client for redis server.
//
// Redis Client.
//
// Redis Commands Official: https://redis.io/commands
//
// Redis Chinese Documentation: http://redisdoc.com/
package gbredis

import (
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
)

// AdapterFunc is the function creating redis adapter.
type AdapterFunc func(config *Config) Adapter

type DecodeFunc func(config *Config)

var (
	// defaultAdapterFunc is the default adapter function creating redis adapter.
	defaultAdapterFunc AdapterFunc = func(config *Config) Adapter {
		return nil
	}
	decode DecodeFunc = nil
)

// New creates and returns a redis client.
// It creates a default redis adapter of go-redis.
func New(config ...*Config) (*Redis, error) {
	var (
		usedConfig  *Config
		usedAdapter Adapter
	)
	if len(config) > 0 && config[0] != nil {
		// Redis client with go redis implements adapter from given configuration.
		usedConfig = config[0]
		usedAdapter = defaultAdapterFunc(config[0])
	} else if configFromGlobal, ok := GetConfig(); ok {
		// Redis client with go redis implements adapter from package configuration.
		usedConfig = configFromGlobal
		usedAdapter = defaultAdapterFunc(configFromGlobal)
	}
	if usedConfig == nil {
		return nil, gberror.NewCode(
			gbcode.CodeInvalidConfiguration,
			`no configuration found for creating Redis client`,
		)
	}
	if usedAdapter == nil {
		return nil, gberror.NewCode(
			gbcode.CodeNecessaryPackageNotImport,
			errorNilAdapter,
		)
	}
	if decode != nil {
		decode(usedConfig)
	}
	redis := &Redis{
		config:       usedConfig,
		localAdapter: usedAdapter,
	}
	return redis.initGroup(), nil
}

// NewWithAdapter creates and returns a redis client with given adapter.
func NewWithAdapter(adapter Adapter) (*Redis, error) {
	if adapter == nil {
		return nil, gberror.NewCodef(gbcode.CodeInvalidParameter, `adapter cannot be nil`)
	}
	redis := &Redis{localAdapter: adapter}
	return redis.initGroup(), nil
}

// RegisterAdapterFunc registers default function creating redis adapter.
func RegisterAdapterFunc(adapterFunc AdapterFunc) {
	defaultAdapterFunc = adapterFunc
}

func SetDecode(fn DecodeFunc) {
	decode = fn
}
