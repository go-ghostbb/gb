package gbhttp

import (
	"context"
	"fmt"
	gbvar "ghostbb.io/gb/container/gb_var"
	gberror "ghostbb.io/gb/errors/gb_error"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strings"
)

// GetBearerToken Get bearer token from header (Authorization: Bearer xxx)
func GetBearerToken(c *gin.Context) string {
	var token string
	auth := c.GetHeader("Authorization")
	prefix := "Bearer "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

// ParseJSON Parse body json data to struct
func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return gberror.New(fmt.Sprintf("Parse request json failed: %s", err.Error()))
	}
	return nil
}

// ParseParam Param returns the value of the URL param
func ParseParam(c *gin.Context, key string) string {
	val := c.Param(key)
	return val
}

// ParseQuery Parse query parameter to struct
func ParseQuery(c *gin.Context, obj any) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return gberror.New(fmt.Sprintf("Parse request query failed: %s", err.Error()))
	}
	return nil
}

// ParseForm Parse body form data to struct
func ParseForm(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		return gberror.New(fmt.Sprintf("Parse request form failed: %s", err.Error()))
	}
	return nil
}

func Ctx(c *gin.Context) context.Context {
	return Get(c, ServerContextKey).Interface().(context.Context)
}

func Set(c *gin.Context, key string, value any) {
	c.Set(key, value)
}

func Get(c *gin.Context, key string) *gbvar.Var {
	value, _ := c.Get(key)
	return gbvar.New(value)
}
