package gbfile_test

import (
	"fmt"
	gbfile "ghostbb.io/gb/os/gb_file"
)

func ExampleGetContents() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_content")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example content")

	// It reads and returns the file content as string.
	// It returns empty string if it fails reading, for example, with permission or IO error.
	fmt.Println(gbfile.GetContents(tempFile))

	// Output:
	// ghostbb example content
}

func ExampleGetBytes() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_content")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example content")

	// It reads and returns the file content as []byte.
	// It returns nil if it fails reading, for example, with permission or IO error.
	fmt.Println(gbfile.GetBytes(tempFile))

	// Output:
	// [103 104 111 115 116 98 98 32 101 120 97 109 112 108 101 32 99 111 110 116 101 110 116]
}

func ExamplePutContents() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_content")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// It creates and puts content string into specifies file path.
	// It automatically creates directory recursively if it does not exist.
	gbfile.PutContents(tempFile, "ghostbb example content")

	// read contents
	fmt.Println(gbfile.GetContents(tempFile))

	// Output:
	// ghostbb example content
}

func ExamplePutBytes() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_content")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutBytes(tempFile, []byte("ghostbb example content"))

	// read contents
	fmt.Println(gbfile.GetContents(tempFile))

	// Output:
	// ghostbb example content
}

func ExamplePutContentsAppend() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_content")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example content")

	// read contents
	fmt.Println(gbfile.GetContents(tempFile))

	// It creates and append content string into specifies file path.
	// It automatically creates directory recursively if it does not exist.
	gbfile.PutContentsAppend(tempFile, " append content")

	// read contents
	fmt.Println(gbfile.GetContents(tempFile))

	// Output:
	// ghostbb example content
	// ghostbb example content append content
}

func ExamplePutBytesAppend() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_content")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example content")

	// read contents
	fmt.Println(gbfile.GetContents(tempFile))

	// write contents
	gbfile.PutBytesAppend(tempFile, []byte(" append"))

	// read contents
	fmt.Println(gbfile.GetContents(tempFile))

	// Output:
	// ghostbb example content
	// ghostbb example content append
}

func ExampleGetNextCharOffsetByPath() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_content")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example content")

	// read contents
	index := gbfile.GetNextCharOffsetByPath(tempFile, 'f', 0)
	fmt.Println(index)

	// Output:
	// 2
}

func ExampleGetBytesTilCharByPath() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_content")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example content")

	// read contents
	fmt.Println(gbfile.GetBytesTilCharByPath(tempFile, 'f', 0))

	// Output:
	// [103 111 102] 2
}

func ExampleGetBytesByTwoOffsetsByPath() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_content")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example content")

	// read contents
	fmt.Println(gbfile.GetBytesByTwoOffsetsByPath(tempFile, 0, 7))

	// Output:
	// [103 104 111 115 116 98 98]
}

func ExampleReadLines() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_content")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "L1 ghostbb example content\nL2 ghostbb example content")

	// read contents
	gbfile.ReadLines(tempFile, func(text string) error {
		// Process each line
		fmt.Println(text)
		return nil
	})

	// Output:
	// L1 ghostbb example content
	// L2 ghostbb example content
}

func ExampleReadLinesBytes() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_content")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "L1 ghostbb example content\nL2 ghostbb example content")

	// read contents
	gbfile.ReadLinesBytes(tempFile, func(bytes []byte) error {
		// Process each line
		fmt.Println(bytes)
		return nil
	})

	// Output:
	// [76 49 32 103 104 111 115 116 98 98 32 101 120 97 109 112 108 101 32 99 111 110 116 101 110 116]
	// [76 50 32 103 104 111 115 116 98 98 32 101 120 97 109 112 108 101 32 99 111 110 116 101 110 116]
}
