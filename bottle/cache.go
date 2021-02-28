package bottle

import (
	"bottle/lru"
	"sync"
)

type cache struct {
	mu  sync.RWMutex
	lru *lru.Cache
	cap int
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cap, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, exist bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.lru == nil {
		return
	}
	if value, exist := c.lru.Get(key); exist {
		return value.(ByteView), exist
	}
	return
}

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}


