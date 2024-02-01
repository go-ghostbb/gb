package mssql

import (
	"fmt"
	gbdb "github.com/Ghostbb-io/gb/database/gb_db"
	gbregex "github.com/Ghostbb-io/gb/text/gb_regex"
	gbstr "github.com/Ghostbb-io/gb/text/gb_str"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func init() {
	if err := gbdb.Register("mssql", New()); err != nil {
		panic(err)
	}
}

type Driver struct {
	*gbdb.Core
}

func (d *Driver) New(core *gbdb.Core, node *gbdb.ConfigNode) (db *gorm.DB, err error) {
	var (
		source string
		config = core.GetConfig()
	)
	if config.Link != "" {
		source = config.Link
		// Custom changing the schema in runtime.
		if config.Name != "" {
			source, _ = gbregex.ReplaceString(`database=([\w\.\-]+)+`, "database="+config.Name, source)
		}
	} else {
		source = fmt.Sprintf(
			"sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable",
			config.User, config.Pass, config.Host, config.Port, config.Name,
		)
		if config.Extra != "" {
			var extraMap map[string]interface{}
			if extraMap, err = gbstr.Parse(config.Extra); err != nil {
				return nil, err
			}
			for k, v := range extraMap {
				source += fmt.Sprintf(`;%s=%s`, k, v)
			}
		}

	}
	mssqlConfig := sqlserver.Config{
		DSN:               source, // DSN data source name
		DefaultStringSize: 191,    // string 類型字段默認長度
	}

	if db, err := gorm.Open(sqlserver.New(mssqlConfig), core.GormConfig()); err != nil {
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(config.MaxIdleConnCount)
		sqlDB.SetMaxOpenConns(config.MaxOpenConnCount)
		return db, nil
	}

}

func New() gbdb.Driver {
	return &Driver{}
}
