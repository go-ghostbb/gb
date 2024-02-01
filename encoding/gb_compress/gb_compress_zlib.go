// Package gbcompress provides kinds of compression algorithms for binary/bytes data.
package gbcompress

import (
	"bytes"
	"compress/zlib"
	gberror "github.com/Ghostbb-io/gb/errors/gb_error"
	"io"
)

// Zlib compresses `data` with zlib algorithm.
func Zlib(data []byte) ([]byte, error) {
	if data == nil || len(data) < 13 {
		return data, nil
	}
	var (
		err    error
		in     bytes.Buffer
		writer = zlib.NewWriter(&in)
	)

	if _, err = writer.Write(data); err != nil {
		err = gberror.Wrapf(err, `zlib.Writer.Write failed`)
		return nil, err
	}
	if err = writer.Close(); err != nil {
		err = gberror.Wrapf(err, `zlib.Writer.Close failed`)
		return in.Bytes(), err
	}
	return in.Bytes(), nil
}

// UnZlib decompresses `data` with zlib algorithm.
func UnZlib(data []byte) ([]byte, error) {
	if data == nil || len(data) < 13 {
		return data, nil
	}
	var (
		out             bytes.Buffer
		bytesReader     = bytes.NewReader(data)
		zlibReader, err = zlib.NewReader(bytesReader)
	)
	if err != nil {
		err = gberror.Wrapf(err, `zlib.NewReader failed`)
		return nil, err
	}
	if _, err = io.Copy(&out, zlibReader); err != nil {
		err = gberror.Wrapf(err, `io.Copy failed`)
		return nil, err
	}
	return out.Bytes(), nil
}
