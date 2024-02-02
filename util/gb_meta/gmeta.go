// Package gbmeta provides embedded meta data feature for struct.
package gbmeta

import (
	gbvar "ghostbb.io/gb/container/gb_var"
	gbstructs "ghostbb.io/gb/os/gb_structs"
)

// Meta is used as an embedded attribute for struct to enabled metadata feature.
type Meta struct{}

const (
	metaAttributeName = "Meta"        // metaAttributeName is the attribute name of metadata in struct.
	metaTypeName      = "gbmeta.Meta" // metaTypeName is for type string comparison.
)

// Data retrieves and returns all metadata from `object`.
func Data(object interface{}) map[string]string {
	reflectType, err := gbstructs.StructType(object)
	if err != nil {
		return nil
	}
	if field, ok := reflectType.FieldByName(metaAttributeName); ok {
		if field.Type.String() == metaTypeName {
			return gbstructs.ParseTag(string(field.Tag))
		}
	}
	return map[string]string{}
}

// Get retrieves and returns specified metadata by `key` from `object`.
func Get(object interface{}, key string) *gbvar.Var {
	v, ok := Data(object)[key]
	if !ok {
		return nil
	}
	return gbvar.New(v)
}
