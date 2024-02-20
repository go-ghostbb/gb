package gbfile_test

import (
	gbdebug "ghostbb.io/gb/debug/gb_debug"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTestFile(filename, content string) error {
	TempDir := testpath()
	err := os.WriteFile(TempDir+filename, []byte(content), 0666)
	return err
}

func delTestFiles(filenames string) {
	os.RemoveAll(testpath() + filenames)
}

func createDir(paths string) {
	TempDir := testpath()
	os.Mkdir(TempDir+paths, 0777)
}

func formatpaths(paths []string) []string {
	for k, v := range paths {
		paths[k] = filepath.ToSlash(v)
		paths[k] = strings.Replace(paths[k], "./", "/", 1)
	}

	return paths
}

func formatpath(paths string) string {
	paths = filepath.ToSlash(paths)
	paths = strings.Replace(paths, "./", "/", 1)
	return paths
}

func testpath() string {
	return gbstr.TrimRight(os.TempDir(), "\\/")
}

func Test_GetContents(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {

		var (
			filepaths string = "/testfile_t1.txt"
		)
		createTestFile(filepaths, "my name is jroam")
		defer delTestFiles(filepaths)

		t.Assert(gbfile.GetContents(testpath()+filepaths), "my name is jroam")
		t.Assert(gbfile.GetContents(""), "")

	})
}

func Test_GetBinContents(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			filepaths1  = "/testfile_t1.txt"
			filepaths2  = testpath() + "/testfile_t1_no.txt"
			readcontent []byte
			str1        = "my name is jroam"
		)
		createTestFile(filepaths1, str1)
		defer delTestFiles(filepaths1)
		readcontent = gbfile.GetBytes(testpath() + filepaths1)
		t.Assert(readcontent, []byte(str1))

		readcontent = gbfile.GetBytes(filepaths2)
		t.Assert(string(readcontent), "")

		t.Assert(string(gbfile.GetBytes(filepaths2)), "")

	})
}

func Test_Truncate(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			filepaths1 = "/testfile_GetContentsyyui.txt"
			err        error
			files      *os.File
		)
		createTestFile(filepaths1, "abcdefghijkmln")
		defer delTestFiles(filepaths1)
		err = gbfile.Truncate(testpath()+filepaths1, 10)
		t.AssertNil(err)

		files, err = os.Open(testpath() + filepaths1)
		t.AssertNil(err)
		defer files.Close()
		fileinfo, err2 := files.Stat()
		t.Assert(err2, nil)
		t.Assert(fileinfo.Size(), 10)

		err = gbfile.Truncate("", 10)
		t.AssertNE(err, nil)

	})
}

func Test_PutContents(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			filepaths   = "/testfile_PutContents.txt"
			err         error
			readcontent []byte
		)
		createTestFile(filepaths, "a")
		defer delTestFiles(filepaths)

		err = gbfile.PutContents(testpath()+filepaths, "test!")
		t.AssertNil(err)

		readcontent, err = os.ReadFile(testpath() + filepaths)
		t.AssertNil(err)
		t.Assert(string(readcontent), "test!")

		err = gbfile.PutContents("", "test!")
		t.AssertNE(err, nil)

	})
}

func Test_PutContentsAppend(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			filepaths   = "/testfile_PutContents.txt"
			err         error
			readcontent []byte
		)

		createTestFile(filepaths, "a")
		defer delTestFiles(filepaths)
		err = gbfile.PutContentsAppend(testpath()+filepaths, "hello")
		t.AssertNil(err)

		readcontent, err = os.ReadFile(testpath() + filepaths)
		t.AssertNil(err)
		t.Assert(string(readcontent), "ahello")

		err = gbfile.PutContentsAppend("", "hello")
		t.AssertNE(err, nil)

	})

}

func Test_PutBinContents(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			filepaths   = "/testfile_PutContents.txt"
			err         error
			readcontent []byte
		)
		createTestFile(filepaths, "a")
		defer delTestFiles(filepaths)

		err = gbfile.PutBytes(testpath()+filepaths, []byte("test!!"))
		t.AssertNil(err)

		readcontent, err = os.ReadFile(testpath() + filepaths)
		t.AssertNil(err)
		t.Assert(string(readcontent), "test!!")

		err = gbfile.PutBytes("", []byte("test!!"))
		t.AssertNE(err, nil)

	})
}

func Test_PutBinContentsAppend(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			filepaths   = "/testfile_PutContents.txt"
			err         error
			readcontent []byte
		)
		createTestFile(filepaths, "test!!")
		defer delTestFiles(filepaths)
		err = gbfile.PutBytesAppend(testpath()+filepaths, []byte("word"))
		t.AssertNil(err)

		readcontent, err = os.ReadFile(testpath() + filepaths)
		t.AssertNil(err)
		t.Assert(string(readcontent), "test!!word")

		err = gbfile.PutBytesAppend("", []byte("word"))
		t.AssertNE(err, nil)

	})
}

func Test_GetBinContentsByTwoOffsetsByPath(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			filepaths   = "/testfile_GetContents.txt"
			readcontent []byte
		)

		createTestFile(filepaths, "abcdefghijk")
		defer delTestFiles(filepaths)
		readcontent = gbfile.GetBytesByTwoOffsetsByPath(testpath()+filepaths, 2, 5)

		t.Assert(string(readcontent), "cde")

		readcontent = gbfile.GetBytesByTwoOffsetsByPath("", 2, 5)
		t.Assert(len(readcontent), 0)

	})

}

func Test_GetNextCharOffsetByPath(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			filepaths  = "/testfile_GetContents.txt"
			localindex int64
		)
		createTestFile(filepaths, "abcdefghijk")
		defer delTestFiles(filepaths)
		localindex = gbfile.GetNextCharOffsetByPath(testpath()+filepaths, 'd', 1)
		t.Assert(localindex, 3)

		localindex = gbfile.GetNextCharOffsetByPath("", 'd', 1)
		t.Assert(localindex, -1)

	})
}

func Test_GetNextCharOffset(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			localindex int64
		)
		reader := strings.NewReader("helloword")

		localindex = gbfile.GetNextCharOffset(reader, 'w', 1)
		t.Assert(localindex, 5)

		localindex = gbfile.GetNextCharOffset(reader, 'j', 1)
		t.Assert(localindex, -1)

	})
}

func Test_GetBinContentsByTwoOffsets(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			reads []byte
		)
		reader := strings.NewReader("helloword")

		reads = gbfile.GetBytesByTwoOffsets(reader, 1, 3)
		t.Assert(string(reads), "el")

		reads = gbfile.GetBytesByTwoOffsets(reader, 10, 30)
		t.Assert(string(reads), "")

	})
}

func Test_GetBinContentsTilChar(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			reads  []byte
			indexs int64
		)
		reader := strings.NewReader("helloword")

		reads, _ = gbfile.GetBytesTilChar(reader, 'w', 2)
		t.Assert(string(reads), "llow")

		_, indexs = gbfile.GetBytesTilChar(reader, 'w', 20)
		t.Assert(indexs, -1)

	})
}

func Test_GetBinContentsTilCharByPath(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			reads     []byte
			indexs    int64
			filepaths = "/testfile_GetContents.txt"
		)

		createTestFile(filepaths, "abcdefghijklmn")
		defer delTestFiles(filepaths)

		reads, _ = gbfile.GetBytesTilCharByPath(testpath()+filepaths, 'c', 2)
		t.Assert(string(reads), "c")

		reads, _ = gbfile.GetBytesTilCharByPath(testpath()+filepaths, 'y', 1)
		t.Assert(string(reads), "")

		_, indexs = gbfile.GetBytesTilCharByPath(testpath()+filepaths, 'x', 1)
		t.Assert(indexs, -1)

	})
}

func Test_Home(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			reads string
			err   error
		)

		reads, err = gbfile.Home("a", "b")
		t.AssertNil(err)
		t.AssertNE(reads, "")
	})
}

func Test_NotFound(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		teatFile := gbfile.Dir(gbdebug.CallerFilePath()) + gbfile.Separator + "testdata/readline/error.log"
		callback := func(line string) error {
			return nil
		}
		err := gbfile.ReadLines(teatFile, callback)
		t.AssertNE(err, nil)
	})
}

func Test_ReadLines(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			expectList = []string{"a", "b", "c", "d", "e"}
			getList    = make([]string, 0)
			callback   = func(line string) error {
				getList = append(getList, line)
				return nil
			}
			teatFile = gbfile.Dir(gbdebug.CallerFilePath()) + gbfile.Separator + "testdata/readline/file.log"
		)
		err := gbfile.ReadLines(teatFile, callback)
		t.AssertEQ(getList, expectList)
		t.AssertEQ(err, nil)
	})
}

func Test_ReadLines_Error(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			callback = func(line string) error {
				return gberror.New("custom error")
			}
			teatFile = gbfile.Dir(gbdebug.CallerFilePath()) + gbfile.Separator + "testdata/readline/file.log"
		)
		err := gbfile.ReadLines(teatFile, callback)
		t.AssertEQ(err.Error(), "custom error")
	})
}

func Test_ReadLinesBytes(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			expectList = [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e")}
			getList    = make([][]byte, 0)
			callback   = func(line []byte) error {
				getList = append(getList, line)
				return nil
			}
			teatFile = gbfile.Dir(gbdebug.CallerFilePath()) + gbfile.Separator + "testdata/readline/file.log"
		)
		err := gbfile.ReadLinesBytes(teatFile, callback)
		t.AssertEQ(getList, expectList)
		t.AssertEQ(err, nil)
	})
}

func Test_ReadLinesBytes_Error(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			callback = func(line []byte) error {
				return gberror.New("custom error")
			}
			teatFile = gbfile.Dir(gbdebug.CallerFilePath()) + gbfile.Separator + "testdata/readline/file.log"
		)
		err := gbfile.ReadLinesBytes(teatFile, callback)
		t.AssertEQ(err.Error(), "custom error")
	})
}
