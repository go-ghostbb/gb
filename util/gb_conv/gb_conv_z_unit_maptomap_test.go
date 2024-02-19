package gbconv_test

import (
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
)

func Test_MapToMap1(t *testing.T) {
	// map[int]int -> map[string]string
	// empty original map.
	gbtest.C(t, func(t *gbtest.T) {
		m1 := g.MapIntInt{}
		m2 := g.MapStrStr{}
		t.Assert(gbconv.MapToMap(m1, &m2), nil)
		t.Assert(len(m1), len(m2))
	})
	// map[int]int -> map[string]string
	gbtest.C(t, func(t *gbtest.T) {
		m1 := g.MapIntInt{
			1: 100,
			2: 200,
		}
		m2 := g.MapStrStr{}
		t.Assert(gbconv.MapToMap(m1, &m2), nil)
		t.Assert(m2["1"], m1[1])
		t.Assert(m2["2"], m1[2])
	})
	// map[string]interface{} -> map[string]string
	gbtest.C(t, func(t *gbtest.T) {
		m1 := g.Map{
			"k1": "v1",
			"k2": "v2",
		}
		m2 := g.MapStrStr{}
		t.Assert(gbconv.MapToMap(m1, &m2), nil)
		t.Assert(m2["k1"], m1["k1"])
		t.Assert(m2["k2"], m1["k2"])
	})
	// map[string]string -> map[string]interface{}
	gbtest.C(t, func(t *gbtest.T) {
		m1 := g.MapStrStr{
			"k1": "v1",
			"k2": "v2",
		}
		m2 := g.Map{}
		t.Assert(gbconv.MapToMap(m1, &m2), nil)
		t.Assert(m2["k1"], m1["k1"])
		t.Assert(m2["k2"], m1["k2"])
	})
	// map[string]interface{} -> map[interface{}]interface{}
	gbtest.C(t, func(t *gbtest.T) {
		m1 := g.MapStrStr{
			"k1": "v1",
			"k2": "v2",
		}
		m2 := g.MapAnyAny{}
		t.Assert(gbconv.MapToMap(m1, &m2), nil)
		t.Assert(m2["k1"], m1["k1"])
		t.Assert(m2["k2"], m1["k2"])
	})
	// string -> map[string]interface{}
	gbtest.C(t, func(t *gbtest.T) {
		jsonStr := `{"id":100, "name":"john"}`

		m1 := g.MapStrAny{}
		t.Assert(gbconv.MapToMap(jsonStr, &m1), nil)
		t.Assert(m1["id"], 100)

		m2 := g.MapStrAny{}
		t.Assert(gbconv.MapToMap([]byte(jsonStr), &m2), nil)
		t.Assert(m2["id"], 100)
	})
}

func Test_MapToMap2(t *testing.T) {
	type User struct {
		Id   int
		Name string
	}
	params := g.Map{
		"key": g.Map{
			"id":   1,
			"name": "john",
		},
	}
	gbtest.C(t, func(t *gbtest.T) {
		m := make(map[string]User)
		err := gbconv.MapToMap(params, &m)
		t.AssertNil(err)
		t.Assert(len(m), 1)
		t.Assert(m["key"].Id, 1)
		t.Assert(m["key"].Name, "john")
	})
	gbtest.C(t, func(t *gbtest.T) {
		m := (map[string]User)(nil)
		err := gbconv.MapToMap(params, &m)
		t.AssertNil(err)
		t.Assert(len(m), 1)
		t.Assert(m["key"].Id, 1)
		t.Assert(m["key"].Name, "john")
	})
	gbtest.C(t, func(t *gbtest.T) {
		m := make(map[string]*User)
		err := gbconv.MapToMap(params, &m)
		t.AssertNil(err)
		t.Assert(len(m), 1)
		t.Assert(m["key"].Id, 1)
		t.Assert(m["key"].Name, "john")
	})
	gbtest.C(t, func(t *gbtest.T) {
		m := (map[string]*User)(nil)
		err := gbconv.MapToMap(params, &m)
		t.AssertNil(err)
		t.Assert(len(m), 1)
		t.Assert(m["key"].Id, 1)
		t.Assert(m["key"].Name, "john")
	})
}

func Test_MapToMapDeep(t *testing.T) {
	type Ids struct {
		Id  int
		Uid int
	}
	type Base struct {
		Ids
		Time string
	}
	type User struct {
		Base
		Name string
	}
	params := g.Map{
		"key": g.Map{
			"id":   1,
			"name": "john",
		},
	}
	gbtest.C(t, func(t *gbtest.T) {
		m := (map[string]*User)(nil)
		err := gbconv.MapToMap(params, &m)
		t.AssertNil(err)
		t.Assert(len(m), 1)
		t.Assert(m["key"].Id, 1)
		t.Assert(m["key"].Name, "john")
	})
}

func Test_MapToMaps(t *testing.T) {
	params := g.Slice{
		g.Map{"id": 1, "name": "john"},
		g.Map{"id": 2, "name": "smith"},
	}
	gbtest.C(t, func(t *gbtest.T) {
		var s []g.Map
		err := gbconv.MapToMaps(params, &s)
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s, params)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var s []*g.Map
		err := gbconv.MapToMaps(params, &s)
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s, params)
	})
	gbtest.C(t, func(t *gbtest.T) {
		jsonStr := `[{"id":100, "name":"john"},{"id":200, "name":"smith"}]`

		var m1 []g.Map
		t.Assert(gbconv.MapToMaps(jsonStr, &m1), nil)
		t.Assert(m1[0]["id"], 100)
		t.Assert(m1[1]["id"], 200)

		t.Assert(gbconv.MapToMaps([]byte(jsonStr), &m1), nil)
		t.Assert(m1[0]["id"], 100)
		t.Assert(m1[1]["id"], 200)
	})
}

func Test_MapToMaps_StructParams(t *testing.T) {
	type User struct {
		Id   int
		Name string
	}
	params := g.Slice{
		User{1, "name1"},
		User{2, "name2"},
	}
	gbtest.C(t, func(t *gbtest.T) {
		var s []g.Map
		err := gbconv.MapToMaps(params, &s)
		t.AssertNil(err)
		t.Assert(len(s), 2)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var s []*g.Map
		err := gbconv.MapToMaps(params, &s)
		t.AssertNil(err)
		t.Assert(len(s), 2)
	})
}
