package gbhttp

import (
	"context"
	"crypto/tls"
	"fmt"
	gbtype "ghostbb.io/gb/container/gb_type"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbproc "ghostbb.io/gb/os/gb_proc"
	gbres "ghostbb.io/gb/os/gb_res"
	gbstr "ghostbb.io/gb/text/gb_str"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type internalServer struct {
	server      *Server      // Belonged server.
	fd          uintptr      // File descriptor for passing to the child process when graceful reload.
	address     string       // Listening address like:":80", ":8080".
	httpServer  *http.Server // Underlying http.Server.
	rawListener net.Listener // Underlying net.Listener.
	rawLnMu     sync.RWMutex // Concurrent safety mutex for `rawListener`.
	isHttps     bool         // Is HTTPS.
	status      *gbtype.Int  // Status of current server. Using `gbtype` to ensure concurrent safety.
	listener    net.Listener // Wrapped net.Listener.
}

func (s *Server) newInternalServer(address string, fd ...int) *internalServer {
	// Change port to address like: 80 -> :80
	if gbstr.IsNumeric(address) {
		address = ":" + address
	}

	is := &internalServer{
		server:     s,
		address:    address,
		httpServer: s.newHttpServer(address),
		status:     gbtype.NewInt(),
	}

	if len(fd) > 0 && fd[0] > 0 {
		is.fd = uintptr(fd[0])
	}
	if s.config.Listeners != nil {
		addrArray := gbstr.SplitAndTrim(address, ":")
		addrPort, err := strconv.Atoi(addrArray[len(addrArray)-1])
		if err == nil {
			for _, v := range s.config.Listeners {
				if listenerPort := (v.Addr().(*net.TCPAddr)).Port; listenerPort == addrPort {
					is.rawListener = v
					break
				}
			}
		}
	}
	return is
}

// newHttpServer creates and returns an underlying http.Server with a given address.
func (s *Server) newHttpServer(address string) *http.Server {
	server := &http.Server{
		Addr:           address,
		Handler:        s.handler,
		ReadTimeout:    s.config.ReadTimeout,
		WriteTimeout:   s.config.WriteTimeout,
		IdleTimeout:    s.config.IdleTimeout,
		MaxHeaderBytes: s.config.MaxHeaderBytes,
	}
	server.SetKeepAlivesEnabled(s.config.KeepAlive)
	return server
}

// Fd retrieves and returns the file descriptor of the current server.
// It is available ony in *nix like operating systems like linux, unix, darwin.
func (s *internalServer) Fd() uintptr {
	if ln := s.getRawListener(); ln != nil {
		file, err := ln.(*net.TCPListener).File()
		if err == nil {
			return file.Fd()
		}
	}
	return 0
}

// CreateListenerTLS creates listener on configured address with HTTPS.
// The parameter `certFile` and `keyFile` specify the necessary certification and key files for HTTPS.
// The optional parameter `tlsConfig` specifies the custom TLS configuration.
func (s *internalServer) CreateListenerTLS(certFile, keyFile string, tlsConfig ...*tls.Config) error {
	var config *tls.Config
	if len(tlsConfig) > 0 && tlsConfig[0] != nil {
		config = tlsConfig[0]
	} else if s.httpServer.TLSConfig != nil {
		config = s.httpServer.TLSConfig
	} else {
		config = &tls.Config{}
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}
	var err error
	if len(config.Certificates) == 0 {
		config.Certificates = make([]tls.Certificate, 1)
		if gbres.Contains(certFile) && gbres.Contains(keyFile) {
			config.Certificates[0], err = tls.X509KeyPair(
				gbres.GetContent(certFile),
				gbres.GetContent(keyFile),
			)
		} else {
			config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		}
	}
	if err != nil {
		return gberror.Wrapf(err, `open certFile "%s" and keyFile "%s" failed`, certFile, keyFile)
	}
	ln, err := s.getNetListener()
	if err != nil {
		return err
	}

	s.listener = tls.NewListener(ln, config)
	s.setRawListener(ln)
	return nil
}

// getNetListener retrieves and returns the wrapped net.Listener.
func (s *internalServer) getNetListener() (net.Listener, error) {
	if s.rawListener != nil {
		return s.rawListener, nil
	}

	var (
		ln  net.Listener
		err error
	)
	if s.fd > 0 {
		f := os.NewFile(s.fd, "")
		ln, err = net.FileListener(f)
		if err != nil {
			err = gberror.Wrap(err, "net.FileListener failed")
			return nil, err
		}
	} else {
		ln, err = net.Listen("tcp", s.httpServer.Addr)
		if err != nil {
			err = gberror.Wrapf(err, `net.Listen address "%s" failed`, s.httpServer.Addr)
		}
	}

	return ln, err
}

// setRawListener sets `rawListener` with given net.Listener.
func (s *internalServer) setRawListener(ln net.Listener) {
	s.rawLnMu.Lock()
	defer s.rawLnMu.Unlock()
	s.rawListener = ln
}

// setRawListener returns the `rawListener` of current server.
func (s *internalServer) getRawListener() net.Listener {
	s.rawLnMu.RLock()
	defer s.rawLnMu.RUnlock()
	return s.rawListener
}

// Serve starts the serving with blocking way.
func (s *internalServer) Serve(ctx context.Context) error {
	if s.rawListener == nil {
		return gberror.NewCode(gbcode.CodeInvalidOperation, `call CreateListener/CreateListenerTLS before Serve`)
	}

	action := "started"
	if s.fd != 0 {
		action = "reloaded"
	}

	s.server.Logger().Infof(
		ctx,
		`pid[%d]: %s server %s listening on [%s]`,
		gbproc.Pid(), s.getProto(), action, s.GetListenedAddress(),
	)
	s.status.Set(ServerStatusRunning)
	err := s.httpServer.Serve(s.listener)
	s.status.Set(ServerStatusStopped)
	return err
}

// GetListenedAddress retrieves and returns the address string which are listened by current server.
func (s *internalServer) GetListenedAddress() string {
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

// getProto retrieves and returns the proto string of current server.
func (s *internalServer) getProto() string {
	proto := "http"
	if s.isHttps {
		proto = "https"
	}
	return proto
}

// GetListenedPort retrieves and returns one port which is listened to by current server.
// Note that this method is only available if the server is listening on one port.
func (s *internalServer) GetListenedPort() int {
	if ln := s.getRawListener(); ln != nil {
		return ln.Addr().(*net.TCPAddr).Port
	}
	return -1
}

// CreateListener creates listener on configured address.
func (s *internalServer) CreateListener() error {
	ln, err := s.getNetListener()
	if err != nil {
		return err
	}
	s.listener = ln
	s.setRawListener(ln)
	return nil
}

// shutdown shuts down the server gracefully.
func (s *internalServer) shutdown(ctx context.Context) {
	if s.status.Val() == ServerStatusStopped {
		return
	}
	timeoutCtx, cancelFunc := context.WithTimeout(
		ctx,
		time.Duration(s.server.config.GracefulShutdownTimeout)*time.Second,
	)
	defer cancelFunc()
	if err := s.httpServer.Shutdown(timeoutCtx); err != nil {
		s.server.Logger().Errorf(
			ctx,
			"%d: %s server [%s] shutdown error: %v",
			gbproc.Pid(), s.getProto(), s.address, err,
		)
	}
}
