package gbhttp

import (
	"context"
	gblog "ghostbb.io/gb/os/gb_log"
	gbproc "ghostbb.io/gb/os/gb_proc"
	"os"
)

func handleProcessSignal() {
	var ctx = context.TODO()
	gbproc.AddSigHandlerShutdown(func(sig os.Signal) {
		showdownServers(ctx, sig)
	})

	gbproc.Listen()
}

func showdownServers(ctx context.Context, signal os.Signal) {
	if signal != nil {
		gblog.Printf(
			ctx,
			"%d: server gracefully shutting down by signal: %s",
			gbproc.Pid(), signal.String(),
		)
	}
	serverMapping.RLockFunc(func(m map[string]interface{}) {
		for _, v := range m {
			server := v.(*Server)
			server.Showdown(ctx)
		}
	})
}
