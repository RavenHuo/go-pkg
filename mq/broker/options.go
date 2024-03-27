package broker

import (
	"context"
	"crypto/tls"
	"github.com/Shopify/sarama"
	"strings"
	"time"
)

// Options 配置Opt
type Options struct {
	Type        string // mq type
	ProducerNum int    // TODO number of producer
	Addresses   []string
	Secure      bool
	TLSConfig   *tls.Config
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type Option func(*Options)

func (o Options) AddressStr() string {
	return strings.Join(o.Addresses, ",")
}

func WithType(typeName string) Option {
	return func(o *Options) {
		o.Type = typeName
	}
}

func ProducerNum(num int) Option {
	return func(o *Options) {
		o.ProducerNum = num
	}
}

// Addresses sets the host addresses to be used by the mq
func Address(address ...string) Option {
	return func(o *Options) {
		o.Addresses = address
	}
}

// Secure communication with the mq
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// Specify TLS Config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

// PublishOptions 发布opt
type PublishOptions struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context       context.Context
	MaxRetryTimes int
	RetryInterval time.Duration
}

func DefaultSPubOptions() *PublishOptions {
	return &PublishOptions{
		Context:       context.Background(),
		MaxRetryTimes: 3,
		RetryInterval: time.Millisecond * 500,
	}
}

type PublishOption func(*PublishOptions)

type RetryFinallyFailedProcessFunc func(topic string, msg []byte)

// SubscribeOptions 订阅Opt
type SubscribeOptions struct {
	// AutoAck defaults to true. When a handler returns
	// with a nil error the message is acked.
	AutoAck bool
	// Subscribers with the same group ID name
	// will create a shared subscription where each
	// receives a subset of messages.
	GroupID string

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	// retry consume policy
	MaxRetryTimes          int
	RetryIntervalSecond    int
	RetryFailedProcessFunc RetryFinallyFailedProcessFunc

	// 自动提交时间间隔
	OffsetsAutoCommitInterval time.Duration
	// The initial offset to use if no offset was previously committed.
	// Should be OffsetNewest or OffsetOldest. Defaults to OffsetNewest.
	OffsetsInitial int64

	MaxPollRecord int64         // 拉取消息的数量
	PollInternal  time.Duration // 拉取消息的时间间隔
}

type SubscribeOption func(*SubscribeOptions)

func DefaultSubscribeOptions() *SubscribeOptions {
	return &SubscribeOptions{
		AutoAck:                   true,
		OffsetsInitial:            sarama.OffsetNewest,
		OffsetsAutoCommitInterval: time.Millisecond * 500,
	}
}

func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	opt := SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// DisableAutoAck will disable auto acking of messages
// after they have been handled.
func DisableAutoAck() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = false
	}
}

// SubscribeGroupID sets the name of the queue to share messages on
func SubscribeGroupID(name string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.GroupID = name
	}
}

// SubscribeContext set context
func SubscribeContext(ctx context.Context) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Context = ctx
	}
}

// Retry set retry consume policy
func Retry(maxTimes int, interval int, processFunc RetryFinallyFailedProcessFunc) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.MaxRetryTimes = maxTimes
		o.RetryIntervalSecond = interval
		o.RetryFailedProcessFunc = processFunc
	}
}

func MaxPollRecord(record int64) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.MaxPollRecord = record
	}
}

func PollInternal(PollInternal time.Duration) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.PollInternal = PollInternal
	}
}
