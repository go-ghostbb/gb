package gbredis

import (
	"context"
	gbvar "ghostbb.io/container/gb_var"
)

// IGroupList manages redis list operations.
// Implements see redis.GroupList.
type IGroupList interface {
	LPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	LPushX(ctx context.Context, key string, element interface{}, elements ...interface{}) (int64, error)
	RPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	RPushX(ctx context.Context, key string, value interface{}) (int64, error)
	LPop(ctx context.Context, key string, count ...int) (*gbvar.Var, error)
	RPop(ctx context.Context, key string, count ...int) (*gbvar.Var, error)
	LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error)
	LLen(ctx context.Context, key string) (int64, error)
	LIndex(ctx context.Context, key string, index int64) (*gbvar.Var, error)
	LInsert(ctx context.Context, key string, op LInsertOp, pivot, value interface{}) (int64, error)
	LSet(ctx context.Context, key string, index int64, value interface{}) (*gbvar.Var, error)
	LRange(ctx context.Context, key string, start, stop int64) (gbvar.Vars, error)
	LTrim(ctx context.Context, key string, start, stop int64) error
	BLPop(ctx context.Context, timeout int64, keys ...string) (gbvar.Vars, error)
	BRPop(ctx context.Context, timeout int64, keys ...string) (gbvar.Vars, error)
	RPopLPush(ctx context.Context, source, destination string) (*gbvar.Var, error)
	BRPopLPush(ctx context.Context, source, destination string, timeout int64) (*gbvar.Var, error)
}

// LInsertOp defines the operation name for function LInsert.
type LInsertOp string

const (
	LInsertBefore LInsertOp = "BEFORE"
	LInsertAfter  LInsertOp = "AFTER"
)
