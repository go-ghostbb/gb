package gbsvc

import (
	gbstr "github.com/Ghostbb-io/gb/text/gb_str"
)

// NewEndpoints creates and returns Endpoints from multiple addresses like:
// "192.168.1.100:80,192.168.1.101:80".
func NewEndpoints(addresses string) Endpoints {
	endpoints := make([]Endpoint, 0)
	for _, address := range gbstr.SplitAndTrim(addresses, EndpointsDelimiter) {
		endpoints = append(endpoints, NewEndpoint(address))
	}
	return endpoints
}

// String formats and returns the Endpoints as a string like:
// "192.168.1.100:80,192.168.1.101:80"
func (es Endpoints) String() string {
	var s string
	for _, endpoint := range es {
		if s != "" {
			s += EndpointsDelimiter
		}
		s += endpoint.String()
	}
	return s
}
