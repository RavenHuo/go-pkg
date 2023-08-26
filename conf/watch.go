package conf

type KeyValue struct {
	Key   []byte
	Value []byte
}
type Event struct {
	ChangeType ChangeType
	KeyValue   KeyValue
}

type ChangeType int32

const (
	Insert ChangeType = iota + 1
	Update
	Delete
)

// WatchHandler 监听处理器
type WatchHandler func(*Event, *Configuration) error

func (conf *Configuration) SetWatchHandler(handlers ...WatchHandler) {
	conf.watchHandlers = append(conf.watchHandlers, handlers...)
}

func (conf *Configuration) OnChange(event *Event) {
	for _, handler := range conf.watchHandlers {
		_ = handler(event, conf)
	}
}
