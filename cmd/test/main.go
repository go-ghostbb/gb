package main

import (
	"bytes"
	_ "github.com/Ghostbb-io/gb/contrib/drivers/mssql"
	_ "github.com/Ghostbb-io/gb/contrib/nosql/redis"
	gbdb "github.com/Ghostbb-io/gb/database/gb_db"
	gberror "github.com/Ghostbb-io/gb/errors/gb_error"
	"github.com/Ghostbb-io/gb/frame/g"
	gbcache "github.com/Ghostbb-io/gb/os/gb_cache"
	gbctx "github.com/Ghostbb-io/gb/os/gb_ctx"
	gbutil "github.com/Ghostbb-io/gb/util/gb_util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	db := g.DB()
	db.SetCacheAdapter(gbcache.NewAdapterRedis(g.Redis("dbCache")))
	server := g.Server()
	server.GET("/hello", func(c *gin.Context) {
		var u User
		err := db.WithContext(gbdb.WithCacheCtx()).Where("id = 1").First(&u).Error
		if err != nil && !gberror.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}
		buff := bytes.NewBuffer(nil)
		g.DumpTo(buff, u, gbutil.DumpOption{})
		c.String(200, buff.String())
	})

	g.Log().Info(gbctx.New(), "test")

	testServer := g.Server("test")
	g.RunMultiple(server, testServer)
}
