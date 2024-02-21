package etcd

import (
	"context"
	gbmap "ghostbb.io/gb/container/gb_map"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbstr "ghostbb.io/gb/text/gb_str"

	etcd3 "go.etcd.io/etcd/client/v3"
)

// Search searches and returns services with specified condition.
func (r *Registry) Search(ctx context.Context, in gbsvc.SearchInput) ([]gbsvc.Service, error) {
	if in.Prefix == "" && in.Name != "" {
		in.Prefix = gbsvc.NewServiceWithName(in.Name).GetPrefix()
	}

	res, err := r.kv.Get(ctx, in.Prefix, etcd3.WithPrefix())
	if err != nil {
		return nil, err
	}
	services, err := extractResponseToServices(res)
	if err != nil {
		return nil, err
	}
	// Service filter.
	filteredServices := make([]gbsvc.Service, 0)
	for _, service := range services {
		if in.Prefix != "" && !gbstr.HasPrefix(service.GetKey(), in.Prefix) {
			continue
		}
		if in.Name != "" && service.GetName() != in.Name {
			continue
		}
		if in.Version != "" && service.GetVersion() != in.Version {
			continue
		}
		if len(in.Metadata) != 0 {
			m1 := gbmap.NewStrAnyMapFrom(in.Metadata)
			m2 := gbmap.NewStrAnyMapFrom(service.GetMetadata())
			if !m1.IsSubOf(m2) {
				continue
			}
		}
		resultItem := service
		filteredServices = append(filteredServices, resultItem)
	}
	return filteredServices, nil
}

// Watch watches specified condition changes.
// The `key` is the prefix of service key.
func (r *Registry) Watch(ctx context.Context, key string) (gbsvc.Watcher, error) {
	return newWatcher(key, r.client)
}
