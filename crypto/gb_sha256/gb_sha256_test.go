package gbsha256_test

import (
	gbsha256 "ghostbb.io/gb/crypto/gb_sha256"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func TestEncrypt256(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		test1 := gbsha256.Encrypt256("ghostbb.io")
		t.Assert(test1, "5fe0b4a7aafb5015c64fba6ed4edca99d65b3d88c733810f6bd2c03d4e0ab0a2")
	})
}
