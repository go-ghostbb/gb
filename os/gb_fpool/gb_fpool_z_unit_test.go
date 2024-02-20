package gbfpool_test

import (
	"context"
	gbfile "ghostbb.io/gb/os/gb_file"
	gbfpool "ghostbb.io/gb/os/gb_fpool"
	gblog "ghostbb.io/gb/os/gb_log"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	"os"
	"testing"
	"time"
)

// TestOpen test open file cache
func TestOpen(t *testing.T) {
	testFile := start("TestOpen.txt")

	gbtest.C(t, func(t *gbtest.T) {
		f, err := gbfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.AssertEQ(err, nil)
		f.Close()

		f2, err1 := gbfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.AssertEQ(err1, nil)
		t.AssertEQ(f, f2)
		f2.Close()
	})

	stop(testFile)
}

// TestOpenErr test open file error
func TestOpenErr(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		testErrFile := "errorPath"
		_, err := gbfpool.Open(testErrFile, os.O_RDWR, 0666)
		t.AssertNE(err, nil)

		// delete file error
		testFile := start("TestOpenDeleteErr.txt")
		pool := gbfpool.New(testFile, os.O_RDWR, 0666)
		_, err1 := pool.File()
		t.AssertEQ(err1, nil)
		stop(testFile)
		_, err1 = pool.File()
		t.AssertNE(err1, nil)

		// append mode delete file error and create again
		testFile = start("TestOpenCreateErr.txt")
		pool = gbfpool.New(testFile, os.O_CREATE, 0666)
		_, err1 = pool.File()
		t.AssertEQ(err1, nil)
		stop(testFile)
		_, err1 = pool.File()
		t.AssertEQ(err1, nil)

		// append mode delete file error
		testFile = start("TestOpenAppendErr.txt")
		pool = gbfpool.New(testFile, os.O_APPEND, 0666)
		_, err1 = pool.File()
		t.AssertEQ(err1, nil)
		stop(testFile)
		_, err1 = pool.File()
		t.AssertNE(err1, nil)

		// trunc mode delete file error
		testFile = start("TestOpenTruncErr.txt")
		pool = gbfpool.New(testFile, os.O_TRUNC, 0666)
		_, err1 = pool.File()
		t.AssertEQ(err1, nil)
		stop(testFile)
		_, err1 = pool.File()
		t.AssertNE(err1, nil)
	})
}

// TestOpenExpire test open file cache expire
func TestOpenExpire(t *testing.T) {
	testFile := start("TestOpenExpire.txt")

	gbtest.C(t, func(t *gbtest.T) {
		f, err := gbfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666, 100*time.Millisecond)
		t.AssertEQ(err, nil)
		f.Close()

		time.Sleep(150 * time.Millisecond)
		f2, err1 := gbfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666, 100*time.Millisecond)
		t.AssertEQ(err1, nil)
		//t.AssertNE(f, f2)
		f2.Close()
	})

	stop(testFile)
}

// TestNewPool test gbfpool new function
func TestNewPool(t *testing.T) {
	testFile := start("TestNewPool.txt")

	gbtest.C(t, func(t *gbtest.T) {
		f, err := gbfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.AssertEQ(err, nil)
		f.Close()

		pool := gbfpool.New(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		f2, err1 := pool.File()
		// pool not equal
		t.AssertEQ(err1, nil)
		//t.AssertNE(f, f2)
		f2.Close()

		pool.Close()
	})

	stop(testFile)
}

// test before
func start(name string) string {
	testFile := os.TempDir() + string(os.PathSeparator) + name
	if gbfile.Exists(testFile) {
		gbfile.Remove(testFile)
	}
	content := "123"
	gbfile.PutContents(testFile, content)
	return testFile
}

// test after
func stop(testFile string) {
	if gbfile.Exists(testFile) {
		err := gbfile.Remove(testFile)
		if err != nil {
			gblog.Error(context.TODO(), err)
		}
	}
}

func Test_ConcurrentOS(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		defer gbfile.Remove(path)
		f1, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.AssertNil(err)
		defer f1.Close()

		f2, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.AssertNil(err)
		defer f2.Close()

		for i := 0; i < 100; i++ {
			_, err = f1.Write([]byte("@1234567890#"))
			t.AssertNil(err)
		}
		for i := 0; i < 100; i++ {
			_, err = f2.Write([]byte("@1234567890#"))
			t.AssertNil(err)
		}

		for i := 0; i < 1000; i++ {
			_, err = f1.Write([]byte("@1234567890#"))
			t.AssertNil(err)
		}
		for i := 0; i < 1000; i++ {
			_, err = f2.Write([]byte("@1234567890#"))
			t.AssertNil(err)
		}
		t.Assert(gbstr.Count(gbfile.GetContents(path), "@1234567890#"), 2200)
	})

	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		defer gbfile.Remove(path)
		f1, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.AssertNil(err)
		defer f1.Close()

		f2, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.AssertNil(err)
		defer f2.Close()

		for i := 0; i < 1000; i++ {
			_, err = f1.Write([]byte("@1234567890#"))
			t.AssertNil(err)
		}
		for i := 0; i < 1000; i++ {
			_, err = f2.Write([]byte("@1234567890#"))
			t.AssertNil(err)
		}
		t.Assert(gbstr.Count(gbfile.GetContents(path), "@1234567890#"), 2000)
	})
	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		defer gbfile.Remove(path)
		f1, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.AssertNil(err)
		defer f1.Close()

		f2, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.AssertNil(err)
		defer f2.Close()

		s1 := ""
		for i := 0; i < 1000; i++ {
			s1 += "@1234567890#"
		}
		_, err = f2.Write([]byte(s1))
		t.AssertNil(err)

		s2 := ""
		for i := 0; i < 1000; i++ {
			s2 += "@1234567890#"
		}
		_, err = f2.Write([]byte(s2))
		t.AssertNil(err)

		t.Assert(gbstr.Count(gbfile.GetContents(path), "@1234567890#"), 2000)
	})
	// DATA RACE
	// gbtest.C(t, func(t *gbtest.T) {
	//	path := gbfile.Temp(gbtime.TimestampNanoStr())
	//	defer gbfile.Remove(path)
	//	f1, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	//	t.AssertNil(err)
	//	defer f1.Close()
	//
	//	f2, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	//	t.AssertNil(err)
	//	defer f2.Close()
	//
	//	wg := sync.WaitGroup{}
	//	ch := make(chan struct{})
	//	for i := 0; i < 1000; i++ {
	//		wg.Add(1)
	//		go func() {
	//			defer wg.Done()
	//			<-ch
	//			_, err = f1.Write([]byte("@1234567890#"))
	//			t.AssertNil(err)
	//		}()
	//	}
	//	for i := 0; i < 1000; i++ {
	//		wg.Add(1)
	//		go func() {
	//			defer wg.Done()
	//			<-ch
	//			_, err = f2.Write([]byte("@1234567890#"))
	//			t.AssertNil(err)
	//		}()
	//	}
	//	close(ch)
	//	wg.Wait()
	//	t.Assert(gbstr.Count(gbfile.GetContents(path), "@1234567890#"), 2000)
	// })
}

func Test_ConcurrentGFPool(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		path := gbfile.Temp(gbtime.TimestampNanoStr())
		defer gbfile.Remove(path)
		f1, err := gbfpool.Open(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.AssertNil(err)
		defer f1.Close()

		f2, err := gbfpool.Open(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.AssertNil(err)
		defer f2.Close()

		for i := 0; i < 1000; i++ {
			_, err = f1.Write([]byte("@1234567890#"))
			t.AssertNil(err)
		}
		for i := 0; i < 1000; i++ {
			_, err = f2.Write([]byte("@1234567890#"))
			t.AssertNil(err)
		}
		t.Assert(gbstr.Count(gbfile.GetContents(path), "@1234567890#"), 2000)
	})
	// DATA RACE
	// gbtest.C(t, func(t *gbtest.T) {
	//	path := gbfile.Temp(gbtime.TimestampNanoStr())
	//	defer gbfile.Remove(path)
	//	f1, err := gbfpool.Open(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	//	t.AssertNil(err)
	//	defer f1.Close()
	//
	//	f2, err := gbfpool.Open(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	//	t.AssertNil(err)
	//	defer f2.Close()
	//
	//	wg := sync.WaitGroup{}
	//	ch := make(chan struct{})
	//	for i := 0; i < 1000; i++ {
	//		wg.Add(1)
	//		go func() {
	//			defer wg.Done()
	//			<-ch
	//			_, err = f1.Write([]byte("@1234567890#"))
	//			t.AssertNil(err)
	//		}()
	//	}
	//	for i := 0; i < 1000; i++ {
	//		wg.Add(1)
	//		go func() {
	//			defer wg.Done()
	//			<-ch
	//			_, err = f2.Write([]byte("@1234567890#"))
	//			t.AssertNil(err)
	//		}()
	//	}
	//	close(ch)
	//	wg.Wait()
	//	t.Assert(gbstr.Count(gbfile.GetContents(path), "@1234567890#"), 2000)
	// })
}
