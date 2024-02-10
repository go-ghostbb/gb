package gbdb

import (
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/intlog"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
	config DatabaseConfig // Current config.
}

// Register registers custom database driver to gbdb.
func Register(name string, driver IDriver) error {
	driverMap[name] = newDriverWrapper(driver)
	return nil
}

func NewDBByConfig(name string, config DatabaseConfig) (*DB, error) {
	var (
		db  = new(DB)
		err error
	)

	if err = db.setConfig(config); err != nil {
		intlog.Printf(gbctx.New(), "%s | %s | database set config error:%v.", name, config.Type, err)
	}
	db.config.instance = name

	if v, ok := driverMap[config.Type]; ok {
		if db.DB, err = v.New(db.config); err != nil {
			return nil, err
		}
		intlog.Printf(gbctx.New(), "%s | %s | database connection successful.", name, config.Type)
		return db, nil
	}

	errorMsg := `cannot find database driver for specified database type "%s"`
	errorMsg += `, did you misspell type name "%s" or forget importing the database driver? `
	return nil, gberror.NewCodef(gbcode.CodeInvalidConfiguration, errorMsg, config.Type, config.Type)
}

var (
	// driverMap manages all custom registered driver.
	driverMap = map[string]IDriver{}
)

const (
	DefaultGroupName = "default" // Default group name.
	defaultCharset   = `utf8`
	defaultProtocol  = `tcp`
)
