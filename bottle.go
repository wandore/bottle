package bottle

import (
	"fmt"
	"log"
	"sync"
)

type Bottle struct {
	name string
	getter Getter
	bottle cache
}

var (
	mu sync.RWMutex
	bottleMap = make(map[string]*Bottle, 0)
)

func NewBottle(name string, getter Getter, cap int) *Bottle {
	if getter == nil {
		log.Fatal("No getter func!")
	}

	mu.Lock()
	mu.RUnlock()

	g := &Bottle{
		name:   name,
		getter: getter,
		bottle: cache{cap: cap},
	}
	bottleMap[name] = g
	return g
}

func GetBottle(name string) *Bottle {
	mu.RLock()
	defer mu.RUnlock()

	g := bottleMap[name]
	return g
}

func (g *Bottle) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is empty")
	}

	if v, exist := g.bottle.get(key); exist {
		log.Println(key + " hit")
		return v, nil
	}

	return g.getLocally(key)
}

func (g *Bottle) getLocally(key string) (ByteView, error) {
	b, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	s := ByteView{b: cloneBytes(b)}
	g.bottle.add(key, s)
	return s, nil
}


