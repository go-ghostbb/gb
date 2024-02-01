package gbdb

import "gorm.io/gorm"

// DriverWrapper is a driver wrapper for extending features with embedded driver.
type DriverWrapper struct {
	driver Driver
}

func (d *DriverWrapper) New(core *Core, node *ConfigNode) (*gorm.DB, error) {
	return d.driver.New(core, node)
}

// newDriverWrapper creates and returns a driver wrapper.
func newDriverWrapper(driver Driver) Driver {
	return &DriverWrapper{
		driver: driver,
	}
}
