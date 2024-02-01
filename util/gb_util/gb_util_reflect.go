package gbutil

import (
	"ghostbb.io/internal/reflection"
)

type (
	OriginValueAndKindOutput = reflection.OriginValueAndKindOutput
	OriginTypeAndKindOutput  = reflection.OriginTypeAndKindOutput
)

// OriginValueAndKind retrieves and returns the original reflect value and kind.
func OriginValueAndKind(value interface{}) (out OriginValueAndKindOutput) {
	return reflection.OriginValueAndKind(value)
}

// OriginTypeAndKind retrieves and returns the original reflect type and kind.
func OriginTypeAndKind(value interface{}) (out OriginTypeAndKindOutput) {
	return reflection.OriginTypeAndKind(value)
}
