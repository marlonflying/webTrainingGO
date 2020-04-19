package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the server response")
}

func main() {
	http.HandleFunc("/", home)
	err := http.ListenAndServe(CONN_HOST+":"+CONN_PORT, nil)
	if err != nil {
		log.Fatal("error satarting http server:", err)
		return
	}
}
