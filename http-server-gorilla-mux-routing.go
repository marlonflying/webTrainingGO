package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

//GetRequestHandler Home message
var GetRequestHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("From GetRequestHandler function..."))
	})

//PostRequestHandler Post request
var PostRequestHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It's a post Request..."))
	})

//PathVariableHandler Variable URL
var PathVariableHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		w.Write([]byte("Hi " + name))
	})

func main() {
	router := mux.NewRouter()
	router.Handle("/", GetRequestHandler).Methods("GET")
	router.Handle("/post", PostRequestHandler).Methods("POST")
	router.Handle("/hello/{name}", PathVariableHandler).Methods("GET", "PUT")
	http.ListenAndServe(connHost+":"+connPort, router)
}
