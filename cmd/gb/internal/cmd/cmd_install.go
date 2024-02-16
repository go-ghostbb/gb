package cmd

import (
	"context"
	"ghostbb.io/gb/cmd/gb/internal/service"
	"ghostbb.io/gb/frame/g"
)

var (
	Install = cInstall{}
)

type cInstall struct {
	g.Meta `name:"install" brief:"install gb binary to system (might need root/admin permission)"`
}

type cInstallInput struct {
	g.Meta `name:"install"`
}

type cInstallOutput struct{}

func (c cInstall) Index(ctx context.Context, in cInstallInput) (out *cInstallOutput, err error) {
	err = service.Install.Run(ctx)
	return
}
