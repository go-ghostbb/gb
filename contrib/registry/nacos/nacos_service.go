package nacos

import (
	"fmt"
	gbmap "ghostbb.io/gb/container/gb_map"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbstr "ghostbb.io/gb/text/gb_str"

	"github.com/joy999/nacos-sdk-go/model"
)

// NewServiceFromInstance new one service from instance
func NewServiceFromInstance(instance []model.Instance) gbsvc.Service {
	n := len(instance)
	if n == 0 {
		return nil
	}
	serviceName := instance[0].ServiceName
	endpoints := make(gbsvc.Endpoints, 0, n)
	for i := 0; i < n; i++ {
		if instance[0].ServiceName != serviceName {
			return nil
		}
		endpoints = append(endpoints, gbsvc.NewEndpoint(fmt.Sprintf("%s%s%d", instance[i].Ip, gbsvc.EndpointHostPortDelimiter, int(instance[i].Port))))
	}
	if gbstr.Contains(serviceName, cstServiceSeparator) {
		arr := gbstr.SplitAndTrim(serviceName, cstServiceSeparator)
		serviceName = arr[1]
	}

	return &gbsvc.LocalService{
		Endpoints: endpoints,
		Name:      serviceName,
		Metadata:  gbmap.NewStrStrMapFrom(instance[0].Metadata).MapStrAny(),
		Version:   gbsvc.DefaultVersion,
	}
}

// NewServicesFromInstances new some services from some instances
func NewServicesFromInstances(instances []model.Instance) []gbsvc.Service {
	serviceMap := map[string][]model.Instance{}
	for _, inst := range instances {
		serviceMap[inst.ServiceName] = append(serviceMap[inst.ServiceName], inst)
	}

	services := make([]gbsvc.Service, 0, len(serviceMap))
	for _, insts := range serviceMap {
		services = append(services, NewServiceFromInstance(insts))
	}

	return services
}
