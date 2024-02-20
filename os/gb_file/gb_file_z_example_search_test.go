package gbfile_test

import (
	"fmt"
	gbfile "ghostbb.io/gb/os/gb_file"
)

func ExampleSearch() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_search")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example content")

	// search file
	realPath, _ := gbfile.Search(fileName, tempDir)
	fmt.Println(gbfile.Basename(realPath))

	// Output:
	// gbfile_example.txt
}
