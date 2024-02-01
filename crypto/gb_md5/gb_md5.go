// Package gbmd5 provides useful API for MD5 encryption algorithms.
package gbmd5

import (
	"crypto/md5"
	"fmt"
	gberror "github.com/Ghostbb-io/gb/errors/gb_error"
	gbconv "github.com/Ghostbb-io/gb/util/gb_conv"
	"io"
	"os"
)

// Encrypt encrypts any type of variable using MD5 algorithms.
// It uses gbconv package to convert `v` to its bytes type.
func Encrypt(data interface{}) (encrypt string, err error) {
	return EncryptBytes(gbconv.Bytes(data))
}

// MustEncrypt encrypts any type of variable using MD5 algorithms.
// It uses gbconv package to convert `v` to its bytes type.
// It panics if any error occurs.
func MustEncrypt(data interface{}) string {
	result, err := Encrypt(data)
	if err != nil {
		panic(err)
	}
	return result
}

// EncryptBytes encrypts `data` using MD5 algorithms.
func EncryptBytes(data []byte) (encrypt string, err error) {
	h := md5.New()
	if _, err = h.Write(data); err != nil {
		err = gberror.Wrap(err, `hash.Write failed`)
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// MustEncryptBytes encrypts `data` using MD5 algorithms.
// It panics if any error occurs.
func MustEncryptBytes(data []byte) string {
	result, err := EncryptBytes(data)
	if err != nil {
		panic(err)
	}
	return result
}

// EncryptString encrypts string `data` using MD5 algorithms.
func EncryptString(data string) (encrypt string, err error) {
	return EncryptBytes([]byte(data))
}

// MustEncryptString encrypts string `data` using MD5 algorithms.
// It panics if any error occurs.
func MustEncryptString(data string) string {
	result, err := EncryptString(data)
	if err != nil {
		panic(err)
	}
	return result
}

// EncryptFile encrypts file content of `path` using MD5 algorithms.
func EncryptFile(path string) (encrypt string, err error) {
	f, err := os.Open(path)
	if err != nil {
		err = gberror.Wrapf(err, `os.Open failed for name "%s"`, path)
		return "", err
	}
	defer f.Close()
	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		err = gberror.Wrap(err, `io.Copy failed`)
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// MustEncryptFile encrypts file content of `path` using MD5 algorithms.
// It panics if any error occurs.
func MustEncryptFile(path string) string {
	result, err := EncryptFile(path)
	if err != nil {
		panic(err)
	}
	return result
}
