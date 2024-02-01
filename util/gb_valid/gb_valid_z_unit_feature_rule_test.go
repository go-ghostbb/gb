package gbvalid_test

import (
	"ghostbb.io/frame/g"
	gbctx "ghostbb.io/os/gb_ctx"
	gbtest "ghostbb.io/test/gb_test"
	"testing"
)

var (
	ctx = gbctx.New()
)

func Test_Check(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		rule := "abc:6,16"
		val1 := 0
		val2 := 7
		val3 := 20
		err1 := g.Validator().Data(val1).Rules(rule).Run(ctx)
		err2 := g.Validator().Data(val2).Rules(rule).Run(ctx)
		err3 := g.Validator().Data(val3).Rules(rule).Run(ctx)
		t.Assert(err1, "InvalidRules: abc:6,16")
		t.Assert(err2, "InvalidRules: abc:6,16")
		t.Assert(err3, "InvalidRules: abc:6,16")
	})
}
