package main

import (
	"bottle/bottle"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	bottle.NewBottle("scores", 2<<10, bottle.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[LoaclDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("value: %s not exist", key)
		}))

	addr := "localhost:8080"
	peers := bottle.NewHttpPool(addr)
	log.Println("bottle is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}