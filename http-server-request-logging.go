package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

//GetRequestHandler Home message
var GetRequestHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Home message..."))
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
	logFile, err := os.OpenFile("server.log",
		os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error opnening log file : ", err)
		return
	}
	router := mux.NewRouter()
	router.Handle("/", handlers.LoggingHandler(logFile,
		GetRequestHandler)).Methods("GET")
	router.Handle("/post", handlers.LoggingHandler(logFile,
		PostRequestHandler)).Methods("POST")
	router.Handle("/hello/{name}",
		handlers.CombinedLoggingHandler(logFile,
			PathVariableHandler)).Methods("GET", "PUT")
	http.ListenAndServe(connHost+":"+connPort, router)
}
