// Package gbcmd provides console operations, like options/arguments reading and command running.
package gbcmd

import (
	gbvar "github.com/Ghostbb-io/gb/container/gb_var"
	"github.com/Ghostbb-io/gb/internal/command"
	"github.com/Ghostbb-io/gb/internal/utils"
	gbctx "github.com/Ghostbb-io/gb/os/gb_ctx"
	"os"
)

const (
	CtxKeyParser    gbctx.StrKey = `CtxKeyParser`
	CtxKeyCommand   gbctx.StrKey = `CtxKeyCommand`
	CtxKeyArguments gbctx.StrKey = `CtxKeyArguments`
)

const (
	helpOptionName        = "help"
	helpOptionNameShort   = "h"
	maxLineChars          = 120
	tracingInstrumentName = "github.com/gogf/gf/v2/os/gcmd.Command"
	tagNameName           = "name"
	tagNameShort          = "short"
)

// Init does custom initialization.
func Init(args ...string) {
	command.Init(args...)
}

// GetOpt returns the option value named `name` as gbvar.Var.
func GetOpt(name string, def ...string) *gbvar.Var {
	if v := command.GetOpt(name, def...); v != "" {
		return gbvar.New(v)
	}
	if command.ContainsOpt(name) {
		return gbvar.New("")
	}
	return nil
}

// GetOptAll returns all parsed options.
func GetOptAll() map[string]string {
	return command.GetOptAll()
}

// GetArg returns the argument at `index` as gbvar.Var.
func GetArg(index int, def ...string) *gbvar.Var {
	if v := command.GetArg(index, def...); v != "" {
		return gbvar.New(v)
	}
	return nil
}

// GetArgAll returns all parsed arguments.
func GetArgAll() []string {
	return command.GetArgAll()
}

// GetOptWithEnv returns the command line argument of the specified `key`.
// If the argument does not exist, then it returns the environment variable with specified `key`.
// It returns the default value `def` if none of them exists.
//
// Fetching Rules:
// 1. Command line arguments are in lowercase format, eg: gf.`package name`.<variable name>;
// 2. Environment arguments are in uppercase format, eg: GF_`package name`_<variable name>；
func GetOptWithEnv(key string, def ...interface{}) *gbvar.Var {
	cmdKey := utils.FormatCmdKey(key)
	if command.ContainsOpt(cmdKey) {
		return gbvar.New(GetOpt(cmdKey))
	} else {
		envKey := utils.FormatEnvKey(key)
		if r, ok := os.LookupEnv(envKey); ok {
			return gbvar.New(r)
		} else {
			if len(def) > 0 {
				return gbvar.New(def[0])
			}
		}
	}
	return nil
}

// BuildOptions builds the options as string.
func BuildOptions(m map[string]string, prefix ...string) string {
	options := ""
	leadStr := "-"
	if len(prefix) > 0 {
		leadStr = prefix[0]
	}
	for k, v := range m {
		if len(options) > 0 {
			options += " "
		}
		options += leadStr + k
		if v != "" {
			options += "=" + v
		}
	}
	return options
}
