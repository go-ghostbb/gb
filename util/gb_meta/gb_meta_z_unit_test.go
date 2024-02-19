package gbmeta_test

import (
	"ghostbb.io/gb/internal/json"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbmeta "ghostbb.io/gb/util/gb_meta"
	"testing"
)

func TestMeta_Basic(t *testing.T) {
	type A struct {
		gbmeta.Meta `tag:"123" orm:"456"`
		Id          int
		Name        string
	}

	gbtest.C(t, func(t *gbtest.T) {
		a := &A{
			Id:   100,
			Name: "john",
		}
		t.Assert(len(gbmeta.Data(a)), 2)
		t.AssertEQ(gbmeta.Get(a, "tag").String(), "123")
		t.AssertEQ(gbmeta.Get(a, "orm").String(), "456")
		t.AssertEQ(gbmeta.Get(a, "none"), nil)

		b, err := json.Marshal(a)
		t.AssertNil(err)
		t.Assert(b, `{"Id":100,"Name":"john"}`)
	})
}

func TestMeta_Convert_Map(t *testing.T) {
	type A struct {
		gbmeta.Meta `tag:"123" orm:"456"`
		Id          int
		Name        string
	}

	gbtest.C(t, func(t *gbtest.T) {
		a := &A{
			Id:   100,
			Name: "john",
		}
		m := gbconv.Map(a)
		t.Assert(len(m), 2)
		t.Assert(m[`Meta`], nil)
	})
}

func TestMeta_Json(t *testing.T) {
	type A struct {
		gbmeta.Meta `tag:"123" orm:"456"`
		Id          int
	}

	gbtest.C(t, func(t *gbtest.T) {
		a := &A{
			Id: 100,
		}
		b, err := json.Marshal(a)
		t.AssertNil(err)
		t.Assert(string(b), `{"Id":100}`)
	})
}
