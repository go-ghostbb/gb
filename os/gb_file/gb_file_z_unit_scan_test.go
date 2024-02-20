package gbfile_test

import (
	gbarray "ghostbb.io/gb/container/gb_array"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func Test_ScanDir(t *testing.T) {
	teatPath := gbtest.DataPath()
	gbtest.C(t, func(t *gbtest.T) {
		files, err := gbfile.ScanDir(teatPath, "*", false)
		t.AssertNil(err)
		t.AssertIN(teatPath+gbfile.Separator+"dir1", files)
		t.AssertIN(teatPath+gbfile.Separator+"dir2", files)
		t.AssertNE(teatPath+gbfile.Separator+"dir1"+gbfile.Separator+"file1", files)
	})
	gbtest.C(t, func(t *gbtest.T) {
		files, err := gbfile.ScanDir(teatPath, "*", true)
		t.AssertNil(err)
		t.AssertIN(teatPath+gbfile.Separator+"dir1", files)
		t.AssertIN(teatPath+gbfile.Separator+"dir2", files)
		t.AssertIN(teatPath+gbfile.Separator+"dir1"+gbfile.Separator+"file1", files)
		t.AssertIN(teatPath+gbfile.Separator+"dir2"+gbfile.Separator+"file2", files)
	})
}

func Test_ScanDirFunc(t *testing.T) {
	teatPath := gbtest.DataPath()
	gbtest.C(t, func(t *gbtest.T) {
		files, err := gbfile.ScanDirFunc(teatPath, "*", true, func(path string) string {
			if gbfile.Name(path) != "file1" {
				return ""
			}
			return path
		})
		t.AssertNil(err)
		t.Assert(len(files), 1)
		t.Assert(gbfile.Name(files[0]), "file1")
	})
}

func Test_ScanDirFile(t *testing.T) {
	teatPath := gbtest.DataPath()
	gbtest.C(t, func(t *gbtest.T) {
		files, err := gbfile.ScanDirFile(teatPath, "*", false)
		t.AssertNil(err)
		t.Assert(len(files), 0)
	})
	gbtest.C(t, func(t *gbtest.T) {
		files, err := gbfile.ScanDirFile(teatPath, "*", true)
		t.AssertNil(err)
		t.AssertNI(teatPath+gbfile.Separator+"dir1", files)
		t.AssertNI(teatPath+gbfile.Separator+"dir2", files)
		t.AssertIN(teatPath+gbfile.Separator+"dir1"+gbfile.Separator+"file1", files)
		t.AssertIN(teatPath+gbfile.Separator+"dir2"+gbfile.Separator+"file2", files)
	})
}

func Test_ScanDirFileFunc(t *testing.T) {
	teatPath := gbtest.DataPath()
	gbtest.C(t, func(t *gbtest.T) {
		array := gbarray.New()
		files, err := gbfile.ScanDirFileFunc(teatPath, "*", false, func(path string) string {
			array.Append(1)
			return path
		})
		t.AssertNil(err)
		t.Assert(len(files), 0)
		t.Assert(array.Len(), 0)
	})
	gbtest.C(t, func(t *gbtest.T) {
		array := gbarray.New()
		files, err := gbfile.ScanDirFileFunc(teatPath, "*", true, func(path string) string {
			array.Append(1)
			if gbfile.Basename(path) == "file1" {
				return path
			}
			return ""
		})
		t.AssertNil(err)
		t.Assert(len(files), 1)
		t.Assert(array.Len(), 3)
	})
}
