package gbhttp

import (
	gbmap "ghostbb.io/container/gb_map"
	gbctx "ghostbb.io/os/gb_ctx"
	gblog "ghostbb.io/os/gb_log"
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
