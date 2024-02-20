package redis_test

import (
	gbredis "ghostbb.io/gb/database/gb_redis"
	gbctx "ghostbb.io/gb/os/gb_ctx"
)

var (
	ctx    = gbctx.GetInitCtx()
	config = &gbredis.Config{
		Address: `:6379`,
		Db:      1,
	}
	redis, _ = gbredis.New(config)
)
