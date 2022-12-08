/**
 * @Author raven
 * @Description
 * @Date 2022/7/28
 **/
package conn

import (
	"context"
	"time"

	"github.com/RavenHuo/go-kit/internal"
)

func (c *RedisConn) Ping(ctx context.Context) *StringCmder {
	cmder := c.warpDo(ctx, PingCommand)
	return cmder.(*StringCmder)
}

func (c *RedisConn) Touch(ctx context.Context, keys ...string) *IntCmder {
	cmder := c.warpDo(ctx, TouchCommand, internal.StringSlice2InterfaceSlice(keys)...)
	return cmder.(*IntCmder)
}

func (c *RedisConn) TTL(ctx context.Context, key string) *DurationCmder {
	cmder := c.warpDo(ctx, TtlCommand, key)
	return cmder.(*DurationCmder)
}

func (c *RedisConn) Type(ctx context.Context, key string) *StringCmder {
	cmder := c.warpDo(ctx, TypeCommand, key)
	return cmder.(*StringCmder)
}

func (c *RedisConn) Append(ctx context.Context, key, value string) *IntCmder {
	cmder := c.warpDo(ctx, AppendCommand, key, value)
	return cmder.(*IntCmder)
}

func (c *RedisConn) Decr(ctx context.Context, key string) *IntCmder {
	cmder := c.warpDo(ctx, DecrCommand, key)
	return cmder.(*IntCmder)
}

func (c *RedisConn) DecrBy(ctx context.Context, key string, decrement int64) *IntCmder {
	cmder := c.warpDo(ctx, DecrByCommand, key, decrement)
	return cmder.(*IntCmder)
}

func (c *RedisConn) GetRange(ctx context.Context, key string, start, end int64) *StringCmder {
	cmder := c.warpDo(ctx, GetRangeCommand, key, start, end)
	return cmder.(*StringCmder)
}

func (c *RedisConn) GetSet(ctx context.Context, key string, value interface{}) *StringCmder {
	cmder := c.warpDo(ctx, GetSetCommand, key, value)
	return cmder.(*StringCmder)
}

func (c *RedisConn) Incr(ctx context.Context, key string) *IntCmder {
	cmder := c.warpDo(ctx, IncrCommand, key)
	return cmder.(*IntCmder)
}

func (c *RedisConn) IncrBy(ctx context.Context, key string, value int64) *IntCmder {
	cmder := c.warpDo(ctx, IncrByCommand, key, value)
	return cmder.(*IntCmder)
}

func (c *RedisConn) IncrByFloat(ctx context.Context, key string, value float64) *FloatCmder {
	cmder := c.warpDo(ctx, IncrByFloatCommand, key, value)
	return cmder.(*FloatCmder)
}

func (c *RedisConn) MGet(ctx context.Context, keys ...string) *InterfaceSliceCmder {
	cmder := c.warpDo(ctx, MGetCommand, internal.StringSlice2InterfaceSlice(keys)...)
	return cmder.(*InterfaceSliceCmder)
}

func (c *RedisConn) MSet(ctx context.Context, values ...interface{}) *StringCmder {
	cmder := c.warpDo(ctx, MSetCommand, values...)
	return cmder.(*StringCmder)
}

func (c *RedisConn) MSetNX(ctx context.Context, values ...interface{}) *BoolCmder {
	cmder := c.warpDo(ctx, MSetNXCommand, values...)
	return cmder.(*BoolCmder)
}

func (c *RedisConn) SetRange(ctx context.Context, key string, offset int64, value string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) StrLen(ctx context.Context, key string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) GetBit(ctx context.Context, key string, offset int64) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) SetBit(ctx context.Context, key string, offset int64, value int) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) BitCount(ctx context.Context, key string, Start, End int64) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) BitOpAnd(ctx context.Context, destKey string, keys ...string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) BitOpOr(ctx context.Context, destKey string, keys ...string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) BitOpXor(ctx context.Context, destKey string, keys ...string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) BitOpNot(ctx context.Context, destKey string, key string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) BitField(ctx context.Context, key string, args ...interface{}) *IntSliceCmder {
	panic("implement me")
}

func (c *RedisConn) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	panic("implement me")
}

func (c *RedisConn) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmder {
	panic("implement me")
}

func (c *RedisConn) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmder {
	panic("implement me")
}

func (c *RedisConn) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmder {
	panic("implement me")
}

func (c *RedisConn) HDel(ctx context.Context, key string, fields ...string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) HExists(ctx context.Context, key, field string) *BoolCmder {
	panic("implement me")
}

func (c *RedisConn) HGet(ctx context.Context, key, field string) *StringCmder {
	panic("implement me")
}

func (c *RedisConn) HGetAll(ctx context.Context, key string) *StringStructMapCmder {
	panic("implement me")
}

func (c *RedisConn) HIncrBy(ctx context.Context, key, field string, incr int64) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) HIncrByFloat(ctx context.Context, key, field string, incr float64) *FloatCmder {
	panic("implement me")
}

func (c *RedisConn) HKeys(ctx context.Context, key string) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) HLen(ctx context.Context, key string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) HMGet(ctx context.Context, key string, fields ...string) *InterfaceSliceCmder {
	panic("implement me")
}

func (c *RedisConn) HSet(ctx context.Context, key string, values ...interface{}) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) HMSet(ctx context.Context, key string, values ...interface{}) *BoolCmder {
	panic("implement me")
}

func (c *RedisConn) HSetNX(ctx context.Context, key, field string, value interface{}) *BoolCmder {
	panic("implement me")
}

func (c *RedisConn) HVals(ctx context.Context, key string) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *StringCmder {
	panic("implement me")
}

func (c *RedisConn) LIndex(ctx context.Context, key string, index int64) *StringCmder {
	panic("implement me")
}

func (c *RedisConn) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) LLen(ctx context.Context, key string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) LPop(ctx context.Context, key string) *StringCmder {
	panic("implement me")
}

func (c *RedisConn) LPush(ctx context.Context, key string, values ...interface{}) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) LPushX(ctx context.Context, key string, values ...interface{}) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) LRange(ctx context.Context, key string, start, stop int64) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) LRem(ctx context.Context, key string, count int64, value interface{}) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) LSet(ctx context.Context, key string, index int64, value interface{}) *StringCmder {
	panic("implement me")
}

func (c *RedisConn) LTrim(ctx context.Context, key string, start, stop int64) *StringCmder {
	panic("implement me")
}

func (c *RedisConn) RPop(ctx context.Context, key string) *StringCmder {
	panic("implement me")
}

func (c *RedisConn) RPopLPush(ctx context.Context, source, destination string) *StringCmder {
	panic("implement me")
}

func (c *RedisConn) RPush(ctx context.Context, key string, values ...interface{}) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) RPushX(ctx context.Context, key string, values ...interface{}) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) SAdd(ctx context.Context, key string, members ...interface{}) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) SCard(ctx context.Context, key string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) SDiff(ctx context.Context, keys ...string) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) SDiffStore(ctx context.Context, destination string, keys ...string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) SInter(ctx context.Context, keys ...string) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) SInterStore(ctx context.Context, destination string, keys ...string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) SIsMember(ctx context.Context, key string, member interface{}) *BoolCmder {
	panic("implement me")
}

func (c *RedisConn) SMembers(ctx context.Context, key string) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) SMembersMap(ctx context.Context, key string) *StringStructMapCmder {
	panic("implement me")
}

func (c *RedisConn) SMove(ctx context.Context, source, destination string, member interface{}) *BoolCmder {
	panic("implement me")
}

func (c *RedisConn) SPop(ctx context.Context, key string) *StringCmder {
	panic("implement me")
}

func (c *RedisConn) SPopN(ctx context.Context, key string, count int64) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) SRandMember(ctx context.Context, key string) *StringCmder {
	panic("implement me")
}

func (c *RedisConn) SRandMemberN(ctx context.Context, key string, count int64) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) SRem(ctx context.Context, key string, members ...interface{}) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) SUnion(ctx context.Context, keys ...string) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) SUnionStore(ctx context.Context, destination string, keys ...string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZAdd(ctx context.Context, key string, members ...*ZMember) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZAddNX(ctx context.Context, key string, members ...*ZMember) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZAddXX(ctx context.Context, key string, members ...*ZMember) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZAddCh(ctx context.Context, key string, members ...*ZMember) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZAddNXCh(ctx context.Context, key string, members ...*ZMember) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZAddXXCh(ctx context.Context, key string, members ...*ZMember) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZIncr(ctx context.Context, key string, member *ZMember) *FloatCmder {
	panic("implement me")
}

func (c *RedisConn) ZIncrNX(ctx context.Context, key string, member *ZMember) *FloatCmder {
	panic("implement me")
}

func (c *RedisConn) ZIncrXX(ctx context.Context, key string, member *ZMember) *FloatCmder {
	panic("implement me")
}

func (c *RedisConn) ZCard(ctx context.Context, key string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZCount(ctx context.Context, key, min, max string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZLexCount(ctx context.Context, key, min, max string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZIncrBy(ctx context.Context, key string, increment float64, member string) *FloatCmder {
	panic("implement me")
}

func (c *RedisConn) ZPopMax(ctx context.Context, key string, count ...int64) *ZMemberSliceCmder {
	panic("implement me")
}

func (c *RedisConn) ZPopMin(ctx context.Context, key string, count ...int64) *ZMemberSliceCmder {
	panic("implement me")
}

func (c *RedisConn) ZRange(ctx context.Context, key string, start, stop int64) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *ZMemberSliceCmder {
	panic("implement me")
}

func (c *RedisConn) ZRangeByScore(ctx context.Context, key string, opt *ZMemberRangeBy) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) ZRangeByLex(ctx context.Context, key string, opt *ZMemberRangeBy) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) ZRangeByScoreWithScores(ctx context.Context, key string, opt *ZMemberRangeBy) *ZMemberSliceCmder {
	panic("implement me")
}

func (c *RedisConn) ZRank(ctx context.Context, key, member string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZRem(ctx context.Context, key string, members ...interface{}) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZRemRangeByScore(ctx context.Context, key, min, max string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZRemRangeByLex(ctx context.Context, key, min, max string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZRevRange(ctx context.Context, key string, start, stop int64) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *ZMemberSliceCmder {
	panic("implement me")
}

func (c *RedisConn) ZRevRangeByScore(ctx context.Context, key string, opt *ZMemberRangeBy) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) ZRevRangeByLex(ctx context.Context, key string, opt *ZMemberRangeBy) *StringSliceCmder {
	panic("implement me")
}

func (c *RedisConn) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *ZMemberRangeBy) *ZMemberSliceCmder {
	panic("implement me")
}

func (c *RedisConn) ZRevRank(ctx context.Context, key, member string) *IntCmder {
	panic("implement me")
}

func (c *RedisConn) ZScore(ctx context.Context, key, member string) *FloatCmder {
	panic("implement me")
}

func (c *RedisConn) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *BaseCmder {
	return c.Do(ctx, EvalCommand, script, keys, args)
}
