package gbhttp

import (
	gbmap "ghostbb.io/gb/container/gb_map"
	gbtype "ghostbb.io/gb/container/gb_type"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

const (
	DefaultServerName                = "default"
	ServerContextKey                 = "gb.server.ctx"
	ServerStatusStopped ServerStatus = 0
	ServerStatusRunning ServerStatus = 1

	// FreePortAddress marks the server listens using random free port.
	FreePortAddress = ":0"
)

var (
	// serverMapping stores more than one server instances for current processes.
	// The key is the name of the server, and the value is its instance.
	serverMapping = gbmap.NewStrAnyMap(true)

	// serverRunning marks the running server counts.
	// If there is no successful server running or all servers' shutdown, this value is 0.
	serverRunning = gbtype.NewInt()

	// allShutdownChan is the event for all servers have done its serving and exit.
	// It is used for process blocking purpose.
	allShutdownChan = make(chan struct{}, 1000)
)

type (
	Server struct {
		*gin.Engine
		instance    string            // Instance name of current HTTP server.
		config      ServerConfig      // Server configuration.
		servers     []*internalServer // Underlying http.Server array.
		serverCount *gbtype.Int       // Underlying http.Server number for internal usage.
		closeChan   chan struct{}     // Used for underlying server closing event notification.
	}

	// ServerStatus is the server status enum type.
	ServerStatus = int

	// Listening file descriptor mapping.
	// The key is either "http" or "https" and the value is its FD.
	listenerFdMap = map[string]string

	IBind interface {
		Init(group *gin.RouterGroup)
	}
)
