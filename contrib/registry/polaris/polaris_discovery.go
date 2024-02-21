package polaris

import (
	"bytes"
	"context"
	"fmt"
	gbmap "ghostbb.io/gb/container/gb_map"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"strings"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/model"
)

// Search returns the service instances in memory according to the service name.
func (r *Registry) Search(ctx context.Context, in gbsvc.SearchInput) ([]gbsvc.Service, error) {
	if in.Prefix == "" && in.Name != "" {
		service := &Service{
			Service: gbsvc.NewServiceWithName(in.Name),
		}
		in.Prefix = service.GetPrefix()
	}
	in.Prefix = trimAndReplace(in.Prefix)
	// get instances
	instancesResponse, err := r.consumer.GetInstances(&polaris.GetInstancesRequest{
		GetInstancesRequest: model.GetInstancesRequest{
			Service:    in.Prefix,
			Namespace:  r.opt.Namespace,
			Timeout:    &r.opt.Timeout,
			RetryCount: &r.opt.RetryCount,
		},
	})
	if err != nil {
		return nil, err
	}

	serviceInstances := instancesToServiceInstances(instancesResponse.GetInstances())
	// Service filter.
	filteredServices := make([]gbsvc.Service, 0)
	for _, service := range serviceInstances {
		if in.Prefix != "" && !gbstr.HasPrefix(trimAndReplace(service.GetKey()), in.Prefix) {
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

// Watch creates a watcher according to the service name.
func (r *Registry) Watch(ctx context.Context, key string) (gbsvc.Watcher, error) {
	return newWatcher(ctx, r.opt.Namespace, trimAndReplace(key), r.consumer)
}

func instancesToServiceInstances(instances []model.Instance) []gbsvc.Service {
	var (
		serviceInstances = make([]gbsvc.Service, 0, len(instances))
		endpointStr      bytes.Buffer
	)

	for _, instance := range instances {
		if instance.IsHealthy() {
			endpointStr.WriteString(fmt.Sprintf("%s:%d%s", instance.GetHost(), instance.GetPort(), gbsvc.EndpointsDelimiter))
		}
	}
	if endpointStr.Len() > 0 {
		for _, instance := range instances {
			if instance.IsHealthy() {
				serviceInstances = append(serviceInstances, instanceToServiceInstance(instance, gbstr.TrimRight(endpointStr.String(), gbsvc.EndpointsDelimiter), ""))
			}
		}
	}
	return serviceInstances
}

// instanceToServiceInstance converts the instance to service instance.
// instanceID Must be null when creating and adding, and non-null when updating and deleting
func instanceToServiceInstance(instance model.Instance, endpointStr, instanceID string) gbsvc.Service {
	var (
		s         *gbsvc.LocalService
		metadata  = instance.GetMetadata()
		names     = strings.Split(instance.GetService(), instanceIDSeparator)
		endpoints = gbsvc.NewEndpoints(endpointStr)
	)
	if names != nil && len(names) > 4 {
		var name bytes.Buffer
		for i := 3; i < len(names)-1; i++ {
			name.WriteString(names[i])
			if i < len(names)-2 {
				name.WriteString(instanceIDSeparator)
			}
		}
		s = &gbsvc.LocalService{
			Head:       names[0],
			Deployment: names[1],
			Namespace:  names[2],
			Name:       name.String(),
			Version:    metadata[metadataKeyVersion],
			Metadata:   gbconv.Map(metadata),
			Endpoints:  endpoints,
		}
	} else {
		s = &gbsvc.LocalService{
			Name:      instance.GetService(),
			Namespace: instance.GetNamespace(),
			Version:   metadata[metadataKeyVersion],
			Metadata:  gbconv.Map(metadata),
			Endpoints: endpoints,
		}
	}
	service := &Service{
		Service: s,
	}
	if instance.GetId() != "" {
		service.ID = instance.GetId()
	}
	if gbstr.Trim(instanceID) != "" {
		service.ID = instanceID
	}
	return service
}

// trimAndReplace trims the prefix and suffix separator and replaces the separator in the middle.
func trimAndReplace(key string) string {
	key = gbstr.Trim(key, gbsvc.DefaultSeparator)
	key = gbstr.Replace(key, gbsvc.DefaultSeparator, instanceIDSeparator)
	return key
}
