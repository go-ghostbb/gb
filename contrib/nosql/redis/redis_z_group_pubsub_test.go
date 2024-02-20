package redis_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func Test_GroupPubSub_Publish(t *testing.T) {
	defer redis.FlushAll(ctx)
	gbtest.C(t, func(t *gbtest.T) {
		conn, subs, err := redis.Subscribe(ctx, "gb")
		t.AssertNil(err)
		t.Assert(subs[0].Channel, "gb")

		defer conn.Close(ctx)

		_, err = redis.Publish(ctx, "gb", "test")
		t.AssertNil(err)

		msg, err := conn.ReceiveMessage(ctx)
		t.AssertNil(err)
		t.Assert(msg.Channel, "gb")
		t.Assert(msg.Payload, "test")
	})
}

func Test_GroupPubSub_Subscribe(t *testing.T) {
	defer redis.FlushAll(ctx)
	gbtest.C(t, func(t *gbtest.T) {
		conn, subs, err := redis.Subscribe(ctx, "aa", "bb", "gb")
		t.AssertNil(err)
		t.Assert(len(subs), 3)
		t.Assert(subs[0].Channel, "aa")
		t.Assert(subs[1].Channel, "bb")
		t.Assert(subs[2].Channel, "gb")

		defer conn.Close(ctx)

		_, err = redis.Publish(ctx, "gb", "test")
		t.AssertNil(err)

		msg, err := conn.ReceiveMessage(ctx)
		t.AssertNil(err)
		t.Assert(msg.Channel, "gb")
		t.Assert(msg.Payload, "test")
	})
}

func Test_GroupPubSub_PSubscribe(t *testing.T) {
	defer redis.FlushAll(ctx)
	gbtest.C(t, func(t *gbtest.T) {
		conn, subs, err := redis.PSubscribe(ctx, "aa", "bb", "g?")
		t.AssertNil(err)
		t.Assert(len(subs), 3)
		t.Assert(subs[0].Channel, "aa")
		t.Assert(subs[1].Channel, "bb")
		t.Assert(subs[2].Channel, "g?")

		defer conn.Close(ctx)

		_, err = redis.Publish(ctx, "gb", "test")
		t.AssertNil(err)

		msg, err := conn.ReceiveMessage(ctx)
		t.AssertNil(err)
		t.Assert(msg.Channel, "gb")
		t.Assert(msg.Payload, "test")
	})
}
