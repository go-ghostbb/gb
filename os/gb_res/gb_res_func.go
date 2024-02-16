package gbres

import (
	"archive/zip"
	"bytes"
	"encoding/hex"
	"fmt"
	gbbase64 "ghostbb.io/gb/encoding/gb_base64"
	gbcompress "ghostbb.io/gb/encoding/gb_compress"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbstr "ghostbb.io/gb/text/gb_str"
)

const (
	packedGoSourceTemplate = `
package %s

import gbres "ghostbb.io/gb/os/gb_res"

func init() {
	if err := gbres.Add("%s"); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}
`
)

// Option contains the extra options for Pack functions.
type Option struct {
	Prefix   string // The file path prefix for each file item in resource manager.
	KeepPath bool   // Keep the passed path when packing, usually for relative path.
}

// Pack packs the path specified by `srcPaths` into bytes.
// The unnecessary parameter `keyPrefix` indicates the prefix for each file
// packed into the result bytes.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
//
// Deprecated: use PackWithOption instead.
func Pack(srcPaths string, keyPrefix ...string) ([]byte, error) {
	option := Option{}
	if len(keyPrefix) > 0 && keyPrefix[0] != "" {
		option.Prefix = keyPrefix[0]
	}
	return PackWithOption(srcPaths, option)
}

// PackWithOption packs the path specified by `srcPaths` into bytes.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
func PackWithOption(srcPaths string, option Option) ([]byte, error) {
	var buffer = bytes.NewBuffer(nil)
	err := zipPathWriter(srcPaths, buffer, option)
	if err != nil {
		return nil, err
	}
	// Gzip the data bytes to reduce the size.
	return gbcompress.Gzip(buffer.Bytes(), 9)
}

// PackToFile packs the path specified by `srcPaths` to target file `dstPath`.
// The unnecessary parameter `keyPrefix` indicates the prefix for each file
// packed into the result bytes.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
//
// Deprecated: use PackToFileWithOption instead.
func PackToFile(srcPaths, dstPath string, keyPrefix ...string) error {
	data, err := Pack(srcPaths, keyPrefix...)
	if err != nil {
		return err
	}
	return gbfile.PutBytes(dstPath, data)
}

// PackToFileWithOption packs the path specified by `srcPaths` to target file `dstPath`.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
func PackToFileWithOption(srcPaths, dstPath string, option Option) error {
	data, err := PackWithOption(srcPaths, option)
	if err != nil {
		return err
	}
	return gbfile.PutBytes(dstPath, data)
}

// PackToGoFile packs the path specified by `srcPaths` to target go file `goFilePath`
// with given package name `pkgName`.
//
// The unnecessary parameter `keyPrefix` indicates the prefix for each file
// packed into the result bytes.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
//
// Deprecated: use PackToGoFileWithOption instead.
func PackToGoFile(srcPath, goFilePath, pkgName string, keyPrefix ...string) error {
	data, err := Pack(srcPath, keyPrefix...)
	if err != nil {
		return err
	}
	return gbfile.PutContents(
		goFilePath,
		fmt.Sprintf(gbstr.TrimLeft(packedGoSourceTemplate), pkgName, gbbase64.EncodeToString(data)),
	)
}

// PackToGoFileWithOption packs the path specified by `srcPaths` to target go file `goFilePath`
// with given package name `pkgName`.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
func PackToGoFileWithOption(srcPath, goFilePath, pkgName string, option Option) error {
	data, err := PackWithOption(srcPath, option)
	if err != nil {
		return err
	}
	return gbfile.PutContents(
		goFilePath,
		fmt.Sprintf(gbstr.TrimLeft(packedGoSourceTemplate), pkgName, gbbase64.EncodeToString(data)),
	)
}

// Unpack unpacks the content specified by `path` to []*File.
func Unpack(path string) ([]*File, error) {
	realPath, err := gbfile.Search(path)
	if err != nil {
		return nil, err
	}
	return UnpackContent(gbfile.GetContents(realPath))
}

// UnpackContent unpacks the content to []*File.
func UnpackContent(content string) ([]*File, error) {
	var (
		err  error
		data []byte
	)
	if isHexStr(content) {
		// It here keeps compatible with old version packing string using hex string.
		// TODO remove this support in the future.
		data, err = gbcompress.UnGzip(hexStrToBytes(content))
		if err != nil {
			return nil, err
		}
	} else if isBase64(content) {
		// New version packing string using base64.
		b, err := gbbase64.DecodeString(content)
		if err != nil {
			return nil, err
		}
		data, err = gbcompress.UnGzip(b)
		if err != nil {
			return nil, err
		}
	} else {
		data, err = gbcompress.UnGzip([]byte(content))
		if err != nil {
			return nil, err
		}
	}
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		err = gberror.Wrapf(err, `create zip reader failed`)
		return nil, err
	}
	array := make([]*File, len(reader.File))
	for i, file := range reader.File {
		array[i] = &File{file: file}
	}
	return array, nil
}

// isBase64 checks and returns whether given content `s` is base64 string.
// It returns true if `s` is base64 string, or false if not.
func isBase64(s string) bool {
	var r bool
	for i := 0; i < len(s); i++ {
		r = (s[i] >= '0' && s[i] <= '9') ||
			(s[i] >= 'a' && s[i] <= 'z') ||
			(s[i] >= 'A' && s[i] <= 'Z') ||
			(s[i] == '+' || s[i] == '-') ||
			(s[i] == '_' || s[i] == '/') || s[i] == '='
		if !r {
			return false
		}
	}
	return true
}

// isHexStr checks and returns whether given content `s` is hex string.
// It returns true if `s` is hex string, or false if not.
func isHexStr(s string) bool {
	var r bool
	for i := 0; i < len(s); i++ {
		r = (s[i] >= '0' && s[i] <= '9') ||
			(s[i] >= 'a' && s[i] <= 'f') ||
			(s[i] >= 'A' && s[i] <= 'F')
		if !r {
			return false
		}
	}
	return true
}

// hexStrToBytes converts hex string content to []byte.
func hexStrToBytes(s string) []byte {
	src := []byte(s)
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, _ = hex.Decode(dst, src)
	return dst
}
