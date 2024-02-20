package gbfile_test

import (
	"fmt"
	gbfile "ghostbb.io/gb/os/gb_file"
)

func ExampleSize() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_size")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "0123456789")
	fmt.Println(gbfile.Size(tempFile))

	// Output:
	// 10
}

func ExampleSizeFormat() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_size")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "0123456789")
	fmt.Println(gbfile.SizeFormat(tempFile))

	// Output:
	// 10.00B
}

func ExampleReadableSize() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_size")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "01234567899876543210")
	fmt.Println(gbfile.ReadableSize(tempFile))

	// Output:
	// 20.00B
}

func ExampleStrToSize() {
	size := gbfile.StrToSize("100MB")
	fmt.Println(size)

	// Output:
	// 104857600
}

func ExampleFormatSize() {
	sizeStr := gbfile.FormatSize(104857600)
	fmt.Println(sizeStr)
	sizeStr0 := gbfile.FormatSize(1024)
	fmt.Println(sizeStr0)
	sizeStr1 := gbfile.FormatSize(999999999999999999)
	fmt.Println(sizeStr1)

	// Output:
	// 100.00M
	// 1.00K
	// 888.18P
}
