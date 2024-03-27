package mq

import (
	"github.com/RavenHuo/go-pkg/mq/broker"
)

var localEventMQ *EventMQ

// Init 初始化参数，申请资源
func Init(opts ...Option) error {
	localEventMQ = NewEventMQ()
	return localEventMQ.Init(opts...)
}

// Destroy Destroy释放Init申请的资源
func Destroy() {
	localEventMQ.Destroy()
}

// Publish 发布事件到MQ
// key为此消息的有序标识，当为""时，消息为无序，非""时消息会Hash到某一个Partition
func Publish(eventType EventType, key string, msg []byte) error {
	return localEventMQ.Publish(eventType, key, msg)
}

func PublishWithOpt(eventType EventType, key string, msg []byte, opts ...broker.PublishOption) error {
	return localEventMQ.Publish(eventType, key, msg, opts...)
}

// AddEventHandler 注册MQ事件处理回调函数
func AddEventHandler(eventType EventType, handler Handler, retryConfig *RetryConsumeConfig) error {
	return localEventMQ.AddEventHandler(eventType, handler, retryConfig)
}

// StartHandler 开始接受MQ事件通知并调用相应的事件处理回调函数
func StartHandler(opts ...broker.SubscribeOption) error {
	return localEventMQ.StartHandler(opts...)
}
