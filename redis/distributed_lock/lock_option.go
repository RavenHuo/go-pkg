package distributed_lock

import (
	"errors"
)

/*
*注：
*1、ExpireSeconds不可以设置小于0的数
*2、CheckIntervalSeconds、RenewSeconds大于0时才会启动看门狗协程为锁续约
*3、MaxLiveSeconds设置为0时，不会运行MaxLiveCb回调函数，即看门狗协程没有最大存活时间，只有显式调用UnLock函数解锁且解锁成功才会结束看门狗协程的生命周期
*4、WaitDuration设置为0时不会重试加锁，只会加一次锁
 */
type CallBackFunc func()

//Redis 分布式锁的配置结构体
type DistributedLockOption struct {
	//锁的检查间隔，也可以理解为续约间隔，因此RenewSeconds必须大于CheckIntervalSeconds
	CheckIntervalSeconds int
	//续约时间，可设置为0
	RenewSeconds int
	//看门狗协程最长存活时间，可设置为0（除非是必须要成功的业务，否则强烈建议设置看门狗协程的最大存活时间，并且在回调函数中进行资源回收，报警等操作）
	MaxLiveSeconds int
	//看门狗协程达到最长存活时间关闭前的回调函数
	MaxLiveCb CallBackFunc

}

func (opt *DistributedLockOption) IsValid() (bool, error) {
	if opt.RenewSeconds < opt.CheckIntervalSeconds {
		return false, errors.New("[DistributedLockOption.IsValid]RenewSeconds can not less than CheckIntervalSeconds")
	}
	return true, nil
}




func (opt *DistributedLockOption) SetRenewSeconds(renewSeconds int) *DistributedLockOption {
	opt.RenewSeconds = renewSeconds
	return opt
}


func (opt *DistributedLockOption) SetCheckIntervalSeconds(checkInterval int) *DistributedLockOption {
	opt.CheckIntervalSeconds = checkInterval
	return opt
}

func (opt *DistributedLockOption) SetMaxLiveSeconds(maxLiveSeconds int) *DistributedLockOption {
	opt.MaxLiveSeconds = maxLiveSeconds
	return opt
}

func (opt *DistributedLockOption) SetMaxLiveCb(cb CallBackFunc) *DistributedLockOption {
	opt.MaxLiveCb = cb
	return opt
}
