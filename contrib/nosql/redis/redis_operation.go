package redis

import (
	"context"
	gbvar "ghostbb.io/container/gb_var"
	gbredis "ghostbb.io/database/gb_redis"
	gberror "ghostbb.io/errors/gb_error"
)

// Do send a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
func (r *Redis) Do(ctx context.Context, command string, args ...interface{}) (*gbvar.Var, error) {
	conn, err := r.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close(ctx)
	}()
	return conn.Do(ctx, command, args...)
}

// Close closes the redis connection pool, which will release all connections reserved by this pool.
// It is commonly not necessary to call Close manually.
func (r *Redis) Close(ctx context.Context) (err error) {
	if err = r.client.Close(); err != nil {
		err = gberror.Wrap(err, `Operation Client Close failed`)
	}
	return
}

// Conn retrieves and returns a connection object for continuous operations.
// Note that you should call Close function manually if you do not use this connection any further.
func (r *Redis) Conn(ctx context.Context) (gbredis.Conn, error) {
	return &Conn{
		redis: r,
	}, nil
}
