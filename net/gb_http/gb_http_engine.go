package gbhttp

import (
	gbmap "github.com/Ghostbb-io/gb/container/gb_map"
	gbctx "github.com/Ghostbb-io/gb/os/gb_ctx"
	gblog "github.com/Ghostbb-io/gb/os/gb_log"
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
