package mq

import (
	"errors"
	"fmt"
	"github.com/RavenHuo/go-pkg/mq/broker"
	"github.com/RavenHuo/go-pkg/mq/broker/kafka"
	"github.com/RavenHuo/go-pkg/mq/broker/redis"
)

type EventMQ struct {
	broker              broker.IBroker
	eventTypeHandlerMap map[EventType]Handler
	exitNotify          chan struct{}
	eventTypeRetryMap   map[EventType]*RetryConsumeConfig
	opt                 *Options
}

func NewEventMQ() *EventMQ {
	return &EventMQ{
		eventTypeHandlerMap: make(map[EventType]Handler),
		exitNotify:          make(chan struct{}),
		eventTypeRetryMap:   make(map[EventType]*RetryConsumeConfig),
	}
}

func (e *EventMQ) Init(opts ...Option) (err error) {
	e.opt = NewOptions()
	for _, o := range opts {
		o(e.opt)
	}
	var bOpts []broker.Option
	bOpts = append(bOpts, broker.Address(e.opt.BrokerAddress...))

	if len(e.opt.BrokerType) == 0 {
		e.opt.BrokerType = BrokerTypeKafka // default is Kafka
	}

	if e.broker, err = newBroker(e.opt.BrokerType); err != nil {
		return
	}
	if err = e.broker.Init(bOpts...); err != nil {
		return
	}

	if err = e.broker.Connect(); err != nil {
		return
	}
	return nil
}

func (e *EventMQ) Destroy() {
	close(e.exitNotify)
	e.broker.Disconnect()
}

func (e *EventMQ) Publish(eventType EventType, key string, msg []byte, opts ...broker.PublishOption) (err error) {
	return e.broker.Publish(string(eventType), key, msg, opts...)
}

func (e *EventMQ) AddEventHandler(eventType EventType, handler Handler, retryConfig *RetryConsumeConfig) (err error) {
	if _, ok := e.eventTypeHandlerMap[eventType]; ok {
		return errors.New(fmt.Sprintf("event type:%v had handler yet", eventType))
	}
	e.eventTypeHandlerMap[eventType] = handler
	if retryConfig != nil {
		e.eventTypeRetryMap[eventType] = retryConfig
	}
	return nil
}

func (e *EventMQ) StartHandler(opts ...broker.SubscribeOption) (err error) {
	for k, v := range e.eventTypeHandlerMap {

		if retry, exists := e.eventTypeRetryMap[k]; exists {
			opts = append(opts, RetryConsumeConfigToSubOption(retry))
		}
		_, err = e.broker.Subscribe(string(k), WrapperHandler(v), opts...)
		if err != nil {
			return
		}

	}
	return nil
}

func newBroker(brokerType string) (broker broker.IBroker, err error) {
	switch brokerType {
	case BrokerTypeKafka:
		broker = kafka.NewBroker()
	case BrokerTypeRedis:
		broker = redis.NewBroker()
	default:
		err = errors.New(fmt.Sprintf("unsupported broker type:%v", brokerType))
	}
	return
}
