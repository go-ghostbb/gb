package gbhttp

import (
	"context"
	gbtype "ghostbb.io/gb/container/gb_type"
	gblog "ghostbb.io/gb/os/gb_log"
	gbproc "ghostbb.io/gb/os/gb_proc"
	"os"
)

const (
	// Allow executing management command after server starts after this interval in milliseconds.
	adminActionIntervalLimit = 2000
	adminActionNone          = 0
	adminActionRestarting    = 1
	adminActionShuttingDown  = 2
	adminActionReloadEnvKey  = "GB_SERVER_RELOAD"
	adminActionRestartEnvKey = "GB_SERVER_RESTART"
	adminGProcCommGroup      = "GB_GPROC_HTTP_SERVER"
)

var (
	// serverProcessStatus is the server status for operation of current process.
	serverProcessStatus = gbtype.NewInt()
)

// shutdownWebServersGracefully gracefully shuts down all servers.
func shutdownWebServersGracefully(ctx context.Context, signal os.Signal) {
	serverProcessStatus.Set(adminActionShuttingDown)
	if signal != nil {
		gblog.Printf(
			ctx,
			"%d: server gracefully shutting down by signal: %s",
			gbproc.Pid(), signal.String(),
		)
	} else {
		gblog.Printf(ctx, "%d: server gracefully shutting down by api", gbproc.Pid())
	}
	serverMapping.RLockFunc(func(m map[string]interface{}) {
		for _, v := range m {
			server := v.(*Server)
			for _, s := range server.servers {
				s.shutdown(ctx)
			}
		}
	})
}
