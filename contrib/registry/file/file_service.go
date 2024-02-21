package file

import (
	gbjson "ghostbb.io/gb/encoding/gb_json"
	gbsvc "ghostbb.io/gb/net/gb_svc"
)

// Service wrapper.
type Service struct {
	gbsvc.Service
	Endpoints gbsvc.Endpoints
	Metadata  gbsvc.Metadata
}

// NewService creates and returns local Service from gbsvc.Service interface object.
func NewService(service gbsvc.Service) *Service {
	s, ok := service.(*Service)
	if ok {
		if s.Endpoints == nil {
			s.Endpoints = make(gbsvc.Endpoints, 0)
		}
		if s.Metadata == nil {
			s.Metadata = make(gbsvc.Metadata)
		}
		return s
	}
	s = &Service{
		Service:   service,
		Endpoints: make(gbsvc.Endpoints, 0),
		Metadata:  make(gbsvc.Metadata),
	}
	if len(service.GetEndpoints()) > 0 {
		s.Endpoints = service.GetEndpoints()
	}
	if len(service.GetMetadata()) > 0 {
		s.Metadata = service.GetMetadata()
	}
	return s
}

// GetMetadata returns the Metadata map of service.
// The Metadata is key-value pair map specifying extra attributes of a service.
func (s *Service) GetMetadata() gbsvc.Metadata {
	return s.Metadata
}

// GetEndpoints returns the Endpoints of service.
// The Endpoints contain multiple host/port information of service.
func (s *Service) GetEndpoints() gbsvc.Endpoints {
	return s.Endpoints
}

// GetValue formats and returns the value of the service.
// The result value is commonly used for key-value registrar server.
func (s *Service) GetValue() string {
	b, _ := gbjson.Marshal(s.Metadata)
	return string(b)
}
