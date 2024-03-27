package mq

import (
	"context"
	"encoding/json"
	"github.com/RavenHuo/go-pkg/log"
	"github.com/RavenHuo/go-pkg/mq/broker"
)

// WrapperHandler 将Handler 包装成 kafka.HandlerFunc
func WrapperHandler(eventHandler Handler) broker.HandlerFunc {
	return func(event broker.Event) error {
		msg := json.RawMessage(event.Message())
		eventInfo := &EventInfo{
			Type: EventType(event.Topic()),
			Info: &msg,
			Seq:  event.Sequence(),
		}
		if err := json.Unmarshal(event.Message(), eventInfo); err != nil {
			log.Warnf(context.Background(), "[kafka] invalid message:%+v", string(event.Message()))
			return err
		}
		return eventHandler(eventInfo)
	}
}

func RetryConsumeConfigToSubOption(config *RetryConsumeConfig) broker.SubscribeOption {
	return broker.Retry(config.MaxRetryTimes, config.RetryIntervalSeconds, config.RetryFailedProcessFunc)
}
