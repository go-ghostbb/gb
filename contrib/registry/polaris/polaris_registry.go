package polaris

import (
	"context"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/model"
)

// Register the registration.
func (r *Registry) Register(ctx context.Context, service gbsvc.Service) (gbsvc.Service, error) {
	// Replace input service to custom service types.
	service = &Service{
		Service: service,
	}
	// Register logic.
	var ids = make([]string, 0, len(service.GetEndpoints()))
	for _, endpoint := range service.GetEndpoints() {
		// medata
		var (
			rmd            map[string]interface{}
			serviceName    = service.GetPrefix()
			serviceVersion = service.GetVersion()
		)
		if service.GetMetadata().IsEmpty() {
			rmd = map[string]interface{}{
				metadataKeyKind:    gbsvc.DefaultProtocol,
				metadataKeyVersion: serviceVersion,
			}
		} else {
			rmd = make(map[string]interface{}, len(service.GetMetadata())+2)
			rmd[metadataKeyKind] = gbsvc.DefaultProtocol
			if protocol, ok := service.GetMetadata()[gbsvc.MDProtocol]; ok {
				rmd[metadataKeyKind] = gbconv.String(protocol)
			}
			rmd[metadataKeyVersion] = serviceVersion
			for k, v := range service.GetMetadata() {
				rmd[k] = v
			}
		}
		// Register RegisterInstance Service registration is performed synchronously,
		// and heartbeat reporting is automatically performed
		registeredService, err := r.provider.RegisterInstance(
			&polaris.InstanceRegisterRequest{
				InstanceRegisterRequest: model.InstanceRegisterRequest{
					Service:      serviceName,
					ServiceToken: r.opt.ServiceToken,
					Namespace:    r.opt.Namespace,
					Host:         endpoint.Host(),
					Port:         endpoint.Port(),
					Protocol:     r.opt.Protocol,
					Weight:       &r.opt.Weight,
					Priority:     &r.opt.Priority,
					Version:      &serviceVersion,
					Metadata:     gbconv.MapStrStr(rmd),
					Healthy:      &r.opt.Healthy,
					Isolate:      &r.opt.Isolate,
					TTL:          &r.opt.TTL,
					Timeout:      &r.opt.Timeout,
					RetryCount:   &r.opt.RetryCount,
				},
			})
		if err != nil {
			return nil, err
		}
		ids = append(ids, registeredService.InstanceID)
	}
	// need to set InstanceID for Deregister
	service.(*Service).ID = gbstr.Join(ids, instanceIDSeparator)
	return service, nil
}

// Deregister the registration.
func (r *Registry) Deregister(ctx context.Context, service gbsvc.Service) error {
	var (
		err   error
		split = gbstr.Split(service.(*Service).ID, instanceIDSeparator)
	)
	for i, endpoint := range service.GetEndpoints() {
		// Deregister
		err = r.provider.Deregister(
			&polaris.InstanceDeRegisterRequest{
				InstanceDeRegisterRequest: model.InstanceDeRegisterRequest{
					Service:      service.GetPrefix(),
					ServiceToken: r.opt.ServiceToken,
					Namespace:    r.opt.Namespace,
					InstanceID:   split[i],
					Host:         endpoint.Host(),
					Port:         endpoint.Port(),
					Timeout:      &r.opt.Timeout,
					RetryCount:   &r.opt.RetryCount,
				},
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}
