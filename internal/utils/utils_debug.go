package utils

import (
	"github.com/Ghostbb-io/gb/internal/command"
)

const (
	// Debug key for checking if in debug mode.
	commandEnvKeyForDebugKey = "gb.debug"
)

var (
	// isDebugEnabled marks whether GoFrame debug mode is enabled.
	isDebugEnabled = false
)

func init() {
	// Debugging configured.
	value := command.GetOptWithEnv(commandEnvKeyForDebugKey)
	if value == "" || value == "0" || value == "false" {
		isDebugEnabled = false
	} else {
		isDebugEnabled = true
	}
}

// IsDebugEnabled checks and returns whether debug mode is enabled.
// The debug mode is enabled when command argument "gb.debug" or environment "GB_DEBUG" is passed.
func IsDebugEnabled() bool {
	return isDebugEnabled
}

// SetDebugEnabled enables/disables the internal debug info.
func SetDebugEnabled(enabled bool) {
	isDebugEnabled = enabled
}
