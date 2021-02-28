package main

import (
	"bottle/bottle"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"a": "0",
	"b": "1",
	"c": "2",
}

func createBottle() *bottle.Bottle {
	return bottle.NewBottle("table", 3, bottle.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, b *bottle.Bottle) {
	nodes := bottle.NewHttpPool(addr)
	nodes.Set(addrs...)
	b.Register(nodes)
	log.Println("bottle is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], nodes))
}

func startAPIServer(apiAddr string, b *bottle.Bottle) {
	http.Handle("/bottle", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := b.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.Clone())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	b := createBottle()
	if api {
		go startAPIServer(apiAddr, b)
	}
	startCacheServer(addrMap[port], []string(addrs), b)
}
