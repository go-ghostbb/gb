package gbclient

import (
	"fmt"
	"ghostbb.io/internal/intlog"
	"ghostbb.io/internal/utils"
	"io"
	"net/http"
	"net/http/httputil"
)

// dumpTextFormat is the format of the dumped raw string
const dumpTextFormat = `+---------------------------------------------+
|                   %s                  |
+---------------------------------------------+
%s
%s
`

// getResponseBody returns the text of the response body.
func getResponseBody(res *http.Response) string {
	if res.Body == nil {
		return ""
	}
	bodyContent, _ := io.ReadAll(res.Body)
	res.Body = utils.NewReadCloser(bodyContent, true)
	return string(bodyContent)
}

// RawRequest returns the raw content of the request.
func (r *Response) RawRequest() string {
	// Response can be nil.
	if r == nil || r.request == nil {
		return ""
	}
	// DumpRequestOut writes more request headers than DumpRequest, such as User-Agent.
	bs, err := httputil.DumpRequestOut(r.request, false)
	if err != nil {
		intlog.Errorf(r.request.Context(), `%+v`, err)
		return ""
	}
	return fmt.Sprintf(
		dumpTextFormat,
		"REQUEST ",
		string(bs),
		r.requestBody,
	)
}

// RawResponse returns the raw content of the response.
func (r *Response) RawResponse() string {
	// Response might be nil.
	if r == nil || r.Response == nil {
		return ""
	}
	bs, err := httputil.DumpResponse(r.Response, false)
	if err != nil {
		intlog.Errorf(r.request.Context(), `%+v`, err)
		return ""
	}

	return fmt.Sprintf(
		dumpTextFormat,
		"RESPONSE",
		string(bs),
		getResponseBody(r.Response),
	)
}

// Raw returns the raw text of the request and the response.
func (r *Response) Raw() string {
	return fmt.Sprintf("%s\n%s", r.RawRequest(), r.RawResponse())
}

// RawDump outputs the raw text of the request and the response to stdout.
func (r *Response) RawDump() {
	fmt.Println(r.Raw())
}
