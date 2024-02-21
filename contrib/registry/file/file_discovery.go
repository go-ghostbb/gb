package file

import (
	"context"
	gbmap "ghostbb.io/gb/container/gb_map"
	"ghostbb.io/gb/frame/g"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbfsnotify "ghostbb.io/gb/os/gb_fsnotify"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbstr "ghostbb.io/gb/text/gb_str"
)

// Search searches and returns services with specified condition.
func (r *Registry) Search(ctx context.Context, in gbsvc.SearchInput) (result []gbsvc.Service, err error) {
	services, err := r.getServices(ctx)
	if err != nil {
		return nil, err
	}
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
		result = append(result, resultItem)
	}
	result = r.mergeServices(result)
	return
}

// Watch watches specified condition changes.
// The `key` is the prefix of service key.
func (r *Registry) Watch(ctx context.Context, key string) (watcher gbsvc.Watcher, err error) {
	fileWatcher := &Watcher{
		prefix:    key,
		discovery: r,
		ch:        make(chan gbsvc.Service, 100),
	}
	_, err = gbfsnotify.Add(r.path, func(event *gbfsnotify.Event) {
		if event.IsChmod() {
			return
		}
		if !gbstr.HasPrefix(gbfile.Basename(event.Path), r.getServiceKeyForFile(key)) {
			return
		}
		service, err := r.getServiceByFilePath(event.Path)
		if err != nil {
			return
		}
		fileWatcher.ch <- service
	})
	return fileWatcher, err
}

func (r *Registry) getServices(ctx context.Context) (services []gbsvc.Service, err error) {
	filePaths, err := gbfile.ScanDirFile(r.path, "*", false)
	if err != nil {
		return nil, err
	}
	for _, filePath := range filePaths {
		s, e := r.getServiceByFilePath(filePath)
		if e != nil {
			return nil, e
		}
		// Check service TTL.
		var (
			updateAt    = s.GetMetadata().Get(updateAtKey).GBTime()
			nowTime     = gbtime.Now()
			subDuration = nowTime.Sub(updateAt)
		)
		if updateAt.IsZero() || subDuration > serviceTTL {
			g.Log().Debugf(
				ctx,
				`service "%s" is expired, update at: %s, current: %s, sub duration: %s`,
				s.GetKey(), updateAt.String(), nowTime.String(), subDuration.String(),
			)
			_ = gbfile.Remove(filePath)
			continue
		}
		services = append(services, s)
	}
	services = r.mergeServices(services)
	return
}

func (r *Registry) getServiceByFilePath(filePath string) (gbsvc.Service, error) {
	var (
		fileName    = gbfile.Basename(filePath)
		fileContent = gbfile.GetContents(filePath)
		serviceKey  = gbstr.Replace(fileName, defaultSeparator, gbsvc.DefaultSeparator)
	)

	serviceKey = gbstr.Replace(serviceKey, defaultEndpointHostPortDelimiter, gbsvc.EndpointHostPortDelimiter)
	serviceKey = gbsvc.DefaultSeparator + serviceKey
	return gbsvc.NewServiceWithKV(serviceKey, fileContent)
}

func (r *Registry) mergeServices(services []gbsvc.Service) []gbsvc.Service {
	if len(services) == 0 {
		return services
	}

	var (
		servicePrefixMap = make(map[string]*Service)
		mergeServices    = make([]gbsvc.Service, 0)
	)
	for _, service := range services {
		if v, ok := servicePrefixMap[service.GetPrefix()]; ok {
			v.Endpoints = append(v.Endpoints, service.GetEndpoints()...)
		} else {
			s := NewService(service)
			servicePrefixMap[s.GetPrefix()] = s
			mergeServices = append(mergeServices, s)
		}
	}
	return mergeServices
}
