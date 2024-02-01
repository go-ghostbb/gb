// Package gbclient provides convenient http client functionalities.
package gbclient

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"github.com/Ghostbb-io/gb"
	gberror "github.com/Ghostbb-io/gb/errors/gb_error"
	gbsel "github.com/Ghostbb-io/gb/net/gb_sel"
	gbsvc "github.com/Ghostbb-io/gb/net/gb_svc"
	gbfile "github.com/Ghostbb-io/gb/os/gb_file"
	"net/http"
	"os"
	"time"
)

// Client is the HTTP client for HTTP request management.
type Client struct {
	http.Client                         // Underlying HTTP Client.
	header            map[string]string // Custom header map.
	cookies           map[string]string // Custom cookie map.
	prefix            string            // Prefix for request.
	authUser          string            // HTTP basic authentication: user.
	authPass          string            // HTTP basic authentication: pass.
	retryCount        int               // Retry count when request fails.
	noUrlEncode       bool              // No url encoding for request parameters.
	retryInterval     time.Duration     // Retry interval when request fails.
	middlewareHandler []HandlerFunc     // Interceptor handlers
	discovery         gbsvc.Discovery   // Discovery for service.
	builder           gbsel.Builder     // Builder for request balance.
}

const (
	httpProtocolName          = `http`
	httpParamFileHolder       = `@file:`
	httpRegexParamJson        = `^[\w\[\]]+=.+`
	httpRegexHeaderRaw        = `^([\w\-]+):\s*(.+)`
	httpHeaderHost            = `Host`
	httpHeaderCookie          = `Cookie`
	httpHeaderUserAgent       = `User-Agent`
	httpHeaderContentType     = `Content-Type`
	httpHeaderContentTypeJson = `application/json`
	httpHeaderContentTypeXml  = `application/xml`
	httpHeaderContentTypeForm = `application/x-www-form-urlencoded`
)

var (
	hostname, _        = os.Hostname()
	defaultClientAgent = fmt.Sprintf(`GClient %s at %s`, gb.VERSION, hostname)
)

// New creates and returns a new HTTP client object.
func New() *Client {
	c := &Client{
		Client: http.Client{
			Transport: &http.Transport{
				// No validation for https certification of the server in default.
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				DisableKeepAlives: true,
			},
		},
		header:    make(map[string]string),
		cookies:   make(map[string]string),
		builder:   gbsel.GetBuilder(),
		discovery: gbsvc.GetRegistry(),
	}
	c.header[httpHeaderUserAgent] = defaultClientAgent
	// It enables OpenTelemetry for client in default.
	c.Use(internalMiddlewareTracing, internalMiddlewareDiscovery)
	return c
}

// Clone deeply clones current client and returns a new one.
func (c *Client) Clone() *Client {
	newClient := New()
	*newClient = *c
	if len(c.header) > 0 {
		newClient.header = make(map[string]string)
		for k, v := range c.header {
			newClient.header[k] = v
		}
	}
	if len(c.cookies) > 0 {
		newClient.cookies = make(map[string]string)
		for k, v := range c.cookies {
			newClient.cookies[k] = v
		}
	}
	return newClient
}

// LoadKeyCrt creates and returns a TLS configuration object with given certificate and key files.
func LoadKeyCrt(crtFile, keyFile string) (*tls.Config, error) {
	crtPath, err := gbfile.Search(crtFile)
	if err != nil {
		return nil, err
	}
	keyPath, err := gbfile.Search(keyFile)
	if err != nil {
		return nil, err
	}
	crt, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		err = gberror.Wrapf(err, `tls.LoadX509KeyPair failed for certFile "%s", keyFile "%s"`, crtPath, keyPath)
		return nil, err
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	tlsConfig.Time = time.Now
	tlsConfig.Rand = rand.Reader
	return tlsConfig, nil
}
