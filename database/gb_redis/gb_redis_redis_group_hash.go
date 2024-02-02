package gbredis

import (
	"context"
	gbvar "ghostbb.io/gb/container/gb_var"
)

// IGroupHash manages redis hash operations.
// Implements see redis.GroupHash.
type IGroupHash interface {
	HSet(ctx context.Context, key string, fields map[string]interface{}) (int64, error)
	HSetNX(ctx context.Context, key, field string, value interface{}) (int64, error)
	HGet(ctx context.Context, key, field string) (*gbvar.Var, error)
	HStrLen(ctx context.Context, key, field string) (int64, error)
	HExists(ctx context.Context, key, field string) (int64, error)
	HDel(ctx context.Context, key string, fields ...string) (int64, error)
	HLen(ctx context.Context, key string) (int64, error)
	HIncrBy(ctx context.Context, key, field string, increment int64) (int64, error)
	HIncrByFloat(ctx context.Context, key, field string, increment float64) (float64, error)
	HMSet(ctx context.Context, key string, fields map[string]interface{}) error
	HMGet(ctx context.Context, key string, fields ...string) (gbvar.Vars, error)
	HKeys(ctx context.Context, key string) ([]string, error)
	HVals(ctx context.Context, key string) (gbvar.Vars, error)
	HGetAll(ctx context.Context, key string) (*gbvar.Var, error)
}
