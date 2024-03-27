package kafka

import (
	"context"
	"errors"
	"github.com/RavenHuo/go-pkg/log"
	"github.com/RavenHuo/go-pkg/mq/broker"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type Broker struct {
	address []string

	c sarama.Client
	p sarama.SyncProducer

	sc []sarama.Client

	scMutex sync.Mutex
	opts    *broker.Options
}

func (k *Broker) Options() *broker.Options {
	return k.opts
}

func (k *Broker) Connect() error {
	if k.c != nil {
		return nil
	}

	config := k.getBrokerConfig()
	// For implementation reasons, the SyncProducer requires
	// `Producer.Return.Errors` and `Producer.Return.Successes`
	// to be set to true in its configuration.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Consumer.MaxWaitTime = 500 * time.Millisecond

	c, err := sarama.NewClient(k.address, config)
	if err != nil {
		return err
	}

	k.c = c

	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		return err
	}

	k.p = p
	k.scMutex.Lock()
	defer k.scMutex.Unlock()
	k.sc = make([]sarama.Client, 0)

	return nil
}

func (k *Broker) Disconnect() error {
	k.scMutex.Lock()
	defer k.scMutex.Unlock()
	for _, client := range k.sc {
		client.Close()
	}
	k.sc = nil
	k.p.Close()
	return k.c.Close()
}

func (k *Broker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(k.opts)
	}
	var addresses []string
	for _, addr := range k.opts.Addresses {
		if len(addr) == 0 {
			continue
		}
		addresses = append(addresses, addr)
	}
	if len(addresses) == 0 {
		return errors.New("mq mq addresses must specify")
	}
	k.address = addresses

	return nil
}

// Publish Publish发布事件到MQ
// key为此消息的有序标识，当为""时，消息为无序，非""时消息会Hash到某一个Partition
func (k *Broker) Publish(topic string, key string, msg []byte, opts ...broker.PublishOption) error {
	var keyEncoder sarama.Encoder
	if len(key) > 0 {
		keyEncoder = sarama.StringEncoder(key)
	}
	pubOpt := broker.DefaultSPubOptions()
	for _, o := range opts {
		o(pubOpt)
	}
	m := &sarama.ProducerMessage{
		Topic: topic,
		Key:   keyEncoder,
		Value: sarama.ByteEncoder(msg),
	}
	partition, offset, err := k.p.SendMessage(m)

	if err == nil {
		log.Infof(ctx, "[kafka] publish success, topic:%s, message:%s, partition:%d, offset:%d", topic, string(msg), partition, offset)
		return nil
	}
	log.Errorf(ctx, "[kafka] publish failed, topic:%s, message:%+v, err:%s", topic, msg, err)

	retryTime := 0
	for retryTime < pubOpt.MaxRetryTimes {
		partition, offset, err = k.p.SendMessage(m)
		retryTime += 1
		if err == nil {
			log.Infof(ctx, "[kafka] publish success, topic:%s, message:%+v, partition:%d, offset:%d, retry:%d", topic, msg, partition, offset, retryTime)
			break
		}
		log.Errorf(ctx, "[kafka] publish failed, topic:%s, message:%+v, retry:%d, err:%s", topic, msg, retryTime, err)
	}
	return err
}

func (k *Broker) getClusterClient(opts *broker.SubscribeOptions) (sarama.Client, error) {
	config := k.getClusterConfig(opts)
	cs, err := sarama.NewClient(k.address, config)
	if err != nil {
		return nil, err
	}
	k.scMutex.Lock()
	defer k.scMutex.Unlock()
	k.sc = append(k.sc, cs)
	return cs, nil
}

func (k *Broker) Subscribe(topic string, handler broker.HandlerFunc, opts ...broker.SubscribeOption) (broker.ISubscriber, error) {
	opt := broker.DefaultSubscribeOptions()
	for _, o := range opts {
		o(opt)
	}
	// we need to create a new client per consumer
	c, err := k.getClusterClient(opt)
	if err != nil {
		return nil, err
	}
	cg, err := sarama.NewConsumerGroupFromClient(opt.GroupID, c)
	if err != nil {
		return nil, err
	}

	h := &consumerGroupHandler{
		handler: handler,
		subOpts: opt,
		cg:      cg,
	}
	topics := []string{topic}
	go func() {
		for {
			select {
			case err := <-cg.Errors():
				if err != nil {
					log.Errorf(ctx, "[kafka] consumer group:%s topic:%s error:%s", opt.GroupID, topics, err.Error())
				}
			default:
				err := cg.Consume(ctx, topics, h)
				if err != nil {
					log.Errorf(ctx, "[kafka] Consume() group:%s topic:%s error:%s", opt.GroupID, topics, err.Error())
					time.Sleep(time.Millisecond * 500) // 连不上时避免太多的error日志
				}
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					log.Errorf(ctx, "[kafka] consumer group:%s topic:%s error:%s return!!!!!!!!!!!!!!!!!!", opt.GroupID, topics, err.Error())
					return
				}
			}
		}
	}()
	return &Subscriber{Cg: cg, Opts: opt, T: topic}, nil
}

func (k *Broker) String() string {
	return "kafka"
}

func NewBroker() *Broker {
	options := &broker.Options{
		Context: context.Background(),
	}
	return &Broker{
		opts: options,
	}
}

func (k *Broker) getBrokerConfig() *sarama.Config {
	return DefaultBrokerConfig
}

func (k *Broker) getClusterConfig(opts *broker.SubscribeOptions) *sarama.Config {
	clusterConfig := DefaultClusterConfig
	clusterConfig.Version = sarama.V2_2_0_0
	clusterConfig.Consumer.Return.Errors = true
	// 从最新的开始消费
	clusterConfig.Consumer.Offsets.Initial = opts.OffsetsInitial
	// commitInterval
	clusterConfig.Consumer.Offsets.AutoCommit.Interval = opts.OffsetsAutoCommitInterval // 500毫秒commit一次
	return clusterConfig
}
