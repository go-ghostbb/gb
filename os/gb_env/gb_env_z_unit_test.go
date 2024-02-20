package gbenv_test

import (
	"ghostbb.io/gb/frame/g"
	gbcmd "ghostbb.io/gb/os/gb_cmd"
	gbenv "ghostbb.io/gb/os/gb_env"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"os"
	"testing"
)

func Test_GBEnv_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(os.Environ(), gbenv.All())
	})
}

func Test_GBEnv_Map(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		value := gbconv.String(gbtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		t.AssertNil(err)
		t.Assert(gbenv.Map()[key], "TEST")
	})
}

func Test_GBEnv_Get(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		value := gbconv.String(gbtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		t.AssertNil(err)
		t.AssertEQ(gbenv.Get(key).String(), "TEST")
	})
}

func Test_GBEnv_GetVar(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		value := gbconv.String(gbtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		t.AssertNil(err)
		t.AssertEQ(gbenv.Get(key).String(), "TEST")
	})
}

func Test_GBEnv_Contains(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		value := gbconv.String(gbtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		t.AssertNil(err)
		t.AssertEQ(gbenv.Contains(key), true)
		t.AssertEQ(gbenv.Contains("none"), false)
	})
}

func Test_GBEnv_Set(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		value := gbconv.String(gbtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := gbenv.Set(key, "TEST")
		t.AssertNil(err)
		t.AssertEQ(os.Getenv(key), "TEST")
	})
}

func Test_GBEnv_SetMap(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		err := gbenv.SetMap(g.MapStrStr{
			"K1": "TEST1",
			"K2": "TEST2",
		})
		t.AssertNil(err)
		t.AssertEQ(os.Getenv("K1"), "TEST1")
		t.AssertEQ(os.Getenv("K2"), "TEST2")
	})
}

func Test_GBEnv_Build(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := gbenv.Build(map[string]string{
			"k1": "v1",
			"k2": "v2",
		})
		t.AssertIN("k1=v1", s)
		t.AssertIN("k2=v2", s)
	})
}

func Test_GBEnv_Remove(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		value := gbconv.String(gbtime.TimestampNano())
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		t.AssertNil(err)
		err = gbenv.Remove(key)
		t.AssertNil(err)
		t.AssertEQ(os.Getenv(key), "")
	})
}

func Test_GetWithCmd(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		gbcmd.Init("-test", "2")
		t.Assert(gbenv.GetWithCmd("TEST"), 2)
	})
	gbtest.C(t, func(t *gbtest.T) {
		gbenv.Set("TEST", "1")
		defer gbenv.Remove("TEST")
		gbcmd.Init("-test", "2")
		t.Assert(gbenv.GetWithCmd("test"), 1)
	})
}

func Test_MapFromEnv(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		m := gbenv.MapFromEnv([]string{"a=1", "b=2"})
		t.Assert(m, g.Map{"a": 1, "b": 2})
	})
}

func Test_MapToEnv(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := gbenv.MapToEnv(g.MapStrStr{"a": "1"})
		t.Assert(s, []string{"a=1"})
	})
}

func Test_Filter(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := gbenv.Filter([]string{"a=1", "a=3"})
		t.Assert(s, []string{"a=3"})
	})
}
