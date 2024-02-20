package cmd

import (
	"context"
	"ghostbb.io/gb"
	"ghostbb.io/gb/cmd/gb/internal/service"
	"ghostbb.io/gb/cmd/gb/internal/utility/mlog"
	"ghostbb.io/gb/frame/g"
	gbcmd "ghostbb.io/gb/os/gb_cmd"
	gbtag "ghostbb.io/gb/util/gb_tag"
	"strings"
)

var (
	GB = cGB{}
)

type cGB struct {
	g.Meta `name:"gb" ad:"{cGBAd}"`
}

const (
	cGBAd = `
ADDITIONAL
    Use "gb COMMAND -h" for details about a command.
`
)

func init() {
	gbtag.Sets(g.MapStrStr{
		"cGBAd": cGBAd,
	})
}

type cGBInput struct {
	g.Meta  `name:"gb"`
	Yes     bool `short:"y" name:"yes"     brief:"all yes for all command without prompt ask"   orphan:"true"`
	Version bool `short:"v" name:"version" brief:"show version information of current binary"   orphan:"true"`
	Debug   bool `short:"d" name:"debug"   brief:"show internal detailed debugging information" orphan:"true"`
}

type cGBOutput struct{}

func (c cGB) Index(ctx context.Context, in cGBInput) (out *cGBOutput, err error) {
	// Version.
	if in.Version {
		_, err = Version.Index(ctx, cVersionInput{})
		return
	}

	answer := "n"
	// No argument or option, do installation checks.
	if data, isInstalled := service.Install.IsInstalled(); !isInstalled {
		mlog.Print("hi, it seams it's the first time you installing gb cli.")
		answer = gbcmd.Scanf("do you want to install gb(%s) binary to your system? [y/n]: ", gb.VERSION)
	} else if !data.IsSelf {
		mlog.Print("hi, you have installed gb cli.")
		answer = gbcmd.Scanf("do you want to install gb(%s) binary to your system? [y/n]: ", gb.VERSION)
	}
	if strings.EqualFold(answer, "y") {
		if err = service.Install.Run(ctx); err != nil {
			return
		}
		gbcmd.Scan("press `Enter` to exit...")
		return
	}

	// Print help content.
	gbcmd.CommandFromCtx(ctx).Print()
	return
}
