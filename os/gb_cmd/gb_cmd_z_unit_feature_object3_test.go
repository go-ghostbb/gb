package gbcmd_test

import (
	"context"
	"ghostbb.io/gb/frame/g"
	gbcmd "ghostbb.io/gb/os/gb_cmd"
	gbctx "ghostbb.io/gb/os/gb_ctx"
	gbtest "ghostbb.io/gb/test/gb_test"
	"os"
	"testing"
)

type TestParamsCase struct {
	g.Meta `name:"root" root:"root"`
}

type TestParamsCaseRootInput struct {
	g.Meta `name:"root"`
	Name   string
}

type TestParamsCaseRootOutput struct {
	Content string
}

func (c *TestParamsCase) Root(ctx context.Context, in TestParamsCaseRootInput) (out *TestParamsCaseRootOutput, err error) {
	out = &TestParamsCaseRootOutput{
		Content: in.Name,
	}
	return
}

func Test_Command_ParamsCase(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var ctx = gbctx.New()
		cmd, err := gbcmd.NewFromObject(TestParamsCase{})
		t.AssertNil(err)

		os.Args = []string{"root", "-name=john"}
		value, err := cmd.RunWithValueError(ctx)
		t.AssertNil(err)
		t.Assert(value, `{"Content":"john"}`)
	})
}
