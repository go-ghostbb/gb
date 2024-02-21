package file

import (
	"context"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtimer "ghostbb.io/gb/os/gb_timer"
	gbstr "ghostbb.io/gb/text/gb_str"
)

// Register registers `service` to Registry.
// Note that it returns a new Service if it changes the input Service with custom one.
func (r *Registry) Register(ctx context.Context, service gbsvc.Service) (registered gbsvc.Service, err error) {
	service = NewService(service)
	service.GetMetadata().Set(updateAtKey, gbtime.Now())
	var (
		filePath    = r.getServiceFilePath(service)
		fileContent = service.GetValue()
	)
	err = gbfile.PutContents(filePath, fileContent)
	if err == nil {
		gbtimer.Add(ctx, serviceUpdateInterval, func(ctx context.Context) {
			if !gbfile.Exists(filePath) {
				gbtimer.Exit()
			}
			// Update TTL in timer.
			service, _ = r.getServiceByFilePath(filePath)
			if service != nil {
				service.GetMetadata().Set(updateAtKey, gbtime.Now())
			}
			_ = gbfile.PutContents(filePath, service.GetValue())
		})
	}
	return service, err
}

// Deregister off-lines and removes `service` from the Registry.
func (r *Registry) Deregister(ctx context.Context, service gbsvc.Service) error {
	return gbfile.Remove(r.getServiceFilePath(service))
}

func (r *Registry) getServiceFilePath(service gbsvc.Service) string {
	return gbfile.Join(r.path, r.getServiceFileName(service))
}

func (r *Registry) getServiceFileName(service gbsvc.Service) string {
	return r.getServiceKeyForFile(service.GetKey())
}

func (r *Registry) getServiceKeyForFile(key string) string {
	key = gbstr.Replace(key, gbsvc.DefaultSeparator, defaultSeparator)
	key = gbstr.Trim(key, defaultSeparator)
	key = gbstr.Replace(key, gbsvc.EndpointHostPortDelimiter, defaultEndpointHostPortDelimiter)
	return key
}
