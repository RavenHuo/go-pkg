package mq

const (
	DefaultMaxRetryTimes        = 5
	DefaultRetryIntervalSeconds = 5
)

const (
	BrokerTypeKafka = "Kafka"
	BrokerTypeRedis = "redis" // 如果是redis类型的需要执行redis.Init()
)

type Options struct {
	BrokerAddress         []string
	BrokerConsumerGroupID string
	RetryConfig           *RetryConsumeConfig
	BrokerType            string
}

func NewOptions() (opts *Options) {
	return &Options{
		RetryConfig: DefaultRetryConsumerConfig(),
	}
}

type Option func(*Options)

func BrokerAddressOption(address []string) Option {
	return func(o *Options) {
		o.BrokerAddress = address
	}
}

func BrokerConsumerGroupIDOption(id string) Option {
	return func(o *Options) {
		o.BrokerConsumerGroupID = id
	}
}

func DefaultRetryConsumerConfig() *RetryConsumeConfig {
	return &RetryConsumeConfig{
		MaxRetryTimes:        DefaultMaxRetryTimes,
		RetryIntervalSeconds: DefaultRetryIntervalSeconds,
	}
}
func BrokerTypeOption(BrokerType string) Option {
	return func(o *Options) {
		o.BrokerType = BrokerType
	}
}
