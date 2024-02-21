package nacos

import (
	"context"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbconv "ghostbb.io/gb/util/gb_conv"

	"github.com/joy999/nacos-sdk-go/vo"
)

// Register registers `service` to Registry.
// Note that it returns a new Service if it changes the input Service with custom one.
func (reg *Registry) Register(ctx context.Context, service gbsvc.Service) (registered gbsvc.Service, err error) {
	metadata := map[string]string{}
	endpoints := service.GetEndpoints()
	p := vo.BatchRegisterInstanceParam{
		ServiceName: service.GetName(),
		GroupName:   reg.groupName,
		Instances:   make([]vo.RegisterInstanceParam, 0, len(endpoints)),
	}

	for k, v := range service.GetMetadata() {
		metadata[k] = gbconv.String(v)
	}

	for _, endpoint := range endpoints {
		p.Instances = append(p.Instances, vo.RegisterInstanceParam{
			Ip:          endpoint.Host(),
			Port:        uint64(endpoint.Port()),
			ServiceName: service.GetName(),
			Metadata:    metadata,
			Weight:      100,
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
			ClusterName: reg.clusterName,
			GroupName:   reg.groupName,
		})
	}

	if _, err = reg.client.BatchRegisterInstance(p); err != nil {
		return
	}

	registered = service

	return
}

// Deregister off-lines and removes `service` from the Registry.
func (reg *Registry) Deregister(ctx context.Context, service gbsvc.Service) (err error) {
	c := reg.client

	for _, endpoint := range service.GetEndpoints() {
		if _, err = c.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          endpoint.Host(),
			Port:        uint64(endpoint.Port()),
			ServiceName: service.GetName(),
			Ephemeral:   true,
			Cluster:     reg.clusterName,
			GroupName:   reg.groupName,
		}); err != nil {
			return
		}
	}

	return
}
