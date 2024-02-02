package gbaes_test

import (
	gbaes "ghostbb.io/gb/crypto/gb_aes"
	gbbase64 "ghostbb.io/gb/encoding/gb_base64"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

var (
	iv, _  = gbbase64.Decode([]byte("TDp4nh5Bc1+9jloLgU3nMA=="))
	key, _ = gbbase64.Decode([]byte("dd79wy1M9ZCVGTcHMJB2bSlmuScjSaOcnj7hxv02aWc="))
)

func TestEncrypt(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		data, err := gbaes.Encrypt([]byte("Gj52676089"), key, iv)
		t.AssertNil(err)
		t.Assert(string(gbbase64.Encode(data)), "wL9/qsUE/SdbIInpE9H7YQ==")
	})
}

func TestDecrypt(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		text, err := gbbase64.Decode([]byte("wL9/qsUE/SdbIInpE9H7YQ=="))
		t.AssertNil(err)
		data, err := gbaes.Decrypt(text, key, iv)
		t.AssertNil(err)
		t.Assert(string(data), "Gj52676089")
	})
}
