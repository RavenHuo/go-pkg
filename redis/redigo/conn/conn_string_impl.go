/**
 * @Author raven
 * @Description
 * @Date 2022/8/17
 **/
package conn

import (
	"context"
	"time"

	"github.com/RavenHuo/go-kit/internal"
)

func (c *RedisConn) Get(ctx context.Context, key string) *StringCmder {
	cmder := c.warpDo(ctx, GetCommand, key)
	return cmder.(*StringCmder)
}

func (c *RedisConn) Set(ctx context.Context, key string, value interface{}) *StringCmder {
	cmder := c.warpDo(ctx, SetCommand, key, value)
	return cmder.(*StringCmder)
}

func (c *RedisConn) Del(ctx context.Context, keys ...string) *IntCmder {
	cmder := c.warpDo(ctx, DelCommand, internal.StringSlice2InterfaceSlice(keys)...)
	return cmder.(*IntCmder)
}

func (c *RedisConn) Exists(ctx context.Context, keys ...string) *IntCmder {
	cmder := c.warpDo(ctx, ExistCommand, internal.StringSlice2InterfaceSlice(keys)...)
	return cmder.(*IntCmder)
}

func (c *RedisConn) Expire(ctx context.Context, key string, expireSeconds uint64) *IntCmder {
	cmder := c.warpDo(ctx, ExpireCommand, key, expireSeconds)
	return cmder.(*IntCmder)
}

func (c *RedisConn) ExpireAt(ctx context.Context, key string, tm time.Time) *IntCmder {
	cmder := c.warpDo(ctx, ExpireAtCommand, key, tm.Unix())
	return cmder.(*IntCmder)
}

func (c *RedisConn) Keys(ctx context.Context, pattern string) *StringSliceCmder {
	cmder := c.warpDo(ctx, KeysCommand, pattern)
	return cmder.(*StringSliceCmder)
}

func (c *RedisConn) Persist(ctx context.Context, key string) *BoolCmder {
	cmder := c.warpDo(ctx, PersistCommand, key)
	return cmder.(*BoolCmder)
}

func (c *RedisConn) PExpire(ctx context.Context, key string, expiration time.Duration) *BoolCmder {
	cmder := c.warpDo(ctx, PExpireCommand, key, expiration.Seconds())
	return cmder.(*BoolCmder)
}

func (c *RedisConn) PExpireAt(ctx context.Context, key string, tm time.Time) *DurationCmder {
	cmder := c.warpDo(ctx, PExpireAtCommand, key, tm.Unix())
	return cmder.(*DurationCmder)
}

func (c *RedisConn) PTTL(ctx context.Context, key string) *DurationCmder {
	cmder := c.warpDo(ctx, PttlCommand, key)
	return cmder.(*DurationCmder)
}

func (c *RedisConn) Rename(ctx context.Context, key, newKey string) *StringCmder {
	cmder := c.warpDo(ctx, RenameCommand, key, newKey)
	return cmder.(*StringCmder)
}

func (c *RedisConn) RenameNX(ctx context.Context, key, newKey string) *BoolCmder {
	cmder := c.warpDo(ctx, RenameNxCommand, key, newKey)
	return cmder.(*BoolCmder)
}

func (c *RedisConn) Restore(ctx context.Context, key string, ttl time.Duration, value string) *StringCmder {
	cmder := c.warpDo(ctx, RestoreCommand, key, ttl.Seconds(), value)
	return cmder.(*StringCmder)
}

func (c *RedisConn) SetNX(ctx context.Context, key string, value interface{}, expireSeconds int64) *BoolCmder {
	cmder := c.warpDo(ctx, SetNxCommand, key, value, expireSeconds)
	return cmder.(*BoolCmder)
}

func (c *RedisConn) SetEX(ctx context.Context, key string, value interface{}, expireSeconds uint64) *BoolCmder {
	cmder := c.warpDo(ctx, SetExCommand, key, expireSeconds, value)
	return cmder.(*BoolCmder)
}