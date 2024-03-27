package redis

import (
	"context"
	"fmt"
	"github.com/RavenHuo/go-pkg/log"
	"github.com/RavenHuo/go-pkg/mq/broker"
	"sync"
	"time"
)

type consumerGroupHandler struct {
	handler      broker.HandlerFunc
	retryInfoMap sync.Map
	subOpts      *broker.SubscribeOptions
}

func retryInfoKey(ts int64, offset int64) string {
	return fmt.Sprintf("%v-%v", ts, offset)
}

func (h *consumerGroupHandler) needRetry() bool {
	return h.subOpts.MaxRetryTimes > 0
}

func (h *consumerGroupHandler) Consume(topics string, msgChan chan *Member) {
	var retryKey string
	var curRetryTimes int
	ctx := context.Background()
	for msg := range msgChan {
		if h.needRetry() {
			retryKey = retryInfoKey(msg.ts, msg.Offset)
			if times, exists := h.retryInfoMap.Load(retryKey); exists {
				curRetryTimes = times.(int)
			}
		}

		log.Infof(ctx, "[redis_mq] Topic:%v Value:%v seq:%v Retry:%v",
			msg.Topic(), string(msg.M), msg.Sequence(), curRetryTimes)
		if curRetryTimes > 0 && curRetryTimes < h.subOpts.MaxRetryTimes {
			time.Sleep(time.Second * time.Duration(h.subOpts.RetryIntervalSecond)) // sleep for retry
		}

		if err := h.handler(msg); err == nil {

			if curRetryTimes > 0 {
				log.Infof(ctx, "[redis_mq] successful to retry consume Topic:%v seq:%v total retry:%v",
					msg.Topic(), msg.Sequence(), curRetryTimes+1)
				h.retryInfoMap.Delete(retryKey)
			}
		} else if h.needRetry() {
			curRetryTimes++
			if curRetryTimes >= h.subOpts.MaxRetryTimes {

				h.retryInfoMap.Delete(retryKey)
				log.Errorf(ctx, "[redis_mq] Topic:%v Partition:%v Offset:%v Msg:%v had retried max times:%v",
					msg.Topic(), msg.Offset, string(msg.Message()), h.subOpts.MaxRetryTimes)

				if h.subOpts.RetryFailedProcessFunc != nil {
					h.subOpts.RetryFailedProcessFunc(msg.Topic(), msg.Message())
				}
			} else {

				log.Warnf(ctx, "[redis_mq] Topic:%v req:%v Msg:%v consume failed, retry it:%v",
					msg.Topic(), msg.Sequence(), string(msg.Message()), curRetryTimes)
				h.retryInfoMap.Store(retryKey, curRetryTimes)
				continue // restart
			}
		} else {

			log.Errorf(ctx, "[redis_mq] Topic:%v seq:%v Msg:%v consume failed, not retry",
				msg.Topic(), msg.Sequence(), string(msg.Message()))
		}
	}
}
