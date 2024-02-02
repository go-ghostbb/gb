package gbclient

import (
	"context"
	gbvar "ghostbb.io/gb/container/gb_var"
	"ghostbb.io/gb/internal/intlog"
	"net/http"
)

// GetVar sends a GET request, retrieves and converts the result content to *gbvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gbvar.Var can be conveniently converted to any type you want.
func (c *Client) GetVar(ctx context.Context, url string, data ...interface{}) *gbvar.Var {
	return c.RequestVar(ctx, http.MethodGet, url, data...)
}

// PutVar sends a PUT request, retrieves and converts the result content to *gbvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gbvar.Var can be conveniently converted to any type you want.
func (c *Client) PutVar(ctx context.Context, url string, data ...interface{}) *gbvar.Var {
	return c.RequestVar(ctx, http.MethodPut, url, data...)
}

// PostVar sends a POST request, retrieves and converts the result content to *gbvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gbvar.Var can be conveniently converted to any type you want.
func (c *Client) PostVar(ctx context.Context, url string, data ...interface{}) *gbvar.Var {
	return c.RequestVar(ctx, http.MethodPost, url, data...)
}

// DeleteVar sends a DELETE request, retrieves and converts the result content to *gbvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gbvar.Var can be conveniently converted to any type you want.
func (c *Client) DeleteVar(ctx context.Context, url string, data ...interface{}) *gbvar.Var {
	return c.RequestVar(ctx, http.MethodDelete, url, data...)
}

// HeadVar sends a HEAD request, retrieves and converts the result content to *gbvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gbvar.Var can be conveniently converted to any type you want.
func (c *Client) HeadVar(ctx context.Context, url string, data ...interface{}) *gbvar.Var {
	return c.RequestVar(ctx, http.MethodHead, url, data...)
}

// PatchVar sends a PATCH request, retrieves and converts the result content to *gbvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gbvar.Var can be conveniently converted to any type you want.
func (c *Client) PatchVar(ctx context.Context, url string, data ...interface{}) *gbvar.Var {
	return c.RequestVar(ctx, http.MethodPatch, url, data...)
}

// ConnectVar sends a CONNECT request, retrieves and converts the result content to *gbvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gbvar.Var can be conveniently converted to any type you want.
func (c *Client) ConnectVar(ctx context.Context, url string, data ...interface{}) *gbvar.Var {
	return c.RequestVar(ctx, http.MethodConnect, url, data...)
}

// OptionsVar sends an OPTIONS request, retrieves and converts the result content to *gbvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gbvar.Var can be conveniently converted to any type you want.
func (c *Client) OptionsVar(ctx context.Context, url string, data ...interface{}) *gbvar.Var {
	return c.RequestVar(ctx, http.MethodOptions, url, data...)
}

// TraceVar sends a TRACE request, retrieves and converts the result content to *gbvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gbvar.Var can be conveniently converted to any type you want.
func (c *Client) TraceVar(ctx context.Context, url string, data ...interface{}) *gbvar.Var {
	return c.RequestVar(ctx, http.MethodTrace, url, data...)
}

// RequestVar sends request using given HTTP method and data, retrieves converts the result to *gbvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gbvar.Var can be conveniently converted to any type you want.
func (c *Client) RequestVar(ctx context.Context, method string, url string, data ...interface{}) *gbvar.Var {
	response, err := c.DoRequest(ctx, method, url, data...)
	if err != nil {
		intlog.Errorf(ctx, `%+v`, err)
		return gbvar.New(nil)
	}
	defer func() {
		if err = response.Close(); err != nil {
			intlog.Errorf(ctx, `%+v`, err)
		}
	}()
	return gbvar.New(response.ReadAll())
}
