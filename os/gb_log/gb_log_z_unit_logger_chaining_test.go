package gblog

import (
	"bytes"
	"fmt"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"testing"
	"time"
)

func Test_To(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		w := bytes.NewBuffer(nil)
		To(w).Error(ctx, 1, 2, 3)
		To(w).Errorf(ctx, "%d %d %d", 1, 2, 3)
		t.Assert(gbstr.Count(w.String(), defaultLevelPrefixes[LEVEL_ERRO]), 2)
		t.Assert(gbstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Path(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gbtime.TimestampNano())

		err := gbfile.Mkdir(path)
		t.AssertNil(err)
		defer gbfile.Remove(path)

		Path(path).File(file).Stdout(false).Error(ctx, 1, 2, 3)
		Path(path).File(file).Stdout(false).Errorf(ctx, "%d %d %d", 1, 2, 3)
		content := gbfile.GetContents(gbfile.Join(path, file))
		t.Assert(gbstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 2)
		t.Assert(gbstr.Count(content, "1 2 3"), 2)
	})
}

func Test_Cat(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		cat := "category"
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gbtime.TimestampNano())

		err := gbfile.Mkdir(path)
		t.AssertNil(err)
		defer gbfile.Remove(path)

		Path(path).File(file).Cat(cat).Stdout(false).Error(ctx, 1, 2, 3)
		Path(path).File(file).Cat(cat).Stdout(false).Errorf(ctx, "%d %d %d", 1, 2, 3)
		content := gbfile.GetContents(gbfile.Join(path, cat, file))
		t.Assert(gbstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 2)
		t.Assert(gbstr.Count(content, "1 2 3"), 2)
	})
}

func Test_Level(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gbtime.TimestampNano())

		err := gbfile.Mkdir(path)
		t.AssertNil(err)
		defer gbfile.Remove(path)

		Path(path).File(file).Level(LEVEL_PROD).Stdout(false).Debug(ctx, 1, 2, 3)
		Path(path).File(file).Level(LEVEL_PROD).Stdout(false).Debug(ctx, "%d %d %d", 1, 2, 3)
		content := gbfile.GetContents(gbfile.Join(path, file))
		t.Assert(gbstr.Count(content, defaultLevelPrefixes[LEVEL_DEBU]), 0)
		t.Assert(gbstr.Count(content, "1 2 3"), 0)
	})
}

func Test_Skip(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gbtime.TimestampNano())

		err := gbfile.Mkdir(path)
		t.AssertNil(err)
		defer gbfile.Remove(path)

		Path(path).File(file).Skip(10).Stdout(false).Error(ctx, 1, 2, 3)
		Path(path).File(file).Stdout(false).Errorf(ctx, "%d %d %d", 1, 2, 3)
		content := gbfile.GetContents(gbfile.Join(path, file))
		fmt.Println(content)
		t.Assert(gbstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 2)
		t.Assert(gbstr.Count(content, "1 2 3"), 2)
		//t.Assert(gbstr.Count(content, "Stack"), 1)
	})
}

func Test_Stack(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gbtime.TimestampNano())

		err := gbfile.Mkdir(path)
		t.AssertNil(err)
		defer gbfile.Remove(path)

		Path(path).File(file).Stack(false).Stdout(false).Error(ctx, 1, 2, 3)
		Path(path).File(file).Stdout(false).Errorf(ctx, "%d %d %d", 1, 2, 3)
		content := gbfile.GetContents(gbfile.Join(path, file))
		fmt.Println(content)
		t.Assert(gbstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 2)
		t.Assert(gbstr.Count(content, "1 2 3"), 2)
		//t.Assert(gbstr.Count(content, "Stack"), 1)
	})
}

func Test_StackWithFilter(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gbtime.TimestampNano())

		err := gbfile.Mkdir(path)
		t.AssertNil(err)
		defer gbfile.Remove(path)

		Path(path).File(file).StackWithFilter("none").Stdout(false).Error(ctx, 1, 2, 3)
		content := gbfile.GetContents(gbfile.Join(path, file))
		fmt.Println(ctx, content)
		t.Assert(gbstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 1)
		t.Assert(gbstr.Count(content, "1 2 3"), 1)
		//t.Assert(gbstr.Count(content, "Stack"), 1)

	})
	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gbtime.TimestampNano())

		err := gbfile.Mkdir(path)
		t.AssertNil(err)
		defer gbfile.Remove(path)

		Path(path).File(file).StackWithFilter("/gf/").Stdout(false).Error(ctx, 1, 2, 3)
		content := gbfile.GetContents(gbfile.Join(path, file))
		fmt.Println(ctx, content)
		t.Assert(gbstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 1)
		t.Assert(gbstr.Count(content, "1 2 3"), 1)
		//t.Assert(gbstr.Count(content, "Stack"), 0)
	})
}

func Test_Header(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gbtime.TimestampNano())

		err := gbfile.Mkdir(path)
		t.AssertNil(err)
		defer gbfile.Remove(path)

		Path(path).File(file).Header(true).Stdout(false).Error(ctx, 1, 2, 3)
		content := gbfile.GetContents(gbfile.Join(path, file))
		t.Assert(gbstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 1)
		t.Assert(gbstr.Count(content, "1 2 3"), 1)
	})
	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gbtime.TimestampNano())

		err := gbfile.Mkdir(path)
		t.AssertNil(err)
		defer gbfile.Remove(path)

		Path(path).File(file).Header(false).Stdout(false).Error(ctx, 1, 2, 3)
		content := gbfile.GetContents(gbfile.Join(path, file))
		t.Assert(gbstr.Count(content, defaultLevelPrefixes[LEVEL_ERRO]), 0)
		t.Assert(gbstr.Count(content, "1 2 3"), 1)
	})
}

func Test_Line(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gbtime.TimestampNano())

		err := gbfile.Mkdir(path)
		t.AssertNil(err)
		defer gbfile.Remove(path)

		Path(path).File(file).Line(true).Stdout(false).Debug(ctx, 1, 2, 3)
		content := gbfile.GetContents(gbfile.Join(path, file))
		fmt.Println(content)
		t.Assert(gbstr.Count(content, defaultLevelPrefixes[LEVEL_DEBU]), 1)
		t.Assert(gbstr.Count(content, "1 2 3"), 1)
		//t.Assert(gbstr.Count(content, ".go"), 1)
		//t.Assert(gbstr.Contains(content, gbfile.Separator), true)
	})
	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gbtime.TimestampNano())

		err := gbfile.Mkdir(path)
		t.AssertNil(err)
		defer gbfile.Remove(path)

		Path(path).File(file).Line(false).Stdout(false).Debug(ctx, 1, 2, 3)
		content := gbfile.GetContents(gbfile.Join(path, file))
		t.Assert(gbstr.Count(content, defaultLevelPrefixes[LEVEL_DEBU]), 1)
		t.Assert(gbstr.Count(content, "1 2 3"), 1)
		//t.Assert(gbstr.Count(content, ".go"), 1)
		//t.Assert(gbstr.Contains(content, gbfile.Separator), false)
	})
}

func Test_Async(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gbtime.TimestampNano())

		err := gbfile.Mkdir(path)
		t.AssertNil(err)
		defer gbfile.Remove(path)

		Path(path).File(file).Async().Stdout(false).Debug(ctx, 1, 2, 3)
		time.Sleep(1000 * time.Millisecond)

		content := gbfile.GetContents(gbfile.Join(path, file))
		t.Assert(gbstr.Count(content, defaultLevelPrefixes[LEVEL_DEBU]), 1)
		t.Assert(gbstr.Count(content, "1 2 3"), 1)
	})

	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, gbtime.TimestampNano())

		err := gbfile.Mkdir(path)
		t.AssertNil(err)
		defer gbfile.Remove(path)

		Path(path).File(file).Async(false).Stdout(false).Debug(ctx, 1, 2, 3)
		content := gbfile.GetContents(gbfile.Join(path, file))
		t.Assert(gbstr.Count(content, defaultLevelPrefixes[LEVEL_DEBU]), 1)
		t.Assert(gbstr.Count(content, "1 2 3"), 1)
	})
}
