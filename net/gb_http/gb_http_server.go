package gbhttp

import (
	"bytes"
	"context"
	gbtype "ghostbb.io/gb/container/gb_type"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/intlog"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbfile "ghostbb.io/gb/os/gb_file"
	gblog "ghostbb.io/gb/os/gb_log"
	gbproc "ghostbb.io/gb/os/gb_proc"
	gbtimer "ghostbb.io/gb/os/gb_timer"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbutil "ghostbb.io/gb/util/gb_util"
	"github.com/gin-gonic/gin"
	"github.com/olekukonko/tablewriter"
	"net/http"
	"strings"
	"sync"
	"time"
)

func GetServer(name ...interface{}) *Server {
	serverName := DefaultServerName
	if len(name) > 0 && name[0] != "" {
		serverName = gbconv.String(name[0])
	}

	v := serverMapping.GetOrSetFuncLock(serverName, func() interface{} {
		e := gin.New()
		s := &Server{
			Engine:      e,
			instance:    serverName,
			servers:     make([]*internalServer, 0),
			closeChan:   make(chan struct{}, 10000),
			serverCount: gbtype.NewInt(),
		}
		// Initialize the server using default configurations.
		if err := s.SetConfig(NewConfig()); err != nil {
			panic(gberror.WrapCode(gbcode.CodeInvalidConfiguration, err, ""))
		}

		e.Use(func(c *gin.Context) {
			c.Set(ServerContextKey, gbctx.New())
		})

		e.Use(s.loggerMiddleware(), s.terminal(), s.Recovery())
		return s
	})

	return v.(*Server)
}

// GetName returns the name of the server.
func (s *Server) GetName() string {
	return s.config.Name
}

// SetName sets the name for the server.
func (s *Server) SetName(name string) {
	s.config.Name = name
}

// SetConfigWithMap sets the configuration for the server using map.
func (s *Server) SetConfigWithMap(m map[string]interface{}) error {
	// The m now is a shallow copy of m.
	// Any changes to m does not affect the original one.
	// A little tricky, isn't it?
	m = gbutil.MapCopy(m)
	// Allow setting the size configuration items using string size like:
	// 1m, 100mb, 512kb, etc.
	if k, v := gbutil.MapPossibleItemByKey(m, "MaxHeaderBytes"); k != "" {
		m[k] = gbfile.StrToSize(gbconv.String(v))
	}
	// Update the current configuration object.
	// It only updates the configured keys not all the object.
	if err := gbconv.Struct(m, &s.config); err != nil {
		return err
	}
	return s.SetConfig(s.config)
}

// Start starts listening on configured port.
// This function does not block the process, you can use function Wait blocking the process.
func (s *Server) Start() error {
	var ctx = gbctx.GetInitCtx()

	// Server can only be run once.
	if s.Status() == ServerStatusRunning {
		return gberror.NewCode(gbcode.CodeInvalidOperation, "server is already running")
	}

	// Logging path setting check.
	if s.config.LogPath != "" && s.config.LogPath != s.config.Logger.GetPath() {
		if err := s.config.Logger.SetPath(s.config.LogPath); err != nil {
			return err
		}
	}

	// ================================================================================================
	// Start the HTTP server.
	// ===============================================================================================
	s.startServer()

	// If this is a child process, it then notifies its parent exit.
	if gbproc.IsChild() {
		gbtimer.SetTimeout(ctx, time.Duration(s.config.GracefulTimeout)*time.Second, func(ctx context.Context) {
			if err := gbproc.Send(gbproc.PPid(), []byte("exit"), adminGProcCommGroup); err != nil {
				intlog.Errorf(ctx, `server error in process communication: %+v`, err)
			}
		})
	}

	s.doRouterMapDump()

	return nil
}

func (s *Server) startServer() {
	var (
		ctx          = context.TODO()
		httpsEnabled bool
	)
	// HTTPS
	if s.config.HTTPSCertPath != "" && s.config.HTTPSKeyPath != "" {
		if len(s.config.HTTPSAddr) == 0 {
			if len(s.config.Address) > 0 {
				s.config.HTTPSAddr = s.config.Address
				s.config.Address = ""
			} else {
				s.config.HTTPSAddr = defaultHttpsAddr
			}
			httpsEnabled = len(s.config.HTTPSAddr) > 0

			for _, v := range strings.Split(s.config.HTTPSAddr, ",") {
				if len(v) == 0 {
					continue
				}

				s.servers = append(s.servers, s.newInternalServer(v))
				s.servers[len(s.servers)-1].isHttps = true
			}
		}
	}
	// HTTP
	if !httpsEnabled && len(s.config.Address) == 0 {
		s.config.Address = defaultHttpAddr
	}

	for _, v := range gbstr.SplitAndTrim(s.config.Address, ",") {
		if len(v) == 0 {
			continue
		}
		s.servers = append(s.servers, s.newInternalServer(v))
	}

	// Start listening asynchronously.
	serverRunning.Add(1)
	var wg = sync.WaitGroup{}
	for _, v := range s.servers {
		wg.Add(1)
		go func(server *internalServer) {
			s.serverCount.Add(1)
			var err error
			// Create listener.
			if server.isHttps {
				err = server.CreateListenerTLS(
					s.config.HTTPSCertPath, s.config.HTTPSKeyPath, s.config.TLSConfig,
				)
			} else {
				err = server.CreateListener()
			}
			if err != nil {
				s.Logger().Fatalf(ctx, `%+v`, err)
			}
			wg.Done()
			// Start listening and serving in blocking way.
			err = server.Serve(ctx)
			// The process exits if the server is closed with none closing error.
			if err != nil && !strings.EqualFold(http.ErrServerClosed.Error(), err.Error()) {
				s.Logger().Fatalf(ctx, `%+v`, err)
			}
			// If all the underlying servers' shutdown, the process exits.
			if s.serverCount.Add(-1) < 1 {
				s.closeChan <- struct{}{}
				if serverRunning.Add(-1) < 1 {
					serverMapping.Remove(s.instance)
					allShutdownChan <- struct{}{}
				}
			}
		}(v)
	}
	wg.Wait()
}

// Wait blocks to wait for all servers done.
// It's commonly used in multiple server situation.
func Wait() {
	var ctx = context.TODO()

	// Signal handler in asynchronous way.
	go handleProcessSignal()

	<-allShutdownChan

	gblog.Infof(ctx, "pid[%d]: all servers shutdown", gbproc.Pid())
}

// Status retrieves and returns the server status.
func (s *Server) Status() ServerStatus {
	if serverRunning.Val() == 0 {
		return ServerStatusStopped
	}
	// If any underlying server is running, the server status is running.
	for _, v := range s.servers {
		if v.status.Val() == ServerStatusRunning {
			return ServerStatusRunning
		}
	}
	return ServerStatusStopped
}

func (s *Server) GetRoutes() gin.RoutesInfo {
	return s.Engine.Routes()
}

// doRouterMapDump checks and dumps the router map to the log.
func (s *Server) doRouterMapDump() {
	if !s.config.DumpRouterMap {
		return
	}

	var (
		headers = []string{"GROUP", "ADDRESS", "METHOD", "ROUTE", "HANDLER"}
		routes  = s.GetRoutes()
	)
	if len(routes) > 0 {
		buffer := bytes.NewBuffer(nil)
		table := tablewriter.NewWriter(buffer)
		table.SetHeader(headers)
		table.SetRowLine(true)
		table.SetBorder(false)
		table.SetCenterSeparator("|")

		for _, route := range routes {
			table.Append([]string{
				s.instance, s.config.Address, route.Method, route.Path, route.Handler,
			})
		}

		table.Render()
		s.Logger().Header(false).Stdout(true).Printf(context.TODO(), "\n%s", buffer.String())
	}
}

// Run starts server listening in blocking way.
// It's commonly used for single server situation.
func (s *Server) Run() {
	var ctx = context.TODO()

	if err := s.Start(); err != nil {
		s.Logger().Fatalf(ctx, `%+v`, err)
	}

	// Signal handler in asynchronous way.
	go handleProcessSignal()

	// Blocking using channel.
	<-s.closeChan

	s.Logger().Infof(ctx, "pid[%d]: all servers shutdown", gbproc.Pid())
}

func (s *Server) Bind(obj ...IBind) {
	for _, o := range obj {
		o.Init(s.Group(""))
	}
}
