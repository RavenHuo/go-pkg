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
	option    *DistributedLockOption
	redisPool *redigo.Pool
}

func NewRedisDistributedLockClient(option *DistributedLockOption, redisOptions *redigo.Options) (*RedisDistributedLockClient, error) {
	redisPool, err := redigo.NewPool(redisOptions)
	if err != nil {
		return nil, err
	}
	return &RedisDistributedLockClient{
		option:    option,
		redisPool: redisPool,
	}, nil
}

func (c *RedisDistributedLockClient) GetLock() (*RedisDistributedLock, error) {
	_, err := c.option.IsValid()
	if err != nil {
		return nil, err
	}
	redisConn, err := c.redisPool.GetConn()
	if err != nil {
		return nil, err
	}
	return getRedisDistributedLock(redisConn, c.option)
}
