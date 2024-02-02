// Package gbtoml provides accessing and converting for TOML content.
package gbtoml

import (
	"bytes"
	gberror "ghostbb.io/gb/errors/gb_error"
	"ghostbb.io/gb/internal/json"

	"github.com/BurntSushi/toml"
)

func Encode(v interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if err := toml.NewEncoder(buffer).Encode(v); err != nil {
		err = gberror.Wrap(err, `toml.Encoder.Encode failed`)
		return nil, err
	}
	return buffer.Bytes(), nil
}

func Decode(v []byte) (interface{}, error) {
	var result interface{}
	if err := toml.Unmarshal(v, &result); err != nil {
		err = gberror.Wrap(err, `toml.Unmarshal failed`)
		return nil, err
	}
	return result, nil
}

func DecodeTo(v []byte, result interface{}) (err error) {
	err = toml.Unmarshal(v, result)
	if err != nil {
		err = gberror.Wrap(err, `toml.Unmarshal failed`)
	}
	return err
}

func ToJson(v []byte) ([]byte, error) {
	if r, err := Decode(v); err != nil {
		return nil, err
	} else {
		return json.Marshal(r)
	}
}
