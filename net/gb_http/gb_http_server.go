package gbhttp

import (
	"bytes"
	"context"
	gberror "ghostbb.io/errors/gb_error"
	gbctx "ghostbb.io/os/gb_ctx"
	gblog "ghostbb.io/os/gb_log"
	gbproc "ghostbb.io/os/gb_proc"
	gbstr "ghostbb.io/text/gb_str"
	gbconv "ghostbb.io/util/gb_conv"
	gbutil "ghostbb.io/util/gb_util"
	"github.com/gin-gonic/gin"
	"github.com/olekukonko/tablewriter"
	"net/http"
	"time"
)

func NewServer(engin *Engine) *Server {
	return &Server{
		server: &http.Server{
			Handler: engin,
		},
		Engine:    engin,
		closeChan: make(chan struct{}, 10000),
	}
}

func GetServer(name ...string) *Server {
	serverName := DefaultServerName
	if len(name) > 0 && name[0] != "" {
		serverName = gbconv.String(name[0])
	}

	v := serverMapping.GetOrSetFuncLock(serverName, func() interface{} {
		engine := Default()
		engine.instance = serverName
		server := NewServer(engine)

		server.Use(gin.Recovery(), engine.GBCtxMiddleware(), server.loggerMiddleware(), gin.LoggerWithFormatter(server.debugLog))
		return server
	})
	return v.(*Server)
}

func (s *Server) SetConfigWithMap(m map[string]interface{}) error {
	m = gbutil.MapCopy(m)

	if err := gbconv.Struct(m, &s.config); err != nil {
		return err
	}

	if _, kValue := gbutil.MapPossibleItemByKey(m, "KeepAlive"); kValue == nil {
		s.server.SetKeepAlivesEnabled(true)
	}
	if _, sValue := gbutil.MapPossibleItemByKey(m, "LogStdout"); sValue == nil {
		s.SetStdout(true)
	}

	return s.SetConfig(s.config)
}

func (s *Server) SetConfig(c ServerConfig) error {
	s.config = c
	if s.config.Address != "" && !gbstr.Contains(s.config.Address, ":") {
		s.config.Address = ":" + s.config.Address
	}
	s.server.Addr = s.config.Address
	s.server.ReadTimeout = s.config.ReadTimeout
	s.server.WriteTimeout = s.config.WriteTimeout
	s.server.IdleTimeout = s.config.IdleTimeout
	s.server.MaxHeaderBytes = s.config.MaxHeaderBytes
	s.server.SetKeepAlivesEnabled(s.config.KeepAlive)
	return nil
}

func (s *Server) getProto() string {
	proto := "http"
	if s.config.Https {
		proto = "https"
	}
	return proto
}

func (s *Server) Logger() *gblog.Logger {
	return s.logger
}

func (s *Server) Config() *ServerConfig {
	return &s.config
}

func (s *Server) SetStdout(stdout bool) {
	s.config.LogStdout = false
}

func (s *Server) Run() {
	go s.startServer()
	go handleProcessSignal()

	<-s.closeChan
	s.Logger().Stdout(true).Infof(context.TODO(), "pid[%d]: server shutdown", gbproc.Pid())
}

func (s *Server) startServer() {
	var err error

	s.Logger().Stdout(true).Infof(context.TODO(),
		"pid[%d]: %s server started listening on [%s]",
		gbproc.Pid(), s.getProto(), s.Config().Address,
	)
	s.doRouterMapDump()

	if s.Config().Https {
		err = s.server.ListenAndServeTLS(s.Config().CertFile, s.Config().KeyFile)
	} else {
		err = s.server.ListenAndServe()
	}
	if err != nil && !gberror.Is(err, http.ErrServerClosed) {
		s.Logger().Error(gbctx.New(), err)
	}
}

func (s *Server) GetEngine() *gin.Engine {
	return s.Engine.Engine
}

func (s *Server) doRouterMapDump() {
	var (
		headers = []string{"GROUP", "ADDRESS", "METHOD", "ROUTE", "HANDLER"}
		routes  = s.GetEngine().Routes()
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
				s.instance, s.Config().Address, route.Method, route.Path, route.Handler,
			})
		}

		table.Render()
		s.Logger().Header(false).Stdout(true).Printf(context.TODO(), "\n%s", buffer.String())
	}
}

func (s *Server) Bind(groups ...IBind) {
	for _, group := range groups {
		group.Register(s.Group("/"))
	}
}

func (s *Server) SetPort(p string) {
	if p == "" {
		return
	}
	if p != "" && !gbstr.Contains(p, ":") {
		p = ":" + p
	}
	s.server.Addr = p
}

func (s *Server) Showdown(ctx context.Context) {
	timeoutCtx, cancelFunc := context.WithTimeout(
		ctx,
		time.Duration(5*time.Second),
	)
	defer cancelFunc()
	if err := s.server.Shutdown(timeoutCtx); err != nil {
		s.Logger().Errorf(
			ctx,
			"%d: %s server [%s] shutdown error: %v",
			gbproc.Pid(), s.getProto(), s.config.Address, err,
		)
	}
	s.closeChan <- struct{}{}
}
