package broker

// IBroker is an interface used for asynchronous messaging.
type IBroker interface {
	Init(...Option) error
	Options() *Options
	Connect() error
	Disconnect() error
	Publish(topic string, key string, msg []byte, opts ...PublishOption) error
	Subscribe(topic string, h HandlerFunc, opts ...SubscribeOption) (ISubscriber, error)
	String() string
}

// HandlerFunc is used to process messages via a subscription of a topic.
// The handler is passed a publication interface which contains the
// message and optional Ack method to acknowledge receipt of the message.
type HandlerFunc func(Event) error

// Event is given to a subscription handler for processing
type Event interface {
	Topic() string
	Message() []byte
	Ack() error
	Sequence() int64
}

// ISubscriber is a convenience return type for the Subscribe method
type ISubscriber interface {
	Options() *SubscribeOptions
	Topic() string
	Unsubscribe() error
}
