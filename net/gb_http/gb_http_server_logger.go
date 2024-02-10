package gbhttp

import (
	"context"
	"fmt"
	"ghostbb.io/gb/internal/instance"
	gblog "ghostbb.io/gb/os/gb_log"
	gbstr "ghostbb.io/gb/text/gb_str"
	"github.com/gin-gonic/gin"
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
		err := gbstr.TrimRight(errorsString(c.Errors.ByType(gin.ErrorTypePrivate)), "\n")

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
			if !s.IsErrorLogEnabled() {
				return
			}
			logger = logger.File(s.config.ErrorLogPattern).Stack(false)
			logger.Error(ctx, msg+"\n    ", m)
			logger.Header(false).Error(context.Background(), err+"\n")
		} else {
			if !s.IsAccessLogEnabled() {
				return
			}
			logger.Info(ctx, msg+"\n    ", m, "\n")
		}
	}
}

func (s *Server) terminal() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !s.config.Terminal {
			return
		}

		var (
			start = time.Now()
			path  = c.Request.URL.Path
			raw   = c.Request.URL.RawQuery
		)

		// Process request
		c.Next()

		// Stop timer
		var (
			timeStamp    = time.Now()
			latency      = timeStamp.Sub(start).Truncate(time.Second)
			clientIP     = c.ClientIP()
			method       = c.Request.Method
			statusCode   = c.Writer.Status()
			errorMessage = gbstr.TrimRight(errorsString(c.Errors.ByType(gin.ErrorTypePrivate)), "\n")

			statusColor = statusCodeColor(statusCode)
			methodColor = methodColor(method)
			resetColor  = resetColor()
		)

		if raw != "" {
			path = path + "?" + raw
		}

		msg := fmt.Sprintf("%v [SERVER] %s |%s %3d %s| %13v | %15s |%s %-7s %s %#v",
			timeStamp.Format("2006/01/02 15:04:05"),
			s.instance,
			statusColor, statusCode, resetColor,
			latency,
			clientIP,
			methodColor, method, resetColor,
			path,
		)
		if errorMessage != "" {
			msg += "\n" + errorMessage + "\n"
		}
		fmt.Println(msg)
	}
}
