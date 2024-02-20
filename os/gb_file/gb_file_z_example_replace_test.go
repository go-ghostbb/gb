package gbfile_test

import (
	"fmt"
	gbfile "ghostbb.io/gb/os/gb_file"
	"regexp"
)

func ExampleReplaceFile() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_replace")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example content")

	// read contents
	fmt.Println(gbfile.GetContents(tempFile))

	// It replaces content directly by file path.
	gbfile.ReplaceFile("content", "replace word", tempFile)

	fmt.Println(gbfile.GetContents(tempFile))

	// Output:
	// ghostbb example content
	// ghostbb example replace word
}

func ExampleReplaceFileFunc() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_replace")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example 123")

	// read contents
	fmt.Println(gbfile.GetContents(tempFile))

	// It replaces content directly by file path and callback function.
	gbfile.ReplaceFileFunc(func(path, content string) string {
		// Replace with regular match
		reg, _ := regexp.Compile(`\d{3}`)
		return reg.ReplaceAllString(content, "[num]")
	}, tempFile)

	fmt.Println(gbfile.GetContents(tempFile))

	// Output:
	// ghostbb example 123
	// ghostbb example [num]
}

func ExampleReplaceDir() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_replace")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example content")

	// read contents
	fmt.Println(gbfile.GetContents(tempFile))

	// It replaces content of all files under specified directory recursively.
	gbfile.ReplaceDir("content", "replace word", tempDir, "gbfile_example.txt", true)

	// read contents
	fmt.Println(gbfile.GetContents(tempFile))

	// Output:
	// ghostbb example content
	// ghostbb example replace word
}

func ExampleReplaceDirFunc() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_replace")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example 123")

	// read contents
	fmt.Println(gbfile.GetContents(tempFile))

	// It replaces content of all files under specified directory with custom callback function recursively.
	gbfile.ReplaceDirFunc(func(path, content string) string {
		// Replace with regular match
		reg, _ := regexp.Compile(`\d{3}`)
		return reg.ReplaceAllString(content, "[num]")
	}, tempDir, "gbfile_example.txt", true)

	fmt.Println(gbfile.GetContents(tempFile))

	// Output:
	// ghostbb example 123
	// ghostbb example [num]

}
