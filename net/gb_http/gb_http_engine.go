package gbhttp

import (
	gbmap "ghostbb.io/gb/container/gb_map"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gblog "ghostbb.io/gb/os/gb_log"
	"github.com/gin-gonic/gin"
)

func Default() *Engine {
	return &Engine{
		Engine:       gin.New(),
		groupMapping: gbmap.NewStrAnyMap(true),
		logger:       gblog.New(),
	}
}

func (e *Engine) GBCtxMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ServerContextKey, gbctx.New())
	}
}
