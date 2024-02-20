package gbres_test

import (
	_ "ghostbb.io/gb/os/gb_res/testdata/data"

	gbfile "ghostbb.io/gb/os/gb_file"
	gbres "ghostbb.io/gb/os/gb_res"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"

	"ghostbb.io/gb/frame/g"
	"strings"
	"testing"
)

func Test_PackFolderToGoFile(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			srcPath    = gbtest.DataPath("files")
			goFilePath = gbfile.Temp(gbtime.TimestampNanoStr(), "testdata.go")
			pkgName    = "testdata"
			err        = gbres.PackToGoFile(srcPath, goFilePath, pkgName)
		)
		t.AssertNil(err)
		_ = gbfile.Remove(goFilePath)
	})
}

func Test_PackMultiFilesToGoFile(t *testing.T) {
	gbres.Dump()
	gbtest.C(t, func(t *gbtest.T) {
		var (
			srcPath    = gbtest.DataPath("files")
			goFilePath = gbfile.Temp(gbtime.TimestampNanoStr(), "data.go")
			pkgName    = "data"
			array, err = gbfile.ScanDir(srcPath, "*", false)
		)
		t.AssertNil(err)
		err = gbres.PackToGoFile(strings.Join(array, ","), goFilePath, pkgName)
		t.AssertNil(err)
		defer func() {
			t.AssertNil(gbfile.Remove(goFilePath))
		}()

		t.AssertNil(gbfile.CopyFile(goFilePath, gbtest.DataPath("data/data.go")))
	})
}

func Test_Pack(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			srcPath   = gbtest.DataPath("files")
			data, err = gbres.Pack(srcPath)
		)
		t.AssertNil(err)

		r := gbres.New()
		err = r.Add(string(data))
		t.AssertNil(err)
		t.Assert(r.Contains("files/"), true)
	})
}

func Test_PackToFile(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			srcPath = gbtest.DataPath("files")
			dstPath = gbfile.Temp(gbtime.TimestampNanoStr())
			err     = gbres.PackToFile(srcPath, dstPath)
		)
		t.AssertNil(err)

		defer func() {
			t.AssertNil(gbfile.Remove(dstPath))
		}()

		r := gbres.New()
		err = r.Load(dstPath)
		t.AssertNil(err)
		t.Assert(r.Contains("files"), true)
	})
}

func Test_PackWithPrefix1(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			srcPath    = gbtest.DataPath("files")
			goFilePath = gbfile.Temp(gbtime.TimestampNanoStr(), "testdata.go")
			pkgName    = "testdata"
			err        = gbres.PackToGoFile(srcPath, goFilePath, pkgName, "www/gf-site/test")
		)
		t.AssertNil(err)
		_ = gbfile.Remove(goFilePath)
	})
}

func Test_PackWithPrefix2(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			srcPath    = gbtest.DataPath("files")
			goFilePath = gbfile.Temp(gbtime.TimestampNanoStr(), "testdata.go")
			pkgName    = "testdata"
			err        = gbres.PackToGoFile(srcPath, goFilePath, pkgName, "/var/www/gf-site/test")
		)
		t.AssertNil(err)
		_ = gbfile.Remove(goFilePath)
	})
}

func Test_Unpack(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			srcPath    = gbtest.DataPath("testdata.txt")
			files, err = gbres.Unpack(srcPath)
		)
		t.AssertNil(err)
		t.Assert(len(files), 63)
	})
}

func Test_Basic(t *testing.T) {
	gbres.Dump()
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbres.Get("none"), nil)
		t.Assert(gbres.Contains("none"), false)
		t.Assert(gbres.Contains("dir1"), true)
	})

	gbtest.C(t, func(t *gbtest.T) {
		path := "dir1/test1"
		file := gbres.Get(path)
		t.AssertNE(file, nil)
		t.Assert(file.Name(), path)

		info := file.FileInfo()
		t.AssertNE(info, nil)
		t.Assert(info.IsDir(), false)
		t.Assert(info.Name(), "test1")

		rc, err := file.Open()
		t.AssertNil(err)
		defer rc.Close()

		b := make([]byte, 5)
		n, err := rc.Read(b)
		t.Assert(n, 5)
		t.AssertNil(err)
		t.Assert(string(b), "test1")

		t.Assert(file.Content(), "test1 content")
	})

	gbtest.C(t, func(t *gbtest.T) {
		path := "dir2"
		file := gbres.Get(path)
		t.AssertNE(file, nil)
		t.Assert(file.Name(), path)

		info := file.FileInfo()
		t.AssertNE(info, nil)
		t.Assert(info.IsDir(), true)
		t.Assert(info.Name(), "dir2")

		rc, err := file.Open()
		t.AssertNil(err)
		defer rc.Close()

		t.Assert(file.Content(), nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		path := "dir2/test2"
		file := gbres.Get(path)
		t.AssertNE(file, nil)
		t.Assert(file.Name(), path)
		t.Assert(file.Content(), "test2 content")
	})
}

func Test_Get(t *testing.T) {
	gbres.Dump()
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertNE(gbres.Get("dir1/test1"), nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		file := gbres.GetWithIndex("dir1", g.SliceStr{"test1"})
		t.AssertNE(file, nil)
		t.Assert(file.Name(), "dir1/test1")
	})
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbres.GetContent("dir1"), "")
		t.Assert(gbres.GetContent("dir1/test1"), "test1 content")
	})
}

func Test_ScanDir(t *testing.T) {
	gbres.Dump()
	gbtest.C(t, func(t *gbtest.T) {
		path := "dir1"
		files := gbres.ScanDir(path, "*", false)
		t.AssertNE(files, nil)
		t.Assert(len(files), 2)
	})
	gbtest.C(t, func(t *gbtest.T) {
		path := "dir1"
		files := gbres.ScanDir(path, "*", true)
		t.AssertNE(files, nil)
		t.Assert(len(files), 3)
	})

	gbtest.C(t, func(t *gbtest.T) {
		path := "dir1"
		files := gbres.ScanDir(path, "*.*", true)
		t.AssertNE(files, nil)
		t.Assert(len(files), 1)
		t.Assert(files[0].Name(), "dir1/sub/sub-test1.txt")
		t.Assert(files[0].Content(), "sub-test1 content")
	})
}

func Test_ScanDirFile(t *testing.T) {
	gbres.Dump()
	gbtest.C(t, func(t *gbtest.T) {
		path := "dir2"
		files := gbres.ScanDirFile(path, "*", false)
		t.AssertNE(files, nil)
		t.Assert(len(files), 1)
	})
	gbtest.C(t, func(t *gbtest.T) {
		path := "dir2"
		files := gbres.ScanDirFile(path, "*", true)
		t.AssertNE(files, nil)
		t.Assert(len(files), 2)
	})

	gbtest.C(t, func(t *gbtest.T) {
		path := "dir2"
		files := gbres.ScanDirFile(path, "*.*", true)
		t.AssertNE(files, nil)
		t.Assert(len(files), 1)
		t.Assert(files[0].Name(), "dir2/sub/sub-test2.txt")
		t.Assert(files[0].Content(), "sub-test2 content")
	})
}

func Test_Export(t *testing.T) {
	gbres.Dump()
	gbtest.C(t, func(t *gbtest.T) {
		var (
			src = `template-res`
			dst = gbfile.Temp(gbtime.TimestampNanoStr())
			err = gbres.Export(src, dst)
		)
		defer gbfile.Remove(dst)
		t.AssertNil(err)
		files, err := gbfile.ScanDir(dst, "*", true)
		t.AssertNil(err)
		t.Assert(len(files), 14)

		name := `template-res/index.html`
		t.Assert(gbfile.GetContents(gbfile.Join(dst, name)), gbres.GetContent(name))
	})
	gbtest.C(t, func(t *gbtest.T) {
		var (
			src = `template-res`
			dst = gbfile.Temp(gbtime.TimestampNanoStr())
			err = gbres.Export(src, dst, gbres.ExportOption{
				RemovePrefix: `template-res`,
			})
		)
		defer gbfile.Remove(dst)
		t.AssertNil(err)
		files, err := gbfile.ScanDir(dst, "*", true)
		t.AssertNil(err)
		t.Assert(len(files), 13)

		nameInRes := `template-res/index.html`
		nameInSys := `index.html`
		t.Assert(gbfile.GetContents(gbfile.Join(dst, nameInSys)), gbres.GetContent(nameInRes))
	})
	gbtest.C(t, func(t *gbtest.T) {
		var (
			src = `template-res/layout1/container.html`
			dst = gbfile.Temp(gbtime.TimestampNanoStr())
			err = gbres.Export(src, dst, gbres.ExportOption{
				RemovePrefix: `template-res`,
			})
		)
		defer gbfile.Remove(dst)
		t.AssertNil(err)
		files, err := gbfile.ScanDir(dst, "*", true)
		t.AssertNil(err)
		t.Assert(len(files), 2)
	})
}

func Test_IsEmpty(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbres.IsEmpty(), false)
	})
}

func TestFile_Name(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			src = `template-res`
		)
		t.Assert(gbres.Get(src).Name(), src)
	})
}

func TestFile_Export(t *testing.T) {
	gbres.Dump()
	gbtest.C(t, func(t *gbtest.T) {
		var (
			src = `template-res`
			dst = gbfile.Temp(gbtime.TimestampNanoStr())
			err = gbres.Get(src).Export(dst)
		)
		defer gbfile.Remove(dst)
		t.AssertNil(err)
		files, err := gbfile.ScanDir(dst, "*", true)
		t.AssertNil(err)
		t.Assert(len(files), 14)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var (
			src = `template-res`
			dst = gbfile.Temp(gbtime.TimestampNanoStr())
			err = gbres.Get(src).Export(dst, gbres.ExportOption{
				RemovePrefix: `template-res`,
			})
		)
		defer gbfile.Remove(dst)
		t.AssertNil(err)
		files, err := gbfile.ScanDir(dst, "*", true)
		t.AssertNil(err)
		t.Assert(len(files), 13)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var (
			src = `template-res/layout1/container.html`
			dst = gbfile.Temp(gbtime.TimestampNanoStr())
			err = gbres.Get(src).Export(dst, gbres.ExportOption{
				RemovePrefix: `template-res`,
			})
		)
		defer gbfile.Remove(dst)
		t.AssertNil(err)
		files, err := gbfile.ScanDir(dst, "*", true)
		t.AssertNil(err)
		t.Assert(len(files), 2)
	})
}
