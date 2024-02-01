package gbvar

import gbconv "github.com/Ghostbb-io/gb/util/gb_conv"

// Ints converts and returns `v` as []int.
func (v *Var) Ints() []int {
	return gbconv.Ints(v.Val())
}

// Int64s converts and returns `v` as []int64.
func (v *Var) Int64s() []int64 {
	return gbconv.Int64s(v.Val())
}

// Uints converts and returns `v` as []uint.
func (v *Var) Uints() []uint {
	return gbconv.Uints(v.Val())
}

// Uint64s converts and returns `v` as []uint64.
func (v *Var) Uint64s() []uint64 {
	return gbconv.Uint64s(v.Val())
}

// Floats is alias of Float64s.
func (v *Var) Floats() []float64 {
	return gbconv.Floats(v.Val())
}

// Float32s converts and returns `v` as []float32.
func (v *Var) Float32s() []float32 {
	return gbconv.Float32s(v.Val())
}

// Float64s converts and returns `v` as []float64.
func (v *Var) Float64s() []float64 {
	return gbconv.Float64s(v.Val())
}

// Strings converts and returns `v` as []string.
func (v *Var) Strings() []string {
	return gbconv.Strings(v.Val())
}

// Interfaces converts and returns `v` as []interfaces{}.
func (v *Var) Interfaces() []interface{} {
	return gbconv.Interfaces(v.Val())
}

// Slice is alias of Interfaces.
func (v *Var) Slice() []interface{} {
	return v.Interfaces()
}

// Array is alias of Interfaces.
func (v *Var) Array() []interface{} {
	return v.Interfaces()
}

// Vars converts and returns `v` as []Var.
func (v *Var) Vars() []*Var {
	array := gbconv.Interfaces(v.Val())
	if len(array) == 0 {
		return nil
	}
	vars := make([]*Var, len(array))
	for k, v := range array {
		vars[k] = New(v)
	}
	return vars
}
