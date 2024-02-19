package gbconv_test

import (
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
)

func Test_Structs_WithTag(t *testing.T) {
	type User struct {
		Uid      int    `json:"id"`
		NickName string `json:"name"`
	}
	gbtest.C(t, func(t *gbtest.T) {
		var users []User
		params := g.Slice{
			g.Map{
				"id":   1,
				"name": "name1",
			},
			g.Map{
				"id":   2,
				"name": "name2",
			},
		}
		err := gbconv.Structs(params, &users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Uid, 1)
		t.Assert(users[0].NickName, "name1")
		t.Assert(users[1].Uid, 2)
		t.Assert(users[1].NickName, "name2")
	})
	gbtest.C(t, func(t *gbtest.T) {
		var users []*User
		params := g.Slice{
			g.Map{
				"id":   1,
				"name": "name1",
			},
			g.Map{
				"id":   2,
				"name": "name2",
			},
		}
		err := gbconv.Structs(params, &users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Uid, 1)
		t.Assert(users[0].NickName, "name1")
		t.Assert(users[1].Uid, 2)
		t.Assert(users[1].NickName, "name2")
	})
}

func Test_Structs_WithoutTag(t *testing.T) {
	type User struct {
		Uid      int
		NickName string
	}
	gbtest.C(t, func(t *gbtest.T) {
		var users []User
		params := g.Slice{
			g.Map{
				"uid":       1,
				"nick-name": "name1",
			},
			g.Map{
				"uid":       2,
				"nick-name": "name2",
			},
		}
		err := gbconv.Structs(params, &users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Uid, 1)
		t.Assert(users[0].NickName, "name1")
		t.Assert(users[1].Uid, 2)
		t.Assert(users[1].NickName, "name2")
	})
	gbtest.C(t, func(t *gbtest.T) {
		var users []*User
		params := g.Slice{
			g.Map{
				"uid":       1,
				"nick-name": "name1",
			},
			g.Map{
				"uid":       2,
				"nick-name": "name2",
			},
		}
		err := gbconv.Structs(params, &users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Uid, 1)
		t.Assert(users[0].NickName, "name1")
		t.Assert(users[1].Uid, 2)
		t.Assert(users[1].NickName, "name2")
	})
}

func Test_Structs_SliceParameter(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Uid      int
			NickName string
		}
		var users []User
		params := g.Slice{
			g.Map{
				"uid":       1,
				"nick-name": "name1",
			},
			g.Map{
				"uid":       2,
				"nick-name": "name2",
			},
		}
		err := gbconv.Structs(params, users)
		t.AssertNE(err, nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Uid      int
			NickName string
		}
		type A struct {
			Users []User
		}
		var a A
		params := g.Slice{
			g.Map{
				"uid":       1,
				"nick-name": "name1",
			},
			g.Map{
				"uid":       2,
				"nick-name": "name2",
			},
		}
		err := gbconv.Structs(params, a.Users)
		t.AssertNE(err, nil)
	})
}

func Test_Structs_DirectReflectSet(t *testing.T) {
	type A struct {
		Id   int
		Name string
	}
	gbtest.C(t, func(t *gbtest.T) {
		var (
			a = []*A{
				{Id: 1, Name: "john"},
				{Id: 2, Name: "smith"},
			}
			b []*A
		)
		err := gbconv.Structs(a, &b)
		t.AssertNil(err)
		t.AssertEQ(a, b)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var (
			a = []A{
				{Id: 1, Name: "john"},
				{Id: 2, Name: "smith"},
			}
			b []A
		)
		err := gbconv.Structs(a, &b)
		t.AssertNil(err)
		t.AssertEQ(a, b)
	})
}

func Test_Structs_IntSliceAttribute(t *testing.T) {
	type A struct {
		Id []int
	}
	type B struct {
		*A
		Name string
	}
	gbtest.C(t, func(t *gbtest.T) {
		var (
			array []*B
		)
		err := gbconv.Structs(g.Slice{
			g.Map{"id": nil, "name": "john"},
			g.Map{"id": nil, "name": "smith"},
		}, &array)
		t.AssertNil(err)
		t.Assert(len(array), 2)
		t.Assert(array[0].Name, "john")
		t.Assert(array[1].Name, "smith")
	})
}
