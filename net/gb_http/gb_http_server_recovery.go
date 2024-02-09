package gbhttp

import (
	gberror "ghostbb.io/gb/errors/gb_error"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) Recovery() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var gberr error
				switch v := err.(type) {
				case *gberror.Error:
					gberr = v
				default:
					gberr = gberror.NewSkip(1, gbconv.String(v))
				}
				_ = c.Error(gberr)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next() // 调用下一个处理
	}
}
