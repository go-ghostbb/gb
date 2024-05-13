package mssql

import (
	"fmt"
	gbdb "ghostbb.io/gb/database/gb_db"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"net/url"
)

func init() {
	if err := gbdb.Register("mssql", New()); err != nil {
		panic(err)
	}
}

type Driver struct{}

func (d *Driver) New(config gbdb.DatabaseConfig) (db *gorm.DB, err error) {
	var (
		source      *url.URL
		query       = url.Values{}
		mssqlConfig sqlserver.Config
	)

	source = &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(config.User, config.Pass),
	}

	query.Add("database", config.Name)

	if config.Extra != "" {
		var extraMap map[string]interface{}
		if extraMap, err = gbstr.Parse(config.Extra); err != nil {
			return nil, err
		}
		for k, v := range extraMap {
			query.Add(k, gbconv.String(v))
		}

		source.RawQuery = query.Encode()
	}

	if config.Instance == "" {
		source.Host = fmt.Sprintf("%s:%s", config.Host, config.Port)
	} else {
		source.Host = config.Host
		source.Path = config.Instance
	}

	mssqlConfig = sqlserver.Config{
		DSN:               source.String(), // DSN data source name
		DefaultStringSize: 191,             // string 類型字段默認長度
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
