package polaris

import (
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbstr "ghostbb.io/gb/text/gb_str"
)

// Service for wrapping gbsvc.Server and extends extra attributes for polaris purpose.
type Service struct {
	gbsvc.Service        // Common service object.
	ID            string // ID is the unique instance ID as registered, for some registrar server.
}

// GetKey overwrites the GetKey function of gbsvc.Service for replacing separator string.
func (s *Service) GetKey() string {
	key := s.Service.GetKey()
	key = gbstr.Replace(key, gbsvc.DefaultSeparator, instanceIDSeparator)
	key = gbstr.TrimLeft(key, instanceIDSeparator)
	return key
}

// GetPrefix overwrites the GetPrefix function of gbsvc.Service for replacing separator string.
func (s *Service) GetPrefix() string {
	prefix := s.Service.GetPrefix()
	prefix = gbstr.Replace(prefix, gbsvc.DefaultSeparator, instanceIDSeparator)
	prefix = gbstr.TrimLeft(prefix, instanceIDSeparator)
	return prefix
}
