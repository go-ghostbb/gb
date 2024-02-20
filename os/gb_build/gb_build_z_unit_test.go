package gbbuild_test

import (
	"ghostbb.io/gb/frame/g"
	gbbuild "ghostbb.io/gb/os/gb_build"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
)

func Test_Info(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbconv.Map(gbbuild.Info()), g.Map{
			"Gb":      "",
			"Golang":  "",
			"Git":     "",
			"Time":    "",
			"Version": "",
			"Data":    g.Map{},
		})
	})
}

func Test_Get(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbbuild.Get(`none`), nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbbuild.Get(`none`, 1), 1)
	})
}

func Test_Map(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbbuild.Data(), map[string]interface{}{})
	})
}
