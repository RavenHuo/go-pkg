package main

import (
	"context"
	"encoding/json"
	"github.com/RavenHuo/go-pkg/log"
	"github.com/RavenHuo/go-pkg/mq"
	"github.com/RavenHuo/go-pkg/mq/broker"
	"github.com/RavenHuo/go-pkg/mq/middleware"
	"github.com/RavenHuo/go-pkg/redis/go_redis"
	"github.com/rs/xid"
	"time"
)

func main() {

	log.Infof(context.Background(), "start---")
	go_redis.Init(&go_redis.Config{})
	topic := mq.EventType("raven-demo")
	topic2 := mq.EventType("raven-demo2")

	err := initMq()
	if err != nil {
		return
	}

	consumer(topic, topic2)

	startHandler()

	err = publish(topic)
	err = publish(topic2)
	if err != nil {
		log.Errorf(context.Background(), "mq push failed, err:%s", err)
		return
	}
	time.Sleep(10 * time.Second)
}

// 初始化
func initMq() error {
	err := mq.Init(mq.BrokerTypeOption(mq.BrokerTypeRedis))
	if err != nil {
		log.Errorf(context.Background(), "init mq failed, err:%s", err)
		return err
	}
	return err
}

// 先初始化，然后加载完消费者之后，再启动服务
func startHandler() {
	opts := []broker.SubscribeOption{
		broker.MaxPollRecord(100),
		broker.PollInternal(5 * time.Second),
	}
	err := mq.StartHandler(opts...)
	if err != nil {
		log.Errorf(context.Background(), "start mq failed ,err:%s", err)
	}
}

// 消费
func consumer(topic, topic2 mq.EventType) {

	mq.AddEventHandler(topic, middleware.WrapperMultiHandler(&SimpleHandler{}), mq.DefaultRetryConsumerConfig())
	mq.AddEventHandler(topic2, middleware.WrapperMultiHandler(&SimpleHandler{}), mq.DefaultRetryConsumerConfig())
}

// 启动之后执行，生产者产生消息
func publish(topic mq.EventType) error {
	now := time.Now()
	info := &mq.EventInfo{
		MsgID:     xid.New().String(),
		Timestamp: &now,
		Type:      topic,
		Info:      &Person{Name: "hello"},
	}
	msg, _ := json.Marshal(info)
	return mq.Publish(topic, "123", msg)

}

type SimpleHandler struct {
}

func (s SimpleHandler) Handler(event *mq.EventInfo) error {
	time.Sleep(time.Second)
	info := &Person{}
	if infoStr, ok := event.Info.(*json.RawMessage); ok {
		if err := json.Unmarshal(*infoStr, &info); err != nil {
			log.Warnf(context.Background(), "json unmarshal failed:%v info:%v", err.Error(), string(*infoStr))
			return nil
		}
	} else {
		log.Warnf(context.Background(), "invalid detail info type:%T", event.Info)
		return nil
	}
	return nil
}

func (s SimpleHandler) Name() string {
	return "simple_handler"
}

type Person struct {
	Name string `json:"name"`
}
