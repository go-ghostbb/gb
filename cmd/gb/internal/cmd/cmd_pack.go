package cmd

import (
	"context"
	"ghostbb.io/gb/cmd/gb/internal/utility/allyes"
	"ghostbb.io/gb/cmd/gb/internal/utility/mlog"
	"ghostbb.io/gb/frame/g"
	gbcmd "ghostbb.io/gb/os/gb_cmd"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbres "ghostbb.io/gb/os/gb_res"
	gbtag "ghostbb.io/gb/util/gb_tag"
	"strings"
)

var (
	Pack = cPack{}
)

type cPack struct {
	g.Meta `name:"pack" usage:"{cPackUsage}" brief:"{cPackBrief}" eg:"{cPackEg}"`
}

const (
	cPackUsage = `gb pack SRC DST`
	cPackBrief = `packing any file/directory to a resource file, or a go file`
	cPackEg    = `
gb pack public data.bin
gb pack public,template data.bin
gb pack public,template packed/data.go
gb pack public,template,config packed/data.go
gb pack public,template,config packed/data.go -n=packed -p=/var/www/my-app
gb pack /var/www/public packed/data.go -n=packed
`
	cPackSrcBrief = `source path for packing, which can be multiple source paths.`
	cPackDstBrief = `
destination file path for packed file. if extension of the filename is ".go" and "-n" option is given, 
it enables packing SRC to go file, or else it packs SRC into a binary file.
`
	cPackNameBrief     = `package name for output go file, it's set as its directory name if no name passed`
	cPackPrefixBrief   = `prefix for each file packed into the resource file`
	cPackKeepPathBrief = `keep the source path from system to resource file, usually for relative path`
)

func init() {
	gbtag.Sets(g.MapStrStr{
		`cPackUsage`:         cPackUsage,
		`cPackBrief`:         cPackBrief,
		`cPackEg`:            cPackEg,
		`cPackSrcBrief`:      cPackSrcBrief,
		`cPackDstBrief`:      cPackDstBrief,
		`cPackNameBrief`:     cPackNameBrief,
		`cPackPrefixBrief`:   cPackPrefixBrief,
		`cPackKeepPathBrief`: cPackKeepPathBrief,
	})
}

type cPackInput struct {
	g.Meta   `name:"pack"`
	Src      string `name:"SRC" arg:"true" v:"required" brief:"{cPackSrcBrief}"`
	Dst      string `name:"DST" arg:"true" v:"required" brief:"{cPackDstBrief}"`
	Name     string `name:"name"     short:"n" brief:"{cPackNameBrief}"`
	Prefix   string `name:"prefix"   short:"p" brief:"{cPackPrefixBrief}"`
	KeepPath bool   `name:"keepPath" short:"k" brief:"{cPackKeepPathBrief}" orphan:"true"`
}

type cPackOutput struct{}

func (c cPack) Index(ctx context.Context, in cPackInput) (out *cPackOutput, err error) {
	if gbfile.Exists(in.Dst) && gbfile.IsDir(in.Dst) {
		mlog.Fatalf("DST path '%s' cannot be a directory", in.Dst)
	}
	if !gbfile.IsEmpty(in.Dst) && !allyes.Check() {
		s := gbcmd.Scanf("path '%s' is not empty, files might be overwrote, continue? [y/n]: ", in.Dst)
		if strings.EqualFold(s, "n") {
			return
		}
	}
	if in.Name == "" && gbfile.ExtName(in.Dst) == "go" {
		in.Name = gbfile.Basename(gbfile.Dir(in.Dst))
	}
	var option = gbres.Option{
		Prefix:   in.Prefix,
		KeepPath: in.KeepPath,
	}
	if in.Name != "" {
		if err = gbres.PackToGoFileWithOption(in.Src, in.Dst, in.Name, option); err != nil {
			mlog.Fatalf("pack failed: %v", err)
		}
	} else {
		if err = gbres.PackToFileWithOption(in.Src, in.Dst, option); err != nil {
			mlog.Fatalf("pack failed: %v", err)
		}
	}
	mlog.Print("done!")
	return
}
