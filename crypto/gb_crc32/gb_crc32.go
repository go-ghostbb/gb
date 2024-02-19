// Package gbcrc32 provides useful API for CRC32 encryption algorithms.
package gbcrc32

import (
	gbconv "ghostbb.io/gb/util/gb_conv"
	"hash/crc32"
)

// Encrypt encrypts any type of variable using CRC32 algorithms.
// It uses gbconv package to convert `v` to its bytes type.
func Encrypt(v interface{}) uint32 {
	return crc32.ChecksumIEEE(gbconv.Bytes(v))
}
