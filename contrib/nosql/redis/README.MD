# redis
gb `gbredis.Adapter` implements using `go-redis`.

# Installation
```
go get -u ghostbb.io/contrib/nosql/redis
```

Commonly imported at top of `main.go`:
```go
package main

import (
	_ "ghostbb.io/contrib/nosql/redis"

	// Other imported packages.
)

func main() {
	// Main logics.
}
```
