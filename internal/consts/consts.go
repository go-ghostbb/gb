// Package consts defines constants that are shared all among packages of framework.
package consts

const (
	ConfigNodeNameDatabase        = "database"
	ConfigNodeNameLogger          = "logger"
	ConfigNodeNameRedis           = "redis"
	ConfigNodeNameCache           = "cache"
	ConfigNodeNameServer          = "server"
	ConfigNodeNameServerSecondary = "httpserver"

	// StackFilterKeyForGoFrame is the stack filtering key for all GoFrame module paths.
	// Eg: .../pkg/mod/ghostbb.io/gb/@v2.0.0-20211011134327-54dd11f51122/debug/gbdebug/gbdebug_caller.go
	StackFilterKeyForGoFrame = "ghostbb.io/gb/"
)
