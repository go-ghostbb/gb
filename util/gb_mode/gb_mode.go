// Package gbmode provides release mode management for project.
//
// It uses string to mark the mode instead of integer, which is convenient for configuration.
package gbmode

import (
	gbdebug "ghostbb.io/debug/gb_debug"
	"ghostbb.io/internal/command"
	gbfile "ghostbb.io/os/gb_file"
)

const (
	NOT_SET       = "not-set"
	DEVELOP       = "develop"
	TESTING       = "testing"
	STAGING       = "staging"
	PRODUCT       = "product"
	commandEnvKey = "gb.mode"
)

var (
	// Note that `currentMode` is not concurrent safe.
	currentMode = NOT_SET
)

// Set sets the mode for current application.
func Set(mode string) {
	currentMode = mode
}

// SetDevelop sets current mode DEVELOP for current application.
func SetDevelop() {
	Set(DEVELOP)
}

// SetTesting sets current mode TESTING for current application.
func SetTesting() {
	Set(TESTING)
}

// SetStaging sets current mode STAGING for current application.
func SetStaging() {
	Set(STAGING)
}

// SetProduct sets current mode PRODUCT for current application.
func SetProduct() {
	Set(PRODUCT)
}

// Mode returns current application mode set.
func Mode() string {
	// If current mode is not set, do this auto check.
	if currentMode == NOT_SET {
		if v := command.GetOptWithEnv(commandEnvKey); v != "" {
			// Mode configured from command argument of environment.
			currentMode = v
		} else {
			// If there are source codes found, it's in develop mode, or else in product mode.
			if gbfile.Exists(gbdebug.CallerFilePath()) {
				currentMode = DEVELOP
			} else {
				currentMode = PRODUCT
			}
		}
	}
	return currentMode
}

// IsDevelop checks and returns whether current application is running in DEVELOP mode.
func IsDevelop() bool {
	return Mode() == DEVELOP
}

// IsTesting checks and returns whether current application is running in TESTING mode.
func IsTesting() bool {
	return Mode() == TESTING
}

// IsStaging checks and returns whether current application is running in STAGING mode.
func IsStaging() bool {
	return Mode() == STAGING
}

// IsProduct checks and returns whether current application is running in PRODUCT mode.
func IsProduct() bool {
	return Mode() == PRODUCT
}
