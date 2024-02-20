package gbcmd_test

import (
	"context"
	"fmt"
	gbcode "ghostbb.io/gb/errors/gb_code"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/frame/g"
	gbcmd "ghostbb.io/gb/os/gb_cmd"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbenv "ghostbb.io/gb/os/gb_env"
	gbtest "ghostbb.io/gb/test/gb_test"
	"os"
	"strings"
	"testing"
)

func Test_Default(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		gbcmd.Init([]string{"gb", "--force", "remove", "-fq", "-p=www", "path", "-n", "root"}...)
		t.Assert(len(gbcmd.GetArgAll()), 2)
		t.Assert(gbcmd.GetArg(1), "path")
		t.Assert(gbcmd.GetArg(100, "test"), "test")
		t.Assert(gbcmd.GetOpt("force"), "remove")
		t.Assert(gbcmd.GetOpt("n"), "root")
		t.Assert(gbcmd.GetOpt("fq").IsNil(), false)
		t.Assert(gbcmd.GetOpt("p").IsNil(), false)
		t.Assert(gbcmd.GetOpt("none").IsNil(), true)
		t.Assert(gbcmd.GetOpt("none", "value"), "value")
	})
	gbtest.C(t, func(t *gbtest.T) {
		gbcmd.Init([]string{"gb", "gen", "-h"}...)
		t.Assert(len(gbcmd.GetArgAll()), 2)
		t.Assert(gbcmd.GetOpt("h"), "")
		t.Assert(gbcmd.GetOpt("h").IsNil(), false)
	})
}

func Test_BuildOptions(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := gbcmd.BuildOptions(g.MapStrStr{
			"n": "john",
		})
		t.Assert(s, "-n=john")
	})

	gbtest.C(t, func(t *gbtest.T) {
		s := gbcmd.BuildOptions(g.MapStrStr{
			"n": "john",
		}, "-test")
		t.Assert(s, "-testn=john")
	})

	gbtest.C(t, func(t *gbtest.T) {
		s := gbcmd.BuildOptions(g.MapStrStr{
			"n1": "john",
			"n2": "huang",
		})
		t.Assert(strings.Contains(s, "-n1=john"), true)
		t.Assert(strings.Contains(s, "-n2=huang"), true)
	})
}

func Test_GetWithEnv(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		gbenv.Set("TEST", "1")
		defer gbenv.Remove("TEST")
		t.Assert(gbcmd.GetOptWithEnv("test"), 1)
	})
	gbtest.C(t, func(t *gbtest.T) {
		gbenv.Set("TEST", "1")
		defer gbenv.Remove("TEST")
		gbcmd.Init("-test", "2")
		t.Assert(gbcmd.GetOptWithEnv("test"), 2)
	})
}

func Test_Command(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			ctx = gbctx.New()
			err error
		)
		commandRoot := &gbcmd.Command{
			Name: "gb",
		}
		// env
		commandEnv := &gbcmd.Command{
			Name: "env",
			Func: func(ctx context.Context, parser *gbcmd.Parser) error {
				fmt.Println("env")
				return nil
			},
		}
		// test
		commandTest := &gbcmd.Command{
			Name:        "test",
			Brief:       "test brief",
			Description: "test description current Golang environment variables",
			Examples: `
gb get ghostbb.io/gb
gb get ghostbb.io/gb@latest
gb get ghostbb.io/gb@master
gb get golang.org/x/sys
`,
			Arguments: []gbcmd.Argument{
				{
					Name:   "my-option",
					Short:  "o",
					Brief:  "It's my custom option",
					Orphan: true,
				},
				{
					Name:   "another",
					Short:  "a",
					Brief:  "It's my another custom option",
					Orphan: true,
				},
			},
			Func: func(ctx context.Context, parser *gbcmd.Parser) error {
				fmt.Println("test")
				return nil
			},
		}
		err = commandRoot.AddCommand(
			commandEnv,
		)
		if err != nil {
			g.Log().Fatal(ctx, err)
		}
		err = commandRoot.AddObject(
			commandTest,
		)
		if err != nil {
			g.Log().Fatal(ctx, err)
		}

		if err = commandRoot.RunWithError(ctx); err != nil {
			if gberror.Code(err) == gbcode.CodeNotFound {
				commandRoot.Print()
			}
		}
	})
}

func Test_Command_Print(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			ctx = gbctx.New()
			err error
		)
		c := &gbcmd.Command{
			Name:        "gb",
			Description: `GoFrame Command Line Interface, which is your helpmate for building GoFrame application with convenience.`,
			Additional: `
Use 'gb help COMMAND' or 'gb COMMAND -h' for detail about a command, which has '...' in the tail of their comments.`,
		}
		// env
		commandEnv := &gbcmd.Command{
			Name:        "env",
			Brief:       "show current Golang environment variables, long brief.long brief.long brief.long brief.long brief.long brief.long brief.long brief.",
			Description: "show current Golang environment variables",
			Func: func(ctx context.Context, parser *gbcmd.Parser) error {
				return nil
			},
		}
		if err = c.AddCommand(commandEnv); err != nil {
			g.Log().Fatal(ctx, err)
		}
		// get
		commandGet := &gbcmd.Command{
			Name:        "get",
			Brief:       "install or update GF to system in default...",
			Description: "show current Golang environment variables",

			Examples: `
gb get ghostbb.io/gb
gb get ghostbb.io/gb@latest
gb get ghostbb.io/gb@master
gb get golang.org/x/sys
`,
			Func: func(ctx context.Context, parser *gbcmd.Parser) error {
				return nil
			},
		}
		if err = c.AddCommand(commandGet); err != nil {
			g.Log().Fatal(ctx, err)
		}
		// build
		//-n, --name       output binary name
		//-v, --version    output binary version
		//-a, --arch       output binary architecture, multiple arch separated with ','
		//-s, --system     output binary system, multiple os separated with ','
		//-o, --output     output binary path, used when building single binary file
		//-p, --path       output binary directory path, default is './bin'
		//-e, --extra      extra custom "go build" options
		//-m, --mod        like "-mod" option of "go build", use "-m none" to disable go module
		//-c, --cgo        enable or disable cgo feature, it's disabled in default

		commandBuild := gbcmd.Command{
			Name:  "build",
			Usage: "gb build FILE [OPTION]",
			Brief: "cross-building go project for lots of platforms...",
			Description: `
The "build" command is most commonly used command, which is designed as a powerful wrapper for
"go build" command for convenience cross-compiling usage.
It provides much more features for building binary:
1. Cross-Compiling for many platforms and architectures.
2. Configuration file support for compiling.
3. Build-In Variables.
`,
			Examples: `
gb build main.go
gb build main.go --swagger
gb build main.go --pack public,template
gb build main.go --cgo
gb build main.go -m none 
gb build main.go -n my-app -a all -s all
gb build main.go -n my-app -a amd64,386 -s linux -p .
gb build main.go -n my-app -v 1.0 -a amd64,386 -s linux,windows,darwin -p ./docker/bin
`,
			Func: func(ctx context.Context, parser *gbcmd.Parser) error {
				return nil
			},
		}
		if err = c.AddCommand(&commandBuild); err != nil {
			g.Log().Fatal(ctx, err)
		}
		_ = c.RunWithError(ctx)
	})
}

func Test_Command_NotFound(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		c0 := &gbcmd.Command{
			Name: "c0",
		}
		c1 := &gbcmd.Command{
			Name: "c1",
			FuncWithValue: func(ctx context.Context, parser *gbcmd.Parser) (interface{}, error) {
				return nil, nil
			},
		}
		c21 := &gbcmd.Command{
			Name: "c21",
			FuncWithValue: func(ctx context.Context, parser *gbcmd.Parser) (interface{}, error) {
				return nil, nil
			},
		}
		c22 := &gbcmd.Command{
			Name: "c22",
			FuncWithValue: func(ctx context.Context, parser *gbcmd.Parser) (interface{}, error) {
				return nil, nil
			},
		}
		t.AssertNil(c0.AddCommand(c1))
		t.AssertNil(c1.AddCommand(c21, c22))

		os.Args = []string{"c0", "c1", "c23", `--test="abc"`}
		err := c0.RunWithError(gbctx.New())
		t.Assert(err.Error(), `command "c1 c23" not found for command "c0", command line: c0 c1 c23 --test="abc"`)
	})
}
