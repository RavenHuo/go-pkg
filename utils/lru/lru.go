package lru

import (
	"container/list"
	"sync"
)

type Lru struct {
	list   *list.List
	lruMap map[interface{}]*list.Element
	size   int
	rwLock *sync.RWMutex
}

func New(size int) *Lru {
	return &Lru{
		list:   list.New(),
		lruMap: make(map[interface{}]*list.Element),
		size:   size,
		rwLock: &sync.RWMutex{},
	}
}
func (lru *Lru) Get(Key interface{}) interface{} {
	lru.rwLock.RLock()
	defer lru.rwLock.RUnlock()
	val, ok := lru.lruMap[Key]
	if !ok {
		return nil
	}
	lru.list.MoveToFront(val)
	return val.Value
}
func (lru *Lru) Put(Key, value interface{}) {
	lru.rwLock.Lock()
	defer lru.rwLock.Unlock()
	val, ok := lru.lruMap[Key]
	if ok {
		val.Value = value
	} else {
		// 超过lru长度限制
		if lru.size < lru.list.Len()+1 {
			lru.list.Back()
		}
		val = &list.Element{Value: value}
	}
	lru.list.MoveToFront(val)
	lru.lruMap[Key] = val
}

func (lru *Lru) Exist(key interface{}) bool {
	lru.rwLock.RLock()
	defer lru.rwLock.RUnlock()
	_, ok := lru.lruMap[key]
	return ok
}
