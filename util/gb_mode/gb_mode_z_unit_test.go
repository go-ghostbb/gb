package gbmode_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbmode "ghostbb.io/gb/util/gb_mode"
	"testing"
)

func Test_AutoCheckSourceCodes(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbmode.IsDevelop(), true)
	})
}

func Test_Set(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		oldMode := gbmode.Mode()
		defer gbmode.Set(oldMode)
		gbmode.SetDevelop()
		t.Assert(gbmode.IsDevelop(), true)
		t.Assert(gbmode.IsTesting(), false)
		t.Assert(gbmode.IsStaging(), false)
		t.Assert(gbmode.IsProduct(), false)
	})
	gbtest.C(t, func(t *gbtest.T) {
		oldMode := gbmode.Mode()
		defer gbmode.Set(oldMode)
		gbmode.SetTesting()
		t.Assert(gbmode.IsDevelop(), false)
		t.Assert(gbmode.IsTesting(), true)
		t.Assert(gbmode.IsStaging(), false)
		t.Assert(gbmode.IsProduct(), false)
	})
	gbtest.C(t, func(t *gbtest.T) {
		oldMode := gbmode.Mode()
		defer gbmode.Set(oldMode)
		gbmode.SetStaging()
		t.Assert(gbmode.IsDevelop(), false)
		t.Assert(gbmode.IsTesting(), false)
		t.Assert(gbmode.IsStaging(), true)
		t.Assert(gbmode.IsProduct(), false)
	})
	gbtest.C(t, func(t *gbtest.T) {
		oldMode := gbmode.Mode()
		defer gbmode.Set(oldMode)
		gbmode.SetProduct()
		t.Assert(gbmode.IsDevelop(), false)
		t.Assert(gbmode.IsTesting(), false)
		t.Assert(gbmode.IsStaging(), false)
		t.Assert(gbmode.IsProduct(), true)
	})
}
