// Package gbaes provides useful API for AES encryption/decryption algorithms.
package gbaes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	gbcode "github.com/Ghostbb-io/gb/errors/gb_code"
	gberror "github.com/Ghostbb-io/gb/errors/gb_error"
)

const (
	// IVDefaultValue is the default value for IV.
	IVDefaultValue = "Ghostbb frame"
)

// Encrypt is alias of EncryptCBC.
func Encrypt(plainText []byte, key []byte, iv ...[]byte) ([]byte, error) {
	return EncryptCBC(plainText, key, iv...)
}

// Decrypt is alias of DecryptCBC.
func Decrypt(cipherText []byte, key []byte, iv ...[]byte) ([]byte, error) {
	return DecryptCBC(cipherText, key, iv...)
}

// EncryptCBC encrypts `plainText` using CBC mode.
// Note that the key must be 16/24/32 bit length.
// The parameter `iv` initialization vector is unnecessary.
func EncryptCBC(plainText []byte, key []byte, iv ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		err = gberror.WrapCodef(gbcode.CodeInvalidParameter, err, `aes.NewCipher failed for key "%s"`, key)
		return nil, err
	}
	blockSize := block.BlockSize()
	plainText = PKCS7Padding(plainText, blockSize)
	ivValue := ([]byte)(nil)
	if len(iv) > 0 {
		ivValue = iv[0]
	} else {
		ivValue = []byte(IVDefaultValue)
	}
	blockMode := cipher.NewCBCEncrypter(block, ivValue)
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)

	return cipherText, nil
}

// DecryptCBC decrypts `cipherText` using CBC mode.
// Note that the key must be 16/24/32 bit length.
// The parameter `iv` initialization vector is unnecessary.
func DecryptCBC(cipherText []byte, key []byte, iv ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		err = gberror.WrapCodef(gbcode.CodeInvalidParameter, err, `aes.NewCipher failed for key "%s"`, key)
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(cipherText) < blockSize {
		return nil, gberror.NewCode(gbcode.CodeInvalidParameter, "cipherText too short")
	}
	ivValue := ([]byte)(nil)
	if len(iv) > 0 {
		ivValue = iv[0]
	} else {
		ivValue = []byte(IVDefaultValue)
	}
	if len(cipherText)%blockSize != 0 {
		return nil, gberror.NewCode(gbcode.CodeInvalidParameter, "cipherText is not a multiple of the block size")
	}
	blockModel := cipher.NewCBCDecrypter(block, ivValue)
	plainText := make([]byte, len(cipherText))
	blockModel.CryptBlocks(plainText, cipherText)
	plainText, e := PKCS7UnPadding(plainText, blockSize)
	if e != nil {
		return nil, e
	}
	return plainText, nil
}

// PKCS5Padding applies PKCS#5 padding to the source byte slice to match the given block size.
//
// If the block size is not provided, it defaults to 8.
func PKCS5Padding(src []byte, blockSize ...int) []byte {
	blockSizeTemp := 8
	if len(blockSize) > 0 {
		blockSizeTemp = blockSize[0]
	}
	return PKCS7Padding(src, blockSizeTemp)
}

// PKCS5UnPadding removes PKCS#5 padding from the source byte slice based on the given block size.
//
// If the block size is not provided, it defaults to 8.
func PKCS5UnPadding(src []byte, blockSize ...int) ([]byte, error) {
	blockSizeTemp := 8
	if len(blockSize) > 0 {
		blockSizeTemp = blockSize[0]
	}
	return PKCS7UnPadding(src, blockSizeTemp)
}

// PKCS7Padding applies PKCS#7 padding to the source byte slice to match the given block size.
func PKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// PKCS7UnPadding removes PKCS#7 padding from the source byte slice based on the given block size.
func PKCS7UnPadding(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	if blockSize <= 0 {
		return nil, gberror.NewCode(gbcode.CodeInvalidParameter, fmt.Sprintf("invalid blockSize: %d", blockSize))
	}

	if length%blockSize != 0 || length == 0 {
		return nil, gberror.NewCode(gbcode.CodeInvalidParameter, "invalid data len")
	}

	unpadding := int(src[length-1])
	if unpadding > blockSize || unpadding == 0 {
		return nil, gberror.NewCode(gbcode.CodeInvalidParameter, "invalid unpadding")
	}

	padding := src[length-unpadding:]
	for i := 0; i < unpadding; i++ {
		if padding[i] != byte(unpadding) {
			return nil, gberror.NewCode(gbcode.CodeInvalidParameter, "invalid padding")
		}
	}

	return src[:(length - unpadding)], nil
}

// EncryptCFB encrypts `plainText` using CFB mode.
// Note that the key must be 16/24/32 bit length.
// The parameter `iv` initialization vector is unnecessary.
func EncryptCFB(plainText []byte, key []byte, padding *int, iv ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		err = gberror.WrapCodef(gbcode.CodeInvalidParameter, err, `aes.NewCipher failed for key "%s"`, key)
		return nil, err
	}
	blockSize := block.BlockSize()
	plainText, *padding = ZeroPadding(plainText, blockSize)
	ivValue := ([]byte)(nil)
	if len(iv) > 0 {
		ivValue = iv[0]
	} else {
		ivValue = []byte(IVDefaultValue)
	}
	stream := cipher.NewCFBEncrypter(block, ivValue)
	cipherText := make([]byte, len(plainText))
	stream.XORKeyStream(cipherText, plainText)
	return cipherText, nil
}

// DecryptCFB decrypts `plainText` using CFB mode.
// Note that the key must be 16/24/32 bit length.
// The parameter `iv` initialization vector is unnecessary.
func DecryptCFB(cipherText []byte, key []byte, unPadding int, iv ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		err = gberror.WrapCodef(gbcode.CodeInvalidParameter, err, `aes.NewCipher failed for key "%s"`, key)
		return nil, err
	}
	if len(cipherText) < aes.BlockSize {
		return nil, gberror.NewCode(gbcode.CodeInvalidParameter, "cipherText too short")
	}
	ivValue := ([]byte)(nil)
	if len(iv) > 0 {
		ivValue = iv[0]
	} else {
		ivValue = []byte(IVDefaultValue)
	}
	stream := cipher.NewCFBDecrypter(block, ivValue)
	plainText := make([]byte, len(cipherText))
	stream.XORKeyStream(plainText, cipherText)
	plainText = ZeroUnPadding(plainText, unPadding)
	return plainText, nil
}

func ZeroPadding(cipherText []byte, blockSize int) ([]byte, int) {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(0)}, padding)
	return append(cipherText, padText...), padding
}

func ZeroUnPadding(plaintext []byte, unPadding int) []byte {
	length := len(plaintext)
	return plaintext[:(length - unPadding)]
}
