package gbstructs_test

import (
	"ghostbb.io/gb/frame/g"
	gbstructs "ghostbb.io/gb/os/gb_structs"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func Test_Basic(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user User
		m, _ := gbstructs.TagMapName(user, []string{"params"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})
		m, _ = gbstructs.TagMapName(&user, []string{"params"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})

		m, _ = gbstructs.TagMapName(&user, []string{"params", "my-tag1"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})
		m, _ = gbstructs.TagMapName(&user, []string{"my-tag1", "params"})
		t.Assert(m, g.Map{"name": "Name", "pass1": "Pass"})
		m, _ = gbstructs.TagMapName(&user, []string{"my-tag2", "params"})
		t.Assert(m, g.Map{"name": "Name", "pass2": "Pass"})
	})

	gbtest.C(t, func(t *gbtest.T) {
		type Base struct {
			Pass1 string `params:"password1"`
			Pass2 string `params:"password2"`
		}
		type UserWithBase struct {
			Id   int
			Name string
			Base `params:"base"`
		}
		user := new(UserWithBase)
		m, _ := gbstructs.TagMapName(user, []string{"params"})
		t.Assert(m, g.Map{
			"base":      "Base",
			"password1": "Pass1",
			"password2": "Pass2",
		})
	})

	gbtest.C(t, func(t *gbtest.T) {
		type Base struct {
			Pass1 string `params:"password1"`
			Pass2 string `params:"password2"`
		}
		type UserWithEmbeddedAttribute struct {
			Id   int
			Name string
			Base
		}
		type UserWithoutEmbeddedAttribute struct {
			Id   int
			Name string
			Pass Base
		}
		user1 := new(UserWithEmbeddedAttribute)
		user2 := new(UserWithoutEmbeddedAttribute)
		m, _ := gbstructs.TagMapName(user1, []string{"params"})
		t.Assert(m, g.Map{"password1": "Pass1", "password2": "Pass2"})
		m, _ = gbstructs.TagMapName(user2, []string{"params"})
		t.Assert(m, g.Map{})
	})
}

func Test_StructOfNilPointer(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user *User
		m, _ := gbstructs.TagMapName(user, []string{"params"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})
		m, _ = gbstructs.TagMapName(&user, []string{"params"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})

		m, _ = gbstructs.TagMapName(&user, []string{"params", "my-tag1"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})
		m, _ = gbstructs.TagMapName(&user, []string{"my-tag1", "params"})
		t.Assert(m, g.Map{"name": "Name", "pass1": "Pass"})
		m, _ = gbstructs.TagMapName(&user, []string{"my-tag2", "params"})
		t.Assert(m, g.Map{"name": "Name", "pass2": "Pass"})
	})
}

func Test_Fields(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user *User
		fields, _ := gbstructs.Fields(gbstructs.FieldsInput{
			Pointer:         user,
			RecursiveOption: 0,
		})
		t.Assert(len(fields), 3)
		t.Assert(fields[0].Name(), "Id")
		t.Assert(fields[1].Name(), "Name")
		t.Assert(fields[1].Tag("params"), "name")
		t.Assert(fields[2].Name(), "Pass")
		t.Assert(fields[2].Tag("my-tag1"), "pass1")
		t.Assert(fields[2].Tag("my-tag2"), "pass2")
		t.Assert(fields[2].Tag("params"), "pass")
	})
}

func Test_Fields_WithEmbedded1(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type B struct {
			Name string
			Age  int
		}
		type A struct {
			Site  string
			B     // Should be put here to validate its index.
			Score int64
		}
		r, err := gbstructs.Fields(gbstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: gbstructs.RecursiveOptionEmbeddedNoTag,
		})
		t.AssertNil(err)
		t.Assert(len(r), 4)
		t.Assert(r[0].Name(), `Site`)
		t.Assert(r[1].Name(), `Name`)
		t.Assert(r[2].Name(), `Age`)
		t.Assert(r[3].Name(), `Score`)
	})
}

func Test_Fields_WithEmbedded2(t *testing.T) {
	type MetaNode struct {
		Id          uint   `orm:"id,primary"  description:""`
		Capacity    string `orm:"capacity"    description:"Capacity string"`
		Allocatable string `orm:"allocatable" description:"Allocatable string"`
		Status      string `orm:"status"      description:"Status string"`
	}
	type MetaNodeZone struct {
		Nodes    uint
		Clusters uint
		Disk     uint
		Cpu      uint
		Memory   uint
		Zone     string
	}

	type MetaNodeItem struct {
		MetaNode
		Capacity    []MetaNodeZone `dc:"Capacity []MetaNodeZone"`
		Allocatable []MetaNodeZone `dc:"Allocatable []MetaNodeZone"`
	}

	gbtest.C(t, func(t *gbtest.T) {
		r, err := gbstructs.Fields(gbstructs.FieldsInput{
			Pointer:         new(MetaNodeItem),
			RecursiveOption: gbstructs.RecursiveOptionEmbeddedNoTag,
		})
		t.AssertNil(err)
		t.Assert(len(r), 4)
		t.Assert(r[0].Name(), `Id`)
		t.Assert(r[1].Name(), `Capacity`)
		t.Assert(r[1].TagStr(), `dc:"Capacity []MetaNodeZone"`)
		t.Assert(r[2].Name(), `Allocatable`)
		t.Assert(r[2].TagStr(), `dc:"Allocatable []MetaNodeZone"`)
		t.Assert(r[3].Name(), `Status`)
	})
}

// Filter repeated fields when there is embedded struct.
func Test_Fields_WithEmbedded_Filter(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type B struct {
			Name string
			Age  int
		}
		type A struct {
			Name  string
			Site  string
			Age   string
			B     // Should be put here to validate its index.
			Score int64
		}
		r, err := gbstructs.Fields(gbstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: gbstructs.RecursiveOptionEmbeddedNoTag,
		})
		t.AssertNil(err)
		t.Assert(len(r), 4)
		t.Assert(r[0].Name(), `Name`)
		t.Assert(r[1].Name(), `Site`)
		t.Assert(r[2].Name(), `Age`)
		t.Assert(r[3].Name(), `Score`)
	})
}

func Test_FieldMap(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user *User
		m, _ := gbstructs.FieldMap(gbstructs.FieldMapInput{
			Pointer:          user,
			PriorityTagArray: []string{"params"},
			RecursiveOption:  gbstructs.RecursiveOptionEmbedded,
		})
		t.Assert(len(m), 3)
		_, ok := m["Id"]
		t.Assert(ok, true)
		_, ok = m["Name"]
		t.Assert(ok, false)
		_, ok = m["name"]
		t.Assert(ok, true)
		_, ok = m["Pass"]
		t.Assert(ok, false)
		_, ok = m["pass"]
		t.Assert(ok, true)
	})
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user *User
		m, _ := gbstructs.FieldMap(gbstructs.FieldMapInput{
			Pointer:          user,
			PriorityTagArray: nil,
			RecursiveOption:  gbstructs.RecursiveOptionEmbedded,
		})
		t.Assert(len(m), 3)
		_, ok := m["Id"]
		t.Assert(ok, true)
		_, ok = m["Name"]
		t.Assert(ok, true)
		_, ok = m["name"]
		t.Assert(ok, false)
		_, ok = m["Pass"]
		t.Assert(ok, true)
		_, ok = m["pass"]
		t.Assert(ok, false)
	})
}

func Test_StructType(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			B
		}
		r, err := gbstructs.StructType(new(A))
		t.AssertNil(err)
		t.Assert(r.Signature(), `ghostbb.io/gb/os/gb_structs_test/gbstructs_test.A`)
	})
	gbtest.C(t, func(t *gbtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			B
		}
		r, err := gbstructs.StructType(new(A).B)
		t.AssertNil(err)
		t.Assert(r.Signature(), `ghostbb.io/gb/os/gb_structs_test/gbstructs_test.B`)
	})
	gbtest.C(t, func(t *gbtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			*B
		}
		r, err := gbstructs.StructType(new(A).B)
		t.AssertNil(err)
		t.Assert(r.String(), `gbstructs_test.B`)
	})
	// Error.
	gbtest.C(t, func(t *gbtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			*B
			Id int
		}
		_, err := gbstructs.StructType(new(A).Id)
		t.AssertNE(err, nil)
	})
}

func Test_StructTypeBySlice(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			Array []*B
		}
		r, err := gbstructs.StructType(new(A).Array)
		t.AssertNil(err)
		t.Assert(r.Signature(), `ghostbb.io/gb/os/gb_structs_test/gbstructs_test.B`)
	})
	gbtest.C(t, func(t *gbtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			Array []B
		}
		r, err := gbstructs.StructType(new(A).Array)
		t.AssertNil(err)
		t.Assert(r.Signature(), `ghostbb.io/gb/os/gb_structs_test/gbstructs_test.B`)
	})
	gbtest.C(t, func(t *gbtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			Array *[]B
		}
		r, err := gbstructs.StructType(new(A).Array)
		t.AssertNil(err)
		t.Assert(r.Signature(), `ghostbb.io/gb/os/gb_structs_test/gbstructs_test.B`)
	})
}

func TestType_FieldKeys(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type B struct {
			Id   int
			Name string
		}
		type A struct {
			Array []*B
		}
		r, err := gbstructs.StructType(new(A).Array)
		t.AssertNil(err)
		t.Assert(r.FieldKeys(), g.Slice{"Id", "Name"})
	})
}

func TestType_TagMap(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type A struct {
			Id   int    `d:"123" description:"I love gb"`
			Name string `v:"required" description:"應用Id"`
		}
		r, err := gbstructs.Fields(gbstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0].TagMap()["d"], `123`)
		t.Assert(r[0].TagMap()["description"], `I love gb`)
		t.Assert(r[1].TagMap()["v"], `required`)
		t.Assert(r[1].TagMap()["description"], `應用Id`)
	})
}

func TestType_TagJsonName(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type A struct {
			Name string `json:"name,omitempty"`
		}
		r, err := gbstructs.Fields(gbstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 1)
		t.Assert(r[0].TagJsonName(), `name`)
	})
}

func TestType_TagDefault(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type A struct {
			Name  string `default:"john"`
			Name2 string `d:"john"`
		}
		r, err := gbstructs.Fields(gbstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0].TagDefault(), `john`)
		t.Assert(r[1].TagDefault(), `john`)
	})
}

func TestType_TagParam(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type A struct {
			Name  string `param:"name"`
			Name2 string `p:"name"`
		}
		r, err := gbstructs.Fields(gbstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0].TagParam(), `name`)
		t.Assert(r[1].TagParam(), `name`)
	})
}

func TestType_TagValid(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type A struct {
			Name  string `valid:"required"`
			Name2 string `v:"required"`
		}
		r, err := gbstructs.Fields(gbstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0].TagValid(), `required`)
		t.Assert(r[1].TagValid(), `required`)
	})
}

func TestType_TagDescription(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type A struct {
			Name  string `description:"my name"`
			Name2 string `des:"my name"`
			Name3 string `dc:"my name"`
		}
		r, err := gbstructs.Fields(gbstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0].TagDescription(), `my name`)
		t.Assert(r[1].TagDescription(), `my name`)
		t.Assert(r[2].TagDescription(), `my name`)
	})
}

func TestType_TagSummary(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type A struct {
			Name  string `summary:"my name"`
			Name2 string `sum:"my name"`
			Name3 string `sm:"my name"`
		}
		r, err := gbstructs.Fields(gbstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0].TagSummary(), `my name`)
		t.Assert(r[1].TagSummary(), `my name`)
		t.Assert(r[2].TagSummary(), `my name`)
	})
}

func TestType_TagAdditional(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type A struct {
			Name  string `additional:"my name"`
			Name2 string `ad:"my name"`
		}
		r, err := gbstructs.Fields(gbstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0].TagAdditional(), `my name`)
		t.Assert(r[1].TagAdditional(), `my name`)
	})
}

func TestType_TagExample(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type A struct {
			Name  string `example:"john"`
			Name2 string `eg:"john"`
		}
		r, err := gbstructs.Fields(gbstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0].TagExample(), `john`)
		t.Assert(r[1].TagExample(), `john`)
	})
}
