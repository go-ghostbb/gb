package gbfsnotify_test

import (
	gbarray "ghostbb.io/gb/container/gb_array"
	gbtype "ghostbb.io/gb/container/gb_type"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbfsnotify "ghostbb.io/gb/os/gb_fsnotify"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
	"time"
)

func TestWatcher_AddOnce(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		value := gbtype.New()
		path := gbfile.Temp(gbconv.String(gbtime.TimestampNano()))
		err := gbfile.PutContents(path, "init")
		t.AssertNil(err)
		defer gbfile.Remove(path)

		time.Sleep(100 * time.Millisecond)
		callback1, err := gbfsnotify.AddOnce("mywatch", path, func(event *gbfsnotify.Event) {
			value.Set(1)
		})
		t.AssertNil(err)
		callback2, err := gbfsnotify.AddOnce("mywatch", path, func(event *gbfsnotify.Event) {
			value.Set(2)
		})
		t.AssertNil(err)
		t.Assert(callback2, nil)

		err = gbfile.PutContents(path, "1")
		t.AssertNil(err)

		time.Sleep(100 * time.Millisecond)
		t.Assert(value, 1)

		err = gbfsnotify.RemoveCallback(callback1.Id)
		t.AssertNil(err)

		err = gbfile.PutContents(path, "2")
		t.AssertNil(err)

		time.Sleep(100 * time.Millisecond)
		t.Assert(value, 1)
	})
}

func TestWatcher_AddRemove(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		path1 := gbfile.Temp() + gbfile.Separator + gbconv.String(gbtime.TimestampNano())
		path2 := gbfile.Temp() + gbfile.Separator + gbconv.String(gbtime.TimestampNano()) + "2"
		gbfile.PutContents(path1, "1")
		defer func() {
			gbfile.Remove(path1)
			gbfile.Remove(path2)
		}()
		v := gbtype.NewInt(1)
		callback, err := gbfsnotify.Add(path1, func(event *gbfsnotify.Event) {
			if event.IsWrite() {
				v.Set(2)
				return
			}
			if event.IsRename() {
				v.Set(3)
				gbfsnotify.Exit()
				return
			}
		})
		t.AssertNil(err)
		t.AssertNE(callback, nil)

		gbfile.PutContents(path1, "2")
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 2)

		gbfile.Rename(path1, path2)
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 3)
	})

	gbtest.C(t, func(t *gbtest.T) {
		path1 := gbfile.Temp() + gbfile.Separator + gbconv.String(gbtime.TimestampNano())
		gbfile.PutContents(path1, "1")
		defer func() {
			gbfile.Remove(path1)
		}()
		v := gbtype.NewInt(1)
		callback, err := gbfsnotify.Add(path1, func(event *gbfsnotify.Event) {
			if event.IsWrite() {
				v.Set(2)
				return
			}
			if event.IsRemove() {
				v.Set(4)
				return
			}
		})
		t.AssertNil(err)
		t.AssertNE(callback, nil)

		gbfile.PutContents(path1, "2")
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 2)

		gbfile.Remove(path1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 4)

		gbfile.PutContents(path1, "1")
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 4)
	})
}

func TestWatcher_Callback1(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		path1 := gbfile.Temp(gbtime.TimestampNanoStr())
		gbfile.PutContents(path1, "1")
		defer func() {
			gbfile.Remove(path1)
		}()
		v := gbtype.NewInt(1)
		callback, err := gbfsnotify.Add(path1, func(event *gbfsnotify.Event) {
			if event.IsWrite() {
				v.Set(2)
				return
			}
		})
		t.AssertNil(err)
		t.AssertNE(callback, nil)

		gbfile.PutContents(path1, "2")
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 2)

		v.Set(3)
		gbfsnotify.RemoveCallback(callback.Id)
		gbfile.PutContents(path1, "3")
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 3)
	})
}

func TestWatcher_Callback2(t *testing.T) {
	// multiple callbacks
	gbtest.C(t, func(t *gbtest.T) {
		path1 := gbfile.Temp(gbtime.TimestampNanoStr())
		t.Assert(gbfile.PutContents(path1, "1"), nil)
		defer func() {
			gbfile.Remove(path1)
		}()
		v1 := gbtype.NewInt(1)
		v2 := gbtype.NewInt(1)
		callback1, err1 := gbfsnotify.Add(path1, func(event *gbfsnotify.Event) {
			if event.IsWrite() {
				v1.Set(2)
				return
			}
		})
		callback2, err2 := gbfsnotify.Add(path1, func(event *gbfsnotify.Event) {
			if event.IsWrite() {
				v2.Set(2)
				return
			}
		})
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.AssertNE(callback1, nil)
		t.AssertNE(callback2, nil)

		t.Assert(gbfile.PutContents(path1, "2"), nil)
		time.Sleep(100 * time.Millisecond)
		t.Assert(v1.Val(), 2)
		t.Assert(v2.Val(), 2)

		v1.Set(3)
		v2.Set(3)
		gbfsnotify.RemoveCallback(callback1.Id)
		t.Assert(gbfile.PutContents(path1, "3"), nil)
		time.Sleep(100 * time.Millisecond)
		t.Assert(v1.Val(), 3)
		t.Assert(v2.Val(), 2)
	})
}

func TestWatcher_WatchFolderWithoutRecursively(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			err     error
			array   = gbarray.New(true)
			dirPath = gbfile.Temp(gbtime.TimestampNanoStr())
		)
		err = gbfile.Mkdir(dirPath)
		t.AssertNil(err)

		_, err = gbfsnotify.Add(dirPath, func(event *gbfsnotify.Event) {
			// fmt.Println(event.String())
			array.Append(1)
		}, false)
		t.AssertNil(err)
		time.Sleep(time.Millisecond * 100)
		t.Assert(array.Len(), 0)

		f, err := gbfile.Create(gbfile.Join(dirPath, "1"))
		t.AssertNil(err)
		t.AssertNil(f.Close())
		time.Sleep(time.Millisecond * 100)
		t.Assert(array.Len(), 1)
	})
}
