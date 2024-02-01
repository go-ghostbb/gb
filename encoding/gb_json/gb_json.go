// Package gbjson provides convenient API for JSON/XML/INI/YAML/TOML data handling.
package gbjson

import (
	gbcode "github.com/Ghostbb-io/gb/errors/gb_code"
	gberror "github.com/Ghostbb-io/gb/errors/gb_error"
	"github.com/Ghostbb-io/gb/internal/reflection"
	"github.com/Ghostbb-io/gb/internal/rwmutex"
	"github.com/Ghostbb-io/gb/internal/utils"
	gbstr "github.com/Ghostbb-io/gb/text/gb_str"
	gbconv "github.com/Ghostbb-io/gb/util/gb_conv"
	"reflect"
	"strconv"
	"strings"
)

type ContentType string

const (
	ContentTypeJson       ContentType = `json`
	ContentTypeJs         ContentType = `js`
	ContentTypeXml        ContentType = `xml`
	ContentTypeIni        ContentType = `ini`
	ContentTypeYaml       ContentType = `yaml`
	ContentTypeYml        ContentType = `yml`
	ContentTypeToml       ContentType = `toml`
	ContentTypeProperties ContentType = `properties`
)

const (
	defaultSplitChar = '.' // Separator char for hierarchical data access.
)

// Json is the customized JSON struct.
type Json struct {
	mu rwmutex.RWMutex
	p  *interface{} // Pointer for hierarchical data access, it's the root of data in default.
	c  byte         // Char separator('.' in default).
	vc bool         // Violence Check(false in default), which is used to access data when the hierarchical data key contains separator char.
}

// Options for Json object creating/loading.
type Options struct {
	Safe      bool        // Mark this object is for in concurrent-safe usage. This is especially for Json object creating.
	Tags      string      // Custom priority tags for decoding, eg: "json,yaml,MyTag". This is especially for struct parsing into Json object.
	Type      ContentType // Type specifies the data content type, eg: json, xml, yaml, toml, ini.
	StrNumber bool        // StrNumber causes the Decoder to unmarshal a number into an interface{} as a string instead of as a float64.
}

// iInterfaces is used for type assert api for Interfaces().
type iInterfaces interface {
	Interfaces() []interface{}
}

// iMapStrAny is the interface support for converting struct parameter to map.
type iMapStrAny interface {
	MapStrAny() map[string]interface{}
}

// iVal is the interface for underlying interface{} retrieving.
type iVal interface {
	Val() interface{}
}

// setValue sets `value` to `j` by `pattern`.
// Note:
// 1. If value is nil and removed is true, means deleting this value;
// 2. It's quite complicated in hierarchical data search, node creating and data assignment;
func (j *Json) setValue(pattern string, value interface{}, removed bool) error {
	var (
		err    error
		array  = strings.Split(pattern, string(j.c))
		length = len(array)
	)
	if value, err = j.convertValue(value); err != nil {
		return err
	}
	// Initialization checks.
	if *j.p == nil {
		if gbstr.IsNumeric(array[0]) {
			*j.p = make([]interface{}, 0)
		} else {
			*j.p = make(map[string]interface{})
		}
	}
	var (
		pparent *interface{} = nil // Parent pointer.
		pointer *interface{} = j.p // Current pointer.
	)
	j.mu.Lock()
	defer j.mu.Unlock()
	for i := 0; i < length; i++ {
		switch (*pointer).(type) {
		case map[string]interface{}:
			if i == length-1 {
				if removed && value == nil {
					// Delete item from map.
					delete((*pointer).(map[string]interface{}), array[i])
				} else {
					if (*pointer).(map[string]interface{}) == nil {
						*pointer = map[string]interface{}{}
					}
					(*pointer).(map[string]interface{})[array[i]] = value
				}
			} else {
				// If the key does not exit in the map.
				if v, ok := (*pointer).(map[string]interface{})[array[i]]; !ok {
					if removed && value == nil {
						goto done
					}
					// Creating new node.
					if gbstr.IsNumeric(array[i+1]) {
						// Creating array node.
						n, _ := strconv.Atoi(array[i+1])
						var v interface{} = make([]interface{}, n+1)
						pparent = j.setPointerWithValue(pointer, array[i], v)
						pointer = &v
					} else {
						// Creating map node.
						var v interface{} = make(map[string]interface{})
						pparent = j.setPointerWithValue(pointer, array[i], v)
						pointer = &v
					}
				} else {
					pparent = pointer
					pointer = &v
				}
			}

		case []interface{}:
			// A string key.
			if !gbstr.IsNumeric(array[i]) {
				if i == length-1 {
					*pointer = map[string]interface{}{array[i]: value}
				} else {
					var v interface{} = make(map[string]interface{})
					*pointer = v
					pparent = pointer
					pointer = &v
				}
				continue
			}
			// Numeric index.
			valueNum, err := strconv.Atoi(array[i])
			if err != nil {
				err = gberror.WrapCodef(gbcode.CodeInvalidParameter, err, `strconv.Atoi failed for string "%s"`, array[i])
				return err
			}

			if i == length-1 {
				// Leaf node.
				if len((*pointer).([]interface{})) > valueNum {
					if removed && value == nil {
						// Deleting element.
						if pparent == nil {
							*pointer = append((*pointer).([]interface{})[:valueNum], (*pointer).([]interface{})[valueNum+1:]...)
						} else {
							j.setPointerWithValue(pparent, array[i-1], append((*pointer).([]interface{})[:valueNum], (*pointer).([]interface{})[valueNum+1:]...))
						}
					} else {
						(*pointer).([]interface{})[valueNum] = value
					}
				} else {
					if removed && value == nil {
						goto done
					}
					if pparent == nil {
						// It is the root node.
						j.setPointerWithValue(pointer, array[i], value)
					} else {
						// It is not the root node.
						s := make([]interface{}, valueNum+1)
						copy(s, (*pointer).([]interface{}))
						s[valueNum] = value
						j.setPointerWithValue(pparent, array[i-1], s)
					}
				}
			} else {
				// Branch node.
				if gbstr.IsNumeric(array[i+1]) {
					n, _ := strconv.Atoi(array[i+1])
					pSlice := (*pointer).([]interface{})
					if len(pSlice) > valueNum {
						item := pSlice[valueNum]
						if s, ok := item.([]interface{}); ok {
							for i := 0; i < n-len(s); i++ {
								s = append(s, nil)
							}
							pparent = pointer
							pointer = &pSlice[valueNum]
						} else {
							if removed && value == nil {
								goto done
							}
							var v interface{} = make([]interface{}, n+1)
							pparent = j.setPointerWithValue(pointer, array[i], v)
							pointer = &v
						}
					} else {
						if removed && value == nil {
							goto done
						}
						var v interface{} = make([]interface{}, n+1)
						pparent = j.setPointerWithValue(pointer, array[i], v)
						pointer = &v
					}
				} else {
					pSlice := (*pointer).([]interface{})
					if len(pSlice) > valueNum {
						pparent = pointer
						pointer = &(*pointer).([]interface{})[valueNum]
					} else {
						s := make([]interface{}, valueNum+1)
						copy(s, pSlice)
						s[valueNum] = make(map[string]interface{})
						if pparent != nil {
							// i > 0
							j.setPointerWithValue(pparent, array[i-1], s)
							pparent = pointer
							pointer = &s[valueNum]
						} else {
							// i = 0
							var v interface{} = s
							*pointer = v
							pparent = pointer
							pointer = &s[valueNum]
						}
					}
				}
			}

		// If the variable pointed to by the `pointer` is not of a reference type,
		// then it modifies the variable via its the parent, ie: pparent.
		default:
			if removed && value == nil {
				goto done
			}
			if gbstr.IsNumeric(array[i]) {
				n, _ := strconv.Atoi(array[i])
				s := make([]interface{}, n+1)
				if i == length-1 {
					s[n] = value
				}
				if pparent != nil {
					pparent = j.setPointerWithValue(pparent, array[i-1], s)
				} else {
					*pointer = s
					pparent = pointer
				}
			} else {
				var v1, v2 interface{}
				if i == length-1 {
					v1 = map[string]interface{}{
						array[i]: value,
					}
				} else {
					v1 = map[string]interface{}{
						array[i]: nil,
					}
				}
				if pparent != nil {
					pparent = j.setPointerWithValue(pparent, array[i-1], v1)
				} else {
					*pointer = v1
					pparent = pointer
				}
				v2 = v1.(map[string]interface{})[array[i]]
				pointer = &v2
			}
		}
	}
done:
	return nil
}

// convertValue converts `value` to map[string]interface{} or []interface{},
// which can be supported for hierarchical data access.
func (j *Json) convertValue(value interface{}) (convertedValue interface{}, err error) {
	if value == nil {
		return
	}

	switch value.(type) {
	case map[string]interface{}:
		return value, nil

	case []interface{}:
		return value, nil

	default:
		var (
			reflectInfo = reflection.OriginValueAndKind(value)
		)
		switch reflectInfo.OriginKind {
		case reflect.Array:
			return gbconv.Interfaces(value), nil

		case reflect.Slice:
			return gbconv.Interfaces(value), nil

		case reflect.Map:
			return gbconv.Map(value), nil

		case reflect.Struct:
			if v, ok := value.(iMapStrAny); ok {
				convertedValue = v.MapStrAny()
			}
			if utils.IsNil(convertedValue) {
				if v, ok := value.(iInterfaces); ok {
					convertedValue = v.Interfaces()
				}
			}
			if utils.IsNil(convertedValue) {
				convertedValue = gbconv.Map(value)
			}
			if utils.IsNil(convertedValue) {
				err = gberror.NewCodef(gbcode.CodeInvalidParameter, `unsupported value type "%s"`, reflect.TypeOf(value))
			}
			return

		default:
			return value, nil
		}
	}
}

// setPointerWithValue sets `key`:`value` to `pointer`, the `key` may be a map key or slice index.
// It returns the pointer to the new value set.
func (j *Json) setPointerWithValue(pointer *interface{}, key string, value interface{}) *interface{} {
	switch (*pointer).(type) {
	case map[string]interface{}:
		(*pointer).(map[string]interface{})[key] = value
		return &value
	case []interface{}:
		n, _ := strconv.Atoi(key)
		if len((*pointer).([]interface{})) > n {
			(*pointer).([]interface{})[n] = value
			return &(*pointer).([]interface{})[n]
		} else {
			s := make([]interface{}, n+1)
			copy(s, (*pointer).([]interface{}))
			s[n] = value
			*pointer = s
			return &s[n]
		}
	default:
		*pointer = value
	}
	return pointer
}

// getPointerByPattern returns a pointer to the value by specified `pattern`.
func (j *Json) getPointerByPattern(pattern string) *interface{} {
	if j.p == nil {
		return nil
	}
	if j.vc {
		return j.getPointerByPatternWithViolenceCheck(pattern)
	} else {
		return j.getPointerByPatternWithoutViolenceCheck(pattern)
	}
}

// getPointerByPatternWithViolenceCheck returns a pointer to the value of specified `pattern` with violence check.
func (j *Json) getPointerByPatternWithViolenceCheck(pattern string) *interface{} {
	if !j.vc {
		return j.getPointerByPatternWithoutViolenceCheck(pattern)
	}

	// It returns nil if pattern is empty.
	if pattern == "" {
		return nil
	}
	// It returns all if pattern is ".".
	if pattern == "." {
		return j.p
	}

	var (
		index   = len(pattern)
		start   = 0
		length  = 0
		pointer = j.p
	)
	if index == 0 {
		return pointer
	}
	for {
		if r := j.checkPatternByPointer(pattern[start:index], pointer); r != nil {
			if length += index - start; start > 0 {
				length += 1
			}
			start = index + 1
			index = len(pattern)
			if length == len(pattern) {
				return r
			} else {
				pointer = r
			}
		} else {
			// Get the position for next separator char.
			index = strings.LastIndexByte(pattern[start:index], j.c)
			if index != -1 && length > 0 {
				index += length + 1
			}
		}
		if start >= index {
			break
		}
	}
	return nil
}

// getPointerByPatternWithoutViolenceCheck returns a pointer to the value of specified `pattern`, with no violence check.
func (j *Json) getPointerByPatternWithoutViolenceCheck(pattern string) *interface{} {
	if j.vc {
		return j.getPointerByPatternWithViolenceCheck(pattern)
	}

	// It returns nil if pattern is empty.
	if pattern == "" {
		return nil
	}
	// It returns all if pattern is ".".
	if pattern == "." {
		return j.p
	}

	pointer := j.p
	if len(pattern) == 0 {
		return pointer
	}
	array := strings.Split(pattern, string(j.c))
	for k, v := range array {
		if r := j.checkPatternByPointer(v, pointer); r != nil {
			if k == len(array)-1 {
				return r
			} else {
				pointer = r
			}
		} else {
			break
		}
	}
	return nil
}

// checkPatternByPointer checks whether there's value by `key` in specified `pointer`.
// It returns a pointer to the value.
func (j *Json) checkPatternByPointer(key string, pointer *interface{}) *interface{} {
	switch (*pointer).(type) {
	case map[string]interface{}:
		if v, ok := (*pointer).(map[string]interface{})[key]; ok {
			return &v
		}
	case []interface{}:
		if gbstr.IsNumeric(key) {
			n, err := strconv.Atoi(key)
			if err == nil && len((*pointer).([]interface{})) > n {
				return &(*pointer).([]interface{})[n]
			}
		}
	}
	return nil
}
