/**
 * @Author raven
 * @Description
 * @Date 2022/7/28
 **/
package conn

type RedisCommand string

const (
	oKResponse = "OK"
)

// string
const (
	GetCommand         RedisCommand = "get"
	SetCommand         RedisCommand = "set"
	DelCommand         RedisCommand = "del"
	ExistCommand       RedisCommand = "exists"
	ExpireCommand      RedisCommand = "expire"
	ExpireAtCommand    RedisCommand = "expireat"
	KeysCommand        RedisCommand = "keys"
	PersistCommand     RedisCommand = "persist"
	PExpireCommand     RedisCommand = "pexpire"
	PExpireAtCommand   RedisCommand = "pexpireat"
	PttlCommand        RedisCommand = "pttl"
	RenameCommand      RedisCommand = "rename"
	RenameNxCommand    RedisCommand = "renamenx"
	RestoreCommand     RedisCommand = "restore"
	AppendCommand      RedisCommand = "append"
	DecrCommand        RedisCommand = "decr"
	DecrByCommand      RedisCommand = "decrby"
	SetNxCommand       RedisCommand = "setnx"
	SetExCommand       RedisCommand = "setex"
	GetRangeCommand    RedisCommand = "getrange"
	GetSetCommand      RedisCommand = "getset"
	IncrCommand        RedisCommand = "Incr"
	IncrByCommand      RedisCommand = "IncrBy"
	IncrByFloatCommand RedisCommand = "IncrByFloat"
	MGetCommand        RedisCommand = "MGet"
	MSetCommand        RedisCommand = "MSet"
	MSetNXCommand      RedisCommand = "MSetNX"
	HSetCommand        RedisCommand = "HSet"
	HGetCommand        RedisCommand = "HGet"
	HGetAllCommand     RedisCommand = "HGetAll"
)

// common
const (
	PingCommand  RedisCommand = "ping"
	TouchCommand RedisCommand = "touch"
	TtlCommand   RedisCommand = "ttl"
	TypeCommand  RedisCommand = "type"
	EvalCommand  RedisCommand = "eval"
)

type wrapperCmderFunc func(cmd *BaseCmder) ICmder

var redisCommandWrapperMap = map[RedisCommand]wrapperCmderFunc{
	SetCommand:       wrapperStringCmder,
	GetCommand:       wrapperStringCmder,
	DelCommand:       wrapperIntCmder,
	ExistCommand:     wrapperIntCmder,
	ExpireAtCommand:  wrapperIntCmder,
	KeysCommand:      wrapperStringSliceCmder,
	PersistCommand:   wrapperBoolCmder,
	PExpireCommand:   wrapperBoolCmder,
	PExpireAtCommand: wrapperDurationCmder,
	PttlCommand:      wrapperDurationCmder,
	RenameCommand:    wrapperStringCmder,
	RestoreCommand:   wrapperStringCmder,
	RenameNxCommand:  wrapperBoolCmder,
	AppendCommand:    wrapperIntCmder,
	DecrCommand:      wrapperIntCmder,
	DecrByCommand:    wrapperIntCmder,
	SetNxCommand:     wrapperBoolCmder,
	SetExCommand:     wrapperBoolCmder,

	// common
	PingCommand:  wrapperStringCmder,
	TouchCommand: wrapperIntCmder,
	TtlCommand:   wrapperDurationCmder,
	TypeCommand:  wrapperStringCmder,
}

func wrapperCmder(command RedisCommand, cmder *BaseCmder) ICmder {
	wrapperFunc := redisCommandWrapperMap[command]
	if wrapperFunc == nil {
		return nil
	}
	return wrapperFunc(cmder)
}
