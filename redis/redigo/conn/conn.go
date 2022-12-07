/**
 * @Author raven
 * @Description
 * @Date 2022/7/27
 **/
package conn

import (
	"context"
	"time"
)

type IConn interface {
	Ping(ctx context.Context) *StringCmder
	Del(ctx context.Context, keys ...string) *IntCmder
	Exists(ctx context.Context, keys ...string) *IntCmder
	Expire(ctx context.Context, key string, expireSeconds uint64) *IntCmder
	ExpireAt(ctx context.Context, key string, tm time.Time) *IntCmder
	// WARNING Keys does not support pipelines
	Keys(ctx context.Context, pattern string) *StringSliceCmder
	Persist(ctx context.Context, key string) *BoolCmder
	PExpire(ctx context.Context, key string, expiration time.Duration) *BoolCmder
	PExpireAt(ctx context.Context, key string, tm time.Time) *DurationCmder
	PTTL(ctx context.Context, key string) *DurationCmder
	Rename(ctx context.Context, key, newkey string) *StringCmder
	RenameNX(ctx context.Context, key, newkey string) *BoolCmder
	Restore(ctx context.Context, key string, ttl time.Duration, value string) *StringCmder
	Touch(ctx context.Context, keys ...string) *IntCmder
	TTL(ctx context.Context, key string) *DurationCmder
	Type(ctx context.Context, key string) *StringCmder
	Append(ctx context.Context, key, value string) *IntCmder
	Decr(ctx context.Context, key string) *IntCmder
	DecrBy(ctx context.Context, key string, decrement int64) *IntCmder
	Get(ctx context.Context, key string) *StringCmder
	GetRange(ctx context.Context, key string, start, end int64) *StringCmder
	GetSet(ctx context.Context, key string, value interface{}) *StringCmder
	Incr(ctx context.Context, key string) *IntCmder
	IncrBy(ctx context.Context, key string, value int64) *IntCmder
	IncrByFloat(ctx context.Context, key string, value float64) *FloatCmder
	MGet(ctx context.Context, keys ...string) *InterfaceSliceCmder
	MSet(ctx context.Context, values ...interface{}) *StringCmder
	MSetNX(ctx context.Context, values ...interface{}) *BoolCmder
	Set(ctx context.Context, key string, value interface{}) *StringCmder
	SetNX(ctx context.Context, key string, value interface{}, expireSeconds int64) *BoolCmder
	SetEX(ctx context.Context, key string, value interface{}, expireSeconds uint64) *BoolCmder
	SetRange(ctx context.Context, key string, offset int64, value string) *IntCmder
	StrLen(ctx context.Context, key string) *IntCmder

	GetBit(ctx context.Context, key string, offset int64) *IntCmder
	SetBit(ctx context.Context, key string, offset int64, value int) *IntCmder
	BitCount(ctx context.Context, key string, Start, End int64) *IntCmder
	BitOpAnd(ctx context.Context, destKey string, keys ...string) *IntCmder
	BitOpOr(ctx context.Context, destKey string, keys ...string) *IntCmder
	BitOpXor(ctx context.Context, destKey string, keys ...string) *IntCmder
	BitOpNot(ctx context.Context, destKey string, key string) *IntCmder
	BitPos(ctx context.Context, key string, bit int64, pos ...int64) *IntCmder
	BitField(ctx context.Context, key string, args ...interface{}) *IntSliceCmder

	// WARNING Scan does not support pipelines
	Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error)
	SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmder
	HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmder
	ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmder

	HDel(ctx context.Context, key string, fields ...string) *IntCmder
	HExists(ctx context.Context, key, field string) *BoolCmder
	HGet(ctx context.Context, key, field string) *StringCmder
	HGetAll(ctx context.Context, key string) *StringStructMapCmder
	HIncrBy(ctx context.Context, key, field string, incr int64) *IntCmder
	HIncrByFloat(ctx context.Context, key, field string, incr float64) *FloatCmder
	HKeys(ctx context.Context, key string) *StringSliceCmder
	HLen(ctx context.Context, key string) *IntCmder
	HMGet(ctx context.Context, key string, fields ...string) *InterfaceSliceCmder
	HSet(ctx context.Context, key string, values ...interface{}) *IntCmder
	HMSet(ctx context.Context, key string, values ...interface{}) *BoolCmder
	HSetNX(ctx context.Context, key, field string, value interface{}) *BoolCmder
	HVals(ctx context.Context, key string) *StringSliceCmder

	BLPop(ctx context.Context, timeout time.Duration, keys ...string) *StringSliceCmder
	BRPop(ctx context.Context, timeout time.Duration, keys ...string) *StringSliceCmder
	BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *StringCmder
	LIndex(ctx context.Context, key string, index int64) *StringCmder
	LInsert(ctx context.Context, key, op string, pivot, value interface{}) *IntCmder
	LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *IntCmder
	LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *IntCmder
	LLen(ctx context.Context, key string) *IntCmder
	LPop(ctx context.Context, key string) *StringCmder
	LPush(ctx context.Context, key string, values ...interface{}) *IntCmder
	LPushX(ctx context.Context, key string, values ...interface{}) *IntCmder
	LRange(ctx context.Context, key string, start, stop int64) *StringSliceCmder
	LRem(ctx context.Context, key string, count int64, value interface{}) *IntCmder
	LSet(ctx context.Context, key string, index int64, value interface{}) *StringCmder
	LTrim(ctx context.Context, key string, start, stop int64) *StringCmder
	RPop(ctx context.Context, key string) *StringCmder
	RPopLPush(ctx context.Context, source, destination string) *StringCmder
	RPush(ctx context.Context, key string, values ...interface{}) *IntCmder
	RPushX(ctx context.Context, key string, values ...interface{}) *IntCmder

	SAdd(ctx context.Context, key string, members ...interface{}) *IntCmder
	SCard(ctx context.Context, key string) *IntCmder
	SDiff(ctx context.Context, keys ...string) *StringSliceCmder
	SDiffStore(ctx context.Context, destination string, keys ...string) *IntCmder
	SInter(ctx context.Context, keys ...string) *StringSliceCmder
	SInterStore(ctx context.Context, destination string, keys ...string) *IntCmder
	SIsMember(ctx context.Context, key string, member interface{}) *BoolCmder
	SMembers(ctx context.Context, key string) *StringSliceCmder
	SMembersMap(ctx context.Context, key string) *StringStructMapCmder
	SMove(ctx context.Context, source, destination string, member interface{}) *BoolCmder
	SPop(ctx context.Context, key string) *StringCmder
	SPopN(ctx context.Context, key string, count int64) *StringSliceCmder
	SRandMember(ctx context.Context, key string) *StringCmder
	SRandMemberN(ctx context.Context, key string, count int64) *StringSliceCmder
	SRem(ctx context.Context, key string, members ...interface{}) *IntCmder
	SUnion(ctx context.Context, keys ...string) *StringSliceCmder
	SUnionStore(ctx context.Context, destination string, keys ...string) *IntCmder

	ZAdd(ctx context.Context, key string, members ...*ZMember) *IntCmder
	ZAddNX(ctx context.Context, key string, members ...*ZMember) *IntCmder
	ZAddXX(ctx context.Context, key string, members ...*ZMember) *IntCmder
	ZAddCh(ctx context.Context, key string, members ...*ZMember) *IntCmder
	ZAddNXCh(ctx context.Context, key string, members ...*ZMember) *IntCmder
	ZAddXXCh(ctx context.Context, key string, members ...*ZMember) *IntCmder
	ZIncr(ctx context.Context, key string, member *ZMember) *FloatCmder
	ZIncrNX(ctx context.Context, key string, member *ZMember) *FloatCmder
	ZIncrXX(ctx context.Context, key string, member *ZMember) *FloatCmder
	ZCard(ctx context.Context, key string) *IntCmder
	ZCount(ctx context.Context, key, min, max string) *IntCmder
	ZLexCount(ctx context.Context, key, min, max string) *IntCmder
	ZIncrBy(ctx context.Context, key string, increment float64, member string) *FloatCmder
	ZPopMax(ctx context.Context, key string, count ...int64) *ZMemberSliceCmder
	ZPopMin(ctx context.Context, key string, count ...int64) *ZMemberSliceCmder
	ZRange(ctx context.Context, key string, start, stop int64) *StringSliceCmder
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) *ZMemberSliceCmder
	ZRangeByScore(ctx context.Context, key string, opt *ZMemberRangeBy) *StringSliceCmder
	ZRangeByLex(ctx context.Context, key string, opt *ZMemberRangeBy) *StringSliceCmder
	ZRangeByScoreWithScores(ctx context.Context, key string, opt *ZMemberRangeBy) *ZMemberSliceCmder
	ZRank(ctx context.Context, key, member string) *IntCmder
	ZRem(ctx context.Context, key string, members ...interface{}) *IntCmder
	ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *IntCmder
	ZRemRangeByScore(ctx context.Context, key, min, max string) *IntCmder
	ZRemRangeByLex(ctx context.Context, key, min, max string) *IntCmder
	ZRevRange(ctx context.Context, key string, start, stop int64) *StringSliceCmder
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *ZMemberSliceCmder
	ZRevRangeByScore(ctx context.Context, key string, opt *ZMemberRangeBy) *StringSliceCmder
	ZRevRangeByLex(ctx context.Context, key string, opt *ZMemberRangeBy) *StringSliceCmder
	ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *ZMemberRangeBy) *ZMemberSliceCmder
	ZRevRank(ctx context.Context, key, member string) *IntCmder
	ZScore(ctx context.Context, key, member string) *FloatCmder

	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *BaseCmder
}
