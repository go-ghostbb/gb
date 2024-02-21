package gbhttp

import (
	"context"
	"crypto/tls"
	gbsvc "ghostbb.io/gb/net/gb_svc"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbres "ghostbb.io/gb/os/gb_res"
	"strconv"
	"time"
)

// SetDumpRouterMap sets the DumpRouterMap for server.
// If DumpRouterMap is enabled, it automatically dumps the route map when server starts.
func (s *Server) SetDumpRouterMap(enabled bool) {
	s.config.DumpRouterMap = enabled
}

// SetAddr sets the listening address for the server.
// The address is like ':80', '0.0.0.0:80', '127.0.0.1:80', '180.18.99.10:80', etc.
func (s *Server) SetAddr(address string) {
	s.config.Address = address
}

// SetHTTPSAddr sets the HTTPS listening ports for the server.
func (s *Server) SetHTTPSAddr(address string) {
	s.config.HTTPSAddr = address
}

// SetPort sets the listening ports for the server.
// The listening ports can be multiple like: SetPort(80, 8080).
func (s *Server) SetPort(port ...int) {
	if len(port) > 0 {
		s.config.Address = ""
		for _, v := range port {
			if len(s.config.Address) > 0 {
				s.config.Address += ","
			}
			s.config.Address += ":" + strconv.Itoa(v)
		}
	}
}

// SetHTTPSPort sets the HTTPS listening ports for the server.
// The listening ports can be multiple like: SetHTTPSPort(443, 500).
func (s *Server) SetHTTPSPort(port ...int) {
	if len(port) > 0 {
		s.config.HTTPSAddr = ""
		for _, v := range port {
			if len(s.config.HTTPSAddr) > 0 {
				s.config.HTTPSAddr += ","
			}
			s.config.HTTPSAddr += ":" + strconv.Itoa(v)
		}
	}
}

// EnableHTTPS enables HTTPS with given certification and key files for the server.
// The optional parameter `tlsConfig` specifies custom TLS configuration.
func (s *Server) EnableHTTPS(certFile, keyFile string, tlsConfig ...*tls.Config) {
	var ctx = context.TODO()
	certFileRealPath := gbfile.RealPath(certFile)
	if certFileRealPath == "" {
		certFileRealPath = gbfile.RealPath(gbfile.Pwd() + gbfile.Separator + certFile)
		if certFileRealPath == "" {
			certFileRealPath = gbfile.RealPath(gbfile.MainPkgPath() + gbfile.Separator + certFile)
		}
	}
	// Resource.
	if certFileRealPath == "" && gbres.Contains(certFile) {
		certFileRealPath = certFile
	}
	if certFileRealPath == "" {
		s.Logger().Fatalf(ctx, `EnableHTTPS failed: certFile "%s" does not exist`, certFile)
	}
	keyFileRealPath := gbfile.RealPath(keyFile)
	if keyFileRealPath == "" {
		keyFileRealPath = gbfile.RealPath(gbfile.Pwd() + gbfile.Separator + keyFile)
		if keyFileRealPath == "" {
			keyFileRealPath = gbfile.RealPath(gbfile.MainPkgPath() + gbfile.Separator + keyFile)
		}
	}
	// Resource.
	if keyFileRealPath == "" && gbres.Contains(keyFile) {
		keyFileRealPath = keyFile
	}
	if keyFileRealPath == "" {
		s.Logger().Fatal(ctx, `EnableHTTPS failed: keyFile "%s" does not exist`, keyFile)
	}
	s.config.HTTPSCertPath = certFileRealPath
	s.config.HTTPSKeyPath = keyFileRealPath
	if len(tlsConfig) > 0 {
		s.config.TLSConfig = tlsConfig[0]
	}
}

// SetTLSConfig sets custom TLS configuration and enables HTTPS feature for the server.
func (s *Server) SetTLSConfig(tlsConfig *tls.Config) {
	s.config.TLSConfig = tlsConfig
}

// SetReadTimeout sets the ReadTimeout for the server.
func (s *Server) SetReadTimeout(t time.Duration) {
	s.config.ReadTimeout = t
}

// SetWriteTimeout sets the WriteTimeout for the server.
func (s *Server) SetWriteTimeout(t time.Duration) {
	s.config.WriteTimeout = t
}

// SetIdleTimeout sets the IdleTimeout for the server.
func (s *Server) SetIdleTimeout(t time.Duration) {
	s.config.IdleTimeout = t
}

// SetMaxHeaderBytes sets the MaxHeaderBytes for the server.
func (s *Server) SetMaxHeaderBytes(b int) {
	s.config.MaxHeaderBytes = b
}

// SetKeepAlive sets the KeepAlive for the server.
func (s *Server) SetKeepAlive(enabled bool) {
	s.config.KeepAlive = enabled
}

func (s *Server) SetTerminal(enabled bool) {
	s.config.Terminal = enabled
}

// SetRegistrar sets the Registrar for server.
func (s *Server) SetRegistrar(registrar gbsvc.Registrar) {
	s.registrar = registrar
}

// SetEndpoints sets the Endpoints for the server.
func (s *Server) SetEndpoints(endpoints []string) {
	s.config.Endpoints = endpoints
}
