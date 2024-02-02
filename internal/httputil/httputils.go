// Package httputil provides HTTP functions for internal usage only.
package httputil

import (
	gburl "ghostbb.io/gb/encoding/gb_url"
	"ghostbb.io/gb/internal/empty"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"net/http"
	"strings"
)

const (
	fileUploadingKey = "@file:"
)

// BuildParams builds the request string for the http client. The `params` can be type of:
// string/[]byte/map/struct/*struct.
//
// The optional parameter `noUrlEncode` specifies whether ignore the url encoding for the data.
func BuildParams(params interface{}, noUrlEncode ...bool) (encodedParamStr string) {
	// If given string/[]byte, converts and returns it directly as string.
	switch v := params.(type) {
	case string, []byte:
		return gbconv.String(params)
	case []interface{}:
		if len(v) > 0 {
			params = v[0]
		} else {
			params = nil
		}
	}
	// Else converts it to map and does the url encoding.
	m, urlEncode := gbconv.Map(params), true
	if len(m) == 0 {
		return gbconv.String(params)
	}
	if len(noUrlEncode) == 1 {
		urlEncode = !noUrlEncode[0]
	}
	// If there's file uploading, it ignores the url encoding.
	if urlEncode {
		for k, v := range m {
			if gbstr.Contains(k, fileUploadingKey) || gbstr.Contains(gbconv.String(v), fileUploadingKey) {
				urlEncode = false
				break
			}
		}
	}
	s := ""
	for k, v := range m {
		// Ignore nil attributes.
		if empty.IsNil(v) {
			continue
		}
		if len(encodedParamStr) > 0 {
			encodedParamStr += "&"
		}
		s = gbconv.String(v)
		if urlEncode {
			if strings.HasPrefix(s, fileUploadingKey) && len(s) > len(fileUploadingKey) {
				// No url encoding if uploading file.
			} else {
				s = gburl.Encode(s)
			}
		}
		encodedParamStr += k + "=" + s
	}
	return
}

// HeaderToMap coverts request headers to map.
func HeaderToMap(header http.Header) map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range header {
		if len(v) > 1 {
			m[k] = v
		} else if len(v) == 1 {
			m[k] = v[0]
		}
	}
	return m
}
