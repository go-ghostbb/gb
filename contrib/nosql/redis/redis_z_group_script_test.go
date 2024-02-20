package redis_test

import (
	gbsha1 "ghostbb.io/gb/crypto/gb_sha1"
	gbredis "ghostbb.io/gb/database/gb_redis"
	"ghostbb.io/gb/frame/g"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func Test_GroupScript_Eval(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			script  = `return ARGV[1]`
			numKeys int64
			keys    = []string{"hello"}
			args    = []interface{}(nil)
		)
		v, err := redis.GroupScript().Eval(ctx, script, numKeys, keys, args)
		t.AssertNil(err)
		t.Assert(v.String(), "hello")
	})
}

func Test_GroupScript_EvalSha(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			script  = gbsha1.Encrypt(`return ARGV[1]`)
			numKeys int64
			keys    = []string{"hello"}
			args    = []interface{}(nil)
		)
		v, err := redis.GroupScript().EvalSha(ctx, script, numKeys, keys, args)
		t.AssertNil(err)
		t.Assert(v.String(), "hello")
	})
}

// https://redis.io/docs/manual/programmability/eval-intro/
func Test_GroupScript_ScriptLoad(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			script     = "return 'Immabe a cached script'"
			scriptSha1 = gbsha1.Encrypt(script)
		)
		_, err := redis.GroupScript().ScriptLoad(ctx, script)
		t.AssertNil(err)

		v, err := redis.GroupScript().EvalSha(ctx, scriptSha1, 0, nil, nil)
		t.AssertNil(err)
		t.Assert(v.String(), "Immabe a cached script")
	})
}

func Test_GroupScript_ScriptExists(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			script     = "return 'Immabe a cached script'"
			scriptSha1 = gbsha1.Encrypt(script)
			scriptSha2 = gbsha1.Encrypt("none")
		)
		_, err := redis.GroupScript().ScriptLoad(ctx, script)
		t.AssertNil(err)

		v, err := redis.GroupScript().ScriptExists(ctx, scriptSha1, scriptSha2)
		t.AssertNil(err)
		t.Assert(v, g.MapStrBool{
			scriptSha1: true,
			scriptSha2: false,
		})
	})
}

func Test_GroupScript_ScriptFlush(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		var (
			script     = "return 'Immabe a cached script'"
			scriptSha1 = gbsha1.Encrypt(script)
			scriptSha2 = gbsha1.Encrypt("none")
		)
		_, err := redis.GroupScript().ScriptLoad(ctx, script)
		t.AssertNil(err)

		v, err := redis.GroupScript().ScriptExists(ctx, scriptSha1, scriptSha2)
		t.AssertNil(err)
		t.Assert(v, g.MapStrBool{
			scriptSha1: true,
			scriptSha2: false,
		})

		err = redis.GroupScript().ScriptFlush(ctx, gbredis.ScriptFlushOption{SYNC: true})
		t.AssertNil(err)

		v, err = redis.GroupScript().ScriptExists(ctx, scriptSha1, scriptSha2)
		t.AssertNil(err)
		t.Assert(v, g.MapStrBool{
			scriptSha1: false,
			scriptSha2: false,
		})
	})
}

func Test_GroupScript_ScriptKill(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		defer redis.FlushDB(ctx)
		err := redis.GroupScript().ScriptKill(ctx)
		t.Assert(err.Error(), `Redis Client Do failed with arguments "[Script Kill]": NOTBUSY No scripts in execution right now.`)
	})
}
