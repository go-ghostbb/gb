// Package file implements service Registry and Discovery using file.
package file

import (
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbfile "ghostbb.io/gb/os/gb_file"
	"time"
)

var (
	_ gbsvc.Registry = &Registry{}
)

const (
	updateAtKey                      = "UpdateAt"
	serviceTTL                       = 20 * time.Second
	serviceUpdateInterval            = 10 * time.Second
	defaultSeparator                 = "#"
	defaultEndpointHostPortDelimiter = "-"
)

// Registry implements interface Registry using file.
// This implement is usually for testing only.
type Registry struct {
	path string // Local storing folder path for Services.
}

// New creates and returns a gbsvc.Registry implements using file.
func New(path string) gbsvc.Registry {
	if !gbfile.Exists(path) {
		_ = gbfile.Mkdir(path)
	}
	return &Registry{
		path: path,
	}
}
