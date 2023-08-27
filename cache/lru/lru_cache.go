package nonolru

import (
	lru "github.com/hashicorp/golang-lru"
	"time"
)

type LruCache struct {
	Size  int
	cache *lru.Cache
}

type LruCacheItem struct {
	data      interface{}
	expireAt  time.Time
	perpetual bool
}

func New(size int) *LruCache {
	cache := &LruCache{}
	cache.cache, _ = lru.New(cache.Size)
	cache.Size = size
	return cache
}

// Put
// 存储数据
func (cache *LruCache) Put(key string, data interface{}, duration time.Duration) {

	var item = LruCacheItem{
		data:     data,
		expireAt: time.Now().Add(duration),
	}
	//存储
	if duration <= 0 {
		item.perpetual = true
	}

	cache.cache.Add(key, item)
}

// Get
// 获取数据
func (cache *LruCache) Get(key string) interface{} {
	if cache.cache == nil {
		return nil
	}
	var intf, ok = cache.cache.Get(key)
	// 未命中
	if !ok {
		return nil
	}
	if intf == nil {
		return nil
	}
	var item = intf.(LruCacheItem)

	// 命中但过期
	if item.expireAt.Sub(time.Now()) < 0 && !item.perpetual {
		cache.cache.Remove(key)
		return nil
	}

	// 命中且未过期
	return item.data
}

// Get
// 获取数据
func (cache *LruCache) Keys() (keys []string) {
	if cache.cache == nil {
		return nil
	}

	is := cache.cache.Keys()
	keys = make([]string, len(is))
	for i, v := range is {
		if key, ok := v.(string); ok {
			keys[i] = key
		}
	}

	return keys
}

func (cache *LruCache) Delete(key string) bool {
	if cache.cache == nil {
		return true
	}

	return cache.cache.Remove(key)
}
