package pgsql

import (
	"fmt"
	gbdb "ghostbb.io/gb/database/gb_db"
	gbstr "ghostbb.io/gb/text/gb_str"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Driver struct{}

func (d *Driver) New(config gbdb.DatabaseConfig) (db *gorm.DB, err error) {
	var (
		source string
	)

	// host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
	// host=%s user=%s password=%s dbname=%s port=%s sslmode=disable
	source = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Host, config.User, config.Pass, config.Name, config.Port,
	)
	if config.Timezone != "" {
		source = fmt.Sprintf("%s timezone=%s", source, config.Timezone)
	}

	if config.Extra != "" {
		var extraMap map[string]interface{}
		if extraMap, err = gbstr.Parse(config.Extra); err != nil {
			return nil, err
		}
		for k, v := range extraMap {
			source += fmt.Sprintf(` %s=%s`, k, v)
		}
	}

	pgsqlConfig := postgres.Config{
		DSN:                  source,
		PreferSimpleProtocol: false,
	}

	if db, err = gorm.Open(postgres.New(pgsqlConfig), config.GormConfig()); err != nil {
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(config.MaxIdle)
		sqlDB.SetMaxOpenConns(config.MaxOpen)
		return db, nil
	}
}
