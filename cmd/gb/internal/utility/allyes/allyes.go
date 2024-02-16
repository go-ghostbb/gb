package allyes

import (
	gbcmd "ghostbb.io/gb/os/gb_cmd"
	gbenv "ghostbb.io/gb/os/gb_env"
)

const (
	EnvName = "GB_CLI_ALL_YES"
)

// Init initializes the package manually.
func Init() {
	if gbcmd.GetOpt("y") != nil {
		gbenv.MustSet(EnvName, "1")
	}
}

// Check checks whether option allow all yes for command.
func Check() bool {
	return gbenv.Get(EnvName).String() == "1"
}
