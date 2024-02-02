// Package gbyaml provides accessing and converting for YAML content.
package gbyaml

import (
	"bytes"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/json"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Encode encodes `value` to an YAML format content as bytes.
func Encode(value interface{}) (out []byte, err error) {
	if out, err = yaml.Marshal(value); err != nil {
		err = gberror.Wrap(err, `yaml.Marshal failed`)
	}
	return
}

// EncodeIndent encodes `value` to an YAML format content with indent as bytes.
func EncodeIndent(value interface{}, indent string) (out []byte, err error) {
	out, err = Encode(value)
	if err != nil {
		return
	}
	if indent != "" {
		var (
			buffer = bytes.NewBuffer(nil)
			array  = strings.Split(strings.TrimSpace(string(out)), "\n")
		)
		for _, v := range array {
			buffer.WriteString(indent)
			buffer.WriteString(v)
			buffer.WriteString("\n")
		}
		out = buffer.Bytes()
	}
	return
}

// Decode parses `content` into and returns as map.
func Decode(content []byte) (map[string]interface{}, error) {
	var (
		result map[string]interface{}
		err    error
	)
	if err = yaml.Unmarshal(content, &result); err != nil {
		err = gberror.Wrap(err, `yaml.Unmarshal failed`)
		return nil, err
	}
	return gbconv.MapDeep(result), nil
}

// DecodeTo parses `content` into `result`.
func DecodeTo(value []byte, result interface{}) (err error) {
	err = yaml.Unmarshal(value, result)
	if err != nil {
		err = gberror.Wrap(err, `yaml.Unmarshal failed`)
	}
	return
}

// ToJson converts `content` to JSON format content.
func ToJson(content []byte) (out []byte, err error) {
	var (
		result interface{}
	)
	if result, err = Decode(content); err != nil {
		return nil, err
	} else {
		return json.Marshal(result)
	}
}
