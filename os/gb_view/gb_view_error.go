package gbview

import (
	gbcmd "ghostbb.io/gb/os/gb_cmd"
)

const (
	// commandEnvKeyForErrorPrint is used to specify the key controlling error printing to stdout.
	// This error is designed not to be returned by functions.
	commandEnvKeyForErrorPrint = "gb.view.errorprint"
)

// errorPrint checks whether printing error to stdout.
func errorPrint() bool {
	return gbcmd.GetOptWithEnv(commandEnvKeyForErrorPrint, true).Bool()
}
