package conf

var conf = Init()

func InitConf(opts ...Opt) {
	conf = Init(opts...)
}

func Get(key string) interface{} {
	return conf.Get(key)
}

func GetInt(key string) int {
	return conf.Get(key).(int)
}
func GetString(key string) string {
	return conf.Get(key).(string)
}

func GetObject(key string, result interface{}) {
	marshalByte, _ := conf.encoder.Encode(conf.Get(key))
	conf.encoder.Decode(marshalByte, result)
}

func Set(key string, value interface{}) {
	conf.Set(key, value)
}

func Watch(eventChan chan *Event) {
	for e := range eventChan {
		conf.OnChange(e)
	}
}
