package gbfile_test

import (
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtest "ghostbb.io/gb/test/gb_test"
	"path/filepath"
	"testing"
)

func Test_Search(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			paths1  string = "/testfiless"
			paths2  string = "./testfile/dirfiles_no"
			tpath   string
			tpath2  string
			tempstr string
			ypaths1 string
			err     error
		)

		createDir(paths1)
		defer delTestFiles(paths1)
		ypaths1 = paths1

		tpath, err = gbfile.Search(testpath() + paths1)
		t.AssertNil(err)

		tpath = filepath.ToSlash(tpath)

		// 自定义优先路径
		tpath2, err = gbfile.Search(testpath() + paths1)
		t.AssertNil(err)
		tpath2 = filepath.ToSlash(tpath2)

		tempstr = testpath()
		paths1 = tempstr + paths1
		paths1 = filepath.ToSlash(paths1)

		t.Assert(tpath, paths1)

		t.Assert(tpath2, tpath)

		// 测试给定目录
		tpath2, err = gbfile.Search(paths1, "testfiless")
		tpath2 = filepath.ToSlash(tpath2)
		tempss := filepath.ToSlash(paths1)
		t.Assert(tpath2, tempss)

		// 测试当前目录
		tempstr, _ = filepath.Abs("./")
		tempstr = testpath()
		paths1 = tempstr + ypaths1
		paths1 = filepath.ToSlash(paths1)

		t.Assert(tpath2, paths1)

		// 测试目录不存在时
		_, err = gbfile.Search(paths2)
		t.AssertNE(err, nil)

	})
}
