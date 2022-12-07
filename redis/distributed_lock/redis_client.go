/**
 * @Author raven
 * @Description
 * @Date 2022/7/22
 **/
package distributed_lock

import (
	"github.com/RavenHuo/go-kit/redis/redigo"
)

type RedisDistributedLockClient struct {
	redisPool *redigo.Pool
}

func NewRedisDistributedLockClient(redisOptions *redigo.Options) (*RedisDistributedLockClient, error) {
	redisPool, err := redigo.NewPool(redisOptions)
	if err != nil {
		return nil, err
	}
	return &RedisDistributedLockClient{
		redisPool: redisPool,
	}, nil
}

func (c *RedisDistributedLockClient) GetLock(option *DistributedLockOption) (*RedisDistributedLock, error) {
	_, err := option.IsValid()
	if err != nil {
		return nil, err
	}
	return getRedisDistributedLock(c.redisPool, option)
}
