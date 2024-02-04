package dbcache_test

import (
	_ "ghostbb.io/gb/contrib/drivers/mssql"

	"ghostbb.io/gb/contrib/dbcache"
	"ghostbb.io/gb/frame/g"
	"gorm.io/gorm"
	"testing"
)

func TestExample(t *testing.T) {
	type User struct {
		gorm.Model
		Username string
		Password string
	}

	db := g.DB("test")
	db.Use(dbcache.New())
	u := new(User)
	db.WithContext(dbcache.WithLevelAllCtx()).Where("id = 1").First(u)
}
