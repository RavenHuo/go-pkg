package redis

import (
	"github.com/RavenHuo/go-pkg/mq/broker"
)

type Subscriber struct {
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
	return nil
}
