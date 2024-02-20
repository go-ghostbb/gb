package gbconv_test

import (
	gbvar "ghostbb.io/gb/container/gb_var"
	"ghostbb.io/gb/frame/g"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"math"
	"testing"
	"time"
)

type iString interface {
	String() string
}

type S struct {
}

func (s S) String() string {
	return "22222"
}

type iError interface {
	Error() string
}

type S1 struct {
}

func (s1 S1) Error() string {
	return "22222"
}

func Test_Bool_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var any interface{} = nil
		t.AssertEQ(gbconv.Bool(any), false)
		t.AssertEQ(gbconv.Bool(false), false)
		t.AssertEQ(gbconv.Bool(nil), false)
		t.AssertEQ(gbconv.Bool(0), false)
		t.AssertEQ(gbconv.Bool("0"), false)
		t.AssertEQ(gbconv.Bool(""), false)
		t.AssertEQ(gbconv.Bool("false"), false)
		t.AssertEQ(gbconv.Bool("off"), false)
		t.AssertEQ(gbconv.Bool([]byte{}), false)
		t.AssertEQ(gbconv.Bool([]string{}), false)
		t.AssertEQ(gbconv.Bool([2]int{1, 2}), true)
		t.AssertEQ(gbconv.Bool([]interface{}{}), false)
		t.AssertEQ(gbconv.Bool([]map[int]int{}), false)

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value對,各個國家對應的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "羅馬"
		countryCapitalMap["Japan"] = "東京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gbconv.Bool(countryCapitalMap), true)

		t.AssertEQ(gbconv.Bool("1"), true)
		t.AssertEQ(gbconv.Bool("on"), true)
		t.AssertEQ(gbconv.Bool(1), true)
		t.AssertEQ(gbconv.Bool(123.456), true)
		t.AssertEQ(gbconv.Bool(boolStruct{}), true)
		t.AssertEQ(gbconv.Bool(&boolStruct{}), true)
	})
}

func Test_Int_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var any interface{} = nil
		t.AssertEQ(gbconv.Int(any), 0)
		t.AssertEQ(gbconv.Int(false), 0)
		t.AssertEQ(gbconv.Int(nil), 0)
		t.Assert(gbconv.Int(nil), 0)
		t.AssertEQ(gbconv.Int(0), 0)
		t.AssertEQ(gbconv.Int("0"), 0)
		t.AssertEQ(gbconv.Int(""), 0)
		t.AssertEQ(gbconv.Int("false"), 0)
		t.AssertEQ(gbconv.Int("off"), 0)
		t.AssertEQ(gbconv.Int([]byte{}), 0)
		t.AssertEQ(gbconv.Int([]string{}), 0)
		t.AssertEQ(gbconv.Int([2]int{1, 2}), 0)
		t.AssertEQ(gbconv.Int([]interface{}{}), 0)
		t.AssertEQ(gbconv.Int([]map[int]int{}), 0)

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value對,各個國家對應的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "羅馬"
		countryCapitalMap["Japan"] = "東京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gbconv.Int(countryCapitalMap), 0)

		t.AssertEQ(gbconv.Int("1"), 1)
		t.AssertEQ(gbconv.Int("on"), 0)
		t.AssertEQ(gbconv.Int(1), 1)
		t.AssertEQ(gbconv.Int(123.456), 123)
		t.AssertEQ(gbconv.Int(boolStruct{}), 0)
		t.AssertEQ(gbconv.Int(&boolStruct{}), 0)
		t.AssertEQ(gbconv.Int("NaN"), 0)
	})
}

func Test_Int8_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var any interface{} = nil
		t.Assert(gbconv.Int8(any), int8(0))
		t.AssertEQ(gbconv.Int8(false), int8(0))
		t.AssertEQ(gbconv.Int8(nil), int8(0))
		t.AssertEQ(gbconv.Int8(0), int8(0))
		t.AssertEQ(gbconv.Int8("0"), int8(0))
		t.AssertEQ(gbconv.Int8(""), int8(0))
		t.AssertEQ(gbconv.Int8("false"), int8(0))
		t.AssertEQ(gbconv.Int8("off"), int8(0))
		t.AssertEQ(gbconv.Int8([]byte{}), int8(0))
		t.AssertEQ(gbconv.Int8([]string{}), int8(0))
		t.AssertEQ(gbconv.Int8([2]int{1, 2}), int8(0))
		t.AssertEQ(gbconv.Int8([]interface{}{}), int8(0))
		t.AssertEQ(gbconv.Int8([]map[int]int{}), int8(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value對,各個國家對應的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "羅馬"
		countryCapitalMap["Japan"] = "東京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gbconv.Int8(countryCapitalMap), int8(0))

		t.AssertEQ(gbconv.Int8("1"), int8(1))
		t.AssertEQ(gbconv.Int8("on"), int8(0))
		t.AssertEQ(gbconv.Int8(int8(1)), int8(1))
		t.AssertEQ(gbconv.Int8(123.456), int8(123))
		t.AssertEQ(gbconv.Int8(boolStruct{}), int8(0))
		t.AssertEQ(gbconv.Int8(&boolStruct{}), int8(0))
		t.AssertEQ(gbconv.Int8("NaN"), int8(0))

	})
}

func Test_Int16_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var any interface{} = nil
		t.Assert(gbconv.Int16(any), int16(0))
		t.AssertEQ(gbconv.Int16(false), int16(0))
		t.AssertEQ(gbconv.Int16(nil), int16(0))
		t.AssertEQ(gbconv.Int16(0), int16(0))
		t.AssertEQ(gbconv.Int16("0"), int16(0))
		t.AssertEQ(gbconv.Int16(""), int16(0))
		t.AssertEQ(gbconv.Int16("false"), int16(0))
		t.AssertEQ(gbconv.Int16("off"), int16(0))
		t.AssertEQ(gbconv.Int16([]byte{}), int16(0))
		t.AssertEQ(gbconv.Int16([]string{}), int16(0))
		t.AssertEQ(gbconv.Int16([2]int{1, 2}), int16(0))
		t.AssertEQ(gbconv.Int16([]interface{}{}), int16(0))
		t.AssertEQ(gbconv.Int16([]map[int]int{}), int16(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value對,各個國家對應的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "羅馬"
		countryCapitalMap["Japan"] = "東京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gbconv.Int16(countryCapitalMap), int16(0))

		t.AssertEQ(gbconv.Int16("1"), int16(1))
		t.AssertEQ(gbconv.Int16("on"), int16(0))
		t.AssertEQ(gbconv.Int16(int16(1)), int16(1))
		t.AssertEQ(gbconv.Int16(123.456), int16(123))
		t.AssertEQ(gbconv.Int16(boolStruct{}), int16(0))
		t.AssertEQ(gbconv.Int16(&boolStruct{}), int16(0))
		t.AssertEQ(gbconv.Int16("NaN"), int16(0))
	})
}

func Test_Int32_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var any interface{} = nil
		t.Assert(gbconv.Int32(any), int32(0))
		t.AssertEQ(gbconv.Int32(false), int32(0))
		t.AssertEQ(gbconv.Int32(nil), int32(0))
		t.AssertEQ(gbconv.Int32(0), int32(0))
		t.AssertEQ(gbconv.Int32("0"), int32(0))
		t.AssertEQ(gbconv.Int32(""), int32(0))
		t.AssertEQ(gbconv.Int32("false"), int32(0))
		t.AssertEQ(gbconv.Int32("off"), int32(0))
		t.AssertEQ(gbconv.Int32([]byte{}), int32(0))
		t.AssertEQ(gbconv.Int32([]string{}), int32(0))
		t.AssertEQ(gbconv.Int32([2]int{1, 2}), int32(0))
		t.AssertEQ(gbconv.Int32([]interface{}{}), int32(0))
		t.AssertEQ(gbconv.Int32([]map[int]int{}), int32(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value對,各個國家對應的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "羅馬"
		countryCapitalMap["Japan"] = "東京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gbconv.Int32(countryCapitalMap), int32(0))

		t.AssertEQ(gbconv.Int32("1"), int32(1))
		t.AssertEQ(gbconv.Int32("on"), int32(0))
		t.AssertEQ(gbconv.Int32(int32(1)), int32(1))
		t.AssertEQ(gbconv.Int32(123.456), int32(123))
		t.AssertEQ(gbconv.Int32(boolStruct{}), int32(0))
		t.AssertEQ(gbconv.Int32(&boolStruct{}), int32(0))
		t.AssertEQ(gbconv.Int32("NaN"), int32(0))
	})
}

func Test_Int64_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var any interface{} = nil
		t.AssertEQ(gbconv.Int64("0x00e"), int64(14))
		t.Assert(gbconv.Int64("022"), int64(22))

		t.Assert(gbconv.Int64(any), int64(0))
		t.Assert(gbconv.Int64(true), 1)
		t.Assert(gbconv.Int64("1"), int64(1))
		t.Assert(gbconv.Int64("0"), int64(0))
		t.Assert(gbconv.Int64("X"), int64(0))
		t.Assert(gbconv.Int64("x"), int64(0))
		t.Assert(gbconv.Int64(int64(1)), int64(1))
		t.Assert(gbconv.Int64(int(0)), int64(0))
		t.Assert(gbconv.Int64(int8(0)), int64(0))
		t.Assert(gbconv.Int64(int16(0)), int64(0))
		t.Assert(gbconv.Int64(int32(0)), int64(0))
		t.Assert(gbconv.Int64(uint64(0)), int64(0))
		t.Assert(gbconv.Int64(uint32(0)), int64(0))
		t.Assert(gbconv.Int64(uint16(0)), int64(0))
		t.Assert(gbconv.Int64(uint8(0)), int64(0))
		t.Assert(gbconv.Int64(uint(0)), int64(0))
		t.Assert(gbconv.Int64(float32(0)), int64(0))

		t.AssertEQ(gbconv.Int64(false), int64(0))
		t.AssertEQ(gbconv.Int64(nil), int64(0))
		t.AssertEQ(gbconv.Int64(0), int64(0))
		t.AssertEQ(gbconv.Int64("0"), int64(0))
		t.AssertEQ(gbconv.Int64(""), int64(0))
		t.AssertEQ(gbconv.Int64("false"), int64(0))
		t.AssertEQ(gbconv.Int64("off"), int64(0))
		t.AssertEQ(gbconv.Int64([]byte{}), int64(0))
		t.AssertEQ(gbconv.Int64([]string{}), int64(0))
		t.AssertEQ(gbconv.Int64([2]int{1, 2}), int64(0))
		t.AssertEQ(gbconv.Int64([]interface{}{}), int64(0))
		t.AssertEQ(gbconv.Int64([]map[int]int{}), int64(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value對,各個國家對應的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "羅馬"
		countryCapitalMap["Japan"] = "東京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gbconv.Int64(countryCapitalMap), int64(0))

		t.AssertEQ(gbconv.Int64("1"), int64(1))
		t.AssertEQ(gbconv.Int64("on"), int64(0))
		t.AssertEQ(gbconv.Int64(int64(1)), int64(1))
		t.AssertEQ(gbconv.Int64(123.456), int64(123))
		t.AssertEQ(gbconv.Int64(boolStruct{}), int64(0))
		t.AssertEQ(gbconv.Int64(&boolStruct{}), int64(0))
		t.AssertEQ(gbconv.Int64("NaN"), int64(0))
	})
}

func Test_Uint_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var any interface{} = nil
		t.AssertEQ(gbconv.Uint(any), uint(0))
		t.AssertEQ(gbconv.Uint(false), uint(0))
		t.AssertEQ(gbconv.Uint(nil), uint(0))
		t.Assert(gbconv.Uint(nil), uint(0))
		t.AssertEQ(gbconv.Uint(uint(0)), uint(0))
		t.AssertEQ(gbconv.Uint("0"), uint(0))
		t.AssertEQ(gbconv.Uint(""), uint(0))
		t.AssertEQ(gbconv.Uint("false"), uint(0))
		t.AssertEQ(gbconv.Uint("off"), uint(0))
		t.AssertEQ(gbconv.Uint([]byte{}), uint(0))
		t.AssertEQ(gbconv.Uint([]string{}), uint(0))
		t.AssertEQ(gbconv.Uint([2]int{1, 2}), uint(0))
		t.AssertEQ(gbconv.Uint([]interface{}{}), uint(0))
		t.AssertEQ(gbconv.Uint([]map[int]int{}), uint(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value對,各個國家對應的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "羅馬"
		countryCapitalMap["Japan"] = "東京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gbconv.Uint(countryCapitalMap), uint(0))

		t.AssertEQ(gbconv.Uint("1"), uint(1))
		t.AssertEQ(gbconv.Uint("on"), uint(0))
		t.AssertEQ(gbconv.Uint(1), uint(1))
		t.AssertEQ(gbconv.Uint(123.456), uint(123))
		t.AssertEQ(gbconv.Uint(boolStruct{}), uint(0))
		t.AssertEQ(gbconv.Uint(&boolStruct{}), uint(0))
		t.AssertEQ(gbconv.Uint("NaN"), uint(0))
	})
}

func Test_Uint8_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var any interface{} = nil
		t.Assert(gbconv.Uint8(any), uint8(0))
		t.AssertEQ(gbconv.Uint8(uint8(1)), uint8(1))
		t.AssertEQ(gbconv.Uint8(false), uint8(0))
		t.AssertEQ(gbconv.Uint8(nil), uint8(0))
		t.AssertEQ(gbconv.Uint8(0), uint8(0))
		t.AssertEQ(gbconv.Uint8("0"), uint8(0))
		t.AssertEQ(gbconv.Uint8(""), uint8(0))
		t.AssertEQ(gbconv.Uint8("false"), uint8(0))
		t.AssertEQ(gbconv.Uint8("off"), uint8(0))
		t.AssertEQ(gbconv.Uint8([]byte{}), uint8(0))
		t.AssertEQ(gbconv.Uint8([]string{}), uint8(0))
		t.AssertEQ(gbconv.Uint8([2]int{1, 2}), uint8(0))
		t.AssertEQ(gbconv.Uint8([]interface{}{}), uint8(0))
		t.AssertEQ(gbconv.Uint8([]map[int]int{}), uint8(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value對,各個國家對應的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "羅馬"
		countryCapitalMap["Japan"] = "東京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gbconv.Uint8(countryCapitalMap), uint8(0))

		t.AssertEQ(gbconv.Uint8("1"), uint8(1))
		t.AssertEQ(gbconv.Uint8("on"), uint8(0))
		t.AssertEQ(gbconv.Uint8(int8(1)), uint8(1))
		t.AssertEQ(gbconv.Uint8(123.456), uint8(123))
		t.AssertEQ(gbconv.Uint8(boolStruct{}), uint8(0))
		t.AssertEQ(gbconv.Uint8(&boolStruct{}), uint8(0))
		t.AssertEQ(gbconv.Uint8("NaN"), uint8(0))
	})
}

func Test_Uint16_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var any interface{} = nil
		t.Assert(gbconv.Uint16(any), uint16(0))
		t.AssertEQ(gbconv.Uint16(uint16(1)), uint16(1))
		t.AssertEQ(gbconv.Uint16(false), uint16(0))
		t.AssertEQ(gbconv.Uint16(nil), uint16(0))
		t.AssertEQ(gbconv.Uint16(0), uint16(0))
		t.AssertEQ(gbconv.Uint16("0"), uint16(0))
		t.AssertEQ(gbconv.Uint16(""), uint16(0))
		t.AssertEQ(gbconv.Uint16("false"), uint16(0))
		t.AssertEQ(gbconv.Uint16("off"), uint16(0))
		t.AssertEQ(gbconv.Uint16([]byte{}), uint16(0))
		t.AssertEQ(gbconv.Uint16([]string{}), uint16(0))
		t.AssertEQ(gbconv.Uint16([2]int{1, 2}), uint16(0))
		t.AssertEQ(gbconv.Uint16([]interface{}{}), uint16(0))
		t.AssertEQ(gbconv.Uint16([]map[int]int{}), uint16(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value對,各個國家對應的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "羅馬"
		countryCapitalMap["Japan"] = "東京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gbconv.Uint16(countryCapitalMap), uint16(0))

		t.AssertEQ(gbconv.Uint16("1"), uint16(1))
		t.AssertEQ(gbconv.Uint16("on"), uint16(0))
		t.AssertEQ(gbconv.Uint16(int16(1)), uint16(1))
		t.AssertEQ(gbconv.Uint16(123.456), uint16(123))
		t.AssertEQ(gbconv.Uint16(boolStruct{}), uint16(0))
		t.AssertEQ(gbconv.Uint16(&boolStruct{}), uint16(0))
		t.AssertEQ(gbconv.Uint16("NaN"), uint16(0))
	})
}

func Test_Uint32_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var any interface{} = nil
		t.Assert(gbconv.Uint32(any), uint32(0))
		t.AssertEQ(gbconv.Uint32(uint32(1)), uint32(1))
		t.AssertEQ(gbconv.Uint32(false), uint32(0))
		t.AssertEQ(gbconv.Uint32(nil), uint32(0))
		t.AssertEQ(gbconv.Uint32(0), uint32(0))
		t.AssertEQ(gbconv.Uint32("0"), uint32(0))
		t.AssertEQ(gbconv.Uint32(""), uint32(0))
		t.AssertEQ(gbconv.Uint32("false"), uint32(0))
		t.AssertEQ(gbconv.Uint32("off"), uint32(0))
		t.AssertEQ(gbconv.Uint32([]byte{}), uint32(0))
		t.AssertEQ(gbconv.Uint32([]string{}), uint32(0))
		t.AssertEQ(gbconv.Uint32([2]int{1, 2}), uint32(0))
		t.AssertEQ(gbconv.Uint32([]interface{}{}), uint32(0))
		t.AssertEQ(gbconv.Uint32([]map[int]int{}), uint32(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value對,各個國家對應的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "羅馬"
		countryCapitalMap["Japan"] = "東京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gbconv.Uint32(countryCapitalMap), uint32(0))

		t.AssertEQ(gbconv.Uint32("1"), uint32(1))
		t.AssertEQ(gbconv.Uint32("on"), uint32(0))
		t.AssertEQ(gbconv.Uint32(int32(1)), uint32(1))
		t.AssertEQ(gbconv.Uint32(123.456), uint32(123))
		t.AssertEQ(gbconv.Uint32(boolStruct{}), uint32(0))
		t.AssertEQ(gbconv.Uint32(&boolStruct{}), uint32(0))
		t.AssertEQ(gbconv.Uint32("NaN"), uint32(0))
	})
}

func Test_Uint64_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var any interface{} = nil
		t.AssertEQ(gbconv.Uint64("0x00e"), uint64(14))
		t.Assert(gbconv.Uint64("022"), uint64(22))

		t.AssertEQ(gbconv.Uint64(any), uint64(0))
		t.AssertEQ(gbconv.Uint64(true), uint64(1))
		t.Assert(gbconv.Uint64("1"), int64(1))
		t.Assert(gbconv.Uint64("0"), uint64(0))
		t.Assert(gbconv.Uint64("X"), uint64(0))
		t.Assert(gbconv.Uint64("x"), uint64(0))
		t.Assert(gbconv.Uint64(int64(1)), uint64(1))
		t.Assert(gbconv.Uint64(int(0)), uint64(0))
		t.Assert(gbconv.Uint64(int8(0)), uint64(0))
		t.Assert(gbconv.Uint64(int16(0)), uint64(0))
		t.Assert(gbconv.Uint64(int32(0)), uint64(0))
		t.Assert(gbconv.Uint64(uint64(0)), uint64(0))
		t.Assert(gbconv.Uint64(uint32(0)), uint64(0))
		t.Assert(gbconv.Uint64(uint16(0)), uint64(0))
		t.Assert(gbconv.Uint64(uint8(0)), uint64(0))
		t.Assert(gbconv.Uint64(uint(0)), uint64(0))
		t.Assert(gbconv.Uint64(float32(0)), uint64(0))

		t.AssertEQ(gbconv.Uint64(false), uint64(0))
		t.AssertEQ(gbconv.Uint64(nil), uint64(0))
		t.AssertEQ(gbconv.Uint64(0), uint64(0))
		t.AssertEQ(gbconv.Uint64("0"), uint64(0))
		t.AssertEQ(gbconv.Uint64(""), uint64(0))
		t.AssertEQ(gbconv.Uint64("false"), uint64(0))
		t.AssertEQ(gbconv.Uint64("off"), uint64(0))
		t.AssertEQ(gbconv.Uint64([]byte{}), uint64(0))
		t.AssertEQ(gbconv.Uint64([]string{}), uint64(0))
		t.AssertEQ(gbconv.Uint64([2]int{1, 2}), uint64(0))
		t.AssertEQ(gbconv.Uint64([]interface{}{}), uint64(0))
		t.AssertEQ(gbconv.Uint64([]map[int]int{}), uint64(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value對,各個國家對應的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "羅馬"
		countryCapitalMap["Japan"] = "東京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gbconv.Uint64(countryCapitalMap), uint64(0))

		t.AssertEQ(gbconv.Uint64("1"), uint64(1))
		t.AssertEQ(gbconv.Uint64("on"), uint64(0))
		t.AssertEQ(gbconv.Uint64(int64(1)), uint64(1))
		t.AssertEQ(gbconv.Uint64(123.456), uint64(123))
		t.AssertEQ(gbconv.Uint64(boolStruct{}), uint64(0))
		t.AssertEQ(gbconv.Uint64(&boolStruct{}), uint64(0))
		t.AssertEQ(gbconv.Uint64("NaN"), uint64(0))
	})
}

func Test_Float32_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var any interface{} = nil
		t.Assert(gbconv.Float32(any), float32(0))
		t.AssertEQ(gbconv.Float32(false), float32(0))
		t.AssertEQ(gbconv.Float32(nil), float32(0))
		t.AssertEQ(gbconv.Float32(0), float32(0))
		t.AssertEQ(gbconv.Float32("0"), float32(0))
		t.AssertEQ(gbconv.Float32(""), float32(0))
		t.AssertEQ(gbconv.Float32("false"), float32(0))
		t.AssertEQ(gbconv.Float32("off"), float32(0))
		t.AssertEQ(gbconv.Float32([]byte{}), float32(0))
		t.AssertEQ(gbconv.Float32([]string{}), float32(0))
		t.AssertEQ(gbconv.Float32([2]int{1, 2}), float32(0))
		t.AssertEQ(gbconv.Float32([]interface{}{}), float32(0))
		t.AssertEQ(gbconv.Float32([]map[int]int{}), float32(0))
		t.AssertEQ(gbconv.Float32(gbvar.New(float32(0))), float32(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value對,各個國家對應的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "羅馬"
		countryCapitalMap["Japan"] = "東京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gbconv.Float32(countryCapitalMap), float32(0))

		t.AssertEQ(gbconv.Float32("1"), float32(1))
		t.AssertEQ(gbconv.Float32("on"), float32(0))
		t.AssertEQ(gbconv.Float32(float32(1)), float32(1))
		t.AssertEQ(gbconv.Float32(123.456), float32(123.456))
		t.AssertEQ(gbconv.Float32(boolStruct{}), float32(0))
		t.AssertEQ(gbconv.Float32(&boolStruct{}), float32(0))
		t.AssertEQ(gbconv.Float32("NaN"), float32(math.NaN()))
	})
}

func Test_Float64_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var any interface{} = nil
		t.Assert(gbconv.Float64(any), float64(0))
		t.AssertEQ(gbconv.Float64(false), float64(0))
		t.AssertEQ(gbconv.Float64(nil), float64(0))
		t.AssertEQ(gbconv.Float64(0), float64(0))
		t.AssertEQ(gbconv.Float64("0"), float64(0))
		t.AssertEQ(gbconv.Float64(""), float64(0))
		t.AssertEQ(gbconv.Float64("false"), float64(0))
		t.AssertEQ(gbconv.Float64("off"), float64(0))
		t.AssertEQ(gbconv.Float64([]byte{}), float64(0))
		t.AssertEQ(gbconv.Float64([]string{}), float64(0))
		t.AssertEQ(gbconv.Float64([2]int{1, 2}), float64(0))
		t.AssertEQ(gbconv.Float64([]interface{}{}), float64(0))
		t.AssertEQ(gbconv.Float64([]map[int]int{}), float64(0))
		t.AssertEQ(gbconv.Float64(gbvar.New(float64(0))), float64(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value對,各個國家對應的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "羅馬"
		countryCapitalMap["Japan"] = "東京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gbconv.Float64(countryCapitalMap), float64(0))

		t.AssertEQ(gbconv.Float64("1"), float64(1))
		t.AssertEQ(gbconv.Float64("on"), float64(0))
		t.AssertEQ(gbconv.Float64(float64(1)), float64(1))
		t.AssertEQ(gbconv.Float64(123.456), float64(123.456))
		t.AssertEQ(gbconv.Float64(boolStruct{}), float64(0))
		t.AssertEQ(gbconv.Float64(&boolStruct{}), float64(0))
		t.AssertEQ(gbconv.Float64("NaN"), float64(math.NaN()))
	})
}

func Test_String_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var s []rune
		t.AssertEQ(gbconv.String(s), "")
		var any interface{} = nil
		t.AssertEQ(gbconv.String(any), "")
		t.AssertEQ(gbconv.String("1"), "1")
		t.AssertEQ(gbconv.String("0"), string("0"))
		t.Assert(gbconv.String("X"), string("X"))
		t.Assert(gbconv.String("x"), string("x"))
		t.Assert(gbconv.String(int64(1)), uint64(1))
		t.Assert(gbconv.String(int(0)), string("0"))
		t.Assert(gbconv.String(int8(0)), string("0"))
		t.Assert(gbconv.String(int16(0)), string("0"))
		t.Assert(gbconv.String(int32(0)), string("0"))
		t.Assert(gbconv.String(uint64(0)), string("0"))
		t.Assert(gbconv.String(uint32(0)), string("0"))
		t.Assert(gbconv.String(uint16(0)), string("0"))
		t.Assert(gbconv.String(uint8(0)), string("0"))
		t.Assert(gbconv.String(uint(0)), string("0"))
		t.Assert(gbconv.String(float32(0)), string("0"))
		t.AssertEQ(gbconv.String(true), "true")
		t.AssertEQ(gbconv.String(false), "false")
		t.AssertEQ(gbconv.String(nil), "")
		t.AssertEQ(gbconv.String(0), string("0"))
		t.AssertEQ(gbconv.String("0"), string("0"))
		t.AssertEQ(gbconv.String(""), "")
		t.AssertEQ(gbconv.String("false"), "false")
		t.AssertEQ(gbconv.String("off"), string("off"))
		t.AssertEQ(gbconv.String([]byte{}), "")
		t.AssertEQ(gbconv.String([]string{}), "[]")
		t.AssertEQ(gbconv.String([2]int{1, 2}), "[1,2]")
		t.AssertEQ(gbconv.String([]interface{}{}), "[]")
		t.AssertEQ(gbconv.String(map[int]int{}), "{}")

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value對,各個國家對應的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "羅馬"
		countryCapitalMap["Japan"] = "東京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gbconv.String(countryCapitalMap), `{"France":"巴黎","India ":"新德里","Italy":"羅馬","Japan":"東京"}`)
		t.AssertEQ(gbconv.String(int64(1)), "1")
		t.AssertEQ(gbconv.String(123.456), "123.456")
		t.AssertEQ(gbconv.String(boolStruct{}), "{}")
		t.AssertEQ(gbconv.String(&boolStruct{}), "{}")

		var info = new(S)
		t.AssertEQ(gbconv.String(info), "22222")
		var errInfo = new(S1)
		t.AssertEQ(gbconv.String(errInfo), "22222")
	})
}

func Test_Runes_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Runes("www"), []int32{119, 119, 119})
		var s []rune
		t.AssertEQ(gbconv.Runes(s), nil)
	})
}

func Test_Rune_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Rune("www"), int32(0))
		t.AssertEQ(gbconv.Rune(int32(0)), int32(0))
		var s []rune
		t.AssertEQ(gbconv.Rune(s), int32(0))
	})
}

func Test_Bytes_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Bytes(nil), nil)
		t.AssertEQ(gbconv.Bytes(int32(0)), []uint8{0, 0, 0, 0})
		t.AssertEQ(gbconv.Bytes("s"), []uint8{115})
		t.AssertEQ(gbconv.Bytes([]byte("s")), []uint8{115})
		t.AssertEQ(gbconv.Bytes(gbvar.New([]byte("s"))), []uint8{115})
	})
}

func Test_Byte_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertEQ(gbconv.Byte(uint8(0)), uint8(0))
		t.AssertEQ(gbconv.Byte("s"), uint8(0))
		t.AssertEQ(gbconv.Byte([]byte("s")), uint8(115))
	})
}

func Test_Convert_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var any interface{} = nil
		t.AssertEQ(gbconv.Convert(any, "string"), "")
		t.AssertEQ(gbconv.Convert("1", "string"), "1")
		t.Assert(gbconv.Convert(int64(1), "int64"), int64(1))
		t.Assert(gbconv.Convert(int(0), "int"), int(0))
		t.Assert(gbconv.Convert(int8(0), "int8"), int8(0))
		t.Assert(gbconv.Convert(int16(0), "int16"), int16(0))
		t.Assert(gbconv.Convert(int32(0), "int32"), int32(0))
		t.Assert(gbconv.Convert(uint64(0), "uint64"), uint64(0))
		t.Assert(gbconv.Convert(uint32(0), "uint32"), uint32(0))
		t.Assert(gbconv.Convert(uint16(0), "uint16"), uint16(0))
		t.Assert(gbconv.Convert(uint8(0), "uint8"), uint8(0))
		t.Assert(gbconv.Convert(uint(0), "uint"), uint(0))
		t.Assert(gbconv.Convert(float32(0), "float32"), float32(0))
		t.Assert(gbconv.Convert(float64(0), "float64"), float64(0))
		t.AssertEQ(gbconv.Convert(true, "bool"), true)
		t.AssertEQ(gbconv.Convert([]byte{}, "[]byte"), []uint8{})
		t.AssertEQ(gbconv.Convert([]string{}, "[]string"), []string{})
		t.AssertEQ(gbconv.Convert([2]int{1, 2}, "[]int"), []int{1, 2})
		t.AssertEQ(gbconv.Convert([2]uint8{1, 2}, "[]uint8"), []uint8{1, 2})
		t.AssertEQ(gbconv.Convert("1989-01-02", "Time", "Y-m-d"), gbconv.Time("1989-01-02", "Y-m-d"))
		t.AssertEQ(gbconv.Convert(1989, "Time"), gbconv.Time("1970-01-01 08:33:09 +0800 CST"))
		t.AssertEQ(gbconv.Convert(gbtime.Now(), "gbtime.Time", 1), *gbtime.New())
		t.AssertEQ(gbconv.Convert(1989, "gbtime.Time"), *gbconv.GBTime("1970-01-01 08:33:09 +0800 CST"))
		t.AssertEQ(gbconv.Convert(gbtime.Now(), "*gbtime.Time", 1), gbtime.New())
		t.AssertEQ(gbconv.Convert(gbtime.Now(), "GBTime", 1), *gbtime.New())
		t.AssertEQ(gbconv.Convert(1989, "*gbtime.Time"), gbconv.GBTime(1989))
		t.AssertEQ(gbconv.Convert(1989, "Duration"), time.Duration(int64(1989)))
		t.AssertEQ(gbconv.Convert("1989", "Duration"), time.Duration(int64(1989)))
		t.AssertEQ(gbconv.Convert("1989", ""), "1989")

		var intNum int = 1
		t.Assert(gbconv.Convert(&intNum, "*int"), int(1))
		var int8Num int8 = 1
		t.Assert(gbconv.Convert(int8Num, "*int8"), int(1))
		t.Assert(gbconv.Convert(&int8Num, "*int8"), int(1))
		var int16Num int16 = 1
		t.Assert(gbconv.Convert(int16Num, "*int16"), int(1))
		t.Assert(gbconv.Convert(&int16Num, "*int16"), int(1))
		var int32Num int32 = 1
		t.Assert(gbconv.Convert(int32Num, "*int32"), int(1))
		t.Assert(gbconv.Convert(&int32Num, "*int32"), int(1))
		var int64Num int64 = 1
		t.Assert(gbconv.Convert(int64Num, "*int64"), int(1))
		t.Assert(gbconv.Convert(&int64Num, "*int64"), int(1))

		var uintNum uint = 1
		t.Assert(gbconv.Convert(&uintNum, "*uint"), int(1))
		var uint8Num uint8 = 1
		t.Assert(gbconv.Convert(uint8Num, "*uint8"), int(1))
		t.Assert(gbconv.Convert(&uint8Num, "*uint8"), int(1))
		var uint16Num uint16 = 1
		t.Assert(gbconv.Convert(uint16Num, "*uint16"), int(1))
		t.Assert(gbconv.Convert(&uint16Num, "*uint16"), int(1))
		var uint32Num uint32 = 1
		t.Assert(gbconv.Convert(uint32Num, "*uint32"), int(1))
		t.Assert(gbconv.Convert(&uint32Num, "*uint32"), int(1))
		var uint64Num uint64 = 1
		t.Assert(gbconv.Convert(uint64Num, "*uint64"), int(1))
		t.Assert(gbconv.Convert(&uint64Num, "*uint64"), int(1))

		var float32Num float32 = 1.1
		t.Assert(gbconv.Convert(float32Num, "*float32"), float32(1.1))
		t.Assert(gbconv.Convert(&float32Num, "*float32"), float32(1.1))

		var float64Num float64 = 1.1
		t.Assert(gbconv.Convert(float64Num, "*float64"), float64(1.1))
		t.Assert(gbconv.Convert(&float64Num, "*float64"), float64(1.1))

		var boolValue bool = true
		t.Assert(gbconv.Convert(boolValue, "*bool"), true)
		t.Assert(gbconv.Convert(&boolValue, "*bool"), true)

		var stringValue string = "1"
		t.Assert(gbconv.Convert(stringValue, "*string"), "1")
		t.Assert(gbconv.Convert(&stringValue, "*string"), "1")

		var durationValue time.Duration = 1989
		var expectDurationValue = time.Duration(int64(1989))
		t.AssertEQ(gbconv.Convert(&durationValue, "*time.Duration"), &expectDurationValue)
		t.AssertEQ(gbconv.Convert(durationValue, "*time.Duration"), &expectDurationValue)

		var string_interface_map = map[string]interface{}{"k1": 1}
		var string_int_map = map[string]int{"k1": 1}
		var string_string_map = map[string]string{"k1": "1"}
		t.AssertEQ(gbconv.Convert(string_int_map, "map[string]string"), string_string_map)
		t.AssertEQ(gbconv.Convert(string_int_map, "map[string]interface{}"), string_interface_map)
	})
}

func Test_Slice_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		value := 123.456
		t.AssertEQ(gbconv.Ints(value), []int{123})
		t.AssertEQ(gbconv.Ints(nil), nil)
		t.AssertEQ(gbconv.Ints([]string{"1", "2"}), []int{1, 2})
		t.AssertEQ(gbconv.Ints([]int{}), []int{})
		t.AssertEQ(gbconv.Ints([]int8{1, 2}), []int{1, 2})
		t.AssertEQ(gbconv.Ints([]int16{1, 2}), []int{1, 2})
		t.AssertEQ(gbconv.Ints([]int32{1, 2}), []int{1, 2})
		t.AssertEQ(gbconv.Ints([]int64{1, 2}), []int{1, 2})
		t.AssertEQ(gbconv.Ints([]uint{1}), []int{1})
		t.AssertEQ(gbconv.Ints([]uint8{1, 2}), []int{1, 2})
		t.AssertEQ(gbconv.Ints([]uint16{1, 2}), []int{1, 2})
		t.AssertEQ(gbconv.Ints([]uint32{1, 2}), []int{1, 2})
		t.AssertEQ(gbconv.Ints([]uint64{1, 2}), []int{1, 2})
		t.AssertEQ(gbconv.Ints([]bool{true}), []int{1})
		t.AssertEQ(gbconv.Ints([]float32{1, 2}), []int{1, 2})
		t.AssertEQ(gbconv.Ints([]float64{1, 2}), []int{1, 2})
		var inter []interface{} = make([]interface{}, 2)
		t.AssertEQ(gbconv.Ints(inter), []int{0, 0})

		t.AssertEQ(gbconv.Strings(value), []string{"123.456"})
		t.AssertEQ(gbconv.Strings(nil), nil)
		t.AssertEQ(gbconv.Strings([]string{"1", "2"}), []string{"1", "2"})
		t.AssertEQ(gbconv.Strings([]int{1}), []string{"1"})
		t.AssertEQ(gbconv.Strings([]int8{1, 2}), []string{"1", "2"})
		t.AssertEQ(gbconv.Strings([]int16{1, 2}), []string{"1", "2"})
		t.AssertEQ(gbconv.Strings([]int32{1, 2}), []string{"1", "2"})
		t.AssertEQ(gbconv.Strings([]int64{1, 2}), []string{"1", "2"})
		t.AssertEQ(gbconv.Strings([]uint{1}), []string{"1"})
		t.AssertEQ(gbconv.Strings([]uint8{1, 2}), []string{"1", "2"})
		t.AssertEQ(gbconv.Strings([]uint16{1, 2}), []string{"1", "2"})
		t.AssertEQ(gbconv.Strings([]uint32{1, 2}), []string{"1", "2"})
		t.AssertEQ(gbconv.Strings([]uint64{1, 2}), []string{"1", "2"})
		t.AssertEQ(gbconv.Strings([]bool{true}), []string{"true"})
		t.AssertEQ(gbconv.Strings([]float32{1, 2}), []string{"1", "2"})
		t.AssertEQ(gbconv.Strings([]float64{1, 2}), []string{"1", "2"})
		var strer = make([]interface{}, 2)
		t.AssertEQ(gbconv.Strings(strer), []string{"", ""})

		t.AssertEQ(gbconv.Floats(value), []float64{123.456})
		t.AssertEQ(gbconv.Floats(nil), nil)
		t.AssertEQ(gbconv.Floats([]string{"1", "2"}), []float64{1, 2})
		t.AssertEQ(gbconv.Floats([]int{1}), []float64{1})
		t.AssertEQ(gbconv.Floats([]int8{1, 2}), []float64{1, 2})
		t.AssertEQ(gbconv.Floats([]int16{1, 2}), []float64{1, 2})
		t.AssertEQ(gbconv.Floats([]int32{1, 2}), []float64{1, 2})
		t.AssertEQ(gbconv.Floats([]int64{1, 2}), []float64{1, 2})
		t.AssertEQ(gbconv.Floats([]uint{1}), []float64{1})
		t.AssertEQ(gbconv.Floats([]uint8{1, 2}), []float64{1, 2})
		t.AssertEQ(gbconv.Floats([]uint16{1, 2}), []float64{1, 2})
		t.AssertEQ(gbconv.Floats([]uint32{1, 2}), []float64{1, 2})
		t.AssertEQ(gbconv.Floats([]uint64{1, 2}), []float64{1, 2})
		t.AssertEQ(gbconv.Floats([]bool{true}), []float64{0})
		t.AssertEQ(gbconv.Floats([]float32{1, 2}), []float64{1, 2})
		t.AssertEQ(gbconv.Floats([]float64{1, 2}), []float64{1, 2})
		var floer = make([]interface{}, 2)
		t.AssertEQ(gbconv.Floats(floer), []float64{0, 0})

		t.AssertEQ(gbconv.Interfaces(value), []interface{}{123.456})
		t.AssertEQ(gbconv.Interfaces(nil), nil)
		t.AssertEQ(gbconv.Interfaces([]interface{}{1}), []interface{}{1})
		t.AssertEQ(gbconv.Interfaces([]string{"1"}), []interface{}{"1"})
		t.AssertEQ(gbconv.Interfaces([]int{1}), []interface{}{1})
		t.AssertEQ(gbconv.Interfaces([]int8{1}), []interface{}{1})
		t.AssertEQ(gbconv.Interfaces([]int16{1}), []interface{}{1})
		t.AssertEQ(gbconv.Interfaces([]int32{1}), []interface{}{1})
		t.AssertEQ(gbconv.Interfaces([]int64{1}), []interface{}{1})
		t.AssertEQ(gbconv.Interfaces([]uint{1}), []interface{}{1})
		t.AssertEQ(gbconv.Interfaces([]uint8{1}), []interface{}{1})
		t.AssertEQ(gbconv.Interfaces([]uint16{1}), []interface{}{1})
		t.AssertEQ(gbconv.Interfaces([]uint32{1}), []interface{}{1})
		t.AssertEQ(gbconv.Interfaces([]uint64{1}), []interface{}{1})
		t.AssertEQ(gbconv.Interfaces([]bool{true}), []interface{}{true})
		t.AssertEQ(gbconv.Interfaces([]float32{1}), []interface{}{1})
		t.AssertEQ(gbconv.Interfaces([]float64{1}), []interface{}{1})
		t.AssertEQ(gbconv.Interfaces([1]int{1}), []interface{}{1})

		type interSlice []int
		slices := interSlice{1}
		t.AssertEQ(gbconv.Interfaces(slices), []interface{}{1})

		t.AssertEQ(gbconv.Maps(nil), nil)
		t.AssertEQ(gbconv.Maps([]map[string]interface{}{{"a": "1"}}), []map[string]interface{}{{"a": "1"}})
		t.AssertEQ(gbconv.Maps(1223), []map[string]interface{}{nil})
		t.AssertEQ(gbconv.Maps([]int{}), nil)
	})
}

// 私有属性不会进行转换
func Test_Slice_PrivateAttribute_All(t *testing.T) {
	type User struct {
		Id   int           `json:"id"`
		name string        `json:"name"`
		Ad   []interface{} `json:"ad"`
	}
	gbtest.C(t, func(t *gbtest.T) {
		user := &User{1, "john", []interface{}{2}}
		array := gbconv.Interfaces(user)
		t.Assert(len(array), 1)
		t.Assert(array[0].(*User).Id, 1)
		t.Assert(array[0].(*User).name, "john")
		t.Assert(array[0].(*User).Ad, []interface{}{2})
	})
}

func Test_Map_Basic_All(t *testing.T) {
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
		t.AssertEQ(gbconv.Map(nil), nil)
		t.AssertEQ(gbconv.Map(map[string]interface{}{"a": 1}), map[string]interface{}{"a": 1})
		t.AssertEQ(gbconv.Map(map[int]interface{}{1: 1}), map[string]interface{}{"1": 1})
		t.AssertEQ(gbconv.Map(map[uint]interface{}{1: 1}), map[string]interface{}{"1": 1})
		t.AssertEQ(gbconv.Map(map[uint]string{1: "1"}), map[string]interface{}{"1": "1"})

		t.AssertEQ(gbconv.Map(map[interface{}]interface{}{"a": 1}), map[interface{}]interface{}{"a": 1})
		t.AssertEQ(gbconv.Map(map[interface{}]string{"a": "1"}), map[interface{}]string{"a": "1"})
		t.AssertEQ(gbconv.Map(map[interface{}]int{"a": 1}), map[interface{}]int{"a": 1})
		t.AssertEQ(gbconv.Map(map[interface{}]uint{"a": 1}), map[interface{}]uint{"a": 1})
		t.AssertEQ(gbconv.Map(map[interface{}]float32{"a": 1}), map[interface{}]float32{"a": 1})
		t.AssertEQ(gbconv.Map(map[interface{}]float64{"a": 1}), map[interface{}]float64{"a": 1})

		t.AssertEQ(gbconv.Map(map[string]bool{"a": true}), map[string]interface{}{"a": true})
		t.AssertEQ(gbconv.Map(map[string]int{"a": 1}), map[string]interface{}{"a": 1})
		t.AssertEQ(gbconv.Map(map[string]uint{"a": 1}), map[string]interface{}{"a": 1})
		t.AssertEQ(gbconv.Map(map[string]float32{"a": 1}), map[string]interface{}{"a": 1})
		t.AssertEQ(gbconv.Map(map[string]float64{"a": 1}), map[string]interface{}{"a": 1})

	})
}

func Test_Map_StructWithGconvTag_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Uid      int
			Name     string
			SiteUrl  string   `gbconv:"-"`
			NickName string   `gbconv:"nickname,omitempty"`
			Pass1    string   `gbconv:"password1"`
			Pass2    string   `gbconv:"password2"`
			Ss       []string `gbconv:"ss"`
		}
		user1 := User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
			Ss:      []string{"sss", "2222"},
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

func Test_Map_StructWithJsonTag_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Uid      int
			Name     string
			SiteUrl  string   `json:"-"`
			NickName string   `json:"nickname, omitempty"`
			Pass1    string   `json:"password1,newpassword"`
			Pass2    string   `json:"password2"`
			Ss       []string `json:"omitempty"`
			ssb, ssa string
		}
		user1 := User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
			Ss:      []string{"sss", "2222"},
			ssb:     "11",
			ssa:     "222",
		}
		user3 := User{
			Uid:      100,
			Name:     "john",
			NickName: "SSS",
			SiteUrl:  "https://goframe.org",
			Pass1:    "123",
			Pass2:    "456",
			Ss:       []string{"sss", "2222"},
			ssb:      "11",
			ssa:      "222",
		}
		user2 := &user1
		_ = gbconv.Map(user1, gbconv.MapOption{Tags: []string{"Ss"}})
		map1 := gbconv.Map(user1, gbconv.MapOption{Tags: []string{"json", "json2"}})
		map2 := gbconv.Map(user2)
		map3 := gbconv.Map(user3)
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
		t.Assert(map3["NickName"], nil)
	})
}

func Test_Map_PrivateAttribute_All(t *testing.T) {
	type User struct {
		Id   int
		name string
	}
	gbtest.C(t, func(t *gbtest.T) {
		user := &User{1, "john"}
		t.Assert(gbconv.Map(user), g.Map{"Id": 1})
	})
}

func Test_Map_StructInherit_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Ids struct {
			Id  int `json:"id"`
			Uid int `json:"uid"`
		}
		type Base struct {
			Ids
			CreateTime string `json:"create_time"`
		}
		type User struct {
			Base
			Passport string  `json:"passport"`
			Password string  `json:"password"`
			Nickname string  `json:"nickname"`
			S        *string `json:"nickname2"`
		}

		user := new(User)
		user.Id = 100
		user.Nickname = "john"
		user.CreateTime = "2019"
		var s = "s"
		user.S = &s

		m := gbconv.MapDeep(user)
		t.Assert(m["id"], user.Id)
		t.Assert(m["nickname"], user.Nickname)
		t.Assert(m["create_time"], user.CreateTime)
		t.Assert(m["nickname2"], user.S)
	})
}

func Test_Struct_Basic1_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Score struct {
			Name   int
			Result string
		}

		type Score2 struct {
			Name   int
			Result string
		}

		type User struct {
			Uid      int
			Name     string
			Site_Url string
			NickName string
			Pass1    string `gbconv:"password1"`
			Pass2    string `gbconv:"password2"`
			As       *Score
			Ass      Score
			Assb     []interface{}
		}
		// 使用默認映射规则绑定属性值到对象
		user := new(User)
		params1 := g.Map{
			"uid":       1,
			"Name":      "john",
			"siteurl":   "https://goframe.org",
			"nick_name": "johng",
			"PASS1":     "123",
			"PASS2":     "456",
			"As":        g.Map{"Name": 1, "Result": "22222"},
			"Ass":       &Score{11, "11"},
			"Assb":      []string{"wwww"},
		}
		_ = gbconv.Struct(nil, user)
		_ = gbconv.Struct(params1, nil)
		_ = gbconv.Struct([]interface{}{nil}, user)
		_ = gbconv.Struct(user, []interface{}{nil})

		var a = []interface{}{nil}
		ab := &a
		_ = gbconv.Struct(params1, *ab)
		var pi *int = nil
		_ = gbconv.Struct(params1, pi)

		_ = gbconv.Struct(params1, user)
		_ = gbconv.Struct(params1, user, map[string]string{"uid": "Names"})
		_ = gbconv.Struct(params1, user, map[string]string{"uid": "as"})

		// 使用struct tag映射绑定属性值到对象
		user = new(User)
		params2 := g.Map{
			"uid":       2,
			"name":      "smith",
			"site-url":  "https://goframe.org",
			"nick name": "johng",
			"password1": "111",
			"password2": "222",
		}
		if err := gbconv.Struct(params2, user); err != nil {
			gbtest.Error(err)
		}
		t.Assert(user, &User{
			Uid:      2,
			Name:     "smith",
			Site_Url: "https://goframe.org",
			NickName: "johng",
			Pass1:    "111",
			Pass2:    "222",
		})
	})
}

// 使用默認映射规则绑定属性值到对象
func Test_Struct_Basic2_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Uid     int
			Name    string
			SiteUrl string
			Pass1   string
			Pass2   string
		}
		user := new(User)
		params := g.Map{
			"uid":      1,
			"Name":     "john",
			"site_url": "https://goframe.org",
			"PASS1":    "123",
			"PASS2":    "456",
		}
		if err := gbconv.Struct(params, user); err != nil {
			gbtest.Error(err)
		}
		t.Assert(user, &User{
			Uid:     1,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
		})
	})
}

// 带有指针的基础类型属性
func Test_Struct_Basic3_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Uid  int
			Name *string
		}
		user := new(User)
		params := g.Map{
			"uid":  1,
			"Name": "john",
		}
		if err := gbconv.Struct(params, user); err != nil {
			gbtest.Error(err)
		}
		t.Assert(user.Uid, 1)
		t.Assert(*user.Name, "john")
	})
}

// slice类型属性的赋值
func Test_Struct_Attr_Slice_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Scores []int
		}
		scores := []interface{}{99, 100, 60, 140}
		user := new(User)
		if err := gbconv.Struct(g.Map{"Scores": scores}, user); err != nil {
			gbtest.Error(err)
		} else {
			t.Assert(user, &User{
				Scores: []int{99, 100, 60, 140},
			})
		}
	})
}

// 属性为struct对象
func Test_Struct_Attr_Struct_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Scores Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Scores": map[string]interface{}{
				"Name":   "john",
				"Result": 100,
			},
		}

		// 嵌套struct转换
		if err := gbconv.Struct(scores, user); err != nil {
			gbtest.Error(err)
		} else {
			t.Assert(user, &User{
				Scores: Score{
					Name:   "john",
					Result: 100,
				},
			})
		}
	})
}

// 属性为struct对象指针
func Test_Struct_Attr_Struct_Ptr_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Scores *Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Scores": map[string]interface{}{
				"Name":   "john",
				"Result": 100,
			},
		}

		// 嵌套struct转换
		if err := gbconv.Struct(scores, user); err != nil {
			gbtest.Error(err)
		} else {
			t.Assert(user.Scores, &Score{
				Name:   "john",
				Result: 100,
			})
		}
	})
}

// 属性为struct对象slice
func Test_Struct_Attr_Struct_Slice1_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Scores []Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Scores": map[string]interface{}{
				"Name":   "john",
				"Result": 100,
			},
		}

		// 嵌套struct转换，属性为slice类型，數值为map类型
		if err := gbconv.Struct(scores, user); err != nil {
			gbtest.Error(err)
		} else {
			t.Assert(user.Scores, []Score{
				{
					Name:   "john",
					Result: 100,
				},
			})
		}
	})
}

// 属性为struct对象slice
func Test_Struct_Attr_Struct_Slice2_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Scores []Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Scores": []interface{}{
				map[string]interface{}{
					"Name":   "john",
					"Result": 100,
				},
				map[string]interface{}{
					"Name":   "smith",
					"Result": 60,
				},
			},
		}

		// 嵌套struct转换，属性为slice类型，數值为slice map类型
		if err := gbconv.Struct(scores, user); err != nil {
			gbtest.Error(err)
		} else {
			t.Assert(user.Scores, []Score{
				{
					Name:   "john",
					Result: 100,
				},
				{
					Name:   "smith",
					Result: 60,
				},
			})
		}
	})
}

// 属性为struct对象slice ptr
func Test_Struct_Attr_Struct_Slice_Ptr_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Scores []*Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Scores": []interface{}{
				map[string]interface{}{
					"Name":   "john",
					"Result": 100,
				},
				map[string]interface{}{
					"Name":   "smith",
					"Result": 60,
				},
			},
		}

		// 嵌套struct转换，属性为slice类型，數值为slice map类型
		if err := gbconv.Struct(scores, user); err != nil {
			gbtest.Error(err)
		} else {
			t.Assert(len(user.Scores), 2)
			t.Assert(user.Scores[0], &Score{
				Name:   "john",
				Result: 100,
			})
			t.Assert(user.Scores[1], &Score{
				Name:   "smith",
				Result: 60,
			})
		}
	})
}

func Test_Struct_PrivateAttribute_All(t *testing.T) {
	type User struct {
		Id   int
		name string
	}
	gbtest.C(t, func(t *gbtest.T) {
		user := new(User)
		err := gbconv.Struct(g.Map{"id": 1, "name": "john"}, user)
		t.AssertNil(err)
		t.Assert(user.Id, 1)
		t.Assert(user.name, "")
	})
}

func Test_Struct_Embedded_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Ids struct {
			Id  int `json:"id"`
			Uid int `json:"uid"`
		}
		type Base struct {
			Ids
			CreateTime string `json:"create_time"`
		}
		type User struct {
			Base
			Passport string `json:"passport"`
			Password string `json:"password"`
			Nickname string `json:"nickname"`
		}
		data := g.Map{
			"id":          100,
			"uid":         101,
			"passport":    "t1",
			"password":    "123456",
			"nickname":    "T1",
			"create_time": "2019",
		}
		user := new(User)
		gbconv.Struct(data, user)
		t.Assert(user.Id, 100)
		t.Assert(user.Uid, 101)
		t.Assert(user.Nickname, "T1")
		t.Assert(user.CreateTime, "2019")
	})
}

func Test_Struct_Time_All(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			CreateTime time.Time
		}
		now := time.Now()
		user := new(User)
		gbconv.Struct(g.Map{
			"create_time": now,
		}, user)
		t.Assert(user.CreateTime.UTC().String(), now.UTC().String())
	})

	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			CreateTime *time.Time
		}
		now := time.Now()
		user := new(User)
		gbconv.Struct(g.Map{
			"create_time": &now,
		}, user)
		t.Assert(user.CreateTime.UTC().String(), now.UTC().String())
	})

	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			CreateTime *gbtime.Time
		}
		now := time.Now()
		user := new(User)
		gbconv.Struct(g.Map{
			"create_time": &now,
		}, user)
		t.Assert(user.CreateTime.Time.UTC().String(), now.UTC().String())
	})

	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			CreateTime gbtime.Time
		}
		now := time.Now()
		user := new(User)
		gbconv.Struct(g.Map{
			"create_time": &now,
		}, user)
		t.Assert(user.CreateTime.Time.UTC().String(), now.UTC().String())
	})

	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			CreateTime gbtime.Time
		}
		now := time.Now()
		user := new(User)
		gbconv.Struct(g.Map{
			"create_time": now,
		}, user)
		t.Assert(user.CreateTime.Time.UTC().String(), now.UTC().String())
	})
}
