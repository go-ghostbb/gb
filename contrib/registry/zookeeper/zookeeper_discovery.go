package zookeeper

import (
	"context"
	gbmap "ghostbb.io/gb/container/gb_map"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbstr "ghostbb.io/gb/text/gb_str"
	"path"
	"strings"
)

// Search searches and returns services with specified condition.
func (r *Registry) Search(_ context.Context, in gbsvc.SearchInput) ([]gbsvc.Service, error) {
	prefix := strings.TrimPrefix(strings.ReplaceAll(in.Prefix, "/", "-"), "-")
	instances, err, _ := r.group.Do(prefix, func() (interface{}, error) {
		serviceNamePath := path.Join(r.opts.namespace, prefix)
		servicesID, _, err := r.conn.Children(serviceNamePath)
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
			byteData, _, err := r.conn.Get(servicePath)
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
	// Service filter.
	filteredServices := make([]gbsvc.Service, 0)
	for _, service := range instances.([]gbsvc.Service) {
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
	return newWatcher(ctx, r.opts.namespace, key, r.conn)
}
