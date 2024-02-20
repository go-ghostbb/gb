package gbcmd_test

import (
	"context"
	"ghostbb.io/gb/frame/g"
	gbcmd "ghostbb.io/gb/os/gb_cmd"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbtag "ghostbb.io/gb/util/gb_tag"
	"os"
	"testing"
)

type commandBuild struct {
	g.Meta               `name:"build" root:"build" args:"true" brief:"{commandBuildBrief}" dc:"{commandBuildDc}" eg:"{commandBuildEg}" ad:"{commandBuildAd}"`
	nodeNameInConfigFile string // nodeNameInConfigFile is the node name for compiler configurations in configuration file.
	packedGoFileName     string // packedGoFileName specifies the file name for packing common folders into one single go file.
}

const (
	commandBuildBrief = `cross-building go project for lots of platforms`
	commandBuildEg    = `
gb build main.go
gb build main.go --pack public,template
gb build main.go --cgo
gb build main.go -m none 
gb build main.go -n my-app -a all -s all
gb build main.go -n my-app -a amd64,386 -s linux -p .
gb build main.go -n my-app -v 1.0 -a amd64,386 -s linux,windows,darwin -p ./docker/bin
`
	commandBuildDc = `
The "build" command is most commonly used command, which is designed as a powerful wrapper for 
"go build" command for convenience cross-compiling usage. 
It provides much more features for building binary:
1. Cross-Compiling for many platforms and architectures.
2. Configuration file support for compiling.
3. Build-In Variables.
`
	commandBuildAd = `
PLATFORMS
    darwin    amd64,arm64
    freebsd   386,amd64,arm
    linux     386,amd64,arm,arm64,ppc64,ppc64le,mips,mipsle,mips64,mips64le
    netbsd    386,amd64,arm
    openbsd   386,amd64,arm
    windows   386,amd64
`
	// https://golang.google.cn/doc/install/source
	commandBuildPlatforms = `
    darwin    amd64
    darwin    arm64
    ios       amd64
    ios       arm64
    freebsd   386
    freebsd   amd64
    freebsd   arm
    linux     386
    linux     amd64
    linux     arm
    linux     arm64
    linux     ppc64
    linux     ppc64le
    linux     mips
    linux     mipsle
    linux     mips64
    linux     mips64le
    netbsd    386
    netbsd    amd64
    netbsd    arm
    openbsd   386
    openbsd   amd64
    openbsd   arm
    windows   386
    windows   amd64
	android   arm
	dragonfly amd64
	plan9     386
	plan9     amd64
	solaris   amd64
`
	commandBuildBriefPack = `
destination file path for packed file. if extension of the filename is ".go" and "-n" option is given, 
it enables packing SRC to go file, or else it packs SRC into a binary file.

`
	commandGenDaoBriefJsonCase = `
generated json tag case for model struct, cases are as follows:
| Case            | Example            |
|---------------- |--------------------|
| Camel           | AnyKindOfString    | 
| CamelLower      | anyKindOfString    | default
| Snake           | any_kind_of_string |
| SnakeScreaming  | ANY_KIND_OF_STRING |
| SnakeFirstUpper | rgb_code_md5       |
| Kebab           | any-kind-of-string |
| KebabScreaming  | ANY-KIND-OF-STRING |
`
)

func init() {
	gbtag.Sets(map[string]string{
		`commandBuildBrief`:          commandBuildBrief,
		`commandBuildDc`:             commandBuildDc,
		`commandBuildEg`:             commandBuildEg,
		`commandBuildAd`:             commandBuildAd,
		`commandBuildBriefPack`:      commandBuildBriefPack,
		`commandGenDaoBriefJsonCase`: commandGenDaoBriefJsonCase,
	})
}

type commandBuildInput struct {
	g.Meta   `name:"build" config:"gbcli.build"`
	Name     string `short:"n" name:"name"     brief:"output binary name"`
	Version  string `short:"v" name:"version"  brief:"output binary version"`
	Arch     string `short:"a" name:"arch"     brief:"output binary architecture, multiple arch separated with ','"`
	System   string `short:"s" name:"system"   brief:"output binary system, multiple os separated with ','"`
	Output   string `short:"o" name:"output"   brief:"output binary path, used when building single binary file"`
	Path     string `short:"p" name:"path"     brief:"output binary directory path, default is './bin'" d:"./bin"`
	Extra    string `short:"e" name:"extra"    brief:"extra custom \"go build\" options"`
	Mod      string `short:"m" name:"mod"      brief:"like \"-mod\" option of \"go build\", use \"-m none\" to disable go module"`
	Cgo      bool   `short:"c" name:"cgo"      brief:"enable or disable cgo feature, it's disabled in default" orphan:"true"`
	JsonCase string `short:"j" name:"jsonCase" brief:"{commandGenDaoBriefJsonCase}" d:"CamelLower"`
	Pack     string `name:"pack" brief:"{commandBuildBriefPack}"`
}

type commandBuildOutput struct{}

func (c commandBuild) Index(ctx context.Context, in commandBuildInput) (out *commandBuildOutput, err error) {
	return
}

func TestNewFromObject(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			ctx = gbctx.New()
		)
		cmd, err := gbcmd.NewFromObject(commandBuild{
			nodeNameInConfigFile: "gbcli.build",
			packedGoFileName:     "build_pack_data.go",
		})
		t.AssertNil(err)

		os.Args = []string{"build", "-h"}
		err = cmd.RunWithError(ctx)
		t.AssertNil(err)
	})
}
