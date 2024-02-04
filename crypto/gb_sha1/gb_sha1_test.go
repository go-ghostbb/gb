package gbsha1_test

import (
	gbsha1 "ghostbb.io/gb/crypto/gb_sha1"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func TestEncrypt(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		test1 := gbsha1.Encrypt("ghostbb.io")
		t.Assert(test1, "fe6696b8cce1f8f4562f21e6b59d5f71d9941a1a")
	})
}
