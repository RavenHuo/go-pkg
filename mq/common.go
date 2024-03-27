package mq

type RetryConsumeConfig struct {
	MaxRetryTimes          int
	RetryIntervalSeconds   int
	RetryFailedProcessFunc func(topic string, msg []byte)
}

// EventType topic
type EventType string

// MQ事件处理回调函数定义
type EventHandler interface {
	Handler(event *EventInfo) error // 处理方法
	Name() string                   // 处理器的名字
}

type Handler func(*EventInfo) error
