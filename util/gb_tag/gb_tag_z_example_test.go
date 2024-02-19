package gbtag_test

import (
	"fmt"
	"ghostbb.io/gb/frame/g"
	gbmeta "ghostbb.io/gb/util/gb_meta"
	gbtag "ghostbb.io/gb/util/gb_tag"
)

func ExampleSet() {
	type User struct {
		g.Meta `name:"User Struct" description:"{UserDescription}"`
	}
	gbtag.Sets(g.MapStrStr{
		`UserDescription`: `This is a demo struct named "User Struct"`,
	})
	fmt.Println(gbmeta.Get(User{}, `description`))

	// Output:
	// This is a demo struct named "User Struct"
}
