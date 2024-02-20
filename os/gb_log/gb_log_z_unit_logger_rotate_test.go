package gblog_test

import (
	"context"
	"ghostbb.io/gb/frame/g"
	gbfile "ghostbb.io/gb/os/gb_file"
	gblog "ghostbb.io/gb/os/gb_log"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"testing"
	"time"
)

var (
	ctx = context.TODO()
)

func Test_Rotate_Size(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		p := gbfile.Temp(gbtime.TimestampNanoStr())
		err := l.SetConfigWithMap(g.Map{
			"Path":                 p,
			"File":                 "access.log",
			"StdoutPrint":          false,
			"RotateSize":           10,
			"RotateBackupLimit":    2,
			"RotateBackupExpire":   5 * time.Second,
			"RotateBackupCompress": 9,
			"RotateCheckInterval":  time.Second, // For unit testing only.
		})
		t.AssertNil(err)
		defer gbfile.Remove(p)

		s := "1234567890abcdefg"
		for i := 0; i < 8; i++ {
			l.Print(ctx, s)
			time.Sleep(time.Second)
		}

		logFiles, err := gbfile.ScanDirFile(p, "access*")
		t.AssertNil(err)

		for _, v := range logFiles {
			content := gbfile.GetContents(v)
			t.AssertIN(gbstr.Count(content, s), []int{1, 2})
		}

		time.Sleep(time.Second * 3)

		files, err := gbfile.ScanDirFile(p, "*.gz")
		t.AssertNil(err)
		t.Assert(len(files), 2)

		time.Sleep(time.Second * 5)
		files, err = gbfile.ScanDirFile(p, "*.gz")
		t.AssertNil(err)
		t.Assert(len(files), 0)
	})
}

func Test_Rotate_Expire(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		l := gblog.New()
		p := gbfile.Temp(gbtime.TimestampNanoStr())
		err := l.SetConfigWithMap(g.Map{
			"Path":                 p,
			"File":                 "access.log",
			"StdoutPrint":          false,
			"RotateExpire":         time.Second / 2,
			"RotateBackupLimit":    2,
			"RotateBackupExpire":   5 * time.Second,
			"RotateBackupCompress": 9,
			"RotateCheckInterval":  time.Second, // For unit testing only.
		})
		t.AssertNil(err)
		defer gbfile.Remove(p)
		s := "1234567890abcdefg"
		for i := 0; i < 10; i++ {
			l.Print(ctx, s)
		}

		files, err := gbfile.ScanDirFile(p, "*.gz")
		t.AssertNil(err)
		t.Assert(len(files), 0)

		t.Assert(gbstr.Count(gbfile.GetContents(gbfile.Join(p, "access.log")), s), 10)

		time.Sleep(time.Second * 3)

		filenames, err := gbfile.ScanDirFile(p, "*")
		t.Log(filenames, err)
		files, err = gbfile.ScanDirFile(p, "*.gz")
		t.AssertNil(err)
		t.Assert(len(files), 1)

		t.Assert(gbstr.Count(gbfile.GetContents(gbfile.Join(p, "access.log")), s), 0)

		time.Sleep(time.Second * 5)
		files, err = gbfile.ScanDirFile(p, "*.gz")
		t.AssertNil(err)
		t.Assert(len(files), 0)
	})
}
