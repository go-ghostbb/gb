package gbutil_test

import (
	"context"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbutil "ghostbb.io/gb/util/gb_util"
	"reflect"
	"testing"
)

var (
	ctx = context.TODO()
)

func Test_Try(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := `gbutil Try test`
		t.Assert(gbutil.Try(ctx, func(ctx context.Context) {
			panic(s)
		}), s)
	})
	gbtest.C(t, func(t *gbtest.T) {
		s := `gbutil Try test`
		t.Assert(gbutil.Try(ctx, func(ctx context.Context) {
			panic(gberror.New(s))
		}), s)
	})
}

func Test_TryCatch(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		gbutil.TryCatch(ctx, func(ctx context.Context) {
			panic("gbutil TryCatch test")
		}, nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		gbutil.TryCatch(ctx, func(ctx context.Context) {
			panic("gbutil TryCatch test")

		}, func(ctx context.Context, err error) {
			t.Assert(err, "gbutil TryCatch test")
		})
	})

	gbtest.C(t, func(t *gbtest.T) {
		gbutil.TryCatch(ctx, func(ctx context.Context) {
			panic(gberror.New("gbutil TryCatch test"))

		}, func(ctx context.Context, err error) {
			t.Assert(err, "gbutil TryCatch test")
		})
	})
}

func Test_IsEmpty(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.IsEmpty(1), false)
	})
}

func Test_Throw(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer func() {
			t.Assert(recover(), "gbutil Throw test")
		}()

		gbutil.Throw("gbutil Throw test")
	})
}

func Test_Keys(t *testing.T) {
	// not support int
	gbtest.C(t, func(t *gbtest.T) {
		var val int = 1
		keys := gbutil.Keys(reflect.ValueOf(val))
		t.AssertEQ(len(keys), 0)
	})
	// map
	gbtest.C(t, func(t *gbtest.T) {
		keys := gbutil.Keys(map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", keys)
		t.AssertIN("2", keys)

		strKeys := gbutil.Keys(map[string]interface{}{
			"key1": 1,
			"key2": 2,
		})
		t.AssertIN("key1", strKeys)
		t.AssertIN("key2", strKeys)
	})
	// *map
	gbtest.C(t, func(t *gbtest.T) {
		keys := gbutil.Keys(&map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", keys)
		t.AssertIN("2", keys)
	})
	// *struct
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			A string
			B int
		}
		keys := gbutil.Keys(new(T))
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
	// *struct nil
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			A string
			B int
		}
		var pointer *T
		keys := gbutil.Keys(pointer)
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
	// **struct nil
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			A string
			B int
		}
		var pointer *T
		keys := gbutil.Keys(&pointer)
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
}

func Test_Values(t *testing.T) {
	// not support int
	gbtest.C(t, func(t *gbtest.T) {
		var val int = 1
		keys := gbutil.Values(reflect.ValueOf(val))
		t.AssertEQ(len(keys), 0)
	})
	// map
	gbtest.C(t, func(t *gbtest.T) {
		values := gbutil.Values(map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN(10, values)
		t.AssertIN(20, values)

		values = gbutil.Values(map[string]interface{}{
			"key1": 10,
			"key2": 20,
		})
		t.AssertIN(10, values)
		t.AssertIN(20, values)
	})
	// *map
	gbtest.C(t, func(t *gbtest.T) {
		keys := gbutil.Values(&map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN(10, keys)
		t.AssertIN(20, keys)
	})
	// struct
	gbtest.C(t, func(t *gbtest.T) {
		type T struct {
			A string
			B int
		}
		keys := gbutil.Values(T{
			A: "1",
			B: 2,
		})
		t.Assert(keys, g.Slice{"1", 2})
	})
}

func TestListToMapByKey(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		listMap := []map[string]interface{}{
			{"key1": 1, "key2": 2},
			{"key3": 3, "key4": 4},
		}
		t.Assert(gbutil.ListToMapByKey(listMap, "key1"), "{\"1\":{\"key1\":1,\"key2\":2}}")
	})
}

func Test_GetOrDefaultStr(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.GetOrDefaultStr("a", "b"), "b")
		t.Assert(gbutil.GetOrDefaultStr("a", "b", "c"), "b")
		t.Assert(gbutil.GetOrDefaultStr("a"), "a")
	})
}

func Test_GetOrDefaultAny(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbutil.GetOrDefaultAny("a", "b"), "b")
		t.Assert(gbutil.GetOrDefaultAny("a", "b", "c"), "b")
		t.Assert(gbutil.GetOrDefaultAny("a"), "a")
	})
}
