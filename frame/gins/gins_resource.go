package gins

import (
	gbres "ghostbb.io/os/gb_res"
)

// Resource returns an instance of Resource.
// The parameter `name` is the name for the instance.
func Resource(name ...string) *gbres.Resource {
	return gbres.Instance(name...)
}
