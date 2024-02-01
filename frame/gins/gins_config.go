package gins

import (
	gbcfg "ghostbb.io/os/gb_cfg"
)

// Config returns an instance of View with default settings.
// The parameter `name` is the name for the instance.
func Config(name ...string) *gbcfg.Config {
	return gbcfg.Instance(name...)
}
