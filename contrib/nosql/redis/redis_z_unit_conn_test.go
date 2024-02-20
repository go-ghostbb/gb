package redis_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func TestConn_DoWithTimeout(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := redis.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)

		_, err = conn.Do(ctx, "set", "test", "123")
		t.AssertNil(err)
		defer conn.Do(ctx, "del", "test")

		r, err := conn.Do(ctx, "get", "test")
		t.AssertNil(err)
		t.Assert(r.String(), "123")
	})
}

func TestConn_ReceiveVarWithTimeout(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		conn, err := redis.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)

		sub, err := conn.Subscribe(ctx, "gb")
		t.AssertNil(err)
		t.Assert(sub[0].Channel, "gb")

		_, err = redis.Publish(ctx, "gb", "test")
		t.AssertNil(err)

		msg, err := conn.ReceiveMessage(ctx)
		t.AssertNil(err)
		t.Assert(msg.Channel, "gb")
		t.Assert(msg.Payload, "test")
	})
}
