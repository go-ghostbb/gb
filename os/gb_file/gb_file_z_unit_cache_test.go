package gbfile_test

import (
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtest "ghostbb.io/gb/test/gb_test"
	"os"
	"testing"
	"time"
)

func Test_GetContentsWithCache(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var f *os.File
		var err error
		fileName := "test"
		strTest := "123"

		if !gbfile.Exists(fileName) {
			f, err = os.CreateTemp("", fileName)
			if err != nil {
				t.Error("create file fail")
			}
		}

		defer f.Close()
		defer os.Remove(f.Name())

		if gbfile.Exists(f.Name()) {
			f, err = gbfile.OpenFile(f.Name(), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				t.Error("file open fail", err)
			}

			err = gbfile.PutContents(f.Name(), strTest)
			if err != nil {
				t.Error("write error", err)
			}

			cache := gbfile.GetContentsWithCache(f.Name(), 1)
			t.Assert(cache, strTest)
		}
	})

	gbtest.C(t, func(t *gbtest.T) {

		var f *os.File
		var err error
		fileName := "test2"
		strTest := "123"

		if !gbfile.Exists(fileName) {
			f, err = os.CreateTemp("", fileName)
			if err != nil {
				t.Error("create file fail")
			}
		}

		defer f.Close()
		defer os.Remove(f.Name())

		if gbfile.Exists(f.Name()) {
			cache := gbfile.GetContentsWithCache(f.Name())

			f, err = gbfile.OpenFile(f.Name(), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				t.Error("file open fail", err)
			}

			err = gbfile.PutContents(f.Name(), strTest)
			if err != nil {
				t.Error("write error", err)
			}

			t.Assert(cache, "")

			time.Sleep(100 * time.Millisecond)
		}
	})
}
