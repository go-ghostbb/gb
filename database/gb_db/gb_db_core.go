package gbdb

import (
	gblog "ghostbb.io/gb/os/gb_log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func (c *Core) GetLogger() *gblog.Logger {
	return c.logger.(*Logger).logger
}

func (c *Core) GetConfig() *ConfigNode {
	return c.config
}

func (c *Core) GormConfig() *gorm.Config {
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.config.TablePrefix,
			SingularTable: c.config.SingularTable,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	config.Logger = c.logger.LogMode(logger.Info)
	return config
}

func (c *Core) SetSlowThreshold(t time.Duration) {
	c.logger.SetSlowThreshold(t)
}

func (c *Core) SetIgnoreRecordNotFoundError(b bool) {
	c.logger.SetIgnoreRecordNotFoundError(b)
}
