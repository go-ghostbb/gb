package gbcfg

import (
	"context"
	gbvar "ghostbb.io/gb/container/gb_var"
	gbjson "ghostbb.io/gb/encoding/gb_json"
	gberror "ghostbb.io/gb/errors/gb_error"
)

// AdapterContent implements interface Adapter using content.
// The configuration content supports the coding types as package `gbjson`.
type AdapterContent struct {
	jsonVar *gbvar.Var // The pared JSON object for configuration content, type: *gbjson.Json.
}

// NewAdapterContent returns a new configuration management object using custom content.
// The parameter `content` specifies the default configuration content for reading.
func NewAdapterContent(content ...string) (*AdapterContent, error) {
	a := &AdapterContent{
		jsonVar: gbvar.New(nil, true),
	}
	if len(content) > 0 {
		if err := a.SetContent(content[0]); err != nil {
			return nil, err
		}
	}
	return a, nil
}

// SetContent sets customized configuration content for specified `file`.
// The `file` is unnecessary param, default is DefaultConfigFile.
func (a *AdapterContent) SetContent(content string) error {
	j, err := gbjson.LoadContent(content, true)
	if err != nil {
		return gberror.Wrap(err, `load configuration content failed`)
	}
	a.jsonVar.Set(j)
	return nil
}

// Available checks and returns the backend configuration service is available.
// The optional parameter `resource` specifies certain configuration resource.
//
// Note that this function does not return error as it just does simply check for
// backend configuration service.
func (a *AdapterContent) Available(ctx context.Context, resource ...string) (ok bool) {
	return !a.jsonVar.IsNil()
}

// Get retrieves and returns value by specified `pattern` in current resource.
// Pattern like:
// "x.y.z" for map item.
// "x.0.y" for slice item.
func (a *AdapterContent) Get(ctx context.Context, pattern string) (value interface{}, err error) {
	if a.jsonVar.IsNil() {
		return nil, nil
	}
	return a.jsonVar.Val().(*gbjson.Json).Get(pattern).Val(), nil
}

// Data retrieves and returns all configuration data in current resource as map.
// Note that this function may lead lots of memory usage if configuration data is too large,
// you can implement this function if necessary.
func (a *AdapterContent) Data(ctx context.Context) (data map[string]interface{}, err error) {
	if a.jsonVar.IsNil() {
		return nil, nil
	}
	return a.jsonVar.Val().(*gbjson.Json).Var().Map(), nil
}
