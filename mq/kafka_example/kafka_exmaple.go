package main

import (
	"context"
	"encoding/json"
	"github.com/RavenHuo/go-pkg/log"
	"github.com/RavenHuo/go-pkg/mq"
	"github.com/RavenHuo/go-pkg/mq/broker"
	"github.com/RavenHuo/go-pkg/mq/middleware"
	"github.com/rs/xid"
	"time"
)

var ctx = context.Background()

func main() {

	kafkaUrl := []string{
		"127.0.0.1:9091", "127.0.0.1:9092",
	}

	topic := mq.EventType("raven-demo")

	err := initMq(kafkaUrl)
	if err != nil {
		return
	}

	consumer(topic)

	startHandler()

	err = publish(topic)
	if err != nil {
		log.Errorf(ctx, "mq push failed, err:%s", err)
		return
	}
	time.Sleep(10 * time.Second)
}

// 初始化
func initMq(kafkaUrl []string) error {
	err := mq.Init(mq.BrokerAddressOption(kafkaUrl))
	if err != nil {
		log.Errorf(ctx, "init mq failed, err:%s", err)
		return err
	}
	return err
}

// 先初始化，然后加载完消费者之后，再启动服务
func startHandler() {
	groupId := "raven"
	err := mq.StartHandler(broker.SubscribeGroupID(groupId))
	if err != nil {
		log.Errorf(ctx, "start mq failed ,err:%s", err)
	}
}

// 消费
func consumer(topic mq.EventType) {

	mq.AddEventHandler(topic, middleware.WrapperMultiHandler(&SimpleHandler{}), mq.DefaultRetryConsumerConfig())
}

// 启动之后执行，生产者产生消息
func publish(topic mq.EventType) error {
	now := time.Now()
	info := &mq.EventInfo{
		MsgID:     xid.New().String(),
		Timestamp: &now,
		Type:      topic,
		Info:      []byte("hello world"),
	}
	msg, _ := json.Marshal(info)
	return mq.Publish(topic, "123", msg)

}

type SimpleHandler struct {
}

func (s SimpleHandler) Handler(event *mq.EventInfo) error {
	time.Sleep(time.Second)
	log.Infof(ctx, "success listen mq :%+v", event)
	return nil
}

func (s SimpleHandler) Name() string {
	return "simple_handler"
}
