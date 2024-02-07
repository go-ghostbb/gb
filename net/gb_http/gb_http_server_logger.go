package gbhttp

import (
	"context"
	"fmt"
	"ghostbb.io/gb/internal/instance"
	gblog "ghostbb.io/gb/os/gb_log"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// Logger is alias of GetLogger.
func (s *Server) Logger() *gblog.Logger {
	return s.config.Logger
}

func (s *Server) loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			loggerInstanceKey = fmt.Sprintf(`Logger Of Server:%s`, s.instance)
			start             = time.Now()
			path              = c.Request.URL.Path
			method            = c.Request.Method
			m                 = map[string]interface{}{
				"query":     c.Request.URL.RawQuery,
				"ip":        c.ClientIP(),
				"userAgent": c.Request.UserAgent(),
			}
			ctx = context.Background()
		)
		c.Next()
		m["body"], _ = c.GetRawData()
		err := strings.TrimRight(c.Errors.ByType(gin.ErrorTypePrivate).String(), "\n")

		logger := instance.GetOrSetFuncLock(loggerInstanceKey, func() interface{} {
			l := s.config.Logger.Clone()
			l.SetFile(s.config.AccessLogPattern)
			l.SetStdoutPrint(s.config.LogStdout)
			l.SetLevelPrint(false)
			return l
		}).(*gblog.Logger)

		if ctxValue, ok := c.Get(ServerContextKey); ok {
			if v, ok1 := ctxValue.(context.Context); ok1 {
				ctx = v
			}
		}

		msg := fmt.Sprintf("【%s】｜%d｜%13v｜%s ==> \"%s\"",
			s.instance,
			c.Writer.Status(),
			time.Since(start),
			method,
			path,
		)
		if err != "" {
			msg += fmt.Sprintf("｜ERROR｜%s", err)
			logger.Error(ctx, msg+"\n    ", m, "\n")
		} else {
			logger.Info(ctx, msg+"\n    ", m, "\n")
		}
	}
}

func (s *Server) debugLog(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	statusColor = param.StatusCodeColor()
	methodColor = param.MethodColor()
	resetColor = param.ResetColor()

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	return fmt.Sprintf("%v [SERVER] %s |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		param.TimeStamp.Format("2006/01/02 15:04:05"),
		s.instance,
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}
