package nonolru

import "time"

const defaultSize = 1000

var defaultCache = New(defaultSize)

func Get(key string) interface{} {
	return defaultCache.Get(key)
}

func Keys() []string {
	return defaultCache.Keys()
}

func Put(key string, data interface{}, duration time.Duration) {
	defaultCache.Put(key, data, duration)
}

func Delete(key string) bool {
	return defaultCache.cache.Remove(key)
}
