package tapbag

import (
	"container/list"
	"sync"
)

// least recently used cache
type LruCache struct {
	Lock sync.RWMutex
	MaxLength int
	Map     map[string]*list.Element
	Cache list.List
	OnEvicted func(key string, value Value)
}

type Value interface {
	Len()
}

type entry struct {
	key string
	value Value
}

func NewLruCache(maxLength int, OnEvicted func(string, Value)) *LruCache {
	return &LruCache{
		MaxLength: maxLength,
		OnEvicted: OnEvicted,
	}
}


func (lc *LruCache) Get(key string) (Value, bool) {
	lc.Lock.Lock()
	defer lc.Lock.Unlock()
	if ele, ok := lc.Map[key]; ok {
		lc.Cache.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}


func (lc *LruCache) RemoveOldest() {
	lc.Lock.Lock()
	defer lc.Lock.Unlock()
	ele := lc.Cache.Back()
	if ele != nil {
		lc.Cache.Remove(ele)
		kv := ele.Value.(*entry)
		delete(lc.Map, kv.key)
		if lc.OnEvicted != nil {
			lc.OnEvicted(kv.key, kv.value)
		}
	}
}

func (lc *LruCache) Add(key string, value Value) {
	lc.Lock.Lock()
	defer lc.Lock.Unlock()
	if ele, ok := lc.Map[key]; ok {
		lc.Cache.MoveToFront(ele)
		kv := ele.Value.(*entry)
		kv.value = value
	} else {
		ele := lc.Cache.PushFront(&entry{key: key, value: value})
		lc.Map[key] = ele
	}
	for lc.Cache.Len() > lc.MaxLength {
		lc.RemoveOldest()
	}
}

func (lc *LruCache) Len() int {
	return lc.Cache.Len()
}