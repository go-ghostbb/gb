package gbsvc

import (
	"context"
	gbcode "ghostbb.io/errors/gb_code"
	gberror "ghostbb.io/errors/gb_error"
)

// Register registers `service` to default registry..
func Register(ctx context.Context, service Service) (Service, error) {
	if defaultRegistry == nil {
		return nil, gberror.NewCodef(gbcode.CodeNotImplemented, `no Registry is registered`)
	}
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return defaultRegistry.Register(ctx, service)
}

// Deregister removes `service` from default registry.
func Deregister(ctx context.Context, service Service) error {
	if defaultRegistry == nil {
		return gberror.NewCodef(gbcode.CodeNotImplemented, `no Registry is registered`)
	}
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	return defaultRegistry.Deregister(ctx, service)
}
