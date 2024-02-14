package mysql

import (
	"fmt"
	gbdb "ghostbb.io/gb/database/gb_db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/url"
	"strings"
)

func init() {
	if err := gbdb.Register("mysql", New()); err != nil {
		panic(err)
	}
}

type Driver struct{}

func (d *Driver) New(config gbdb.DatabaseConfig) (db *gorm.DB, err error) {
	var (
		source string
	)
	source = fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s?charset=%s",
		config.User, config.Pass, config.Protocol, config.Host, config.Port, config.Name, config.Charset,
	)
	if config.Timezone != "" {
		if strings.Contains(config.Timezone, "/") {
			config.Timezone = url.QueryEscape(config.Timezone)
		}
		source = fmt.Sprintf("%s&loc=%s", source, config.Timezone)
	}
	if config.Extra != "" {
		source = fmt.Sprintf("%s&%s", source, config.Extra)
	}

	mysqlConfig := mysql.Config{
		DSN:               source, // DSN data source name
		DefaultStringSize: 191,    // string 類型字段默認長度
	}

	if db, err = gorm.Open(mysql.New(mysqlConfig), config.GormConfig()); err != nil {
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
