package gbhttp

import (
	"fmt"
	gberror "ghostbb.io/gb/errors/gb_error"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

func errorsString(errs []*gin.Error) string {
	if len(errs) == 0 {
		return ""
	}
	var buffer strings.Builder
	for i, msg := range errs {
		switch v := msg.Err.(type) {
		case *gberror.Error:
			fmt.Fprintf(&buffer, "Error #%02d: %+v\n", i+1, v)
		default:
			fmt.Fprintf(&buffer, "Error #%02d: %s\n", i+1, msg.Err)
		}

		if msg.Meta != nil {
			fmt.Fprintf(&buffer, "     Meta: %v\n", msg.Meta)
		}
	}
	return buffer.String()
}

func statusCodeColor(statusCode int) string {
	switch {
	case statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices:
		return green
	case statusCode >= http.StatusMultipleChoices && statusCode < http.StatusBadRequest:
		return white
	case statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// methodColor is the ANSI color for appropriately logging http method to a terminal.
func methodColor(method string) string {
	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

func resetColor() string {
	return reset
}
