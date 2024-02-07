package gbhttp

import (
	"context"
	gbproc "ghostbb.io/gb/os/gb_proc"
	"os"
)

// handleProcessSignal handles all signals from system in blocking way.
func handleProcessSignal() {
	var ctx = context.TODO()
	gbproc.AddSigHandlerShutdown(func(sig os.Signal) {
		shutdownWebServersGracefully(ctx, sig)
	})

	gbproc.Listen()
}
