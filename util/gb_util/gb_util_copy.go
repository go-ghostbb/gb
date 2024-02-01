package gbutil

import (
	"ghostbb.io/internal/deepcopy"
)

// Copy returns a deep copy of v.
//
// Copy is unable to copy unexported fields in a struct (lowercase field names).
// Unexported fields can't be reflected by the Go runtime and therefore
// they can't perform any data copies.
func Copy(src interface{}) (dst interface{}) {
	return deepcopy.Copy(src)
}
