package gbspath_test

import (
	gbfile "ghostbb.io/gb/os/gb_file"
	gbspath "ghostbb.io/gb/os/gb_spath"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func TestSPath_Api(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		pwd := gbfile.Pwd()
		root := pwd
		file, _ := gbfile.Create(gbfile.Join(root, "gb_tmp", "gb.txt"))
		defer func() {
			t.AssertNil(file.Close())
			t.AssertNil(gbfile.Remove(gbfile.Join(root, "gb_tmp")))
		}()
		fp, isDir := gbspath.Search(root, "gb_tmp")
		t.Assert(fp, gbfile.Join(root, "gb_tmp"))
		t.Assert(isDir, true)
		fp, isDir = gbspath.Search(root, "gb_tmp", "gb.txt")
		t.Assert(fp, gbfile.Join(root, "gb_tmp", "gb.txt"))
		t.Assert(isDir, false)

		fp, isDir = gbspath.SearchWithCache(root, "gb_tmp")
		t.Assert(fp, gbfile.Join(root, "gb_tmp"))
		t.Assert(isDir, true)
		fp, isDir = gbspath.SearchWithCache(root, "gb_tmp", "gb.txt")
		t.Assert(fp, gbfile.Join(root, "gb_tmp", "gb.txt"))
		t.Assert(isDir, false)
	})
}

func TestSPath_Basic(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		pwd := gbfile.Pwd()
		root := pwd

		gbfile.Create(gbfile.Join(root, "gb_tmp", "gb.txt"))
		defer gbfile.Remove(gbfile.Join(root, "gb_tmp"))
		gsp := gbspath.New(root, false)
		realPath, err := gsp.Add(gbfile.Join(root, "gb_tmp"))
		t.AssertNil(err)
		t.Assert(realPath, gbfile.Join(root, "gb_tmp"))
		realPath, err = gsp.Add("gb_tmp1")
		t.Assert(err != nil, true)
		t.Assert(realPath, "")
		realPath, err = gsp.Add(gbfile.Join(root, "gb_tmp", "gb.txt"))
		t.Assert(err != nil, true)
		t.Assert(realPath, "")

		gsp.Remove("gb_tmp1")

		t.Assert(gsp.Size(), 2)
		t.Assert(len(gsp.Paths()), 2)
		t.Assert(len(gsp.AllPaths()), 0)
		realPath, err = gsp.Set(gbfile.Join(root, "gb_tmp1"))
		t.Assert(err != nil, true)
		t.Assert(realPath, "")
		realPath, err = gsp.Set(gbfile.Join(root, "gb_tmp", "gb.txt"))
		t.AssertNE(err, nil)
		t.Assert(realPath, "")

		realPath, err = gsp.Set(root)
		t.AssertNil(err)
		t.Assert(realPath, root)

		fp, isDir := gsp.Search("gb_tmp")
		t.Assert(fp, gbfile.Join(root, "gb_tmp"))
		t.Assert(isDir, true)
		fp, isDir = gsp.Search("gb_tmp", "gb.txt")
		t.Assert(fp, gbfile.Join(root, "gb_tmp", "gb.txt"))
		t.Assert(isDir, false)
		fp, isDir = gsp.Search("/", "gb.txt")
		t.Assert(fp, root)
		t.Assert(isDir, true)

		gsp = gbspath.New(root, true)
		realPath, err = gsp.Add(gbfile.Join(root, "gb_tmp"))
		t.AssertNil(err)
		t.Assert(realPath, gbfile.Join(root, "gb_tmp"))

		gbfile.Mkdir(gbfile.Join(root, "gb_tmp1"))
		gbfile.Rename(gbfile.Join(root, "gb_tmp1"), gbfile.Join(root, "gb_tmp2"))
		gbfile.Rename(gbfile.Join(root, "gb_tmp2"), gbfile.Join(root, "gb_tmp1"))
		defer gbfile.Remove(gbfile.Join(root, "gb_tmp1"))
		realPath, err = gsp.Add("gb_tmp1")
		t.Assert(err != nil, false)
		t.Assert(realPath, gbfile.Join(root, "gb_tmp1"))

		realPath, err = gsp.Add("gb_tmp3")
		t.Assert(err != nil, true)
		t.Assert(realPath, "")

		gsp.Remove(gbfile.Join(root, "gb_tmp"))
		gsp.Remove(gbfile.Join(root, "gb_tmp1"))
		gsp.Remove(gbfile.Join(root, "gb_tmp3"))
		t.Assert(gsp.Size(), 3)
		t.Assert(len(gsp.Paths()), 3)

		gsp.AllPaths()
		gsp.Set(root)
		fp, isDir = gsp.Search("gb_tmp")
		t.Assert(fp, gbfile.Join(root, "gb_tmp"))
		t.Assert(isDir, true)

		fp, isDir = gsp.Search("gb_tmp", "gb.txt")
		t.Assert(fp, gbfile.Join(root, "gb_tmp", "gb.txt"))
		t.Assert(isDir, false)
	})
}
