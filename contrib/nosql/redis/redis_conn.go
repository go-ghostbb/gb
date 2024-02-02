package redis

import (
	"context"
	"fmt"
	"ghostbb.io/gb"
	gbvar "ghostbb.io/gb/container/gb_var"
	gbredis "ghostbb.io/gb/database/gb_redis"
	gbjson "ghostbb.io/gb/encoding/gb_json"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbtrace "ghostbb.io/gb/net/gb_trace"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbconv "ghostbb.io/gb/util/gb_conv"
	gbutil "ghostbb.io/gb/util/gb_util"
	"reflect"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Conn manages the connection operations.
type Conn struct {
	ps    *redis.PubSub
	redis *Redis
}

// traceItem holds the information for redis trace.
type traceItem struct {
	err       error
	command   string
	args      []interface{}
	costMilli int64
}

const (
	traceInstrumentName               = "ghostbb.io/gb/database/gbredis"
	traceAttrRedisAddress             = "redis.address"
	traceAttrRedisDb                  = "redis.db"
	traceEventRedisExecution          = "redis.execution"
	traceEventRedisExecutionCommand   = "redis.execution.command"
	traceEventRedisExecutionCost      = "redis.execution.cost"
	traceEventRedisExecutionArguments = "redis.execution.arguments"
)

// Do send a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
func (c *Conn) Do(ctx context.Context, command string, args ...interface{}) (reply *gbvar.Var, err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	for k, v := range args {
		var (
			reflectInfo = gbutil.OriginTypeAndKind(v)
		)
		switch reflectInfo.OriginKind {
		case
			reflect.Struct,
			reflect.Map,
			reflect.Slice,
			reflect.Array:
			// Ignore slice types of: []byte.
			if _, ok := v.([]byte); !ok {
				if args[k], err = gbjson.Marshal(v); err != nil {
					return nil, err
				}
			}
		}
	}

	// Trace span start.
	tr := otel.GetTracerProvider().Tracer(traceInstrumentName, trace.WithInstrumentationVersion(gb.VERSION))
	_, span := tr.Start(ctx, "Redis."+command, trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	timestampMilli1 := gbtime.TimestampMilli()
	reply, err = c.doCommand(ctx, command, args...)
	timestampMilli2 := gbtime.TimestampMilli()

	// Trace span end.
	c.traceSpanEnd(ctx, span, &traceItem{
		err:       err,
		command:   command,
		args:      args,
		costMilli: timestampMilli2 - timestampMilli1,
	})
	return
}

// Do send a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
func (c *Conn) doCommand(ctx context.Context, command string, args ...interface{}) (reply *gbvar.Var, err error) {
	argStrSlice := gbconv.Strings(args)
	switch gbstr.ToLower(command) {
	case `subscribe`:
		c.ps = c.redis.client.Subscribe(ctx, argStrSlice...)

	case `psubscribe`:
		c.ps = c.redis.client.PSubscribe(ctx, argStrSlice...)

	case `unsubscribe`:
		if c.ps != nil {
			err = c.ps.Unsubscribe(ctx, argStrSlice...)
			if err != nil {
				err = gberror.Wrapf(err, `Redis PubSub Unsubscribe failed with arguments "%v"`, argStrSlice)
			}
		}

	case `punsubscribe`:
		if c.ps != nil {
			err = c.ps.PUnsubscribe(ctx, argStrSlice...)
			if err != nil {
				err = gberror.Wrapf(err, `Redis PubSub PUnsubscribe failed with arguments "%v"`, argStrSlice)
			}
		}

	default:
		arguments := make([]interface{}, len(args)+1)
		copy(arguments, []interface{}{command})
		copy(arguments[1:], args)
		reply, err = c.resultToVar(c.redis.client.Do(ctx, arguments...).Result())
		if err != nil {
			err = gberror.Wrapf(err, `Redis Client Do failed with arguments "%v"`, arguments)
		}
	}
	return
}

// resultToVar converts redis operation result to gbvar.Var.
func (c *Conn) resultToVar(result interface{}, err error) (*gbvar.Var, error) {
	if err == redis.Nil {
		err = nil
	}
	if err == nil {
		switch v := result.(type) {
		case []byte:
			return gbvar.New(string(v)), err

		case []interface{}:
			return gbvar.New(gbconv.Strings(v)), err

		case *redis.Message:
			result = &gbredis.Message{
				Channel:      v.Channel,
				Pattern:      v.Pattern,
				Payload:      v.Payload,
				PayloadSlice: v.PayloadSlice,
			}

		case *redis.Subscription:
			result = &gbredis.Subscription{
				Kind:    v.Kind,
				Channel: v.Channel,
				Count:   v.Count,
			}
		}
	}

	return gbvar.New(result), err
}

// Receive receives a single reply as gbvar.Var from the Redis server.
func (c *Conn) Receive(ctx context.Context) (*gbvar.Var, error) {
	if c.ps != nil {
		v, err := c.resultToVar(c.ps.Receive(ctx))
		if err != nil {
			err = gberror.Wrapf(err, `Redis PubSub Receive failed`)
		}
		return v, err
	}
	return nil, nil
}

// Close closes current PubSub or puts the connection back to connection pool.
func (c *Conn) Close(ctx context.Context) (err error) {
	if c.ps != nil {
		err = c.ps.Close()
		if err != nil {
			err = gberror.Wrapf(err, `Redis PubSub Close failed`)
		}
	}
	return
}

// Subscribe subscribes the client to the specified channels.
//
// https://redis.io/commands/subscribe/
func (c *Conn) Subscribe(ctx context.Context, channel string, channels ...string) ([]*gbredis.Subscription, error) {
	args := append([]interface{}{channel}, gbconv.Interfaces(channels)...)
	_, err := c.Do(ctx, "Subscribe", args...)
	if err != nil {
		return nil, err
	}
	subs := make([]*gbredis.Subscription, len(args))
	for i := 0; i < len(subs); i++ {
		v, err := c.Receive(ctx)
		if err != nil {
			return nil, err
		}
		subs[i] = v.Val().(*gbredis.Subscription)
	}
	return subs, err
}

// PSubscribe subscribes the client to the given patterns.
//
// Supported glob-style patterns:
// - h?llo subscribes to hello, hallo and hxllo
// - h*llo subscribes to hllo and heeeello
// - h[ae]llo subscribes to hello and hallo, but not hillo
//
// Use \ to escape special characters if you want to match them verbatim.
//
// https://redis.io/commands/psubscribe/
func (c *Conn) PSubscribe(ctx context.Context, pattern string, patterns ...string) ([]*gbredis.Subscription, error) {
	args := append([]interface{}{pattern}, gbconv.Interfaces(patterns)...)
	_, err := c.Do(ctx, "PSubscribe", args...)
	if err != nil {
		return nil, err
	}
	subs := make([]*gbredis.Subscription, len(args))
	for i := 0; i < len(subs); i++ {
		v, err := c.Receive(ctx)
		if err != nil {
			return nil, err
		}
		subs[i] = v.Val().(*gbredis.Subscription)
	}
	return subs, err
}

// ReceiveMessage receives a single message of subscription from the Redis server.
func (c *Conn) ReceiveMessage(ctx context.Context) (*gbredis.Message, error) {
	v, err := c.Receive(ctx)
	if err != nil {
		return nil, err
	}
	return v.Val().(*gbredis.Message), nil
}

// traceSpanEnd checks and adds redis trace information to OpenTelemetry.
func (c *Conn) traceSpanEnd(ctx context.Context, span trace.Span, item *traceItem) {
	if gbtrace.IsUsingDefaultProvider() || !gbtrace.IsTracingInternal() {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if item.err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, item.err))
	}

	span.SetAttributes(gbtrace.CommonLabels()...)

	span.SetAttributes(
		attribute.String(traceAttrRedisAddress, c.redis.config.Address),
		attribute.Int(traceAttrRedisDb, c.redis.config.Db),
	)

	jsonBytes, _ := gbjson.Marshal(item.args)
	span.AddEvent(traceEventRedisExecution, trace.WithAttributes(
		attribute.String(traceEventRedisExecutionCommand, item.command),
		attribute.String(traceEventRedisExecutionCost, fmt.Sprintf(`%d ms`, item.costMilli)),
		attribute.String(traceEventRedisExecutionArguments, string(jsonBytes)),
	))
}
