package gbfile_test

import (
	"fmt"
	gbfile "ghostbb.io/gb/os/gb_file"
	"time"
)

func ExampleGetContentsWithCache() {
	// init
	var (
		fileName = "gbfile_example.txt"
		tempDir  = gbfile.Temp("gbfile_example_cache")
		tempFile = gbfile.Join(tempDir, fileName)
	)

	// write contents
	gbfile.PutContents(tempFile, "ghostbb example content")

	// It reads the file content with cache duration of one minute,
	// which means it reads from cache after then without any IO operations within on minute.
	fmt.Println(gbfile.GetContentsWithCache(tempFile, time.Minute))

	// write new contents will clear its cache
	gbfile.PutContents(tempFile, "new ghostbb example content")

	// There's some delay for cache clearing after file content change.
	time.Sleep(time.Second * 1)

	// read contents
	fmt.Println(gbfile.GetContentsWithCache(tempFile))

	// May Output:
	// ghostbb example content
	// new ghostbb example content
}
