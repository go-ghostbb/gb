package gbdebug

import (
	"crypto/md5"
	"fmt"
	gbhash "ghostbb.io/encoding/gb_hash"
	gberror "ghostbb.io/errors/gb_error"
	"io"
	"os"
	"strconv"
)

// BinVersion returns the version of current running binary.
// It uses gbhash.BKDRHash+BASE36 algorithm to calculate the unique version of the binary.
func BinVersion() string {
	if binaryVersion == "" {
		binaryContent, _ := os.ReadFile(selfPath)
		binaryVersion = strconv.FormatInt(
			int64(gbhash.BKDR(binaryContent)),
			36,
		)
	}
	return binaryVersion
}

// BinVersionMd5 returns the version of current running binary.
// It uses MD5 algorithm to calculate the unique version of the binary.
func BinVersionMd5() string {
	if binaryVersionMd5 == "" {
		binaryVersionMd5, _ = md5File(selfPath)
	}
	return binaryVersionMd5
}

// md5File encrypts file content of `path` using MD5 algorithms.
func md5File(path string) (encrypt string, err error) {
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
