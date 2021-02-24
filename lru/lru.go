package lru

import (
	"container/list"
	"log"
)

type Cache struct {
	cap int
	used int
	cacheList *list.List
	cacheMap map[string]*list.Element
	CallBack func(key string, value valueType)
}

type entry struct {
	key string
	value valueType
}

type valueType interface {
	Len() int
}

func (c *Cache) Get(key string) (value valueType, exist bool) {
	if elem, exist := c.cacheMap[key]; exist {
		c.cacheList.MoveToBack(elem)
		kv := elem.Value.(*entry)
		return kv.value, exist
	}
	return
}

func (c *Cache) Remove() {
	elem := c.cacheList.Front()
	if elem != nil {
		c.cacheList.Remove(elem)
		kv := elem.Value.(*entry)
		delete(c.cacheMap, kv.key)
		c.used -= len(kv.key) + kv.value.Len()
		if c.CallBack != nil {
			c.CallBack(kv.key, kv.value)
		}
	}
	return
}

func (c *Cache) Add(key string, value valueType) {
	if elem, exist := c.cacheMap[key]; exist {
		kv := elem.Value.(*entry)
		extra := value.Len() - kv.value.Len()
		if c.cap - c.used < extra {
			return
		} else {
			c.cacheList.MoveToBack(elem)
			kv.value = value
			c.used += extra
			log.Println("key: " + key + " recorded")
		}
	} else {
		need := len(key) + value.Len()
		if need > c.cap {
			return
		}
		for c.cap - c.used < need {
			c.Remove()
		}
		c.cacheList.PushBack(&entry{
			key:   key,
			value: value,
		})
		c.used += need
	}
}

func New(cap int, CallBack func(key string, value valueType)) *Cache {
	return &Cache{
		cap:       cap,
		used:      0,
		cacheList: list.New(),
		cacheMap:  make(map[string]*list.Element, 0),
		CallBack: CallBack,
	}
}
