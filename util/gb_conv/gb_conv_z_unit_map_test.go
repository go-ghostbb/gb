package gbconv_test

import (
	"encoding/json"
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbutil "ghostbb.io/gb/util/gb_util"
	"testing"

	"gopkg.in/yaml.v3"
)

func Test_Map_Basic(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		m1 := map[string]string{
			"k": "v",
		}
		m2 := map[int]string{
			3: "v",
		}
		m3 := map[float64]float32{
			1.22: 3.1,
		}
		t.Assert(gbconv.Map(m1), g.Map{
			"k": "v",
		})
		t.Assert(gbconv.Map(m2), g.Map{
			"3": "v",
		})
		t.Assert(gbconv.Map(m3), g.Map{
			"1.22": "3.1",
		})
		t.Assert(gbconv.Map(`{"name":"goframe"}`), g.Map{
			"name": "goframe",
		})
		t.Assert(gbconv.Map(`{"name":"goframe"`), nil)
		t.Assert(gbconv.Map(`{goframe}`), nil)
		t.Assert(gbconv.Map([]byte(`{"name":"goframe"}`)), g.Map{
			"name": "goframe",
		})
		t.Assert(gbconv.Map([]byte(`{"name":"goframe"`)), nil)
		t.Assert(gbconv.Map([]byte(`{goframe}`)), nil)
	})
}

func Test_Map_Slice(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		slice1 := g.Slice{"1", "2", "3", "4"}
		slice2 := g.Slice{"1", "2", "3"}
		slice3 := g.Slice{}
		t.Assert(gbconv.Map(slice1), g.Map{
			"1": "2",
			"3": "4",
		})
		t.Assert(gbconv.Map(slice2), g.Map{
			"1": "2",
			"3": nil,
		})
		t.Assert(gbconv.Map(slice3), g.Map{})
	})
}

func Test_Maps_Basic(t *testing.T) {
	params := g.Slice{
		g.Map{"id": 100, "name": "john"},
		g.Map{"id": 200, "name": "smith"},
	}
	gbtest.C(t, func(t *gbtest.T) {
		list := gbconv.Maps(params)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})
	gbtest.C(t, func(t *gbtest.T) {
		list := gbconv.SliceMap(params)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})
	gbtest.C(t, func(t *gbtest.T) {
		list := gbconv.SliceMapDeep(params)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})
	gbtest.C(t, func(t *gbtest.T) {
		type Base struct {
			Age int
		}
		type User struct {
			Id   int
			Name string
			Base
		}

		users := make([]User, 0)
		params := []g.Map{
			{"id": 1, "name": "john", "age": 18},
			{"id": 2, "name": "smith", "age": 20},
		}
		err := gbconv.SliceStruct(params, &users)
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

func Test_Maps_JsonStr(t *testing.T) {
	jsonStr := `[{"id":100, "name":"john"},{"id":200, "name":"smith"}]`
	gbtest.C(t, func(t *gbtest.T) {
		list := gbconv.Maps(jsonStr)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)

		list = gbconv.Maps([]byte(jsonStr))
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbconv.Maps(`[id]`), nil)
		t.Assert(gbconv.Maps(`test`), nil)
		t.Assert(gbconv.Maps([]byte(`[id]`)), nil)
		t.Assert(gbconv.Maps([]byte(`test`)), nil)
	})
}

func Test_Map_StructWithGConvTag(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Uid      int
			Name     string
			SiteUrl  string `gbconv:"-"`
			NickName string `gbconv:"nickname, omitempty"`
			Pass1    string `gbconv:"password1"`
			Pass2    string `gbconv:"password2"`
		}
		user1 := User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
		}
		user2 := &user1
		map1 := gbconv.Map(user1)
		map2 := gbconv.Map(user2)
		t.Assert(map1["Uid"], 100)
		t.Assert(map1["Name"], "john")
		t.Assert(map1["SiteUrl"], nil)
		t.Assert(map1["NickName"], nil)
		t.Assert(map1["nickname"], nil)
		t.Assert(map1["password1"], "123")
		t.Assert(map1["password2"], "456")

		t.Assert(map2["Uid"], 100)
		t.Assert(map2["Name"], "john")
		t.Assert(map2["SiteUrl"], nil)
		t.Assert(map2["NickName"], nil)
		t.Assert(map2["nickname"], nil)
		t.Assert(map2["password1"], "123")
		t.Assert(map2["password2"], "456")
	})
}

func Test_Map_StructWithJsonTag(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Uid      int
			Name     string
			SiteUrl  string `json:"-"`
			NickName string `json:"nickname, omitempty"`
			Pass1    string `json:"password1"`
			Pass2    string `json:"password2"`
		}
		user1 := User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
		}
		user2 := &user1
		map1 := gbconv.Map(user1)
		map2 := gbconv.Map(user2)
		t.Assert(map1["Uid"], 100)
		t.Assert(map1["Name"], "john")
		t.Assert(map1["SiteUrl"], nil)
		t.Assert(map1["NickName"], nil)
		t.Assert(map1["nickname"], nil)
		t.Assert(map1["password1"], "123")
		t.Assert(map1["password2"], "456")

		t.Assert(map2["Uid"], 100)
		t.Assert(map2["Name"], "john")
		t.Assert(map2["SiteUrl"], nil)
		t.Assert(map2["NickName"], nil)
		t.Assert(map2["nickname"], nil)
		t.Assert(map2["password1"], "123")
		t.Assert(map2["password2"], "456")
	})
}

func Test_Map_StructWithCTag(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Uid      int
			Name     string
			SiteUrl  string `c:"-"`
			NickName string `c:"nickname, omitempty"`
			Pass1    string `c:"password1"`
			Pass2    string `c:"password2"`
		}
		user1 := User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
		}
		user2 := &user1
		map1 := gbconv.Map(user1)
		map2 := gbconv.Map(user2)
		t.Assert(map1["Uid"], 100)
		t.Assert(map1["Name"], "john")
		t.Assert(map1["SiteUrl"], nil)
		t.Assert(map1["NickName"], nil)
		t.Assert(map1["nickname"], nil)
		t.Assert(map1["password1"], "123")
		t.Assert(map1["password2"], "456")

		t.Assert(map2["Uid"], 100)
		t.Assert(map2["Name"], "john")
		t.Assert(map2["SiteUrl"], nil)
		t.Assert(map2["NickName"], nil)
		t.Assert(map2["nickname"], nil)
		t.Assert(map2["password1"], "123")
		t.Assert(map2["password2"], "456")
	})
}

func Test_Map_PrivateAttribute(t *testing.T) {
	type User struct {
		Id   int
		name string
	}
	gbtest.C(t, func(t *gbtest.T) {
		user := &User{1, "john"}
		t.Assert(gbconv.Map(user), g.Map{"Id": 1})
	})
}

func Test_Map_Embedded(t *testing.T) {
	type Base struct {
		Id int
	}
	type User struct {
		Base
		Name string
	}
	type UserDetail struct {
		User
		Brief string
	}
	gbtest.C(t, func(t *gbtest.T) {
		user := &User{}
		user.Id = 1
		user.Name = "john"

		m := gbconv.Map(user)
		t.Assert(len(m), 2)
		t.Assert(m["Id"], user.Id)
		t.Assert(m["Name"], user.Name)
	})
	gbtest.C(t, func(t *gbtest.T) {
		user := &UserDetail{}
		user.Id = 1
		user.Name = "john"
		user.Brief = "john guo"

		m := gbconv.Map(user)
		t.Assert(len(m), 3)
		t.Assert(m["Id"], user.Id)
		t.Assert(m["Name"], user.Name)
		t.Assert(m["Brief"], user.Brief)
	})
}

func Test_Map_Embedded2(t *testing.T) {
	type Ids struct {
		Id  int `c:"id"`
		Uid int `c:"uid"`
	}
	type Base struct {
		Ids
		CreateTime string `c:"create_time"`
	}
	type User struct {
		Base
		Passport string `c:"passport"`
		Password string `c:"password"`
		Nickname string `c:"nickname"`
	}
	gbtest.C(t, func(t *gbtest.T) {
		user := new(User)
		user.Id = 100
		user.Nickname = "john"
		user.CreateTime = "2019"
		m := gbconv.Map(user)
		t.Assert(m["id"], "100")
		t.Assert(m["nickname"], user.Nickname)
		t.Assert(m["create_time"], "2019")
	})
	gbtest.C(t, func(t *gbtest.T) {
		user := new(User)
		user.Id = 100
		user.Nickname = "john"
		user.CreateTime = "2019"
		m := gbconv.MapDeep(user)
		t.Assert(m["id"], user.Id)
		t.Assert(m["nickname"], user.Nickname)
		t.Assert(m["create_time"], user.CreateTime)
	})
}

func Test_MapDeep2(t *testing.T) {
	type A struct {
		F string
		G string
	}

	type B struct {
		A
		H string
	}

	type C struct {
		A A
		F string
	}

	type D struct {
		I A
		F string
	}

	gbtest.C(t, func(t *gbtest.T) {
		b := new(B)
		c := new(C)
		d := new(D)
		mb := gbconv.MapDeep(b)
		mc := gbconv.MapDeep(c)
		md := gbconv.MapDeep(d)
		t.Assert(gbutil.MapContains(mb, "F"), true)
		t.Assert(gbutil.MapContains(mb, "G"), true)
		t.Assert(gbutil.MapContains(mb, "H"), true)
		t.Assert(gbutil.MapContains(mc, "A"), true)
		t.Assert(gbutil.MapContains(mc, "F"), true)
		t.Assert(gbutil.MapContains(mc, "G"), false)
		t.Assert(gbutil.MapContains(md, "F"), true)
		t.Assert(gbutil.MapContains(md, "I"), true)
		t.Assert(gbutil.MapContains(md, "H"), false)
		t.Assert(gbutil.MapContains(md, "G"), false)
	})
}

func Test_MapDeep3(t *testing.T) {
	type Base struct {
		Id   int    `c:"id"`
		Date string `c:"date"`
	}
	type User struct {
		UserBase Base   `c:"base"`
		Passport string `c:"passport"`
		Password string `c:"password"`
		Nickname string `c:"nickname"`
	}

	gbtest.C(t, func(t *gbtest.T) {
		user := &User{
			UserBase: Base{
				Id:   1,
				Date: "2019-10-01",
			},
			Passport: "john",
			Password: "123456",
			Nickname: "JohnGuo",
		}
		m := gbconv.MapDeep(user)
		t.Assert(m, g.Map{
			"base": g.Map{
				"id":   user.UserBase.Id,
				"date": user.UserBase.Date,
			},
			"passport": user.Passport,
			"password": user.Password,
			"nickname": user.Nickname,
		})
	})

	gbtest.C(t, func(t *gbtest.T) {
		user := &User{
			UserBase: Base{
				Id:   1,
				Date: "2019-10-01",
			},
			Passport: "john",
			Password: "123456",
			Nickname: "JohnGuo",
		}
		m := gbconv.Map(user)
		t.Assert(m, g.Map{
			"base":     user.UserBase,
			"passport": user.Passport,
			"password": user.Password,
			"nickname": user.Nickname,
		})
	})
}

func Test_MapDeepWithAttributeTag(t *testing.T) {
	type Ids struct {
		Id  int `c:"id"`
		Uid int `c:"uid"`
	}
	type Base struct {
		Ids        `json:"ids"`
		CreateTime string `c:"create_time"`
	}
	type User struct {
		Base     `json:"base"`
		Passport string `c:"passport"`
		Password string `c:"password"`
		Nickname string `c:"nickname"`
	}
	gbtest.C(t, func(t *gbtest.T) {
		user := new(User)
		user.Id = 100
		user.Nickname = "john"
		user.CreateTime = "2019"
		m := gbconv.Map(user)
		t.Assert(m["id"], "")
		t.Assert(m["nickname"], user.Nickname)
		t.Assert(m["create_time"], "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		user := new(User)
		user.Id = 100
		user.Nickname = "john"
		user.CreateTime = "2019"
		m := gbconv.MapDeep(user)
		t.Assert(m["base"].(map[string]interface{})["ids"].(map[string]interface{})["id"], user.Id)
		t.Assert(m["nickname"], user.Nickname)
		t.Assert(m["base"].(map[string]interface{})["create_time"], user.CreateTime)
	})
}

func Test_MapDeepWithNestedMapAnyAny(t *testing.T) {
	type User struct {
		ExtraAttributes g.Map `c:"extra_attributes"`
	}

	gbtest.C(t, func(t *gbtest.T) {
		user := &User{
			ExtraAttributes: g.Map{
				"simple_attribute": 123,
				"map_string_attribute": g.Map{
					"inner_value": 456,
				},
				"map_interface_attribute": g.MapAnyAny{
					"inner_value": 456,
					123:           "integer_key_should_be_converted_to_string",
				},
			},
		}
		m := gbconv.MapDeep(user)
		t.Assert(m, g.Map{
			"extra_attributes": g.Map{
				"simple_attribute": 123,
				"map_string_attribute": g.Map{
					"inner_value": user.ExtraAttributes["map_string_attribute"].(g.Map)["inner_value"],
				},
				"map_interface_attribute": g.Map{
					"inner_value": user.ExtraAttributes["map_interface_attribute"].(g.MapAnyAny)["inner_value"],
					"123":         "integer_key_should_be_converted_to_string",
				},
			},
		})
	})

	type Outer struct {
		OuterStruct map[string]interface{} `c:"outer_struct" yaml:"outer_struct"`
		Field3      map[string]interface{} `c:"field3" yaml:"field3"`
	}

	gbtest.C(t, func(t *gbtest.T) {
		problemYaml := []byte(`
outer_struct:
  field1: &anchor1
    inner1: 123
    inner2: 345
  field2: 
    inner3: 456
    inner4: 789
    <<: *anchor1
field3:
  123: integer_key
`)
		parsed := &Outer{}

		err := yaml.Unmarshal(problemYaml, parsed)
		t.AssertNil(err)

		_, err = json.Marshal(parsed)
		t.AssertNil(err)

		converted := gbconv.MapDeep(parsed)
		jsonData, err := json.Marshal(converted)
		t.AssertNil(err)

		t.Assert(string(jsonData), `{"field3":{"123":"integer_key"},"outer_struct":{"field1":{"inner1":123,"inner2":345},"field2":{"inner1":123,"inner2":345,"inner3":456,"inner4":789}}}`)
	})
}

func TestMapStrStr(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbconv.MapStrStr(map[string]string{"k": "v"}), map[string]string{"k": "v"})
		t.Assert(gbconv.MapStrStr(`{}`), nil)
	})
}

func TestMapStrStrDeep(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbconv.MapStrStrDeep(map[string]string{"k": "v"}), map[string]string{"k": "v"})
		t.Assert(gbconv.MapStrStrDeep(`{"k":"v"}`), map[string]string{"k": "v"})
		t.Assert(gbconv.MapStrStrDeep(`{}`), nil)
	})
}

func TestMapsDeep(t *testing.T) {
	jsonStr := `[{"id":100, "name":"john"},{"id":200, "name":"smith"}]`
	params := g.Slice{
		g.Map{"id": 100, "name": "john"},
		g.Map{"id": 200, "name": "smith"},
	}

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbconv.MapsDeep(nil), nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		list := gbconv.MapsDeep(params)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})

	gbtest.C(t, func(t *gbtest.T) {
		list := gbconv.MapsDeep(jsonStr)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)

		list = gbconv.MapsDeep([]byte(jsonStr))
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})

	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbconv.MapsDeep(`[id]`), nil)
		t.Assert(gbconv.MapsDeep(`test`), nil)
		t.Assert(gbconv.MapsDeep([]byte(`[id]`)), nil)
		t.Assert(gbconv.MapsDeep([]byte(`test`)), nil)
		t.Assert(gbconv.MapsDeep([]string{}), nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		stringInterfaceMapList := make([]map[string]interface{}, 0)
		stringInterfaceMapList = append(stringInterfaceMapList, map[string]interface{}{"id": 100})
		stringInterfaceMapList = append(stringInterfaceMapList, map[string]interface{}{"id": 200})
		list := gbconv.MapsDeep(stringInterfaceMapList)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)

		list = gbconv.MapsDeep([]byte(jsonStr))
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})
}

func TestMapWithJsonOmitEmpty(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type S struct {
			Key   string      `json:",omitempty"`
			Value interface{} `json:",omitempty"`
		}
		s := S{
			Key:   "",
			Value: 1,
		}
		m1 := gbconv.Map(s)
		t.Assert(m1, g.Map{
			"Key":   "",
			"Value": 1,
		})

		m2 := gbconv.Map(s, gbconv.MapOption{
			Deep:      false,
			OmitEmpty: true,
			Tags:      nil,
		})
		t.Assert(m2, g.Map{
			"Value": 1,
		})
	})

	gbtest.C(t, func(t *gbtest.T) {
		type ProductConfig struct {
			Pid      int `v:"required" json:"pid,omitempty"`
			TimeSpan int `v:"required" json:"timeSpan,omitempty"`
		}
		type CreateGoodsDetail struct {
			ProductConfig
			AutoRenewFlag int `v:"required" json:"autoRenewFlag"`
		}
		s := &CreateGoodsDetail{
			ProductConfig: ProductConfig{
				Pid:      1,
				TimeSpan: 0,
			},
			AutoRenewFlag: 0,
		}
		m1 := gbconv.Map(s)
		t.Assert(m1, g.Map{
			"pid":           1,
			"timeSpan":      0,
			"autoRenewFlag": 0,
		})

		m2 := gbconv.Map(s, gbconv.MapOption{
			Deep:      false,
			OmitEmpty: true,
			Tags:      nil,
		})
		t.Assert(m2, g.Map{
			"pid":           1,
			"autoRenewFlag": 0,
		})
	})
}
