package gbconv

import (
	"github.com/Ghostbb-io/gb/internal/json"
	"github.com/Ghostbb-io/gb/internal/reflection"
	"reflect"
)

// SliceAny is alias of Interfaces.
func SliceAny(any interface{}) []interface{} {
	return Interfaces(any)
}

// Interfaces converts `any` to []interface{}.
func Interfaces(any interface{}) []interface{} {
	if any == nil {
		return nil
	}
	var array []interface{}
	switch value := any.(type) {
	case []interface{}:
		array = value
	case []string:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []int:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []int8:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []int16:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []int32:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []int64:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []uint:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []uint8:
		if json.Valid(value) {
			_ = json.UnmarshalUseNumber(value, &array)
		} else {
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		}
	case []uint16:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []uint32:
		for _, v := range value {
			array = append(array, v)
		}
	case []uint64:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []bool:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []float32:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	case []float64:
		array = make([]interface{}, len(value))
		for k, v := range value {
			array[k] = v
		}
	}
	if array != nil {
		return array
	}
	if v, ok := any.(iInterfaces); ok {
		return v.Interfaces()
	}
	// JSON format string value converting.
	if checkJsonAndUnmarshalUseNumber(any, &array) {
		return array
	}
	// Not a common type, it then uses reflection for conversion.
	originValueAndKind := reflection.OriginValueAndKind(any)
	switch originValueAndKind.OriginKind {
	case reflect.Slice, reflect.Array:
		var (
			length = originValueAndKind.OriginValue.Len()
			slice  = make([]interface{}, length)
		)
		for i := 0; i < length; i++ {
			slice[i] = originValueAndKind.OriginValue.Index(i).Interface()
		}
		return slice

	default:
		return []interface{}{any}
	}
}
