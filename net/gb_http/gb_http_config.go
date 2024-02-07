package gbhttp

import (
	"context"
	"crypto/tls"
	"ghostbb.io/gb/internal/intlog"
	gblog "ghostbb.io/gb/os/gb_log"
	gbstr "ghostbb.io/gb/text/gb_str"
	"net"
	"time"
)

const (
	defaultHttpAddr  = ":80"  // Default listening port for HTTP.
	defaultHttpsAddr = ":443" // Default listening port for HTTPS.
)

type ServerConfig struct {
	Name           string         `json:"name"`
	Address        string         `json:"address"`
	HTTPSAddr      string         `json:"httpsAddr"`
	Listeners      []net.Listener `json:"listeners"`
	HTTPSCertPath  string         `json:"httpsCertPath"`
	HTTPSKeyPath   string         `json:"httpsKeyPath"`
	TLSConfig      *tls.Config    `json:"tlsConfig"`
	ReadTimeout    time.Duration  `json:"read-timeout"`
	WriteTimeout   time.Duration  `json:"write-timeout"`
	IdleTimeout    time.Duration  `json:"idle-timeout"`
	MaxHeaderBytes int            `json:"max-header-bytes"`
	KeepAlive      bool           `json:"keep-alive"`

	// ======================================================================================================
	// Logging.
	// ======================================================================================================

	Logger           *gblog.Logger `json:"logger"`           // Logger specifies the logger for server.
	LogPath          string        `json:"logPath"`          // LogPath specifies the directory for storing logging files.
	LogLevel         string        `json:"logLevel"`         // LogLevel specifies the logging level for logger.
	LogStdout        bool          `json:"logStdout"`        // LogStdout specifies whether printing logging content to stdout.
	ErrorStack       bool          `json:"errorStack"`       // ErrorStack specifies whether logging stack information when error.
	ErrorLogEnabled  bool          `json:"errorLogEnabled"`  // ErrorLogEnabled enables error logging content to files.
	ErrorLogPattern  string        `json:"errorLogPattern"`  // ErrorLogPattern specifies the error log file pattern like: error-{Ymd}.log
	AccessLogEnabled bool          `json:"accessLogEnabled"` // AccessLogEnabled enables access logging content to files.
	AccessLogPattern string        `json:"accessLogPattern"` // AccessLogPattern specifies the error log file pattern like: access-{Ymd}.log

	// DumpRouterMap specifies whether automatically dumps router map when server starts.
	DumpRouterMap bool `json:"dumpRouterMap"`

	// Graceful enables graceful reload feature for all servers of the process.
	Graceful bool `json:"graceful"`

	// GracefulTimeout set the maximum survival time (seconds) of the parent process.
	GracefulTimeout uint8 `json:"gracefulTimeout"`

	// GracefulShutdownTimeout set the maximum survival time (seconds) before stopping the server.
	GracefulShutdownTimeout uint8 `json:"gracefulShutdownTimeout"`
}

func NewConfig() ServerConfig {
	return ServerConfig{
		Name:                    DefaultServerName,
		Address:                 ":0",
		HTTPSAddr:               "",
		ReadTimeout:             60 * time.Second,
		WriteTimeout:            0, // No timeout.
		IdleTimeout:             60 * time.Second,
		MaxHeaderBytes:          10240, // 10KB
		KeepAlive:               true,
		Logger:                  gblog.New(),
		LogLevel:                "all",
		LogStdout:               true,
		ErrorStack:              true,
		ErrorLogEnabled:         true,
		ErrorLogPattern:         "error-{Ymd}.log",
		AccessLogEnabled:        false,
		AccessLogPattern:        "access-{Ymd}.log",
		DumpRouterMap:           true,
		Graceful:                false,
		GracefulTimeout:         2, // seconds
		GracefulShutdownTimeout: 5, // seconds
	}
}

func (s *Server) SetConfig(c ServerConfig) error {
	s.config = c
	// Automatically add ':' prefix for address if it is missed.
	if s.config.Address != "" && !gbstr.Contains(s.config.Address, ":") {
		s.config.Address = ":" + s.config.Address
	}

	// Logging.
	if s.config.LogPath != "" && s.config.LogPath != s.config.Logger.GetPath() {
		if err := s.config.Logger.SetPath(s.config.LogPath); err != nil {
			return err
		}
	}
	if err := s.config.Logger.SetLevelStr(s.config.LogLevel); err != nil {
		intlog.Errorf(context.TODO(), `%+v`, err)
	}
	intlog.Printf(context.TODO(), "SetConfig: %+v", s.config)
	return nil
}
