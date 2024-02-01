# redis
gb `gbredis.Adapter` implements using `go-redis`.

# Installation
```
go get -u github.com/Ghostbb-io/gb/contrib/nosql/redis
```

Commonly imported at top of `main.go`:
```go
package main

import (
	_ "github.com/Ghostbb-io/gb/contrib/nosql/redis"

	// Other imported packages.
)

func main() {
	// Main logics.
}
```
