package gbsvc

import (
	"context"
	gbjson "ghostbb.io/gb/encoding/gb_json"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/intlog"
	gbcmd "ghostbb.io/gb/os/gb_cmd"
	gbstr "ghostbb.io/gb/text/gb_str"
)

// LocalService provides a default implements for interface Service.
type LocalService struct {
	Head       string    // Service custom head string in service key.
	Deployment string    // Service deployment name, eg: dev, qa, staging, prod, etc.
	Namespace  string    // Service Namespace, to indicate different services in the same environment with the same Name.
	Name       string    // Name for the service.
	Version    string    // Service version, eg: v1.0.0, v2.1.1, etc.
	Endpoints  Endpoints // Service Endpoints, pattern: IP:port, eg: 192.168.1.2:8000.
	Metadata   Metadata  // Custom data for this service, which can be set using JSON by environment or command-line.
}

// NewServiceWithName creates and returns a default implements for interface Service by service name.
func NewServiceWithName(name string) Service {
	s := &LocalService{
		Name:     name,
		Metadata: make(Metadata),
	}
	s.autoFillDefaultAttributes()
	return s
}

// NewServiceWithKV creates and returns a default implements for interface Service by key-value pair string.
func NewServiceWithKV(key, value string) (Service, error) {
	var (
		err   error
		array = gbstr.Split(gbstr.Trim(key, DefaultSeparator), DefaultSeparator)
	)
	if len(array) < 6 {
		err = gberror.NewCodef(gbcode.CodeInvalidParameter, `invalid service key "%s"`, key)
		return nil, err
	}
	s := &LocalService{
		Head:       array[0],
		Deployment: array[1],
		Namespace:  array[2],
		Name:       array[3],
		Version:    array[4],
		Endpoints:  NewEndpoints(array[5]),
		Metadata:   make(Metadata),
	}
	s.autoFillDefaultAttributes()
	if len(value) > 0 {
		if err = gbjson.Unmarshal([]byte(value), &s.Metadata); err != nil {
			err = gberror.WrapCodef(gbcode.CodeInvalidParameter, err, `invalid service value "%s"`, value)
			return nil, err
		}
	}
	return s, nil
}

// GetName returns the name of the service.
// The name is necessary for a service, and should be unique among services.
func (s *LocalService) GetName() string {
	return s.Name
}

// GetVersion returns the version of the service.
// It is suggested using GNU version naming like: v1.0.0, v2.0.1, v2.1.0-rc.
// A service can have multiple versions deployed at once.
// If no version set in service, the default version of service is "latest".
func (s *LocalService) GetVersion() string {
	return s.Version
}

// GetKey formats and returns a unique key string for service.
// The result key is commonly used for key-value registrar server.
func (s *LocalService) GetKey() string {
	serviceNameUnique := s.GetPrefix()
	serviceNameUnique += DefaultSeparator + s.Endpoints.String()
	return serviceNameUnique
}

// GetValue formats and returns the value of the service.
// The result value is commonly used for key-value registrar server.
func (s *LocalService) GetValue() string {
	b, err := gbjson.Marshal(s.Metadata)
	if err != nil {
		intlog.Errorf(context.TODO(), `%+v`, err)
	}
	return string(b)
}

// GetPrefix formats and returns the key prefix string.
// The result prefix string is commonly used in key-value registrar server
// for service searching.
//
// Take etcd server for example, the prefix string is used like:
// `etcdctl get /services/prod/hello.svc --prefix`
func (s *LocalService) GetPrefix() string {
	s.autoFillDefaultAttributes()
	return DefaultSeparator + gbstr.Join(
		[]string{
			s.Head,
			s.Deployment,
			s.Namespace,
			s.Name,
			s.Version,
		},
		DefaultSeparator,
	)
}

// GetMetadata returns the Metadata map of service.
// The Metadata is key-value pair map specifying extra attributes of a service.
func (s *LocalService) GetMetadata() Metadata {
	return s.Metadata
}

// GetEndpoints returns the Endpoints of service.
// The Endpoints contain multiple host/port information of service.
func (s *LocalService) GetEndpoints() Endpoints {
	return s.Endpoints
}

func (s *LocalService) autoFillDefaultAttributes() {
	if s.Head == "" {
		s.Head = gbcmd.GetOptWithEnv(EnvPrefix, DefaultHead).String()
	}
	if s.Deployment == "" {
		s.Deployment = gbcmd.GetOptWithEnv(EnvDeployment, DefaultDeployment).String()
	}
	if s.Namespace == "" {
		s.Namespace = gbcmd.GetOptWithEnv(EnvNamespace, DefaultNamespace).String()
	}
	if s.Name == "" {
		s.Name = gbcmd.GetOptWithEnv(EnvName).String()
	}
	if s.Version == "" {
		s.Version = gbcmd.GetOptWithEnv(EnvVersion, DefaultVersion).String()
	}
}
