// Package gbtype provides high performance and concurrent-safe basic variable types.
package gbtype

// New is alias of NewInterface.
// See NewInterface.
func New(value ...interface{}) *Interface {
	return NewInterface(value...)
}
