package conf

var conf = Init()

func InitConf(opts ...Opt) {
	conf = Init(opts...)
}

func Get(key string) interface{} {
	return conf.Get(key)
}

func Set(key string, value interface{}) {
	conf.Set(key, value)
}

func Watch(eventChan chan *Event) {
	for e := range eventChan {
		conf.OnChange(e)
	}
}
