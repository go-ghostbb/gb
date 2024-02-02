package gbhttp

import (
	"context"
	gbmap "ghostbb.io/gb/container/gb_map"
	gbtype "ghostbb.io/gb/container/gb_type"
	gblog "ghostbb.io/gb/os/gb_log"
	gbproc "ghostbb.io/gb/os/gb_proc"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type (
	Server struct {
		server *http.Server
		*Engine
		config    ServerConfig
		closeChan chan struct{}
	}

	Engine struct {
		*gin.Engine                   // origin engine
		instance     string           // Instance name of current HTTP server.
		groupMapping *gbmap.StrAnyMap // version group
		logger       *gblog.Logger
	}

	ServerConfig struct {
		Address        string        `json:"address"`
		Https          bool          `json:"https"`
		CertFile       string        `json:"cert-file"`
		KeyFile        string        `json:"key-file"`
		ReadTimeout    time.Duration `json:"read-timeout"`
		WriteTimeout   time.Duration `json:"write-timeout"`
		IdleTimeout    time.Duration `json:"idle-timeout"`
		MaxHeaderBytes int           `json:"max-header-bytes"`
		KeepAlive      bool          `json:"keep-alive"`

		LogCat    string `json:"log-cat"`
		LogStdout bool   `json:"log-stdout"`
	}

	IBind interface {
		Register(*gin.RouterGroup)
	}
)

const (
	DefaultServerName = "default"
	ServerContextKey  = "gb.server.ctx"
)

var (
	// serverMapping stores more than one server instances for current processes.
	// The key is the name of the server, and the value is its instance.
	serverMapping = gbmap.NewStrAnyMap(true)

	// serverRunning marks the running server counts.
	// If there is no successful server running or all servers' shutdown, this value is 0.
	serverRunning = gbtype.NewInt()
)

func RunMultiple(servers ...*Server) {
	chanList := make([]chan struct{}, 0)
	for _, server := range servers {
		chanList = append(chanList, server.closeChan)
		serverRunning.Add(1)
		go server.startServer()
		go handleProcessSignal()
	}
	for _, closeChan := range chanList {
		<-closeChan
		serverRunning.Add(-1)
	}

	gblog.Stdout(true).Infof(context.TODO(), "pid[%d]: servers shutdown", gbproc.Pid())
}
