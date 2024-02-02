package gbclient

import (
	"context"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbregex "ghostbb.io/gb/text/gb_regex"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbmeta "ghostbb.io/gb/util/gb_meta"
	gbtag "ghostbb.io/gb/util/gb_tag"
	gbutil "ghostbb.io/gb/util/gb_util"
	"net/http"
	"reflect"
)

// DoRequestObj does HTTP request using standard request/response object.
// The request object `req` is defined like:
//
//	type UseCreateReq struct {
//	    g.Meta `path:"/user" method:"put"`
//	    // other fields....
//	}
//
// The response object `res` should be a pointer type. It automatically converts result
// to given object `res` is success.
//
// Example:
// var (
//
//	req = UseCreateReq{}
//	res *UseCreateRes
//
// )
//
// err := DoRequestObj(ctx, req, &res)
func (c *Client) DoRequestObj(ctx context.Context, req, res interface{}) error {
	var (
		method = gbmeta.Get(req, gbtag.Method).String()
		path   = gbmeta.Get(req, gbtag.Path).String()
	)
	if method == "" {
		return gberror.NewCodef(
			gbcode.CodeInvalidParameter,
			`no "%s" tag found in request object: %s`,
			gbtag.Method, reflect.TypeOf(req).String(),
		)
	}
	if path == "" {
		return gberror.NewCodef(
			gbcode.CodeInvalidParameter,
			`no "%s" tag found in request object: %s`,
			gbtag.Path, reflect.TypeOf(req).String(),
		)
	}
	path = c.handlePathForObjRequest(path, req)
	switch gbstr.ToUpper(method) {
	case
		http.MethodGet,
		http.MethodPut,
		http.MethodPost,
		http.MethodDelete,
		http.MethodHead,
		http.MethodPatch,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace:
		if result := c.RequestVar(ctx, method, path, req); res != nil && !result.IsEmpty() {
			return result.Scan(res)
		}
		return nil

	default:
		return gberror.Newf(`invalid HTTP method "%s"`, method)
	}
}

// handlePathForObjRequest replaces parameters in `path` with parameters from request object.
// Eg:
// /order/{id}  -> /order/1
// /user/{name} -> /order/john
func (c *Client) handlePathForObjRequest(path string, req interface{}) string {
	if gbstr.Contains(path, "{") {
		requestParamsMap := gbconv.Map(req)
		if len(requestParamsMap) > 0 {
			path, _ = gbregex.ReplaceStringFuncMatch(`\{(\w+)\}`, path, func(match []string) string {
				foundKey, foundValue := gbutil.MapPossibleItemByKey(requestParamsMap, match[1])
				if foundKey != "" {
					return gbconv.String(foundValue)
				}
				return match[0]
			})
		}
	}
	return path
}
