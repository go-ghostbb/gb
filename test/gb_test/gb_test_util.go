package gbtest

import (
	"fmt"
	gbdebug "github.com/Ghostbb-io/gb/debug/gb_debug"
	"github.com/Ghostbb-io/gb/internal/empty"
	gbconv "github.com/Ghostbb-io/gb/util/gb_conv"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const (
	pathFilterKey = "/test/gtest/gtest"
)

// C creates a unit testing case.
// The parameter `t` is the pointer to testing.T of stdlib (*testing.T).
// The parameter `f` is the closure function for unit testing case.
func C(t *testing.T, f func(t *T)) {
	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n%s", err, gbdebug.StackWithFilter([]string{pathFilterKey}))
			t.Fail()
		}
	}()
	f(&T{t})
}

// Assert checks `value` and `expect` EQUAL.
func Assert(value, expect interface{}) {
	rvExpect := reflect.ValueOf(expect)
	if empty.IsNil(value) {
		value = nil
	}
	if rvExpect.Kind() == reflect.Map {
		if err := compareMap(value, expect); err != nil {
			panic(err)
		}
		return
	}
	var (
		strValue  = gbconv.String(value)
		strExpect = gbconv.String(expect)
	)
	if strValue != strExpect {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v == %v`, strValue, strExpect))
	}
}

// AssertEQ checks `value` and `expect` EQUAL, including their TYPES.
func AssertEQ(value, expect interface{}) {
	// Value assert.
	rvExpect := reflect.ValueOf(expect)
	if empty.IsNil(value) {
		value = nil
	}
	if rvExpect.Kind() == reflect.Map {
		if err := compareMap(value, expect); err != nil {
			panic(err)
		}
		return
	}
	strValue := gbconv.String(value)
	strExpect := gbconv.String(expect)
	if strValue != strExpect {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v == %v`, strValue, strExpect))
	}
	// Type assert.
	t1 := reflect.TypeOf(value)
	t2 := reflect.TypeOf(expect)
	if t1 != t2 {
		panic(fmt.Sprintf(`[ASSERT] EXPECT TYPE %v[%v] == %v[%v]`, strValue, t1, strExpect, t2))
	}
}

// AssertNE checks `value` and `expect` NOT EQUAL.
func AssertNE(value, expect interface{}) {
	rvExpect := reflect.ValueOf(expect)
	if empty.IsNil(value) {
		value = nil
	}
	if rvExpect.Kind() == reflect.Map {
		if err := compareMap(value, expect); err == nil {
			panic(fmt.Sprintf(`[ASSERT] EXPECT %v != %v`, value, expect))
		}
		return
	}
	var (
		strValue  = gbconv.String(value)
		strExpect = gbconv.String(expect)
	)
	if strValue == strExpect {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v != %v`, strValue, strExpect))
	}
}

// AssertNQ checks `value` and `expect` NOT EQUAL, including their TYPES.
func AssertNQ(value, expect interface{}) {
	// Type assert.
	t1 := reflect.TypeOf(value)
	t2 := reflect.TypeOf(expect)
	if t1 == t2 {
		panic(
			fmt.Sprintf(
				`[ASSERT] EXPECT TYPE %v[%v] != %v[%v]`,
				gbconv.String(value), t1, gbconv.String(expect), t2,
			),
		)
	}
	// Value assert.
	AssertNE(value, expect)
}

// AssertGT checks `value` is GREATER THAN `expect`.
// Notice that, only string, integer and float types can be compared by AssertGT,
// others are invalid.
func AssertGT(value, expect interface{}) {
	passed := false
	switch reflect.ValueOf(expect).Kind() {
	case reflect.String:
		passed = gbconv.String(value) > gbconv.String(expect)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		passed = gbconv.Int(value) > gbconv.Int(expect)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		passed = gbconv.Uint(value) > gbconv.Uint(expect)

	case reflect.Float32, reflect.Float64:
		passed = gbconv.Float64(value) > gbconv.Float64(expect)
	}
	if !passed {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v > %v`, value, expect))
	}
}

// AssertGE checks `value` is GREATER OR EQUAL THAN `expect`.
// Notice that, only string, integer and float types can be compared by AssertGTE,
// others are invalid.
func AssertGE(value, expect interface{}) {
	passed := false
	switch reflect.ValueOf(expect).Kind() {
	case reflect.String:
		passed = gbconv.String(value) >= gbconv.String(expect)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		passed = gbconv.Int64(value) >= gbconv.Int64(expect)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		passed = gbconv.Uint64(value) >= gbconv.Uint64(expect)

	case reflect.Float32, reflect.Float64:
		passed = gbconv.Float64(value) >= gbconv.Float64(expect)
	}
	if !passed {
		panic(fmt.Sprintf(
			`[ASSERT] EXPECT %v(%v) >= %v(%v)`,
			value, reflect.ValueOf(value).Kind(),
			expect, reflect.ValueOf(expect).Kind(),
		))
	}
}

// AssertLT checks `value` is LESS EQUAL THAN `expect`.
// Notice that, only string, integer and float types can be compared by AssertLT,
// others are invalid.
func AssertLT(value, expect interface{}) {
	passed := false
	switch reflect.ValueOf(expect).Kind() {
	case reflect.String:
		passed = gbconv.String(value) < gbconv.String(expect)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		passed = gbconv.Int(value) < gbconv.Int(expect)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		passed = gbconv.Uint(value) < gbconv.Uint(expect)

	case reflect.Float32, reflect.Float64:
		passed = gbconv.Float64(value) < gbconv.Float64(expect)
	}
	if !passed {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v < %v`, value, expect))
	}
}

// AssertLE checks `value` is LESS OR EQUAL THAN `expect`.
// Notice that, only string, integer and float types can be compared by AssertLTE,
// others are invalid.
func AssertLE(value, expect interface{}) {
	passed := false
	switch reflect.ValueOf(expect).Kind() {
	case reflect.String:
		passed = gbconv.String(value) <= gbconv.String(expect)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		passed = gbconv.Int(value) <= gbconv.Int(expect)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		passed = gbconv.Uint(value) <= gbconv.Uint(expect)

	case reflect.Float32, reflect.Float64:
		passed = gbconv.Float64(value) <= gbconv.Float64(expect)
	}
	if !passed {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v <= %v`, value, expect))
	}
}

// AssertIN checks `value` is IN `expect`.
// The `expect` should be a slice,
// but the `value` can be a slice or a basic type variable.
// TODO map support.
// TODO: gbconv.Strings(0) is not [0]
func AssertIN(value, expect interface{}) {
	var (
		passed     = true
		expectKind = reflect.ValueOf(expect).Kind()
	)
	switch expectKind {
	case reflect.Slice, reflect.Array:
		expectSlice := gbconv.Strings(expect)
		for _, v1 := range gbconv.Strings(value) {
			result := false
			for _, v2 := range expectSlice {
				if v1 == v2 {
					result = true
					break
				}
			}
			if !result {
				passed = false
				break
			}
		}
	default:
		panic(fmt.Sprintf(`[ASSERT] INVALID EXPECT VALUE TYPE: %v`, expectKind))
	}
	if !passed {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v IN %v`, value, expect))
	}
}

// AssertNI checks `value` is NOT IN `expect`.
// The `expect` should be a slice,
// but the `value` can be a slice or a basic type variable.
// TODO map support.
func AssertNI(value, expect interface{}) {
	var (
		passed     = true
		expectKind = reflect.ValueOf(expect).Kind()
	)
	switch expectKind {
	case reflect.Slice, reflect.Array:
		for _, v1 := range gbconv.Strings(value) {
			result := true
			for _, v2 := range gbconv.Strings(expect) {
				if v1 == v2 {
					result = false
					break
				}
			}
			if !result {
				passed = false
				break
			}
		}
	default:
		panic(fmt.Sprintf(`[ASSERT] INVALID EXPECT VALUE TYPE: %v`, expectKind))
	}
	if !passed {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v NOT IN %v`, value, expect))
	}
}

// Error panics with given `message`.
func Error(message ...interface{}) {
	panic(fmt.Sprintf("[ERROR] %s", fmt.Sprint(message...)))
}

// Fatal prints `message` to stderr and exit the process.
func Fatal(message ...interface{}) {
	_, _ = fmt.Fprintf(
		os.Stderr, "[FATAL] %s\n%s", fmt.Sprint(message...),
		gbdebug.StackWithFilter([]string{pathFilterKey}),
	)
	os.Exit(1)
}

// compareMap compares two maps, returns nil if they are equal, or else returns error.
func compareMap(value, expect interface{}) error {
	var (
		rvValue  = reflect.ValueOf(value)
		rvExpect = reflect.ValueOf(expect)
	)
	if rvExpect.Kind() == reflect.Map {
		if rvValue.Kind() == reflect.Map {
			if rvExpect.Len() == rvValue.Len() {
				// Turn two interface maps to the same type for comparison.
				// Direct use of rvValue.MapIndex(key).Interface() will panic
				// when the key types are inconsistent.
				mValue := make(map[string]string)
				mExpect := make(map[string]string)
				ksValue := rvValue.MapKeys()
				ksExpect := rvExpect.MapKeys()
				for _, key := range ksValue {
					mValue[gbconv.String(key.Interface())] = gbconv.String(rvValue.MapIndex(key).Interface())
				}
				for _, key := range ksExpect {
					mExpect[gbconv.String(key.Interface())] = gbconv.String(rvExpect.MapIndex(key).Interface())
				}
				for k, v := range mExpect {
					if v != mValue[k] {
						return fmt.Errorf(`[ASSERT] EXPECT VALUE map["%v"]:%v == map["%v"]:%v`+
							"\nGIVEN : %v\nEXPECT: %v", k, mValue[k], k, v, mValue, mExpect)
					}
				}
			} else {
				return fmt.Errorf(`[ASSERT] EXPECT MAP LENGTH %d == %d`, rvValue.Len(), rvExpect.Len())
			}
		} else {
			return fmt.Errorf(`[ASSERT] EXPECT VALUE TO BE A MAP, BUT GIVEN "%s"`, rvValue.Kind())
		}
	}
	return nil
}

// AssertNil asserts `value` is nil.
func AssertNil(value interface{}) {
	if empty.IsNil(value) {
		return
	}
	if err, ok := value.(error); ok {
		panic(fmt.Sprintf(`%+v`, err))
	}
	Assert(value, nil)
}

// DataPath retrieves and returns the testdata path of current package,
// which is used for unit testing cases only.
// The optional parameter `names` specifies the sub-folders/sub-files,
// which will be joined with current system separator and returned with the path.
func DataPath(names ...string) string {
	_, path, _ := gbdebug.CallerWithFilter([]string{pathFilterKey})
	path = filepath.Dir(path) + string(filepath.Separator) + "testdata"
	for _, name := range names {
		path += string(filepath.Separator) + name
	}
	return path
}

// DataContent retrieves and returns the file content for specified testdata path of current package
func DataContent(names ...string) string {
	path := DataPath(names...)
	if path != "" {
		data, err := os.ReadFile(path)
		if err == nil {
			return string(data)
		}
	}
	return ""
}
