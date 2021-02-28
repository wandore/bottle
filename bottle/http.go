package bottle

import (
	"bottle/hash"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	PrefixPath = "/bottle/"
	Backup     = 10
)

type HttpPool struct {
	addr    string
	prefix  string
	mu      sync.RWMutex
	conHash *hash.ConHash
	httpMap map[string]*httpGetter
}

func NewHttpPool(addr string) *HttpPool {
	return &HttpPool{
		addr:   addr,
		prefix: PrefixPath,
	}
}

func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.prefix) {
		log.Fatal("unexpected prefix: " + r.URL.Path)
	}

	log.Println("http method: " + r.Method + " prefix: " + r.URL.Path)

	adds := strings.SplitN(r.URL.Path[len(p.prefix):], "/", 2)
	if len(adds) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	bottleName := adds[0]
	bottle := GetBottle(bottleName)
	if bottle == nil {
		http.Error(w, "no bottle: "+bottleName, http.StatusNotFound)
		return
	}

	key := adds[1]
	value, err := bottle.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(value.Clone())
}

func (p *HttpPool) Set(nodes ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.conHash = hash.New(Backup, nil)
	p.conHash.Add(nodes...)

	p.httpMap = make(map[string]*httpGetter, len(nodes))

	for _, node := range nodes {
		p.httpMap[node] = &httpGetter{addr: node + PrefixPath}
	}
}

func (p *HttpPool) NodeRoute(key string) (NodeGetter, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if node := p.conHash.Get(key); node != "" && node != p.addr {
		log.Println("Choose node: " + node)
		return p.httpMap[node], true
	}

	return nil, false
}
