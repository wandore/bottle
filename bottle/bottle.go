package bottle

import (
	"bottle/protect"
	"fmt"
	"log"
	"sync"
)

type Bottle struct {
	name    string
	getter  Getter
	bottle  cache
	nodes   NodeRouter
	handler *protect.Handler
}

var (
	mu        sync.RWMutex
	bottleMap = make(map[string]*Bottle, 0)
)

func NewBottle(name string, cap int, getter Getter) *Bottle {
	if getter == nil {
		log.Fatal("No getter func!")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Bottle{
		name:    name,
		getter:  getter,
		bottle:  cache{cap: cap},
		handler: &protect.Handler{},
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

	v, err := g.handler.Query(key, func() (interface{}, error) {
		if g.nodes != nil {
			if node, ok := g.nodes.NodeRoute(key); ok {
				if value, err := g.getFromNode(node, key); err == nil {
					return value, nil
				} else {
					log.Println("Failed to get from node: ", node)
				}
			}
		}
		return g.getLocally(key)
	})

	if err != nil {
		return ByteView{}, err
	}

	return v.(ByteView), nil
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

func (g *Bottle) Register(nodes NodeRouter) {
	if g.nodes != nil {
		log.Fatal("Nodes has been registered already.")
	}

	g.nodes = nodes
}

func (g *Bottle) getFromNode(node NodeGetter, key string) (ByteView, error) {
	bytes, err := node.NodeGet(g.name, key)
	if err != nil {
		return ByteView{}, err
	}

	return ByteView{b: bytes}, nil
}
