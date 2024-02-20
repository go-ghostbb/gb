package gbfile_test

import (
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtest "ghostbb.io/gb/test/gb_test"
	"os"
	"testing"
	"time"
)

func Test_MTime(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {

		var (
			file1   = "/testfile_t1.txt"
			err     error
			fileobj os.FileInfo
		)

		createTestFile(file1, "")
		defer delTestFiles(file1)
		fileobj, err = os.Stat(testpath() + file1)
		t.AssertNil(err)

		t.Assert(gbfile.MTime(testpath()+file1), fileobj.ModTime())
		t.Assert(gbfile.MTime(""), "")
	})
}

func Test_MTimeMillisecond(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			file1   = "/testfile_t1.txt"
			err     error
			fileobj os.FileInfo
		)

		createTestFile(file1, "")
		defer delTestFiles(file1)
		fileobj, err = os.Stat(testpath() + file1)
		t.AssertNil(err)

		time.Sleep(time.Millisecond * 100)
		t.AssertGE(
			gbfile.MTimestampMilli(testpath()+file1),
			fileobj.ModTime().UnixNano()/1000000,
		)
		t.Assert(gbfile.MTimestampMilli(""), -1)
	})
}
