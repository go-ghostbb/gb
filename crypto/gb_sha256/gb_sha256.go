package gbsha256

import (
	"crypto/sha256"
	"encoding/hex"
	gbconv "ghostbb.io/gb/util/gb_conv"
)

// Encrypt256 encrypts any type of variable using SHA256 algorithms.
// It uses package gbconv to convert `v` to its bytes type.
func Encrypt256(v interface{}) string {
	r := sha256.Sum256(gbconv.Bytes(v))
	return hex.EncodeToString(r[:])
}
