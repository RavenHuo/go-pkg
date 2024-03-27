package kafka

import (
	"context"
	"fmt"
	"github.com/RavenHuo/go-pkg/log"
	"github.com/RavenHuo/go-pkg/mq/broker"
	"github.com/Shopify/sarama"
	"sync"
	"time"
)

var (
	DefaultBrokerConfig  = sarama.NewConfig()
	DefaultClusterConfig = sarama.NewConfig()
	ctx                  = context.Background()
)

// consumerGroupHandler is the implementation of sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	handler      broker.HandlerFunc
	cg           sarama.ConsumerGroup
	retryInfoMap sync.Map
	subOpts      *broker.SubscribeOptions
}

func retryInfoKey(partition int32, offset int64) string {
	return fmt.Sprintf("%v-%v", partition, offset)
}

func (h *consumerGroupHandler) needRetry() bool {
	return h.subOpts.MaxRetryTimes > 0
}
func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var retryKey string
	var curRetryTimes int
	for msg := range claim.Messages() {
		curRetryTimes = 0
		if h.needRetry() {
			retryKey = retryInfoKey(msg.Partition, msg.Offset)
			if times, exists := h.retryInfoMap.Load(retryKey); exists {
				curRetryTimes = times.(int)
			}
		}

		log.Infof(ctx, "[kafka] Topic:%v Key:%v Value:%v Offset:%v Retry:%v",
			msg.Topic, string(msg.Key), string(msg.Value), msg.Offset, curRetryTimes)
		if curRetryTimes > 0 && curRetryTimes < h.subOpts.MaxRetryTimes {
			time.Sleep(time.Second * time.Duration(h.subOpts.RetryIntervalSecond)) // sleep for retry
		}

		p := &Publication{
			M:      msg.Value,
			T:      msg.Topic,
			Km:     msg,
			Cg:     h.cg,
			Sess:   sess,
			Offset: msg.Offset,
		}

		if err := h.handler(p); err == nil {

			if h.subOpts.AutoAck {
				p.Sess.MarkMessage(p.Km, "")
			}
			if curRetryTimes > 0 {
				log.Infof(ctx, "[kafka] successful to retry consume Topic:%v Msg:%v Offset:%v total retry:%v",
					msg.Topic, string(msg.Value), msg.Offset, curRetryTimes+1)
				h.retryInfoMap.Delete(retryKey)
			}
		} else if h.needRetry() {
			curRetryTimes++
			if curRetryTimes >= h.subOpts.MaxRetryTimes {

				p.Sess.MarkMessage(p.Km, "") // had tried max times, commit and not retry anymore
				h.retryInfoMap.Delete(retryKey)
				log.Errorf(ctx, "[kafka] Topic:%v Partition:%v Offset:%v Msg:%v had retried max times:%v",
					msg.Topic, msg.Partition, msg.Offset, string(msg.Value), h.subOpts.MaxRetryTimes)

				if h.subOpts.RetryFailedProcessFunc != nil {
					h.subOpts.RetryFailedProcessFunc(msg.Topic, msg.Value)
				}
			} else {

				log.Warnf(ctx, "[kafka] Topic:%v Partition:%v Offset:%v Msg:%v consume failed, retry it:%v",
					msg.Topic, msg.Partition, msg.Offset, string(msg.Value), curRetryTimes)
				h.retryInfoMap.Store(retryKey, curRetryTimes)
				break // restart
			}
		} else {

			log.Errorf(ctx, "[kafka] Topic:%v Partition:%v Offset:%v Msg:%v consume failed, not retry",
				msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
			p.Sess.MarkMessage(p.Km, "")
		}
	}
	return nil
}
