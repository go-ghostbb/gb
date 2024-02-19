package gbutil_test

import (
	"fmt"
	"ghostbb.io/gb/frame/g"
	gbutil "ghostbb.io/gb/util/gb_util"
)

func ExampleSliceInsertBefore() {
	s1 := g.Slice{
		0, 1, 2, 3, 4,
	}
	s2 := gbutil.SliceInsertBefore(s1, 1, 8, 9)
	fmt.Println(s1)
	fmt.Println(s2)

	// Output:
	// [0 1 2 3 4]
	// [0 8 9 1 2 3 4]
}

func ExampleSliceInsertAfter() {
	s1 := g.Slice{
		0, 1, 2, 3, 4,
	}
	s2 := gbutil.SliceInsertAfter(s1, 1, 8, 9)
	fmt.Println(s1)
	fmt.Println(s2)

	// Output:
	// [0 1 2 3 4]
	// [0 1 8 9 2 3 4]
}
