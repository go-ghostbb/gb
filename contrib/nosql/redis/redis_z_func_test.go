package redis

import (
	gbredis "ghostbb.io/gb/database/gb_redis"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
)

func Test_mustMergeOptionToArgs(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var args []interface{}
		newArgs := mustMergeOptionToArgs(args, gbredis.SetOption{
			NX:  true,
			Get: true,
		})
		t.Assert(newArgs, []interface{}{"NX", "Get"})
	})
	gbtest.C(t, func(t *gbtest.T) {
		var args []interface{}
		newArgs := mustMergeOptionToArgs(args, gbredis.SetOption{
			NX:  true,
			Get: true,
			TTLOption: gbredis.TTLOption{
				EX: gbconv.PtrInt64(60),
			},
		})
		t.Assert(newArgs, []interface{}{"EX", 60, "NX", "Get"})
	})
}
