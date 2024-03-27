package kafka

import (
	"github.com/Shopify/sarama"
)

type Publication struct {
	T      string
	Cg     sarama.ConsumerGroup
	Km     *sarama.ConsumerMessage
	M      []byte
	Sess   sarama.ConsumerGroupSession
	Offset int64
}

func (p *Publication) Topic() string {
	return p.T
}

func (p *Publication) Message() []byte {
	return p.M
}

func (p *Publication) Ack() error {
	p.Sess.MarkMessage(p.Km, "")
	return nil
}

func (p *Publication) Sequence() int64 {
	return p.Offset
}
