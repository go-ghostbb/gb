package mssql

import (
	"fmt"
	gbdb "ghostbb.io/gb/database/gb_db"
	gbstr "ghostbb.io/gb/text/gb_str"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func init() {
	if err := gbdb.Register("mssql", New()); err != nil {
		panic(err)
	}
}

type Driver struct{}

func (d *Driver) New(config gbdb.DatabaseConfig) (db *gorm.DB, err error) {
	var (
		source = config.Link
	)
	if source == "" {
		source = fmt.Sprintf(
			"sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable",
			config.User, config.Pass, config.Host, config.Port, config.Name,
		)
		source = gbstr.Replace(source, "\\", `\`)
		if config.Extra != "" {
			var extraMap map[string]interface{}
			if extraMap, err = gbstr.Parse(config.Extra); err != nil {
				return nil, err
			}
			for k, v := range extraMap {
				source += fmt.Sprintf(`&%s=%s`, k, v)
			}
		}
	}

	mssqlConfig := sqlserver.Config{
		DSN:               source, // DSN data source name
		DefaultStringSize: 191,    // string 類型字段默認長度
	}

	if db, err = gorm.Open(sqlserver.New(mssqlConfig), config.GormConfig()); err != nil {
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(config.MaxIdle)
		sqlDB.SetMaxOpenConns(config.MaxOpen)
		return db, nil
	}

}

func New() gbdb.IDriver {
	return &Driver{}
}
