package go_redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(opt *redis.Options) *RedisClient {
	client := redis.NewClient(opt)
	return &RedisClient{client}
}

func (c *RedisClient) Lock(ctx context.Context, key string, seconds int, waitTime time.Duration) error {
	var startTime = time.Now()
	var err error

	client := c.Conn(ctx)

	if err != nil {
		return err
	}
	defer client.Close()
	var lockOk bool
	var sleepTime = time.Millisecond * 10
	for {
		setNxRes := client.SetNX(ctx, key, 1, time.Second*time.Duration(seconds))
		lockOk = setNxRes.Val()
		// 结束等待
		if lockOk || waitTime == 0 || time.Since(startTime) > waitTime {
			break
		}
		// 等待
		time.Sleep(sleepTime)
		sleepTime *= 2
	}
	if !lockOk {
		return errors.New(fmt.Sprintf("lock %s failed", key))
	}
	return err
}

func (c *RedisClient) UnLock(ctx context.Context, key string) {
	c.Del(ctx, key)
}
