/**
 * @Author raven
 * @Description
 * @Date 2022/7/28
 **/
package conn

import (
	"context"

	"github.com/garyburd/redigo/redis"
)

type RedisConn struct {
	redis.Conn
	hooks hooks
}

func WrapperRedisConn(conn redis.Conn) *RedisConn {
	return &RedisConn{
		Conn: conn,
	}
}
func (c *RedisConn) AddHook(hook Hook) {
	c.hooks.AddHook(hook)
}

func (c *RedisConn) Do(ctx context.Context, redisCommand RedisCommand, args ...interface{}) *BaseCmder {
	cmder := newBaseCmder(string(redisCommand), args...)
	reply, err := c.process(ctx, cmder)
	cmder.setVal(reply)
	cmder.SetErr(err)
	return cmder
}

func (c *RedisConn) process(ctx context.Context, cmder *BaseCmder) (interface{}, error) {
	return c.hooks.process(ctx, cmder, func(ctx context.Context, cmder *BaseCmder) (interface{}, error) {
		return c.do(ctx, cmder)
	})
}

func (c *RedisConn) do(ctx context.Context, cmder *BaseCmder) (interface{}, error) {
	return c.Conn.Do(cmder.Cmd(), cmder.Args()...)
}

func (c *RedisConn) warpDo(ctx context.Context, commandName RedisCommand, args ...interface{}) ICmder {
	return wrapperCmder(commandName, c.Do(ctx, commandName, args...))
}

func (c *RedisConn) Close() error {
	return c.Conn.Close()
}
