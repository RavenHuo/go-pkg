/**
 * @Author raven
 * @Description
 * @Date 2022/7/22
 **/
package distributed_lock

import (
	"github.com/RavenHuo/go-kit/redis/go_redis"
	"github.com/go-redis/redis/v8"
)

type RedisDistributedLockClient struct {
	redisClient *go_redis.RedisClient
}

func New(opt *redis.Options) *RedisDistributedLockClient {
	return &RedisDistributedLockClient{
		redisClient: go_redis.NewRedisClient(opt),
	}
}

func NewWithClient(redisClient *go_redis.RedisClient) *RedisDistributedLockClient {
	return &RedisDistributedLockClient{
		redisClient: redisClient,
	}
}

func (c *RedisDistributedLockClient) GetLock(option *DistributedLockOption) (*RedisDistributedLock, error) {
	_, err := option.IsValid()
	if err != nil {
		return nil, err
	}
	return getRedisDistributedLock(c.redisClient, option)
}
