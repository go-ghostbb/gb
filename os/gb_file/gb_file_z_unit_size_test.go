package gbfile_test

import (
	gbfile "ghostbb.io/gb/os/gb_file"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
)

func Test_Size(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			paths1 string = "/testfile_t1.txt"
			sizes  int64
		)

		createTestFile(paths1, "abcdefghijklmn")
		defer delTestFiles(paths1)

		sizes = gbfile.Size(testpath() + paths1)
		t.Assert(sizes, 14)

		sizes = gbfile.Size("")
		t.Assert(sizes, 0)

	})
}

func Test_SizeFormat(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var (
			paths1 = "/testfile_t1.txt"
			sizes  string
		)

		createTestFile(paths1, "abcdefghijklmn")
		defer delTestFiles(paths1)

		sizes = gbfile.SizeFormat(testpath() + paths1)
		t.Assert(sizes, "14.00B")

		sizes = gbfile.SizeFormat("")
		t.Assert(sizes, "0.00B")

	})
}

func Test_StrToSize(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbfile.StrToSize("0.00B"), 0)
		t.Assert(gbfile.StrToSize("16.00B"), 16)
		t.Assert(gbfile.StrToSize("1.00K"), 1024)
		t.Assert(gbfile.StrToSize("1.00KB"), 1024)
		t.Assert(gbfile.StrToSize("1.00KiloByte"), 1024)
		t.Assert(gbfile.StrToSize("15.26M"), gbconv.Int64(15.26*1024*1024))
		t.Assert(gbfile.StrToSize("15.26MB"), gbconv.Int64(15.26*1024*1024))
		t.Assert(gbfile.StrToSize("1.49G"), gbconv.Int64(1.49*1024*1024*1024))
		t.Assert(gbfile.StrToSize("1.49GB"), gbconv.Int64(1.49*1024*1024*1024))
		t.Assert(gbfile.StrToSize("8.73T"), gbconv.Int64(8.73*1024*1024*1024*1024))
		t.Assert(gbfile.StrToSize("8.73TB"), gbconv.Int64(8.73*1024*1024*1024*1024))
		t.Assert(gbfile.StrToSize("8.53P"), gbconv.Int64(8.53*1024*1024*1024*1024*1024))
		t.Assert(gbfile.StrToSize("8.53PB"), gbconv.Int64(8.53*1024*1024*1024*1024*1024))
		t.Assert(gbfile.StrToSize("8.01EB"), gbconv.Int64(8.01*1024*1024*1024*1024*1024*1024))
		t.Assert(gbfile.StrToSize("0.01ZB"), gbconv.Int64(0.01*1024*1024*1024*1024*1024*1024*1024))
		t.Assert(gbfile.StrToSize("0.01YB"), gbconv.Int64(0.01*1024*1024*1024*1024*1024*1024*1024*1024))
		t.Assert(gbfile.StrToSize("0.01BB"), gbconv.Int64(0.01*1024*1024*1024*1024*1024*1024*1024*1024*1024))
		t.Assert(gbfile.StrToSize("0.01AB"), gbconv.Int64(-1))
		t.Assert(gbfile.StrToSize("123456789"), 123456789)
	})
}

func Test_FormatSize(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbfile.FormatSize(0), "0.00B")
		t.Assert(gbfile.FormatSize(16), "16.00B")

		t.Assert(gbfile.FormatSize(1024), "1.00K")

		t.Assert(gbfile.FormatSize(16000000), "15.26M")

		t.Assert(gbfile.FormatSize(1600000000), "1.49G")

		t.Assert(gbfile.FormatSize(9600000000000), "8.73T")
		t.Assert(gbfile.FormatSize(9600000000000000), "8.53P")
	})
}

func Test_ReadableSize(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {

		var (
			paths1 string = "/testfile_t1.txt"
		)
		createTestFile(paths1, "abcdefghijklmn")
		defer delTestFiles(paths1)
		t.Assert(gbfile.ReadableSize(testpath()+paths1), "14.00B")
		t.Assert(gbfile.ReadableSize(""), "0.00B")

	})
}
