/**
 * @Author raven
 * @Description
 * @Date 2022/7/22
 **/
package distributed_lock

import (
	"context"
	"time"
)

// 分布式锁
type IDistributedLock interface {
	Lock(ctx context.Context, key string, lockDuration time.Duration, waitDuration time.Duration) (bool, error)
	Unlock(ctx context.Context, key string) (bool, error)
	renew(ctx context.Context)
}

type IDistributedLockClient interface {
	GetLock() (IDistributedLock,error)
}