package gbcfg_test

import (
	"fmt"
	"ghostbb.io/gb/frame/g"
	gbcfg "ghostbb.io/gb/os/gb_cfg"
	gbcmd "ghostbb.io/gb/os/gb_cmd"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbenv "ghostbb.io/gb/os/gb_env"
	"os"
)

func ExampleConfig_GetWithEnv() {
	var (
		key = `ENV_TEST`
		ctx = gbctx.New()
	)
	v, err := g.Cfg().GetWithEnv(ctx, key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("env:%s\n", v)
	if err = gbenv.Set(key, "gb"); err != nil {
		panic(err)
	}
	v, err = g.Cfg().GetWithEnv(ctx, key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("env:%s", v)

	// Output:
	// env:
	// env:gb
}

func ExampleConfig_GetWithCmd() {
	var (
		key = `cmd.test`
		ctx = gbctx.New()
	)
	v, err := g.Cfg().GetWithCmd(ctx, key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("cmd:%s\n", v)
	// Re-Initialize custom command arguments.
	os.Args = append(os.Args, fmt.Sprintf(`--%s=yes`, key))
	gbcmd.Init(os.Args...)
	// Retrieve the configuration and command option again.
	v, err = g.Cfg().GetWithCmd(ctx, key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("cmd:%s", v)

	// Output:
	// cmd:
	// cmd:yes
}

func Example_NewWithAdapter() {
	var (
		ctx          = gbctx.New()
		content      = `{"a":"b", "c":1}`
		adapter, err = gbcfg.NewAdapterContent(content)
	)
	if err != nil {
		panic(err)
	}
	config := gbcfg.NewWithAdapter(adapter)
	fmt.Println(config.MustGet(ctx, "a"))
	fmt.Println(config.MustGet(ctx, "c"))

	// Output:
	// b
	// 1
}
