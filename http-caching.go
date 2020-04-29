package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

var newCache *cache.Cache

func init() {
	newCache = cache.New(5*time.Minute, 10*time.Minute)
	newCache.Set("foo", "bar", cache.DefaultExpiration)
}

func getFromCache(w http.ResponseWriter, r *http.Request) {
	foo, found := newCache.Get("foo")
	if found {
		log.Print("Key FOUND in cache, value as :: ", foo.(string))
		fmt.Fprintf(w, "Hello "+foo.(string))
	} else {
		log.Print("Key NOT found in cache :: ", "foo")
		fmt.Fprintf(w, "Key NOT found in cache...")
	}
}

func main() {
	http.HandleFunc("/", getFromCache)
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("Error starting http server: ", err)
		return
	}
}
