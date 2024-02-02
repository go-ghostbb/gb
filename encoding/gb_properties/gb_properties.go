// Package gbproperties provides accessing and converting for .properties content.
package gbproperties

import (
	"bytes"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/json"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"sort"
	"strings"

	"github.com/magiconair/properties"
)

// Decode converts properties format to map.
func Decode(data []byte) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})
	pr, err := properties.Load(data, properties.UTF8)
	if err != nil || pr == nil {
		err = gberror.Wrapf(err, `Lib magiconair load Properties data failed.`)
		return nil, err
	}
	for _, key := range pr.Keys() {
		// ignore existence check: we know it's there
		value, _ := pr.Get(key)
		// recursively build nested maps
		path := strings.Split(key, ".")
		lastKey := strings.ToLower(path[len(path)-1])
		deepestMap := deepSearch(res, path[0:len(path)-1])

		// set innermost value
		deepestMap[lastKey] = value
	}
	return res, nil
}

// Encode converts map to properties format.
func Encode(data map[string]interface{}) (res []byte, err error) {
	pr := properties.NewProperties()

	flattened := map[string]interface{}{}

	flattened = flattenAndMergeMap(flattened, data, "", ".")

	keys := make([]string, 0, len(flattened))

	for key := range flattened {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		_, _, err := pr.Set(key, gbconv.String(flattened[key]))
		if err != nil {
			err = gberror.Wrapf(err, `Sets the property key to the corresponding value failed.`)
			return nil, err
		}
	}

	var buf bytes.Buffer

	_, err = pr.Write(&buf, properties.UTF8)
	if err != nil {
		err = gberror.Wrapf(err, `Properties Write buf failed.`)
		return nil, err
	}

	return buf.Bytes(), nil
}

// ToJson convert .properties format to JSON.
func ToJson(data []byte) (res []byte, err error) {
	prMap, err := Decode(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(prMap)
}

// deepSearch scans deep maps, following the key indexes listed in the sequence "path".
// The last value is expected to be another map, and is returned.
func deepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			// intermediate key does not exist
			// => create it and continue from there
			m3 := make(map[string]interface{})
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]interface{})
		if !ok {
			m3 = make(map[string]interface{})
			m[k] = m3
		}
		// continue search from here
		m = m3
	}
	return m
}

// flattenAndMergeMap recursively flattens the given map into a new map
func flattenAndMergeMap(shadow map[string]interface{}, m map[string]interface{}, prefix string, delimiter string) map[string]interface{} {
	if shadow != nil && prefix != "" && shadow[prefix] != nil {
		return shadow
	}

	var m2 map[string]interface{}
	if prefix != "" {
		prefix += delimiter
	}
	for k, val := range m {
		fullKey := prefix + k
		switch v := val.(type) {
		case map[string]interface{}:
			m2 = v
		case map[interface{}]interface{}:
			m2 = gbconv.Map(v)
		default:
			// immediate value
			shadow[strings.ToLower(fullKey)] = v
			continue
		}
		// recursively merge to shadow map
		shadow = flattenAndMergeMap(shadow, m2, fullKey, delimiter)
	}
	return shadow
}
