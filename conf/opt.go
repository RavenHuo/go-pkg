package conf

const (
	defaultKeyDelimiter = "."
	defaultPath         = "etc/conf.yaml"
)

type Options struct {
	path          string
	keyDelimiter  string
	watchHandlers []WatchHandler
}

type Opt func(*Options)

func WithPath(path string) Opt {
	return func(options *Options) {
		options.path = path
	}
}

func WithKeyDelimiter(keyDelimiter string) Opt {
	return func(options *Options) {
		options.keyDelimiter = keyDelimiter
	}
}

func WithWatchHandlers(watchHandlers []WatchHandler) Opt {
	return func(options *Options) {
		options.watchHandlers = watchHandlers
	}
}

func defaultOptions() *Options {
	opts := &Options{
		path:          defaultPath,
		keyDelimiter:  defaultKeyDelimiter,
		watchHandlers: make([]WatchHandler, 0),
	}
	return opts
}
