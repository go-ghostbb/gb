package gbvar

import gbconv "ghostbb.io/util/gb_conv"

// MapOption specifies the option for map converting.
type MapOption = gbconv.MapOption

// Map converts and returns `v` as map[string]interface{}.
func (v *Var) Map(option ...MapOption) map[string]interface{} {
	return gbconv.Map(v.Val(), option...)
}

// MapStrAny is like function Map, but implements the interface of MapStrAny.
func (v *Var) MapStrAny(option ...MapOption) map[string]interface{} {
	return v.Map(option...)
}

// MapStrStr converts and returns `v` as map[string]string.
func (v *Var) MapStrStr(option ...MapOption) map[string]string {
	return gbconv.MapStrStr(v.Val(), option...)
}

// MapStrVar converts and returns `v` as map[string]Var.
func (v *Var) MapStrVar(option ...MapOption) map[string]*Var {
	m := v.Map(option...)
	if len(m) > 0 {
		vMap := make(map[string]*Var, len(m))
		for k, v := range m {
			vMap[k] = New(v)
		}
		return vMap
	}
	return nil
}

// MapDeep converts and returns `v` as map[string]interface{} recursively.
// Deprecated: used Map instead.
func (v *Var) MapDeep(tags ...string) map[string]interface{} {
	return gbconv.MapDeep(v.Val(), tags...)
}

// MapStrStrDeep converts and returns `v` as map[string]string recursively.
// Deprecated: used MapStrStr instead.
func (v *Var) MapStrStrDeep(tags ...string) map[string]string {
	return gbconv.MapStrStrDeep(v.Val(), tags...)
}

// MapStrVarDeep converts and returns `v` as map[string]*Var recursively.
// Deprecated: used MapStrVar instead.
func (v *Var) MapStrVarDeep(tags ...string) map[string]*Var {
	m := v.MapDeep(tags...)
	if len(m) > 0 {
		vMap := make(map[string]*Var, len(m))
		for k, v := range m {
			vMap[k] = New(v)
		}
		return vMap
	}
	return nil
}

// Maps converts and returns `v` as map[string]string.
// See gbconv.Maps.
func (v *Var) Maps(option ...MapOption) []map[string]interface{} {
	return gbconv.Maps(v.Val(), option...)
}

// MapsDeep converts `value` to []map[string]interface{} recursively.
// Deprecated: used Maps instead.
func (v *Var) MapsDeep(tags ...string) []map[string]interface{} {
	return gbconv.MapsDeep(v.Val(), tags...)
}

// MapToMap converts any map type variable `params` to another map type variable `pointer`.
// See gbconv.MapToMap.
func (v *Var) MapToMap(pointer interface{}, mapping ...map[string]string) (err error) {
	return gbconv.MapToMap(v.Val(), pointer, mapping...)
}

// MapToMaps converts any map type variable `params` to another map type variable `pointer`.
// See gbconv.MapToMaps.
func (v *Var) MapToMaps(pointer interface{}, mapping ...map[string]string) (err error) {
	return gbconv.MapToMaps(v.Val(), pointer, mapping...)
}

// MapToMapsDeep converts any map type variable `params` to another map type variable
// `pointer` recursively.
// See gbconv.MapToMapsDeep.
func (v *Var) MapToMapsDeep(pointer interface{}, mapping ...map[string]string) (err error) {
	return gbconv.MapToMaps(v.Val(), pointer, mapping...)
}
