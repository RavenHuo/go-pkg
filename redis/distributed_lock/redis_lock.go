/**
 * @Author raven
 * @Description
 * @Date 2022/7/22
 **/
package distributed_lock

import (
	"context"
	"errors"
	"time"

	"github.com/RavenHuo/go-kit/redis/redigo/conn"
	"github.com/RavenHuo/go-kit/utils"
	"github.com/sirupsen/logrus"
)

type RedisDistributedLock struct {
	redisConn *conn.RedisConn

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

func getRedisDistributedLock(redisConn *conn.RedisConn, option *DistributedLockOption) (*RedisDistributedLock, error) {
	lock := &RedisDistributedLock{
		redisConn: redisConn,
		option:    option,
	}
	lock.uuId = utils.GetUuid()
	return lock, nil
}

func (r RedisDistributedLock) Lock(ctx context.Context, key string, expiredDuration time.Duration, waitDuration time.Duration) (bool, error) {
	if expiredDuration <= 0 || expiredDuration.Seconds() <= 0{
		return false, errors.New("[RedisDistributedLock]expiredDuration can not less than 0")
	}
	if r.redisConn == nil {
		return false, errors.New("[RedisDistributedLock.Lock]cannot get pb")
	}

	r.expiredDuration = expiredDuration
	r.waitDuration = waitDuration

	startTime := time.Now()
	sleepTime := time.Millisecond * 10
	waitTime := waitDuration
	success := false
	var err error
	expiredTime := int64(r.expiredDuration.Seconds())
	for {
		success, err = r.redisConn.SetNX(ctx, key, r.uuId, expiredTime).Result()
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

func (r RedisDistributedLock) Unlock(ctx context.Context, key string) (bool, error) {
	panic("implement me")
}

func (r RedisDistributedLock) renew(ctx context.Context) {
	panic("implement me")
}
func (r RedisDistributedLock) Close(ctx context.Context) error {
	return r.redisConn.Close()
}
