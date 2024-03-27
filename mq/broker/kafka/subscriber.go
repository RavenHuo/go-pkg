package kafka

import (
	"github.com/RavenHuo/go-pkg/mq/broker"
	"github.com/Shopify/sarama"
)

type Subscriber struct {
	Cg   sarama.ConsumerGroup
	T    string
	Opts *broker.SubscribeOptions
}

func (s *Subscriber) Options() *broker.SubscribeOptions {
	return s.Opts
}

func (s *Subscriber) Topic() string {
	return s.T
}

func (s *Subscriber) Unsubscribe() error {
	return s.Cg.Close()
}
