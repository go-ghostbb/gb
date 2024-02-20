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

type TestNoNameTagCase struct {
	g.Meta `name:"root"`
}

type TestNoNameTagCaseRootInput struct {
	Name string
}

type TestNoNameTagCaseRootOutput struct {
	Content string
}

func (c *TestNoNameTagCase) TEST(ctx context.Context, in TestNoNameTagCaseRootInput) (out *TestNoNameTagCaseRootOutput, err error) {
	out = &TestNoNameTagCaseRootOutput{
		Content: in.Name,
	}
	return
}

func Test_Command_NoNameTagCase(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var ctx = gbctx.New()
		cmd, err := gbcmd.NewFromObject(TestNoNameTagCase{})
		t.AssertNil(err)

		os.Args = []string{"root", "TEST", "-name=john"}
		value, err := cmd.RunWithValueError(ctx)
		t.AssertNil(err)
		t.Assert(value, `{"Content":"john"}`)
	})
}
