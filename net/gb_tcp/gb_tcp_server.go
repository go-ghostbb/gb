package gbtcp

import (
	"crypto/tls"
	"fmt"
	gbmap "ghostbb.io/gb/container/gb_map"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"net"
	"sync"
)

const (
	// FreePortAddress marks the server listens using random free port.
	FreePortAddress = ":0"
)

const (
	defaultServer = "default"
)

// Server is a TCP server.
type Server struct {
	mu        sync.Mutex   // Used for Server.listen concurrent safety. -- The golang test with data race checks this.
	listen    net.Listener // TCP address listener.
	address   string       // Server listening address.
	handler   func(*Conn)  // Connection handler.
	tlsConfig *tls.Config  // TLS configuration.
}

// Map for name to server, for singleton purpose.
var serverMapping = gbmap.NewStrAnyMap(true)

// GetServer returns the TCP server with specified `name`,
// or it returns a new normal TCP server named `name` if it does not exist.
// The parameter `name` is used to specify the TCP server
func GetServer(name ...interface{}) *Server {
	serverName := defaultServer
	if len(name) > 0 && name[0] != "" {
		serverName = gbconv.String(name[0])
	}
	return serverMapping.GetOrSetFuncLock(serverName, func() interface{} {
		return NewServer("", nil)
	}).(*Server)
}

// NewServer creates and returns a new normal TCP server.
// The parameter `name` is optional, which is used to specify the instance name of the server.
func NewServer(address string, handler func(*Conn), name ...string) *Server {
	s := &Server{
		address: address,
		handler: handler,
	}
	if len(name) > 0 && name[0] != "" {
		serverMapping.Set(name[0], s)
	}
	return s
}

// NewServerTLS creates and returns a new TCP server with TLS support.
// The parameter `name` is optional, which is used to specify the instance name of the server.
func NewServerTLS(address string, tlsConfig *tls.Config, handler func(*Conn), name ...string) *Server {
	s := NewServer(address, handler, name...)
	s.SetTLSConfig(tlsConfig)
	return s
}

// NewServerKeyCrt creates and returns a new TCP server with TLS support.
// The parameter `name` is optional, which is used to specify the instance name of the server.
func NewServerKeyCrt(address, crtFile, keyFile string, handler func(*Conn), name ...string) (*Server, error) {
	s := NewServer(address, handler, name...)
	if err := s.SetTLSKeyCrt(crtFile, keyFile); err != nil {
		return nil, err
	}
	return s, nil
}

// SetAddress sets the listening address for server.
func (s *Server) SetAddress(address string) {
	s.address = address
}

// GetAddress get the listening address for server.
func (s *Server) GetAddress() string {
	return s.address
}

// SetHandler sets the connection handler for server.
func (s *Server) SetHandler(handler func(*Conn)) {
	s.handler = handler
}

// SetTLSKeyCrt sets the certificate and key file for TLS configuration of server.
func (s *Server) SetTLSKeyCrt(crtFile, keyFile string) error {
	tlsConfig, err := LoadKeyCrt(crtFile, keyFile)
	if err != nil {
		return err
	}
	s.tlsConfig = tlsConfig
	return nil
}

// SetTLSConfig sets the TLS configuration of server.
func (s *Server) SetTLSConfig(tlsConfig *tls.Config) {
	s.tlsConfig = tlsConfig
}

// Close closes the listener and shutdowns the server.
func (s *Server) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listen == nil {
		return nil
	}
	return s.listen.Close()
}

// Run starts running the TCP Server.
func (s *Server) Run() (err error) {
	if s.handler == nil {
		err = gberror.NewCode(gbcode.CodeMissingConfiguration, "start running failed: socket handler not defined")
		return
	}
	if s.tlsConfig != nil {
		// TLS Server
		s.mu.Lock()
		s.listen, err = tls.Listen("tcp", s.address, s.tlsConfig)
		s.mu.Unlock()
		if err != nil {
			err = gberror.Wrapf(err, `tls.Listen failed for address "%s"`, s.address)
			return
		}
	} else {
		// Normal Server
		var tcpAddr *net.TCPAddr
		if tcpAddr, err = net.ResolveTCPAddr("tcp", s.address); err != nil {
			err = gberror.Wrapf(err, `net.ResolveTCPAddr failed for address "%s"`, s.address)
			return err
		}
		s.mu.Lock()
		s.listen, err = net.ListenTCP("tcp", tcpAddr)
		s.mu.Unlock()
		if err != nil {
			err = gberror.Wrapf(err, `net.ListenTCP failed for address "%s"`, s.address)
			return err
		}
	}
	// Listening loop.
	for {
		var conn net.Conn
		if conn, err = s.listen.Accept(); err != nil {
			err = gberror.Wrapf(err, `Listener.Accept failed`)
			return err
		} else if conn != nil {
			go s.handler(NewConnByNetConn(conn))
		}
	}
}

// GetListenedAddress retrieves and returns the address string which are listened by current server.
func (s *Server) GetListenedAddress() string {
	if !gbstr.Contains(s.address, FreePortAddress) {
		return s.address
	}
	var (
		address      = s.address
		listenedPort = s.GetListenedPort()
	)
	address = gbstr.Replace(address, FreePortAddress, fmt.Sprintf(`:%d`, listenedPort))
	return address
}

// GetListenedPort retrieves and returns one port which is listened to by current server.
func (s *Server) GetListenedPort() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ln := s.listen; ln != nil {
		return ln.Addr().(*net.TCPAddr).Port
	}
	return -1
}
