package nats_client

import (
	"context"
	"errors"
	"github.com/RavenHuo/go-pkg/log"
	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	NatsConn *nats.Conn
	opts     *nats.Options
	Subject  string
	Observer []func(msg *nats.Msg)
}

func InitClient(opts *nats.Options, subject string) (*NatsClient, error) {
	natsClient := &NatsClient{
		opts:    opts,
		Subject: subject,
	}
	if opts == nil {
		return nil, errors.New("opt nil err")
	}
	natsConn, err := opts.Connect()
	if err != nil {
		log.Infof(context.Background(), "nats connect failed, opts:%+v, err:%s", opts, err)
		return nil, err
	}
	natsClient.Subject = subject
	natsClient.NatsConn = natsConn
	return natsClient, err
}

//发布订阅模型, 一个发布者, 多个订阅者, 多个订阅者都可以收到同一个消息
func (natsClient *NatsClient) RegisterHandler(cb func(msg *nats.Msg)) {
	natsClient.NatsConn.Subscribe(natsClient.Subject, cb)
}

// 队列模型, 一个发布者, 多个订阅者, 消息在多个消息中负载均衡分配, 分配给 A 消费者, 这个消息就不会再分配给其他消费者了
func (natsClient *NatsClient) RegisterQueueHandler(queue string, cb func(msg *nats.Msg)) {
	natsClient.NatsConn.QueueSubscribe(natsClient.Subject, queue, cb)
}

func (natsClient *NatsClient) Push(msg []byte) error {
	return natsClient.NatsConn.Publish(natsClient.Subject, msg)
}
