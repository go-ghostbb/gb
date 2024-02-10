package gbdb

import "gorm.io/gorm"

type IDriver interface {
	// New creates and returns a database object for specified database server.
	New(config DatabaseConfig) (*gorm.DB, error)
}

type DriverWrapper struct {
	driver IDriver
}

func (d *DriverWrapper) New(config DatabaseConfig) (*gorm.DB, error) {
	return d.driver.New(config)
}

// newDriverWrapper creates and returns a driver wrapper.
func newDriverWrapper(driver IDriver) IDriver {
	return &DriverWrapper{
		driver: driver,
	}
}
