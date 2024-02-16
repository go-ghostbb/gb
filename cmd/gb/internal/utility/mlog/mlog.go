package mlog

import (
	"context"
	gbcmd "ghostbb.io/gb/os/gb_cmd"
	gbenv "ghostbb.io/gb/os/gb_env"
	gblog "ghostbb.io/gb/os/gb_log"
)

const (
	headerPrintEnvName = "GB_CLI_MLOG_HEADER"
)

var (
	ctx    = context.TODO()
	logger = gblog.New()
)

func init() {
	if gbenv.Get(headerPrintEnvName).String() == "1" {
		logger.SetHeaderPrint(true)
	} else {
		logger.SetHeaderPrint(false)
	}

	if gbcmd.GetOpt("debug") != nil || gbcmd.GetOpt("gf.debug") != nil {
		logger.SetHeaderPrint(true)
		logger.SetStackSkip(4)
		logger.SetFlags(logger.GetFlags() | gblog.F_FILE_LONG)
		logger.SetDebug(true)
	} else {
		logger.SetStack(false)
		logger.SetDebug(false)
	}
}

// SetHeaderPrint enables/disables header printing to stdout.
func SetHeaderPrint(enabled bool) {
	logger.SetHeaderPrint(enabled)
	if enabled {
		_ = gbenv.Set(headerPrintEnvName, "1")
	} else {
		_ = gbenv.Set(headerPrintEnvName, "0")
	}
}

func Print(v ...interface{}) {
	logger.Print(ctx, v...)
}

func Printf(format string, v ...interface{}) {
	logger.Printf(ctx, format, v...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(ctx, v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(ctx, format, v...)
}

func Debug(v ...interface{}) {
	logger.Debug(ctx, v...)
}

func Debugf(format string, v ...interface{}) {
	logger.Debugf(ctx, format, v...)
}
