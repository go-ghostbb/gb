package gbbuild

import (
	"context"
	"ghostbb.io/gb"
	gbvar "ghostbb.io/gb/container/gb_var"
	gbbase64 "ghostbb.io/gb/encoding/gb_base64"
	"ghostbb.io/gb/internal/intlog"
	"ghostbb.io/gb/internal/json"
	"runtime"
)

// BuildInfo maintains the built info of current binary.
type BuildInfo struct {
	GB      string                 // Built used GB version.
	Golang  string                 // Built used Golang version.
	Git     string                 // Built used git repo. commit id and datetime.
	Time    string                 // Built datetime.
	Version string                 // Built version.
	Data    map[string]interface{} // All custom built data key-value pairs.
}

const (
	gbVersion    = `gbVersion`
	goVersion    = `goVersion`
	BuiltGit     = `builtGit`
	BuiltTime    = `builtTime`
	BuiltVersion = `builtVersion`
)

var (
	builtInVarStr = ""                       // Raw variable base64 string, which is injected by go build flags.
	builtInVarMap = map[string]interface{}{} // Binary custom variable map decoded.
)

func init() {
	// The `builtInVarStr` is injected by go build flags.
	if builtInVarStr != "" {
		err := json.UnmarshalUseNumber(gbbase64.MustDecodeString(builtInVarStr), &builtInVarMap)
		if err != nil {
			intlog.Errorf(context.TODO(), `%+v`, err)
		}
		builtInVarMap[gbVersion] = gb.VERSION
		builtInVarMap[goVersion] = runtime.Version()
		intlog.Printf(context.TODO(), "build variables: %+v", builtInVarMap)
	} else {
		intlog.Print(context.TODO(), "no build variables")
	}
}

// Info returns the basic built information of the binary as map.
// Note that it should be used with gb-cli tool "gb build",
// which automatically injects necessary information into the binary.
func Info() BuildInfo {
	return BuildInfo{
		GB:      Get(gbVersion).String(),
		Golang:  Get(goVersion).String(),
		Git:     Get(BuiltGit).String(),
		Time:    Get(BuiltTime).String(),
		Version: Get(BuiltVersion).String(),
		Data:    Data(),
	}
}

// Get retrieves and returns the build-in binary variable with given name.
func Get(name string, def ...interface{}) *gbvar.Var {
	if v, ok := builtInVarMap[name]; ok {
		return gbvar.New(v)
	}
	if len(def) > 0 {
		return gbvar.New(def[0])
	}
	return nil
}

// Data returns the custom build-in variables as map.
func Data() map[string]interface{} {
	return builtInVarMap
}
