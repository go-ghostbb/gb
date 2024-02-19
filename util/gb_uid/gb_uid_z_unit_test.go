package gbuid_test

import (
	gbset "ghostbb.io/gb/container/gb_set"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbuid "ghostbb.io/gb/util/gb_uid"
	"testing"
)

func Test_S(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		set := gbset.NewStrSet()
		for i := 0; i < 1000000; i++ {
			s := gbuid.S()
			t.Assert(set.AddIfNotExist(s), true)
			t.Assert(len(s), 32)
		}
	})
}

func Test_S_Data(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(len(gbuid.S([]byte("123"))), 32)
	})
}
