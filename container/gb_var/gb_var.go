// Package gbvar provides an universal variable type, like generics.
package gbvar

import (
	gbtype "ghostbb.io/gb/container/gb_type"
	"ghostbb.io/gb/internal/deepcopy"
	"ghostbb.io/gb/internal/json"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbutil "ghostbb.io/gb/util/gb_util"
	"time"
)

// Var is an universal variable type implementer.
type Var struct {
	value interface{} // Underlying value.
	safe  bool        // Concurrent safe or not.
}

// New creates and returns a new Var with given `value`.
// The optional parameter `safe` specifies whether Var is used in concurrent-safety,
// which is false in default.
func New(value interface{}, safe ...bool) *Var {
	if len(safe) > 0 && safe[0] {
		return &Var{
			value: gbtype.NewInterface(value),
			safe:  true,
		}
	}
	return &Var{
		value: value,
	}
}

// Copy does a deep copy of current Var and returns a pointer to this Var.
func (v *Var) Copy() *Var {
	return New(gbutil.Copy(v.Val()), v.safe)
}

// Clone does a shallow copy of current Var and returns a pointer to this Var.
func (v *Var) Clone() *Var {
	return New(v.Val(), v.safe)
}

// Set sets `value` to `v`, and returns the old value.
func (v *Var) Set(value interface{}) (old interface{}) {
	if v.safe {
		if t, ok := v.value.(*gbtype.Interface); ok {
			old = t.Set(value)
			return
		}
	}
	old = v.value
	v.value = value
	return
}

// Val returns the current value of `v`.
func (v *Var) Val() interface{} {
	if v == nil {
		return nil
	}
	if v.safe {
		if t, ok := v.value.(*gbtype.Interface); ok {
			return t.Val()
		}
	}
	return v.value
}

// Interface is alias of Val.
func (v *Var) Interface() interface{} {
	return v.Val()
}

// Bytes converts and returns `v` as []byte.
func (v *Var) Bytes() []byte {
	return gbconv.Bytes(v.Val())
}

// String converts and returns `v` as string.
func (v *Var) String() string {
	return gbconv.String(v.Val())
}

// Bool converts and returns `v` as bool.
func (v *Var) Bool() bool {
	return gbconv.Bool(v.Val())
}

// Int converts and returns `v` as int.
func (v *Var) Int() int {
	return gbconv.Int(v.Val())
}

// Int8 converts and returns `v` as int8.
func (v *Var) Int8() int8 {
	return gbconv.Int8(v.Val())
}

// Int16 converts and returns `v` as int16.
func (v *Var) Int16() int16 {
	return gbconv.Int16(v.Val())
}

// Int32 converts and returns `v` as int32.
func (v *Var) Int32() int32 {
	return gbconv.Int32(v.Val())
}

// Int64 converts and returns `v` as int64.
func (v *Var) Int64() int64 {
	return gbconv.Int64(v.Val())
}

// Uint converts and returns `v` as uint.
func (v *Var) Uint() uint {
	return gbconv.Uint(v.Val())
}

// Uint8 converts and returns `v` as uint8.
func (v *Var) Uint8() uint8 {
	return gbconv.Uint8(v.Val())
}

// Uint16 converts and returns `v` as uint16.
func (v *Var) Uint16() uint16 {
	return gbconv.Uint16(v.Val())
}

// Uint32 converts and returns `v` as uint32.
func (v *Var) Uint32() uint32 {
	return gbconv.Uint32(v.Val())
}

// Uint64 converts and returns `v` as uint64.
func (v *Var) Uint64() uint64 {
	return gbconv.Uint64(v.Val())
}

// Float32 converts and returns `v` as float32.
func (v *Var) Float32() float32 {
	return gbconv.Float32(v.Val())
}

// Float64 converts and returns `v` as float64.
func (v *Var) Float64() float64 {
	return gbconv.Float64(v.Val())
}

// Time converts and returns `v` as time.Time.
// The parameter `format` specifies the format of the time string using gtime,
// eg: Y-m-d H:i:s.
func (v *Var) Time(format ...string) time.Time {
	return gbconv.Time(v.Val(), format...)
}

// Duration converts and returns `v` as time.Duration.
// If value of `v` is string, then it uses time.ParseDuration for conversion.
func (v *Var) Duration() time.Duration {
	return gbconv.Duration(v.Val())
}

// GTime converts and returns `v` as *gtime.Time.
// The parameter `format` specifies the format of the time string using gtime,
// eg: Y-m-d H:i:s.
func (v *Var) GTime(format ...string) *gbtime.Time {
	return gbconv.GBTime(v.Val(), format...)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (v Var) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Val())
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (v *Var) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.UnmarshalUseNumber(b, &i); err != nil {
		return err
	}
	v.Set(i)
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for Var.
func (v *Var) UnmarshalValue(value interface{}) error {
	v.Set(value)
	return nil
}

// DeepCopy implements interface for deep copy of current type.
func (v *Var) DeepCopy() interface{} {
	if v == nil {
		return nil
	}
	return New(deepcopy.Copy(v.Val()), v.safe)
}
