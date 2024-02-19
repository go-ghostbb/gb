package gbvalid_test

import (
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func Test_CI(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		err := g.Validator().Data("id").Rules("in:Id,Name").Run(ctx)
		t.AssertNE(err, nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		err := g.Validator().Data("id").Rules("ci|in:Id,Name").Run(ctx)
		t.AssertNil(err)
	})
	gbtest.C(t, func(t *gbtest.T) {
		err := g.Validator().Ci().Rules("in:Id,Name").Data("id").Run(ctx)
		t.AssertNil(err)
	})
}
