package gbconv_test

import (
	"fmt"
	gbvar "ghostbb.io/gb/container/gb_var"
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"math/big"
	"testing"
)

func Test_Scan_WithMapParameter(t *testing.T) {
	type User struct {
		Uid  int
		Name string
	}
	gbtest.C(t, func(t *gbtest.T) {
		for i := 0; i < 100; i++ {
			var (
				user   = new(User)
				params = g.Map{
					"uid":    1,
					"myname": "john",
					"name":   "smith",
				}
			)
			err := gbconv.Scan(params, user, g.MapStrStr{
				"myname": "Name",
			})
			t.AssertNil(err)
			t.Assert(user, &User{
				Uid:  1,
				Name: "john",
			})
		}
	})
}

func Test_Scan_StructStructs(t *testing.T) {
	type User struct {
		Uid   int
		Name  string
		Pass1 string `gbconv:"password1"`
		Pass2 string `gbconv:"password2"`
	}
	gbtest.C(t, func(t *gbtest.T) {
		var (
			user   = new(User)
			params = g.Map{
				"uid":   1,
				"name":  "john",
				"PASS1": "123",
				"PASS2": "456",
			}
		)
		err := gbconv.Scan(params, user)
		t.AssertNil(err)
		t.Assert(user, &User{
			Uid:   1,
			Name:  "john",
			Pass1: "123",
			Pass2: "456",
		})
	})
	gbtest.C(t, func(t *gbtest.T) {
		var (
			users  []User
			params = g.Slice{
				g.Map{
					"uid":   1,
					"name":  "john1",
					"PASS1": "111",
					"PASS2": "222",
				},
				g.Map{
					"uid":   2,
					"name":  "john2",
					"PASS1": "333",
					"PASS2": "444",
				},
			}
		)
		err := gbconv.Scan(params, &users)
		t.AssertNil(err)
		t.Assert(users, g.Slice{
			&User{
				Uid:   1,
				Name:  "john1",
				Pass1: "111",
				Pass2: "222",
			},
			&User{
				Uid:   2,
				Name:  "john2",
				Pass1: "333",
				Pass2: "444",
			},
		})
	})
}

func Test_Scan_StructStr(t *testing.T) {
	type User struct {
		Uid   int
		Name  string
		Pass1 string `gbconv:"password1"`
		Pass2 string `gbconv:"password2"`
	}
	gbtest.C(t, func(t *gbtest.T) {
		var (
			user   = new(User)
			params = `{"uid":1,"name":"john", "pass1":"123","pass2":"456"}`
		)
		err := gbconv.Scan(params, user)
		t.AssertNil(err)
		t.Assert(user, &User{
			Uid:   1,
			Name:  "john",
			Pass1: "123",
			Pass2: "456",
		})
	})
	gbtest.C(t, func(t *gbtest.T) {
		var (
			users  []User
			params = `[
{"uid":1,"name":"john1", "pass1":"111","pass2":"222"},
{"uid":2,"name":"john2", "pass1":"333","pass2":"444"}
]`
		)
		err := gbconv.Scan(params, &users)
		t.AssertNil(err)
		t.Assert(users, g.Slice{
			&User{
				Uid:   1,
				Name:  "john1",
				Pass1: "111",
				Pass2: "222",
			},
			&User{
				Uid:   2,
				Name:  "john2",
				Pass1: "333",
				Pass2: "444",
			},
		})
	})
}

func Test_Scan_Map(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var m map[string]string
		data := g.Map{
			"k1": "v1",
			"k2": "v2",
		}
		err := gbconv.Scan(data, &m)
		t.AssertNil(err)
		t.Assert(data, m)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var m map[int]int
		data := g.Map{
			"1": "11",
			"2": "22",
		}
		err := gbconv.Scan(data, &m)
		t.AssertNil(err)
		t.Assert(data, m)
	})
	// json string parameter.
	gbtest.C(t, func(t *gbtest.T) {
		var m map[string]string
		data := `{"k1":"v1","k2":"v2"}`
		err := gbconv.Scan(data, &m)
		t.AssertNil(err)
		t.Assert(m, g.Map{
			"k1": "v1",
			"k2": "v2",
		})
	})
}

func Test_Scan_Maps(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var maps []map[string]string
		data := g.Slice{
			g.Map{
				"k1": "v1",
				"k2": "v2",
			},
			g.Map{
				"k3": "v3",
				"k4": "v4",
			},
		}
		err := gbconv.Scan(data, &maps)
		t.AssertNil(err)
		t.Assert(data, maps)
	})
	// json string parameter.
	gbtest.C(t, func(t *gbtest.T) {
		var maps []map[string]string
		data := `[{"k1":"v1","k2":"v2"},{"k3":"v3","k4":"v4"}]`
		err := gbconv.Scan(data, &maps)
		t.AssertNil(err)
		t.Assert(maps, g.Slice{
			g.Map{
				"k1": "v1",
				"k2": "v2",
			},
			g.Map{
				"k3": "v3",
				"k4": "v4",
			},
		})
	})
}

func Test_Scan_JsonAttributes(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Sku struct {
			GiftId      int64  `json:"gift_id"`
			Name        string `json:"name"`
			ScorePrice  int    `json:"score_price"`
			MarketPrice int    `json:"market_price"`
			CostPrice   int    `json:"cost_price"`
			Stock       int    `json:"stock"`
		}
		v := gbvar.New(`
[
{"name": "red", "stock": 10, "gift_id": 1, "cost_price": 80, "score_price": 188, "market_price": 188}, 
{"name": "blue", "stock": 100, "gift_id": 2, "cost_price": 81, "score_price": 200, "market_price": 288}
]`)
		type Product struct {
			Skus []Sku
		}
		var p *Product
		err := gbconv.Scan(g.Map{
			"Skus": v,
		}, &p)
		t.AssertNil(err)
		t.Assert(len(p.Skus), 2)

		t.Assert(p.Skus[0].Name, "red")
		t.Assert(p.Skus[0].Stock, 10)
		t.Assert(p.Skus[0].GiftId, 1)
		t.Assert(p.Skus[0].CostPrice, 80)
		t.Assert(p.Skus[0].ScorePrice, 188)
		t.Assert(p.Skus[0].MarketPrice, 188)

		t.Assert(p.Skus[1].Name, "blue")
		t.Assert(p.Skus[1].Stock, 100)
		t.Assert(p.Skus[1].GiftId, 2)
		t.Assert(p.Skus[1].CostPrice, 81)
		t.Assert(p.Skus[1].ScorePrice, 200)
		t.Assert(p.Skus[1].MarketPrice, 288)
	})
}

func Test_Scan_JsonAttributes_StringArray(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type S struct {
			Array []string
		}
		var s *S
		err := gbconv.Scan(g.Map{
			"Array": `["a", "b"]`,
		}, &s)
		t.AssertNil(err)
		t.Assert(len(s.Array), 2)
		t.Assert(s.Array[0], "a")
		t.Assert(s.Array[1], "b")
	})

	gbtest.C(t, func(t *gbtest.T) {
		type S struct {
			Array []string
		}
		var s *S
		err := gbconv.Scan(g.Map{
			"Array": `[]`,
		}, &s)
		t.AssertNil(err)
		t.Assert(len(s.Array), 0)
	})

	gbtest.C(t, func(t *gbtest.T) {
		type S struct {
			Array []int64
		}
		var s *S
		err := gbconv.Scan(g.Map{
			"Array": `[]`,
		}, &s)
		t.AssertNil(err)
		t.Assert(len(s.Array), 0)
	})
}

func Test_Scan_SameType_Just_Assign(t *testing.T) {
	// Struct.
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Uid     int
			Name    string
			Pass1   string
			Pass2   string
			Pointer *int
		}
		var (
			int1  = 1
			int2  = 1
			user1 = new(User)
			user2 *User
		)
		user1.Pointer = &int1
		err := gbconv.Scan(user1, &user2)
		t.AssertNil(err)
		t.Assert(fmt.Sprintf(`%p`, user1), fmt.Sprintf(`%p`, user2))
		t.Assert(*user1.Pointer, *user2.Pointer)
		user1.Pointer = &int2
		t.Assert(*user1.Pointer, *user2.Pointer)
	})
	// Map.
	gbtest.C(t, func(t *gbtest.T) {
		var (
			int1 = 1
			int2 = 1
			m1   = map[string]*int{
				"int": &int1,
			}
			m2 map[string]*int
		)
		err := gbconv.Scan(m1, &m2)
		t.AssertNil(err)
		t.Assert(fmt.Sprintf(`%p`, m1), fmt.Sprintf(`%p`, m2))
		t.Assert(*m1["int"], *m2["int"])
		m1["int"] = &int2
		t.Assert(*m1["int"], *m2["int"])
	})
}

func Test_ScanList_Basic(t *testing.T) {
	// Struct attribute.
	gbtest.C(t, func(t *gbtest.T) {
		type EntityUser struct {
			Uid  int
			Name string
		}

		type EntityUserDetail struct {
			Uid     int
			Address string
		}

		type EntityUserScores struct {
			Id    int
			Uid   int
			Score int
		}

		type Entity struct {
			User       EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}

		var (
			err         error
			entities    []Entity
			entityUsers = []EntityUser{
				{Uid: 1, Name: "name1"},
				{Uid: 2, Name: "name2"},
				{Uid: 3, Name: "name3"},
			}
			userDetails = []EntityUserDetail{
				{Uid: 1, Address: "address1"},
				{Uid: 2, Address: "address2"},
			}
			userScores = []EntityUserScores{
				{Id: 10, Uid: 1, Score: 100},
				{Id: 11, Uid: 1, Score: 60},
				{Id: 20, Uid: 2, Score: 99},
			}
		)
		err = gbconv.ScanList(entityUsers, &entities, "User")
		t.AssertNil(err)

		err = gbconv.ScanList(userDetails, &entities, "UserDetail", "User", "uid")
		t.AssertNil(err)

		err = gbconv.ScanList(userScores, &entities, "UserScores", "User", "uid")
		t.AssertNil(err)

		t.Assert(len(entities), 3)
		t.Assert(entities[0].User, entityUsers[0])
		t.Assert(entities[1].User, entityUsers[1])
		t.Assert(entities[2].User, entityUsers[2])

		t.Assert(entities[0].UserDetail, userDetails[0])
		t.Assert(entities[1].UserDetail, userDetails[1])
		t.Assert(entities[2].UserDetail, EntityUserDetail{})

		t.Assert(len(entities[0].UserScores), 2)
		t.Assert(entities[0].UserScores[0], userScores[0])
		t.Assert(entities[0].UserScores[1], userScores[1])

		t.Assert(len(entities[1].UserScores), 1)
		t.Assert(entities[1].UserScores[0], userScores[2])

		t.Assert(len(entities[2].UserScores), 0)
	})
	// Pointer attribute.
	gbtest.C(t, func(t *gbtest.T) {
		type EntityUser struct {
			Uid  int
			Name string
		}

		type EntityUserDetail struct {
			Uid     int
			Address string
		}

		type EntityUserScores struct {
			Id    int
			Uid   int
			Score int
		}

		type Entity struct {
			User       *EntityUser
			UserDetail *EntityUserDetail
			UserScores []*EntityUserScores
		}

		var (
			err         error
			entities    []*Entity
			entityUsers = []*EntityUser{
				{Uid: 1, Name: "name1"},
				{Uid: 2, Name: "name2"},
				{Uid: 3, Name: "name3"},
			}
			userDetails = []*EntityUserDetail{
				{Uid: 1, Address: "address1"},
				{Uid: 2, Address: "address2"},
			}
			userScores = []*EntityUserScores{
				{Id: 10, Uid: 1, Score: 100},
				{Id: 11, Uid: 1, Score: 60},
				{Id: 20, Uid: 2, Score: 99},
			}
		)
		err = gbconv.ScanList(entityUsers, &entities, "User")
		t.AssertNil(err)

		err = gbconv.ScanList(userDetails, &entities, "UserDetail", "User", "uid")
		t.AssertNil(err)

		err = gbconv.ScanList(userScores, &entities, "UserScores", "User", "uid")
		t.AssertNil(err)

		t.Assert(len(entities), 3)
		t.Assert(entities[0].User, entityUsers[0])
		t.Assert(entities[1].User, entityUsers[1])
		t.Assert(entities[2].User, entityUsers[2])

		t.Assert(entities[0].UserDetail, userDetails[0])
		t.Assert(entities[1].UserDetail, userDetails[1])
		t.Assert(entities[2].UserDetail, nil)

		t.Assert(len(entities[0].UserScores), 2)
		t.Assert(entities[0].UserScores[0], userScores[0])
		t.Assert(entities[0].UserScores[1], userScores[1])

		t.Assert(len(entities[1].UserScores), 1)
		t.Assert(entities[1].UserScores[0], userScores[2])

		t.Assert(len(entities[2].UserScores), 0)
	})
}

func Test_ScanList_Embedded(t *testing.T) {
	// Struct attribute.
	gbtest.C(t, func(t *gbtest.T) {
		type EntityUser struct {
			Uid  int
			Name string
		}

		type EntityUserDetail struct {
			Uid     int
			Address string
		}

		type EntityUserScores struct {
			Id    int
			Uid   int
			Score int
		}

		type Entity struct {
			EntityUser
			UserDetail EntityUserDetail
			UserScores []EntityUserScores
		}

		var (
			err         error
			entities    []Entity
			entityUsers = []EntityUser{
				{Uid: 1, Name: "name1"},
				{Uid: 2, Name: "name2"},
				{Uid: 3, Name: "name3"},
			}
			userDetails = []EntityUserDetail{
				{Uid: 1, Address: "address1"},
				{Uid: 2, Address: "address2"},
			}
			userScores = []EntityUserScores{
				{Id: 10, Uid: 1, Score: 100},
				{Id: 11, Uid: 1, Score: 60},
				{Id: 20, Uid: 2, Score: 99},
			}
		)
		err = gbconv.Scan(entityUsers, &entities)
		t.AssertNil(err)

		err = gbconv.ScanList(userDetails, &entities, "UserDetail", "uid")
		t.AssertNil(err)

		err = gbconv.ScanList(userScores, &entities, "UserScores", "uid")
		t.AssertNil(err)

		t.Assert(len(entities), 3)
		t.Assert(entities[0].EntityUser, entityUsers[0])
		t.Assert(entities[1].EntityUser, entityUsers[1])
		t.Assert(entities[2].EntityUser, entityUsers[2])

		t.Assert(entities[0].UserDetail, userDetails[0])
		t.Assert(entities[1].UserDetail, userDetails[1])
		t.Assert(entities[2].UserDetail, EntityUserDetail{})

		t.Assert(len(entities[0].UserScores), 2)
		t.Assert(entities[0].UserScores[0], userScores[0])
		t.Assert(entities[0].UserScores[1], userScores[1])

		t.Assert(len(entities[1].UserScores), 1)
		t.Assert(entities[1].UserScores[0], userScores[2])

		t.Assert(len(entities[2].UserScores), 0)
	})
	// Pointer attribute.
	gbtest.C(t, func(t *gbtest.T) {
		type EntityUser struct {
			Uid  int
			Name string
		}

		type EntityUserDetail struct {
			Uid     int
			Address string
		}

		type EntityUserScores struct {
			Id    int
			Uid   int
			Score int
		}

		type Entity struct {
			*EntityUser
			UserDetail *EntityUserDetail
			UserScores []*EntityUserScores
		}

		var (
			err         error
			entities    []Entity
			entityUsers = []EntityUser{
				{Uid: 1, Name: "name1"},
				{Uid: 2, Name: "name2"},
				{Uid: 3, Name: "name3"},
			}
			userDetails = []EntityUserDetail{
				{Uid: 1, Address: "address1"},
				{Uid: 2, Address: "address2"},
			}
			userScores = []EntityUserScores{
				{Id: 10, Uid: 1, Score: 100},
				{Id: 11, Uid: 1, Score: 60},
				{Id: 20, Uid: 2, Score: 99},
			}
		)
		err = gbconv.Scan(entityUsers, &entities)
		t.AssertNil(err)

		err = gbconv.ScanList(userDetails, &entities, "UserDetail", "uid")
		t.AssertNil(err)

		err = gbconv.ScanList(userScores, &entities, "UserScores", "uid")
		t.AssertNil(err)

		t.Assert(len(entities), 3)
		t.Assert(entities[0].EntityUser, entityUsers[0])
		t.Assert(entities[1].EntityUser, entityUsers[1])
		t.Assert(entities[2].EntityUser, entityUsers[2])

		t.Assert(entities[0].UserDetail, userDetails[0])
		t.Assert(entities[1].UserDetail, userDetails[1])
		t.Assert(entities[2].UserDetail, nil)

		t.Assert(len(entities[0].UserScores), 2)
		t.Assert(entities[0].UserScores[0], userScores[0])
		t.Assert(entities[0].UserScores[1], userScores[1])

		t.Assert(len(entities[1].UserScores), 1)
		t.Assert(entities[1].UserScores[0], userScores[2])

		t.Assert(len(entities[2].UserScores), 0)
	})
}

type Float64 float64

func (f *Float64) UnmarshalValue(value interface{}) error {
	if v, ok := value.(*big.Rat); ok {
		f64, _ := v.Float64()
		*f = Float64(f64)
	}
	return nil
}

func Test_Issue1607(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Demo struct {
			B Float64
		}
		rat := &big.Rat{}
		rat.SetFloat64(1.5)

		var demos = make([]Demo, 1)
		err := gbconv.Scan([]map[string]interface{}{
			{"A": 1, "B": rat},
		}, &demos)
		t.AssertNil(err)
		t.Assert(demos[0].B, 1.5)
	})
}
