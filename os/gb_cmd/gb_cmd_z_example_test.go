package gbcmd_test

import (
	"context"
	"fmt"
	"ghostbb.io/gb/frame/g"
	gbcmd "ghostbb.io/gb/os/gb_cmd"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbenv "ghostbb.io/gb/os/gb_env"
	"os"
)

func ExampleInit() {
	gbcmd.Init("gb", "build", "main.go", "-o=gb.exe", "-y")
	fmt.Printf(`%#v`, gbcmd.GetArgAll())

	// Output:
	// []string{"gb", "build", "main.go"}
}

func ExampleGetArg() {
	gbcmd.Init("gb", "build", "main.go", "-o=gb.exe", "-y")
	fmt.Printf(
		`Arg[0]: "%v", Arg[1]: "%v", Arg[2]: "%v", Arg[3]: "%v"`,
		gbcmd.GetArg(0), gbcmd.GetArg(1), gbcmd.GetArg(2), gbcmd.GetArg(3),
	)

	// Output:
	// Arg[0]: "gb", Arg[1]: "build", Arg[2]: "main.go", Arg[3]: ""
}

func ExampleGetArgAll() {
	gbcmd.Init("gb", "build", "main.go", "-o=gb.exe", "-y")
	fmt.Printf(`%#v`, gbcmd.GetArgAll())

	// Output:
	// []string{"gb", "build", "main.go"}
}

func ExampleGetOpt() {
	gbcmd.Init("gb", "build", "main.go", "-o=gb.exe", "-y")
	fmt.Printf(
		`Opt["o"]: "%v", Opt["y"]: "%v", Opt["d"]: "%v"`,
		gbcmd.GetOpt("o"), gbcmd.GetOpt("y"), gbcmd.GetOpt("d", "default value"),
	)

	// Output:
	// Opt["o"]: "gb.exe", Opt["y"]: "", Opt["d"]: "default value"
}

func ExampleGetOpt_Def() {
	gbcmd.Init("gb", "build", "main.go", "-o=gb.exe", "-y")

	fmt.Println(gbcmd.GetOpt("s", "Def").String())

	// Output:
	// Def
}

func ExampleGetOptAll() {
	gbcmd.Init("gb", "build", "main.go", "-o=gb.exe", "-y")
	fmt.Printf(`%#v`, gbcmd.GetOptAll())

	// May Output:
	// map[string]string{"o":"gb.exe", "y":""}
}

func ExampleGetOptWithEnv() {
	fmt.Printf("Opt[gb.test]:%s\n", gbcmd.GetOptWithEnv("gb.test"))
	_ = gbenv.Set("GB_TEST", "YES")
	fmt.Printf("Opt[gb.test]:%s\n", gbcmd.GetOptWithEnv("gb.test"))

	// Output:
	// Opt[gb.test]:
	// Opt[gb.test]:YES
}

func ExampleParse() {
	os.Args = []string{"gb", "build", "main.go", "-o=gb.exe", "-y"}
	p, err := gbcmd.Parse(g.MapStrBool{
		"o,output": true,
		"y,yes":    false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(p.GetOpt("o"))
	fmt.Println(p.GetOpt("output"))
	fmt.Println(p.GetOpt("y") != nil)
	fmt.Println(p.GetOpt("yes") != nil)
	fmt.Println(p.GetOpt("none") != nil)
	fmt.Println(p.GetOpt("none", "Def"))

	// Output:
	// gb.exe
	// gb.exe
	// true
	// true
	// false
	// Def
}

func ExampleCommandFromCtx() {
	var (
		command = gbcmd.Command{
			Name: "start",
		}
	)

	ctx := context.WithValue(gbctx.New(), gbcmd.CtxKeyCommand, &command)
	unAddCtx := context.WithValue(gbctx.New(), gbcmd.CtxKeyCommand, &gbcmd.Command{})
	nonKeyCtx := context.WithValue(gbctx.New(), "Testkey", &gbcmd.Command{})

	fmt.Println(gbcmd.CommandFromCtx(ctx).Name)
	fmt.Println(gbcmd.CommandFromCtx(unAddCtx).Name)
	fmt.Println(gbcmd.CommandFromCtx(nonKeyCtx) == nil)

	// Output:
	// start
	//
	// true
}

func ExampleCommand_AddCommand() {
	commandRoot := &gbcmd.Command{
		Name: "gb",
	}
	commandRoot.AddCommand(&gbcmd.Command{
		Name: "start",
	}, &gbcmd.Command{})

	commandRoot.Print()

	// Output:
	// USAGE
	//     gb COMMAND [OPTION]
	//
	// COMMAND
	//     start
}

func ExampleCommand_AddCommand_Repeat() {
	commandRoot := &gbcmd.Command{
		Name: "gb",
	}
	err := commandRoot.AddCommand(&gbcmd.Command{
		Name: "start",
	}, &gbcmd.Command{
		Name: "stop",
	}, &gbcmd.Command{
		Name: "start",
	})

	fmt.Println(err)

	// Output:
	// command "start" is already added to command "gb"
}

func ExampleCommand_AddObject() {
	var (
		command = gbcmd.Command{
			Name: "start",
		}
	)

	command.AddObject(&TestCmdObject{})

	command.Print()

	// Output:
	// USAGE
	//     start COMMAND [OPTION]
	//
	// COMMAND
	//     root    root env command
}

func ExampleCommand_AddObject_Error() {
	var (
		command = gbcmd.Command{
			Name: "start",
		}
	)

	err := command.AddObject(&[]string{"Test"})

	fmt.Println(err)

	// Output:
	// input object should be type of struct, but got "*[]string"
}

func ExampleCommand_Print() {
	commandRoot := &gbcmd.Command{
		Name: "gb",
	}
	commandRoot.AddCommand(&gbcmd.Command{
		Name: "start",
	}, &gbcmd.Command{})

	commandRoot.Print()

	// Output:
	// USAGE
	//     gb COMMAND [OPTION]
	//
	// COMMAND
	//     start
}

func ExampleScan() {
	fmt.Println(gbcmd.Scan("gb scan"))

	// Output:
	// gb scan
}

func ExampleScanf() {
	fmt.Println(gbcmd.Scanf("gb %s", "scanf"))

	// Output:
	// gb scanf
}

func ExampleParserFromCtx() {
	parser, _ := gbcmd.Parse(nil)

	ctx := context.WithValue(gbctx.New(), gbcmd.CtxKeyParser, parser)
	nilCtx := context.WithValue(gbctx.New(), "NilCtxKeyParser", parser)

	fmt.Println(gbcmd.ParserFromCtx(ctx).GetArgAll())
	fmt.Println(gbcmd.ParserFromCtx(nilCtx) == nil)

	// Output:
	// [gb build main.go]
	// true
}

func ExampleParseArgs() {
	p, _ := gbcmd.ParseArgs([]string{
		"gb", "--force", "remove", "-fq", "-p=www", "path", "-n", "root",
	}, nil)

	fmt.Println(p.GetArgAll())
	fmt.Println(p.GetOptAll())

	// Output:
	// [gb path]
	// map[force:remove fq: n:root p:www]
}

func ExampleParser_GetArg() {
	p, _ := gbcmd.ParseArgs([]string{
		"gb", "--force", "remove", "-fq", "-p=www", "path", "-n", "root",
	}, nil)

	fmt.Println(p.GetArg(-1, "Def").String())
	fmt.Println(p.GetArg(-1) == nil)

	// Output:
	// Def
	// true
}
