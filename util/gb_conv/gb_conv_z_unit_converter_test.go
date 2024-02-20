package gbconv_test

import (
	"encoding/json"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
	"time"
)

func TestConverter_ConvertWithRefer(t *testing.T) {
	type tA struct {
		Val int
	}

	type tB struct {
		Val1 int32
		Val2 string
	}

	gbtest.C(t, func(t *gbtest.T) {
		err := gbconv.RegisterConverter(func(a tA) (b *tB, err error) {
			b = &tB{
				Val1: int32(a.Val),
				Val2: "abcd",
			}
			return
		})
		t.AssertNil(err)
	})

	gbtest.C(t, func(t *gbtest.T) {
		a := &tA{
			Val: 1,
		}
		var b tB
		result := gbconv.ConvertWithRefer(a, &b)
		t.Assert(result.(*tB), &tB{
			Val1: 1,
			Val2: "abcd",
		})
	})

	gbtest.C(t, func(t *gbtest.T) {
		a := &tA{
			Val: 1,
		}
		var b tB
		result := gbconv.ConvertWithRefer(a, b)
		t.Assert(result.(tB), tB{
			Val1: 1,
			Val2: "abcd",
		})
	})
}

func TestConverter_Struct(t *testing.T) {
	type tA struct {
		Val int
	}

	type tB struct {
		Val1 int32
		Val2 string
	}

	type tAA struct {
		ValTop int
		ValTA  tA
	}

	type tBB struct {
		ValTop int32
		ValTB  tB
	}

	type tCC struct {
		ValTop string
		ValTa  *tB
	}

	type tDD struct {
		ValTop string
		ValTa  tB
	}

	type tEE struct {
		Val1 time.Time  `json:"val1"`
		Val2 *time.Time `json:"val2"`
		Val3 *time.Time `json:"val3"`
	}

	type tFF struct {
		Val1 json.RawMessage            `json:"val1"`
		Val2 []json.RawMessage          `json:"val2"`
		Val3 map[string]json.RawMessage `json:"val3"`
	}

	gbtest.C(t, func(t *gbtest.T) {
		a := &tA{
			Val: 1,
		}
		var b *tB
		err := gbconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.Val1, 0)
		t.Assert(b.Val2, "")
	})

	gbtest.C(t, func(t *gbtest.T) {
		err := gbconv.RegisterConverter(func(a tA) (b *tB, err error) {
			b = &tB{
				Val1: int32(a.Val),
				Val2: "abc",
			}
			return
		})
		t.AssertNil(err)
	})

	gbtest.C(t, func(t *gbtest.T) {
		a := &tA{
			Val: 1,
		}
		var b *tB
		err := gbconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.Val1, 1)
		t.Assert(b.Val2, "abc")
	})

	gbtest.C(t, func(t *gbtest.T) {
		a := &tA{
			Val: 1,
		}
		var b *tB
		err := gbconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.Val1, 1)
		t.Assert(b.Val2, "abc")
	})

	gbtest.C(t, func(t *gbtest.T) {
		a := &tA{
			Val: 1,
		}
		var b *tB
		err := gbconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.Val1, 1)
		t.Assert(b.Val2, "abc")
	})

	gbtest.C(t, func(t *gbtest.T) {
		a := &tA{
			Val: 1,
		}
		var b *tB
		err := gbconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.Val1, 1)
		t.Assert(b.Val2, "abc")
	})

	gbtest.C(t, func(t *gbtest.T) {
		aa := &tAA{
			ValTop: 123,
			ValTA:  tA{Val: 234},
		}
		var bb *tBB

		err := gbconv.Scan(aa, &bb)
		t.AssertNil(err)
		t.AssertNE(bb, nil)
		t.Assert(bb.ValTop, 123)
		t.AssertNE(bb.ValTB.Val1, 234)

		err = gbconv.RegisterConverter(func(a tAA) (b *tBB, err error) {
			b = &tBB{
				ValTop: int32(a.ValTop) + 2,
			}
			err = gbconv.Scan(a.ValTA, &b.ValTB)
			return
		})
		t.AssertNil(err)

		err = gbconv.Scan(aa, &bb)
		t.AssertNil(err)
		t.AssertNE(bb, nil)
		t.Assert(bb.ValTop, 125)
		t.Assert(bb.ValTB.Val1, 234)
		t.Assert(bb.ValTB.Val2, "abc")

	})

	gbtest.C(t, func(t *gbtest.T) {
		aa := &tAA{
			ValTop: 123,
			ValTA:  tA{Val: 234},
		}
		var cc *tCC
		err := gbconv.Scan(aa, &cc)
		t.AssertNil(err)
		t.AssertNE(cc, nil)
		t.Assert(cc.ValTop, "123")
		t.AssertNE(cc.ValTa, nil)
		t.Assert(cc.ValTa.Val1, 234)
		t.Assert(cc.ValTa.Val2, "abc")
	})

	gbtest.C(t, func(t *gbtest.T) {
		aa := &tAA{
			ValTop: 123,
			ValTA:  tA{Val: 234},
		}

		var dd *tDD
		err := gbconv.Scan(aa, &dd)
		t.AssertNil(err)
		t.AssertNE(dd, nil)
		t.Assert(dd.ValTop, "123")
		t.Assert(dd.ValTa.Val1, 234)
		t.Assert(dd.ValTa.Val2, "abc")
	})

	gbtest.C(t, func(t *gbtest.T) {
		aa := &tEE{}

		var tmp = map[string]any{
			"val1": "2023-04-15 19:10:00 +0800 CST",
			"val2": "2023-04-15 19:10:00 +0800 CST",
			"val3": "2006-01-02T15:04:05Z07:00",
		}
		err := gbconv.Struct(tmp, aa)
		t.AssertNil(err)
		t.AssertNE(aa, nil)
		t.Assert(aa.Val1.Local(), gbtime.New("2023-04-15 19:10:00 +0800 CST").Local().Time)
		t.Assert(aa.Val2.Local(), gbtime.New("2023-04-15 19:10:00 +0800 CST").Local().Time)
		t.Assert(aa.Val3.Local(), gbtime.New("2006-01-02T15:04:05Z07:00").Local().Time)
	})

	gbtest.C(t, func(t *gbtest.T) {
		ff := &tFF{}
		var tmp = map[string]any{
			"val1": map[string]any{"hello": "world"},
			"val2": []any{map[string]string{"hello": "world"}},
			"val3": map[string]map[string]string{"val3": {"hello": "world"}},
		}

		err := gbconv.Struct(tmp, ff)
		t.AssertNil(err)
		t.AssertNE(ff, nil)
		t.Assert(ff.Val1, []byte(`{"hello":"world"}`))
		t.AssertEQ(len(ff.Val2), 1)
		t.Assert(ff.Val2[0], []byte(`{"hello":"world"}`))
		t.AssertEQ(len(ff.Val3), 1)
		t.Assert(ff.Val3["val3"], []byte(`{"hello":"world"}`))
	})
}

func TestConverter_CustomBasicType_ToStruct(t *testing.T) {
	type CustomString string
	type CustomStruct struct {
		S string
	}
	gbtest.C(t, func(t *gbtest.T) {
		var (
			a CustomString = "abc"
			b *CustomStruct
		)
		err := gbconv.Scan(a, &b)
		t.AssertNE(err, nil)
		t.Assert(b, nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		err := gbconv.RegisterConverter(func(a CustomString) (b *CustomStruct, err error) {
			b = &CustomStruct{
				S: string(a),
			}
			return
		})
		t.AssertNil(err)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var (
			a CustomString = "abc"
			b *CustomStruct
		)
		err := gbconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.S, a)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var (
			a CustomString = "abc"
			b *CustomStruct
		)
		err := gbconv.Scan(&a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.S, a)
	})
}

func TestConverter_CustomTimeType_ToStruct(t *testing.T) {
	type timestamppb struct {
		S string
	}
	type CustomGTime struct {
		T *gbtime.Time
	}
	type CustomPbTime struct {
		T *timestamppb
	}
	gbtest.C(t, func(t *gbtest.T) {
		var (
			a = CustomGTime{
				T: gbtime.NewFromStrFormat("2023-10-26", "Y-m-d"),
			}
			b *CustomPbTime
		)
		err := gbconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.Assert(b.T.S, "")
	})

	gbtest.C(t, func(t *gbtest.T) {
		err := gbconv.RegisterConverter(func(in gbtime.Time) (*timestamppb, error) {
			return &timestamppb{
				S: in.Local().Format("Y-m-d"),
			}, nil
		})
		t.AssertNil(err)
		err = gbconv.RegisterConverter(func(in timestamppb) (*gbtime.Time, error) {
			return gbtime.NewFromStr(in.S), nil
		})
		t.AssertNil(err)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var (
			a = CustomGTime{
				T: gbtime.NewFromStrFormat("2023-10-26", "Y-m-d"),
			}
			b *CustomPbTime
			c *CustomGTime
		)
		err := gbconv.Scan(a, &b)
		t.AssertNil(err)
		t.AssertNE(b, nil)
		t.AssertNE(b.T, nil)

		err = gbconv.Scan(b, &c)
		t.AssertNil(err)
		t.AssertNE(c, nil)
		t.AssertNE(c.T, nil)
		t.AssertEQ(a.T.Timestamp(), c.T.Timestamp())
	})
}
