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
func (natsClient *NatsClient) RegisterHandler(cb func(msg *nats.Msg)) {
	natsClient.NatsConn.Subscribe(natsClient.Subject, cb)
}

func (natsClient *NatsClient) Push(msg []byte) error {
	return natsClient.NatsConn.Publish(natsClient.Subject, msg)
}
