package gbfile_test

import (
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbuid "ghostbb.io/gb/util/gb_uid"
	"os"
	"testing"
)

func Test_Copy(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			paths  = "/testfile_copyfile1.txt"
			topath = "/testfile_copyfile2.txt"
		)

		createTestFile(paths, "")
		defer delTestFiles(paths)

		t.Assert(gbfile.Copy(testpath()+paths, testpath()+topath), nil)
		defer delTestFiles(topath)

		t.Assert(gbfile.IsFile(testpath()+topath), true)
		t.AssertNE(gbfile.Copy(paths, ""), nil)
		t.AssertNE(gbfile.Copy("", topath), nil)
	})
}

func Test_Copy_File_To_Dir(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			src = gbtest.DataPath("dir1", "file1")
			dst = gbfile.Temp(gbuid.S(), "dir2")
		)
		err := gbfile.Mkdir(dst)
		t.AssertNil(err)
		defer gbfile.Remove(dst)

		err = gbfile.Copy(src, dst)
		t.AssertNil(err)

		expectPath := gbfile.Join(dst, "file1")
		t.Assert(gbfile.GetContents(expectPath), gbfile.GetContents(src))
	})
}

func Test_Copy_Dir_To_File(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			src = gbtest.DataPath("dir1")
			dst = gbfile.Temp(gbuid.S(), "file2")
		)
		f, err := gbfile.Create(dst)
		t.AssertNil(err)
		defer f.Close()
		defer gbfile.Remove(dst)

		err = gbfile.Copy(src, dst)
		t.AssertNE(err, nil)
	})
}

func Test_CopyFile(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			paths  = "/testfile_copyfile1.txt"
			topath = "/testfile_copyfile2.txt"
		)

		createTestFile(paths, "")
		defer delTestFiles(paths)

		t.Assert(gbfile.CopyFile(testpath()+paths, testpath()+topath), nil)
		defer delTestFiles(topath)

		t.Assert(gbfile.IsFile(testpath()+topath), true)
		t.AssertNE(gbfile.CopyFile(paths, ""), nil)
		t.AssertNE(gbfile.CopyFile("", topath), nil)
	})
	// Content replacement.
	gbtest.C(t, func(t *gbtest.T) {
		src := gbfile.Temp(gbtime.TimestampNanoStr())
		dst := gbfile.Temp(gbtime.TimestampNanoStr())
		srcContent := "1"
		dstContent := "1"
		t.Assert(gbfile.PutContents(src, srcContent), nil)
		t.Assert(gbfile.PutContents(dst, dstContent), nil)
		t.Assert(gbfile.GetContents(src), srcContent)
		t.Assert(gbfile.GetContents(dst), dstContent)

		t.Assert(gbfile.CopyFile(src, dst), nil)
		t.Assert(gbfile.GetContents(src), srcContent)
		t.Assert(gbfile.GetContents(dst), srcContent)
	})
	// Set mode
	gbtest.C(t, func(t *gbtest.T) {
		var (
			src     = "/testfile_copyfile1.txt"
			dst     = "/testfile_copyfile2.txt"
			dstMode = os.FileMode(0600)
		)
		t.AssertNil(createTestFile(src, ""))
		defer delTestFiles(src)

		t.Assert(gbfile.CopyFile(testpath()+src, testpath()+dst, gbfile.CopyOption{Mode: dstMode}), nil)
		defer delTestFiles(dst)

		dstStat, err := gbfile.Stat(testpath() + dst)
		t.AssertNil(err)
		t.Assert(dstStat.Mode().Perm(), dstMode)
	})
	// Preserve src file's mode
	gbtest.C(t, func(t *gbtest.T) {
		var (
			src = "/testfile_copyfile1.txt"
			dst = "/testfile_copyfile2.txt"
		)
		t.AssertNil(createTestFile(src, ""))
		defer delTestFiles(src)

		t.Assert(gbfile.CopyFile(testpath()+src, testpath()+dst, gbfile.CopyOption{PreserveMode: true}), nil)
		defer delTestFiles(dst)

		srcStat, err := gbfile.Stat(testpath() + src)
		t.AssertNil(err)
		dstStat, err := gbfile.Stat(testpath() + dst)
		t.AssertNil(err)
		t.Assert(srcStat.Mode().Perm(), dstStat.Mode().Perm())
	})
}

func Test_CopyDir(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			dirPath1 = "/test-copy-dir1"
			dirPath2 = "/test-copy-dir2"
		)
		haveList := []string{
			"t1.txt",
			"t2.txt",
		}
		createDir(dirPath1)
		for _, v := range haveList {
			t.Assert(createTestFile(dirPath1+"/"+v, ""), nil)
		}
		defer delTestFiles(dirPath1)

		var (
			yfolder  = testpath() + dirPath1
			tofolder = testpath() + dirPath2
		)

		if gbfile.IsDir(tofolder) {
			t.Assert(gbfile.Remove(tofolder), nil)
			t.Assert(gbfile.Remove(""), nil)
		}

		t.Assert(gbfile.CopyDir(yfolder, tofolder), nil)
		defer delTestFiles(tofolder)

		t.Assert(gbfile.IsDir(yfolder), true)

		for _, v := range haveList {
			t.Assert(gbfile.IsFile(yfolder+"/"+v), true)
		}

		t.Assert(gbfile.IsDir(tofolder), true)

		for _, v := range haveList {
			t.Assert(gbfile.IsFile(tofolder+"/"+v), true)
		}

		t.Assert(gbfile.Remove(tofolder), nil)
		t.Assert(gbfile.Remove(""), nil)
	})
	// Content replacement.
	gbtest.C(t, func(t *gbtest.T) {
		src := gbfile.Temp(gbtime.TimestampNanoStr(), gbtime.TimestampNanoStr())
		dst := gbfile.Temp(gbtime.TimestampNanoStr(), gbtime.TimestampNanoStr())
		defer func() {
			gbfile.Remove(src)
			gbfile.Remove(dst)
		}()
		srcContent := "1"
		dstContent := "1"
		t.Assert(gbfile.PutContents(src, srcContent), nil)
		t.Assert(gbfile.PutContents(dst, dstContent), nil)
		t.Assert(gbfile.GetContents(src), srcContent)
		t.Assert(gbfile.GetContents(dst), dstContent)

		err := gbfile.CopyDir(gbfile.Dir(src), gbfile.Dir(dst))
		t.AssertNil(err)
		t.Assert(gbfile.GetContents(src), srcContent)
		t.Assert(gbfile.GetContents(dst), srcContent)

		t.AssertNE(gbfile.CopyDir(gbfile.Dir(src), ""), nil)
		t.AssertNE(gbfile.CopyDir("", gbfile.Dir(dst)), nil)
	})
}
