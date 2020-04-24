package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

func helloworld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Careloco!")
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login Page!")
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Logout Page!")
}

func main() {
	http.HandleFunc("/", helloworld)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
}
