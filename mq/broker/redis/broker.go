package redis

import (
	"context"
	"fmt"
	"github.com/RavenHuo/go-pkg/log"
	"github.com/RavenHuo/go-pkg/mq/broker"
	"github.com/RavenHuo/go-pkg/redis/go_redis"
	"sync"
	"time"
)

const (
	MqKey = "say_mq:%s"
)

func GetKey(topic string) string {
	return fmt.Sprintf(MqKey, topic)
}

type Broker struct {
	client  *go_redis.RedisClient
	scMutex sync.Mutex
	opts    *broker.Options
}

func (b *Broker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(b.opts)
	}
	return nil
}

func (b *Broker) Options() *broker.Options {
	return b.opts
}

func (b *Broker) Connect() error {
	b.client = go_redis.GetRedis()
	return nil
}

func (b *Broker) Disconnect() error {
	return nil
}

func (b *Broker) Publish(topic string, key string, msg []byte, opts ...broker.PublishOption) error {
	pubOpt := broker.DefaultSPubOptions()
	for _, o := range opts {
		o(pubOpt)
	}
	ctx := context.Background()
	mqKey := GetKey(topic)
	pushResult := b.client.RPush(ctx, mqKey, string(msg))
	if pushResult.Err() == nil {
		log.Infof(ctx, "[redis_mq] publish success, topic:%s, message:%s", topic, string(msg))
		return nil
	}
	log.Errorf(ctx, "[redis_mq] publish failed, message:%+v, err:%s", msg, pushResult.Err())
	retryTime := 0
	for retryTime < pubOpt.MaxRetryTimes {
		pushResult = b.client.LPush(ctx, mqKey, string(msg))
		retryTime += 1
		if pushResult.Err() == nil {
			log.Infof(ctx, "[kafka] publish success, topic:%s, message:%+v, retry:%d", topic, msg, retryTime)
			break
		}
		log.Errorf(ctx, "[kafka] publish failed, topic:%s, message:%+v, retry:%d, err:%s", topic, msg, retryTime, pushResult.Err())
	}
	return pushResult.Err()
}

func (b *Broker) Subscribe(topic string, h broker.HandlerFunc, opts ...broker.SubscribeOption) (broker.ISubscriber, error) {
	opt := broker.DefaultSubscribeOptions()
	for _, o := range opts {
		o(opt)
	}

	mqChan := make(chan *Member, 10)
	go func() {
		b.popMsg(opt, topic, mqChan)
		for {
			select {
			// 每3秒同步一次
			case <-time.After(opt.PollInternal):
				b.popMsg(opt, topic, mqChan)
			}
		}
	}()

	handler := &consumerGroupHandler{
		handler: h,
		subOpts: opt,
	}
	go func() {
		handler.Consume(topic, mqChan)
	}()
	return &Subscriber{Opts: opt, T: topic}, nil
}

func (b *Broker) popMsg(opt *broker.SubscribeOptions, topic string, mqChan chan *Member) {
	ctx := context.Background()
	nowTs := time.Now().UnixNano() / 1e6
	mqKey := GetKey(topic)
	for i := 0; i < int(opt.MaxPollRecord); i++ {
		result := b.client.LPop(ctx, mqKey)
		if result.Err() == nil {
			m := &Member{
				T:      topic,
				M:      []byte(result.Val()),
				Offset: int64(i),
				ts:     nowTs,
			}
			mqChan <- m
		} else {
			// err 不为空证明list为空
			return
		}
	}
}

func (b *Broker) String() string {
	return "redis_mq"
}

func NewBroker() *Broker {
	options := &broker.Options{
		Context: context.Background(),
	}
	return &Broker{
		opts: options,
	}
}
