package main

import (
	"ghostbb.io/gb/cmd/gb/gbcmd"
	"ghostbb.io/gb/cmd/gb/internal/utility/mlog"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbctx "ghostbb.io/gb/os/gb_ctx"
)

func main() {
	var (
		ctx = gbctx.GetInitCtx()
	)
	command, err := gbcmd.GetCommand(ctx)
	if err != nil {
		mlog.Fatalf(`%+v`, err)
	}
	if command == nil {
		panic(gberror.New(`retrieve root command failed for "gb"`))
	}
	command.Run(ctx)
}
