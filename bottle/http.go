package bottle

import (
	"log"
	"net/http"
	"strings"
)

const PrefixPath = "/bottle/"

type HttpPool struct {
	addr string
	prefix string
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
		http.Error(w, "no bottle: " + bottleName, http.StatusNotFound)
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