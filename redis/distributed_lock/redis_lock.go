/**
 * @Author raven
 * @Description
 * @Date 2022/7/22
 **/
package distributed_lock

import (
	"context"
	"errors"
	"github.com/RavenHuo/go-pkg/redis/go_redis"
	"time"

	"github.com/RavenHuo/go-pkg/utils"
	"github.com/sirupsen/logrus"
)

//lua 脚本
const unlockScript = `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`
const checkAndRenewScript = `
        if (redis.call('exists', KEYS[1]) == 1) then
            redis.call('expire', KEYS[1], ARGV[1]); 
            return 1; 
        end;
        return 0;
    `
const RELEASE_SUCCESS = 1

type RedisDistributedLock struct {
	redisClient *go_redis.RedisClient

	// 配置
	option *DistributedLockOption
	// 锁的key
	key string
	// 通知channel
	notify chan struct{}
	// 看门狗 renew 次数
	watchDogLiveSecs int

	//锁持有者id，只有具有同一个id的协程才可以解锁
	uuId string
	//锁期时间，不可设置为0
	expiredDuration time.Duration

	//等待时间，在等待时间内会重试锁
	waitDuration time.Duration
}

func getRedisDistributedLock(redisClient *go_redis.RedisClient, option *DistributedLockOption) (*RedisDistributedLock, error) {
	lock := &RedisDistributedLock{
		redisClient: redisClient,
		option:      option,
	}
	lock.uuId = utils.GetUuid()
	return lock, nil
}

func (r *RedisDistributedLock) Lock(ctx context.Context, key string, expiredDuration time.Duration, waitDuration time.Duration) (bool, error) {
	if expiredDuration <= 0 || expiredDuration.Seconds() <= 0 {
		return false, errors.New("[RedisDistributedLock]expiredDuration can not less than 0")
	}

	r.expiredDuration = expiredDuration
	r.waitDuration = waitDuration

	startTime := time.Now()
	sleepTime := time.Millisecond * 10
	waitTime := waitDuration
	success := false
	var err error
	expiredTime := time.Duration(r.expiredDuration.Seconds())
	for {
		success, err = r.redisClient.SetNX(ctx, key, r.uuId, expiredTime).Result()
		// 结束等待
		if (err == nil && success == true) || waitTime == 0 || time.Since(startTime) > waitTime {
			break
		}
		// 等待
		logrus.Infof("[RedisDistributedLock.Lock]key is locked by others, sleep %s", sleepTime.String())
		time.Sleep(sleepTime)
		sleepTime *= 2
	}

	if err == nil && success {
		lockOpt := r.option
		if lockOpt.CheckIntervalSeconds > 0 && lockOpt.RenewSeconds > 0 {
			//go pLock.startWatchDogRoutine()
		}
	}
	return success, err
}

func (r *RedisDistributedLock) Unlock(ctx context.Context, key string) (bool, error) {
	pass, isValidErr := r.option.IsValid()
	if isValidErr != nil || !pass {
		return false, isValidErr
	}

	result, err := r.redisClient.Eval(ctx, unlockScript, []string{key}, r.uuId).Result()
	if err != nil {
		return false, err
	}
	intResult := result.(int32)
	if intResult == RELEASE_SUCCESS {
		return true, nil
	}
	return false, nil
}

func (r *RedisDistributedLock) renew(ctx context.Context) {

	result, err := r.redisClient.Eval(ctx, checkAndRenewScript, []string{r.key}, r.option.RenewSeconds).Result()
	if err != nil {
		return
	}
	intResult := result.(int32)
	if intResult != RELEASE_SUCCESS {
		r.cancelRenew(ctx)
	}
}

func (r *RedisDistributedLock) startWatchDogRoutine(ctx context.Context) {
	r.watchDogLiveSecs = 0
	ticker := time.NewTicker(time.Duration(r.option.CheckIntervalSeconds) * time.Second)
	for {
		select {
		case <-r.notify:
			return
		case <-ticker.C:
			r.checkAndRenewLock(ctx)
		}
	}
}

func (r *RedisDistributedLock) checkAndRenewLock(ctx context.Context) {
	option := r.option
	r.watchDogLiveSecs += option.CheckIntervalSeconds
	if option == nil {
		return
	}
	if option.MaxLiveSeconds > 0 &&
		r.watchDogLiveSecs >= option.MaxLiveSeconds {
		option.MaxLiveCb()
		r.cancelRenew(ctx)
		return
	}
	r.renew(ctx)
}

func (r *RedisDistributedLock) cancelRenew(ctx context.Context) {
	if r.option.CheckIntervalSeconds > 0 && r.option.RenewSeconds > 0 && r.notify != nil {
		r.notify <- struct{}{}
	}
}
