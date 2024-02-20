package gbconv_test

import (
	gbvar "ghostbb.io/gb/container/gb_var"
	gbjson "ghostbb.io/gb/encoding/gb_json"
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
)

func Test_Slice(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		value := 123.456
		t.AssertEQ(gbconv.Bytes("123"), []byte("123"))
		t.AssertEQ(gbconv.Bytes([]interface{}{1}), []byte{1})
		t.AssertEQ(gbconv.Bytes([]interface{}{300}), []byte("[300]"))
		t.AssertEQ(gbconv.Strings(value), []string{"123.456"})
		t.AssertEQ(gbconv.SliceStr(value), []string{"123.456"})
		t.AssertEQ(gbconv.SliceInt(value), []int{123})
		t.AssertEQ(gbconv.SliceUint(value), []uint{123})
		t.AssertEQ(gbconv.SliceUint32(value), []uint32{123})
		t.AssertEQ(gbconv.SliceUint64(value), []uint64{123})
		t.AssertEQ(gbconv.SliceInt32(value), []int32{123})
		t.AssertEQ(gbconv.SliceInt64(value), []int64{123})
		t.AssertEQ(gbconv.Ints(value), []int{123})
		t.AssertEQ(gbconv.SliceFloat(value), []float64{123.456})
		t.AssertEQ(gbconv.Floats(value), []float64{123.456})
		t.AssertEQ(gbconv.SliceFloat32(value), []float32{123.456})
		t.AssertEQ(gbconv.SliceFloat64(value), []float64{123.456})
		t.AssertEQ(gbconv.Interfaces(value), []interface{}{123.456})
		t.AssertEQ(gbconv.SliceAny(" [26, 27] "), []interface{}{26, 27})
	})
	gbtest.C(t, func(t *gbtest.T) {
		s := gbvar.Vars{
			gbvar.New(1),
			gbvar.New(2),
		}
		t.AssertEQ(gbconv.SliceInt64(s), []int64{1, 2})
	})
}

func Test_Slice_Ints(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Ints(nil), nil)
		t.AssertEQ(gbconv.Ints("[26, 27]"), []int{26, 27})
		t.AssertEQ(gbconv.Ints(" [26, 27] "), []int{26, 27})
		t.AssertEQ(gbconv.Ints([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []int{0, 0})
		t.AssertEQ(gbconv.Ints([]bool{true, false}), []int{1, 0})
		t.AssertEQ(gbconv.Ints([][]byte{{byte(1)}, {byte(2)}}), []int{1, 2})
	})
}

func Test_Slice_Int32s(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Int32s(nil), nil)
		t.AssertEQ(gbconv.Int32s(" [26, 27] "), []int32{26, 27})
		t.AssertEQ(gbconv.Int32s([]string{"1", "2"}), []int32{1, 2})
		t.AssertEQ(gbconv.Int32s([]int{1, 2}), []int32{1, 2})
		t.AssertEQ(gbconv.Int32s([]int8{1, 2}), []int32{1, 2})
		t.AssertEQ(gbconv.Int32s([]int16{1, 2}), []int32{1, 2})
		t.AssertEQ(gbconv.Int32s([]int32{1, 2}), []int32{1, 2})
		t.AssertEQ(gbconv.Int32s([]int64{1, 2}), []int32{1, 2})
		t.AssertEQ(gbconv.Int32s([]uint{1, 2}), []int32{1, 2})
		t.AssertEQ(gbconv.Int32s([]uint8{1, 2}), []int32{1, 2})
		t.AssertEQ(gbconv.Int32s([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []int32{0, 0})
		t.AssertEQ(gbconv.Int32s([]uint16{1, 2}), []int32{1, 2})
		t.AssertEQ(gbconv.Int32s([]uint32{1, 2}), []int32{1, 2})
		t.AssertEQ(gbconv.Int32s([]uint64{1, 2}), []int32{1, 2})
		t.AssertEQ(gbconv.Int32s([]bool{true, false}), []int32{1, 0})
		t.AssertEQ(gbconv.Int32s([]float32{1, 2}), []int32{1, 2})
		t.AssertEQ(gbconv.Int32s([]float64{1, 2}), []int32{1, 2})
		t.AssertEQ(gbconv.Int32s([][]byte{{byte(1)}, {byte(2)}}), []int32{1, 2})

		s := gbvar.Vars{
			gbvar.New(1),
			gbvar.New(2),
		}
		t.AssertEQ(gbconv.SliceInt32(s), []int32{1, 2})
	})
}

func Test_Slice_Int64s(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Int64s(nil), nil)
		t.AssertEQ(gbconv.Int64s(" [26, 27] "), []int64{26, 27})
		t.AssertEQ(gbconv.Int64s([]string{"1", "2"}), []int64{1, 2})
		t.AssertEQ(gbconv.Int64s([]int{1, 2}), []int64{1, 2})
		t.AssertEQ(gbconv.Int64s([]int8{1, 2}), []int64{1, 2})
		t.AssertEQ(gbconv.Int64s([]int16{1, 2}), []int64{1, 2})
		t.AssertEQ(gbconv.Int64s([]int32{1, 2}), []int64{1, 2})
		t.AssertEQ(gbconv.Int64s([]int64{1, 2}), []int64{1, 2})
		t.AssertEQ(gbconv.Int64s([]uint{1, 2}), []int64{1, 2})
		t.AssertEQ(gbconv.Int64s([]uint8{1, 2}), []int64{1, 2})
		t.AssertEQ(gbconv.Int64s([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []int64{0, 0})
		t.AssertEQ(gbconv.Int64s([]uint16{1, 2}), []int64{1, 2})
		t.AssertEQ(gbconv.Int64s([]uint32{1, 2}), []int64{1, 2})
		t.AssertEQ(gbconv.Int64s([]uint64{1, 2}), []int64{1, 2})
		t.AssertEQ(gbconv.Int64s([]bool{true, false}), []int64{1, 0})
		t.AssertEQ(gbconv.Int64s([]float32{1, 2}), []int64{1, 2})
		t.AssertEQ(gbconv.Int64s([]float64{1, 2}), []int64{1, 2})
		t.AssertEQ(gbconv.Int64s([][]byte{{byte(1)}, {byte(2)}}), []int64{1, 2})

		s := gbvar.Vars{
			gbvar.New(1),
			gbvar.New(2),
		}
		t.AssertEQ(gbconv.Int64s(s), []int64{1, 2})
	})
}

func Test_Slice_Uints(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Uints(nil), nil)
		t.AssertEQ(gbconv.Uints("1"), []uint{1})
		t.AssertEQ(gbconv.Uints(" [26, 27] "), []uint{26, 27})
		t.AssertEQ(gbconv.Uints([]string{"1", "2"}), []uint{1, 2})
		t.AssertEQ(gbconv.Uints([]int{1, 2}), []uint{1, 2})
		t.AssertEQ(gbconv.Uints([]int8{1, 2}), []uint{1, 2})
		t.AssertEQ(gbconv.Uints([]int16{1, 2}), []uint{1, 2})
		t.AssertEQ(gbconv.Uints([]int32{1, 2}), []uint{1, 2})
		t.AssertEQ(gbconv.Uints([]int64{1, 2}), []uint{1, 2})
		t.AssertEQ(gbconv.Uints([]uint{1, 2}), []uint{1, 2})
		t.AssertEQ(gbconv.Uints([]uint8{1, 2}), []uint{1, 2})
		t.AssertEQ(gbconv.Uints([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []uint{0, 0})
		t.AssertEQ(gbconv.Uints([]uint16{1, 2}), []uint{1, 2})
		t.AssertEQ(gbconv.Uints([]uint32{1, 2}), []uint{1, 2})
		t.AssertEQ(gbconv.Uints([]uint64{1, 2}), []uint{1, 2})
		t.AssertEQ(gbconv.Uints([]bool{true, false}), []uint{1, 0})
		t.AssertEQ(gbconv.Uints([]float32{1, 2}), []uint{1, 2})
		t.AssertEQ(gbconv.Uints([]float64{1, 2}), []uint{1, 2})
		t.AssertEQ(gbconv.Uints([][]byte{{byte(1)}, {byte(2)}}), []uint{1, 2})

		s := gbvar.Vars{
			gbvar.New(1),
			gbvar.New(2),
		}
		t.AssertEQ(gbconv.Uints(s), []uint{1, 2})
	})
}

func Test_Slice_Uint32s(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Uint32s(nil), nil)
		t.AssertEQ(gbconv.Uint32s("1"), []uint32{1})
		t.AssertEQ(gbconv.Uint32s(" [26, 27] "), []uint32{26, 27})
		t.AssertEQ(gbconv.Uint32s([]string{"1", "2"}), []uint32{1, 2})
		t.AssertEQ(gbconv.Uint32s([]int{1, 2}), []uint32{1, 2})
		t.AssertEQ(gbconv.Uint32s([]int8{1, 2}), []uint32{1, 2})
		t.AssertEQ(gbconv.Uint32s([]int16{1, 2}), []uint32{1, 2})
		t.AssertEQ(gbconv.Uint32s([]int32{1, 2}), []uint32{1, 2})
		t.AssertEQ(gbconv.Uint32s([]int64{1, 2}), []uint32{1, 2})
		t.AssertEQ(gbconv.Uint32s([]uint{1, 2}), []uint32{1, 2})
		t.AssertEQ(gbconv.Uint32s([]uint8{1, 2}), []uint32{1, 2})
		t.AssertEQ(gbconv.Uint32s([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []uint32{0, 0})
		t.AssertEQ(gbconv.Uint32s([]uint16{1, 2}), []uint32{1, 2})
		t.AssertEQ(gbconv.Uint32s([]uint32{1, 2}), []uint32{1, 2})
		t.AssertEQ(gbconv.Uint32s([]uint64{1, 2}), []uint32{1, 2})
		t.AssertEQ(gbconv.Uint32s([]bool{true, false}), []uint32{1, 0})
		t.AssertEQ(gbconv.Uint32s([]float32{1, 2}), []uint32{1, 2})
		t.AssertEQ(gbconv.Uint32s([]float64{1, 2}), []uint32{1, 2})
		t.AssertEQ(gbconv.Uint32s([][]byte{{byte(1)}, {byte(2)}}), []uint32{1, 2})

		s := gbvar.Vars{
			gbvar.New(1),
			gbvar.New(2),
		}
		t.AssertEQ(gbconv.Uint32s(s), []uint32{1, 2})
	})
}

func Test_Slice_Uint64s(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Uint64s(nil), nil)
		t.AssertEQ(gbconv.Uint64s("1"), []uint64{1})
		t.AssertEQ(gbconv.Uint64s(" [26, 27] "), []uint64{26, 27})
		t.AssertEQ(gbconv.Uint64s([]string{"1", "2"}), []uint64{1, 2})
		t.AssertEQ(gbconv.Uint64s([]int{1, 2}), []uint64{1, 2})
		t.AssertEQ(gbconv.Uint64s([]int8{1, 2}), []uint64{1, 2})
		t.AssertEQ(gbconv.Uint64s([]int16{1, 2}), []uint64{1, 2})
		t.AssertEQ(gbconv.Uint64s([]int32{1, 2}), []uint64{1, 2})
		t.AssertEQ(gbconv.Uint64s([]int64{1, 2}), []uint64{1, 2})
		t.AssertEQ(gbconv.Uint64s([]uint{1, 2}), []uint64{1, 2})
		t.AssertEQ(gbconv.Uint64s([]uint8{1, 2}), []uint64{1, 2})
		t.AssertEQ(gbconv.Uint64s([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []uint64{0, 0})
		t.AssertEQ(gbconv.Uint64s([]uint16{1, 2}), []uint64{1, 2})
		t.AssertEQ(gbconv.Uint64s([]uint64{1, 2}), []uint64{1, 2})
		t.AssertEQ(gbconv.Uint64s([]uint64{1, 2}), []uint64{1, 2})
		t.AssertEQ(gbconv.Uint64s([]bool{true, false}), []uint64{1, 0})
		t.AssertEQ(gbconv.Uint64s([]float32{1, 2}), []uint64{1, 2})
		t.AssertEQ(gbconv.Uint64s([]float64{1, 2}), []uint64{1, 2})
		t.AssertEQ(gbconv.Uint64s([][]byte{{byte(1)}, {byte(2)}}), []uint64{1, 2})

		s := gbvar.Vars{
			gbvar.New(1),
			gbvar.New(2),
		}
		t.AssertEQ(gbconv.Uint64s(s), []uint64{1, 2})
	})
}

func Test_Slice_Float32s(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Float32s("123.4"), []float32{123.4})
		t.AssertEQ(gbconv.Float32s([]string{"123.4", "123.5"}), []float32{123.4, 123.5})
		t.AssertEQ(gbconv.Float32s([]int{123}), []float32{123})
		t.AssertEQ(gbconv.Float32s([]int8{123}), []float32{123})
		t.AssertEQ(gbconv.Float32s([]int16{123}), []float32{123})
		t.AssertEQ(gbconv.Float32s([]int32{123}), []float32{123})
		t.AssertEQ(gbconv.Float32s([]int64{123}), []float32{123})
		t.AssertEQ(gbconv.Float32s([]uint{123}), []float32{123})
		t.AssertEQ(gbconv.Float32s([]uint8{123}), []float32{123})
		t.AssertEQ(gbconv.Float32s([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []float32{0, 0})
		t.AssertEQ(gbconv.Float32s([]uint16{123}), []float32{123})
		t.AssertEQ(gbconv.Float32s([]uint32{123}), []float32{123})
		t.AssertEQ(gbconv.Float32s([]uint64{123}), []float32{123})
		t.AssertEQ(gbconv.Float32s([]bool{true, false}), []float32{0, 0})
		t.AssertEQ(gbconv.Float32s([]float32{123}), []float32{123})
		t.AssertEQ(gbconv.Float32s([]float64{123}), []float32{123})

		s := gbvar.Vars{
			gbvar.New(1.1),
			gbvar.New(2.1),
		}
		t.AssertEQ(gbconv.SliceFloat32(s), []float32{1.1, 2.1})
	})
}

func Test_Slice_Float64s(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Float64s("123.4"), []float64{123.4})
		t.AssertEQ(gbconv.Float64s([]string{"123.4", "123.5"}), []float64{123.4, 123.5})
		t.AssertEQ(gbconv.Float64s([]int{123}), []float64{123})
		t.AssertEQ(gbconv.Float64s([]int8{123}), []float64{123})
		t.AssertEQ(gbconv.Float64s([]int16{123}), []float64{123})
		t.AssertEQ(gbconv.Float64s([]int32{123}), []float64{123})
		t.AssertEQ(gbconv.Float64s([]int64{123}), []float64{123})
		t.AssertEQ(gbconv.Float64s([]uint{123}), []float64{123})
		t.AssertEQ(gbconv.Float64s([]uint8{123}), []float64{123})
		t.AssertEQ(gbconv.Float64s([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []float64{0, 0})
		t.AssertEQ(gbconv.Float64s([]uint16{123}), []float64{123})
		t.AssertEQ(gbconv.Float64s([]uint32{123}), []float64{123})
		t.AssertEQ(gbconv.Float64s([]uint64{123}), []float64{123})
		t.AssertEQ(gbconv.Float64s([]bool{true, false}), []float64{0, 0})
		t.AssertEQ(gbconv.Float64s([]float32{123}), []float64{123})
		t.AssertEQ(gbconv.Float64s([]float64{123}), []float64{123})
	})
}

func Test_Slice_Empty(t *testing.T) {
	// Int.
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Ints(""), []int{})
		t.Assert(gbconv.Ints(nil), nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Int32s(""), []int32{})
		t.Assert(gbconv.Int32s(nil), nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Int64s(""), []int64{})
		t.Assert(gbconv.Int64s(nil), nil)
	})
	// Uint.
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Uints(""), []uint{})
		t.Assert(gbconv.Uints(nil), nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Uint32s(""), []uint32{})
		t.Assert(gbconv.Uint32s(nil), nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Uint64s(""), []uint64{})
		t.Assert(gbconv.Uint64s(nil), nil)
	})
	// Float.
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Floats(""), []float64{})
		t.Assert(gbconv.Floats(nil), nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Float32s(""), []float32{})
		t.Assert(gbconv.Float32s(nil), nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Float64s(""), []float64{})
		t.Assert(gbconv.Float64s(nil), nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Strings(""), []string{})
		t.Assert(gbconv.Strings(nil), nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.SliceAny(""), []interface{}{""})
		t.Assert(gbconv.SliceAny(nil), nil)
	})
}

func Test_Strings(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		array := []*g.Var{
			g.NewVar(1),
			g.NewVar(2),
			g.NewVar(3),
		}
		t.AssertEQ(gbconv.Strings(array), []string{"1", "2", "3"})

		t.AssertEQ(gbconv.Strings([]uint8(`["1","2"]`)), []string{"1", "2"})
		t.AssertEQ(gbconv.Strings([][]byte{{byte(0)}, {byte(1)}}), []string{"\u0000", "\u0001"})
	})

	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Strings("123"), []string{"123"})
	})
}

func Test_Slice_Interfaces(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		array := gbconv.Interfaces([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`))
		t.Assert(len(array), 2)
		t.Assert(array[0].(g.Map)["id"], 1)
		t.Assert(array[0].(g.Map)["name"], "john")
	})
	// map
	gbtest.C(t, func(t *gbtest.T) {
		array := gbconv.Interfaces(g.Map{
			"id":   1,
			"name": "john",
		})
		t.Assert(len(array), 1)
		t.Assert(array[0].(g.Map)["id"], 1)
		t.Assert(array[0].(g.Map)["name"], "john")
	})
	// struct
	gbtest.C(t, func(t *gbtest.T) {
		type A struct {
			Id   int `json:"id"`
			Name string
		}
		array := gbconv.Interfaces(&A{
			Id:   1,
			Name: "john",
		})
		t.Assert(len(array), 1)
		t.Assert(array[0].(*A).Id, 1)
		t.Assert(array[0].(*A).Name, "john")
	})
}

func Test_Slice_PrivateAttribute(t *testing.T) {
	type User struct {
		Id   int    `json:"id"`
		name string `json:"name"`
	}
	gbtest.C(t, func(t *gbtest.T) {
		user := &User{1, "john"}
		array := gbconv.Interfaces(user)
		t.Assert(len(array), 1)
		t.Assert(array[0].(*User).Id, 1)
		t.Assert(array[0].(*User).name, "john")
	})
}

func Test_Slice_Structs(t *testing.T) {
	type Base struct {
		Age int
	}
	type User struct {
		Id   int
		Name string
		Base
	}

	gbtest.C(t, func(t *gbtest.T) {
		users := make([]User, 0)
		params := []g.Map{
			{"id": 1, "name": "john", "age": 18},
			{"id": 2, "name": "smith", "age": 20},
		}
		err := gbconv.Structs(params, &users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, params[0]["id"])
		t.Assert(users[0].Name, params[0]["name"])
		t.Assert(users[0].Age, 18)

		t.Assert(users[1].Id, params[1]["id"])
		t.Assert(users[1].Name, params[1]["name"])
		t.Assert(users[1].Age, 20)
	})

	gbtest.C(t, func(t *gbtest.T) {
		users := make([]User, 0)
		params := []g.Map{
			{"id": 1, "name": "john", "age": 18},
			{"id": 2, "name": "smith", "age": 20},
		}
		err := gbconv.StructsTag(params, &users, "")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, params[0]["id"])
		t.Assert(users[0].Name, params[0]["name"])
		t.Assert(users[0].Age, 18)

		t.Assert(users[1].Id, params[1]["id"])
		t.Assert(users[1].Name, params[1]["name"])
		t.Assert(users[1].Age, 20)
	})
}

func Test_EmptyString_To_CustomType(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Status string
		type Req struct {
			Name     string
			Statuses []Status
			Types    []string
		}
		var (
			req  *Req
			data = g.Map{
				"Name":     "john",
				"Statuses": "",
				"Types":    "",
			}
		)
		err := gbconv.Scan(data, &req)
		t.AssertNil(err)
		t.Assert(len(req.Statuses), 0)
		t.Assert(len(req.Types), 0)
	})
	gbtest.C(t, func(t *gbtest.T) {
		type Status string
		type Req struct {
			Name     string
			Statuses []*Status
			Types    []string
		}
		var (
			req  *Req
			data = g.Map{
				"Name":     "john",
				"Statuses": "",
				"Types":    "",
			}
		)
		err := gbconv.Scan(data, &req)
		t.AssertNil(err)
		t.Assert(len(req.Statuses), 0)
		t.Assert(len(req.Types), 0)
	})
}

func Test_SliceMap_WithNilMapValue(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			list1 = []map[string]*gbvar.Var{
				{"name": nil},
			}
			list2 []map[string]any
		)
		list2 = gbconv.SliceMap(list1)
		t.Assert(len(list2), 1)
		t.Assert(list1[0], list2[0])
		t.Assert(gbjson.MustEncodeString(list1), gbjson.MustEncodeString(list2))
	})
}
