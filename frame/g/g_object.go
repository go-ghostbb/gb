package g

import (
	gbdb "ghostbb.io/gb/database/gb_db"
	gbredis "ghostbb.io/gb/database/gb_redis"
	"ghostbb.io/gb/frame/gins"
	gbi18n "ghostbb.io/gb/i18n/gb_i18n"
	gbhttp "ghostbb.io/gb/net/gb_http"
	gbtcp "ghostbb.io/gb/net/gb_tcp"
	gbcfg "ghostbb.io/gb/os/gb_cfg"
	gblog "ghostbb.io/gb/os/gb_log"
	gbres "ghostbb.io/gb/os/gb_res"
	gbvalid "ghostbb.io/gb/util/gb_valid"
)

// Server returns an instance of http server with specified name.
func Server(name ...interface{}) *gbhttp.Server {
	return gins.Server(name...)
}

// TCPServer returns an instance of tcp server with specified name.
func TCPServer(name ...interface{}) *gbtcp.Server {
	return gbtcp.GetServer(name...)
}

// Config returns an instance of config object with specified name.
func Config(name ...string) *gbcfg.Config {
	return gins.Config(name...)
}

// Cfg is alias of Config.
// See Config.
func Cfg(name ...string) *gbcfg.Config {
	return Config(name...)
}

// Resource returns an instance of Resource.
// The parameter `name` is the name for the instance.
func Resource(name ...string) *gbres.Resource {
	return gins.Resource(name...)
}

// I18n returns an instance of gbi18n.Manager.
// The parameter `name` is the name for the instance.
func I18n(name ...string) *gbi18n.Manager {
	return gins.I18n(name...)
}

// Res is alias of Resource.
// See Resource.
func Res(name ...string) *gbres.Resource {
	return Resource(name...)
}

// Log returns an instance of gblog.Logger.
// The parameter `name` is the name for the instance.
func Log(name ...string) *gblog.Logger {
	return gins.Log(name...)
}

// DB returns an instance of database ORM object with specified configuration group name.
func DB(name ...string) *gbdb.DB {
	return gins.Database(name...)
}

// Redis returns an instance of redis client with specified configuration group name.
func Redis(name ...string) *gbredis.Redis {
	return gins.Redis(name...)
}

// Validator is a convenience function, which creates and returns a new validation manager object.
func Validator() *gbvalid.Validator {
	return gbvalid.New()
}
